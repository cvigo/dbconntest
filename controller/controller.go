package controller

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"

	"dbconntest/log"

	"github.com/godror/godror"
	"github.com/google/uuid"
	"github.com/ibmdb/go_ibm_db"
	statsLib "github.com/montanaflynn/stats"
)

const Ping = "ping"
const Txn = "txn"
const SimpleQuery = "query_alone"
const TxnQuery = "query_txn"

var db *sql.DB

type Stats struct {
	Start       time.Time
	End         time.Time
	TotalRuns   int
	TotalErrors int
	Ping        ResponseTimes
	BeginTx     ResponseTimes
	Query       ResponseTimes
	Commit      ResponseTimes
	Total       ResponseTimes
}

type ResponseTimes struct {
	Mean  time.Duration
	Pct95 time.Duration
	Pct99 time.Duration
}

type JobParams struct {
	JobType      string
	DbType       string
	URL          string
	Query        string
	Connections  int
	Timeout      time.Duration
	ThreadLock   bool
	LogFormat    string
	LogLevel     string
	DriverTraces bool
	DriverLogs   bool
}

type runStats struct {
	id          string
	start       time.Time
	end         time.Time
	pingTime    time.Duration
	beginTxTime time.Duration
	queryTime   time.Duration
	commitTime  time.Duration
	totalTime   time.Duration
	err         error
}

func DoWork(params *JobParams) {
	var err error
	var MaxGoRoutines int
	var MaxThreads int

	err = log.LogInit(params.LogLevel, params.LogFormat)
	if err != nil {
		fmt.Println("Error initializing log: ", err)
		os.Exit(1)
	}

	if params.DriverLogs {
		go_ibm_db.SetLogFunc(log.DriverLog)
		godror.SetLog(log.DriverLog)
	}

	if params.DriverTraces {
		go_ibm_db.SetTraceFunc(log.DriverTrace)
		go_ibm_db.SetTraceText("start_time", "end_time", "duration")
		godror.SetTrace(log.DriverTrace)
		godror.SetTraceText("start_time", "end_time", "duration")
	}

	var timings []*runStats

	var threadProfile = pprof.Lookup("threadcreate")
	var goroutineProfile = pprof.Lookup("goroutine")

	// make a buffered channel so the display goroutine can't slow down the workers
	completeCh := make(chan *runStats, params.Connections)
	doneCh := make(chan struct{})
	runs := 0

	db, err = sql.Open(params.DbType, params.URL)
	if err != nil {
		fmt.Printf("Error connecting to %s: %s\n", params.URL, err)
		return
	}

	db.SetMaxIdleConns(0) // Not interested in connection reuse for this type of test. Happy to drop after the operation ends
	db.SetMaxOpenConns(params.Connections * 2)

	// start the display goroutine
	go func() {
		for timing := range completeCh {
			timings = append(timings, timing)
			runs++
			if timing.err != nil {
				log.Logger.Errorf("%s Worker Error: %s", timing.id, timing.err)
			} else {
				log.Logger.With(
					"Ping", timing.pingTime,
					"BeginTx", timing.beginTxTime,
					"Query", timing.queryTime,
					"Commit", timing.commitTime,
					"start_time", timing.start,
					"end_time", timing.end,
					"duration", timing.totalTime,
				).Infof("%s completed in %04d position.", timing.id, runs)
			}

			if gr := runtime.NumGoroutine(); gr > MaxGoRoutines {
				MaxGoRoutines = gr
			}
			if thr := threadProfile.Count(); thr > MaxThreads {
				MaxThreads = thr
			}
			//threadProfile.WriteTo(os.Stderr, 1)
			//goroutineProfile.WriteTo(os.Stderr, 1)
		}
		doneCh <- struct{}{}
	}()

	// Do stuff
	ctx, cancel := context.WithTimeout(context.Background(), params.Timeout)

	// handle ctrl-c and kill signals
	go func() {
		sigchan := make(chan os.Signal)
		signal.Notify(sigchan, os.Interrupt)
		signal.Notify(sigchan, os.Kill)
		select {
		case <-sigchan:
			cancel() // this should cancel all the workers, unless they're stuck into a DB driver call
			log.Logger.Info("Program killed ! Waiting for workers to finish, 10 seconds to kill them...")
			time.AfterFunc(10*time.Second, func() {
				os.Exit(1)
			})
		}

	}()

	// spread the workers
	waitGroup := &sync.WaitGroup{}
	for i := 0; i < params.Connections; i++ {
		waitGroup.Add(1)
		go func(i int) {
			if params.ThreadLock {
				runtime.LockOSThread()
			}
			// Do the work and send the results to the display goroutine
			log.Logger.Debugf("threads in starting: %d", threadProfile.Count())
			log.Logger.Debugf("goroutines in starting: %d", goroutineProfile.Count())
			completeCh <- doWork(ctx, fmt.Sprintf("Worker-%04d", i), params)
			waitGroup.Done()
		}(i)
	}
	waitGroup.Wait()
	close(completeCh)
	<-doneCh
	_ = db.Close()

	// calculate stats
	stats := GetStats(timings)

	log.Logger.Infof("Total runs %d", stats.TotalRuns)
	log.Logger.Infof("Total errors %d", stats.TotalErrors)
	if params.JobType == Ping {
		log.Logger.Infow("PING Command",
			"Mean", stats.Ping.Mean,
			"Pct95", stats.Ping.Pct95,
			"Pct99", stats.Ping.Pct99,
		)
	}
	if params.JobType == Txn || params.JobType == TxnQuery {
		log.Logger.Infow("BEGIN_TXN Command",
			"Mean", stats.BeginTx.Mean,
			"Pct95", stats.BeginTx.Pct95,
			"Pct99", stats.BeginTx.Pct99,
		)
	}

	if params.JobType == SimpleQuery || params.JobType == TxnQuery {
		log.Logger.Infow("QUERY Command",
			"Mean", stats.Query.Mean,
			"Pct95", stats.Query.Pct95,
			"Pct99", stats.Query.Pct99,
		)
	}
	if params.JobType == Txn || params.JobType == TxnQuery {
		log.Logger.Infow("COMMIT Command",
			"Mean", stats.Commit.Mean,
			"Pct95", stats.Commit.Pct95,
			"Pct99", stats.Commit.Pct99,
		)
	}

	log.Logger.Infow("Total Command",
		"Mean", stats.Total.Mean,
		"Pct95", stats.Total.Pct95,
		"Pct99", stats.Total.Pct99,
	)

	log.Logger.Infof("Max. threads count: %d", MaxThreads)
	log.Logger.Infof("Max. goroutines count: %d", MaxGoRoutines)
}

