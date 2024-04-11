package main

import (
	"context"
	"log"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

func main() {
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
}


// type PlaylistInfo struct {
// 	ID          string
// 	Name        string
// 	Description string
// 	ImageURL    string
// }

// var playlistInfo PlaylistInfo

// 	playlistInfo.ID = string(playlist.ID)
// 	playlistInfo.Name = playlist.Name
// 	playlistInfo.Description = playlist.Description

// 	// Extract artist image URL
// 	if len(playlist.Tracks.Tracks) > 0 && len(playlist.Tracks.Tracks[0].Track.Artists) > 0 {
// 		artist := playlist.Tracks.Tracks[0].Track.Artists[0]
// 		if len(artist.Images) > 0 {
// 			playlistInfo.ImageURL = artist.Images[0].URL
// 		}
// 	}

// 	// Print playlist information
// 	log.Println("playlist id:", playlistInfo.ID)
// 	log.Println("playlist name:", playlistInfo.Name)
// 	log.Println("playlist description:", playlistInfo.Description)
// 	log.Println("playlist image URL:", playlistInfo.ImageURL)
