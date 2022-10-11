package cmd

import (
	"github.com/elliotwms/opml-to-spotify/internal/clients"
	"github.com/elliotwms/opml-to-spotify/internal/clients/itunes"
	"net/http"
	"testing"

	"github.com/zmb3/spotify/v2"
	"net/http/httptest"
)

func setupMockSpotify(t *testing.T) *http.ServeMux {
	c, s, mux := testSpotifyClient()
	clients.Spotify = c

	t.Cleanup(s.Close)

	return mux
}

func setupMockItunes(t *testing.T) *http.ServeMux {
	c, s, mux := testiTunesClient()
	clients.ITunes = c

	t.Cleanup(s.Close)

	return mux
}

func testSpotifyClient() (*spotify.Client, *httptest.Server, *http.ServeMux) {
	s, mux := testClient()

	return spotify.New(s.Client(), spotify.WithBaseURL(s.URL+"/")), s, mux
}

func testiTunesClient() (*itunes.Client, *httptest.Server, *http.ServeMux) {
	s, mux := testClient()

	return &itunes.Client{
		HTTPClient: s.Client(),
		BaseURL:    s.URL,
	}, s, mux
}

func testClient() (*httptest.Server, *http.ServeMux) {
	mux := http.NewServeMux()

	return httptest.NewServer(mux), mux
}