func doWork(ctx context.Context, id string, params *JobParams) *runStats {
	start := time.Now()
	stats := &runStats{id: id, start: start}

	defer func() {
		stats.totalTime = stats.pingTime + stats.beginTxTime + stats.queryTime + stats.commitTime
		stats.end = time.Now()
	}()

	ctx = godror.ContextWithTraceTag(ctx, godror.TraceTag{
		ClientIdentifier: id,
		ClientInfo:       "dbconntest",
		DbOp:             params.JobType,
		Module:           "doWork",
		Action:           "read",
	})

	hostname, _ := os.Hostname()
	ctx = go_ibm_db.ContextWithTraceTag(ctx, go_ibm_db.TraceTag{
		ClientApplicationName:  "dbconntest",
		ClientCorrelationToken: uuid.NewString(),
		ClientUserID:           id,
		ClientWorkstationName:  hostname,
		ProgramID:              "doWork",
	})

	switch params.JobType {
	case Ping:
		err := db.PingContext(ctx)
		stats.pingTime = time.Since(start)
		stats.err = err
		return stats

	case Txn:
		tx, err := db.BeginTx(ctx, &sql.TxOptions{
			Isolation: sql.LevelDefault,
			ReadOnly:  false,
		})
		stats.beginTxTime = time.Since(start)
		if err != nil {
			stats.err = err
			return stats
		}
		start2 := time.Now()
		err = tx.Commit()
		stats.commitTime = time.Since(start2)
		stats.err = err
		return stats

	case SimpleQuery:
		rows, err := db.QueryContext(ctx, params.Query)
		stats.queryTime = time.Since(start)
		stats.err = err
		_ = rows.Close()
		return stats

	case TxnQuery:
		tx, err := db.BeginTx(ctx, &sql.TxOptions{
			Isolation: sql.LevelDefault,
			ReadOnly:  false,
		})
		stats.beginTxTime = time.Since(start)
		if err != nil {
			stats.err = err
			return stats
		}
		start2 := time.Now()
		rows, err := tx.QueryContext(ctx, params.Query)
		stats.queryTime = time.Since(start2)
		if err != nil {
			stats.err = err
			return stats
		}
		start3 := time.Now()
		err = tx.Commit()
		stats.commitTime = time.Since(start3)
		stats.err = err
		_ = rows.Close()
		return stats

	default:
		stats.err = fmt.Errorf("Unknown job type: %s", params.JobType)
		return stats
	}
}

func GetStats(runstats []*runStats) *Stats {
	var stats Stats

	var pingTime []time.Duration
	var beginTxTime []time.Duration
	var queryTime []time.Duration
	var commitTime []time.Duration
	var totalTime []time.Duration

	stats.TotalRuns = len(runstats)
	for _, s := range runstats {
		if s.err != nil {
			stats.TotalErrors++
		}
		pingTime = append(pingTime, s.pingTime)
		beginTxTime = append(beginTxTime, s.beginTxTime)
		queryTime = append(queryTime, s.queryTime)
		commitTime = append(commitTime, s.commitTime)
		totalTime = append(totalTime, s.totalTime)
	}

	stats.Ping.Mean, stats.Ping.Pct95, stats.Ping.Pct99 = calculate(pingTime)
	stats.BeginTx.Mean, stats.BeginTx.Pct95, stats.BeginTx.Pct99 = calculate(beginTxTime)
	stats.Query.Mean, stats.Query.Pct95, stats.Query.Pct99 = calculate(queryTime)
	stats.Commit.Mean, stats.Commit.Pct95, stats.Commit.Pct99 = calculate(commitTime)
	stats.Total.Mean, stats.Total.Pct95, stats.Total.Pct99 = calculate(totalTime)

	return &stats
}

func calculate(times []time.Duration) (mean, p95, p99 time.Duration) {
	data := statsLib.LoadRawData(times)
	t, err := statsLib.Mean(data)
	if err == nil {
		mean = time.Duration(t)
	}
	t, err = statsLib.Percentile(data, 95)
	if err == nil {
		p95 = time.Duration(t)
	}
	t, err = statsLib.Percentile(data, 99)
	if err == nil {
		p99 = time.Duration(t)
	}
	return
}
