/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"dbconntest/log"

	"dbconntest/cmd"
)

func main() {
	cmd.Execute()

	if log.BaseLogger != nil {
		log.BaseLogger.Sync()
	}
}
