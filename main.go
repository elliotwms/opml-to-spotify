package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gilliek/go-opml/opml"
	"github.com/google/uuid"
	"github.com/zmb3/spotify/v2"
	"github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

func main() {
	if len(os.Args) < 2 {
		panic("missing argument: filename")
	}

	clientID := os.Getenv("SPOTIFY_ID")
	if clientID == "" {
		panic("missing SPOTIFY_ID env var")
	}

	outlines, err := getOutlines(os.Args[1])
	if err != nil {
		panic(err)
	}

	client := spotify.New(login())

	fmt.Printf("Searching for %d shows", len(outlines))

	ctx := context.Background()
	shows, err := search(ctx, client, outlines)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found %d out of %d shows", len(shows), len(outlines))

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
func login() *http.Client {
	s := http.Server{
		Addr: "localhost:8080",
	}

	state := uuid.New().String()

	clientID := os.Getenv("SPOTIFY_ID")
	auth := spotifyauth.New(
		spotifyauth.WithClientID(clientID),
		spotifyauth.WithScopes(spotifyauth.ScopeUserLibraryModify),
		spotifyauth.WithRedirectURL("http://localhost:8080/callback"),
	)

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

	fmt.Printf("Click here to log in: %s\nWaiting for token...", auth.AuthURL(state))

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
