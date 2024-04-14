package main

import (
	"fmt"
	api "guessthesong/src"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	timerDone = make(chan struct{})
)

func ResetData() {
	api.CurrentSong.Scores = 0
	api.CurrentSong.RemainingAttempts = 5
	api.CurrentSong.Timer = 60
}

func StartTimer() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			api.CurrentSong.Timer--
			if api.CurrentSong.Timer <= 0 {
				close(timerDone)
				api.NextTrack()
				return
			}
		case <-timerDone:
			return
		}
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	html := template.Must(template.ParseFiles("templates/index.html"))

	if r.Method == "POST" {
		input := r.FormValue("value")
		if input == (strings.ToLower(api.CurrentSong.TitleSong)) {
			api.CurrentSong.Scores += 10
		} else {
			api.CurrentSong.RemainingAttempts--
		}
	}

	if api.CurrentSong.RemainingAttempts == 0 {
		http.Redirect(w, r, "/lose", http.StatusSeeOther)
	}

	if api.CurrentSong.Scores == 50 {
		http.Redirect(w, r, "/win", http.StatusSeeOther)
	}

	err := html.Execute(w, api.CurrentSong)
	if err != nil {
		return
	}

}

func Lost(w http.ResponseWriter, r *http.Request) {
	html := template.Must(template.ParseFiles("templates/lose.html"))

	if r.Method == "POST" {
		ResetData()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	err := html.Execute(w, nil)
	if err != nil {
		return
	}
}

func Win(w http.ResponseWriter, r *http.Request) {
	html := template.Must(template.ParseFiles("templates/win.html"))

	if r.Method == "POST" {
		ResetData()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	err := html.Execute(w, nil)
	if err != nil {
		return
	}
}

func SendJqueryJs(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("anim.js")
	if err != nil {
		http.Error(w, "Couldn't read file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	_, err = w.Write(data)
	if err != nil {
		return
	}
}

func TimerHandler(w http.ResponseWriter, r *http.Request) {
	timer := api.CurrentSong.Timer
	w.Header().Set("content-Type", "application/json")
	fmt.Fprintf(w, `{"time": %d}`, timer)
}

func main() {
	api.LoadData()
	go StartTimer()
	http.HandleFunc("/", Index)
	http.HandleFunc("/lose", Lost)
	http.HandleFunc("/win", Win)
	http.HandleFunc("/anim.js", SendJqueryJs)
	http.HandleFunc("/timer", TimerHandler)
	fs := http.FileServer(http.Dir("./static/"))
	fs2 := http.FileServer(http.Dir("images/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/images/", http.StripPrefix("/images/", fs2))
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		return
	}
}
