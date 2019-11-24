package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	s.setupRoutes()

	err := http.ListenAndServe(s.Address, s.Router)
	if err != nil {
		return fmt.Errorf("unable to start server: %w", err)
	}

	return nil
}

// setupRoutes sets up all routes
func (s *Server) setupRoutes() {
	dir := "./web/static"

	if _, err := os.Stat(dir); err == nil {
		// path/to/whatever exists
		log.Println("The path exists")
	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist
		log.Println("The path does not exists")
	}

	s.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	s.Router.PathPrefix("/resources/").Handler(http.StripPrefix("/resources/", http.FileServer(http.Dir("./resources"))))

	s.Router.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s wants to connect to the websocket.", r.RemoteAddr)
	})
	s.Router.HandleFunc("/", HandleIndex)
}
