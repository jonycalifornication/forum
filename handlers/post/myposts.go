package post

import (
	"forum/database"
	"forum/handlers"
	"forum/models"
	"net/http"
)

func MyPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handlers.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
		return
	}

	sessionID := cookie.Value
	sessionData, ok := handlers.SessionStore[sessionID]
	if !ok {
		http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
		return
	}

	username := sessionData["username"]
	userPosts, err := database.GetPostsByUsername(username)
	if err != nil {
		handlers.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	data := struct {
		Authenticated bool
		Username      string
		Posts         []models.Post
	}{
		Authenticated: true,
		Username:      username,
		Posts:         userPosts,
	}

	handlers.RenderTemplate(w, "my_posts.html", data)
}
