package cmd

import (
	"time"

	"dbconntest/controller"

	"github.com/spf13/cobra"
)

func DoWork(jobType string) {
	params := controller.JobParams{
		JobType:      jobType,
		DbType:       Driver,
		URL:          URL,
		Query:        SQL,
		Connections:  Conns,
		Timeout:      Timeout,
		ThreadLock:   ThreadLock,
		LogFormat:    LogFormat,
		LogLevel:     LogLevel,
		DriverTraces: DriverTraces,
		DriverLogs:   DriverLogs,
	}
	controller.DoWork(&params)
}

func setCommonFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&URL, "URL", "u", "", "ODBC connection string")
	cmd.PersistentFlags().StringVarP(&Driver, "driver", "d", "", "DB driver (go_ibm_db for DB2, godror for Oracle)")
	cmd.PersistentFlags().IntVarP(&Conns, "connections", "c", 1, "number of connections")
	cmd.PersistentFlags().DurationVarP(&Timeout, "timeout", "t", 10*time.Second, "timeout")
	cmd.PersistentFlags().BoolVarP(&ThreadLock, "threadlock", "l", false, "If set to \"true\", each connection locks an OS thread")
	cmd.PersistentFlags().StringVarP(&LogFormat, "logformat", "", "console", "log format (console or json)")
	cmd.PersistentFlags().StringVarP(&LogLevel, "loglevel", "", "info", "log level (debug, info, warning, error, fatal, panic)")
	cmd.PersistentFlags().BoolVarP(&DriverLogs, "driverlogs", "", false, "print driver logs (can be very verbose!!")
	cmd.PersistentFlags().BoolVarP(&DriverTraces, "drivertraces", "", false, "print ODBC calls traces (can be very verbose!!")

	_ = cmd.RegisterFlagCompletionFunc("driver", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"go_ibm_db\tDB2 Driver-", "json\tJson format"}, cobra.ShellCompDirectiveDefault
	})

	_ = cmd.RegisterFlagCompletionFunc("logformat", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"console\tHuman readable format", "json\tJson format"}, cobra.ShellCompDirectiveDefault
	})

	_ = cmd.RegisterFlagCompletionFunc("loglevel", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"debug", "info", "warning", "error", "fatal", "panic"}, cobra.ShellCompDirectiveDefault
	})

	_ = cmd.MarkPersistentFlagRequired("URL")
	_ = cmd.MarkPersistentFlagRequired("driver")

}
