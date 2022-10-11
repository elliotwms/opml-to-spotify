package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const flagClientID = "client-id"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "opml-to-spotify",
	Short: "Import your current podcast library into Spotify",
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
	rootCmd.PersistentFlags().StringP(flagClientID, "c", "", "Spotify application Client ID")
}
