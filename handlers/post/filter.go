package post

import (
	"forum/database"
	"forum/handlers"
	"forum/models"
	"log"
	"net/http"
)

func Filter(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/category/" {
		handlers.ErrorHandler(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		handlers.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	category := r.URL.Query().Get("name")

	cookie, err := r.Cookie("session_id")

	authenticated := false
	username := ""
	if err == nil {
		if sessionData, exists := handlers.SessionStore[cookie.Value]; exists {
			authenticated = true
			username = sessionData["username"]
		}
	}

	posts, err := database.GetPostsByCategory(category)
	if err != nil {
		log.Println("Error getting posts")
		handlers.ErrorHandler(w, http.StatusNotFound)
		return
	}

	data := struct {
		Category      string
		Authenticated bool
		Username      string
		Posts         []models.Post
		ErrorMessage  string
	}{
		Category:      category,
		Authenticated: authenticated,
		Username:      username,
		Posts:         posts,
		ErrorMessage:  "", // Или установите значение, если есть ошибка
	}

	handlers.RenderTemplate(w, "category.html", data)
}
