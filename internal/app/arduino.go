package app

import (
	"bufio"
	"fmt"
	"log"
	"sync"

	"github.com/tarm/serial"
	// "go.bug.st/serial"
)

type Arduino struct {
	Port      string
	Baud      int
	conn      *serial.Port
	listeners []chan []byte
	mu        sync.RWMutex
}

func NewArduino(port string, baud int) *Arduino {
	return &Arduino{
		Port: port,
		Baud: baud,
	}
}

func (a *Arduino) Connect() error {
	conf := &serial.Config{
		Name: a.Port,
		Baud: a.Baud,
	}

	var err error
	a.conn, err = serial.OpenPort(conf)
	if err != nil {
		return err
	}

	go a.read()

	return nil
}

func (a *Arduino) Close() {
	defer a.conn.Close()

	for _, l := range a.listeners {
		close(l)
	}
}

func (a *Arduino) read() {
	r := bufio.NewReader(a.conn)

	for {
		// msg, _, err := r.ReadLine()
		msg, err := r.ReadString('z')
		if err != nil {
			err = fmt.Errorf("unable to read from arduino: %w", err)
			log.Println(err)
		}

		for _, l := range a.listeners {
			// go func(msg []byte, l chan []byte) {
			l <- []byte(msg)
			// }([]byte(msg), l)
		}
	}
}

func (a *Arduino) Listen(c chan []byte) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.listeners = append(a.listeners, c)
}

func (a *Arduino) Write(msg string) error {
	if len(msg) != 4 {
		return fmt.Errorf("the message '%s' to be send is to long", msg)
	}
	
	log.Printf("Arduino<- '%s'", msg)

	_, err := fmt.Fprint(a.conn, msg)
	if err != nil {
		err = fmt.Errorf("unable to write to Arduino: %w", err)
		return err
	}

	err = a.conn.Flush()
	if err != nil {
		err = fmt.Errorf("unable to flush message to Arduino: %w", err)
		return err
	}

	return nil
}
