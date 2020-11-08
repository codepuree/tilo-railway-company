package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/codepuree/tilo-railway-company/internal/app"
	"github.com/codepuree/tilo-railway-company/pkg/scriptcontrol"
	"github.com/codepuree/tilo-railway-company/pkg/traincontrol"
	"github.com/codepuree/tilo-railway-company/pkg/trclib"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

var buildDate = "unknown"

func main() {
	var port string
	var baudRate int
	var address string

	flag.StringVar(&port, "serialPort", "COM3", "Serial port name, where the arduino is connected")
	flag.IntVar(&baudRate, "serialBaud", 9600, "Serial baud rate")
	flag.StringVar(&address, "address", ":8080", "Address of the server")

	flag.Parse()

	log.Println("Build date:", buildDate)

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
			err := a.Write(string(msg))
			if err != nil {
				log.Println(fmt.Errorf("unable to send message to arduino: %w", err))
			}
		}
	}()

	ac := make(chan []byte)
	go func() {
		for m := range ac {
			log.Printf("Arduino-> '%s'", string(m))
			server.Websocket.SendToAll(app.Message{From: "arduino", To: "all", Data: string(m)})
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
	msgc := make(chan traincontrol.Message)
	go func() {
		for msg := range msgc {
			err := server.Websocket.SendToAll(app.Message{From: "traincontrol", To: "all", Data: msg})
			if err != nil {
				log.Println(fmt.Errorf("unable to send message from traincontrol to all: %w", err))
			}
		}
	}()

	trc := traincontrol.NewTrainControl(rec, send, msgc, tcConfig)
	interp := interp.New(interp.Options{})
	interp.Use(stdlib.Symbols)
	interp.Use(trclib.Symbols)
	scriptPath := "/var/www/custom/track1.go"
	sctrl, err := scriptcontrol.LoadFromFile(interp, scriptPath)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to load track1: %w", err))
	}
	funcs := make([]string, 0, len(sctrl))
	for k := range sctrl {
		funcs = append(funcs, k)
	}
	log.Printf("Loaded %v functions from '%s':\n\t%+v", len(sctrl), scriptPath, funcs)
	a.Listen(arec)

	go func() {
		l := make(chan string)
		server.Websocket.Listen(l)

		for msg := range l {
			if len(msg) > 3 && msg[0:2] == "s:" {
				startParameter := strings.IndexRune(msg, '(')
				endParameter := strings.IndexRune(msg, ')')

				switch {
				case startParameter == -1:
					name := msg[2:]
					Func, ok := sctrl[name]
					if !ok {
						log.Println("unkown function ", name)
						continue
					}

					f, ok := Func.Func.(func(*traincontrol.TrainControl))
					if !ok {
						log.Println("unable to cast to func")
						continue
					}

					log.Printf("Starting function '%s'...", name)
					go f(trc)
				case startParameter > -1 && endParameter > -1:
					name := msg[2:startParameter]
					Func, ok := sctrl[name]
					if !ok {
						log.Println("unkown function ", name)
						continue
					}
					parameterRaw := msg[(startParameter + 1):endParameter]
					var parameter interface{}
					err := json.Unmarshal([]byte(parameterRaw), &parameter)
					if err != nil {
						log.Println("unable to parse parameters: ", parameterRaw)
					}

					switch p := parameter.(type) {
					case int:
						f, ok := Func.Func.(func(*traincontrol.TrainControl, int))
						if !ok {
							log.Println("unable to cast to func(tc, int)")
							continue
						}

						log.Printf("Starting function '%s(%v)'...", name, p)
						go f(trc, p)
					case float64:
						f, ok := Func.Func.(func(*traincontrol.TrainControl, int))
						if !ok {
							log.Println("unable to cast to func(tc, int)")
							continue
						}

						log.Printf("Starting function '%s(%v)'...", name, p)
						go f(trc, int(p))
					case string:
						f, ok := Func.Func.(func(*traincontrol.TrainControl, string))
						if !ok {
							log.Println("unable to cast to func(tc, int)")
							continue
						}

						log.Printf("Starting function '%s(%v)'...", name, p)
						go f(trc, p)
					default:
						log.Println("unsupported parameter type: ", reflect.ValueOf(p))
						continue
					}
				default:
					log.Println("Unable to parse message ", msg)
					continue
				}
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
