package middleware

import (
	"net/http"

	"github.com/issue-notifier/issue-notifier-api/session"
)

func IsAuthenticated(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ses, err := session.Store.Get(r, session.CookieName)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else {
			// Check if user is Authenticated
			userSession, ok := ses.Values["UserSession"].(session.UserSession)

			if !userSession.IsAuthenticated || !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			} else {
				next.ServeHTTP(w, r)
			}
		}
	})
}
