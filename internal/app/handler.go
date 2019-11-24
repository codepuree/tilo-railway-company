package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// HandleIndex handles and serves the index endpoint
func HandleIndex(http.ResponseWriter, *http.Request) {
	log.Println("Handling index")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var conns []*websocket.Conn

func HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(fmt.Errorf("websocket error: %w", err))
		return
	}
	conns = append(conns, conn)

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		sendToAll(messageType, msg)
	}
}

func sendToAll(messageType int, msg []byte) {
	for _, conn := range conns {
		if err := conn.WriteMessage(messageType, msg); err != nil {
			log.Println(err)
			return
		}
	}
}
