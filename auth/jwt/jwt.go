package jwttoken

import (
	"net/http"
	"time"

	"authapp/models/sessions"

	"github.com/go-chi/jwtauth"
)

type jwtState struct {
	tokenAuth *jwtauth.JWTAuth
}

var state *jwtState

// Init initializes token subsystem
func Init() error {
	state = new(jwtState)
	state.tokenAuth = jwtauth.New("HS256", []byte("fdsfvsdvsvxcvxfgsdgsdfgsd"), nil)
	return nil
}

// New creates new token for userID
func New(userID string, expirationTime time.Time) (res string, e error) {
	_, res, e = state.tokenAuth.Encode(jwtauth.Claims{"user_id": userID, "expired": expirationTime})
	return
}

//GetVerifier return token verifier middleware
func GetVerifier() func(http.Handler) http.Handler {
	return jwtauth.Verifier(state.tokenAuth)
}

// AppAuthenticator JWT authenticator
func AppAuthenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		if token == nil || !token.Valid || sessions.IsExpired(token.Raw) {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		next.ServeHTTP(w, r)
	})
}
