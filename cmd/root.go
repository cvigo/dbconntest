/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"time"

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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dbconntest",
	Short: "Golang DB Connection Tester",
	Long:  `Golang DB Connection Tester`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
	rootCmd.PersistentFlags().StringVarP(&URL, "URL", "u", "", "ODBC connection string")
	rootCmd.PersistentFlags().StringVarP(&Driver, "driver", "d", "", "DB driver (go_ibm_db, godror")
	rootCmd.PersistentFlags().IntVarP(&Conns, "connections", "c", 1, "number of connections (default 1)")
	rootCmd.PersistentFlags().DurationVarP(&Timeout, "timeout", "t", 10*time.Second, "timeout (default 10s)")
	rootCmd.PersistentFlags().BoolVarP(&ThreadLock, "threadlock", "l", false, "each connection locks an OS thread (default false)")
	rootCmd.PersistentFlags().StringVarP(&LogFormat, "logformat", "", "console", "log format (console or json). Default console")
	rootCmd.PersistentFlags().StringVarP(&LogLevel, "loglevel", "", "info", "log level (debug, info, warning,, error, fatal, panic. Default info)")

	rootCmd.RegisterFlagCompletionFunc("driver", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"go_ibm_db\tDB2 Driver-", "json\tJson format"}, cobra.ShellCompDirectiveDefault
	})

	rootCmd.RegisterFlagCompletionFunc("logformat", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"console\tHuman readable format", "json\tJson format"}, cobra.ShellCompDirectiveDefault
	})

	rootCmd.RegisterFlagCompletionFunc("loglevel", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"debug", "info", "warning", "error", "fatal", "panic"}, cobra.ShellCompDirectiveDefault
	})

	_ = rootCmd.MarkPersistentFlagRequired("URL")
	_ = rootCmd.MarkPersistentFlagRequired("driver")
}
