package app

import (
	"encoding/json"
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
	jsonEnc     json.Encoder
}

type Message struct {
	From string      `json:"from,omitempty"`
	To   string      `json:"to,omitempty"`
	Data interface{} `json:"data,omitempty"`
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

	log.Printf("%s connected to the websocket", conn.RemoteAddr().String())

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

func (ws *Websocket) SendToAll(msg Message) error {
	bMsg, err := json.Marshal(&msg)
	if err != nil {
		return fmt.Errorf("unable to marshal message: %w", err)
	}

	ws.SendToAllRaw(websocket.TextMessage, bMsg)

	return nil
}

func (ws *Websocket) SendToAllRaw(msgType int, msg []byte) {
	for _, conn := range ws.connections {
		if err := conn.WriteMessage(msgType, msg); err != nil {
			log.Println(fmt.Errorf("unable to send message to %s: %w", conn.UnderlyingConn().RemoteAddr(), err))

			// TODO: check errors and delete connection if neccessary

			// Remove connection
			err = ws.closeConn(conn)
			if err != nil {
				log.Println(fmt.Errorf("unable to remove connection: %w", err))
			}

			continue
		}
	}
}

func (ws *Websocket) sendToOthers(msgType int, msg []byte, sender *websocket.Conn) {
	for _, conn := range ws.connections {
		if conn == sender {
			continue
		}

		if err := conn.WriteMessage(msgType, msg); err != nil {
			log.Println(fmt.Errorf("unable to write message to %s: %w", conn.RemoteAddr().String(), err))

			// TODO: check errors and delete connection if neccessary

			// Remove connection
			err = ws.closeConn(conn)
			if err != nil {
				log.Println(fmt.Errorf("unable to remove connection: %w", err))
			}

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

func (ws *Websocket) closeConn(conn *websocket.Conn) error {
	defer ws.mu.Unlock()
	ws.mu.Lock()

	err := conn.Close()
	if err != nil {
		return fmt.Errorf("unable to close websocket connection to %s: %w", conn.RemoteAddr().String(), err)
	}

	for i, c := range ws.connections {
		if c == conn {
			ws.connections = append(ws.connections[:i], ws.connections[i+1:]...)
		}
	}

	// TODO: check if connection got removed

	return nil
}
