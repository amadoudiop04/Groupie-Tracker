package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/Blindtest", Blindtest)
	http.HandleFunc("/Deaftest", Deaftest)
	http.HandleFunc("/Petitbac", Petitbac)

	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", nil)
}

func renderTemplate(w http.ResponseWriter, templatePath string, data interface{}) {
	tmpl, err := template.ParseFiles("./html/" + templatePath) //, "./html/templates/header.html", "./html/templates/head.html")
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Routes
func HomePage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Home.html", nil)
}

func Blindtest(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Blindtest.html", nil)
}

func Deaftest(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Deaftest.html", nil)
}

func Petitbac(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Petitbac.html", nil)
}
