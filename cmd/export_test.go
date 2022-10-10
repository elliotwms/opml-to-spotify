package cmd

import (
	"bytes"
	"encoding/json"
	"github.com/elliotwms/opml-to-spotify/internal/clients"
	"github.com/elliotwms/opml-to-spotify/internal/clients/itunes"
	"github.com/zmb3/spotify/v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExport(t *testing.T) {
	buf := new(bytes.Buffer)
	exportCmd.SetOut(buf)

	var s *httptest.Server
	clients.Spotify, s = testSpotifyClient(map[string]http.HandlerFunc{
		"/me/shows": func(w http.ResponseWriter, r *http.Request) {
			bs, _ := json.Marshal(spotify.SavedShowPage{Shows: []spotify.SavedShow{
				{
					FullShow: spotify.FullShow{
						SimpleShow: spotify.SimpleShow{
							Name:             "Hello, World!",
							AvailableMarkets: []string{"GB"},
						},
						Episodes: spotify.SimpleEpisodePage{},
					},
				},
			}})

			_, _ = w.Write(bs)
		},
	})
	defer s.Close()

	var s2 *httptest.Server
	clients.ITunes, s2 = testiTunesClient(map[string]http.HandlerFunc{
		"/search": func(w http.ResponseWriter, r *http.Request) {
			bs, _ := json.Marshal(itunes.Results{
				ResultsCount: 1,
				Results: []itunes.Entry{
					{
						TrackName: "Hello, World!",
						FeedURL:   "https://rss.hello.world",
					},
				},
			})

			_, _ = w.Write(bs)
		},
	})
	defer s2.Close()

	exportCmd.Run(exportCmd, nil)
}
