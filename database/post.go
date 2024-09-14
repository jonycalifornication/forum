package database

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/models"
	"io"
	"log"
	"os"
	"path/filepath"
)

func CreatePost(post *models.PostCreate, file models.File) error {
	saveDir := "web/images"

	err := os.MkdirAll(saveDir, os.ModePerm)
	if err != nil {
		return err
	}
	dst, err := os.Create(filepath.Join("web/images", file.Header.Filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file.FileGiven)
	if err != nil {
		return err
	}

	fullPath := filepath.Join(saveDir, file.Header.Filename)
	fmt.Println("fullllpath", fullPath)
	query := "INSERT INTO posts (username, title, text, image_path) VALUES (?, ?, ?, ?)"
	result, err := DB.Exec(query, post.Username, post.Title, post.Text, fullPath)
	if err != nil {
		log.Println("Error creating new post", err)
		return err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	for _, category := range post.Categories {
		var categoryID int
		if err := DB.QueryRow("SELECT id FROM categories WHERE name = ?", category).Scan(&categoryID); err != nil {
			log.Println("Error retrieving category ID:", err)
			return err
		}
		query2 := "INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)"
		if _, err := DB.Exec(query2, postID, categoryID); err != nil {
			log.Println("Error inserting post category:", err)
			return err
		}
	}

	return nil
}

func GetAllPosts() ([]models.Post, error) {
	var posts []models.Post
	query := "SELECT id, username, title, text, created_at FROM posts ORDER BY created_at DESC;"

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Username, &post.Title, &post.Text, &post.CreatedAt); err != nil {
			return nil, err
		}

		categories, err := getCategoriesForPost(post.ID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories

		posts = append(posts, post)
	}

	return posts, nil
}

func GetPostsByUsername(username string) ([]models.Post, error) {
	var posts []models.Post

	query := "SELECT * FROM posts WHERE username = ? ORDER BY created_at DESC ;"
	rows, err := DB.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Username, &post.Title, &post.Text, &post.CreatedAt, &post.ImagePath); err != nil {
			return nil, err
		}

		categories, err := getCategoriesForPost(post.ID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories
		posts = append(posts, post)
	}

	return posts, nil
}

func GetPostsById(id int) (models.Post, error) {
	var post models.Post

	query := "SELECT id, username, title, text, created_at, image_path FROM posts WHERE id = ?"
	row := DB.QueryRow(query, id)

	if err := row.Scan(&post.ID, &post.Username, &post.Title, &post.Text, &post.CreatedAt, &post.ImagePath); err != nil {
		if err == sql.ErrNoRows {
			return models.Post{}, err
		}
		return models.Post{}, err
	}

	categories, err := getCategoriesForPost(post.ID)
	if err != nil {
		return models.Post{}, err
	}
	post.Categories = categories

	return post, nil
}

func getCategoriesForPost(postID int) ([]string, error) {
	var categories []string

	query := `
        SELECT c.name
        FROM post_categories pc
        JOIN categories c ON pc.category_id = c.id
        WHERE pc.post_id = ?
    `
	rows, err := DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func GetPostsByCategory(category string) ([]models.Post, error) {
	var categoryExists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM categories WHERE name = ?)"
	err := DB.QueryRow(checkQuery, category).Scan(&categoryExists)
	if err != nil {
		return nil, err
	}

	if !categoryExists {
		return nil, errors.New("category doesnt exist")
	}
	var posts []models.Post
	query := `
		SELECT p.id, p.username, p.title, p.text, p.created_at
		FROM posts p
		JOIN post_categories pc ON p.id = pc.post_id
		JOIN categories c ON pc.category_id = c.id
		WHERE c.name = ?
		ORDER BY p.created_at DESC;
	`

	rows, err := DB.Query(query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Username, &post.Title, &post.Text, &post.CreatedAt); err != nil {
			return nil, err
		}
		categories, err := getCategoriesForPost(post.ID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories
		posts = append(posts, post)
	}

	return posts, nil
}

func GetLikedPost(userID int) ([]models.Post, error) {
	var posts []models.Post
	query := `
    SELECT
        p.id,
        p.title,
        p.text,
        p.username,
        p.created_at
    FROM
        posts p
    JOIN
        post_reactions pr ON p.id = pr.post_id
    WHERE
        pr.user_id = ? AND pr.reaction_type = 'like';
    `
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Text, &post.Username, &post.CreatedAt); err != nil {
			return nil, err
		}
		categories, err := getCategoriesForPost(post.ID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func DeletePostByID(postID int) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec("DELETE FROM post_reactions WHERE post_id = ?", postID); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec("DELETE FROM post_categories WHERE post_id = ?", postID); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec("DELETE FROM posts WHERE id = ?", postID); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Println("Post with ID has been deleted.", postID)
	return nil
}
