package clients

import (
	"context"
	"github.com/elliotwms/opml-to-spotify/internal/config"
	"github.com/elliotwms/opml-to-spotify/pkg/pkce"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

var Spotify *spotify.Client

func GetSpotify(cmd *cobra.Command) *spotify.Client {
	if Spotify != nil {
		return Spotify
	}

	Spotify = login(cmd)

	return Spotify
}

// login logs the user into the application via OAuth authorization code flow
// https://developer.spotify.com/documentation/general/guides/authorization/code-flow/
// The user will be presented with a Spotify Login URL via the terminal which they should visit, then be redirected to a
// locally hosted http server which captures the auth code and performs token exchange
func login(cmd *cobra.Command) *spotify.Client {
	verifier := pkce.NewVerifier(pkce.LenMax)

	s := http.Server{
		Addr: "localhost:8080",
	}

	state := uuid.New().String()

	auth := spotifyauth.New(
		spotifyauth.WithClientID(config.ClientID(cmd)),
		spotifyauth.WithScopes(spotifyauth.ScopeUserLibraryRead, spotifyauth.ScopeUserLibraryModify),
		spotifyauth.WithRedirectURL("http://localhost:8080/callback"),
	)

	tokenChan := make(chan *oauth2.Token, 1)
	http.HandleFunc("/callback", func(writer http.ResponseWriter, request *http.Request) {
		token, err := auth.Token(context.Background(), state, request, verifier.Params()...)

		if err != nil {
			log.Fatalf(err.Error())
		}

		tokenChan <- token

		_, _ = writer.Write([]byte("Token retrieved successfully. Please return to your terminal"))
	})

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	challenge := verifier.Challenge()
	url := auth.AuthURL(state, challenge.Params()...)

	cmd.Println("Visit this URL in your browser to log in:", url)

	return spotify.New(auth.Client(context.Background(), <-tokenChan))
}
