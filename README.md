# opml-to-spotify

`opml-to-spotify` allows you to import your current podcast library into Spotify, with a couple of things to note:

* As you're unable to add a podcast to Spotify via RSS, `opml-to-spotify` searches the Spotify API for shows matching the title of your current list. This means you're likely to miss shows (it'll log when this happens)
* There's currently no way to export your shows from Spotify if you ever want to go back, so this is a one-way operation. The API doesn't expose enough information required to build an OPML file, so once you're in you're stuck.
* The only way to mark all your episodes as played is to go through each show _in the mobile app_, tap the gear on the Show page and _Mark as played_ from there

## Usage

> :warning: until `opml-to-spotify` is approved by Spotify you will need to follow the additional instructions for [using your own application](#Using-your-own-application) below 

* Export your podcast library out of your old podcast app as an OPML (`.opml`) file
* Download the binary for your OS from the latest release on the [releases](https://github.com/elliotwms/opml-to-spotify/releases) page
* Run the following:
```shell
$ ./opml-to-spotify import {filename}.opml
```

### Using your own application

During local development, and until the official `opml-to-spotify` application is approved, it will be necessary for you to use your own Spotify developer application to call the spotify API 

* [Create an application on Spotify](https://developer.spotify.com/dashboard/applications)
  * Obtain the Client ID
  * Add a redirect URI for `http://localhost:8080/callback`
* When running the application, specify the client ID ith either flags (`-c {id}`) or environment variables (`SPOTIFY_ID`)
