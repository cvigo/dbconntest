/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"dbconntest/controller"

	"github.com/spf13/cobra"
)

// txnCmd represents the txn command
var txnCmd = &cobra.Command{
	Use:   "txn",
	Short: "Connect to the database, start and commit a DB transaction",
	Long: `
With this command, each goroutine creates a DB transaction, and commits the transaction without running any query.
As all the goroutines are launched at the same time, each one will make a new connection to the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		DoWork(controller.Txn)
	},
}

func init() {
	rootCmd.AddCommand(txnCmd)

	setCommonFlags(txnCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// txnCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// txnCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
