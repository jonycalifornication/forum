package handlers

import (
	"forum/models"
	"log"
	"net/http"
	"text/template"
)

var SessionStore = make(map[string]map[string]string)

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("web/html/" + tmpl)
	if err != nil {
		log.Printf("Error parsing template %s: %v", tmpl, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, data); err != nil {
		log.Printf("Error executing template %s: %v", tmpl, err)
		if !headersSent(w) {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func headersSent(w http.ResponseWriter) bool {
	if hw, ok := w.(interface{ WriteHeader(int) }); ok {
		return hw != nil
	}
	return false
}

func ErrorHandler(w http.ResponseWriter, code int) {
	data := models.Error{
		Text: http.StatusText(code),
		Code: code,
	}

	RenderTemplate(w, "error_page.html", data)
}
