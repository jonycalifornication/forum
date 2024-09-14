package database

import (
	"database/sql"
	"fmt"
	"time"
)

func ToggleReaction(postID, userID int, action string) error {
	existingReaction, err := ReactExistCheck(postID, userID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking existing reaction: %w", err)
	}

	// Remove reaction if it is the same
	if existingReaction == action {
		_, err = DB.Exec("DELETE FROM post_reactions WHERE post_id = ? AND user_id = ?", postID, userID)
		if err != nil {
			return fmt.Errorf("error deleting existing reaction: %w", err)
		}
		return nil
	}

	// Insert or update reaction
	if existingReaction == "" {
		_, err = DB.Exec("INSERT INTO post_reactions (post_id, user_id, reaction_type, created_at) VALUES (?, ?, ?, ?)", postID, userID, action, time.Now())
	} else {
		_, err = DB.Exec("UPDATE post_reactions SET reaction_type = ?, created_at = ? WHERE post_id = ? AND user_id = ?", action, time.Now(), postID, userID)
	}

	if err != nil {
		return fmt.Errorf("error updating reaction: %w", err)
	}

	return nil
}

func ToggleCommentReaction(commentID, userID int, action string) error {
	// Check existing reaction
	existingReaction, err := CommentReactExistCheck(commentID, userID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking existing reaction: %w", err)
	}

	// Remove reaction if it is the same
	if existingReaction == action {
		_, err = DB.Exec("DELETE FROM comment_reactions WHERE comment_id = ? AND user_id = ?", commentID, userID)
		if err != nil {
			return fmt.Errorf("error deleting existing reaction: %w", err)
		}
		return nil
	}

	// Insert or update reaction
	if existingReaction == "" {
		_, err = DB.Exec("INSERT INTO comment_reactions (comment_id, user_id, reaction_type, created_at) VALUES (?, ?, ?, ?)", commentID, userID, action, time.Now())
	} else {
		_, err = DB.Exec("UPDATE comment_reactions SET reaction_type = ?, created_at = ? WHERE comment_id = ? AND user_id = ?", action, time.Now(), commentID, userID)
	}

	if err != nil {
		return fmt.Errorf("error updating reaction: %w", err)
	}

	return nil
}

func ReactExistCheck(postID, userID int) (string, error) {
	var existingReaction string
	err := DB.QueryRow("SELECT reaction_type FROM post_reactions WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&existingReaction)
	if err != nil {
		return "", err
	}

	return existingReaction, nil
}

func CommentReactExistCheck(commentID, userID int) (string, error) {
	var existingReaction string
	err := DB.QueryRow("SELECT reaction_type FROM comment_reactions WHERE comment_id = ? AND user_id = ?", commentID, userID).Scan(&existingReaction)
	if err != nil {
		return "", err
	}

	return existingReaction, nil
}

func GetReactionCounts(postID int) (int, int, error) {
	var likeCount, dislikeCount sql.NullInt64

	query := `SELECT
                COALESCE(SUM(CASE WHEN reaction_type = 'like' THEN 1 ELSE 0 END), 0) AS like_count,
                COALESCE(SUM(CASE WHEN reaction_type = 'dislike' THEN 1 ELSE 0 END), 0) AS dislike_count
              FROM post_reactions
              WHERE post_id = ?;`

	err := DB.QueryRow(query, postID).Scan(&likeCount, &dislikeCount)
	if err != nil {
		return 0, 0, err
	}

	return int(likeCount.Int64), int(dislikeCount.Int64), nil
}

func GetCommentReactionCounts(commentID int) (int, int, error) {
	var likeCount, dislikeCount sql.NullInt64

	query := `SELECT
                COALESCE(SUM(CASE WHEN reaction_type = 'like' THEN 1 ELSE 0 END), 0) AS like_count,
                COALESCE(SUM(CASE WHEN reaction_type = 'dislike' THEN 1 ELSE 0 END), 0) AS dislike_count
              FROM comment_reactions
              WHERE comment_id = ?;`

	err := DB.QueryRow(query, commentID).Scan(&likeCount, &dislikeCount)
	if err != nil {
		return 0, 0, err
	}

	return int(likeCount.Int64), int(dislikeCount.Int64), nil
}
