package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", handleHome)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
    log.Println("Serveur démarré sur :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
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
        http.ServeFile(w, r, "index.html")
    }
}




// var categories = []string{"Artiste", "Album", "Groupe de musique", "Instrument de musique", "Featuring"}

// func contains(s []string, e string) bool {
//     for _, a := range s {
//         if a == e {
//             return true
//         }
//     }
//     return false
// }

// func isUnique(s []string, e string) bool {
//     count := 0
//     for _, a := range s {
//         if a == e {
//             count++
//         }
//     }
//     return count == 1
// }



// func bac() {
    

   
// }
// rand.Seed(time.Now().UnixNano())
//     letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
//     randomLetter := string(letters[rand.Intn(len(letters))])

//     fmt.Println(randomLetter)