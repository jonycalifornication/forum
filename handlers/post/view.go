package post

import (
	"fmt"
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

	cookie, err := r.Cookie("session_id")

	var username string

	authenticated := false
	if err == nil {
		sessionID := cookie.Value
		sessionData, ok := handlers.SessionStore[sessionID]
		if ok {
			username = sessionData["username"]
			authenticated = true
		}
	}

	idStr := r.URL.Query().Get("id")
	id, err1 := strconv.Atoi(idStr)
	if err1 != nil {
		log.Println(err1)
		handlers.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	post, err := database.GetPostsById(id)
	if err != nil {
		log.Println(err1)
		handlers.ErrorHandler(w, http.StatusNotFound)
		return
	}

	fmt.Println("pppppppppppppppp", post)

	likeCount, dislikeCount, err := database.GetReactionCounts(id)
	if err != nil {
		log.Println("Error getting reaction counts:", err)
		likeCount = 0
		dislikeCount = 0
	}

	comments, err := database.GetCommentsByPostId(id, username)
	if err != nil {
		log.Println(err1)
		handlers.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	var delete bool
	if username == post.Username {
		delete = true
	} else {
		delete = false
	}

	data := struct {
		Delete        bool
		Authenticated bool
		Username      string
		Post          models.Post
		Comments      []models.Comment
		LikeCount     int
		DislikeCount  int
	}{
		Delete:        delete,
		Authenticated: authenticated,
		Username:      username,
		Post:          post,
		Comments:      comments,
		LikeCount:     likeCount,
		DislikeCount:  dislikeCount,
	}

	log.Println("Post ID:", id)
	log.Println("Fetched post:", post.Title)
	log.Printf("Post reactions - Likes: %d, Dislikes: %d", likeCount, dislikeCount)
	log.Printf("Fetched comments: %d", len(comments))
	for _, comment := range comments {
		log.Printf("Comment ID %d - Likes: %d, Dislikes: %d", comment.ID, comment.LikeCount, comment.DislikeCount)
	}
	log.Println("Rendering template for post")

	handlers.RenderTemplate(w, "post.html", data)
}
