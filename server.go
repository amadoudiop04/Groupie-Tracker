package main

import(
	"html/template"
	"net/http"
)
	


func index(w http.ResponseWriter, r *http.Request) {
	html := template.Must(template.ParseFiles("templates/index.html"))
	html.Execute(w, nil)
}

	
func main() {
	http.HandleFunc("/", index)
	fs := http.FileServer(http.Dir("./static/"))
	fs2 := http.FileServer(http.Dir("images/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/images/", http.StripPrefix("/images/", fs2))
	http.ListenAndServe(":8000", nil)
}


