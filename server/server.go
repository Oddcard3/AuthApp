package server

import (
	"authapp/api"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	jwttoken "authapp/auth/jwt"

	"github.com/spf13/viper"
)

// Server server instance
type Server struct {
	*http.Server
}

// NewServer creates new server
func NewServer() (*Server, error) {
	port := viper.GetString("port")
	addr := ":" + port

	jwttoken.Init()

	api, err := api.NewAppAPI(false)
	if err != nil {
		return nil, err
	}

	srv := http.Server{
		Addr:    addr,
		Handler: api,
	}

	return &Server{&srv}, nil
}

// Start runs ListenAndServe on the http.Server with graceful shutdown.
func (srv *Server) Start() {
	//log.Println("starting server...")
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()
	fmt.Printf("Listening on %s\n", srv.Addr)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	//sig := <-quit
	//log.Println("Shutting down server... Reason:", sig)
	// teardown logic...

	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}
	log.Println("Server gracefully stopped")
}
