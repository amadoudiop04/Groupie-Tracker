package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

// ---------------------------ROOMS---------------------------\\
type Room struct {
	Id         int
	CreatedBy  int
	MaxPlayers int
	Name       string
	GameID     int
}

func CreateRoom(createdBy int, maxPlayers int, name string, gameID int) int {
	db := InitTable("ROOMS")
	defer db.Close()

	// Create the good ID for the room
	var roomID int
	for i := 0; i < 100; i++ {
		roomID = (gameID * 100) + i
		row := db.QueryRow("SELECT COUNT(*) FROM ROOMS WHERE id = ?", roomID)
		var count int
		err := row.Scan(&count)
		if err != nil {
			log.Println("Error checking room ID:", err)
		}
		if count == 0 {
			break
		}
	}

	// Créer la room avec l'ID déterminé
	_, err := db.Exec("INSERT INTO ROOMS (id, created_by, max_player, number_of_player, name, id_game) VALUES (?, ?, ?, ?, ?, ?)", roomID, createdBy, maxPlayers, 0, name, gameID)
	if err != nil {
		log.Println("Error creating room:", err)
	}

	JoinRoom(roomID, createdBy)

	return roomID
}

func CreateBlindtestRoom(createdBy int, maxPlayers int, name string, gameID int, numberOfGameTurns int, timeOfMusic int, timeToAnswer int) {
	roomID := CreateRoom(createdBy, maxPlayers, name, gameID)

	// Ouvrir une connexion à la table GAME_ROOM
	gameRoomDB := InitTable("GAME_ROOM")
	defer gameRoomDB.Close()

	// Insérer les données dans la table GAME_ROOM
	_, err := gameRoomDB.Exec("INSERT INTO GAME_ROOM (id_room, number_of_game_turns, blindtest_time_of_music, blindtest_time_to_answer) VALUES (?, ?, ?, ?)", roomID, numberOfGameTurns, timeOfMusic, timeToAnswer)
	if err != nil {
		log.Println("Error creating game room:", err)
		return
	}
}

func CreatePetitbacRoom(createdBy int, maxPlayers int, name string, gameID int, numberOfGameTurns int, categories string, timeToAnswer int) {
	roomID := CreateRoom(createdBy, maxPlayers, name, gameID)

	// Ouvrir une connexion à la table GAME_ROOM
	gameRoomDB := InitTable("GAME_ROOM")
	defer gameRoomDB.Close()

	// Insérer les données dans la table GAME_ROOM
	_, err := gameRoomDB.Exec("INSERT INTO GAME_ROOM (id_room, number_of_game_turns, petitbac_categories, petitbac_time_to_answer) VALUES (?, ?, ?, ?)", roomID, numberOfGameTurns, categories, timeToAnswer)
	if err != nil {
		log.Println("Error creating game room:", err)
		return
	}
}

func CreateGuessthesongRoom(db *sql.DB, createdBy int, maxPlayers int, name string, gameID int, numberOfGameTurns int, difficulty string, timeToAnswer int) {
	roomID := CreateRoom(createdBy, maxPlayers, name, gameID)

	// Ouvrir une connexion à la table GAME_ROOM
	gameRoomDB := InitTable("GAME_ROOM")
	defer gameRoomDB.Close()

	// Insérer les données dans la table GAME_ROOM
	_, err := gameRoomDB.Exec("INSERT INTO GAME_ROOM (id_room, number_of_game_turns, guessthesong_difficulty, guessthesong_time_to_answer) VALUES (?, ?, ?, ?)", roomID, numberOfGameTurns, difficulty, timeToAnswer)
	if err != nil {
		log.Println("Error creating game room:", err)
		return
	}
}

func JoinRoom(roomID int, userID int) error {
	//ROOM_USERS
	db := InitTable("ROOM_USERS")
	defer db.Close()

	_, err := db.Exec("INSERT INTO ROOM_USERS (id_room, id_user) VALUES (?, ?)", roomID, userID)
	if err != nil {
		log.Println("Error joining room:", err)
		return err
	}

	//ROOMS
	roomDB := InitTable("ROOMS")
	defer roomDB.Close()

	var playerCount int
	row := roomDB.QueryRow("SELECT number_of_player FROM ROOMS WHERE id = ?", roomID)
	err = row.Scan(&playerCount)
	if err != nil {
		log.Println("Error counting players in room:", err)
		return err
	}

	playerCount++

	_, err = roomDB.Exec("UPDATE ROOMS SET number_of_player = ? WHERE id = ?", playerCount, roomID)
	if err != nil {
		log.Println("Error joining room:", err)
		return err
	}
	return nil
}

