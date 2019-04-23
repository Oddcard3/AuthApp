package jwttoken

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
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

// GetUserID gets user ID from token if it's valid
func GetUserID(tokenStr string) (userID string, err error) {
	token, err := state.tokenAuth.Decode(tokenStr)
	if err != nil {
		return
	}

	claimsMap, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("Failed to get claims from token")
		return
	}

	userIDItem, ok := claimsMap["user_id"]
	if !ok {
		err = errors.New("No claim for user ID")
	}
	userID, ok = userIDItem.(string)
	if !ok {
		err = errors.New("User ID claim incorrect format")
	}
	return
}

// AppAuthenticator JWT authenticator
func AppAuthenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		// TODO: check if user is active
		if token == nil || !token.Valid /*|| sessions.IsExpired(token.Raw)*/ {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		next.ServeHTTP(w, r)
	})
}
