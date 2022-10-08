package cmd

import (
	"context"
	"fmt"
	"github.com/gilliek/go-opml/opml"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify/v2"
	"github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

const flagClientID = "client-id"
const flagClientSecret = "client-secret"
const flagDryRun = "dry-run"

// importCmd represents the import command
var importCmd = &cobra.Command{
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

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringP(flagClientID, "c", "", "Spotify application Client ID")
	importCmd.Flags().StringP(flagClientSecret, "s", "", "Spotify application Client secret")
	importCmd.Flags().BoolP("dry-run", "d", false, "Read-only, will not update your subscriptions")
}

func run(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		panic("missing argument: filename")
	}

	outlines, err := getOutlines(args[0])
	if err != nil {
		panic(err)
	}

	client := spotify.New(login(
		cmd,
		cmd.Flag(flagClientID).Value.String(),
		cmd.Flag(flagClientSecret).Value.String(),
	))

	cmd.Printf("Searching for %d shows", len(outlines))

	ctx := context.Background()
	shows, err := search(ctx, client, outlines)
	if err != nil {
		panic(err)
	}

	cmd.Printf("Found %d out of %d shows\n", len(shows), len(outlines))

	if cmd.Flag(flagDryRun).Value.String() == "true" {
		cmd.Printf("Dry-run. Exiting...\n")
		return
	}

	if err = client.SaveShowsForCurrentUser(ctx, shows); err != nil {
		panic(err)
	}
}

func getOutlines(filename string) ([]opml.Outline, error) {
	f, err := opml.NewOPMLFromFile(filename)
	if err != nil {
		return nil, err
	}
	outlines := f.Outlines()
	return outlines, nil
}

// login
func login(cmd *cobra.Command, clientID, secret string) *http.Client {
	s := http.Server{
		Addr: "localhost:8080",
	}

	state := uuid.New().String()

	opts := []spotifyauth.AuthenticatorOption{
		spotifyauth.WithScopes(spotifyauth.ScopeUserLibraryModify),
		spotifyauth.WithRedirectURL("http://localhost:8080/callback"),
	}

	if clientID != "" {
		opts = append(opts, spotifyauth.WithClientID(clientID))
	}
	if secret != "" {
		opts = append(opts, spotifyauth.WithClientSecret(secret))
	}

	auth := spotifyauth.New(opts...)

	tokenChan := make(chan *oauth2.Token, 1)
	http.HandleFunc("/callback", func(writer http.ResponseWriter, request *http.Request) {
		token, err := auth.Token(context.Background(), state, request)
		if err != nil {
			log.Fatalf(err.Error())
		}

		tokenChan <- token
	})

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	cmd.Printf("Click here to log in: %s\nWaiting for token...", auth.AuthURL(state))

	return auth.Client(context.Background(), <-tokenChan)
}

// search the Spotify API for each of the shows specified in the opml outlines by name, returning the first match of
// each
func search(ctx context.Context, client *spotify.Client, outlines []opml.Outline) ([]spotify.ID, error) {
	var shows []spotify.ID

	for _, o := range outlines {
		fmt.Printf("Searching for %s\n", o.Title)
		res, err := client.Search(ctx, o.Title, spotify.SearchTypeShow)
		if err != nil {
			return nil, err
		}

		s := findShow(res.Shows.Shows, o)

		if s == nil {
			fmt.Printf("Could not find show: %s\n", o.Title)
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
