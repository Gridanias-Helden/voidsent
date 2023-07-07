package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gridanias-helden/voidsent/internal/models"
	"github.com/gridanias-helden/voidsent/internal/services"
)

type sessionKey struct{}

var SessionKey sessionKey

func WithSession(manager services.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var session *models.Session

			c, err := r.Cookie("voidsent_session")
			if err == nil {
				session, err = manager.LoadSessionByID(r.Context(), c.Value)
				if err != nil {
					r = r.WithContext(context.WithValue(r.Context(), SessionKey, session))
				}
			}

			if err != nil && !strings.HasPrefix(r.RequestURI, "/auth") {
				// Log in first
				http.Redirect(w, r, "/auth/login", http.StatusTemporaryRedirect)
				return
			} else if err == nil && strings.HasPrefix(r.RequestURI, "/auth") {
				// Don't log in again ...
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
