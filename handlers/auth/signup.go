package auth

import (
	"forum/database"
	"forum/handlers"
	"forum/internal"
	"forum/models"
	"log"
	"net/http"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/sign_up" {
		handlers.ErrorHandler(w, http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		handlers.RenderTemplate(w, "sign_up.html", nil)
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			log.Println(err)
			handlers.ErrorHandler(w, http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		var errorMessage string

		if err := internal.ValidateUsername(username); err != nil {
			errorMessage = err.Error()
		} else if err := internal.ValidateEmail(email); err != nil {
			errorMessage = err.Error()
		} else if err := internal.ValidatePassword(password); err != nil {
			errorMessage = err.Error()
		} else if err := database.CreateUser(&models.UserCreate{
			Name:     username,
			Email:    email,
			Password: password,
		}); err != nil {
			errorMessage = "User with this email or name already exists"
		}

		if errorMessage != "" {
			data := struct{ ErrorMessage string }{errorMessage}
			handlers.RenderTemplate(w, "sign_up.html", data)
		} else {
			log.Println("User registered successfully")
			http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
		}
	default:
		handlers.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}
