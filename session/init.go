package session

import (
	"encoding/gob"

	"github.com/gorilla/sessions"
)

var (
	// Store will hold all session data
	Store *sessions.CookieStore // TODO: Move to token based authentication in future
	// CookieName cookie key which will be used to get the session from store
	CookieName string = "cookie-name"
)

// UserSession struct to store user session details
type UserSession struct {
	UserID          string `json:"userID"`
	AccessToken     string `json:"accessToken"`
	IsAuthenticated bool   `json:"isAuthenticated"`
}

// Init session store
func Init(sessionAuthKey string) {
	Store = sessions.NewCookieStore(
		[]byte(sessionAuthKey),
	)

	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 15,
		HttpOnly: true,
	}

	gob.Register(UserSession{})
}
