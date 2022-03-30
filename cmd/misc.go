package cmd

import "dbconntest/controller"

func DoWork(jobType string) {
	params := controller.JobParams{
		JobType:     jobType,
		DbType:      Driver,
		URL:         URL,
		Query:       "",
		Connections: Conns,
		Timeout:     Timeout,
		ThreadLock:  ThreadLock,
		LogFormat:   LogFormat,
		LogLevel:    LogLevel,
	}
	controller.DoWork(&params)
}
