package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
)

// Server holds information on the server
type Server struct {
	Address     string
	Router      *mux.Router
	EventSystem *EventSystem
	Websocket   *Websocket
}

// NewServer creates a new server for hosting the website
func NewServer(address string) *Server {
	return &Server{
		Address:     address,
		Router:      mux.NewRouter(),
		EventSystem: NewEventSystem(),
		Websocket:   NewWebsocket(),
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

// Stop shutdowns the server
func (s *Server) Stop() {
	s.EventSystem.Close()
}

// setupRoutes sets up all routes
func (s *Server) setupRoutes() {
	dir := "/var/www/static"

	if _, err := os.Stat(dir); err == nil {
		// path/to/whatever exists
		log.Println("The path exists")
	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist
		log.Println("The path does not exists")
	}

	s.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	s.Router.PathPrefix("/resources/").Handler(http.StripPrefix("/resources/", http.FileServer(http.Dir("./resources"))))

	msgc := make(chan []byte)
	go func() {
		defer close(msgc)

		for msg := range msgc {
			err := s.Websocket.SendToAll(Message{From: "eventsystem", To: "all", Data: string(msg)})
			if err != nil {
				log.Println(fmt.Errorf("unable to send message from event system to all: %w", err))
			}
		}
	}()
	s.EventSystem.Listen(msgc)

	s.Router.Handle("/websocket", s.Websocket)
	s.Router.HandleFunc("/", HandleIndex)

	// System
	s.Router.HandleFunc("/sys/shutdown", func(w http.ResponseWriter, r *http.Request) {
		err := exec.Command("shutdown", "-h", "now").Run()
		if err != nil {
			err = fmt.Errorf("unable to shutdown the system: %w", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	s.Router.HandleFunc("/sys/reboot", func(w http.ResponseWriter, r *http.Request) {
		err := exec.Command("shutdown", "-r", "now").Run()
		if err != nil {
			err = fmt.Errorf("unable to reboot the system: %w", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
