package main

import (
	"fmt"
	games "guessthesong/go"
	"html/template"
	"net/http"
)

func GuessTheSong(w http.ResponseWriter, r *http.Request) {
	games.LoadData()
	html := template.Must(template.ParseFiles("html/guessthesong/index.html"))
	if r.Method == "POST" {

		action := r.FormValue("action")
		if action == "next" {
			fmt.Println(" Next")
			games.NextTrack()
		}
		if action == "previous" {
			fmt.Println(" previous")
			games.CurrentSong.CheatMess = "Please don't cheat ❌"
			// Code pour previous
		}
		if action == "playPause" {
			fmt.Println("play")
			games.CurrentSong.CheatMess = "Please don't cheat❌"
			// Code pour Play
		}

		input := r.FormValue("value")
		if games.CompareStrings(input, games.CurrentSong.TitleSong) {
			games.CurrentSong.Scores += 10
		} else {
			games.CurrentSong.RemainingAttempts--
		}
	}

	if games.CurrentSong.RemainingAttempts == 0 {
		http.Redirect(w, r, "/lose", http.StatusSeeOther)
	}

	if games.CurrentSong.Scores == 50 {
		http.Redirect(w, r, "/win", http.StatusSeeOther)
	}

	err := html.Execute(w, games.CurrentSong)
	if err != nil {
		return
	}
}

func GuessTheSongLose(w http.ResponseWriter, r *http.Request) {
	html := template.Must(template.ParseFiles("html/guessthesong/lose.html"))

	if r.Method == "POST" {
		games.ResetData()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	err := html.Execute(w, nil)
	if err != nil {
		return
	}
}

func GuessTheSongWin(w http.ResponseWriter, r *http.Request) {
	html := template.Must(template.ParseFiles("html/guessthesong/win.html"))

	if r.Method == "POST" {
		games.ResetData()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	err := html.Execute(w, nil)
	if err != nil {
		return
	}
}

func GuessTheSongInfo(w http.ResponseWriter, r *http.Request) {
	html := template.Must(template.ParseFiles("html/guessthesong/info.html"))

	if r.Method == "POST" {
		games.ResetData()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	err := html.Execute(w, nil)
	if err != nil {
		return
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/GuessTheSong", http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/", redirect)
	http.HandleFunc("/GuessTheSong", GuessTheSong)
	http.HandleFunc("/GuessTheSongLose", GuessTheSongLose)
	http.HandleFunc("/GuessTheSongWin", GuessTheSongWin)
	http.HandleFunc("/GuessTheSongInfo", GuessTheSongInfo)

	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		return
	}
}
