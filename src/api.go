package api

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	// "errors"
	"fmt"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
	"log"
)

type dataImage struct {
	imageUrl string
}

var data = dataImage{
	imageUrl: "",
}

var GameIndex int
var MyPlaylist *spotify.FullPlaylist
var track *spotify.FullTrack
var TheLyrics string

func Api() *spotify.FullTrack {
	authConfig := &clientcredentials.Config{
		ClientID:     "a8795237d8ea48a09bc79862cf40c00b",
		ClientSecret: "8847d2dc62fc4948b882a1069e16cb35",
		TokenURL:     spotify.TokenURL,
	}

	accessToken, err := authConfig.Token(context.Background())
	if err != nil {
		log.Fatalf("error retrieve access token: %v", err)
	}

	client := spotify.Authenticator{}.NewClient(accessToken)

	playlistID := spotify.ID("37i9dQZF1E37xpTgH5tN1Z?si=e6ee0165030640af")
	playlist, err := client.GetPlaylist(playlistID)
	MyPlaylist = playlist
	if err != nil {
		log.Fatalf("error retrieve playlist data: %v", err)
	}

	log.Println("playlist id:", playlist.ID)
	log.Println("playlist name:", playlist.Name)
	log.Println("playlist description:", playlist.Description)

	//log.Println(" le fameu playlist", playlist)
	//log.Println(" le fameu playlist", MyPlaylist)

	for {
		Max := len(playlist.Tracks.Tracks)
		randomIndex := GetRandomIndex(Max)

		GameIndex = randomIndex

		track = &playlist.Tracks.Tracks[randomIndex].Track

		artist := GetArtistsNames(track.Artists)
		title := track.Name

		lyrics, err := GetLyrics(artist, title)
		TheLyrics = lyrics

		if err != nil || lyrics == "" {
			log.Printf("No lyrics found for %s by %s, skipping...", title, artist)
			continue
		}
		log.Println("Track Name", track.Name)
		log.Println("Artists(s)", GetArtistsNames(track.Artists))

		artistID := track.Artists[0].ID
		artistDetails, err := client.GetArtist(artistID)
		if err != nil {
			log.Printf("Error retrieving artist details: %v", err)
		} else {
			log.Println("Artist Image URL:", artistDetails.Images[0].URL)
		}

		artistImageURL := ""
		if len(artistDetails.Images) > 0 {
			artistImageURL = artistDetails.Images[0].URL
		}

		data.imageUrl = artistImageURL

		break
	}

	return track
}

func GetRandomIndex(max int) int {
	var randomIndex int
	if max > 0 {
		buf := make([]byte, 8)
		_, err := rand.Read(buf)
		if err != nil {
			log.Fatalf("error generating random index: %v", err)
		}
		randomIndex = int(binary.BigEndian.Uint64(buf) % uint64(max))
	}
	return randomIndex
}

func GetArtistsNames(artists []spotify.SimpleArtist) string {
	var names []string
	for _, artist := range artists {
		names = append(names, artist.Name)
	}
	return fmt.Sprintf("%v", names)
}

func NextTrack(playlist *spotify.FullPlaylist) (*spotify.FullTrack, error) {
	fmt.Println(GameIndex)
	if playlist == nil || len(playlist.Tracks.Tracks) == 0 {
		return nil, errors.New("empty or nil playlist")
	}

	if TheLyrics != "" {
		GameIndex = (GameIndex + 1) % len(playlist.Tracks.Tracks)
		if GameIndex < 0 {
			GameIndex = 0
		}
		fmt.Println(GameIndex)
		fmt.Println(playlist.Tracks.Tracks)
		return &playlist.Tracks.Tracks[GameIndex].Track, nil
	}
	GameIndex = (GameIndex + 1) % len(playlist.Tracks.Tracks)
	if GameIndex < 0 {
		GameIndex = 0
	}
	fmt.Println(GameIndex)
	// fmt.Println(playlist.Tracks.Tracks)
	return &playlist.Tracks.Tracks[GameIndex].Track, nil
}

// func PreviousTrack(playlist *spotify.FullPlaylist) (*spotify.FullTrack, error) {
// 	GameIndex = (GameIndex - 1 + len(playlist.Tracks.Tracks)) % len(playlist.Tracks.Tracks)
// 	return &playlist.Tracks.Tracks[GameIndex].Track, nil
// }

// func RestartPlaylist(playlist *spotify.FullPlaylist) error {
// 	newIndex := GetRandomIndex(len(playlist.Tracks.Tracks))
// 	if newIndex < 0 || newIndex >= len(playlist.Tracks.Tracks) {
// 		return errors.New("invalid index")
// 	}
// 	GameIndex = newIndex
// 	return nil
// }