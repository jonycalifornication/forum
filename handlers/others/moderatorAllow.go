package others

import (
	"fmt"
	"forum/database"
	"forum/handlers"
	"log"
	"net/http"
)

func ModeratorAllow(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/admin_page_allow" {
		handlers.ErrorHandler(w, http.StatusNotFound)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		log.Println("No username provided")
		handlers.ErrorHandler(w, http.StatusBadRequest) // Возвращаем 400 если нет имени пользователя
		return
	}

	fmt.Println("USERNAME:", username)

	err := database.UpdateUserRoleToModerator(username)
	if err != nil {
		log.Println("database.UpdateUserRoleToModerator", err)
		handlers.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin_page", http.StatusFound)
}

func ModeratorDeny(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/admin_page_deny" {
		handlers.ErrorHandler(w, http.StatusNotFound)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		log.Println("No username provided")
		handlers.ErrorHandler(w, http.StatusBadRequest)
		return
	}

	fmt.Println("USERNAME:", username)

	err := database.DenyUpdateUsertoModerator(username)
	if err != nil {
		log.Println(err)
		handlers.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin_page", http.StatusFound)
}

func DemoteToUser(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/admin_page_demote_to_user" {
		handlers.ErrorHandler(w, http.StatusNotFound)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		log.Println("No username provided")
		handlers.ErrorHandler(w, http.StatusBadRequest)
		return
	}

	err := database.DemoteToUser(username)
	if err != nil {
		log.Println(err)
		handlers.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin_page", http.StatusFound)
}
