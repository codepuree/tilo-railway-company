package traincontrol

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"sync"
	"time"
	"unicode"
)

type TrainControl struct {
	rec     <-chan string  // receiving channel
	send    chan<- string  // sending channel
	message chan<- Message // return other message
	sync.Mutex
	listeners []func(string)

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
	Trains   map[string]*Train  `json:"trains,omitempty"`
}

type Message interface{}

type Train struct {
	Name       string `json:"name,omitempty"`
	CrawlSpeed int    `json:"crawl_speed,omitempty"`
	MaxSpeed   int    `json:"max_speed,omitempty"`
	Accelerate struct {
		Step int           `json:"step,omitempty"`
		Time time.Duration `json:"time,omitempty"`
	} `json:"accelerate,omitempty"`
	Brake struct {
		Step int           `json:"step,omitempty"`
		Time time.Duration `json:"time,omitempty"`
	} `json:"brake,omitempty"`
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

// Changed returns the new sensor state as soon as it changes
func (sns *Sensor) Changed() bool {
	lst := sns.getChan()
	defer sns.close(lst)

	// TODO: could be a potential infinite loop that might be solved by adding a context with a timeout
	for s := range lst {
		return s
	}

	return sns.State
}

// Await waits for the sensor to change to the desired state
func (sns *Sensor) Await(state bool) {
	lst := sns.getChan()
	defer sns.close(lst)

	for s := range lst {
		if s == state {
			log.Printf("State changed to %v", s)
			return
		}
	}
}

// CountTo counts the sensor state changes up to the given amount and returns
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
	ID    rune        `json:"id,omitempty"`
	State SwitchState `json:"state,omitempty"`
}

type SwitchState rune

const (
	Straight SwitchState = '0'
	Bent                 = '1'
)

type Signal struct {
	// ID
	State bool
}

type Block struct {
	ID        string         `json:"id,omitempty"`
	Direction BlockDirection `json:"direction,omitempty"`
	Speed     int            `json:"speed,omitempty"`
	Sensors   []int          `json:"sensors,omitempty"` // TODO: custom unmarshal from ids to sensors
	// sensors   []*Sensor
	Distances []float64 `json:"distances,omitempty"`
	Train     Train     // The train that is currently in the block.
}

type BlockDirection rune

const (
	Forward          BlockDirection = 'f'
	Backward                        = 'b'
	Stopped                         = 's'
	EmergencyStopped                = 'x'
)

type Track struct {
	Blocks []*Block
}

func NewTrainControl(rec <-chan string, send chan<- string, message chan<- Message, config TrainControlConfig) *TrainControl {
	tc := &TrainControl{
		rec:     rec,
		send:    send,
		message: message,

		TrainControlConfig: config,
	}

	if tc.Blocks == nil {
		tc.Blocks = make(map[string]*Block)
	}
	if tc.Sensors == nil {
		tc.Sensors = make(map[int]*Sensor)
	}
	if tc.Switches == nil {
		tc.Switches = make(map[rune]*Switch)
	}
	if tc.Signals == nil {
		tc.Signals = make(map[string]*Signal)
	}

	go func() {
		for msg := range tc.rec {
			err := tc.interpret(msg)
			if err != nil {
				log.Println(fmt.Errorf("unable to interpret message '%s': %w", msg, err))
				continue
			}

			for _, lst := range tc.listeners {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("Recovered from:", r)
					}
				}()
				go lst(msg)
			}
		}
	}()

	// go func() {
	// 	tc.getSensorStates()
	// 	time.Sleep(5 * time.Second)
	// 	tc.getBlockDirections()
	// 	time.Sleep(5 * time.Second)
	// 	tc.getBlockSpeeds()
	// 	time.Sleep(5 * time.Second)
	// 	tc.getSignalStates()
	// }()

	return tc
}

func (tc *TrainControl) String() string {
	return "This is the TrainControl"
}

func (tc *TrainControl) SetSwitch(id string, state string) {
	fmt.Printf("Switch '%s' changes to '%s'\n", id, state)
	tc.send <- "y" + id + state + "z"
}

func (tc *TrainControl) SetBlockDirection(id string, state string) {
	tc.send <- id + "d" + state + "z"
}

