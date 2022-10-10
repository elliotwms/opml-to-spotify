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

	FeedURL string `json:"feedURL"`
}

func Search(title, country string) ([]Entry, error) {
	//GET https://itunes.apple.com/search?term=Radiolab&media=podcast&entity=podcast&attribute=titleTerm
	u, err := url.Parse("https://itunes.apple.com/search?media=podcast&entity=podcast&attribute=titleTerm")
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
