/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"dbconntest/controller"

	"github.com/spf13/cobra"
)

// queryTxnCmd represents the queryTxn command
var queryTxnCmd = &cobra.Command{
	Use:   "query_txn",
	Short: "Runs a query inside a transaction",
	Long:  `Runs a query inside a transaction`,
	Run: func(cmd *cobra.Command, args []string) {
		DoWork(controller.TxnQuery)
	},
}

func init() {
	rootCmd.AddCommand(queryTxnCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// queryTxnCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	queryTxnCmd.Flags().StringVarP(&SQL, "sql", "", "", "SQL query to run")
	_ = queryTxnCmd.MarkFlagRequired("sql")
}
