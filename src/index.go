// package main

// import (
// 	"fmt"
// 	"math/rand"
// 	"time"
// )

// type Music struct {
// 	Img    string
// 	Name   string
// 	Artist string
// 	Music  string
// }

// var (
// 	nowPlaying      string
// 	trackArt        string
// 	trackName       string
// 	trackArtist     string
// 	playPauseBtn    string
// 	nextBtn         string
// 	prevBtn         string
// 	seekSlider      string
// 	volumeSlider    string
// 	currTime        string
// 	totalDuration   string
// 	wave            string
// 	randomIcon      string
// 	currTrack       string
// 	trackIndex      int
// 	isPlaying       bool
// 	isRandom        bool
// 	musicList       []Music
// 	updateTimer     *time.Timer
// )

// func mainLoader() {
// 	musicList = []Music{
// 		{Img: "images/stay.png", Name: "Stay", Artist: "The Kid LAROI, Justin Bieber", Music: "musics/stay.mp3"},
// 		{Img: "images/faded.png", Name: "Falling Down", Artist: "Wid Cards", Music: "musics/fallingdown.mp3"},
// 		{Img: "images/faded.png", Name: "Faded", Artist: "Alan Walker", Music: "musics/Faded.mp3"},
// 		{Img: "images/ratherbe.png", Name: "Rather Be", Artist: "Clean Bandit", Music: "musics/RatherBe.mp3"},
// 	}
// 	loadTrack(trackIndex)
// }

// func loadTrack(trackIndex int) {
// 	reset()

// 	currTrack = musicList[trackIndex].Music
// 	trackArt = musicList[trackIndex].Img
// 	trackName = musicList[trackIndex].Name
// 	trackArtist = musicList[trackIndex].Artist
// 	nowPlaying = fmt.Sprintf("Playing music %d of %d", trackIndex+1, len(musicList))

// 	updateTimer = time.NewTimer(1 * time.Second)
// 	go setUpdate()

// 	currTrackListener := make(chan bool)
// 	go func() {
// 		for {
// 			select {
// 			case <-updateTimer.C:
// 				nextTrack()
// 			case <-currTrackListener:
// 				return
// 			}
// 		}
// 	}()

// 	randomBgColor()
// }

// func randomBgColor() {
// 	hex := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e"}
// 	var a string

// 	populate := func(a string) string {
// 		for i := 0; i < 6; i++ {
// 			x := rand.Intn(15)
// 			y := hex[x]
// 			a += y
// 		}
// 		return a
// 	}

// 	Color1 := populate("#")
// 	Color2 := populate("#")
// 	angle := "to right"

// 	gradient := fmt.Sprintf("linear-gradient(%s,%s,%s)", angle, Color1, Color2)
// 	fmt.Println("Background:", gradient)
// }

// func reset() {
// 	currTime = "00:00"
// 	totalDuration = "00:00"
// 	seekSlider = "0"
// }

// func playpauseTrack() {
// 	if isPlaying {
// 		pauseTrack()
// 	} else {
// 		playTrack()
// 	}
// }

// func playTrack() {
// 	isPlaying = true
// 	fmt.Println("Playing track:", currTrack)
// 	fmt.Println("Track name:", trackName)
	
// }

// func pauseTrack() {
// 	isPlaying = false
// 	fmt.Println("Pausing track:", currTrack)

// }

// func nextTrack() {
// 	if trackIndex < len(musicList)-1 && !isRandom {
// 		trackIndex++
// 	} else if trackIndex < len(musicList)-1 && isRandom {
// 		randomIndex := rand.Intn(len(musicList))
// 		trackIndex = randomIndex
// 	} else {
// 		trackIndex = 0
// 	}
// 	loadTrack(trackIndex)
// 	playTrack()
// }

// func setUpdate() {
// 	for {
// 		<-updateTimer.C
// 		fmt.Println("Updating track information...")
		
// 	}
// }


package main

import (
    "fmt"
    "net/http"
)

type PlaybackState struct {
    IsPlaying bool 
  
}

var playbackState = PlaybackState{
    IsPlaying: false,
}


func playHandler(w http.ResponseWriter, r *http.Request) {
    playbackState.IsPlaying = !playbackState.IsPlaying

    if playbackState.IsPlaying {
        fmt.Fprintf(w, "Playback started")
    } else {
        fmt.Fprintf(w, "Playback paused")
    }
}

func nextHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Next track")
}

func previousHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Previous track")
}

func restartHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Restarting track")
}

func main() {
    http.HandleFunc("/play", playHandler)
    http.HandleFunc("/next", nextHandler)
    http.HandleFunc("/previous", previousHandler)
    http.HandleFunc("/restart", restartHandler)
    fmt.Println("Server running on :8080")
    http.ListenAndServe(":8080", nil)
}
