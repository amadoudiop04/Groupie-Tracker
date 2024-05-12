package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"groupieTracker/database"
	"groupieTracker/games"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var Pseudo string
var messages []Message
var ChatDiscours string

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Data struct {
	DatasgameBlindTest games.PageData
	Datasgames         games.Song
	RandomLetter       string
	Info               []Message
}

type Message struct {
	Username    string
	TextMessage string
	GamesDatas  string
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var mutex = sync.Mutex{} // Add mutex to handle concurrent map access

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			return
		}
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		log.Println(msg)
		mutex.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Println(err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
		log.Println(msg)
	}
}

func main() {
	db := database.InitTable("USER")
	defer db.Close()
	roomDB := database.InitTable("ROOMS")
	gameRoomDB := database.InitTable("GAME_ROOM")
	roomUserDB := database.InitTable("ROOM_USERS")

	db.Exec("DELETE FROM USER WHERE id > 1;") // /!\ TO REMOVE BEFORE DEPLOYMENT /!\
	roomDB.Exec("DELETE FROM ROOMS")
	gameRoomDB.Exec("DELETE FROM GAME_ROOM")
	roomUserDB.Exec("DELETE FROM ROOM_USERS")

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
	http.HandleFunc("/BlindtestRoom", BlindtestRoom)
	http.HandleFunc("/BlindtestRoomHandler", BlindtestRoomHandler)
	http.HandleFunc("/PublicBlindtest", JoinPublicBlindtest)
	http.HandleFunc("/CreateBlindtest", CreateBlindtest)
	http.HandleFunc("/CreateBlindtestHandler", CreateBlindtestHandler)
	http.HandleFunc("/JoinBlindtest", JoinPrivateBlindtest)
	http.HandleFunc("/JoinBlindtestHandler", JoinBlindtestHandler)
	http.HandleFunc("/EndBlindtest", EndBlindtest)
	http.HandleFunc("/BlindtestRules", BlindtestRules)
	http.HandleFunc("/GuessTheSongLandingPage", GuessTheSongLandingPage)
	http.HandleFunc("/GuessTheSong", GuessTheSong)
	http.HandleFunc("/GuessTheSongLose", GuessTheSongLose)
	http.HandleFunc("/GuessTheSongWin", GuessTheSongWin)
	http.HandleFunc("/GuessTheSongRules", GuessTheSongRules)
	http.HandleFunc("/PetitBacLandingPage", PetitBacLandingPage)
	http.HandleFunc("/PetitBac", PetitBac)
	http.HandleFunc("/PetitBacHandler", PetitBacHandler)
	http.HandleFunc("/PetitBacRules", PetitBacRules)
	go handleMessages()
	http.HandleFunc("/websocket", websocketHandler)
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
	Pseudo = username
	canConnect, err := database.AuthenticateUser(username, password)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if canConnect {
		sessionID, err := database.GetUserIdByUsername(username)

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
	email := strings.ReplaceAll(r.FormValue("email"), " ", "")
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
			sessionID, err := database.GetUserIdByUsername(username)

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
		//create cookie
		sessionID, err := database.GetUserIdByEmail(email)

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

		//send email
		apiKey := "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"
		sg := sendgrid.NewSendClient(apiKey)

		code := database.GenerateCode()
		_, err = db.Exec("UPDATE USER SET recoveryCode = ? WHERE id = ?", code, sessionID)
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
		fmt.Println(recoveryCode, userCode)
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
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	userID, _ := strconv.Atoi(cookie.Value)

	roomID, err := database.GetRoomIDByUserID(userID)
	//If the user is associate to a room, we have to make him leave
	if err == nil {
		database.LeaveRoom(roomID, userID)
		fmt.Println("user is leaving the room")
	}

	renderTemplate(w, "BlindTest/LandingPage.html", nil)
}

func Blindtest(w http.ResponseWriter, r *http.Request) {
	//Recovering the room and its data
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	userID, _ := strconv.Atoi(cookie.Value)

	roomID, err := database.GetRoomIDByUserID(userID)
	if err != nil {
		http.Redirect(w, r, "/JoinBlindtest", http.StatusSeeOther)
	}

	roomData, err := database.GetRoomData(roomID)
	if err != nil {
		fmt.Println("La room n'est associée à aucune données")
	}

	//Online chat
	if r.Method == "POST" {
		action := r.FormValue("action")
		if action == "ChatMessage" {
			message := r.Form.Get("Message")
			ChatDiscours = message
			messages = append(messages, Message{
				Username:    Pseudo,
				TextMessage: ChatDiscours,
			})

			jsonMessage, err := json.Marshal(messages[len(messages)-1])
			if err != nil {
				fmt.Println("Error marshaling message to JSON:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			broadcast <- Message{
				Username:    Pseudo,
				TextMessage: string(jsonMessage),
			}
		}
		if action == "DeleteMessage" {
			indexStr := r.Form.Get("messageIndex")
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				http.Error(w, "Invalid message index", http.StatusBadRequest)
				return
			}
			if index < 0 || index >= len(messages) {
				http.Error(w, "Invalid message index", http.StatusBadRequest)
				return
			}
			messages = append(messages[:index], messages[index+1:]...)
		}
	}

	//Game
	if userID == database.GetRoomCreator(roomID) {
		w.Header().Set("Refresh", strconv.Itoa(roomData.BlindtestTimeOfMusic+roomData.BlindtestTimeToAnswer))
		tracks := games.Api("6Xf0gjt1YmwvEG5iS8QOfg?si=2de553d01ff84abb")
		tracks = games.RemovePlayedTracks(tracks)
		currentTrack := games.NextTrack(tracks)

		if currentTrack == nil {
			http.Redirect(w, r, "/EndBlindtest", http.StatusSeeOther)
		}

		games.PlayedTracks = append(games.PlayedTracks, currentTrack)

		trackData := games.PageData{
			Track: currentTrack,
		}

		mediasBlindtest := Data{
			DatasgameBlindTest: trackData,
			Info:               messages,
		}

		//Data
		data := struct {
			MediasBlindtest  Data
			DurationOfMusic  int
			DurationOfAnswer int
		}{
			MediasBlindtest:  mediasBlindtest,
			DurationOfMusic:  roomData.BlindtestTimeOfMusic,
			DurationOfAnswer: roomData.BlindtestTimeToAnswer,
		}

		//Ici, recolter les données et les envoyer par websocket à tous les joueurs présents dans la room
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error marshaling message to JSON:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		broadcast <- Message{
			GamesDatas: string(jsonData),
		}
		//Execute html
		renderTemplate(w, "BlindTest/index.html", data)

	} else {
		//Ici, récolter les données reçues par websocket et les enregistrées dans data
		data := struct {
			MediasBlindtest  Data
			DurationOfMusic  int
			DurationOfAnswer int
		}{}
		renderTemplate(w, "BlindTest/index.html", data)
	}
}

func BlindtestRoom(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	userID, _ := strconv.Atoi(cookie.Value)
	roomID, err := database.GetRoomIDByUserID(userID)
	if err != nil {
		http.Redirect(w, r, "/JoinBlindtest", http.StatusSeeOther)
	}

	creatorID := database.GetRoomCreator(roomID)
	numberOfPlayer := database.GetNumberOfPlayer(roomID)

	mediasBlindtest := Data{Info: messages}

	data := struct {
		ButtonVisible   bool
		MediasBlindtest Data
		RoomID          int
		PlayerNumber    int
	}{
		ButtonVisible:   false,
		MediasBlindtest: mediasBlindtest,
		RoomID:          roomID,
		PlayerNumber:    numberOfPlayer,
	}

	if userID == creatorID {
		data.ButtonVisible = true
	}

	//Online chat
	if r.Method == "POST" {
		action := r.FormValue("action")
		if action == "ChatMessage" {
			message := r.Form.Get("Message")
			ChatDiscours = message
			messages = append(messages, Message{
				Username:    Pseudo,
				TextMessage: ChatDiscours,
			})

			jsonMessage, err := json.Marshal(messages[len(messages)-1])
			if err != nil {
				fmt.Println("Error marshaling message to JSON:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			broadcast <- Message{
				Username:    Pseudo,
				TextMessage: string(jsonMessage),
			}
		}
		if action == "DeleteMessage" {
			indexStr := r.Form.Get("messageIndex")
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				http.Error(w, "Invalid message index", http.StatusBadRequest)
				return
			}
			if index < 0 || index >= len(messages) {
				http.Error(w, "Invalid message index", http.StatusBadRequest)
				return
			}
			messages = append(messages[:index], messages[index+1:]...)
		}
	}

	//Execute html
	renderTemplate(w, "BlindTest/Room.html", data)
}

func BlindtestRoomHandler(w http.ResponseWriter, r *http.Request) {

}

func JoinPublicBlindtest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	userID, _ := strconv.Atoi(cookie.Value)
	roomID := 1

	if !database.CheckRoomExistence(roomID) {
		db := database.InitTable("ROOMS")
		defer db.Close()

		createdBy := userID
		maxPlayers := 50
		name := "publicBlindtest"
		gameID := 1

		// Create the public room
		_, err := db.Exec("INSERT INTO ROOMS (id, created_by, max_player, number_of_player, name, id_game) VALUES (?, ?, ?, ?, ?, ?)", roomID, createdBy, maxPlayers, 0, name, gameID)
		if err != nil {
			log.Println("Error creating room:", err)
		}

		//Room Data
		gameRoomDB := database.InitTable("GAME_ROOM")
		defer gameRoomDB.Close()

		numberOfGameTurns := 10
		timeOfMusic := 10
		timeToAnswer := 5

		_, err = gameRoomDB.Exec("INSERT INTO GAME_ROOM (id_room, number_of_game_turns, blindtest_time_of_music, blindtest_time_to_answer) VALUES (?, ?, ?, ?)", roomID, numberOfGameTurns, timeOfMusic, timeToAnswer)
		if err != nil {
			log.Println("Error creating game room:", err)
			return
		}
	}

	database.JoinRoom(roomID, userID)

	http.Redirect(w, r, "/Blindtest", http.StatusSeeOther)
}

