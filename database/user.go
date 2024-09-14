package database

import (
	"errors"
	"forum/internal"
	"forum/models"
	"log"
)

func CreateUser(user *models.UserCreate) error {
	hash, err := internal.PasswordHash(user.Password)
	if err != nil {
		log.Println("Error hashing password")
		return err
	}
	user.Password = string(hash)
	stmt, err := DB.Prepare("INSERT INTO users (username, email, password) VALUES (?, ?, ?)")
	if err != nil {
		log.Println("Error preparing query", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Name, user.Email, user.Password)
	if err != nil {
		log.Println("Error creating new user", err)
	}

	return err
}

func AuthenticateUser(username, password string) error {
	storedpasswordhash, err := GetPasswordHash(username)
	if err != nil {
		log.Println("Error getting password hash from database")
		return errors.New("User doesn't exist")
	}
	if err := internal.CheckPasswordHash(password, storedpasswordhash); err != nil {
		log.Println("Password and stored password hash do not match")
		return errors.New("Incorrect password")
	}
	return nil
}

func GetPasswordHash(username string) ([]byte, error) {
	var hash string
	query := "SELECT password FROM users WHERE username = ?"
	err := DB.QueryRow(query, username).Scan(&hash)
	if err != nil {
		log.Println("Error getting hash:", err)
		return nil, err
	}

	return []byte(hash), nil
}

func GetUserID(username string) (int, error) {
	var userID int
	query := "SELECT id FROM users WHERE username = ?"
	if err := DB.QueryRow(query, username).Scan(&userID); err != nil {
		log.Println(err)
		return 0, err
	}
	return userID, nil
}
