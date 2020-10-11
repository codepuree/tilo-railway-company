package traincontrol

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"unicode"
)

type TrainControl struct {
	rec  <-chan string // receiving channel
	send chan<- string // sending channel

	// Sensors  map[int]*Sensor
	// Switches map[rune]*Switch
	// Signals  map[string]*Signal
	// Blocks   map[string]*Block
	TrainControlConfig
}

type TrainControlConfig struct {
	Sensors  map[int]*Sensor    `json:"sensors,omitempty"`
	Switches map[rune]*Switch   `json:"switches,omitempty"`
	Signals  map[string]*Signal `json:"signals,omitempty"`
	Blocks   map[string]*Block  `json:"blocks,omitempty"`
}

type Sensor struct {
	ID    int  `json:"id,omitempty"`
	State bool `json:"state,omitempty"`

	sync.Mutex
	listeners []chan bool
}

func (sns *Sensor) getChan() <-chan bool {
	sns.Lock()
	defer sns.Unlock()

	lst := make(chan bool)
	sns.listeners = append(sns.listeners, lst)

	return lst
}

func (sns *Sensor) close(lst <-chan bool) {
	sns.Lock()
	defer sns.Unlock()

	for i, l := range sns.listeners {
		if l == lst {
			sns.listeners = append(sns.listeners[:i], sns.listeners[i+1:]...)
			return
		}
	}
}

func (sns *Sensor) Await(state bool) {
	lst := sns.getChan()
	defer sns.close(lst)

	for s := range lst {
		if s == state {
			return
		}
	}
}

func (sns *Sensor) CountTo(state bool, amount int) {
	lst := sns.getChan()
	defer sns.close(lst)

	num := 0

	for s := range lst {
		if s == state {
			num++

			if num == amount {
				return
			}
		}
	}
}

func (sns *Sensor) setState(state bool) {
	if state != sns.State {
		sns.State = state

		for _, lst := range sns.listeners {
			lst <- state
		}
	}
}

type Switch struct {
	ID    rune `json:"id,omitempty"`
	State bool `json:"state,omitempty"`
}

type Signal struct {
	// ID
	State bool
}

type Block struct {
	ID    rune `json:"id,omitempty"`
	State bool `json:"state,omitempty"`
}

func NewTrainControl(rec <-chan string, send chan<- string, config TrainControlConfig) *TrainControl {
	tc := &TrainControl{
		rec:  rec,
		send: send,

		TrainControlConfig: config,
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
	case len(msg) != 4:
		return fmt.Errorf("incomplete message, incorrect number of symbols")

	case string(msg[3]) != "z":
		return fmt.Errorf("incomplete message, last character '%s'", string(msg[3]))

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

		tc.Sensors[id] = &Sensor{State: state}

	case msg[0] == 'y': // Switch
		return fmt.Errorf("Switch interpretation not implemented yet")

	case msg[0] == 'x': // Signal
		return fmt.Errorf("Signal interpretation not implemented yet")

	case unicode.IsLetter(rune(msg[0])) && msg[0] != 'y' && msg[0] != 'x': // Block
		return fmt.Errorf("Block interpretation not implemented yet")

	default:
		return fmt.Errorf("unable to interpret message '%s'", msg)
	}

	return nil
}

func (tc *TrainControl) Close() {

}
