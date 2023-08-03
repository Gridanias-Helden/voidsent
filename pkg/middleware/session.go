package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gridanias-helden/voidsent/pkg/storage"
)

type sessionKey struct{}

var SessionKey sessionKey

func WithSession(sessions storage.Sessions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("voidsent_session")
			if (err != nil || c == nil) && !strings.HasPrefix(r.RequestURI, "/auth") {
				// No session, log in ...
				log.Printf("No session, redirecting to /auth/login/ ...")
				http.Redirect(w, r, "/auth/login/", http.StatusTemporaryRedirect)
				return
			}

			if (err != nil || c == nil) && strings.HasPrefix(r.RequestURI, "/auth") {
				next.ServeHTTP(w, r)
				return
			}

			session, err := sessions.SessionByID(r.Context(), c.Value)
			if err != nil && !strings.HasPrefix(r.RequestURI, "/auth") {
				// Invalid session, log in again ...
				log.Printf("Invalid session, redirecting to /auth/login/ ... %s / %+v", err, session)
				http.Redirect(w, r, "/auth/login/", http.StatusTemporaryRedirect)
				return
			}

			if err == nil && strings.HasPrefix(r.RequestURI, "/auth") {
				// Don't log in again ...
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), SessionKey, session))
			next.ServeHTTP(w, r)
		})
	}
}
