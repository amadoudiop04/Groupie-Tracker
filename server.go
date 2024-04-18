package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"strings"

	"./database"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := database.InitDatabase("USER")
	defer db.Close()

	db.Exec("DELETE FROM USER WHERE id > 1;") //--> Remove some users

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
func LoginPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Login.html", nil)
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
		http.Redirect(w, r, "/Home", http.StatusSeeOther)
	} else {

		data := struct {
			ErrorMessage string
		}{
			ErrorMessage: "Connexion impossible : nom d'utilisateur ou mot de passe incorrect.",
		}
		renderTemplate(w, "Login.html", data)
	}
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Register.html", nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	db := database.InitDatabase("USER")
	defer db.Close()

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	if !database.UniqueUsername(username) { //username already used
		data := struct {
			ErrorMessage string
		}{
			ErrorMessage: "Ce pseudo est déjà utilisé, veuillez choisir un autre pseudo.",
		}
		renderTemplate(w, "Register.html", data)
	} else if !database.UniqueEmail(email) { //email already used
		data := struct {
			ErrorMessage string
		}{
			ErrorMessage: "Cet email est déjà utilisé, veuillez choisir un autre email.",
		}
		renderTemplate(w, "Register.html", data)
	} else if !database.VerifyPassword(password) { //password doesn't follow CNIL recommendations
		data := struct {
			ErrorMessage string
		}{
			ErrorMessage: "Votre mot de passe doit contenir 12 caractères comprenant des majuscules, des minuscules, des chiffres et des caractères spéciaux.",
		}
		renderTemplate(w, "Register.html", data)
	} else if confirmPassword != password { //password and password confirmation don't match
		data := struct {
			ErrorMessage string
		}{
			ErrorMessage: "Le mot de passe et la confirmation du mot de passe ne correspondent pas.",
		}
		renderTemplate(w, "Register.html", data)
	} else { //Account is valid, we can create it
		err := database.RegisterUser(db, username, password, email)
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else {
			rowsUsers := database.SelectAllFromTable(db, "USER")
			database.DisplayUserTable(rowsUsers) //--> show the table with the new user
			http.Redirect(w, r, "/Home", http.StatusSeeOther)
		}
	}
}

func PasswordForgottenPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "PasswordForgotten.html", nil)
}

func PasswordForgottenHandler(w http.ResponseWriter, r *http.Request) {
	db := database.InitDatabase("USER")
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
		renderTemplate(w, "PasswordForgotten.html", data)
	} else {
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
	renderTemplate(w, "AccountRecovery.html", nil)
}

func AccountRecoveryHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/ResetPassword", http.StatusSeeOther)
}

func ResetPasswordPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "ResetPassword.html", nil)
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	db := database.InitDatabase("USER")
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
