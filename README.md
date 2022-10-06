# opml-to-spotify

## Usage

* Export your podcast library out of your old podcast app as an `.opml` file
* [Create an application on Spotify](https://developer.spotify.com/dashboard/applications)
  * Obtain the Client ID
* Build and run the application
```shell
$ go build .
$ SPOTIFY_ID={id} ./opml-to-spotify {filename}.opml
```