package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"optiassign/domain"
)

// AuthHandler handles authentication routes
type AuthHandler struct {
	authService *domain.AuthService
	oauthConfig *oauth2.Config
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService *domain.AuthService) *AuthHandler {
	config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &AuthHandler{
		authService: authService,
		oauthConfig: config,
	}
}

// Login initiates Google OAuth flow
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	state := generateRandomState()
	
	// Store state in session (simplified for MVS)
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production
		SameSite: http.SameSiteLaxMode,
		MaxAge:   600, // 10 minutes
	})
	
	url := h.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Callback handles Google OAuth callback
func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Verify state parameter
	state := r.URL.Query().Get("state")
	cookie, err := r.Cookie("oauth_state")
	if err != nil || cookie.Value != state {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}
	
	// Clear state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Authorization code not provided", http.StatusBadRequest)
		return
	}
	
	// Exchange code for token
	token, err := h.oauthConfig.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}
	
	// Get user info from Google
	userInfo, err := h.getUserInfo(ctx, token)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	
	// Get or create user
	user, err := h.authService.GetOrCreateUser(ctx, userInfo.ID, userInfo.Email, userInfo.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Authentication failed: %v", err), http.StatusForbidden)
		return
	}
	
	// Create session
	session := &domain.Session{
		UserID:    user.ID,
		Email:     user.Email,
		Name:      user.Name,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	
	// Store session (simplified for MVS)
	sessionData, _ := json.Marshal(session)
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    string(sessionData),
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400, // 24 hours
	})
	
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Logout clears the user session
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// getUserInfo retrieves user information from Google
func (h *AuthHandler) getUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	client := h.oauthConfig.Client(ctx, token)
	
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	
	return &userInfo, nil
}

// GoogleUserInfo represents user information from Google
type GoogleUserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// generateRandomState generates a random state parameter
func generateRandomState() string {
	// Simplified for MVS - in production use crypto/rand
	return fmt.Sprintf("%d", time.Now().UnixNano())
}