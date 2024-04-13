package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := initDatabase("USER")
	defer db.Close()

	//db.Exec("DELETE FROM USER WHERE id > 0;") //--> Remove some users

	rowsUsers := selectAllFromTable(db, "USER")
	displayUserTable(rowsUsers) //--> Show the table USER in terminal

	http.HandleFunc("/", LoginPage)
	http.HandleFunc("/LoginHandler", LoginHandler)
	http.HandleFunc("/Register", RegisterPage)
	http.HandleFunc("/RegisterHandler", RegisterHandler)
	http.HandleFunc("/PasswordForgotten", PasswordForgottenPage)
	http.HandleFunc("/PasswordForgottenHandler", PasswordForgottenHandler)
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

	canConnect, err := AuthenticateUser(username, password)
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
	db := initDatabase("USER")
	defer db.Close()

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	if confirmPassword == password {
		err := RegisterUser(db, username, password, email)
		if err != nil {
			if err.Error() == "Password incorrect" {
				data := struct {
					ErrorMessage string
				}{
					ErrorMessage: "Votre mot de passe doit contenir 12 caractères comprenant des majuscules, des minuscules, des chiffres et des caractères spéciaux.",
				}
				renderTemplate(w, "Register.html", data)
			} else {
				log.Print(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		} else {
			rowsUsers := selectAllFromTable(db, "USER")
			displayUserTable(rowsUsers) //--> show the table with the new user
			http.Redirect(w, r, "/Home", http.StatusSeeOther)
		}
	} else {
		data := struct {
			ErrorMessage string
		}{
			ErrorMessage: "Le mot de passe et la confirmation du mot de passe ne correspondent pas.",
		}
		renderTemplate(w, "Register.html", data)
	}
}

func PasswordForgottenPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "PasswordForgotten.html", nil)
}

func PasswordForgottenHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "AccountRecovery.html", nil)
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

func HashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func RegisterUser(db *sql.DB, username, password, email string) error {
	if !VerifyPassword(password) {
		return fmt.Errorf("Password incorrect")
	}

	hashedPassword := HashPassword(password)

	_, err := db.Exec("INSERT INTO USER (pseudo, email, password) VALUES (?, ?, ?)", username, email, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func VerifyPassword(password string) bool {
	if len(password) < 12 || strings.ToUpper(password) == password {
		return false
	}

	if ok, _ := regexp.MatchString(`[!@#$%^&*()_+{}\[\]:;<>,.?/~\-]`, password); !ok {
		return false
	}

	if ok, _ := regexp.MatchString(`[0-9]`, password); !ok {
		return false
	}

	return true
}

func AuthenticateUser(username, password string) (bool, error) {
	db := initDatabase("USER")
	defer db.Close()

	hashedPassword := HashPassword(password)

	var storedPassword string
	query := "SELECT password FROM USER WHERE pseudo = ? OR email = ?"
	err := db.QueryRow(query, username, username).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return storedPassword == hashedPassword, nil
}

func initDatabase(database string) *sql.DB {
	db, err := sql.Open("sqlite3", database)

	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
				CREATE TABLE IF NOT EXISTS USER (
					id INTEGER PRIMARY KEY,
					pseudo TEXT NOT NULL,
					email TEXT NOT NULL,
					password TEXT NOT NULL
				);
				`
	_, err = db.Exec(sqlStmt)

	if err != nil {
		log.Fatal(err)
	}

	return db
}

type User struct {
	id       int
	pseudo   string
	email    string
	password string
}

func displayUserTable(rows *sql.Rows) {
	for rows.Next() {
		var users User
		err := rows.Scan(&users.id, &users.pseudo, &users.email, &users.password)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(users)
	}
}

func selectAllFromTable(db *sql.DB, table string) *sql.Rows {
	query := "SELECT * FROM " + table
	result, _ := db.Query(query)
	return result
}
