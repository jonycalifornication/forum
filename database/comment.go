package database

import (
	"forum/models"
	"log"
)

func InsertComment(userID, postID int, username, text string) error {
	query := `
		INSERT INTO comments (post_id, user_id, username, text)
		VALUES (?, ?, ?, ?)
	`

	if _, err := DB.Exec(query, postID, userID, username, text); err != nil {
		log.Println("Error inserting comment:", err)
		return err
	}

	return nil
}

func GetCommentsByPostId(postID int, username string) ([]models.Comment, error) {
	var comments []models.Comment

	query := `
		SELECT id, username, text, created_at
		FROM comments
		WHERE post_id = ?
		ORDER BY created_at ASC
	`

	rows, err := DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.Username, &comment.Text, &comment.CreatedAt); err != nil {
			log.Println("Error scanning comment row:", err)
			return nil, err
		}

		likeCount, dislikeCount, err := GetCommentReactionCounts(comment.ID)
		if err != nil {
			log.Println("Error getting reaction counts for comment:", err)
			return nil, err
		}
		comment.LikeCount = likeCount
		comment.DislikeCount = dislikeCount
		comment.CanDelete = (comment.Username == username)

		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error with rows iteration:", err)
		return nil, err
	}

	return comments, nil
}

func DeleteCommentByID(commentID int) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec("DELETE FROM comment_reactions WHERE comment_id = ?", commentID); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec("DELETE FROM comments WHERE id = ?", commentID); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Println("Comment with ID has been deleted.", commentID)
	return nil
}
