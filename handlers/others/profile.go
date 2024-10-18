package others

import (
	"forum/database"
	"forum/handlers"
	"forum/models"
	"log"
	"net/http"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user_profile" {
		handlers.ErrorHandler(w, http.StatusNotFound)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
		return
	}

	sessionID := cookie.Value
	user, ok := handlers.SessionStore[sessionID]
	if !ok {
		http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
		return
	}

	username := user["username"]

	userInfo, err := database.GetUserInfoByUsername(username)
	if err != nil {
		log.Println(err)
	}

	posts, err := database.GetPostsByUsername(username)
	if err != nil {
		log.Println(err)
	}

	repliesFromAdmin, err := database.GetRepliesByModeratorUsername(username)
	if err != nil {
		log.Println(err)
	}

	var limitedPosts []models.Post
	if len(posts) > 3 {
		limitedPosts = posts[:3]
	} else {
		limitedPosts = posts
	}

	data := struct {
		RepliesFromAdmin []models.Reply
		UserInfo         models.User
		Posts            []models.Post
		ErrorMessage     string
	}{
		RepliesFromAdmin: repliesFromAdmin,
		UserInfo:         *userInfo,
		Posts:            limitedPosts,
		ErrorMessage:     "",
	}

	switch r.Method {
	case "GET":
		handlers.RenderTemplate(w, "profile.html", data)
	default:
		handlers.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}
