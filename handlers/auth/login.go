package auth

import (
	"forum/database"
	"forum/handlers"
	"forum/internal"
	"log"
	"net/http"
	"time"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/sign_in" {
		handlers.ErrorHandler(w, http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		handlers.RenderTemplate(w, "sign_in.html", nil)
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			log.Println(err)
			handlers.ErrorHandler(w, http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		if err := database.AuthenticateUser(username, password); err != nil {
			data := struct{ ErrorMessage string }{err.Error()}
			handlers.RenderTemplate(w, "sign_in.html", data)
			return
		}

		for sessionID, sessionData := range handlers.SessionStore {
			if sessionData["username"] == username {
				delete(handlers.SessionStore, sessionID)
			}
		}

		sessionID, err := internal.GenerateSessionID()
		if err != nil {
			log.Println(err)
			handlers.ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		handlers.SessionStore[sessionID] = map[string]string{"username": username}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_id",
			Value:   sessionID,
			Expires: time.Now().Add(24 * time.Hour),
			Path:    "/",
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		handlers.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}
