package models

import (
	"mime/multipart"
	"time"
)

type UserCreate struct {
	Name     string
	Email    string
	Password string
}

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type PostCreate struct {
	Username   string   `json:"username"`
	Title      string   `json:"title"`
	Text       string   `json:"text"`
	Categories []string `json:"categories"`
}

type Post struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	Title      string    `json:"title"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"created_at"`
	Categories []string  `json:"categories"`
	ImagePath  string    `json:"image_path"`
}

type Comment struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Text         string    `json:"text"`
	CreatedAt    time.Time `json:"created_at"`
	LikeCount    int       `json:"like_count"`
	DislikeCount int       `json:"dislike_count"`
	CanDelete    bool
}

type Error struct {
	Text string
	Code int
}

type File struct {
	FileGiven multipart.File
	Header    *multipart.FileHeader
}
