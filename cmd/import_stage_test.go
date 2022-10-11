package cmd

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/gilliek/go-opml/opml"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"github.com/zmb3/spotify/v2"
)

func ImportTest(t *testing.T, f func()) {
	_ = os.Remove("test.opml")
	_ = os.Remove("missing.txt")
	t.Cleanup(func() {
		_ = os.Remove("test.opml")
		_ = os.Remove("missing.txt")
	})

	f()
}

func NewImportStage(t *testing.T) (*ImportStage, *ImportStage, *ImportStage) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)

	s := &ImportStage{
		t:       t,
		require: require.New(t),
		cmd:     importCmd,
		out:     out,
		errOut:  err,

		spotifyServer: setupMockSpotify(t),
		itunesServer:  setupMockItunes(t),
	}

	importCmd.SetOut(out)
	importCmd.SetErr(err)

	return s, s, s
}

type ImportStage struct {
	t       *testing.T
	require *require.Assertions
	out     *bytes.Buffer
	errOut  *bytes.Buffer
	cmd     *cobra.Command

	spotifyServer, itunesServer *http.ServeMux
	savedShows                  []string
}

func (s *ImportStage) and() *ImportStage {
	return s
}

func (s *ImportStage) an_opml_file() *ImportStage {
	bs, _ := xml.Marshal(opml.OPML{
		Version: "1.0",
		Head: opml.Head{
			Title: "Test Podcast Subscriptions",
		},
		Body: opml.Body{
			Outlines: []opml.Outline{
				{
					Type:    "rss",
					Text:    "Hello, World!",
					Title:   "Hello, World!",
					XMLURL:  "https://rss.hello.world",
					HTMLURL: "https://hello.world",
				},
			},
		},
	})

	_ = os.WriteFile("test.opml", bs, 0644)

	return s
}

func (s *ImportStage) spotify_will_return_search_results() *ImportStage {
	s.spotifyServer.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		bs, _ := json.Marshal(spotify.SearchResult{Shows: &spotify.SimpleShowPage{
			Shows: []spotify.FullShow{
				{
					SimpleShow: spotify.SimpleShow{
						ID:               "1",
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

func (s *ImportStage) spotify_will_return_no_search_results() *ImportStage {
	s.spotifyServer.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		bs, _ := json.Marshal(spotify.SearchResult{Shows: &spotify.SimpleShowPage{
			Shows: []spotify.FullShow{},
		}})

		_, _ = w.Write(bs)
	})

	return s
}

func (s *ImportStage) spotify_will_return_no_exact_matching_results() {
	s.spotifyServer.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		bs, _ := json.Marshal(spotify.SearchResult{Shows: &spotify.SimpleShowPage{
			Shows: []spotify.FullShow{
				{
					SimpleShow: spotify.SimpleShow{
						ID:               "1",
						Name:             "Hello, World! The Cheap Knock-Off",
						AvailableMarkets: []string{"GB"},
					},
					Episodes: spotify.SimpleEpisodePage{},
				},
			},
		}})

		_, _ = w.Write(bs)
	})
}

func (s *ImportStage) spotify_will_allow_the_user_to_save_the_shows() *ImportStage {
	s.spotifyServer.HandleFunc("/me/shows", func(w http.ResponseWriter, r *http.Request) {
		ids := strings.Split(r.URL.Query().Get("ids"), ",")

		s.savedShows = append(s.savedShows, ids...)
	})

	return s
}

func (s *ImportStage) the_dry_run_flag_is_set() {
	_ = s.cmd.Flags().Set(flagDryRun, "true")
}

func (s *ImportStage) the_missing_flag_is_set() {
	_ = s.cmd.Flags().Set(flagMissing, "missing.txt")
}

func (s *ImportStage) the_command_is_run() {
	s.cmd.Run(s.cmd, []string{"test.opml"})

	s.t.Log(s.out.String())
}

func (s *ImportStage) the_user_is_subscribed_to_the_show() *ImportStage {
	s.require.NotEmpty(s.savedShows)

	return s
}

func (s *ImportStage) the_user_is_not_subscribed_to_any_shows() *ImportStage {
	s.require.Empty(s.savedShows)

	return s
}

func (s *ImportStage) no_errors_are_output() *ImportStage {
	s.require.Empty(s.errOut)

	return s
}

func (s *ImportStage) the_error_is_output(v string) {
	s.require.Contains(s.errOut.String(), v)
}

func (s *ImportStage) the_message_is_output(v string) *ImportStage {
	s.require.Contains(s.out.String(), v)

	return s
}

func (s *ImportStage) the_missing_file_exists() {
	s.require.FileExists("missing.txt")
}
