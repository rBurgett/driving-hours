package handlers

import (
	"net/http"

	"driving-hours/internal/auth"
	"driving-hours/internal/storage"
	"driving-hours/internal/templates"
)

type AuthHandler struct {
	storage   storage.Storage
	sessions  *auth.SessionManager
	renderer  *templates.Renderer
}

func NewAuthHandler(s storage.Storage, sm *auth.SessionManager, r *templates.Renderer) *AuthHandler {
	return &AuthHandler{
		storage:  s,
		sessions: sm,
		renderer: r,
	}
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	// If already logged in, redirect to appropriate dashboard
	user, _ := h.sessions.GetUserFromSession(r)
	if user != nil {
		if user.IsAdmin() {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/driver", http.StatusSeeOther)
		}
		return
	}

	h.renderer.Render(w, r, "auth/login.html", templates.Data{
		"Title": "Login",
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		h.renderer.Render(w, r, "auth/login.html", templates.Data{
			"Title": "Login",
			"Error": "Email and password are required",
			"Email": email,
		})
		return
	}

	user, err := h.storage.GetUserByEmail(email)
	if err != nil {
		h.renderer.Render(w, r, "auth/login.html", templates.Data{
			"Title": "Login",
			"Error": "An error occurred. Please try again.",
			"Email": email,
		})
		return
	}

	if user == nil {
		h.renderer.Render(w, r, "auth/login.html", templates.Data{
			"Title": "Login",
			"Error": "Invalid email or password",
			"Email": email,
		})
		return
	}

	valid, err := auth.VerifyPassword(password, user.PasswordHash)
	if err != nil || !valid {
		h.renderer.Render(w, r, "auth/login.html", templates.Data{
			"Title": "Login",
			"Error": "Invalid email or password",
			"Email": email,
		})
		return
	}

	if err := h.sessions.CreateSession(w, user.ID); err != nil {
		h.renderer.Render(w, r, "auth/login.html", templates.Data{
			"Title": "Login",
			"Error": "An error occurred. Please try again.",
			"Email": email,
		})
		return
	}

	if user.IsAdmin() {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/driver", http.StatusSeeOther)
	}
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	_ = h.sessions.DestroySession(w, r)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
