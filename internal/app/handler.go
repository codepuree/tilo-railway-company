package app

import (
	"log"
	"net/http"
)

// HandleIndex handles and serves the index endpoint
func HandleIndex(http.ResponseWriter, *http.Request) {
	log.Println("Handling index")
}
