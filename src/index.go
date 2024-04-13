package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Song struct {
	Singer     string
	TitleSong  string
	LyricsSong string
}

type LyricsResponse struct {
	Lyrics string `json:"lyrics"`
}

var CurrentSong = Song{
	Singer:     "",
	TitleSong:  "",
	LyricsSong: "",
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

	return lyricsResponse.Lyrics, nil
}

func LoadData() {
	track := Api()
	artist := getArtistsNames(track.Artists)
	title := track.Name

	lyrics, err := GetLyrics(artist, title)
	if err != nil {
		log.Fatalf("error retrieving lyrics: %v", err)
	}

	CurrentSong.Singer = artist
	CurrentSong.TitleSong = title
	CurrentSong.LyricsSong = lyrics

	fmt.Println("artist", artist)
	fmt.Println("Title", title)
	fmt.Println("lyrics :", lyrics)
	// fmt.Println("artist" , currentSong.Singer)
	// fmt.Println("Title", currentSong.TitleSong)
	// fmt.Println("lyrics :", currentSong.LyricsSong)
}
