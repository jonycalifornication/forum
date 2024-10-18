package others

import (
	"forum/database"
	"forum/handlers"
	"log"
	"net/http"
	"strconv"
)

func SendReply(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/send_reply" {
		handlers.ErrorHandler(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		handlers.ErrorHandler(w, http.StatusMethodNotAllowed)
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

	moderatorUsername := r.FormValue("moderatorusername")
	admin := user["username"]
	postIDStr := r.FormValue("postid")
	replyText := r.FormValue("replyText")

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		log.Println("Error converting postId to int:", err)
		handlers.ErrorHandler(w, http.StatusBadRequest)
		return
	}

	if replyText == "" {
		http.Error(w, "Reason for reporting is required", http.StatusBadRequest)
		return
	}

	err = database.SaveReplyToAdmin(moderatorUsername, admin, postID, replyText)
	if err != nil {
		log.Println("Error saving report to admin:", err)
		handlers.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin_page", http.StatusSeeOther)
}

func DeleteReportFromAdminPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/delete_report_from_admin" {
		handlers.ErrorHandler(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		handlers.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	_, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
		return
	}

	idStr := r.FormValue("ID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting report ID:", err)
		handlers.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	err = database.DeleteReportByID(id)
	if err != nil {
		log.Println("Error deleting report:", err)
		handlers.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin_page", http.StatusSeeOther)
}

func DeleteReplyByID(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/delete_reply_from_admin" {
		handlers.ErrorHandler(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		handlers.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	_, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
		return
	}

	idStr := r.FormValue("ID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting report ID:", err)
		handlers.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	err = database.DeleteReplyByID(id)
	if err != nil {
		log.Println("Error deleting report:", err)
		handlers.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/user_profile", http.StatusSeeOther)
}
