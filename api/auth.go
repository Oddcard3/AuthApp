package api

import (
	jwttoken "authapp/auth/jwt"
	"authapp/db/models/users"
	"authapp/models/sessions"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

type credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type sessionToken struct {
	Token string `json:"token"`
}

type apiError struct {
	Code    int    `json:"error"`
	Message string `json:"message"`
}

func createSession(login string) (token string, e error) {
	expiredTs := time.Now().Add(time.Minute * 60)

	var tokenString string
	tokenString, e = jwttoken.New(login, expiredTs)
	if e != nil {
		return
	}
	// if e = sessions.New(login, tokenString, expiredTs); e != nil {
	// 	return
	// }

	token = tokenString
	return
}

func login(w http.ResponseWriter, r *http.Request) {
	// TODO: log

	var req credentials
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	authOk, err := users.Login(req.Login, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !authOk {
		if err = json.NewEncoder(w).Encode(apiError{100, "Incorrect login or password"}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	token, tokenErr := createSession(req.Login)
	if tokenErr != nil {
		http.Error(w, tokenErr.Error(), http.StatusInternalServerError)
		return
	}

	resp := new(sessionToken)
	resp.Token = token

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var userID string
	userID, err = sessions.GetUser(token.Raw)
	fmt.Println("token = " + token.Raw)
	if err == sessions.ErrNotFound {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = sessions.Logout(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// AuthRouter handler for auth (login/logout)
func AuthRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Post("/login", login)
	r.Group(func(r chi.Router) {
		r.Use(jwttoken.GetVerifier())
		r.Use(jwttoken.AppAuthenticator)
		r.Post("/logout", logout)
	})
	return r
}
