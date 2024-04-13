package main

import (
	"fmt"
	api "guessthesong/src"
	"html/template"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	html := template.Must(template.ParseFiles("templates/index.html"))
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			return
		}
		value := r.FormValue("value")
		fmt.Println("value", value)
		_, err = fmt.Fprintf(w, "Value received: %s", value)
		if err != nil {
			return
		}
	}

	err := html.Execute(w, api.CurrentSong)
	if err != nil {
		return
	}
}

func main() {
	api.LoadData()
	http.HandleFunc("/", Index)
	fs := http.FileServer(http.Dir("./static/"))
	fs2 := http.FileServer(http.Dir("images/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/images/", http.StripPrefix("/images/", fs2))
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		return
	}
}
