package comment

import (
	"fmt"
	"forum/database"
	"forum/handlers"
	"log"
	"net/http"
	"strconv"
)

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentIDStr := r.FormValue("commentId")

	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		log.Println("Error converting commentId to int:", err)
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	if err := database.DeleteCommentByID(commentID); err != nil {
		log.Println("Error deleting comment:", err)
		handlers.ErrorHandler(w, http.StatusInternalServerError)
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
