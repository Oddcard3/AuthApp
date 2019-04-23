package api

import (
	jwttoken "authapp/auth/jwt"
	"net/http"

	"encoding/json"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type message struct {
	ID      string `json:"id"`
	Creator int    `json:"creator"`
	Users   []int  `json:"users"`
	Text    string `json:"text"`
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(userList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ChatRouter handler for posts
func ChatRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Group(func(r chi.Router) {
		r.Use(jwttoken.GetVerifier())
		r.Use(jwttoken.AppAuthenticator)
		r.Get("/messages/{chatId}", getMessages)
	})

	return r
}
