package database

import (
	"database/sql"
	"errors"
	"fmt"
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

func GetUserInfoByUsername(username string) (*models.User, error) {
	var user models.User
	query := "SELECT username, email, role FROM users WHERE username = ?;"

	err := DB.QueryRow(query, username).Scan(&user.Name, &user.Email, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func GetAllUsers() ([]models.User, error) {
	query := "SELECT username, email, role FROM users"

	rows, err := DB.Query(query)
	if err != nil {
		log.Println("Error retrieving users:", err)
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Name, &user.Email, &user.Role); err != nil {
			log.Println("Error scanning user data:", err)
			return nil, err
		}
		users = append(users, user)
	}

	// Проверка на ошибки после завершения итерации
	if err := rows.Err(); err != nil {
		log.Println("Error during row iteration:", err)
		return nil, err
	}

	return users, nil
}
