package custom

import "github.com/codepuree/tilo-railway-company/pkg/traincontrol"

// Track stops as soon as sensor 3 is reached
func Track(tc *traincontrol.TrainControl) {
	tc.Sensors[7].Await(false)
	tc.Sensors[3].Await(false)
	Stop(tc)
}

// Track1 sets the switches to the first track
func Track1(tc *traincontrol.TrainControl) {
	tc.SetSwitch("a", "0")
	tc.SetSwitch("b", "0")
	tc.SetSwitch("e", "0")
	tc.SetSwitch("f", "0")
}

// Stop stops all blocks
func Stop(tc *traincontrol.TrainControl) {
	tc.SetBlockDirection("b", "s")
	tc.SetBlockDirection("a", "s")
	tc.SetBlockDirection("c", "s")
	tc.SetBlockDirection("d", "s")
	tc.SetBlockDirection("e", "s")
	tc.SetBlockDirection("f", "s")
}
