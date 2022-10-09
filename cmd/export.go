package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export your Spotify shows as an OPML file",
	Long:  `Export your Spotify shows as an OPML file`,
	Run: func(cmd *cobra.Command, args []string) {
		client := login(cmd)

		// todo paginate
		shows, err := client.CurrentUsersShows(context.Background())
		if err != nil {
			panic(err)
		}

		cmd.Printf("Got %d shows", len(shows.Shows))
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
