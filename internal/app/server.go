package app

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Server holds information on the server
type Server struct {
	Address string
	Router  *mux.Router
}

// NewServer creates a new server for hosting the website
func NewServer(address string) *Server {
	return &Server{
		Address: address,
		Router:  mux.NewRouter(),
	}
}

// Start sets up the server and starts it
func (s *Server) Start() error {

	err := http.ListenAndServe(s.Address, nil)
	if err != nil {
		return fmt.Errorf("unable to start server: %w", err)
	}

	return nil
}

// setupRoutes sets up all routes
func (s *Server) setupRoutes() {
	dir := "./web/static"
	s.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	s.Router.HandleFunc("/", HandleIndex)
}
