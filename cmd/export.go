package cmd

import (
	"context"
	"encoding/xml"
	"os"
	"time"

	"github.com/elliotwms/opml-to-spotify/internal/clients"
	"github.com/elliotwms/opml-to-spotify/internal/clients/itunes"
	"github.com/gilliek/go-opml/opml"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export your Spotify shows as an OPML file",
	Long:  `Export your Spotify shows as an OPML file`,
	Run:   runExport,
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringP("file", "f", "spotify.opml", "Name of the file to output generated OPML to")
}

func runExport(cmd *cobra.Command, _ []string) {
	client := clients.GetSpotify(cmd)

	// todo paginate
	shows, err := client.CurrentUsersShows(context.Background())
	if err != nil {
		panic(err)
	}

	cmd.Printf("Got %d shows", len(shows.Shows))

	var outlines []opml.Outline

	for _, show := range shows.Shows {
		results, err := itunes.Search(show.Name, show.AvailableMarkets[0])
		if err != nil {
			panic(err)
		}

		cmd.Printf("Found %d matches", len(results))

		for _, result := range results {
			// Exact matches on title only
			// Yes I know this is a terrible system
			if result.TrackName != show.Name {
				continue
			}

			outlines = append(outlines, opml.Outline{
				Type:    "rss",
				Text:    show.Name,
				Title:   show.Name,
				XMLURL:  results[0].FeedURL,
				HTMLURL: show.ExternalURLs["spotify"],
			})
		}

		break
	}

	// build OPML
	bs, err := xml.Marshal(buildOPML(outlines))
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(cmd.Flag("file").Value.String(), bs, 0644)
	if err != nil {
		panic(err)
	}
}

func buildOPML(outlines []opml.Outline) opml.OPML {
	return opml.OPML{
		Version: "1.0",
		Head: opml.Head{
			Title:       "Spotify Podcast Subscriptions",
			DateCreated: time.Now().Format(time.RFC3339),
			Docs:        "https://github.com/elliotwms/opml-to-spotify",
		},
		Body: opml.Body{
			Outlines: outlines,
		},
	}
}
