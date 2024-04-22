package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/zmb3/spotify"
)

type Song struct {
	Singer            string
	TitleSong         string
	LyricsSong        string
	ImageURL          string
	CheatMess         string
	Scores            int
	RemainingAttempts int
	Timer             int
	ThePlaylist       *spotify.FullPlaylist
}

type LyricsResponse struct {
	Lyrics string `json:"lyrics"`
}

var CurrentSong = Song{
	Singer:            "",
	TitleSong:         "",
	LyricsSong:        "",
	ImageURL:          "",
	Scores:            0,
	RemainingAttempts: 5,
	Timer:             30,
	ThePlaylist:       nil,
	CheatMess:         "",
}

func GetLyrics(artist, title string) (string, error) {
	url := fmt.Sprintf("https://api.lyrics.ovh/v1/%s/%s", artist, title)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to retrieve lyrics, status code: %d", resp.StatusCode)
	}

	var lyricsResponse LyricsResponse
	if err := json.NewDecoder(resp.Body).Decode(&lyricsResponse); err != nil {
		return "", err
	}

	lines := strings.Split(lyricsResponse.Lyrics, "\n")
	return strings.Join(lines[:20], "\n"), nil
}

func LoadData() {
	track := Api()
	artist := GetArtistsNames(track.Artists)
	title := track.Name

	lyrics, err := GetLyrics(artist, title)
	if err != nil {
		log.Fatalf("error retrieving lyrics: %v", err)
	}

	CurrentSong.Singer = artist
	CurrentSong.TitleSong = title
	CurrentSong.LyricsSong = lyrics
	CurrentSong.ImageURL = data.imageUrl

	fmt.Println("artist", artist)
	fmt.Println("Title", title)
	fmt.Println("lyrics :", lyrics)
	fmt.Println("imageUrl:", data.imageUrl)
	// fmt.Println("artist" , currentSong.Singer)
	// fmt.Println("Title", currentSong.TitleSong)
	// fmt.Println("lyrics :", currentSong.LyricsSong)
}
