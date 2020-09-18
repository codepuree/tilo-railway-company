package app

import (
	"fmt"
	"log"
	"net/http"
	"sync"

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

type Websocket struct {
	connections []*websocket.Conn
	listeners   []chan string
	mu          sync.RWMutex
}

func NewWebsocket() *Websocket {
	return &Websocket{
		connections: make([]*websocket.Conn, 0),
		listeners:   make([]chan string, 0),
	}
}

func (ws *Websocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(fmt.Errorf("websocket error: %w", err))
		return
	}
	ws.connections = append(ws.connections, conn)

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		ws.sendToOthers(messageType, msg, conn)
		ws.publish(string(msg))
	}
}

func (ws *Websocket) SendToAll(msgType int, msg []byte) {
	for _, conn := range ws.connections {
		if err := conn.WriteMessage(msgType, msg); err != nil {
			log.Println(err)
			return
		}
	}
}

func (ws *Websocket) sendToOthers(msgType int, msg []byte, sender *websocket.Conn) {
	for _, conn := range ws.connections {
		if conn == sender {
			continue
		}

		if err := conn.WriteMessage(msgType, msg); err != nil {
			log.Println(err)
			continue
		}
	}
}

func (ws *Websocket) Listen(lc chan string) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	ws.listeners = append(ws.listeners, lc)
}

func (ws *Websocket) publish(msg string) {
	for _, lc := range ws.listeners {
		lc <- msg
		// go func(msg string, lc chan string) {
		// 	select {
		// 	case lc <- msg:
		// 		break
		// 	default:
		// 		log.Println("unable to send message to websocket listener")
		// 	}
		// }(msg, lc)
	}
}
