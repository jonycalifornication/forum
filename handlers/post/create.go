package post

import (
	"forum/database"
	"forum/handlers"
	"forum/internal"
	"forum/models"
	"log"
	"net/http"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/create_post" {
		handlers.ErrorHandler(w, http.StatusNotFound)
		return
	}

	// Проверка аутентификации
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

	authenticated := true
	username := user["username"]

	data := struct {
		Authenticated bool
		Username      string
		ErrorMessage  string
	}{
		Authenticated: authenticated,
		Username:      username,
		ErrorMessage:  "",
	}

	switch r.Method {
	case "GET":
		handlers.RenderTemplate(w, "create_post.html", data)
	case "POST":
		if err := r.ParseMultipartForm(20 << 20); err != nil {
			log.Println("Error parsing multipart form:", err)
			data.ErrorMessage = "The size of the image should not exceed 20mb"
			handlers.RenderTemplate(w, "create_post.html", data)
			return
		}

		log.Println("Form data received:", r.Form)

		categories := r.Form["categories[]"]
		if len(categories) < 1 {
			data.ErrorMessage = "Select at least one category."
			handlers.RenderTemplate(w, "create_post.html", data)
			return
		}

		title := internal.SanitizeComment(r.FormValue("title"))
		text := internal.SanitizeComment(r.FormValue("text"))

		var userfile models.File
		file, header, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			if header.Size > 20*1024*1024 {
				data.ErrorMessage = "File size exceeds 20MB limit."
				handlers.RenderTemplate(w, "create_post.html", data)
				return
			}
			userfile = models.File{
				FileGiven: file,
				Header:    header,
			}
		} else if err != http.ErrMissingFile {
			log.Println(err)
			data.ErrorMessage = "Add a image"
			handlers.RenderTemplate(w, "create_post.html", data)
			return
		}

		post := models.PostCreate{
			Username:   username,
			Title:      title,
			Text:       text,
			Categories: categories,
		}

		if err := database.CreatePost(&post, userfile); err != nil {
			log.Println(err)
			data.ErrorMessage = "Failed to create post. Please try again."
			handlers.RenderTemplate(w, "create_post.html", data)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		handlers.ErrorHandler(w, http.StatusMethodNotAllowed)
	}
}
