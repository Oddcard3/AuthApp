package api

import (
	"time"
	"net/http"
	"strings"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"

	"github.com/spf13/viper"

	"authapp/logging"
)

func corsConfig() *cors.Cors {
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	return cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           86400, // Maximum value not ignored by any of major browsers
	})
}

// NewAppAPI creates HTTP router for the app
func NewAppAPI(enableCORS bool) (*chi.Mux, error) {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	// r.Use(middleware.RealIP)
	// there is issue with websocket if the compression is enabled
	// r.Use(middleware.DefaultCompress)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Timeout(15 * time.Second))

	//r.Use(logging.NewStructuredLogger(logger))
	r.Use(logging.NewStructuredLogger(logging.NewLogger()))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// use CORS middleware if client is not served by this api, e.g. from other domain or CDN
	if enableCORS {
		r.Use(corsConfig().Handler)
	}

	r.Mount("/auth", AuthRouter())
	r.Mount("/api/post/", PostsRouter())
	r.Mount("/api/users", UsersRouter())

	r.Mount("/ws", WebsocketRouter())
	// TODO: move from api?
	// r.Mount("/", http.Handle("/", http.FileServer(http.Dir("G://MYOWN//AuthAppUI//dist//AuthAppUI"))))

	uiDir := viper.GetString("config.uipath")
	fileServer(r, "/", http.Dir(uiDir))

	return r, nil
}

// fileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
