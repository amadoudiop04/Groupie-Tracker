package database

import (
	"database/sql"
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

func CreateRoom(db *sql.DB, createdBy int, maxPlayers int, name string, gameID int) error {
	_, err := db.Exec("INSERT INTO ROOMS (created_by, max_player, name, id_game) VALUES (?, ?, ?, ?)", createdBy, maxPlayers, name, gameID)
	if err != nil {
		log.Println("Error creating room:", err)
	}
	return err
}

func GetAvailableRooms(db *sql.DB) ([]Room, error) {
	var rooms []Room
	rows, err := db.Query("SELECT id, created_by, max_player, name, id_game FROM ROOMS")
	if err != nil {
		log.Println("Error fetching rooms:", err)
		return rooms, err
	}
	defer rows.Close()

	for rows.Next() {
		var room Room
		if err := rows.Scan(&room.Id, &room.CreatedBy, &room.MaxPlayers, &room.Name, &room.GameID); err != nil {
			log.Println("Error scanning room row:", err)
			continue
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func JoinRoom(db *sql.DB, roomID int, userID int) error {
	_, err := db.Exec("INSERT INTO ROOM_USERS (id_room, id_user) VALUES (?, ?)", roomID, userID)
	if err != nil {
		log.Println("Error joining room:", err)
	}
	return err
}

func LeaveRoom(db *sql.DB, roomID int, userID int) error {
	_, err := db.Exec("DELETE FROM ROOM_USERS WHERE id_room = ? AND id_user = ?", roomID, userID)
	if err != nil {
		log.Println("Error leaving room:", err)
	}
	return err
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
