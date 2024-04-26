package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"strings"

	"groupieTracker/database"
	"groupieTracker/games"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := database.InitTable("USER")
	defer db.Close()

	db.Exec("DELETE FROM USER WHERE id > 1;") //--> Remove users with id > 1  /!\ TO REMOVE BEFORE DEPLOYMENT /!\

	rowsUsers := database.SelectAllFromTable(db, "USER")
	database.DisplayUserTable(rowsUsers) //--> Show the table USER in terminal

	http.HandleFunc("/", LoginPage)
	http.HandleFunc("/LoginHandler", LoginHandler)
	http.HandleFunc("/Register", RegisterPage)
	http.HandleFunc("/RegisterHandler", RegisterHandler)
	http.HandleFunc("/PasswordForgotten", PasswordForgottenPage)
	http.HandleFunc("/PasswordForgottenHandler", PasswordForgottenHandler)
	http.HandleFunc("/AccountRecovery", AccountRecoveryPage)
	http.HandleFunc("/AccountRecoveryHandler", AccountRecoveryHandler)
	http.HandleFunc("/ResetPassword", ResetPasswordPage)
	http.HandleFunc("/ResetPasswordHandler", ResetPasswordHandler)
	http.HandleFunc("/Home", HomePage)
	http.HandleFunc("/UserProfile", UserProfile)
	http.HandleFunc("/BlindtestLandingPage", BlindtestLandingPage)
	http.HandleFunc("/Blindtest", Blindtest)
	http.HandleFunc("/EndBlindtest", EndBlindtest)
	http.HandleFunc("/GuessTheSong", GuessTheSong)
	http.HandleFunc("/GuessTheSongLose", GuessTheSongLose)
	http.HandleFunc("/GuessTheSongWin", GuessTheSongWin)
	http.HandleFunc("/GuessTheSongInfo", GuessTheSongInfo)
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
		return
	}
}

// Routes
func LoginPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Login/Login.html", nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	canConnect, err := database.AuthenticateUser(username, password)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if canConnect {
		sessionID, err := database.GetUserID(username)

		if err != nil {
			http.Error(w, "Erreur de récupération de l'ID de session", http.StatusInternalServerError)
			return
		}

		cookie := http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   2592000,
		}

		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/Home", http.StatusSeeOther)
	} else {
		data := struct {
			ErrorMessage string
		}{
			ErrorMessage: "Connexion impossible : nom d'utilisateur ou mot de passe incorrect.",
		}
		renderTemplate(w, "Login/Login.html", data)
	}
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Login/Register.html", nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	db := database.InitTable("USER")
	defer db.Close()

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	if username == "000" && password == "000" { //admin account
		database.RegisterUser(db, username, password, email)
		http.Redirect(w, r, "/Home", http.StatusSeeOther)
		return
	}

	if !database.UniqueUsername(username) { //username already used
		data := struct {
			ErrorMessage string
		}{
			ErrorMessage: "Ce pseudo est déjà utilisé, veuillez choisir un autre pseudo.",
		}
		renderTemplate(w, "Login/Register.html", data)
	} else if !database.UniqueEmail(email) { //email already used
		data := struct {
			ErrorMessage string
		}{
			ErrorMessage: "Cet email est déjà utilisé, veuillez choisir un autre email.",
		}
		renderTemplate(w, "Login/Register.html", data)
	} else if !database.VerifyPassword(password) { //password doesn't follow CNIL recommendations
		data := struct {
			ErrorMessage string
		}{
			ErrorMessage: "Votre mot de passe doit contenir 12 caractères comprenant des majuscules, des minuscules, des chiffres et des caractères spéciaux.",
		}
		renderTemplate(w, "Login/Register.html", data)
	} else if confirmPassword != password { //password and password confirmation don't match
		data := struct {
			ErrorMessage string
		}{
			ErrorMessage: "Le mot de passe et la confirmation du mot de passe ne correspondent pas.",
		}
		renderTemplate(w, "Login/Register.html", data)
	} else { //Account is valid, we can create it
		err := database.RegisterUser(db, username, password, email)
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else {
			sessionID, err := database.GetUserID(username)

			if err != nil {
				http.Error(w, "Erreur de récupération de l'ID de session", http.StatusInternalServerError)
				return
			}

			cookie := http.Cookie{
				Name:     "session_id",
				Value:    sessionID,
				Path:     "/",
				HttpOnly: true,
				MaxAge:   2592000,
			}

			http.SetCookie(w, &cookie)

			rowsUsers := database.SelectAllFromTable(db, "USER")
			database.DisplayUserTable(rowsUsers) //--> show the table with the new user

			http.Redirect(w, r, "/Home", http.StatusSeeOther)
		}
	}
}

func PasswordForgottenPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Login/PasswordForgotten.html", nil)
}

