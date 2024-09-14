package comment

import (
	"fmt"
	"forum/database"
	"forum/handlers"
	"forum/internal"
	"log"
	"net/http"
	"strconv"
)

func Comment(w http.ResponseWriter, r *http.Request) {
	// Проверка наличия куки сессии
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
		return
	}

	sessionID := cookie.Value

	// Проверка валидности sessionID в sessionStore
	sessionData, ok := handlers.SessionStore[sessionID]
	if !ok {
		http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
		return
	}

	username := sessionData["username"]
	userID, err := database.GetUserID(username)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Извлечение ID поста из формы
	postIDStr := r.FormValue("postId")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		log.Println("Error converting postId to int:", err)
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Извлечение текста комментария из формы
	text := r.FormValue("text")
	text = internal.SanitizeComment(text)

	// Вставка комментария в базу данных
	if err := database.InsertComment(userID, postID, username, text); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Перенаправление обратно на страницу поста
	http.Redirect(w, r, fmt.Sprintf("/posts/?id=%d", postID), http.StatusSeeOther)
}
