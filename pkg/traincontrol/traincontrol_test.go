package traincontrol

import (
	"context"
	"sync"
	"testing"
	"time"
)

func echoServer(dur time.Duration) (rec chan string, send chan string) {
	rec = make(chan string)
	send = make(chan string)

	go func() {
		for msg := range send {
			time.Sleep(dur)
			rec <- msg
		}
	}()

	return rec, send
}

func TestTrainControl_interpret(t *testing.T) {
	tc := &TrainControl{
		TrainControlConfig: TrainControlConfig{
			Sensors:  make(map[int]*Sensor),
			Switches: make(map[rune]*Switch),
			Signals:  make(map[string]*Signal),
			Blocks:   make(map[rune]*Block),
		},
	}

	tests := []struct {
		name    string
		msg     string
		wantErr bool
	}{
		{"Incomplete message", "01l", true},
		{"empty message", "", true},
		{"Sensor", "01lz", true},
		{"Switch", "ya0z", true},
		{"Signal", "xa1z", false},
		{"Block direction", "adfz", true},
		{"Block speed", "a50z", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tc.interpret(tt.msg); (err != nil) != tt.wantErr {
				t.Errorf("TrainControl.interpret(\"%s\") error = %v, wantErr %v", tt.msg, err, tt.wantErr)
			}
		})
	}
}

func TestSensor_Close(t *testing.T) {
	sns := Sensor{}

	var lsts []<-chan bool

	for i := 0; i < 10; i++ {
		lsts = append(lsts, sns.getChan())
	}

	for _, lst := range lsts {
		sns.close(lst)
	}

	if len(sns.listeners) != 0 {
		t.Errorf("unable to close all listeners, there are still %d open", len(sns.listeners))
	}
}

func TestSensor_Await(t *testing.T) {
	awaited := false
	defer func() {
		if !awaited {
			t.Error("unable to await")
		}
	}()

	sns := Sensor{
		State: false,
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		sns.Await(true)
		awaited = true
	}()
	time.Sleep(time.Second)
	sns.setState(true)
	wg.Wait()
}

func TestSensor_setState(t *testing.T) {
	sns := Sensor{
		State: false,
	}

	sns.setState(true)

	if sns.State != true {
		t.Error("unable to set state")
	}
}

func TestSensor_CountTo(t *testing.T) {
	// type args struct {
	// 	state  bool
	// 	amount int
	// }
	// tests := []struct {
	// 	name string
	// 	sns  *Sensor
	// 	args args
	// }{
	// 	// TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		tt.sns.CountTo(tt.args.state, tt.args.amount)
	// 	})
	// }

}

func TestTrainControl_sendMessageAwait(t *testing.T) {
	rec, send := echoServer(100 * time.Millisecond)
	tc := NewTrainControl(rec, send, TrainControlConfig{
		Sensors:  map[int]*Sensor{19: {ID: 19, State: false}},
		Switches: make(map[rune]*Switch),
		Signals:  make(map[string]*Signal),
		Blocks:   make(map[rune]*Block),
	})

	ctx := context.Background()
	t.Log("Sending 19hz")
	tc.sendMessageAwait(ctx, "19hz")
	t.Log("Receiving 19hz")
}