func LeaveRoom(roomID int, userID int) error {
	db := InitTable("ROOM_USERS")
	defer db.Close()

	_, err := db.Exec("DELETE FROM ROOM_USERS WHERE id_room = ? AND id_user = ?", roomID, userID)
	if err != nil {
		log.Println("Error leaving room:", err)
	}
	return err
}

func VerifyRoom(roomID int) bool {
	db := InitTable("ROOMS")
	defer db.Close()

	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM ROOMS WHERE id = ?", roomID)
	err := row.Scan(&count)
	if err != nil {
		log.Println("Error verifying room:", err)
		return false
	}

	return count > 0
}

func CheckRoomCapacity(db *sql.DB, roomID int) (bool, error) {
	var currentPlayers int
	err := db.QueryRow("SELECT COUNT(id_user) FROM ROOM_USERS WHERE id_room = ?", roomID).Scan(&currentPlayers)
	if err != nil {
		log.Println("Error checking room capacity:", err)
		return false, err
	}
	var maxPlayers int
	err = db.QueryRow("SELECT max_player FROM ROOMS WHERE id = ?", roomID).Scan(&maxPlayers)
	if err != nil {
		log.Println("Error checking room capacity:", err)
		return false, err
	}
	return currentPlayers < maxPlayers, nil
}

func GetRoomDetails(db *sql.DB, roomID int) (Room, error) {
	var room Room
	err := db.QueryRow("SELECT id, created_by, max_player, name, id_game FROM ROOMS WHERE id = ?", roomID).
		Scan(&room.Id, &room.CreatedBy, &room.MaxPlayers, &room.Name, &room.GameID)
	if err != nil {
		log.Println("Error fetching room details:", err)
		return room, err
	}
	return room, nil
}

func GetRoomIDByUserID(userID int) int {
	db := InitTable("ROOM_USERS")
	defer db.Close()

	var roomID int
	err := db.QueryRow("SELECT id_room FROM ROOM_USERS WHERE id_user = ?", userID).Scan(&roomID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No room found for the user")
		}
		log.Fatal(err)
	}
	return roomID
}

func GetRoomCreator(roomID int) int {
	db := InitTable("ROOMS")
	defer db.Close()

	var creatorID int
	err := db.QueryRow("SELECT created_by FROM ROOMS WHERE id = ?", roomID).Scan(&creatorID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No room found for the user")
		}
		log.Fatal(err)
	}
	return creatorID
}

func CheckRoomExistence(roomID int) bool {
	db := InitTable("ROOMS")
	defer db.Close()

	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM ROOMS WHERE id = ?", roomID)
	err := row.Scan(&count)
	if err != nil {
		log.Println("Error checking room existence:", err)
		return false
	}

	return count > 0
}

func GetNumberOfPlayer(roomID int) int {
	db := InitTable("ROOMS")
	defer db.Close()

	var numberOfPlayer int
	err := db.QueryRow("SELECT number_of_player FROM ROOMS WHERE id = ?", roomID).Scan(&numberOfPlayer)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No room found for the user")
		}
		log.Fatal(err)
	}
	return numberOfPlayer
}

// ---------------------------ROOM_USERS---------------------------\\
type Room_Users struct {
	RoomId int
	UserId int
	Score  int
}

func AddUserToRoom(db *sql.DB, roomID, userID int, initialScore int) error {
	_, err := db.Exec("INSERT INTO ROOM_USERS (id_room, id_user, score) VALUES (?, ?, ?)", roomID, userID, initialScore)
	if err != nil {
		log.Println("Erreur lors de l'ajout de l'utilisateur à la salle de jeu :", err)
		return err
	}
	return nil
}

