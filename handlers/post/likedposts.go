package post

import (
	"forum/database"
	"forum/handlers"
	"forum/models"
	"log"
	"net/http"
)

func LikedPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/liked_posts" {
		handlers.ErrorHandler(w, http.StatusNotFound)
		return
	}

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

	userid, err := database.GetUserID(sessionData["username"])
	if err != nil {
		log.Println(err)
		handlers.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	posts, err := database.GetLikedPost(userid)
	if err != nil {
		log.Println(err)
		handlers.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	data := struct {
		Authenticated bool
		Username      string
		Posts         []models.Post
	}{
		Authenticated: true,
		Username:      sessionData["username"],
		Posts:         posts,
	}

	handlers.RenderTemplate(w, "liked_posts.html", data)
}
