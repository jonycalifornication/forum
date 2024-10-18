package others

import (
	"forum/database"
	"forum/handlers"
	"log"
	"net/http"
)

func Apply(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/apply" {
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

	err = database.SendApplyModeratorRequest(username)
	if err != nil {
		log.Println(err)
		handlers.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/user_profile", http.StatusSeeOther)
}
