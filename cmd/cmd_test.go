package cmd

import (
	"github.com/elliotwms/opml-to-spotify/internal/clients/itunes"
	"net/http"

	"github.com/zmb3/spotify/v2"
	"net/http/httptest"
)

func testSpotifyClient(handlers map[string]http.HandlerFunc) (*spotify.Client, *httptest.Server) {
	s := testClient(handlers)

	return spotify.New(s.Client(), spotify.WithBaseURL(s.URL+"/")), s
}

func testiTunesClient(handlers map[string]http.HandlerFunc) (*itunes.Client, *httptest.Server) {
	s := testClient(handlers)

	return &itunes.Client{
		HTTPClient: s.Client(),
		BaseURL:    s.URL,
	}, s
}

func testClient(handlers map[string]http.HandlerFunc) *httptest.Server {
	mux := http.NewServeMux()

	for path, handlerFunc := range handlers {
		mux.HandleFunc(path, handlerFunc)
	}

	s := httptest.NewServer(mux)
	return s
}
