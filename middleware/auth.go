package middleware

import (
	"net/http"

	"github.com/issue-notifier/issue-notifier-api/session"
)

// IsAuthenticated is an authentication middleware to validate user session when calling authenticated routes
func IsAuthenticated(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ses, err := session.Store.Get(r, session.CookieName)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else {
			// Check if user is Authenticated
			// TODO validate session token, expiry time, etc.
			userSession, ok := ses.Values["UserSession"].(session.UserSession)

			if !userSession.IsAuthenticated || !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			} else {
				next.ServeHTTP(w, r)
			}
		}
	})
}
