package api

import (
	jwttoken "authapp/auth/jwt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/jsonapi"
	// "github.com/google/uuid"
)

type post struct {
	ID    string `jsonapi:"primary,post"`
	Title string `jsonapi:"attr,title"`
	Body  string `jsonapi:"attr,body"`
}

type posts struct {
	ID    int     `jsonapi:"primary,posts"`
	Posts []*post `jsonapi:"attr,items"`
}

func createFakePosts() []*post {
	l := make([]*post, 0)
	l = append(l, &post{"1", "Go became most popular language",
		"According to stackoverflow requests statistics Go is most popular language in the world"})
	l = append(l, &post{"2", "Barselona won Real Madrid",
		"Score is 3:0, Messi made hat trick!"})
	l = append(l, &post{"3", "Navalny is arrested again",
		"The opposition politician has been arrested during illegal meeting in the center of Moscow"})
	return l
}

var postList = createFakePosts()

func listHandler(w http.ResponseWriter, r *http.Request) {
	postsResp := &posts{1, postList}
	if err := jsonapi.MarshalPayload(w, postsResp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func add(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	p := new(post)

	// TODO: check if there is ID, return json api error

	if err := jsonapi.UnmarshalPayload(r.Body, p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	p.ID = postID

	postList = append(postList, p)
}

func get(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")

	var p *post
	for _, v := range postList {
		if v.ID == postID {
			p = v
			break
		}
	}
	if p == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if err := jsonapi.MarshalPayload(w, p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// PostsRouter handler for posts
func PostsRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Group(func(r chi.Router) {
		r.Use(jwttoken.GetVerifier())
		r.Use(jwttoken.AppAuthenticator)
		r.Get("/all", listHandler)
		r.Post("/{postID}", add)
		r.Get("/{postID}", get)
	})

	return r
}
