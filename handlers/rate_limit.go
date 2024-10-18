package handlers

import (
	"forum/models"
	"html/template"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
}

type visitor struct {
	lastSeen map[string]time.Time
	count    map[string]int
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*visitor),
	}
}

func (rl *RateLimiter) getVisitor(ip string) *visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		v = &visitor{
			lastSeen: make(map[string]time.Time),
			count:    make(map[string]int),
		}
		rl.visitors[ip] = v

		go func() {
			time.Sleep(time.Minute)
			rl.mu.Lock()
			delete(rl.visitors, ip)
			rl.mu.Unlock()
		}()
	}
	return v
}

func (rl *RateLimiter) LimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		v := rl.getVisitor(ip)

		page := r.URL.Path

		lastSeen, exists := v.lastSeen[page]
		if !exists || time.Since(lastSeen) > time.Minute {

			v.count[page] = 1
		} else {
			v.count[page]++
		}
		v.lastSeen[page] = time.Now()

		if v.count[page] > 25 {
			rl.renderErrorPage(w, http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) renderErrorPage(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	tmpl, err := template.ParseFiles("web/html/error_page.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := models.Error{
		Code: statusCode,
		Text: "Too Many Requests. Please try again later.",
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
