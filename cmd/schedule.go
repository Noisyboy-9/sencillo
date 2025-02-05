package cmd

import (
	"github.com/spf13/cobra"

	"github.com/noisyboy-9/sencillo/internal/app"
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

func scheduleRunner(*cobra.Command, []string) {
	app.InitApp()
	app.SetupGracefulShutdown()
}
