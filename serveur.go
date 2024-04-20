package main

import (
	"context"
	"fmt"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
	"html/template"
	"log"
	"net/http"
)

type gameData struct {
	Name       string
	ScoreBoard int
}

var Infos = gameData{
	Name:       "",
	ScoreBoard: 0,
}

type PageData struct {
	Track *spotify.SimpleTrack
}

var currentTrackIndex int

func api() []*spotify.SimpleTrack {
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

	playlistID := spotify.ID("6Xf0gjt1YmwvEG5iS8QOfg?si=2de553d01ff84abb")
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

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			input := r.FormValue("value")
			currentTrack := api()[currentTrackIndex]

			if input == currentTrack.Name {
			
				currentTrackIndex++
				if currentTrackIndex >= len(api()) {
					currentTrackIndex = 0
				}
			}
			fmt.Println(input)
		}

		for {
			track := api()[currentTrackIndex]

			if track.PreviewURL != "" {
				break
			}
			currentTrackIndex++
			if currentTrackIndex >= len(api()) {
				currentTrackIndex = 0
			}
		}

		tpl := template.Must(template.ParseFiles("index.html"))
		data := PageData{
			Track: api()[currentTrackIndex],
		}

		if err := tpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
