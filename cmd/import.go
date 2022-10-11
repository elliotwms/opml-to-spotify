package cmd

import (
	"context"

	"github.com/elliotwms/opml-to-spotify/internal/clients"
	"github.com/gilliek/go-opml/opml"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify/v2"
)

const flagDryRun = "dry-run"

const maxSaveShowsBatchSize = 50

var (
	// importCmd represents the import command
	importCmd = &cobra.Command{
		Use:   "import filename.opml",
		Short: "Import your current podcast library into Spotify, using an OPML file",
		Long: `Import your current podcast library into Spotify, using an OPML file.

* Export your podcast library out of your old podcast app as an OPML file
* Create an application on Spotify: https://developer.spotify.com/dashboard/applications
* Run the command specifying your application's Client ID and Secret as either flags 
(--client-id and --client-secret) or environment variables (SPOTIFY_ID and SPOTIFY_SECRET)
`,
		Run: run,
	}
)

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().BoolP(flagDryRun, "d", false, "Read-only, will not update your subscriptions")
}

func run(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		panic("missing argument: filename")
	}

	outlines, err := getOutlines(args[0])
	if err != nil {
		panic(err)
	}

	client := clients.GetSpotify(cmd)

	cmd.Printf("Searching for %d shows\n", len(outlines))

	ctx := context.Background()
	shows, err := searchSpotifyForOutlines(cmd, ctx, client, outlines)
	if err != nil {
		panic(err)
	}

	cmd.Printf("Found %d out of %d shows\n", len(shows), len(outlines))

	if cmd.Flag(flagDryRun).Value.String() == "true" {
		cmd.Printf("Dry-run. Exiting...\n")
		return
	}

	// save shows for current user in batches of maxSaveShowsBatchSize
	for i := 0; i < len(shows); i += maxSaveShowsBatchSize {
		j := i + maxSaveShowsBatchSize
		if j > len(shows) {
			j = len(shows)
		}

		if err = client.SaveShowsForCurrentUser(ctx, shows[i:j]); err != nil {
			panic(err)
		}
	}
}

// getOutlines gets the outlines from a specified OPML file
func getOutlines(filename string) ([]opml.Outline, error) {
	f, err := opml.NewOPMLFromFile(filename)
	if err != nil {
		return nil, err
	}
	outlines := f.Outlines()
	return outlines, nil
}

// searchSpotifyForOutlines searches the Spotify API for each of the shows specified in the opml outlines by name,
// returning the first match of each
// It does not paginate through the search results as the exact match would typically be on the first page
func searchSpotifyForOutlines(cmd *cobra.Command, ctx context.Context, client *spotify.Client, outlines []opml.Outline) ([]spotify.ID, error) {
	var shows []spotify.ID

	for _, o := range outlines {
		cmd.Printf("Searching for %s\n", o.Title)
		res, err := client.Search(ctx, o.Title, spotify.SearchTypeShow)
		if err != nil {
			return nil, err
		}

		s := findShow(res.Shows.Shows, o)

		if s == nil {
			cmd.PrintErrf("Could not find show: %s\n", o.Title)
			continue
		}

		shows = append(shows, s.ID)
	}

	return shows, nil
}

// findShow finds the first show in the list
func findShow(shows []spotify.FullShow, o opml.Outline) *spotify.FullShow {
	for _, s := range shows {
		if s.Name == o.Title {
			return &s
		}
	}
	return nil
}
