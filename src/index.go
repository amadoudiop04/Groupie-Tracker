package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)


// type song struct{
// 	singer string
// 	titleSong string
// 	lyricsSong string
// }


type LyricsResponse struct {
	Lyrics string `json:"lyrics"`
}

func GetLyrics(artist, title string) (string, error) {
	url := fmt.Sprintf("https://api.lyrics.ovh/v1/%s/%s", artist, title)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to retrieve lyrics, status code: %d", resp.StatusCode)
	}

	var lyricsResponse LyricsResponse
	if err := json.NewDecoder(resp.Body).Decode(&lyricsResponse); err != nil {
		return "", err
	}

	return lyricsResponse.Lyrics, nil
}

func main() {
	track := Api()
	artist := getArtistsNames(track.Artists)
	title := track.Name

	lyrics, err := GetLyrics(artist, title)
	if err != nil {
		log.Fatalf("error retrieving lyrics: %v", err)
	}

	fmt.Println("artist" ,artist)
	fmt.Println("Title", title)
	fmt.Println("lyrics :", lyrics)
}




// var game = guessTheSong {
// 	singer: GetLyrics(artist) ,
	
// }




















































































// package main

// import (
// 	"fmt"
// )

// func main(){
// 	playlist := Api()
// 	fmt.Println("PlayList ID", playlist)
// }


