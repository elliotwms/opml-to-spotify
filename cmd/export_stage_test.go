package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/elliotwms/opml-to-spotify/internal/clients"
	"github.com/elliotwms/opml-to-spotify/internal/clients/itunes"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"github.com/zmb3/spotify/v2"
)

// ExportTest is a wrapper to clean up resources usually created by running the export command
func ExportTest(t *testing.T, f func()) {
	_ = os.Remove("spotify.opml")
	t.Cleanup(func() {
		_ = os.Remove("spotify.opml")
	})

	f()
}

func NewExportStage(t *testing.T) (*ExportStage, *ExportStage, *ExportStage) {
	out := new(bytes.Buffer)
	exportCmd.SetOut(out)

	s := &ExportStage{
		t:       t,
		require: require.New(t),
		cmd:     exportCmd,
		out:     out,
	}

	return s, s, s
}

type ExportStage struct {
	t       *testing.T
	require *require.Assertions
	out     *bytes.Buffer
	cmd     *cobra.Command
	opml    []byte
}

func (s *ExportStage) and() *ExportStage {
	return s
}
func (s *ExportStage) spotify_returns_one_show() *ExportStage {
	var server *httptest.Server
	clients.Spotify, server = testSpotifyClient(map[string]http.HandlerFunc{
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
	s.t.Cleanup(server.Close)

	return s
}

func (s *ExportStage) itunes_returns_a_match() {
	var server *httptest.Server
	clients.ITunes, server = testiTunesClient(map[string]http.HandlerFunc{
		"/search": func(w http.ResponseWriter, r *http.Request) {
			bs, _ := json.Marshal(itunes.Results{
				ResultsCount: 1,
				Results: []itunes.Entry{
					{
						TrackName: r.URL.Query().Get("term"),
						FeedURL:   "https://rss.hello.world",
					},
				},
			})

			_, _ = w.Write(bs)
		},
	})
	s.t.Cleanup(server.Close)
}

func (s *ExportStage) the_command_is_run() {
	// run the command
	s.cmd.Run(exportCmd, nil)
}

func (s *ExportStage) the_output_opml_file_exists() *ExportStage {
	var err error
	s.opml, err = os.ReadFile("spotify.opml")
	s.require.NoError(err, "Command should create spotify.opml file")

	return s
}

func (s *ExportStage) the_output_opml_file_contains_the_expected_show() {
	if len(s.opml) == 0 {
		panic("called before file read")
	}

	s.require.Contains(string(s.opml), "Hello, World!")
}
