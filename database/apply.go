package database

import (
	"fmt"
	"forum/models"
	"log"
)

func SendApplyModeratorRequest(username string) error {
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM apply WHERE username = ? LIMIT 1)`

	err := DB.QueryRow(checkQuery, username).Scan(&exists)
	if err != nil {
		log.Println("Error checking for existing application:", err)
		return err
	}

	if exists {
		log.Println("User has already applied for moderator.")
		return nil // Возвращаем nil, чтобы не считать это ошибкой
	}

	insertQuery := `INSERT INTO apply (username) VALUES (?)`

	if _, err := DB.Exec(insertQuery, username); err != nil {
		log.Println("Error inserting application:", err)
		return err
	}

	log.Println("Application for moderator successfully submitted.")
	return nil
}

func GetApplyModeratorRequest() ([]models.ApplyModeratorRequest, error) {
	query := "SELECT username, created_at FROM apply"

	rows, err := DB.Query(query)
	if err != nil {
		log.Println("Error retrieving apply moderator requests:", err)
		return nil, err
	}
	defer rows.Close()

	var requests []models.ApplyModeratorRequest

	for rows.Next() {
		var request models.ApplyModeratorRequest
		if err := rows.Scan(&request.Username, &request.CreatedAt); err != nil {
			log.Println("Error scanning apply request:", err)
			return nil, err
		}
		requests = append(requests, request)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error during row iteration:", err)
		return nil, err
	}

	return requests, nil
}

func UpdateUserRoleToModerator(username string) error {
	// Начинаем транзакцию, чтобы обе операции (обновление и удаление) выполнились атомарно
	tx, err := DB.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return err
	}

	// SQL-запрос для обновления роли пользователя
	updateQuery := "UPDATE users SET role = 'moderator' WHERE username = ?"
	_, err = tx.Exec(updateQuery, username)
	if err != nil {
		log.Println("Error updating user role to moderator:", err)
		tx.Rollback() // Откатываем изменения в случае ошибки
		return err
	}

	// SQL-запрос для удаления записи из таблицы apply
	deleteQuery := "DELETE FROM apply WHERE username = ?"
	_, err = tx.Exec(deleteQuery, username)
	if err != nil {
		log.Println("Error deleting username from apply table:", err)
		tx.Rollback() // Откатываем изменения в случае ошибки
		return err
	}

	// Подтверждаем транзакцию
	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction:", err)
		return err
	}

	log.Printf("User %s has been promoted to moderator and removed from apply table.", username)
	return nil
}

func DenyUpdateUsertoModerator(username string) error {
	// SQL-запрос для удаления записи из таблицы apply
	deleteQuery := "DELETE FROM apply WHERE username = ?"

	// Выполнение запроса для удаления заявки
	_, err := DB.Exec(deleteQuery, username)
	if err != nil {
		log.Println("Error deleting apply request:", err)
		return err
	}

	log.Printf("Moderator application for user %s has been denied and removed.", username)
	return nil
}

func DemoteToUser(username string) error {
	// Проверяем, существует ли пользователь с указанным username
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE username = ? LIMIT 1)`

	err := DB.QueryRow(checkQuery, username).Scan(&exists)
	if err != nil {
		log.Println("Error checking if user exists:", err)
		return err
	}

	if !exists {
		log.Println("User does not exist.")
		return fmt.Errorf("user %s does not exist", username) // Возвращаем ошибку, если пользователь не найден
	}

	// SQL-запрос для понижения роли пользователя
	updateQuery := `UPDATE users SET role = 'user' WHERE username = ?`

	// Выполняем запрос
	_, err = DB.Exec(updateQuery, username)
	if err != nil {
		log.Println("Error demoting user to 'user' role:", err)
		return err
	}

	log.Printf("User %s has been demoted to 'user' role.", username)
	return nil
}
