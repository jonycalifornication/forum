package auth

import (
	"encoding/json"
	"fmt"
	"forum/database"
	"forum/handlers"
	"forum/internal"
	"forum/models"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
}

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	authURL := "https://accounts.google.com/o/oauth2/v2/auth"
	v := url.Values{}
	v.Set("client_id", internal.Cfg.GoogleClientID)
	v.Set("redirect_uri", "http://localhost:8080/auth/callback")
	v.Set("response_type", "code")
	v.Set("scope", "https://www.googleapis.com/auth/userinfo.profile")
	http.Redirect(w, r, authURL+"?"+v.Encode(), http.StatusSeeOther)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := getToken(code)
	if err != nil {
		log.Println("Failed to get token:", err)
		return
	}

	userInfo, err := getUserInfo(token)
	if err != nil {
		log.Println("Failed to get user info:", err)
		return
	}
	var user User
	err = json.Unmarshal(userInfo, &user)
	if err != nil {
		log.Println(err)
	}
	username := user.Name

	newUser := models.UserCreate{
		Name:     user.Name,
		Email:    user.ID,
		Password: user.ID,
	}

	if _, err1 := database.GetUserID(username); err1 != nil {
		err = database.CreateUser(&newUser)
		if err != nil {
			log.Println(err)
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
	} else {
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
	}
}

func getToken(code string) (string, error) {
	tokenURL := "https://oauth2.googleapis.com/token"
	v := url.Values{}
	v.Set("code", code)
	v.Set("client_id", internal.Cfg.GoogleClientID)
	v.Set("client_secret", internal.Cfg.GoogleClientSecret)
	v.Set("redirect_uri", "http://localhost:8080/auth/callback")
	v.Set("grant_type", "authorization_code")

	resp, err := http.PostForm(tokenURL, v)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResponse map[string]interface{}
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", err
	}

	token, ok := tokenResponse["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("no access token found")
	}

	return token, nil
}

func getUserInfo(token string) ([]byte, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
