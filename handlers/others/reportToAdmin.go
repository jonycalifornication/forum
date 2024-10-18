package others

import (
	"forum/database"
	"forum/handlers"
	"log"
	"net/http"
	"strconv"
)

func ReportToAdmin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/report_to_admin" {
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

	postURL := r.FormValue("returnUrl")
	usernameModerator := user["username"]
	username := r.FormValue("username")
	postIDStr := r.FormValue("postId")
	reportReason := r.FormValue("reportReason")

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		log.Println("Error converting postId to int:", err)
		handlers.ErrorHandler(w, http.StatusBadRequest)
		return
	}

	if reportReason == "" {
		http.Error(w, "Reason for reporting is required", http.StatusBadRequest)
		return
	}

	err = database.SaveReportToAdmin(username, usernameModerator, postID, reportReason, postURL)
	if err != nil {
		log.Println("Error saving report to admin:", err)
		handlers.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, postURL, http.StatusSeeOther)
}