func RemoveUserFromRoom(db *sql.DB, roomID, userID int) error {
	_, err := db.Exec("DELETE FROM ROOM_USERS WHERE id_room = ? AND id_user = ?", roomID, userID)
	if err != nil {
		log.Println("Erreur lors de la suppression de l'utilisateur de la salle de jeu :", err)
		return err
	}
	return nil
}

func UpdateUserScore(db *sql.DB, roomID, userID int, newScore int) error {
	_, err := db.Exec("UPDATE ROOM_USERS SET score = ? WHERE id_room = ? AND id_user = ?", newScore, roomID, userID)
	if err != nil {
		log.Println("Erreur lors de la mise à jour du score de l'utilisateur dans la salle de jeu :", err)
		return err
	}
	return nil
}

func GetRoomUsers(db *sql.DB, roomID int) (map[int]int, error) {
	users := make(map[int]int)

	rows, err := db.Query("SELECT id_user, score FROM ROOM_USERS WHERE id_room = ?", roomID)
	if err != nil {
		log.Println("Erreur lors de la récupération des utilisateurs de la salle de jeu :", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID, score int
		err := rows.Scan(&userID, &score)
		if err != nil {
			log.Println("Erreur lors de la lecture des données des utilisateurs de la salle de jeu :", err)
			return nil, err
		}
		users[userID] = score
	}

	if err := rows.Err(); err != nil {
		log.Println("Erreur lors de l'itération sur les résultats de la requête pour récupérer les utilisateurs de la salle de jeu :", err)
		return nil, err
	}

	return users, nil
}

// ---------------------------GAMES---------------------------\\
type Game struct {
	ID   int
	Name string
}

func CreateGame(db *sql.DB, name string) error {
	_, err := db.Exec("INSERT INTO GAMES (name) VALUES (?)", name)
	if err != nil {
		log.Println("Error creating game:", err)
		return err
	}
	return nil
}

func GetGameByID(db *sql.DB, id int) (*Game, error) {
	var game Game
	err := db.QueryRow("SELECT id, name FROM GAMES WHERE id = ?", id).Scan(&game.ID, &game.Name)
	if err != nil {
		log.Println("Error getting game by ID:", err)
		return nil, err
	}
	return &game, nil
}

func GetAllGames(db *sql.DB) ([]*Game, error) {
	rows, err := db.Query("SELECT id, name FROM GAMES")
	if err != nil {
		log.Println("Error getting all games:", err)
		return nil, err
	}
	defer rows.Close()

	var games []*Game
	for rows.Next() {
		var game Game
		if err := rows.Scan(&game.ID, &game.Name); err != nil {
			log.Println("Error scanning game rows:", err)
			return nil, err
		}
		games = append(games, &game)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over game rows:", err)
		return nil, err
	}
	return games, nil
}

// ---------------------------GAME_ROOM---------------------------\\

type GameRoomData struct {
	NumberOfGameTurns        int
	BlindtestTimeOfMusic     int
	BlindtestTimeToAnswer    int
	PetitbacCategories       string
	PetitbacTimeToAnswer     int
	GuessthesongDifficulty   string
	GuessthesongTimeToAnswer int
}

func GetRoomData(roomID int) (GameRoomData, error) {
	db := InitTable("GAME_ROOM")
	defer db.Close()

	var gameRoomData GameRoomData
	row := db.QueryRow("SELECT COALESCE(number_of_game_turns, 0), COALESCE(blindtest_time_of_music, 0), COALESCE(blindtest_time_to_answer, 0), COALESCE(petitbac_categories, ''), COALESCE(petitbac_time_to_answer, 0), COALESCE(guessthesong_difficulty, ''), COALESCE(guessthesong_time_to_answer, 0) FROM GAME_ROOM WHERE id_room = ?", roomID)
	err := row.Scan(&gameRoomData.NumberOfGameTurns, &gameRoomData.BlindtestTimeOfMusic, &gameRoomData.BlindtestTimeToAnswer, &gameRoomData.PetitbacCategories, &gameRoomData.PetitbacTimeToAnswer, &gameRoomData.GuessthesongDifficulty, &gameRoomData.GuessthesongTimeToAnswer)
	if err != nil {
		if err == sql.ErrNoRows {
			return GameRoomData{}, errors.New("No game room data found for the room")
		}
		return GameRoomData{}, err
	}
	return gameRoomData, nil
}
