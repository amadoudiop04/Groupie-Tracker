package database

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
)

func HashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func RegisterUser(db *sql.DB, username, password, email string) error {
	hashedPassword := HashPassword(password)

	_, err := db.Exec("INSERT INTO USER (pseudo, email, password) VALUES (?, ?, ?)", username, strings.ToLower(email), hashedPassword)
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

func UniqueEmail(email string) bool {
	db := InitTable("USER")
	defer db.Close()

	rowsUsers := selectValueFromTable(db, "USER", "email")

	for rowsUsers.Next() {
		var userEmail string
		err := rowsUsers.Scan(&userEmail)
		if err != nil {
			fmt.Println(err)
		}
		if userEmail == email {
			return false
		}
	}

	return true
}

func UniqueUsername(username string) bool {
	db := InitTable("USER")
	defer db.Close()

	rowsUsers := selectValueFromTable(db, "USER", "pseudo")

	for rowsUsers.Next() {
		var userPseudo string
		err := rowsUsers.Scan(&userPseudo)
		if err != nil {
			fmt.Println(err)
		}
		if userPseudo == username {
			return false
		}
	}

	return true
}

func AuthenticateUser(username, password string) (bool, error) {
	db := InitTable("USER")
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

func DisplayUserTable(rows *sql.Rows) {
	for rows.Next() {
		var users User
		err := rows.Scan(&users.Id, &users.Pseudo, &users.Email, &users.Password, &users.recoveryCode)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(users)
	}
}

func SelectAllFromTable(db *sql.DB, table string) *sql.Rows {
	query := "SELECT * FROM " + table
	result, _ := db.Query(query)
	return result
}

func selectValueFromTable(db *sql.DB, table string, value string) *sql.Rows {
	query := "SELECT " + value + " FROM " + table
	result, _ := db.Query(query)
	return result
}

func ReplaceEmptyString(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func GenerateCode() int {
	return rand.Intn(900000) + 100000
}

//structure

type User struct {
	Id           int
	Pseudo       string
	Email        string
	Password     string
	recoveryCode *int
}
