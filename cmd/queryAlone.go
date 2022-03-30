/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"dbconntest/controller"

	"github.com/spf13/cobra"
)

// queryAloneCmd represents the queryAlone command
var queryAloneCmd = &cobra.Command{
	Use:   "query_alone",
	Short: "Runs a query without opening a transaction",
	Long: `
With this command, each goroutine runs the indicated SQL query to the database, without opening a transaction
beforehand.
As all the goroutines are launched at the same time, each one will make a new connection to the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		DoWork(controller.SimpleQuery)
	},
}

func init() {
	rootCmd.AddCommand(queryAloneCmd)

	setCommonFlags(queryAloneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// queryAloneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	queryAloneCmd.Flags().StringVarP(&SQL, "sql", "", "", "SQL query to run")
	_ = queryAloneCmd.MarkFlagRequired("sql")
}
