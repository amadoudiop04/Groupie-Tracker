package database

import (
	"database/sql"
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
	_, err := gameRoomDB.Exec("INSERT INTO GAME_ROOM (id_room, game_state, number_of_game_turns, blindtest_time_of_music, blindtest_time_to_answer) VALUES (?, ?, ?, ?, ?)", roomID, false, numberOfGameTurns, timeOfMusic, timeToAnswer)
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
	_, err := gameRoomDB.Exec("INSERT INTO GAME_ROOM (id_room, game_state, number_of_game_turns, petitbac_categories, petitbac_time_to_answer) VALUES (?, ?, ?, ?, ?)", roomID, false, numberOfGameTurns, categories, timeToAnswer)
	if err != nil {
		log.Println("Error creating game room:", err)
		return
	}
}

func CreateGuessthesongRoom(createdBy int, maxPlayers int, name string, gameID int, numberOfGameTurns int, difficulty string, timeToAnswer int) {
	roomID := CreateRoom(createdBy, maxPlayers, name, gameID)

	// Ouvrir une connexion à la table GAME_ROOM
	gameRoomDB := InitTable("GAME_ROOM")
	defer gameRoomDB.Close()

	// Insérer les données dans la table GAME_ROOM
	_, err := gameRoomDB.Exec("INSERT INTO GAME_ROOM (id_room, game_state, number_of_game_turns, guessthesong_difficulty, guessthesong_time_to_answer) VALUES (?, ?, ?, ?, ?)", roomID, false, numberOfGameTurns, difficulty, timeToAnswer)
	if err != nil {
		log.Println("Error creating game room:", err)
		return
	}
}

func JoinRoom(roomID int, userID int) error {
	//ROOM_USERS
	db := InitTable("ROOM_USERS")
	defer db.Close()

	_, err := db.Exec("INSERT INTO ROOM_USERS (id_room, id_user, score) VALUES (?, ?, ?)", roomID, userID, 0)
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
		log.Println("Error update number of player in the room :", err)
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
		return err
	}

	roomDB := InitTable("ROOMS")
	defer roomDB.Close()
	gameRoomDB := InitTable("ROOMS")
	defer gameRoomDB.Close()

	var playerCount int
	row := roomDB.QueryRow("SELECT number_of_player FROM ROOMS WHERE id = ?", roomID)
	err = row.Scan(&playerCount)
	if err != nil {
		log.Println("Error counting players in room:", err)
		return err
	}

	playerCount--

	if playerCount == 0 {
		_, err = roomDB.Exec("DELETE FROM ROOMS WHERE id = ?", roomID)
		if err != nil {
			log.Println("Error update number of player in the room :", err)
			return err
		}
		_, err = gameRoomDB.Exec("DELETE FROM GAME_ROOM WHERE id_room = ?", roomID)
		if err != nil {
			log.Println("Error update number of player in the room :", err)
			return err
		}
	} else {
		_, err = roomDB.Exec("UPDATE ROOMS SET number_of_player = ? WHERE id = ?", playerCount, roomID)
		if err != nil {
			log.Println("Error update number of player in the room :", err)
			return err
		}
	}
	if userID == GetRoomCreator(roomID) {
		newCreator := GetAllUserInRoom(roomID)[0]
		SetRoomCreator(newCreator, roomID)
	}
	return nil
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

func GetRoomIDByUserID(userID int) (int, error) {
	db := InitTable("ROOM_USERS")
	defer db.Close()

	var roomID int
	err := db.QueryRow("SELECT id_room FROM ROOM_USERS WHERE id_user = ?", userID).Scan(&roomID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No room found for the user")
		}
		return 0, err
	}
	return roomID, nil
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
		fmt.Println(err)
	}
	return creatorID
}

func SetRoomCreator(roomID, userID int) {
	db := InitTable("ROOMS")
	defer db.Close()

	_, err := db.Exec("UPDATE ROOMS SET created_by = ? WHERE id = ?", userID, roomID)
	if err != nil {
		log.Println("Error setting game state:", err)
	}
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
		fmt.Println(err)
	}
	return numberOfPlayer
}

func GetGameState(roomID int) bool {
	db := InitTable("GAME_ROOM")
	defer db.Close()

	var gameState bool
	err := db.QueryRow("SELECT game_state FROM GAME_ROOM WHERE id_room = ?", roomID).Scan(&gameState)
	if err != nil {
		log.Println("Error getting game state:", err)
		return false
	}
	return gameState
}