func CreateBlindtest(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "BlindTest/CreatePrivateRoom.html", nil)
}

func CreateBlindtestHandler(w http.ResponseWriter, r *http.Request) {
	gameTurns, _ := strconv.Atoi(r.FormValue("gameTurns"))
	musicDuration, _ := strconv.Atoi(r.FormValue("musicDuration"))
	answerDuration, _ := strconv.Atoi(r.FormValue("answerDuration"))
	roomName := r.FormValue("roomName")
	maxPlayer, _ := strconv.Atoi(r.FormValue("maxPlayer"))

	if gameTurns <= 0 {
		gameTurns = 5
	}
	if musicDuration <= 0 {
		musicDuration = 10
	}
	if answerDuration <= 0 {
		answerDuration = 5
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	sessionID := cookie.Value

	userID, _ := strconv.Atoi(sessionID)
	gameID := 1

	database.CreateBlindtestRoom(userID, maxPlayer, roomName, gameID, gameTurns, musicDuration, answerDuration)

	http.Redirect(w, r, "/BlindtestRoom", http.StatusSeeOther)
}

func JoinPrivateBlindtest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	userID, _ := strconv.Atoi(cookie.Value)

	roomID, err := database.GetRoomIDByUserID(userID)
	//If the user is associate to a room, we have to make him leave
	if err == nil {
		database.LeaveRoom(roomID, userID)
		fmt.Println("user is leaving the room")
	}

	renderTemplate(w, "BlindTest/JoinPrivateRoom.html", nil)
}

func JoinBlindtestHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}
	userID, _ := strconv.Atoi(cookie.Value)
	roomID, err2 := strconv.Atoi(r.FormValue("roomID"))
	if err2 != nil || roomID <= 0 {
		data := struct{ ErrorMessage string }{ErrorMessage: "Veuillez rentrez un identifiant valide."}
		renderTemplate(w, "BlindTest/JoinPrivateRoom.html", data)
	} else if !database.VerifyRoom(roomID) {
		data := struct{ ErrorMessage string }{ErrorMessage: "Aucune room ne correspond à cet identifiant."}
		renderTemplate(w, "BlindTest/JoinPrivateRoom.html", data)
	} else {
		database.JoinRoom(roomID, userID)
		http.Redirect(w, r, "/BlindtestRoom", http.StatusSeeOther)
	}
}

func EndBlindtest(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "BlindTest/EndGame.html", nil)
}

func GuessTheSongLandingPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "GuessTheSong/LandingPage.html", nil)
}

func GuessTheSong(w http.ResponseWriter, r *http.Request) {
	games.LoadData()
	w.Header().Set("Refresh", "25")
	if r.Method == "POST" {
		action := r.FormValue("action")
		if action == "guessTheSong" {
			input := r.FormValue("value")
			if games.CompareStrings(input, games.CurrentSong.TitleSong) {
				games.CurrentSong.Scores += 10
			} else {
				games.CurrentSong.RemainingAttempts--
			}
			if games.CurrentSong.RemainingAttempts == 0 {
				http.Redirect(w, r, "/GuessTheSongLose", http.StatusSeeOther)
			}
			if games.CurrentSong.Scores == 50 {
				http.Redirect(w, r, "/GuessTheSongWin", http.StatusSeeOther)
			}
		}
		if action == "ChatMessage" {
			message := r.Form.Get("Message")
			ChatDiscours = message
			messages = append(messages, Message{
				Username:    Pseudo,
				TextMessage: ChatDiscours,
			})

			jsonMessage, err := json.Marshal(messages[len(messages)-1])
			if err != nil {
				fmt.Println("Error marshaling message to JSON:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			broadcast <- Message{
				Username:    Pseudo,
				TextMessage: string(jsonMessage),
			}
		}
		if action == "DeleteMessage" {
			indexStr := r.Form.Get("messageIndex")
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				http.Error(w, "Invalid message index", http.StatusBadRequest)
				return
			}
			if index < 0 || index >= len(messages) {
				http.Error(w, "Invalid message index", http.StatusBadRequest)
				return
			}
			messages = append(messages[:index], messages[index+1:]...)
		}
	}

	medias := Data{
		Datasgames: games.CurrentSong,
		Info:       messages,
	}

	html := template.Must(template.ParseFiles("html/GuessTheSong/index.html"))
	err := html.Execute(w, medias)
	if err != nil {
		fmt.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
	w.Header().Set("Refresh", "31")
	if r.Method == "POST" {
		action := r.FormValue("action")
		if action == "true" {
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
		}

		if action == "ChatMessage" {
			log.Println("put your message")
			message := r.Form.Get("Message")
			ChatDiscours = message
			messages = append(messages, Message{
				Username:    Pseudo,
				TextMessage: ChatDiscours,
			})

			jsonMessage, err := json.Marshal(messages[len(messages)-1])
			if err != nil {
				fmt.Println("Error marshaling message to JSON:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			broadcast <- Message{
				Username:    Pseudo,
				TextMessage: string(jsonMessage),
			}
		}
		if action == "DeleteMessage" {
			indexStr := r.Form.Get("messageIndex")
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				http.Error(w, "Invalid message index", http.StatusBadRequest)
				return
			}
			if index < 0 || index >= len(messages) {
				http.Error(w, "Invalid message index", http.StatusBadRequest)
				return
			}
			messages = append(messages[:index], messages[index+1:]...)
		}

	}

	rand.Seed(time.Now().UnixNano())
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	randomLetter := string(letters[rand.Intn(len(letters))])

	data := Data{
		RandomLetter: randomLetter,
		Info:         messages,
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

func PetitBacHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/PetitBac", http.StatusSeeOther)
}

func PetitBacRules(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "PetitBac/Rules.html", nil)
}
