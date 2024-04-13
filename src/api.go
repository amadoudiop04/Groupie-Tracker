package api

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
	"log"
)

var (
	playlist     *spotify.FullPlaylist
	currentIndex int
)

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

	playlistID := spotify.ID("37i9dQZF1DWT0W2NVqAFjU?si=0363f4727b854b2b")
	playlist, err := client.GetPlaylist(playlistID)
	if err != nil {
		log.Fatalf("error retrieve playlist data: %v", err)
	}

	log.Println("playlist id:", playlist.ID)
	log.Println("playlist name:", playlist.Name)
	log.Println("playlist description:", playlist.Description)

	var track *spotify.FullTrack

	for {
		max := len(playlist.Tracks.Tracks)
		randomIndex := getRandomIndex(max)

		track = &playlist.Tracks.Tracks[randomIndex].Track

		artist := getArtistsNames(track.Artists)
		title := track.Name
		lyrics, err := GetLyrics(artist, title)

		if err != nil || lyrics == "" {
			log.Printf("No lyrics found for %s by %s, skipping...", title, artist)
			continue
		}
		log.Println("Track Name", track.Name)
		log.Println("Artists(s)", getArtistsNames(track.Artists))
		break
	}
	return track
}

func nextTrack() {
	currentIndex = (currentIndex + 1) % len(playlist.Tracks.Tracks)
}

func previousTrack() {
	currentIndex = (currentIndex - 1 + len(playlist.Tracks.Tracks)) % len(playlist.Tracks.Tracks)
}

func restartPlaylist() {
	currentIndex = getRandomIndex(len(playlist.Tracks.Tracks))
}

func getRandomIndex(max int) int {
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

func getArtistsNames(artists []spotify.SimpleArtist) string {
	var names []string
	for _, artist := range artists {
		names = append(names, artist.Name)
	}
	return fmt.Sprintf("%v", names)
}
