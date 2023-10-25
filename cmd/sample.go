package cmd

import (
	"github.com/noisyboy-9/golang_api_template/internal/app"
	"github.com/spf13/cobra"
)

// sampleCmd represents the sample command
var sampleCmd = &cobra.Command{
	Use: "sample",
	Run: runSampleCmd,
}

func init() {
	rootCmd.AddCommand(sampleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sampleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sampleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func runSampleCmd(cmd *cobra.Command, args []string) {
	app.InitApp()
	// main logic goes here
	app.SetupGracefulShutdown()
}
