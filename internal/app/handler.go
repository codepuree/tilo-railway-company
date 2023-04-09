package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/codepuree/tilo-railway-company/internal/communication"
	"github.com/gorilla/websocket"
)

// HandleIndex handles and serves the index endpoint
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	log.Printf("Sending index to %s\n", r.RemoteAddr)
	index, err := ioutil.ReadFile("/var/www/static/index.html")
	if err != nil {
		log.Println("unable to load index file")
		return
	}

	_, err = w.Write(index)
	if err != nil {
		log.Println("unable to send index file")
		return
	}
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
	// onInitialConnection is a function that returns data that should be sent to client on their initial connection
	onInitialConnection func() []byte
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

	if ws.onInitialConnection != nil {
		err = conn.WriteMessage(0, ws.onInitialConnection())
		if err != nil {
			log.Println(fmt.Errorf("unable to send inital message to client: %w", err))
			return
		}
		log.Printf("%s got send the initial message", conn.RemoteAddr().String())
	}

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

func (ws *Websocket) SendToAll(msg communication.Message) error {
	bMsg, err := json.Marshal(&msg)
	if err != nil {
		return fmt.Errorf("unable to marshal message: %w", err)
	}

	ws.SendToAllRaw(websocket.TextMessage, bMsg)

	return nil
}

func sendTo(conn *websocket.Conn, msgType int, msg []byte) error {
	var err error
	defer func() {
		if a := recover(); a != nil {
			err = fmt.Errorf("panic during write message: %w", a)
		}
	}()

	err = conn.WriteMessage(msgType, msg)

	return err
}

func (ws *Websocket) SendToAllRaw(msgType int, msg []byte) {
	for _, conn := range ws.connections {
		if err := sendTo(conn, msgType, msg); err != nil {
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
