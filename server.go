package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"groupieTracker/database"
	"groupieTracker/games"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
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
	http.HandleFunc("/UserProfileHandler", UserProfileHandler)
	http.HandleFunc("/BlindtestLandingPage", BlindtestLandingPage)
	http.HandleFunc("/Blindtest", Blindtest)
	http.HandleFunc("/EndBlindtest", EndBlindtest)
	http.HandleFunc("/BlindtestRules", BlindtestRules)
	http.HandleFunc("/GuessTheSongLandingPage", GuessTheSongLandingPage)
	http.HandleFunc("/GuessTheSong", GuessTheSong)
	http.HandleFunc("/GuessTheSongLose", GuessTheSongLose)
	http.HandleFunc("/GuessTheSongWin", GuessTheSongWin)
	http.HandleFunc("/GuessTheSongRules", GuessTheSongRules)
	http.HandleFunc("/PetitBacLandingPage", PetitBacLandingPage)
	http.HandleFunc("/PetitBac", PetitBac)
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
		data := struct{ ErrorMessage string }{ErrorMessage: "Connexion impossible : nom d'utilisateur ou mot de passe incorrect."}
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
		data := struct{ ErrorMessage string }{ErrorMessage: "Ce pseudo est déjà utilisé, veuillez choisir un autre pseudo."}
		renderTemplate(w, "Login/Register.html", data)
	} else if !database.UniqueEmail(email) { //email already used
		data := struct{ ErrorMessage string }{ErrorMessage: "Cet email est déjà utilisé, veuillez choisir un autre email."}
		renderTemplate(w, "Login/Register.html", data)
	} else if !database.VerifyPassword(password) { //password doesn't follow CNIL recommendations
		data := struct{ ErrorMessage string }{ErrorMessage: "Votre mot de passe doit contenir 12 caractères comprenant des majuscules, des minuscules, des chiffres et des caractères spéciaux."}
		renderTemplate(w, "Login/Register.html", data)
	} else if confirmPassword != password { //password and password confirmation don't match
		data := struct{ ErrorMessage string }{ErrorMessage: "Le mot de passe et la confirmation du mot de passe ne correspondent pas."}
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
		data := struct{ ErrorMessage string }{ErrorMessage: "Cet email n'est associé à aucun compte."}
		renderTemplate(w, "Login/PasswordForgotten.html", data)
	} else {
		apiKey := "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"
		sg := sendgrid.NewSendClient(apiKey)

		code := database.GenerateCode()
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Session cookie not found", http.StatusUnauthorized)
			return
		}
		userID := cookie.Value
		_, err = db.Exec("UPDATE USER SET recoveryCode = ? WHERE id = ?", code, userID)
		if err != nil {
			log.Print(err)
		}

		from := mail.NewEmail("Groupie Tracker", "help.groupietracker@gmail.com")
		subject := "Retrouvez votre compte"
		to := mail.NewEmail("Client", email)
		content := mail.NewContent("text/plain", "Bonjour, \n voici le code qui vous permettra de réinitialiser votre mot de passe : "+strconv.Itoa(code)+"\n À bientôt sur Groupie Tracker !")
		message := mail.NewV3MailInit(from, subject, to, content)

		response, err := sg.Send(message)
		if err != nil {
			log.Println("Erreur lors de l'envoi de l'e-mail:", err)
		} else {
			log.Println("Code de statut de l'envoi de l'e-mail:", response.StatusCode)
			log.Println("Réponse de l'API SendGrid:", response.Body)
			http.Redirect(w, r, "/AccountRecovery", http.StatusSeeOther)
		}
	}
}

func AccountRecoveryPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Login/AccountRecovery.html", nil)
}

