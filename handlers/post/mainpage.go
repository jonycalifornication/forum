package post

import (
	"fmt"
	"forum/database"
	"forum/handlers"
	"forum/models"
	"log"
	"net/http"
)

func MainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		handlers.ErrorHandler(w, http.StatusNotFound)
		return
	}

	cookie, err := r.Cookie("session_id")

	authenticated := false
	username := ""
	if err == nil {
		sessionID := cookie.Value
		if sessionData, exists := handlers.SessionStore[sessionID]; exists {
			authenticated = true
			username = sessionData["username"]
		}
	}

	var posts []models.Post
	posts, err = database.GetAllPosts()
	if err != nil {
		log.Println(err)
		handlers.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	userInfo, err := database.GetUserInfoByUsername(username)
	role := ""
	if err != nil {
		log.Println(err)
	} else if userInfo != nil {
		role = userInfo.Role
	}

	fmt.Println("--------------", role)

	var errorMessage string
	if r.URL.Query().Get("error") != "" {
		errorMessage = "An error occurred."
	}

	data := struct {
		Authenticated bool
		Username      string
		Posts         []models.Post
		ErrorMessage  string
		Role          string
	}{
		Authenticated: authenticated,
		Username:      username,
		Posts:         posts,
		ErrorMessage:  errorMessage,
		Role:          role,
	}

	switch r.Method {
	case "GET":
		handlers.RenderTemplate(w, "index.html", data)
	default:
		handlers.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}
