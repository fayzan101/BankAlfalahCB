package middleware

import (
	"net/http"
	"sync"
	"time"

	"bank-ai-chatbot/internal/api/handlers"
	apperrors "bank-ai-chatbot/pkg/errors"
	"bank-ai-chatbot/pkg/response"
)

type RateLimitConfig struct {
	Requests int
	Window   time.Duration
}

type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	limit    int
	window   time.Duration
}

type visitor struct {
	count       int
	windowStart time.Time
}

func NewRateLimiter(cfg RateLimitConfig) *RateLimiter {
	if cfg.Requests <= 0 {
		cfg.Requests = 60
	}
	if cfg.Window <= 0 {
		cfg.Window = time.Minute
	}

	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		limit:    cfg.Requests,
		window:   cfg.Window,
	}

	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	v, ok := rl.visitors[key]
	if !ok || now.Sub(v.windowStart) >= rl.window {
		rl.visitors[key] = &visitor{count: 1, windowStart: now}
		return true
	}

	if v.count >= rl.limit {
		return false
	}

	v.count++
	return true
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, v := range rl.visitors {
			if now.Sub(v.windowStart) >= rl.window {
				delete(rl.visitors, key)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) LimitByIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := "ip:" + r.RemoteAddr
		if !rl.Allow(key) {
			response.Error(w, apperrors.TooManyRequests("rate limit exceeded, try again later"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) LimitByUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := "ip:" + r.RemoteAddr
		if userID, ok := handlers.UserIDFromContext(r.Context()); ok {
			key = "user:" + userID.String()
		}

		if !rl.Allow(key) {
			response.Error(w, apperrors.TooManyRequests("rate limit exceeded, try again later"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
