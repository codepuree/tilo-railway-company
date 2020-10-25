package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/codepuree/tilo-railway-company/internal/app"
	"github.com/codepuree/tilo-railway-company/pkg/scriptcontrol"
	"github.com/codepuree/tilo-railway-company/pkg/traincontrol"
	"github.com/codepuree/tilo-railway-company/pkg/trclib"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
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
			log.Printf("Arduino-> '%s'", string(m))
			server.Websocket.SendToAll(1, m)
		}
	}()
	a.Listen(ac)

	send := make(chan string)
	go func() {
		for msg := range send {
			// server.Websocket.SendToAll(1, []byte(msg))
			a.Write(msg)
		}
	}()

	arec := make(chan []byte)
	rec := make(chan string)
	go func() {
		defer close(rec)

		for m := range arec {
			rec <- string(m)
		}
	}()

	tcFile, err := os.Open("./config.json")
	if err != nil {
		log.Fatal(fmt.Errorf("unable to open config.json file: %w", err))
	}
	defer tcFile.Close()
	byteValue, err := ioutil.ReadAll(tcFile)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to read config.json file: %w", err))
	}
	var tcConfig traincontrol.TrainControlConfig
	err = json.Unmarshal(byteValue, &tcConfig)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to unmarshal config.json: %w", err))
	}

	log.Printf("Loaded configuration: %+v", tcConfig)

	trc := traincontrol.NewTrainControl(rec, send, tcConfig)
	interp := interp.New(interp.Options{})
	interp.Use(stdlib.Symbols)
	interp.Use(trclib.Symbols)
	sctrl, err := scriptcontrol.LoadFromFile(interp, "/var/www/custom/track1.go")
	if err != nil {
		log.Fatal(fmt.Errorf("unable to load track1: %w", err))
	}
	a.Listen(arec)

	go func() {
		l := make(chan string)
		server.Websocket.Listen(l)

		for msg := range l {
			if len(msg) > 3 && msg[0:2] == "s:" {
				name := msg[2:]
				log.Println("Starting Script function: ", name)
				Func, ok := sctrl[name]
				if !ok {
					continue
				}

				go Func.Func(trc)
			}
		}

	}()

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
