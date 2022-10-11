package config

import (
	"os"

	"github.com/spf13/cobra"
)

// clientID is the opml-to-spotify Spotify application Client ID
var clientID string

// ClientID gets the clientID var depending on priority from:
// * initial value (empty but can be set by ldflags)
// * SPOTIFY_ID environment variable
// * client-id flag
func ClientID(cmd *cobra.Command) string {
	if s := os.Getenv("SPOTIFY_ID"); s != "" {
		clientID = s
	}

	// allow overriding of client ID via flag
	if s := cmd.Flag("client-id").Value.String(); s != "" {
		clientID = s
	}

	return clientID
}
