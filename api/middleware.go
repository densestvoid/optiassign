package api

import (
	"context"
	"encoding/json"
	"net/http"
	"optiassign/domain"
)

// AuthMiddleware checks if the user is authenticated
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := getSession(r)
		if err != nil || session == nil || !session.IsValid() {
			// Redirect to login for HTML requests, return 401 for API requests
			if r.Header.Get("Accept") == "application/json" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		
		// Add session to request context
		ctx := context.WithValue(r.Context(), "session", session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// getSession retrieves the session from the request
func getSession(r *http.Request) (*domain.Session, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil, err
	}
	
	var session domain.Session
	if err := json.Unmarshal([]byte(cookie.Value), &session); err != nil {
		return nil, err
	}
	
	return &session, nil
}

// GetSessionFromContext retrieves the session from the request context
func GetSessionFromContext(r *http.Request) *domain.Session {
	if session, ok := r.Context().Value("session").(*domain.Session); ok {
		return session
	}
	return nil
}