func AccountRecoveryHandler(w http.ResponseWriter, r *http.Request) {
	db := database.InitTable("USER")
	defer db.Close()

	recoveryCode := r.FormValue("recoveryCode")
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Session cookie not found", http.StatusUnauthorized)
		return
	}
	userID := cookie.Value
	var userCode int
	_ = db.QueryRow("SELECT recoveryCode FROM USER WHERE id = ?", userID).Scan(&userCode)

	if recoveryCode == strconv.Itoa(userCode) {
		_, _ = db.Exec("UPDATE USER SET recoveryCode = NULL WHERE id = ?", userID)
		http.Redirect(w, r, "/ResetPassword", http.StatusSeeOther)
	} else {
		data := struct{ ErrorMessage string }{ErrorMessage: "Code incorrect."}
		renderTemplate(w, "Login/AccountRecovery.html", data)
	}
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
		data := struct{ ErrorMessage string }{ErrorMessage: "Le mot de passe et la confirmation du mot de passe ne correspondent pas."}
		renderTemplate(w, "Login/ResetPassword.html", data)
	} else if !database.VerifyPassword(password) {
		data := struct{ ErrorMessage string }{ErrorMessage: "Votre mot de passe doit contenir 12 caractères comprenant des majuscules, des minuscules, des chiffres et des caractères spéciaux."}
		renderTemplate(w, "Login/ResetPassword.html", data)
	} else {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Session cookie not found", http.StatusUnauthorized)
			return
		}
		userID := cookie.Value
		_, err = db.Exec("UPDATE USER SET password = ? WHERE id = ?", database.HashPassword(password), userID)
		fmt.Println("Password set : " + password)
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
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
		ID           int
		Username     string
		Email        string
		ErrorMessage string
	}{
		ID:           user.Id,
		Username:     user.Pseudo,
		Email:        user.Email,
		ErrorMessage: "",
	}
	renderTemplate(w, "Home/UserProfile.html", userData)
}

func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
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

	username := database.ReplaceEmptyString(r.FormValue("usernameInput"), user.Pseudo)
	email := database.ReplaceEmptyString(r.FormValue("emailInput"), user.Email)
	currentPassword := r.FormValue("currentPasswordInput")
	newPassword := r.FormValue("newPasswordInput")
	confirmNewPassword := r.FormValue("confirmNewPasswordInput")

	userData := struct {
		ID           int
		Username     string
		Email        string
		ErrorMessage string
	}{
		ID:           user.Id,
		Username:     user.Pseudo,
		Email:        user.Email,
		ErrorMessage: "",
	}

	if currentPassword == "" {
		err = database.SetUserData(sessionID, username, email, user.Password)
		if err != nil {
			log.Println("Erreur lors de la mise à jour des données de l'utilisateur:", err)
			http.Error(w, "Erreur lors de la mise à jour des données de l'utilisateur", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/UserProfile", http.StatusSeeOther)
	} else if database.HashPassword(currentPassword) != user.Password {
		userData.ErrorMessage = "Le mot de passe actuel saisit est incorrect."
		renderTemplate(w, "Home/UserProfile.html", userData)
	} else if newPassword != confirmNewPassword {
		userData.ErrorMessage = "Le mot de passe et la confirmation du mot de passe ne correspondent pas."
		renderTemplate(w, "Home/UserProfile.html", userData)
	} else if !database.VerifyPassword(newPassword) {
		userData.ErrorMessage = "Votre nouveau mot de passe doit contenir 12 caractères comprenant des majuscules, des minuscules, des chiffres et des caractères spéciaux."
		renderTemplate(w, "Home/UserProfile.html", userData)
	} else {
		err = database.SetUserData(sessionID, username, email, database.HashPassword(newPassword))
		if err != nil {
			log.Println("Erreur lors de la mise à jour des données de l'utilisateur:", err)
			http.Error(w, "Erreur lors de la mise à jour des données de l'utilisateur", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/UserProfile", http.StatusSeeOther)
	}
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

func GuessTheSongLandingPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "GuessTheSong/LandingPage.html", nil)
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
	html := template.Must(template.ParseFiles("html/GuessTheSong/Lose.html"))

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
	html := template.Must(template.ParseFiles("html/GuessTheSong/Win.html"))

	if r.Method == "POST" {
		games.ResetData()
		http.Redirect(w, r, "/GuessTheSong", http.StatusSeeOther)
	}

	err := html.Execute(w, nil)
	if err != nil {
		return
	}
}

func GuessTheSongRules(w http.ResponseWriter, r *http.Request) {
	html := template.Must(template.ParseFiles("html/GuessTheSong/Rules.html"))

	if r.Method == "POST" {
		games.ResetData()
		http.Redirect(w, r, "/GuessTheSong", http.StatusSeeOther)
	}

	err := html.Execute(w, nil)
	if err != nil {
		return
	}
}

func BlindtestRules(w http.ResponseWriter, r *http.Request) {
	html := template.Must(template.ParseFiles("html/BlindTest/Rules.html"))

	if r.Method == "POST" {
		games.ResetData()
		http.Redirect(w, r, "/Blindtest", http.StatusSeeOther)
	}

	err := html.Execute(w, nil)
	if err != nil {
		return
	}
}

func PetitBacLandingPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "PetitBac/LandingPage.html", nil)
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

		data := games.Data{
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
