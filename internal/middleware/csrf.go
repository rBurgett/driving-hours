package middleware

import (
	"net/http"

	"github.com/gorilla/csrf"
)

// CSRFProtect returns middleware that protects against CSRF attacks
func CSRFProtect(key []byte, secure bool) func(http.Handler) http.Handler {
	return csrf.Protect(
		key,
		csrf.Secure(secure),
		csrf.Path("/"),
		csrf.SameSite(csrf.SameSiteLaxMode),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Invalid CSRF token", http.StatusForbidden)
		})),
	)
}

// CSRFToken returns the CSRF token for the current request
func CSRFToken(r *http.Request) string {
	return csrf.Token(r)
}

// CSRFTemplateField returns the CSRF template field for forms
func CSRFTemplateField(r *http.Request) string {
	return string(csrf.TemplateField(r))
}
