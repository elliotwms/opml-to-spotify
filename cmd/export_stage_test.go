package cmd

import (
	"bytes"
	"encoding/json"
	"github.com/zmb3/spotify/v2"
	"net/http"
	"os"
	"testing"

	"github.com/elliotwms/opml-to-spotify/internal/clients/itunes"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
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
	out, err := new(bytes.Buffer), new(bytes.Buffer)
	exportCmd.SetOut(out)

	s := &ExportStage{
		t:             t,
		require:       require.New(t),
		cmd:           exportCmd,
		out:           out,
		errOut:        err,
		spotifyServer: setupMockSpotify(t),
		itunesServer:  setupMockItunes(t),
	}

	return s, s, s
}

type ExportStage struct {
	t       *testing.T
	require *require.Assertions
	out     *bytes.Buffer
	errOut  *bytes.Buffer
	cmd     *cobra.Command

	spotifyServer, itunesServer *http.ServeMux

	opml []byte
}

func (s *ExportStage) and() *ExportStage {
	return s
}

func (s *ExportStage) spotify_will_return_one_show() *ExportStage {
	s.spotifyServer.HandleFunc("/me/shows", func(w http.ResponseWriter, r *http.Request) {
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
	})

	return s
}

func (s *ExportStage) itunes_will_return_a_match() {
	s.itunesServer.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
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
	})
}

func (s *ExportStage) the_command_is_run() {
	// run the command
	s.cmd.Run(s.cmd, nil)

	s.t.Log(s.out.String())
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

func (s *ExportStage) no_errors_are_output() {
	s.require.Empty(s.errOut)
}
