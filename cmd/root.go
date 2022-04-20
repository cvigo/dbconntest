/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"time"

	cc "github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"
)

var URL string
var Driver string
var Conns int
var Timeout time.Duration
var SQL string
var ThreadLock bool
var LogFormat string
var LogLevel string
var DriverTraces bool
var DriverLogs bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dbconntest",
	Short: "Golang DB Connection Tester",
	Long: `
This simple tool spreads the number of goroutines indicated by "--connections" and runs simple SQL commands
simultaneously.

Please see the list of available commands below.

Use "completion" command to generate the completion script for your shell (installation is left to you...).`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cc.Init(&cc.Config{
		RootCmd:  rootCmd,
		Headings: cc.HiBlue + cc.Bold + cc.Underline,
		Commands: cc.HiYellow + cc.Bold,
		Example:  cc.Italic,
		ExecName: cc.Bold,
		Flags:    cc.Bold,
	})
}