func SetGameState(roomID int) {
	db := InitTable("GAME_ROOM")
	defer db.Close()

	gameState := GetGameState(roomID)

	_, err := db.Exec("UPDATE GAME_ROOM SET game_state = ? WHERE id_room = ?", !gameState, roomID)
	if err != nil {
		log.Println("Error setting game state:", err)
	}
}

func GetAllUserInRoom(roomID int) []int {
	db := InitTable("ROOM_USERS")
	defer db.Close()

	rows, err := db.Query("SELECT id_user FROM ROOM_USERS WHERE id_room = ?", roomID)
	if err != nil {
		log.Println("Error getting users in room:", err)
		return nil
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			log.Println("Error scanning user ID:", err)
			continue
		}
		userIDs = append(userIDs, userID)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over user rows:", err)
	}
	return userIDs
}

func GetUserScore(roomID, userID int) int {
	db := InitTable("ROOM_USERS")
	defer db.Close()

	var score int
	err := db.QueryRow("SELECT COALESCE(score, 0) FROM ROOM_USERS WHERE id_room = ? AND id_user = ?", roomID, userID).Scan(&score)
	if err != nil {
		log.Println("Error getting user score:", err)
	}
	return score
}

func SetUserScore(roomID, userID, score int) {
	db := InitTable("ROOM_USERS")
	defer db.Close()

	_, err := db.Exec("UPDATE ROOM_USERS SET score = ? WHERE id_room = ? AND id_user = ?", score, roomID, userID)
	if err != nil {
		log.Println("Error setting user score:", err)
	}
}

// ---------------------------GAME_ROOM---------------------------\\

type GameRoomData struct {
	GameState                bool
	NumberOfGameTurns        int
	BlindtestTrackIndex      int
	BlindtestTimeOfMusic     int
	BlindtestTimeToAnswer    int
	PetitbacCategories       string
	PetitbacTimeToAnswer     int
	GuessthesongDifficulty   string
	GuessthesongTimeToAnswer int
}

func GetRoomData(roomID int) GameRoomData {
	db := InitTable("GAME_ROOM")
	defer db.Close()

	var gameRoomData GameRoomData
	row := db.QueryRow("SELECT game_state, COALESCE(number_of_game_turns, 0), COALESCE(blindtest_track_index, -1), COALESCE(blindtest_time_of_music, 0), COALESCE(blindtest_time_to_answer, 0), COALESCE(petitbac_categories, ''), COALESCE(petitbac_time_to_answer, 0), COALESCE(guessthesong_difficulty, ''), COALESCE(guessthesong_time_to_answer, 0) FROM GAME_ROOM WHERE id_room = ?", roomID)
	err := row.Scan(&gameRoomData.GameState, &gameRoomData.NumberOfGameTurns, &gameRoomData.BlindtestTrackIndex, &gameRoomData.BlindtestTimeOfMusic, &gameRoomData.BlindtestTimeToAnswer, &gameRoomData.PetitbacCategories, &gameRoomData.PetitbacTimeToAnswer, &gameRoomData.GuessthesongDifficulty, &gameRoomData.GuessthesongTimeToAnswer)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No game data")
		}
		return GameRoomData{}
	}
	return gameRoomData
}

func UpdateGameData(data GameRoomData, roomID int) {
	db := InitTable("GAME_ROOM")
	defer db.Close()

	query := `
        UPDATE GAME_ROOM SET
            game_state = ?,
            number_of_game_turns = ?,
			blindtest_track_index = ?,
            blindtest_time_of_music = ?,
            blindtest_time_to_answer = ?,
            petitbac_categories = ?,
            petitbac_time_to_answer = ?,
            guessthesong_difficulty = ?,
            guessthesong_time_to_answer = ?
        WHERE id_room = ?
    `

	_, err := db.Exec(query,
		data.GameState,
		data.NumberOfGameTurns,
		data.BlindtestTrackIndex,
		data.BlindtestTimeOfMusic,
		data.BlindtestTimeToAnswer,
		data.PetitbacCategories,
		data.PetitbacTimeToAnswer,
		data.GuessthesongDifficulty,
		data.GuessthesongTimeToAnswer,
		roomID,
	)
	if err != nil {
		fmt.Println(err)
	}
}
