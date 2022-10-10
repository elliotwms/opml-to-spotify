package cmd

import (
	"context"
	"encoding/xml"
	"errors"
	"os"
	"time"

	"github.com/elliotwms/opml-to-spotify/internal/clients"
	"github.com/elliotwms/opml-to-spotify/internal/clients/itunes"
	"github.com/gilliek/go-opml/opml"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify/v2"
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
	var outlines []opml.Outline

	client := clients.GetSpotify(cmd)
	ctx := context.Background()
	shows, err := client.CurrentUsersShows(ctx)
	if err != nil {
		panic(err)
	}

	for {
		cmd.Printf("Got %d shows\n", len(shows.Shows))

		for _, show := range shows.Shows {
			cmd.Println("Searching for", show.Name)

			entry := searchiTunesForMatch(cmd, show)
			if entry == nil {
				cmd.Println("Could not match show:", show.Name)
				continue
			}

			outlines = append(outlines, buildOutline(show, entry))
		}

		cmd.Printf("Matched %d outlines\n", len(outlines))

		if err := client.NextPage(ctx, shows); err != nil {
			if errors.Is(err, spotify.ErrNoMorePages) {
				break
			} else {
				panic(err)
			}
		}
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

func searchiTunesForMatch(cmd *cobra.Command, show spotify.SavedShow) *itunes.Entry {
	results, err := itunes.Search(show.Name, show.AvailableMarkets[0])
	if err != nil {
		panic(err)
	}

	cmd.Printf("Found %d matches\n", len(results))

	for _, result := range results {
		// Exact matches on title only
		// Yes I know this is a terrible system - PRs welcome
		if result.TrackName == show.Name {
			return &result
		}
	}

	return nil
}

func buildOutline(show spotify.SavedShow, entry *itunes.Entry) opml.Outline {
	return opml.Outline{
		Type:    "rss",
		Text:    show.Name,
		Title:   show.Name,
		XMLURL:  entry.FeedURL,
		HTMLURL: show.ExternalURLs["spotify"],
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
