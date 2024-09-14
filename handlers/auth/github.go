package auth

import (
	"encoding/json"
	"fmt"
	"forum/database"
	"forum/handlers"
	"forum/internal"
	"forum/models"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type UserGithub struct {
	ID   int    `json:"id"`
	Name string `json:"login"`
}

func GithubLogin(w http.ResponseWriter, r *http.Request) {
	authURL := "https://github.com/login/oauth/authorize"
	v := url.Values{}
	v.Set("client_id", internal.Cfg.GithubClientID)
	v.Set("redirect_uri", "http://localhost:8080/auth/github/callback")
	v.Set("response_type", "code")
	v.Set("scope", "user:email")
	http.Redirect(w, r, authURL+"?"+v.Encode(), http.StatusSeeOther)
}

func GitHubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := getGitHubToken(code)
	if err != nil {
		fmt.Fprintf(w, "Failed to get token: %s", err.Error())
		return
	}

	userInfo, err := getGitHubUserInfo(token)
	if err != nil {
		fmt.Fprintf(w, "Failed to get user info: %s", err.Error())
		return
	}

	var user UserGithub
	err = json.Unmarshal(userInfo, &user)
	if err != nil {
		fmt.Println("error", err)
	}
	fmt.Println(user)
	id := strconv.Itoa(user.ID)
	newUser := models.UserCreate{
		Name:     user.Name,
		Email:    id,
		Password: id,
	}

	if _, err1 := database.GetUserID(user.Name); err1 != nil {
		err = database.CreateUser(&newUser)
		if err != nil {
			log.Println(err)
		}

		for sessionID, sessionData := range handlers.SessionStore {
			if sessionData["username"] == user.Name {
				delete(handlers.SessionStore, sessionID)
			}
		}

		sessionID, err := internal.GenerateSessionID()
		if err != nil {
			log.Println(err)
			handlers.ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		handlers.SessionStore[sessionID] = map[string]string{"username": user.Name}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_id",
			Value:   sessionID,
			Expires: time.Now().Add(24 * time.Hour),
			Path:    "/",
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		for sessionID, sessionData := range handlers.SessionStore {
			if sessionData["username"] == user.Name {
				delete(handlers.SessionStore, sessionID)
			}
		}

		sessionID, err := internal.GenerateSessionID()
		if err != nil {
			log.Println(err)
			handlers.ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		handlers.SessionStore[sessionID] = map[string]string{"username": user.Name}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_id",
			Value:   sessionID,
			Expires: time.Now().Add(24 * time.Hour),
			Path:    "/",
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func getGitHubToken(code string) (string, error) {
	tokenURL := "https://github.com/login/oauth/access_token"
	v := url.Values{}
	v.Set("client_id", internal.Cfg.GithubClientID)
	v.Set("client_secret", internal.Cfg.GithubClientSecret)
	v.Set("code", code)
	v.Set("redirect_uri", "http://localhost:8080/auth/github/callback")

	req, err := http.NewRequest("POST", tokenURL, nil)
	if err != nil {
		return "", err
	}
	req.URL.RawQuery = v.Encode()
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
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

func getGitHubUserInfo(token string) ([]byte, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
