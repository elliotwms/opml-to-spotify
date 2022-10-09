package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/elliotwms/opml-to-spotify/pkg/pkce"
	"github.com/gilliek/go-opml/opml"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify/v2"
	"github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

const flagClientID = "client-id"
const flagDryRun = "dry-run"

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

	// clientID is the opml-to-spotify Spotify application Client ID, see setClientID
	clientID string
)

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringP(flagClientID, "c", "", "Spotify application Client ID")
	importCmd.Flags().BoolP(flagDryRun, "d", false, "Read-only, will not update your subscriptions")

}

func run(cmd *cobra.Command, args []string) {
	setClientID(cmd)

	if len(args) < 1 {
		panic("missing argument: filename")
	}

	outlines, err := getOutlines(args[0])
	if err != nil {
		panic(err)
	}

	client := spotify.New(login(cmd))

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

// setClientID sets the clientID var depending on priority from:
// * initial value (empty but can be set by ldflags)
// * SPOTIFY_ID environment variable
// * client-id flag
func setClientID(cmd *cobra.Command) {
	if s := os.Getenv("SPOTIFY_ID"); s != "" {
		clientID = s
	}

	// allow overriding of client ID via flag
	if s := cmd.Flag(flagClientID).Value.String(); s != "" {
		clientID = s
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

// login logs the user into the application via OAuth authorization code flow
// https://developer.spotify.com/documentation/general/guides/authorization/code-flow/
// The user will be presented with a Spotify Login URL via the terminal which they should visit, then be redirected to a
// locally hosted http server which captures the auth code and performs token exchange
func login(cmd *cobra.Command) *http.Client {
	verifier := pkce.NewVerifier(pkce.LenMax)

	s := http.Server{
		Addr: "localhost:8080",
	}

	state := uuid.New().String()

	auth := spotifyauth.New(
		spotifyauth.WithClientID(clientID),
		spotifyauth.WithScopes(spotifyauth.ScopeUserLibraryModify),
		spotifyauth.WithRedirectURL("http://localhost:8080/callback"),
	)

	tokenChan := make(chan *oauth2.Token, 1)
	http.HandleFunc("/callback", func(writer http.ResponseWriter, request *http.Request) {
		token, err := auth.Token(context.Background(), state, request, verifier.Params()...)

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

	challenge := verifier.Challenge()
	url := auth.AuthURL(state, challenge.Params()...)

	cmd.Printf("Visit this URL in your browser to log in: %s\n", url)

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
