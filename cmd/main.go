package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/codepuree/tilo-railway-company/internal/app"
)

func main() {
	var port string
	var baudRate int
	var address string

	flag.StringVar(&port, "serialPort", "COM3", "Serial port name, where the arduino is connected")
	flag.IntVar(&baudRate, "serialBaud", 9600, "Serial baud rate")
	flag.StringVar(&address, "address", ":8080", "Address of the server")

	flag.Parse()

	server := app.NewServer(address)
	a := app.NewArduino(port, baudRate)

	err := a.Connect()
	if err != nil {
		err = fmt.Errorf("unable to connect to Arduino: %w", err)
		log.Println(err)
	}

	amsgc := make(chan []byte)
	server.EventSystem.Listen(amsgc)

	go func() {
		defer close(amsgc)

		for msg := range amsgc {
			a.Write(string(msg))
		}
	}()

	ac := make(chan []byte)
	go func() {
		for m := range ac {
			log.Println("Arduino->", string(m))
			server.Websocket.SendToAll(1, m)
		}
	}()
	a.Listen(ac)

	mc := make(chan string)
	go func() {
		for m := range mc {
			a.Write(m)
		}
	}()
	server.Websocket.Listen(mc)

	err = server.Start()
	if err != nil {
		err = fmt.Errorf("unable to start the server: %w", err)
		log.Fatal(err)
	}
	defer server.Stop()
}
