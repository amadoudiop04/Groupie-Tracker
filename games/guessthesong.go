package games

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"unicode"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
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

func ApiTrack() *spotify.FullTrack {
	authConfig := &clientcredentials.Config{
		ClientID:     "a8795237d8ea48a09bc79862cf40c00b",
		ClientSecret: "8847d2dc62fc4948b882a1069e16cb35",
		TokenURL:     spotify.TokenURL,
	}

	accessToken, err := authConfig.Token(context.Background())
	if err != nil {
		log.Printf("error retrieve access token: %v", err)
	}

	client := spotify.Authenticator{}.NewClient(accessToken)

	playlistID := spotify.ID("37i9dQZF1E39Yzkf9WsM5K?si=6ddc1f4641a848ef")
	playlist, err := client.GetPlaylist(playlistID)
	MyPlaylist = playlist
	if err != nil {
		log.Printf("error retrieve playlist data: %v", err)
	}

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

		artistID := track.Artists[0].ID
		artistDetails, err := client.GetArtist(artistID)
		if err != nil {
			log.Printf("Error retrieving artist details: %v", err)
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

func GetTrackInfo(playlist *spotify.FullPlaylist) (*spotify.FullTrack, error) {
	if playlist == nil || len(playlist.Tracks.Tracks) == 0 {
		return nil, errors.New("empty or nil playlist")
	}

	GameIndex = (GameIndex + 1) % len(playlist.Tracks.Tracks)
	if TheLyrics != "" {
		if GameIndex < 0 {
			GameIndex = 0
		} else {
			GameIndex = (GameIndex + 1) % len(playlist.Tracks.Tracks)
		}
		return &playlist.Tracks.Tracks[GameIndex].Track, nil
	}
	return &playlist.Tracks.Tracks[GameIndex].Track, nil
}

func NextTracks() {
	CurrentSong.ThePlaylist = MyPlaylist
	nextTrack, err := GetTrackInfo(CurrentSong.ThePlaylist)
	if err != nil {
		return
	}

	artist := GetArtistsNames(nextTrack.Artists)
	title := nextTrack.Name

	lyrics, err := GetLyrics(artist, title)
	if err != nil || lyrics == "" {
		log.Printf("Failed to retrieve lyrics for %s by %s, skipping...", title, artist)
		NextTracks()
		return
	}

	CurrentSong.Singer = artist
	CurrentSong.TitleSong = title
	CurrentSong.LyricsSong = truncLyrics(lyrics)
	CurrentSong.Timer = 30
}

type Song struct {
	Singer            string
	TitleSong         string
	LyricsSong        string
	ImageURL          string
	Scores            int
	RemainingAttempts int
	Timer             int
	ThePlaylist       *spotify.FullPlaylist
}

type LyricsResponse struct {
	Lyrics string `json:"lyrics"`
}

type gamesData struct {
	Name           string
	NumberOfTurn   int
	AnswerDuration int
}

var GuessTheSongData gamesData
var PetitBacData gamesData

var CurrentSong = Song{
	Singer:            "",
	TitleSong:         "",
	LyricsSong:        "",
	ImageURL:          "",
	Scores:            0,
	RemainingAttempts: 5,
	Timer:             30,
	ThePlaylist:       nil,
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
	track := ApiTrack()
	artist := GetArtistsNames(track.Artists)
	title := track.Name

	lyrics, err := GetLyrics(artist, title)
	if err != nil {
		fmt.Printf("error retrieving lyrics: %v", err)
	}

	CurrentSong.Singer = artist
	CurrentSong.TitleSong = title
	CurrentSong.LyricsSong = truncLyrics(lyrics)
	CurrentSong.ImageURL = data.imageUrl
}

func ResetData() {
	CurrentSong.Scores = 0
	CurrentSong.RemainingAttempts = 5
	CurrentSong.Timer = 30
}

func RemoveAccents(input string) string {
	var output string
	for _, char := range input {
		if unicode.Is(unicode.Mn, char) {
			continue
		}
		output += string(char)
	}
	return output
}

func CompareStrings(input1, input2 string) bool {

	input1 = RemoveAccents(strings.ToLower(input1))
	input2 = RemoveAccents(strings.ToLower(input2))

	return input1 == input2
}

func truncLyrics(lyrics string) string {
	lines := strings.Split(lyrics, "\n")

	skipLineBetweenBrackets := func(line string) bool {
		return strings.Contains(line, "[") && strings.Contains(line, "]")
	}

	filteredLines := make([]string, 0)
	for _, line := range lines[3:] {
		if !skipLineBetweenBrackets(line) && line != "" {
			filteredLines = append(filteredLines, line)
		}
	}

	randomLineIndex := rand.Intn(len(filteredLines) - 6)
	truncatedLyrics := strings.Join(filteredLines[randomLineIndex:randomLineIndex+5], "\n")

	return truncatedLyrics
}