func (tc *TrainControl) SetBlockSpeed(id string, speed int) {
	tc.send <- fmt.Sprintf("%s%02dz", id, speed)
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

		// tc.Sensors[id] = &Sensor{State: state}

		snr, ok := tc.Sensors[id]
		if !ok {
			return fmt.Errorf("sensor %v not initialized", id)
		}
		snr.setState(state)

	case msg[0] == 'y': // Switch
		// example: ya1z
		id := rune(msg[1])
		if _, ok := tc.Switches[id]; !ok {
			tc.Switches[id] = &Switch{ID: id}
		}

		switch {
		case rune(msg[2]) == '0':
			tc.Switches[id].State = Straight
		case rune(msg[2]) == '1':
			tc.Switches[id].State = Bent
		default:
			return fmt.Errorf("unkonw switch state '%s'", string(msg[2]))
		}

		return nil

	case msg[0] == 'x': // Signal
		return fmt.Errorf("Signal interpretation not implemented yet")

	case unicode.IsLetter(rune(msg[0])) && msg[0] != 'y' && msg[0] != 'x': // Block
		id := string(msg[0])
		if _, ok := tc.Blocks[id]; !ok {
			tc.Blocks[id] = &Block{ID: id}
		}

		switch {
		case rune(msg[1]) == 'd':
			switch {
			case rune(msg[2]) == 'f':
				tc.Blocks[id].Direction = Forward
			case rune(msg[2]) == 'b':
				tc.Blocks[id].Direction = Backward
			case rune(msg[2]) == 's':
				tc.Blocks[id].Direction = Stopped
			case rune(msg[2]) == 'x':
				tc.Blocks[id].Direction = EmergencyStopped
			default:
				return fmt.Errorf("unknown block direction '%s'", string(msg[2]))
			}

		case unicode.IsNumber(rune(msg[1])) && unicode.IsNumber(rune(msg[2])):
			speed, err := strconv.Atoi(string(msg[1:3]))
			if err != nil {
				return fmt.Errorf("unable to interprete the speed '%s'", string(msg[1:2]))
			}

			tc.Blocks[id].Speed = speed

		default:
			return fmt.Errorf("unable to interpret block message")
		}
		return nil

	default:
		return fmt.Errorf("unable to interpret message '%s'", msg)
	}

	return nil
}

func (tc *TrainControl) sendMessage(msg string) {
	tc.send <- msg
}

func (tc *TrainControl) sendMessageAwait(ctx context.Context, msg string) {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	func(ctx context.Context) {
		// tc.sendMessage(msg)
		// tc.waitMessage(msg)
		time.Sleep(1 * time.Second)
	}(ctx)
}

func (tc *TrainControl) waitMessage(msg string) {
	c := make(chan bool)
	f := func(newMsg string) {
		if msg == newMsg {
			c <- true
		}
	}

	tc.Lock()
	tc.listeners = append(tc.listeners, f)
	tc.Unlock()

	<-c
}

// GetActiveBlocks returns all the blocks were the power is turned on.
func (tc *TrainControl) GetActiveBlocks() []*Block {
	abs := make([]*Block, 0)

	for _, b := range tc.Blocks {
		if b.Speed > 0 && b.Direction != Stopped && b.Direction != EmergencyStopped {
			abs = append(abs, b)
		}
	}

	return abs
}

// GetOccupiedBlocks returns all blocks that are blocked, e.g. by trains
func (tc *TrainControl) GetOccupiedBlocks() []*Block {
	log.Println("Warn: 'GetOccupiedBlocks' is not implemented yet!")
	return nil
}

// GetActiveTrain returns the train that is currently running
func (tc *TrainControl) GetActiveTrain() *Train {
	// log.Println("Warn: 'GetActiveTrain' is not implemented yet!")
	return tc.Trains["N700"]
}

func (tc *TrainControl) Close() {
	v := int(math.Abs(float64(12)))
	time.Sleep(time.Duration(v) * time.Millisecond)
}

// ResetArduino sends `rstz` to the Arduino, this causes a software reset
func (tc *TrainControl) ResetArduino() {
	tc.send <- "rstz"
}

// PublishMessage publishes a desired message
func (tc *TrainControl) PublishMessage(msg Message) {
	tc.message <- msg
}

// getSensorStates sends `wsez` to retrive all the states of all the sensors
func (tc *TrainControl) getSensorStates() {
	tc.send <- "wsez"
}

// getSwitchStates sends `wswz` to retrive all the states of all the switches
func (tc *TrainControl) getSwitchStates() {
	tc.send <- "wswz"
}

// getBlockDirections sends `wblz` to retrive all the directions of all the blocks
func (tc *TrainControl) getBlockDirections() {
	tc.send <- "wblz"
}

// getBlockSpeeds sends `wspz` to retrive all the speeds of all the blocks
func (tc *TrainControl) getBlockSpeeds() {
	tc.send <- "wspz"
}

// getSignalStates sends `wsiz` to retrive all the states of all the signals
func (tc *TrainControl) getSignalStates() {
	tc.send <- "wsiz"
}
