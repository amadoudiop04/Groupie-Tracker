package database

import (
	"database/sql"
	"log"
)

func InitTable(database string) *sql.DB {
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
				
				CREATE TABLE IF NOT EXISTS ROOMS (
					id INTEGER PRIMARY KEY,
					created_by INTEGER NOT NULL,
					max_player INTEGER NOT NULL,
					number_of_player INTEGER NOT NULL,
					name TEXT NOT NULL,
					id_game INTEGER,
					FOREIGN KEY (created_by) REFERENCES USER(id),
					FOREIGN KEY (id_game) REFERENCES GAMES(id)
				);
				
				CREATE TABLE IF NOT EXISTS ROOM_USERS (
					id_room INTEGER,
					id_user INTEGER,
					score INTEGER,
					FOREIGN KEY (id_room) REFERENCES ROOMS(id),
					FOREIGN KEY (id_user) REFERENCES USER(id),
					PRIMARY KEY (id_room, id_user)
				);
				
				CREATE TABLE IF NOT EXISTS GAMES (
					id INTEGER PRIMARY KEY,
					name TEXT NOT NULL
				);
				
				CREATE TABLE IF NOT EXISTS GAME_ROOM (
					id_room INTEGER PRIMARY KEY,
					number_of_game_turns INTEGER,
					blindtest_time_of_music INTEGER,
					blindtest_time_to_answer INTEGER,
					petitbac_categories TEXT,
					petitbac_time_to_answer INTEGER,
					guessthesong_difficulty TEXT,
					guessthesong_time_to_answer INTEGER,
					FOREIGN KEY (id_room) REFERENCES ROOMS(id)
				);
				`
	_, err = db.Exec(sqlStmt)

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func GetUserIdByUsername(username string) (string, error) {
	db := InitTable("USER")
	defer db.Close()

	var userID string
	err := db.QueryRow("SELECT id FROM USER WHERE pseudo = ?", username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		log.Println("Erreur lors de la récupération de l'ID de l'utilisateur:", err)
		return "", err
	}

	return userID, nil
}

func GetUserIdByEmail(email string) (string, error) {
	db := InitTable("USER")
	defer db.Close()

	var userID string
	err := db.QueryRow("SELECT id FROM USER WHERE email = ?", email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		log.Println("Erreur lors de la récupération de l'ID de l'utilisateur:", err)
		return "", err
	}

	return userID, nil
}

func GetUserData(ID string) (User, error) {
	db := InitTable("USER")
	defer db.Close()

	var user User
	err := db.QueryRow("SELECT id, pseudo, email, password FROM USER WHERE id = ?", ID).
		Scan(&user.Id, &user.Pseudo, &user.Email, &user.Password)
	if err != nil {
		log.Println("Erreur lors de la récupération des données de l'utilisateur:", err)
		return User{}, err
	}

	return user, nil
}

func SetUserData(ID, pseudo, email, password string) error {
	db := InitTable("USER")
	defer db.Close()

	_, err := db.Exec("UPDATE USER SET pseudo = ?, email = ?, password = ? WHERE id = ?", pseudo, email, password, ID)
	if err != nil {
		log.Println("Erreur lors de la mise à jour des données de l'utilisateur:", err)
		return err
	}

	return nil
}
