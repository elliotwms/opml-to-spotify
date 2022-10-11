package itunes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearch(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bs, _ := json.Marshal(Results{
			ResultsCount: 1,
			Results: []Entry{
				{
					TrackName: "James O'Brien's Mystery Hour",
					FeedURL:   "rss",
				},
			},
		})

		_, _ = w.Write(bs)
	}))

	client := &Client{
		HTTPClient: s.Client(),
		BaseURL:    s.URL,
	}

	res, err := client.Search("James O'Brien's Mystery Hour", "GB")

	require.NoError(t, err)
	require.Len(t, res, 1)
}
