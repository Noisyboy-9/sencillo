/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/noisyboy-9/random-k8s-scheduler/internal/app"
	"github.com/spf13/cobra"
)

// scheduleCmd represents the schedule command
var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "start scheduling pods",
	Run:   scheduleRunner,
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
}

func scheduleRunner(cmd *cobra.Command, args []string) {
	app.InitApp()
	app.SetupGracefulShutdown()
}
