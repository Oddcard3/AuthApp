package api

import (
	jwttoken "authapp/auth/jwt"
	"net/http"

	"encoding/json"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func createFakeUsers() []*user {
	l := make([]*user, 0)
	l = append(l, &user{"1", "Dmitry"})
	l = append(l, &user{"2", "Andrey"})
	return l
}

var userList = createFakeUsers()

func usersList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(userList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// UsersRouter handler for posts
func UsersRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Group(func(r chi.Router) {
		r.Use(jwttoken.GetVerifier())
		r.Use(jwttoken.AppAuthenticator)
		r.Get("/", usersList)
	})

	return r
}
