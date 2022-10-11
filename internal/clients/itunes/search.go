// Package itunes is a really small package which sort of does one thing: search for podcasts in the iTunes Search API
package itunes

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Results struct {
	ResultsCount int     `json:"resultsCount"`
	Results      []Entry `json:"results"`
}

type Entry struct {
	TrackName string `json:"trackName"`
	FeedURL   string `json:"feedURL"`
}

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
}

func New() *Client {
	return &Client{
		HTTPClient: http.DefaultClient,
		BaseURL:    "https://itunes.apple.com",
	}
}

// Search searches the iTunes Search API for podcasts matching the title and country specified in the params
func (c *Client) Search(title, country string) ([]Entry, error) {
	u, err := url.Parse(c.BaseURL + "/search?media=podcast&entity=podcast&attribute=titleTerm")
	if err != nil {
		return nil, err
	}

	values := u.Query()
	values.Set("term", title)
	values.Set("country", country)

	u.RawQuery = values.Encode()

	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	results := Results{}
	err = json.NewDecoder(res.Body).Decode(&results)
	if err != nil {
		return nil, err
	}

	return results.Results, nil
}
