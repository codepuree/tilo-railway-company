package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/codepuree/tilo-railway-company/internal/app"
	"github.com/codepuree/tilo-railway-company/internal/communication"
	"github.com/codepuree/tilo-railway-company/pkg/scriptcontrol"
	"github.com/codepuree/tilo-railway-company/pkg/traincontrol"
)

var buildDate = "unknown"

func main() {
	var port string
	var baudRate int
	var address string
	var configPath string

	flag.StringVar(&port, "serialPort", "COM3", "Serial port name, where the arduino is connected")
	flag.IntVar(&baudRate, "serialBaud", 9600, "Serial baud rate")
	flag.StringVar(&address, "address", ":8080", "Address of the server")
	flag.StringVar(&configPath, "config", "./config.json", "path to the config")

	flag.Parse()

	log.Println("Build date:", buildDate)

	server := app.NewServer(address)

	tcConfig, err := traincontrol.ConfigFromFile(configPath)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to load '%s': %w", configPath, err))
	}

	scriptFolder := "/var/www/custom/"
	layer, err := app.NewLayer(port, baudRate, &tcConfig, scriptFolder)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to initialize layer: %w", err))
	}

	// Send message from the event system to the Arduino
	amsgc := make(chan []byte)
	go func() {
		defer close(amsgc)

		for msg := range amsgc {
			err := layer.Arduino.Write(string(msg))
			if err != nil {
				log.Println(fmt.Errorf("unable to send message to arduino: %w", err))
			}
		}
	}()
	server.EventSystem.Listen(amsgc)

	// Send message from the Arduino to the web-socket clients
	ac := make(chan []byte)
	go func() {
		for m := range ac {
			log.Printf("Arduino-> '%s'", string(m))
			server.Websocket.SendToAll(communication.Message{From: "arduino", To: "all", Data: string(m)})
		}
	}()
	layer.Arduino.Listen(ac)

	sendToArduino := make(chan string)
	go func() {
		for msg := range sendToArduino {
			// server.Websocket.SendToAll(1, []byte(msg))
			layer.Arduino.Write(msg)
		}
	}()

	// Convert channel message from `[]byte` to `string`
	arec := make(chan []byte)
	receiveStringFromArduino := make(chan string)
	go func() {
		defer close(receiveStringFromArduino)

		for m := range arec {
			receiveStringFromArduino <- string(m)
		}
	}()

	log.Printf("Loaded configuration: %+v", tcConfig)
	msgc := make(chan traincontrol.Message)
	go func() {
		for msg := range msgc {
			err := server.Websocket.SendToAll(communication.Message{From: "traincontrol", To: "all", Data: msg})
			if err != nil {
				log.Println(fmt.Errorf("unable to send message from traincontrol to all: %w", err))
			}
		}
	}()

	layer.TrainControl.BindTrainControl(receiveStringFromArduino, sendToArduino, msgc)
	layer.ScriptControl.Bind(server.Websocket)
	defer layer.Close()

	// Execute main function
	mainFunction := "ControlRunner"
	if ctr, ok := layer.ScriptControl.Funcs[mainFunction]; ok {
		f, ok := ctr.Func.(func(*traincontrol.TrainControl))
		if ok {
			go f(layer.TrainControl)
		}
	}

	layer.Arduino.Listen(arec)

	go func() {
		l := make(chan string)
		server.Websocket.Listen(l)

		for msg := range l {
			err := evaluateFunctionCallMessage(msg, layer.ScriptControl, layer.TrainControl)
			if err != nil {
				log.Println(fmt.Errorf("unable to parse function call message '%s': %w", msg, err))
			}
		}
	}()

	mc := make(chan string)
	go func() {
		for m := range mc {
			layer.Arduino.Write(m)
		}
	}()
	server.Websocket.Listen(mc)

	err = server.Start()
	if err != nil {
		log.Fatal(fmt.Errorf("unable to start the server: %w", err))
	}
	defer server.Stop()
}

func evaluateFunctionCallMessage(msg string, sc *scriptcontrol.ScriptControl, trc *traincontrol.TrainControl) error {
	if len(msg) <= 3 {
		return fmt.Errorf("message to short")
	}

	if msg[0:2] != "s:" {
		return fmt.Errorf("message does not start with 's:'")
	}

	startParameter := strings.IndexRune(msg, '(')
	endParameter := strings.IndexRune(msg, ')')

	switch {
	case startParameter == -1:
		name := msg[2:]
		Func, ok := sc.Funcs[name]
		if !ok {
			return fmt.Errorf("unknown function '%s'", name)
		}

		f, ok := Func.Func.(func(*traincontrol.TrainControl))
		if !ok {
			return fmt.Errorf("unable to cast to func")
		}

		log.Printf("Starting function '%s'...", name)
		go f(trc)
	case startParameter > -1 && endParameter > -1:
		name := msg[2:startParameter]
		Func, ok := sc.Funcs[name]
		if !ok {
			return fmt.Errorf("unkown function '%s'", name)
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
				return fmt.Errorf("unable to cast to func(tc, int)")
			}

			log.Printf("Starting function '%s(%v)'...", name, p)
			go f(trc, p)
		case float64:
			f, ok := Func.Func.(func(*traincontrol.TrainControl, float64))
			if !ok {
				return fmt.Errorf("unable to cast to func(tc, float64)")
			}

			log.Printf("Starting function '%s(%v)'...", name, p)
			go f(trc, float64(p))
		case string:
			f, ok := Func.Func.(func(*traincontrol.TrainControl, string))
			if !ok {
				return fmt.Errorf("unable to cast to func(tc, string)")
			}

			log.Printf("Starting function '%s(%v)'...", name, p)
			go f(trc, p)
		default:
			return fmt.Errorf("unsupported parameter type '%s'", reflect.ValueOf(p))
		}
	default:
		return fmt.Errorf("Unable to parse message '%s'", msg)
	}

	return nil
}
