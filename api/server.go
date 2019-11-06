package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/pelletier/go-toml"
)

// Server provides an http.Server
type Server struct {
	*http.Server
}

// Creates and configures an API server
func NewServer(config *toml.Tree) (*Server, error) {
	log.Println("Configuring ...")

	api, err := New(config)
	if err != nil {
		return nil, err
	}

	addr := config.Get("http.addr").(string)
	port := config.Get("http.port").(string)
	srv := http.Server{
		Addr:    addr + ":" + port,
		Handler: api,
	}

	return &Server{&srv}, nil
}

// Runs ListenAndServe on the http.Server with graceful shutdown.
func (srv *Server) Start() {
	log.Println("Starting...")
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()
	log.Printf("Listening on %s\n", srv.Addr)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	log.Println("Shutting down... Reason:", sig)

	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}
	log.Println("Gracefully stopped")
}
