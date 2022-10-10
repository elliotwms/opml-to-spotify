package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/elliotwms/opml-to-spotify/internal/clients"
	"github.com/elliotwms/opml-to-spotify/internal/clients/itunes"
	"github.com/zmb3/spotify/v2"
)

func TestExport(t *testing.T) {
	_ = os.Remove("spotify.opml")
	t.Cleanup(func() {
		_ = os.Remove("spotify.opml")
	})

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
	t.Cleanup(s.Close)

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
	t.Cleanup(s2.Close)

	// run the command
	exportCmd.Run(exportCmd, nil)

	bs, err := os.ReadFile("spotify.opml")
	if err != nil {
		t.Fatal("did not create spotify.opml file")
	}

	if !strings.Contains(string(bs), "Hello, World!") {
		t.Fatal("Output OPML file did not contain expected show")
	}
}
