package others

import (
	"forum/database"
	"forum/handlers"
	"forum/models"
	"log"
	"net/http"
)

func AdminPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/admin_page" {
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

	currentUser, err := database.GetUserInfoByUsername(username)
	if err != nil {
		log.Println("Error retrieving current user info:", err)
		handlers.ErrorHandler(w, http.StatusNotFound)
		return
	}

	if currentUser.Role != "admin" {
		log.Printf("Unauthorized access attempt by user %s\n", username)
		handlers.ErrorHandler(w, http.StatusForbidden)
		return
	}

	requests, err := database.GetApplyModeratorRequest()
	if err != nil {
		log.Println("Error retrieving moderator requests:", err)
		handlers.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	userInfo, err := database.GetAllUsers()
	if err != nil {
		log.Println("Error retrieving users info:", err)
		handlers.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	posts, err := database.GetAllPosts()
	if err != nil {
		log.Println("Error retrieving posts info:", err)
		handlers.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	reports, err := database.GetAllReports()
	if err != nil {
		log.Println("Error retrieving reports info:", err)
		handlers.ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	data := struct {
		Requests []models.ApplyModeratorRequest
		UserInfo []models.User
		Posts    []models.Post
		Reports  []models.Report
	}{
		Requests: requests,
		UserInfo: userInfo,
		Posts:    posts,
		Reports:  reports,
	}

	switch r.Method {
	case "GET":
		handlers.RenderTemplate(w, "admin_page.html", data)
	default:
		handlers.ErrorHandler(w, http.StatusMethodNotAllowed)
	}
}
