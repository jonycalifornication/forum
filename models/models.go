package models

import (
	"mime/multipart"
	"time"
)

type UserCreate struct {
	Name     string
	Email    string
	Password string
	Role     string
}

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
	Role     string
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

type ApplyModeratorRequest struct {
	Username  string
	CreatedAt string
}

type Report struct {
	ID                int       `json:"id"`
	Username          string    `json:"username"`
	PostID            int       `json:"post_id"`
	Reason            string    `json:"reason"`
	ModeratorUsername string    `json:"moderator_username"`
	PostURL           string    `json:"post_url"`
	CreatedAt         time.Time `json:"created_at"`
}

type Reply struct {
	ID                int       `json:"id"`
	ModeratorUsername string    `json:"moderator_username"`
	Admin             string    `json:"admin"`
	PostID            int       `json:"post_id"`
	ReplyText         string    `json:"reply_text"`
	CreatedAt         time.Time `json:"created_at"`
}
