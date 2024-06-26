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
	http.HandleFunc("/GuessTheSongRoom", GuessTheSongRoom)
	http.HandleFunc("/GuessTheSongRoomHandler", GuessTheSongRoomHandler)
	http.HandleFunc("/PublicGuessTheSong", JoinPublicGuessTheSong)
	http.HandleFunc("/CreateGuessTheSong", CreateGuessTheSong)
	http.HandleFunc("/CreateGuessTheSongHandler", CreateGuessTheSongHandler)
	http.HandleFunc("/JoinGuessTheSong", JoinPrivateGuessTheSong)
	http.HandleFunc("/JoinGuessTheSongHandler", JoinGuessTheSongHandler)
	http.HandleFunc("/GuessTheSongRules", GuessTheSongRules)
	http.HandleFunc("/PetitBacLandingPage", PetitBacLandingPage)
	http.HandleFunc("/PetitBac", PetitBac)
	http.HandleFunc("/PetitBacHandler", PetitBacHandler)
	http.HandleFunc("/PetitBacRoom", PetitBacRoom)
	http.HandleFunc("/PetitBacRoomHandler", PetitBacRoomHandler)
	http.HandleFunc("/PublicPetitBac", JoinPublicPetitBac)
	http.HandleFunc("/CreatePetitBac", CreatePetitBac)
	http.HandleFunc("/CreatePetitBacHandler", CreatePetitBacHandler)
	http.HandleFunc("/JoinPetitBac", JoinPrivatePetitBac)
	http.HandleFunc("/JoinPetitBacHandler", JoinPetitBacHandler)
	http.HandleFunc("/PetitBacRules", PetitBacRules)
	http.HandleFunc("/Result", Result)
	go handleMessages()
	http.HandleFunc("/websocket", websocketHandler)
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", nil)
}

