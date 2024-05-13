package games

import (
	"context"
	"log"
	"math/rand"
	"strings"
	"unicode"

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
		log.Printf("error retrieving access token: %v", err)
	}

	client := spotify.Authenticator{}.NewClient(accessToken)

	playlistID := spotify.ID(ID)
	playlist, err := client.GetPlaylist(playlistID)
	if err != nil {
		log.Printf("error retrieving playlist data: %v", err)
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

func NextTrack(tracks []*spotify.SimpleTrack) (*spotify.SimpleTrack, int) {

	if len(tracks) == 0 {
		return nil, -1
	}

	var nextTrack *spotify.SimpleTrack
	var index int
	for {
		index = rand.Intn(len(tracks))
		nextTrack = tracks[index]

		if nextTrack.PreviewURL != "" {
			break
		}
	}

	return nextTrack, index
}

func RemovePlayedTracks(tracks []*spotify.SimpleTrack) []*spotify.SimpleTrack {
	updatedTracks := []*spotify.SimpleTrack{}

	for _, track := range tracks {
		if !IsTrackPlayed(track) {
			updatedTracks = append(updatedTracks, track)
		}
	}

	return updatedTracks
}

func IsTrackPlayed(track *spotify.SimpleTrack) bool {
	for _, playedTrack := range PlayedTracks {
		if track.ID == playedTrack.ID {
			return true
		}
	}
	return false
}

func NormalizeString(s string) string {
	s = strings.ToLower(s)

	t := make([]rune, 0, len(s))
	for _, r := range s {
		if unicode.IsLetter(r) {
			t = append(t, unicode.ToLower(r))
		}
	}
	return string(t)
}

type gameData struct {
	Name           string
	NumberOfTurn   int
	MusicDuration  int
	AnswerDuration int
}

var BlindtestData gameData

type PageData struct {
	Track *spotify.SimpleTrack
}

var PlayedTracks []*spotify.SimpleTrack
