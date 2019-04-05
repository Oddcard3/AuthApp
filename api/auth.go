package api

import (
	jwttoken "authapp/auth/jwt"
	"authapp/models/sessions"
	"authapp/models/users"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/google/jsonapi"
)

type loginRequest struct {
	ID       int    `jsonapi:"primary,logins"`
	Login    string `jsonapi:"attr,login"`
	Password string `jsonapi:"attr,password"`
}

type loginResponse struct {
	ID    int    `jsonapi:"primary, tokens"`
	Token string `jsonapi:"attr,token"`
}

func createSession(login string) (token string, e error) {
	expiredTs := time.Now().Add(time.Minute * 60)

	var tokenString string
	tokenString, e = jwttoken.New(login, expiredTs)
	if e != nil {
		return
	}
	if e = sessions.New(login, tokenString, expiredTs); e != nil {
		return
	}

	token = tokenString
	return
}

func login(w http.ResponseWriter, r *http.Request) {
	// TODO: log

	loginReq := new(loginRequest)

	if err := jsonapi.UnmarshalPayload(r.Body, loginReq); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", jsonapi.MediaType)

	u, e := users.Get(loginReq.Login)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
	if u == nil || !u.CheckPassword(loginReq.Password) {
		jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
			Title:  "Incorrect login or password",
			Detail: "Incorrect login or password",
			Status: "400",
			Meta:   nil,
		}})
		return
	}

	token, tokenErr := createSession(loginReq.Login)
	if tokenErr != nil {
		http.Error(w, tokenErr.Error(), http.StatusInternalServerError)
		return
	}

	loginResp := new(loginResponse)
	loginResp.Token = token
	if err := jsonapi.MarshalPayload(w, loginResp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
