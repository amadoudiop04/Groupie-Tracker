package games

import (
	"context"
	"log"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

func Api(ID string) []*spotify.SimpleTrack {
	authConfig := &clientcredentials.Config{
		ClientID:     "42d26f90ce1b486f96349f3f8f9cf94c",
		ClientSecret: "23166304a010453a9a31f5c93e625cd3",
		TokenURL:     spotify.TokenURL,
	}

	accessToken, err := authConfig.Token(context.Background())
	if err != nil {
		log.Fatalf("error retrieving access token: %v", err)
	}

	client := spotify.Authenticator{}.NewClient(accessToken)

	playlistID := spotify.ID(ID)
	playlist, err := client.GetPlaylist(playlistID)
	if err != nil {
		log.Fatalf("error retrieving playlist data: %v", err)
	}

	var tracks []*spotify.SimpleTrack
	for _, playlistTrack := range playlist.Tracks.Tracks {
		track := playlistTrack.Track

		simpleTrack := &spotify.SimpleTrack{
			ID:         track.ID,
			Name:       track.Name,
			Artists:    track.Artists,
			PreviewURL: track.PreviewURL,
		}
		tracks = append(tracks, simpleTrack)
	}
	return tracks
}
