package post

import (
	"forum/database"
	"forum/handlers"
	"log"
	"net/http"
	"strconv"
)

func DeletePost(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	cookie, err := r.Cookie("session_id")
	if err != nil {
		// Перенаправление на страницу входа, если куки не найдены
		http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
		return
	}

	sessionID := cookie.Value
	sessionData, ok := handlers.SessionStore[sessionID]
	if !ok {
		// Перенаправление на страницу входа, если сессия не найдена
		http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
		return
	}

	postIDStr := r.FormValue("postId")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		log.Println("Error converting postId to int:", err)
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Получение информации о посте
	post, err := database.GetPostsById(postID)
	if err != nil {
		log.Println("Error fetching post:", err)
		handlers.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	username := sessionData["username"]

	userInfo, err := database.GetUserInfoByUsername(username)
	if err != nil {
		log.Println("Error fetching user info:", err)
		handlers.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	if post.Username != username && (userInfo.Role != "admin" && userInfo.Role != "moderator") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := database.DeletePostByID(postID); err != nil {
		log.Println("Error deleting post:", err)
		handlers.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
