package traincontrol

import (
	"fmt"
	"log"
	"strconv"
	"unicode"
)

type TrainControl struct {
	rec  <-chan string // receiving channel
	send chan<- string // sending channel

	Sensors  map[int]Sensor
	Switches map[string]Switch
	Signals  map[string]Signal
	Blocks   map[string]Block
}

type Sensor struct {
	State bool
}

type Switch struct{}

type Signal struct {
	State bool
}

type Block struct{}

func NewTrainControl(rec <-chan string, send chan<- string) *TrainControl {
	tc := &TrainControl{
		rec:  rec,
		send: send,

		Sensors:  make(map[int]Sensor),
		Switches: make(map[string]Switch),
		Signals:  make(map[string]Signal),
		Blocks:   make(map[string]Block),
	}

	go func() {
		for msg := range tc.rec {
			err := tc.interpret(msg)
			if err != nil {
				log.Println(fmt.Errorf("unable to interpret message '%s': %w", msg, err))
			}
		}
	}()

	return tc
}

func (tc *TrainControl) String() string {
	return "This is the TrainControl"
}

func (tc *TrainControl) SetSwitch(id string, state string) {
	fmt.Printf("Switch '%s' changes to '%s'\n", id, state)
}

func (tc *TrainControl) interpret(msg string) error {
	switch {
	case len(msg) < 2:
		return fmt.Errorf("incomplete message, too short")
	case string(msg[len(msg)-1]) != "z":
		return fmt.Errorf("incomplete message, last character '%s'", string(msg[len(msg)-1]))
	case unicode.IsDigit(rune(msg[0])): // Sensor
		id, err := strconv.Atoi(msg[0:2])
		if err != nil {
			return fmt.Errorf("unable to decode id from '%s'", msg[0:2])
		}

		var state bool
		switch {
		case msg[2] == 'h':
			state = true
		case msg[2] == 'l':
			state = false
		default:
			return fmt.Errorf("unknown state '%d': %w", msg[2], err)
		}

		tc.Sensors[id] = Sensor{State: state}

	case msg[0] == 'y': // Switch
		return fmt.Errorf("Switch interpretation not implemented yet")

	case unicode.IsLetter(rune(msg[0])) && msg[0] != 'y': // Block
		return fmt.Errorf("Block interpretation not implemented yet")

	default:
		return fmt.Errorf("unable to interpret message '%s'", msg)
	}

	return nil
}

func (tc *TrainControl) Close() {

}
