package clients

import "github.com/elliotwms/opml-to-spotify/internal/clients/itunes"

var ITunes *itunes.Client

func GetiTunesClient() *itunes.Client {
	if ITunes != nil {
		return ITunes
	}

	return itunes.New()
}
