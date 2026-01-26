package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"driving-hours/internal/models"
	"driving-hours/internal/storage"
)

const (
	SessionCookieName = "session"
	SessionDuration   = 7 * 24 * time.Hour // 7 days
	TokenLength       = 32
)

type SessionManager struct {
	storage  storage.Storage
	secure   bool
}

func NewSessionManager(storage storage.Storage, secure bool) *SessionManager {
	return &SessionManager{
		storage: storage,
		secure:  secure,
	}
}

// GenerateToken creates a cryptographically secure session token
func GenerateToken() (string, error) {
	b := make([]byte, TokenLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// CreateSession creates a new session for the user and sets the cookie
func (sm *SessionManager) CreateSession(w http.ResponseWriter, userID string) error {
	token, err := GenerateToken()
	if err != nil {
		return err
	}

	now := time.Now()
	session := &models.Session{
		Token:     token,
		UserID:    userID,
		ExpiresAt: now.Add(SessionDuration),
		CreatedAt: now,
	}

	if err := sm.storage.SaveSession(session); err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   sm.secure,
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}

// GetSession retrieves the current session from the request
func (sm *SessionManager) GetSession(r *http.Request) (*models.Session, error) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return nil, nil
	}

	return sm.storage.GetSession(cookie.Value)
}

// DestroySession removes the session and clears the cookie
func (sm *SessionManager) DestroySession(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return nil
	}

	if err := sm.storage.DeleteSession(cookie.Value); err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   sm.secure,
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}

// GetUserFromSession retrieves the user associated with the current session
func (sm *SessionManager) GetUserFromSession(r *http.Request) (*models.User, error) {
	session, err := sm.GetSession(r)
	if err != nil || session == nil {
		return nil, err
	}

	// Check admin first
	admin, err := sm.storage.GetAdmin()
	if err != nil {
		return nil, err
	}
	if admin != nil && admin.ID == session.UserID {
		return admin, nil
	}

	// Check regular users
	return sm.storage.GetUser(session.UserID)
}
