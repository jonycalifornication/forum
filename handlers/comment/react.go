package comment

import (
	"fmt"
	"forum/database"
	"forum/handlers"
	"log"
	"net/http"
	"strconv"
)

func ReactComment(w http.ResponseWriter, r *http.Request) {
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

	commentIDStr := r.FormValue("commentId")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		log.Println("Error converting commentId to int:", err)
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	userID, err := database.GetUserID(sessionData["username"])
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	action := r.FormValue("action")

	if err := database.ToggleCommentReaction(commentID, userID, action); err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	postIDStr := r.FormValue("postId")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		log.Println("Error converting postId to int:", err)
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/posts/?id=%d", postID), http.StatusSeeOther)
}
