package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Data struct {
	RandomLetter string
}

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/PetitBac", PetitBac)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	log.Println("Serveur démarré sur :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/PetitBac", http.StatusSeeOther)
}

func PetitBac(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()

		artiste := r.Form.Get("artiste")
		album := r.Form.Get("album")
		groupe := r.Form.Get("groupe")
		instrument := r.Form.Get("instrument")
		featuring := r.Form.Get("featuring")

		fmt.Println("Nouvelle entrée ajoutée :")
		fmt.Println("Artiste:", artiste)
		fmt.Println("Album:", album)
		fmt.Println("Groupe de musique:", groupe)
		fmt.Println("Instrument de musique:", instrument)
		fmt.Println("Featuring:", featuring)
	} else {
		rand.Seed(time.Now().UnixNano())
		letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		randomLetter := string(letters[rand.Intn(len(letters))])

		data := Data{
			RandomLetter: randomLetter,
		}

		tmpl, err := template.ParseFiles("html/PetitBac/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
