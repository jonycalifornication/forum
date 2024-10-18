package post

import (
	"forum/database"
	"forum/handlers"
	"forum/models"
	"log"
	"net/http"
	"strconv"
)

func Post(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/posts/" {
		handlers.ErrorHandler(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		handlers.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	// Получение сессии
	cookie, err := r.Cookie("session_id")

	var username string
	authenticated := false
	var userInfo *models.User
	if err == nil {
		sessionID := cookie.Value
		sessionData, ok := handlers.SessionStore[sessionID]
		if ok {
			username = sessionData["username"]
			authenticated = true

			// Получаем информацию о пользователе, только если авторизован
			userInfo, err = database.GetUserInfoByUsername(username)
			if err != nil {
				log.Println("Error getting user info:", err)
				handlers.ErrorHandler(w, http.StatusInternalServerError)
				return
			}
		}
	} // Закрываем блок if

	// Преобразование id из строки в число
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting post ID:", err)
		handlers.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	// Получение поста
	post, err := database.GetPostsById(id)
	if err != nil {
		log.Println("Error getting post by ID:", err)
		handlers.ErrorHandler(w, http.StatusNotFound)
		return
	}

	// Получение реакции
	likeCount, dislikeCount, err := database.GetReactionCounts(id)
	if err != nil {
		log.Println("Error getting reaction counts:", err)
		likeCount = 0
		dislikeCount = 0
	}

	// Получение комментариев
	comments, err := database.GetCommentsByPostId(id, username)
	if err != nil {
		log.Println("Error getting comments by post ID:", err)
		handlers.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	// Проверка, может ли пользователь удалить пост
	var delete bool

	if username == post.Username || (userInfo != nil && (userInfo.Role == "admin" || userInfo.Role == "moderator")) {
		delete = true
	} else {
		delete = false
	}

	currentURL := r.URL.RequestURI()

	data := struct {
		CurrentURL    string
		Role          string
		Delete        bool
		Authenticated bool
		Username      string
		Post          models.Post
		Comments      []models.Comment
		LikeCount     int
		DislikeCount  int
	}{
		CurrentURL:    currentURL,
		Role:          "",
		Delete:        delete,
		Authenticated: authenticated,
		Username:      username,
		Post:          post,
		Comments:      comments,
		LikeCount:     likeCount,
		DislikeCount:  dislikeCount,
	}

	if userInfo != nil {
		data.Role = userInfo.Role
	}

	// Логирование информации
	log.Printf("Post ID: %d, Title: %s", post.ID, post.Title)
	log.Printf("Post reactions - Likes: %d, Dislikes: %d", likeCount, dislikeCount)
	log.Printf("Fetched comments: %d", len(comments))
	for _, comment := range comments {
		log.Printf("Comment ID %d - Likes: %d, Dislikes: %d", comment.ID, comment.LikeCount, comment.DislikeCount)
	}
	log.Println("Rendering template for post")

	// Рендеринг шаблона
	handlers.RenderTemplate(w, "post.html", data)
}