func PasswordForgottenHandler(w http.ResponseWriter, r *http.Request) {
	db := database.InitTable("USER")
	defer db.Close()

	email := strings.ReplaceAll(r.FormValue("email"), " ", "")

	userExists := false
	query := "SELECT email FROM USER WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&email)
	if err == nil {
		userExists = true
	} else if err != sql.ErrNoRows {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !userExists {
		data := struct {
			ErrorMessage string
		}{
			ErrorMessage: "Cet email n'est associé à aucun compte.",
		}
		renderTemplate(w, "Login/PasswordForgotten.html", data)
	} else {
		// 			/!\ IT DOES'NT WORK FOR NOW /!\
		// Sender data
		from := "help.groupietracker@gmail.com"
		password := "GTamyagoGT"

		// Receiver email address
		to := []string{email}

		// smtp server configuration
		smtpHost := "smtp.gmail.com"
		smtpPort := "587"

		// Code
		code := "5846"

		// Message
		message := []byte(code)

		// Authentication
		auth := smtp.PlainAuth("", from, password, smtpHost)

		// Sending email
		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Email Sent Successfully!")
		http.Redirect(w, r, "/AccountRecovery", http.StatusSeeOther)
	}
}

func AccountRecoveryPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Login/AccountRecovery.html", nil)
}

func AccountRecoveryHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/ResetPassword", http.StatusSeeOther)
}

func ResetPasswordPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Login/ResetPassword.html", nil)
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	db := database.InitTable("USER")
	defer db.Close()

	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	if password != confirmPassword {
		data := struct {
			ErrorMessage string
		}{
			ErrorMessage: "Le mot de passe et la confirmation du mot de passe ne correspondent pas.",
		}
		renderTemplate(w, "ResetPassword.html", data)
	} else if !database.VerifyPassword(password) {
		data := struct {
			ErrorMessage string
		}{
			ErrorMessage: "Votre mot de passe doit contenir 12 caractères comprenant des majuscules, des minuscules, des chiffres et des caractères spéciaux.",
		}
		renderTemplate(w, "ResetPassword.html", data)
	} else {
		//updateQuery := "UPDATE USER SET password = ? WHERE email = ?"     // --> find a way to change the wright user (have to register the email)
		//_, err := db.Exec(updateQuery, password, email)
		//if err != nil {
		//	log.Print(err)
		//	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		//	return
		//}
		http.Redirect(w, r, "/Home", http.StatusSeeOther)
	}
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Home/Home.html", nil)
}

func UserProfile(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	sessionID := cookie.Value

	user, err := database.GetUserData(sessionID)
	if err != nil {
		log.Println("Erreur lors de la récupération des données de l'utilisateur:", err)
		http.Error(w, "Erreur lors de la récupération des données de l'utilisateur", http.StatusInternalServerError)
		return
	}

	userData := struct {
		ID       int
		Username string
		Email    string
	}{
		ID:       user.Id,
		Username: user.Pseudo,
		Email:    user.Email,
	}

	renderTemplate(w, "Home/UserProfile.html", userData)
}

func BlindtestLandingPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "BlindTest/LandingPage.html", nil)
}

func Blindtest(w http.ResponseWriter, r *http.Request) {
	tracks := games.Api("6Xf0gjt1YmwvEG5iS8QOfg?si=2de553d01ff84abb")
	tracks = games.RemovePlayedTracks(tracks)
	currentTrack := games.NextTrack(tracks)

	if currentTrack == nil {
		http.Redirect(w, r, "/EndBlindtest", http.StatusSeeOther)
	}

	games.PlayedTracks = append(games.PlayedTracks, currentTrack)

	data := games.PageData{
		Track: currentTrack,
	}

	w.Header().Set("Refresh", "16")
	renderTemplate(w, "BlindTest/index.html", data)

}

func EndBlindtest(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "BlindTest/EndGame.html", nil)
}

func GuessTheSong(w http.ResponseWriter, r *http.Request) {
	games.LoadData()
	w.Header().Set("Refresh", "20")
	html := template.Must(template.ParseFiles("html/GuessTheSong/index.html"))
	if r.Method == "POST" {

		action := r.FormValue("action")
		if action == "next" {
			fmt.Println(" Next")
			games.NextTracks()
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
		http.Redirect(w, r, "/GuessTheSongLose", http.StatusSeeOther)
	}

	if games.CurrentSong.Scores == 50 {
		http.Redirect(w, r, "/GuessTheSongWin", http.StatusSeeOther)
	}

	err := html.Execute(w, games.CurrentSong)
	if err != nil {
		return
	}
}

func GuessTheSongLose(w http.ResponseWriter, r *http.Request) {
	html := template.Must(template.ParseFiles("html/GuessTheSong/lose.html"))

	if r.Method == "POST" {
		games.ResetData()
		http.Redirect(w, r, "/GuessTheSong", http.StatusSeeOther)
	}

	err := html.Execute(w, nil)
	if err != nil {
		return
	}
}

func GuessTheSongWin(w http.ResponseWriter, r *http.Request) {
	html := template.Must(template.ParseFiles("html/GuessTheSong/win.html"))

	if r.Method == "POST" {
		games.ResetData()
		http.Redirect(w, r, "/GuessTheSong", http.StatusSeeOther)
	}

	err := html.Execute(w, nil)
	if err != nil {
		return
	}
}

func GuessTheSongInfo(w http.ResponseWriter, r *http.Request) {
	html := template.Must(template.ParseFiles("html/GuessTheSong/info.html"))

	if r.Method == "POST" {
		games.ResetData()
		http.Redirect(w, r, "/GuessTheSong", http.StatusSeeOther)
	}

	err := html.Execute(w, nil)
	if err != nil {
		return
	}
}

func Petitbac(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "PetitBac/index.html", nil)
}
