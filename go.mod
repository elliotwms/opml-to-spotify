module github.com/elliotwms/opml-to-spotify

go 1.19

replace github.com/zmb3/spotify/v2 v2.3.0 => github.com/elliotwms/spotify/v2 v2.0.0-20221006214212-9854b6945e10

require (
	github.com/gilliek/go-opml v1.0.0
	github.com/google/uuid v1.3.0
	github.com/spf13/cobra v1.5.0
	github.com/zmb3/spotify/v2 v2.3.0
	golang.org/x/oauth2 v0.0.0-20220524215830-622c5d57e401
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)
