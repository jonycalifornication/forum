package auth

import (
	"net/http"
	"time"
)

func SignOut(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/sign_out" {
		http.Error(w, "Page does not exist", http.StatusNotFound)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   "",
		Expires: time.Now().Add(-24 * time.Hour),
		Path:    "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
