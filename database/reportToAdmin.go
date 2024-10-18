package database

import (
	"forum/models"
	"log"
	"time"
)

func SaveReportToAdmin(username string, moderatorUsername string, postID int, reportReason string, postURL string) error {
	insertQuery := `INSERT INTO reports (username, post_id, reason, moderator_username, postURL) VALUES (?, ?, ?, ?, ?)`

	_, err := DB.Exec(insertQuery, username, postID, reportReason, moderatorUsername, postURL)
	if err != nil {
		log.Println("Error saving report to admin:", err)
		return err
	}

	log.Printf("Report by user %s for post ID %d has been saved for moderator %s.", username, postID, moderatorUsername)
	return nil
}

func GetAllReports() ([]models.Report, error) {
	var reports []models.Report
	query := "SELECT id, username, post_id, reason, moderator_username, postURL, created_at FROM reports ORDER BY created_at DESC;"

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var report models.Report
		if err := rows.Scan(&report.ID, &report.Username, &report.PostID, &report.Reason, &report.ModeratorUsername, &report.PostURL, &report.CreatedAt); err != nil {
			return nil, err
		}

		reports = append(reports, report)
	}

	return reports, nil
}

func SaveReplyToAdmin(moderatorUsername string, admin string, postID int, replyText string) error {
	query := `
        INSERT INTO replies (moderator_username, admin, post_id, reply_text, created_at) 
        VALUES (?, ?, ?, ?, ?)
    `

	_, err := DB.Exec(query, moderatorUsername, admin, postID, replyText, time.Now())
	if err != nil {
		log.Println("Error saving reply to moderator:", err)
		return err
	}

	log.Println("Reply saved successfullyl")
	return nil
}

func DeleteReportByID(reportID int) error {
	deleteQuery := `DELETE FROM reports WHERE id = ?`

	_, err := DB.Exec(deleteQuery, reportID)
	if err != nil {
		log.Println("Error deleting report:", err)
		return err
	}

	log.Printf("Report with ID %d has been deleted.", reportID)
	return nil
}

func GetRepliesByModeratorUsername(moderatorUsername string) ([]models.Reply, error) {
	var replies []models.Reply

	query := `
		SELECT id, moderator_username, admin, post_id, reply_text, created_at
		FROM replies
		WHERE moderator_username = ?
		ORDER BY created_at ASC
	`

	rows, err := DB.Query(query, moderatorUsername)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var reply models.Reply
		if err := rows.Scan(&reply.ID, &reply.ModeratorUsername, &reply.Admin, &reply.PostID, &reply.ReplyText, &reply.CreatedAt); err != nil {
			log.Println("Error scanning reply row:", err)
			return nil, err
		}

		replies = append(replies, reply)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error with rows iteration:", err)
		return nil, err
	}

	return replies, nil
}

func DeleteReplyByID(replyID int) error {
	query := `DELETE FROM replies WHERE id = ?`

	// Подготовка запроса
	stmt, err := DB.Prepare(query)
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		return err
	}
	defer stmt.Close()

	// Выполнение запроса с подстановкой ID
	_, err = stmt.Exec(replyID)
	if err != nil {
		log.Println("Error executing SQL statement:", err)
		return err
	}

	log.Println("Reply with ID", replyID, "was successfully deleted.")
	return nil
}
