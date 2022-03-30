/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"dbconntest/controller"

	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "ping",
	Short: "Just connect and ping the database",
	Long: `
With this command, each goroutine will make a connection to the database, which requires the creation of
a new physical connection (typically a TCP Socket, optionally encrypted) to the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		DoWork(controller.Ping)
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)

	setCommonFlags(connectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// connectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// connectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