func renderTemplate(w http.ResponseWriter, templatePath string, data interface{}) {
	tmpl, err := template.ParseFiles("./html/" + templatePath)
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

	roomData := database.GetRoomData(roomID)

	//Online chat
	if r.Method == "POST" {
		action := r.FormValue("action")
		if action == "ChatMessage" {
			cookie, err := r.Cookie("session_id")
			if err != nil {
				http.Redirect(w, r, "/Login", http.StatusSeeOther)
				return
			}

			userID, _ := strconv.Atoi(cookie.Value)
			userdata, _ := database.GetUserData(strconv.Itoa(userID))

			message := r.Form.Get("Message")
			ChatDiscours = message
			messages = append(messages, Message{
				Username:    userdata.Pseudo,
				TextMessage: ChatDiscours,
			})

			jsonMessage, err := json.Marshal(messages[len(messages)-1])
			if err != nil {
				fmt.Println("Error marshaling message to JSON:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			broadcast <- Message{
				Username:    userdata.Pseudo,
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
		currentTrack, indexTrack := games.NextTrack(tracks)

		if currentTrack == nil {
			http.Redirect(w, r, "/Result", http.StatusSeeOther)
		}
		games.PlayedTracks = append(games.PlayedTracks, currentTrack)

		gameData := database.GetRoomData(roomID)

		if games.BlindtestData.NumberOfTurn == -1 {
			games.BlindtestData.NumberOfTurn = gameData.NumberOfGameTurns
		} else {
			games.BlindtestData.NumberOfTurn--
		}
		if games.BlindtestData.NumberOfTurn == 0 {
			http.Redirect(w, r, "/Result", http.StatusSeeOther)
		}

		trackData := games.PageData{
			Track: currentTrack,
		}

		mediasBlindtest := Data{
			DatasgameBlindTest: trackData,
			Info:               messages,
		}

		userData, _ := database.GetUserData(strconv.Itoa(userID))
		userScore := database.GetUserScore(roomID, userID)

		if r.Method == "POST" {
			response := r.FormValue("BlindtestResponse")
			if games.NormalizeString(response) == games.NormalizeString(games.PlayedTracks[len(games.PlayedTracks)-2].Name) {
				userScore++
				database.SetUserScore(roomID, userID, userScore)
			}
		}

		//Data
		data := struct {
			MediasBlindtest  Data
			Username         string
			UserScore        int
			DurationOfMusic  int
			DurationOfAnswer int
			ActualTurn       int
			NumberOfTurns    int
		}{
			MediasBlindtest:  mediasBlindtest,
			Username:         userData.Pseudo,
			UserScore:        userScore,
			DurationOfMusic:  roomData.BlindtestTimeOfMusic,
			DurationOfAnswer: roomData.BlindtestTimeToAnswer,
			ActualTurn:       gameData.NumberOfGameTurns - games.BlindtestData.NumberOfTurn + 1,
			NumberOfTurns:    gameData.NumberOfGameTurns,
		}

		gameData.BlindtestTrackIndex = indexTrack
		gameData.BlindtestTimeOfMusic = data.DurationOfMusic
		gameData.BlindtestTimeToAnswer = data.DurationOfAnswer
		database.UpdateGameData(gameData, roomID)

		//Execute html
		renderTemplate(w, "BlindTest/index.html", data)

	} else {
		time.Sleep(50 * time.Millisecond)
		gameData := database.GetRoomData(roomID)
		w.Header().Set("Refresh", strconv.Itoa(gameData.BlindtestTimeOfMusic+gameData.BlindtestTimeToAnswer))

		tracks := games.Api("6Xf0gjt1YmwvEG5iS8QOfg?si=2de553d01ff84abb")
		tracks = games.RemovePlayedTracks(tracks)
		// currentTrack := tracks[gameData.BlindtestTrackIndex]
		currentTrack, _ := games.NextTrack(tracks)
		if currentTrack == nil {
			http.Redirect(w, r, "/Result", http.StatusSeeOther)
		}
		games.PlayedTracks = append(games.PlayedTracks, currentTrack)

		trackData := games.PageData{
			Track: currentTrack,
		}

		mediasBlindtest := Data{
			DatasgameBlindTest: trackData,
			Info:               messages,
		}

		userData, _ := database.GetUserData(strconv.Itoa(userID))
		userScore := database.GetUserScore(roomID, userID)

		if r.Method == "POST" {
			response := r.FormValue("BlindtestResponse")
			if games.NormalizeString(response) == games.NormalizeString(games.PlayedTracks[len(games.PlayedTracks)-2].Name) {
				userScore++
				database.SetUserScore(roomID, userID, userScore)
			}
		}

		data := struct {
			MediasBlindtest  Data
			Username         string
			UserScore        int
			DurationOfMusic  int
			DurationOfAnswer int
			ActualTurn       int
			NumberOfTurns    int
		}{
			MediasBlindtest:  mediasBlindtest,
			Username:         userData.Pseudo,
			UserScore:        userScore,
			DurationOfMusic:  gameData.BlindtestTimeOfMusic,
			DurationOfAnswer: gameData.BlindtestTimeToAnswer,
			ActualTurn:       gameData.NumberOfGameTurns - games.BlindtestData.NumberOfTurn + 1,
			NumberOfTurns:    gameData.NumberOfGameTurns,
		}
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

	games.BlindtestData.NumberOfTurn = -1

	if userID == creatorID {
		data.ButtonVisible = true
	} else {
		for !database.GetGameState(roomID) {
			time.Sleep(1 * time.Second)
		}

		http.Redirect(w, r, "/Blindtest", http.StatusSeeOther)
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
	database.SetGameState(roomID)
	http.Redirect(w, r, "/Blindtest", http.StatusSeeOther)
}

func JoinPublicBlindtest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	userID, _ := strconv.Atoi(cookie.Value)
	roomID := 1
	games.BlindtestData.NumberOfTurn = -1

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

		_, err = gameRoomDB.Exec("INSERT INTO GAME_ROOM (id_room, game_state, number_of_game_turns, blindtest_time_of_music, blindtest_time_to_answer) VALUES (?, ?, ?, ?, ?)", roomID, true, numberOfGameTurns, timeOfMusic, timeToAnswer)
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

	renderTemplate(w, "GuessTheSong/LandingPage.html", nil)
}

func GuessTheSong(w http.ResponseWriter, r *http.Request) {
	//Recovering the room and its data
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	userID, _ := strconv.Atoi(cookie.Value)

	roomID, err := database.GetRoomIDByUserID(userID)
	if err != nil {
		http.Redirect(w, r, "/JoinGuessTheSong", http.StatusSeeOther)
	}

	roomData := database.GetRoomData(roomID)

	//Game Chat
	if r.Method == "POST" {
		action := r.FormValue("action")
		if action == "ChatMessage" {
			cookie, err := r.Cookie("session_id")
			if err != nil {
				http.Redirect(w, r, "/Login", http.StatusSeeOther)
				return
			}

			userID, _ := strconv.Atoi(cookie.Value)
			userdata, _ := database.GetUserData(strconv.Itoa(userID))

			message := r.Form.Get("Message")
			ChatDiscours = message
			messages = append(messages, Message{
				Username:    userdata.Pseudo,
				TextMessage: ChatDiscours,
			})

			jsonMessage, err := json.Marshal(messages[len(messages)-1])
			if err != nil {
				fmt.Println("Error marshaling message to JSON:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			broadcast <- Message{
				Username:    userdata.Pseudo,
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

	if userID == database.GetRoomCreator(roomID) {
		w.Header().Set("Refresh", strconv.Itoa(roomData.GuessthesongTimeToAnswer))

		games.LoadData()

		medias := Data{
			Datasgames: games.CurrentSong,
			Info:       messages,
		}

		userData, _ := database.GetUserData(strconv.Itoa(userID))
		userScore := database.GetUserScore(roomID, userID)

		if r.Method == "POST" {
			input := r.FormValue("value")
			if games.CompareStrings(games.NormalizeString(input), games.NormalizeString(games.CurrentSong.TitleSong)) {
				userScore++
				database.SetUserScore(roomID, userID, userScore)
			}
		}

		gameData := database.GetRoomData(roomID)

		if games.GuessTheSongData.NumberOfTurn == -1 {
			games.GuessTheSongData.NumberOfTurn = gameData.NumberOfGameTurns
		} else {
			games.GuessTheSongData.NumberOfTurn--
		}
		if games.GuessTheSongData.NumberOfTurn == 0 {
			http.Redirect(w, r, "/Result", http.StatusSeeOther)
		}

		//Data
		data := struct {
			MediasGuessthesong Data
			Username           string
			UserScore          int
			DurationOfMusic    int
			DurationOfAnswer   int
			ActualTurn         int
			NumberOfTurns      int
		}{
			MediasGuessthesong: medias,
			Username:           userData.Pseudo,
			UserScore:          userScore,
			DurationOfMusic:    roomData.GuessthesongTimeToAnswer,
			DurationOfAnswer:   roomData.GuessthesongTimeToAnswer,
			ActualTurn:         gameData.NumberOfGameTurns - games.GuessTheSongData.NumberOfTurn + 1,
			NumberOfTurns:      gameData.NumberOfGameTurns,
		}

		gameData.GuessthesongTimeToAnswer = data.DurationOfAnswer
		database.UpdateGameData(gameData, roomID)

		//Execute html
		renderTemplate(w, "GuessTheSong/index.html", data)
	} else {
		time.Sleep(50 * time.Millisecond)
		gameData := database.GetRoomData(roomID)
		w.Header().Set("Refresh", strconv.Itoa(gameData.GuessthesongTimeToAnswer))

		games.LoadData()

		medias := Data{
			Datasgames: games.CurrentSong,
			Info:       messages,
		}

		userData, _ := database.GetUserData(strconv.Itoa(userID))
		userScore := database.GetUserScore(roomID, userID)

		if r.Method == "POST" {
			input := r.FormValue("value")
			if games.CompareStrings(games.NormalizeString(input), games.NormalizeString(games.CurrentSong.TitleSong)) {
				userScore++
				database.SetUserScore(roomID, userID, userScore)
			}
		}

		if games.GuessTheSongData.NumberOfTurn == -1 {
			games.GuessTheSongData.NumberOfTurn = gameData.NumberOfGameTurns
		} else {
			games.GuessTheSongData.NumberOfTurn--
		}
		if games.GuessTheSongData.NumberOfTurn == 0 {
			http.Redirect(w, r, "/Result", http.StatusSeeOther)
		}

		//Data
		data := struct {
			MediasGuessthesong Data
			Username           string
			UserScore          int
			DurationOfMusic    int
			DurationOfAnswer   int
			ActualTurn         int
			NumberOfTurns      int
		}{
			MediasGuessthesong: medias,
			Username:           userData.Pseudo,
			UserScore:          userScore,
			DurationOfMusic:    roomData.BlindtestTimeOfMusic,
			DurationOfAnswer:   roomData.BlindtestTimeToAnswer,
			ActualTurn:         gameData.NumberOfGameTurns - games.GuessTheSongData.NumberOfTurn + 1,
			NumberOfTurns:      gameData.NumberOfGameTurns,
		}

		gameData.GuessthesongTimeToAnswer = data.DurationOfAnswer
		database.UpdateGameData(gameData, roomID)

		//Execute html
		renderTemplate(w, "GuessTheSong/index.html", data)
	}
}

func GuessTheSongRoom(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	userID, _ := strconv.Atoi(cookie.Value)
	roomID, err := database.GetRoomIDByUserID(userID)
	if err != nil {
		http.Redirect(w, r, "/JoinGuessTheSong", http.StatusSeeOther)
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

	games.GuessTheSongData.NumberOfTurn = -1

	if userID == creatorID {
		data.ButtonVisible = true
	} else {
		for !database.GetGameState(roomID) {
			time.Sleep(1 * time.Second)
		}

		http.Redirect(w, r, "/GuessTheSong", http.StatusSeeOther)
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
	renderTemplate(w, "GuessTheSong/Room.html", data)
}

func GuessTheSongRoomHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	userID, _ := strconv.Atoi(cookie.Value)
	roomID, err := database.GetRoomIDByUserID(userID)
	if err != nil {
		http.Redirect(w, r, "/JoinGuessTheSong", http.StatusSeeOther)
	}
	database.SetGameState(roomID)
	http.Redirect(w, r, "/GuessTheSong", http.StatusSeeOther)
}

func JoinPublicGuessTheSong(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	userID, _ := strconv.Atoi(cookie.Value)
	roomID := 2
	games.GuessTheSongData.NumberOfTurn = -1

	if !database.CheckRoomExistence(roomID) {
		db := database.InitTable("ROOMS")
		defer db.Close()

		createdBy := userID
		maxPlayers := 50
		name := "publicGuessTheSong"
		gameID := 2

		// Create the public room
		_, err := db.Exec("INSERT INTO ROOMS (id, created_by, max_player, number_of_player, name, id_game) VALUES (?, ?, ?, ?, ?, ?)", roomID, createdBy, maxPlayers, 0, name, gameID)
		if err != nil {
			log.Println("Error creating room:", err)
		}

		//Room Data
		gameRoomDB := database.InitTable("GAME_ROOM")
		defer gameRoomDB.Close()

		numberOfGameTurns := 10
		difficulty := "facile"
		timeToAnswer := 30

		_, err = gameRoomDB.Exec("INSERT INTO GAME_ROOM (id_room, game_state, number_of_game_turns, guessthesong_difficulty, guessthesong_time_to_answer) VALUES (?, ?, ?, ?, ?)", roomID, true, numberOfGameTurns, difficulty, timeToAnswer)
		if err != nil {
			log.Println("Error creating game room:", err)
			return
		}
	}

	database.JoinRoom(roomID, userID)

	http.Redirect(w, r, "/GuessTheSong", http.StatusSeeOther)
}

func CreateGuessTheSong(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "GuessTheSong/CreatePrivateRoom.html", nil)
}

func CreateGuessTheSongHandler(w http.ResponseWriter, r *http.Request) {
	gameTurns, _ := strconv.Atoi(r.FormValue("gameTurns"))
	difficulty := r.FormValue("difficulty")
	answerDuration, _ := strconv.Atoi(r.FormValue("answerDuration"))
	roomName := r.FormValue("roomName")
	maxPlayer, _ := strconv.Atoi(r.FormValue("maxPlayer"))

	if gameTurns <= 0 {
		gameTurns = 5
	}
	if answerDuration <= 0 {
		answerDuration = 30
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	sessionID := cookie.Value

	userID, _ := strconv.Atoi(sessionID)
	gameID := 2

	database.CreateGuessthesongRoom(userID, maxPlayer, roomName, gameID, gameTurns, difficulty, answerDuration)

	http.Redirect(w, r, "/GuessTheSongRoom", http.StatusSeeOther)
}

func JoinPrivateGuessTheSong(w http.ResponseWriter, r *http.Request) {
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
	}

	renderTemplate(w, "GuessTheSong/JoinPrivateRoom.html", nil)
}

func JoinGuessTheSongHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}
	userID, _ := strconv.Atoi(cookie.Value)
	roomID, err2 := strconv.Atoi(r.FormValue("roomID"))
	if err2 != nil || roomID <= 0 {
		data := struct{ ErrorMessage string }{ErrorMessage: "Veuillez rentrez un identifiant valide."}
		renderTemplate(w, "GuessTheSong/JoinPrivateRoom.html", data)
	} else if !database.VerifyRoom(roomID) {
		data := struct{ ErrorMessage string }{ErrorMessage: "Aucune room ne correspond à cet identifiant."}
		renderTemplate(w, "GuessTheSong/JoinPrivateRoom.html", data)
	} else {
		database.JoinRoom(roomID, userID)
		http.Redirect(w, r, "/GuessTheSongRoom", http.StatusSeeOther)
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

	renderTemplate(w, "PetitBac/LandingPage.html", nil)
}

func PetitBac(w http.ResponseWriter, r *http.Request) {

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

	roomData := database.GetRoomData(roomID)

	//Game Chat
	if r.Method == "POST" {
		action := r.FormValue("action")
		if action == "ChatMessage" {
			cookie, err := r.Cookie("session_id")
			if err != nil {
				http.Redirect(w, r, "/Login", http.StatusSeeOther)
				return
			}

			userID, _ := strconv.Atoi(cookie.Value)
			userdata, _ := database.GetUserData(strconv.Itoa(userID))

			log.Println("put your message")
			message := r.Form.Get("Message")
			ChatDiscours = message
			messages = append(messages, Message{
				Username:    userdata.Pseudo,
				TextMessage: ChatDiscours,
			})

			jsonMessage, err := json.Marshal(messages[len(messages)-1])
			if err != nil {
				fmt.Println("Error marshaling message to JSON:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			broadcast <- Message{
				Username:    userdata.Pseudo,
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
		w.Header().Set("Refresh", strconv.Itoa(roomData.PetitbacTimeToAnswer))
		letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		randomLetter := string(letters[rand.Intn(len(letters))])

		gameData := database.GetRoomData(roomID)

		if games.PetitBacData.NumberOfTurn == -1 {
			games.PetitBacData.NumberOfTurn = gameData.NumberOfGameTurns
		} else {
			games.PetitBacData.NumberOfTurn--
		}
		if games.PetitBacData.NumberOfTurn == 0 {
			http.Redirect(w, r, "/Result", http.StatusSeeOther)
		}

		userData, _ := database.GetUserData(strconv.Itoa(userID))
		userScore := database.GetUserScore(roomID, userID)

		if r.Method == "POST" {
			response := r.FormValue("action")
			if response == "true" {
				r.ParseForm()
				fields := []string{"artiste", "album", "groupe", "instrument", "featuring"}
				for _, field := range fields {
					value := r.Form.Get(field)
					if strings.HasPrefix(value, randomLetter) {
						userScore++
						database.SetUserScore(roomID, userID, userScore)
					}
				}
			}
		}

		data := struct {
			RandomLetter     string
			Username         string
			UserScore        int
			DurationOfAnswer int
			ActualTurn       int
			NumberOfTurns    int
			Info             []Message
		}{
			RandomLetter:     randomLetter,
			Username:         userData.Pseudo,
			UserScore:        userScore,
			DurationOfAnswer: roomData.PetitbacTimeToAnswer,
			ActualTurn:       gameData.NumberOfGameTurns - games.PetitBacData.NumberOfTurn + 1,
			NumberOfTurns:    gameData.NumberOfGameTurns,
			Info:             messages,
		}

		gameData.PetitbacTimeToAnswer = data.DurationOfAnswer
		database.UpdateGameData(gameData, roomID)

		//Execute html
		renderTemplate(w, "PetitBac/index.html", data)
	} else {
		time.Sleep(50 * time.Millisecond)
		gameData := database.GetRoomData(roomID)
		w.Header().Set("Refresh", strconv.Itoa(gameData.PetitbacTimeToAnswer))
		letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		randomLetter := string(letters[rand.Intn(len(letters))])
		fmt.Println(gameData)

		if games.PetitBacData.NumberOfTurn == -1 {
			games.PetitBacData.NumberOfTurn = gameData.NumberOfGameTurns
		} else {
			games.PetitBacData.NumberOfTurn--
		}
		if games.PetitBacData.NumberOfTurn == 0 {
			http.Redirect(w, r, "/Result", http.StatusSeeOther)
		}

		userData, _ := database.GetUserData(strconv.Itoa(userID))
		userScore := database.GetUserScore(roomID, userID)

		if r.Method == "POST" {
			response := r.FormValue("action")
			if response == "true" {
				r.ParseForm()
				fields := []string{"artiste", "album", "groupe", "instrument", "featuring"}
				for _, field := range fields {
					value := r.Form.Get(field)
					if strings.HasPrefix(value, randomLetter) {
						userScore++
						database.SetUserScore(roomID, userID, userScore)
					}
				}
			}
		}

		data := struct {
			RandomLetter     string
			Username         string
			UserScore        int
			DurationOfAnswer int
			ActualTurn       int
			NumberOfTurns    int
			Info             []Message
		}{
			RandomLetter:     randomLetter,
			Username:         userData.Pseudo,
			UserScore:        userScore,
			DurationOfAnswer: roomData.PetitbacTimeToAnswer,
			ActualTurn:       gameData.NumberOfGameTurns - games.PetitBacData.NumberOfTurn + 1,
			NumberOfTurns:    gameData.NumberOfGameTurns,
			Info:             messages,
		}

		gameData.PetitbacTimeToAnswer = data.DurationOfAnswer
		database.UpdateGameData(gameData, roomID)

		//Execute html
		renderTemplate(w, "PetitBac/index.html", data)
	}
}

func PetitBacHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/PetitBac", http.StatusSeeOther)
}

func PetitBacRoom(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	userID, _ := strconv.Atoi(cookie.Value)
	roomID, err := database.GetRoomIDByUserID(userID)
	if err != nil {
		http.Redirect(w, r, "/JoinPetitBac", http.StatusSeeOther)
	}

	creatorID := database.GetRoomCreator(roomID)
	numberOfPlayer := database.GetNumberOfPlayer(roomID)
	games.PetitBacData.NumberOfTurn = -1

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
	} else {
		for !database.GetGameState(roomID) {
			time.Sleep(1 * time.Second)
		}

		http.Redirect(w, r, "/PetitBac", http.StatusSeeOther)
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
	renderTemplate(w, "PetitBac/Room.html", data)
}

func PetitBacRoomHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	userID, _ := strconv.Atoi(cookie.Value)
	roomID, err := database.GetRoomIDByUserID(userID)
	if err != nil {
		http.Redirect(w, r, "/JoinPetitBac", http.StatusSeeOther)
	}
	database.SetGameState(roomID)
	http.Redirect(w, r, "/PetitBac", http.StatusSeeOther)
}

func JoinPublicPetitBac(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	userID, _ := strconv.Atoi(cookie.Value)
	roomID := 3
	games.PetitBacData.NumberOfTurn = -1

	if !database.CheckRoomExistence(roomID) {
		db := database.InitTable("ROOMS")
		defer db.Close()

		createdBy := userID
		maxPlayers := 50
		name := "publicPetitBac"
		gameID := 3

		// Create the public room
		_, err := db.Exec("INSERT INTO ROOMS (id, created_by, max_player, number_of_player, name, id_game) VALUES (?, ?, ?, ?, ?, ?)", roomID, createdBy, maxPlayers, 0, name, gameID)
		if err != nil {
			log.Println("Error creating room:", err)
		}

		//Room Data
		gameRoomDB := database.InitTable("GAME_ROOM")
		defer gameRoomDB.Close()

		numberOfGameTurns := 10
		categories := "album, groupe de musique, instrument de musique, featuring"
		timeToAnswer := 60

		_, err = gameRoomDB.Exec("INSERT INTO GAME_ROOM (id_room, game_state, number_of_game_turns, petitbac_categories, petitbac_time_to_answer) VALUES (?, ?, ?, ?, ?)", roomID, true, numberOfGameTurns, categories, timeToAnswer)
		if err != nil {
			log.Println("Error creating game room:", err)
			return
		}
	}

	database.JoinRoom(roomID, userID)

	http.Redirect(w, r, "/PetitBac", http.StatusSeeOther)
}

func CreatePetitBac(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "PetitBac/CreatePrivateRoom.html", nil)
}

func CreatePetitBacHandler(w http.ResponseWriter, r *http.Request) {
	gameTurns, _ := strconv.Atoi(r.FormValue("gameTurns"))
	answerDuration, _ := strconv.Atoi(r.FormValue("answerDuration"))
	roomName := r.FormValue("roomName")
	maxPlayer, _ := strconv.Atoi(r.FormValue("maxPlayer"))
	categories := r.Form["categories[]"]

	for i, category := range categories {
		categories[i] = strings.ToLower(category)
	}

	if gameTurns <= 0 {
		gameTurns = 5
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
	gameID := 3

	database.CreatePetitbacRoom(userID, maxPlayer, roomName, gameID, gameTurns, strings.Join(categories, ", "), answerDuration)

	http.Redirect(w, r, "/PetitBacRoom", http.StatusSeeOther)
}

func JoinPrivatePetitBac(w http.ResponseWriter, r *http.Request) {
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
	}

	renderTemplate(w, "PetitBac/JoinPrivateRoom.html", nil)
}

func JoinPetitBacHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}
	userID, _ := strconv.Atoi(cookie.Value)
	roomID, err2 := strconv.Atoi(r.FormValue("roomID"))
	if err2 != nil || roomID <= 0 {
		data := struct{ ErrorMessage string }{ErrorMessage: "Veuillez rentrez un identifiant valide."}
		renderTemplate(w, "PetitBac/JoinPrivateRoom.html", data)
	} else if !database.VerifyRoom(roomID) {
		data := struct{ ErrorMessage string }{ErrorMessage: "Aucune room ne correspond à cet identifiant."}
		renderTemplate(w, "PetitBac/JoinPrivateRoom.html", data)
	} else {
		database.JoinRoom(roomID, userID)
		http.Redirect(w, r, "/PetitBacRoom", http.StatusSeeOther)
	}
}

func PetitBacRules(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "PetitBac/Rules.html", nil)
}

type Classement struct {
	Rank   int
	Name   string
	Points int
}

func Result(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}
	userID, _ := strconv.Atoi(cookie.Value)

	roomID, _ := database.GetRoomIDByUserID(userID)
	userIDs := database.GetAllUserInRoom(roomID)

	var scoreBoard []Classement

	for _, userID := range userIDs {
		userScore := database.GetUserScore(roomID, userID)
		userData, _ := database.GetUserData(strconv.Itoa(userID))
		username := userData.Pseudo
		rank := rand.Intn(2) + 1
		scoreBoard = append(scoreBoard, Classement{
			Rank:   rank,
			Name:   username,
			Points: userScore,
		})
	}

	data := struct {
		ScoreBoard []Classement
	}{
		ScoreBoard: scoreBoard,
	}

	renderTemplate(w, "BlindTest/Result.html", data)
}
