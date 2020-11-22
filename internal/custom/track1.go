package custom

import (
	"log"
	"math"
	"time"

	"github.com/codepuree/tilo-railway-company/pkg/traincontrol"
)

//global variables here
var flag_direction int = 0
var flag_speedLock = 0
var flag_track int = 0
var flag_driveCircle = 1
var blocks [4]string

var direction string = ""
var previousTrack [4]string
var speed int = 0
var previousSpeed int = 0

func ManualControl(tc *traincontrol.TrainControl, direction string, speed int, blocks [4]string) {
	train := tc.GetActiveTrain()
	train = tc.Trains["N700"]

	if blocks == [4]string{"", "", "", ""} { //exit manual Control completely when only direction was set (and no blocks set until now)
		return
	}

	if flag_direction == 1 || flag_track != 0 {

		actualTrack := blocks
		if previousTrack != actualTrack && previousTrack != [4]string{"", "", "", ""} { //Partial reset of tracks (a,b,c,d) in case of track change
			log.Println("----------------Send Reset Command for previous Track to Arduino") //Track need to be set to stop and zero in case another track is choosen
			if previousTrack[0] != actualTrack[0] {                                         // in case one block changed while speedlock = 0
				PartialReset2Arduino(tc, previousTrack[0])               //reset block a,b,c or d (Direction and Speed)
				PartialSet2Arduino(tc, actualTrack[0], direction, speed) //Set Direction and Speed directly from user input
			}
			if previousTrack[1] != actualTrack[1] {
				PartialReset2Arduino(tc, previousTrack[1])               //reset block a,b,c or d (Direction and Speed)
				PartialSet2Arduino(tc, actualTrack[1], direction, speed) //Set Direction and Speed directly from user input
			}
			// } else {
			// 	PartialReset2Arduino(tc, previousTrack[0])               // while track change the previous track is resetted while the new state is written to new track
			// 	PartialReset2Arduino(tc, previousTrack[1])               // hard write happen only when no speed change was done (speedlock = 0) see below
			// 	PartialSet2Arduino(tc, actualTrack[0], direction, speed) // this is to prevent braking while driving and changing tracks
			// 	PartialSet2Arduino(tc, actualTrack[1], direction, speed)
			// }
		}

		previousDirection := string(tc.Blocks[[]rune(blocks[0])[0]].Direction) //gets direction requested from arduino. compare to last input
		actualDirection := direction                                           //send Direction
		if previousDirection != actualDirection {                              //Execution only by change
			log.Println("----------------Manual Control started. (Track & Direction)")
			log.Println("----------------Previous Direction was: ", previousDirection)
			log.Println("----------------Actual Direction is: ", actualDirection)
			for _, block := range blocks {
				Direction2Arduino(tc, block, actualDirection)
			}
		}

		if previousTrack != actualTrack { //send Track to Arduino
			log.Println("----------------Manual Control started. (Track & Direction)") //Execution only by change
			log.Println("----------------Previous Track was: ", previousTrack)
			log.Println("----------------Actual Track is: ", actualTrack)
			// if flag_driveCircle == 1 { //in case of do circle just ssend command once
			// 	Switches2Arduino(tc, blocks[0])
			// } else { //iterate through blocks array to set both tracks
			for _, block := range blocks { //send command to set junctions to new track
				Switches2Arduino(tc, block)
			}
			// }
			previousTrack = blocks //after track was set, store information in previous track for later comparision
		}
	}

	previousSpeed = tc.Blocks[[]rune(blocks[2])[0]].Speed // compare speed in junctions or open track in case junctions were switched in between. shall prevent intermediate full acceleration
	actualSpeed := speed
	if flag_direction == 1 && flag_track != 0 && previousSpeed != actualSpeed && flag_speedLock == 0 { //send Speed to Arduino
		log.Println("----------------Manual Control started. (Speed)")
		log.Println("----------------Send Speed Command to Arduino")
		log.Println("----------------Previous Speed was: ", previousSpeed)
		log.Println("----------------Actual Speed is: ", actualSpeed)

		if previousSpeed != 0 && actualSpeed == 0 { //Execution only by change
			flag_speedLock = 1
			log.Println("----------------Braking and Full Reset of Blocks ...")
			Brake2Arduino(tc, blocks, previousSpeed, actualSpeed, train.Brake.Step, train.Brake.Time)
			FullReset2Arduino(tc, blocks)
			tc.PublishMessage(struct {
				Speed int `json:"speed"`
			}{
				Speed: actualSpeed,
			}) //synchronize all websites with set state. in case sombody plaxed arpund during speedlock
			flag_speedLock = 0

		} else if previousSpeed < actualSpeed {
			flag_speedLock = 1
			log.Println("----------------Accelerating ...")
			Accelerate2Arduino(tc, blocks, previousSpeed, actualSpeed, train.Accelerate.Step, train.Accelerate.Time)
			tc.PublishMessage(struct {
				Speed int `json:"speed"`
			}{
				Speed: actualSpeed,
			}) //synchronize all websites with set state. in case sombody plaxed arpund during speedlock
			flag_speedLock = 0

		} else if previousSpeed > actualSpeed {
			flag_speedLock = 1
			log.Println("----------------Braking ...")
			Brake2Arduino(tc, blocks, previousSpeed, actualSpeed, train.Brake.Step, train.Brake.Time)
			tc.PublishMessage(struct {
				Speed int `json:"speed"`
			}{
				Speed: actualSpeed,
			}) //synchronize all websites with set state. in case sombody plaxed arpund during speedlock
			flag_speedLock = 0

		}

	}
}

func SetDirection(tc *traincontrol.TrainControl, d string) {
	if d == "f" {
		direction = "f"
	}
	if d == "b" {
		direction = "b"
	}
	log.Println("----------------Direction set: ", direction)

	flag_direction = 1

	tc.PublishMessage(struct {
		Direction string `json:"direction"`
	}{
		Direction: d,
	}) //synchronize all websites with set state

	ManualControl(tc, direction, speed, blocks)
}

func SetSpeed(tc *traincontrol.TrainControl, s int) {
	if flag_speedLock == 0 {
		speed = s
		log.Println("----------------Speed set: ", speed)
	}

	tc.PublishMessage(struct {
		Speed int `json:"speed"`
	}{
		Speed: s,
	}) //synchronize all websites with set state

	ManualControl(tc, direction, speed, blocks)
}

func SetTrack(tc *traincontrol.TrainControl, t int) {
	if flag_driveCircle == 1 && flag_speedLock == 0 { // "a,b,c,d" for Track 1-4, "f" for junctions, "g" for open terrain //Track change only possible if no speed was changed, no ramp is running
		if t == 1 {
			blocks = [4]string{"aw", "ao", "f", "g"}
			flag_track = 11 // flag_track: from track -> to track
			log.Println("----------------driveCircle: Blocks set: ", blocks, flag_track)

			tc.PublishMessage(struct {
				Blocks [4]string `json:"blocks"`
			}{
				Blocks: blocks,
			}) //synchronize all websites with set state
		}
		if t == 2 {
			blocks = [4]string{"bw", "bo", "f", "g"}
			flag_track = 22
			log.Println("----------------driveCircle: Blocks set: ", blocks, flag_track)

			tc.PublishMessage(struct {
				Blocks [4]string `json:"blocks"`
			}{
				Blocks: blocks,
			}) //synchronize all websites with set state
		}
		if t == 3 {
			blocks = [4]string{"cw", "co", "f", "g"}
			flag_track = 33
			log.Println("----------------driveCircle: Blocks set: ", blocks, flag_track)

			tc.PublishMessage(struct {
				Blocks [4]string `json:"blocks"`
			}{
				Blocks: blocks,
			}) //synchronize all websites with set state
		}
		if t == 4 {
			blocks = [4]string{"dw", "do", "f", "g"}
			flag_track = 44
			log.Println("----------------driveCircle: Blocks set: ", blocks, flag_track)

			tc.PublishMessage(struct {
				Blocks [4]string `json:"blocks"`
			}{
				Blocks: blocks,
			}) //synchronize all websites with set state
		}
	}

	ManualControl(tc, direction, speed, blocks)
}

func Brake2Arduino(tc *traincontrol.TrainControl, blocks [4]string, start int, target int, step int, dur time.Duration) {
	start_time := time.Now()
	i := start

	log.Println("----------------Brake Ramp started")
	for i >= target {
		for _, block := range blocks { //for each block do action
			tc.SetBlockSpeed(string(block[0]), i)
		}
		time.Sleep(dur * time.Millisecond)
		i = i - step
	}

	end_time := time.Now()
	brake_duration := end_time.Sub(start_time).Seconds()
	log.Println("----------------Brake Ramp done after: ", brake_duration)
}

func Accelerate2Arduino(tc *traincontrol.TrainControl, blocks [4]string, start int, target int, step int, dur time.Duration) {
	start_time := time.Now()

	log.Println("----------------Acceleration Ramp started")
	for i := start; i <= target; i = i + step {
		for _, block := range blocks { //for each block do action
			tc.SetBlockSpeed(string(block[0]), i)
		}
		time.Sleep(dur * time.Millisecond)
	}

	end_time := time.Now()
	acceleration_duration := end_time.Sub(start_time).Seconds()
	log.Println("----------------Acceleration Ramp done after: ", acceleration_duration)
}

func PartialReset2Arduino(tc *traincontrol.TrainControl, block string) {

	log.Println("----------------Track changed. Reset for block: ", block)
	tc.SetBlockSpeed(string(block[0]), 0)
	tc.SetBlockDirection(string(block[0]), "s")
}

func PartialSet2Arduino(tc *traincontrol.TrainControl, block string, direction string, speed int) {

	log.Println("----------------Track changed. Set for block: ", block)
	tc.SetBlockSpeed(string(block[0]), speed)
	tc.SetBlockDirection(string(block[0]), direction)
}

func FullReset2Arduino(tc *traincontrol.TrainControl, blocks [4]string) {

	log.Println("----------------Reset started for blocks: ", blocks)
	for _, block := range blocks {
		tc.SetBlockSpeed(string(block[0]), 0)
		tc.SetBlockDirection(string(block[0]), "s")
	}

	log.Println("----------------Reset done")
}

func Switches2Arduino(tc *traincontrol.TrainControl, block string) {
	if block == "aw" {
		log.Println("----------------Send Switches for Track 1 west in-/outbound to Arduino")
		tc.SetSwitch("a", "0")
		tc.SetSwitch("b", "0")
		tc.SetSwitch("e", "0")
		tc.SetSwitch("f", "0")
	}
	if block == "ae" {
		log.Println("----------------Send Switches for Track 1 to east in-/outbound Arduino")
		tc.SetSwitch("a", "0")
		tc.SetSwitch("b", "0")
		tc.SetSwitch("e", "0")
		tc.SetSwitch("f", "0")
	}
	if block == "bw" {
		log.Println("----------------Send Switches for Track 2 west in-/outbound to Arduino")
		tc.SetSwitch("a", "0")
		tc.SetSwitch("b", "1")
		tc.SetSwitch("d", "0")
		tc.SetSwitch("e", "1")
		tc.SetSwitch("f", "0")
	}
	if block == "be" {
		log.Println("----------------Send Switches for Track 2 east in-/outbound to Arduino")
		tc.SetSwitch("a", "0")
		tc.SetSwitch("b", "1")
		tc.SetSwitch("d", "0")
		tc.SetSwitch("e", "1")
		tc.SetSwitch("f", "0")
	}
	if block == "cw" {
		log.Println("----------------Send Switches for Track 3 west in-/outbound to Arduino")
		tc.SetSwitch("a", "1")
		tc.SetSwitch("c", "0")
		tc.SetSwitch("d", "1")
		tc.SetSwitch("e", "1")
		tc.SetSwitch("f", "0")
	}
	if block == "ce" {
		log.Println("----------------Send Switches for Track 3 east in-/outbound to Arduino")
		tc.SetSwitch("a", "1")
		tc.SetSwitch("c", "0")
		tc.SetSwitch("d", "1")
		tc.SetSwitch("e", "1")
		tc.SetSwitch("f", "0")
	}
	if block == "dw" {
		log.Println("----------------Send Switches for Track 4 west in-/outbound to Arduino")
		tc.SetSwitch("a", "1")
		tc.SetSwitch("c", "1")
		tc.SetSwitch("f", "1")
	}
	if block == "de" {
		log.Println("----------------Send Switches for Track 4 east in-/outbound to Arduino")
		tc.SetSwitch("a", "1")
		tc.SetSwitch("c", "1")
		tc.SetSwitch("f", "1")
	}
}

func Direction2Arduino(tc *traincontrol.TrainControl, block string, direction string) {
	tc.SetBlockDirection(string(block[0]), direction)
}

func EmergencyStop2Arduino(tc *traincontrol.TrainControl) {
	tc.SetBlockDirection("a", "s")
	tc.SetBlockDirection("b", "s")
	tc.SetBlockDirection("c", "s")
	tc.SetBlockDirection("d", "s")
	tc.SetBlockDirection("f", "s")
	tc.SetBlockDirection("g", "s")

	tc.SetBlockSpeed("a", 0)
	tc.SetBlockSpeed("b", 0)
	tc.SetBlockSpeed("c", 0)
	tc.SetBlockSpeed("d", 0)
	tc.SetBlockSpeed("f", 0)
	tc.SetBlockSpeed("g", 0)
}

// func StateFromArduino(tc *traincontrol.TrainControl) {
// 	tc.GetSensorStates()
// 	tc.GetBlockDirections()
// 	tc.GetBlockSpeeds()
// 	tc.GetSwitchStates()
// 	tc.GetSignalStates()
// }

//==============================================================================================================================
//==============================================================================================================================
//==============================================================================================================================
//============================================================ BELOW TEST CODE =================================================
//==============================================================================================================================
//==============================================================================================================================
//==============================================================================================================================

//==============================================================================================================================
//====================================================== S A M P L E S =========================================================
//==============================================================================================================================

// go get -u github.com/codepuree/tilo-railway-company/pkg/traincontrol			//get latest

//==============================================================================================================================
//available functions===========================================================================================================
//==============================================================================================================================
//to Server

// func ManualControl(tc *traincontrol.TrainControl, direction string, speed int, blocks [4]string)
// func SetDirection(tc *traincontrol.TrainControl, d string)
// func SetSpeed(tc *traincontrol.TrainControl, s int)
// func SetTrack(tc *traincontrol.TrainControl, t int)

//to Arduino
// func Brake2Arduino(tc *traincontrol.TrainControl, blocks [4]string, start int, target int, step int, dur time.Duration)
// func Accelerate2Arduino(tc *traincontrol.TrainControl, blocks [4]string, start int, target int, step int, dur time.Duration)
// func PartialReset2Arduino(tc *traincontrol.TrainControl, block string)
// func PartialSet2Arduino(tc *traincontrol.TrainControl, block string, direction string, speed int)
// func FullReset2Arduino(tc *traincontrol.TrainControl, blocks [4]string)
// func Switches2Arduino(tc *traincontrol.TrainControl, block string)
// func Direction2Arduino(tc *traincontrol.TrainControl, block string, direction string)
// func EmergencyStop2Arduino(tc *traincontrol.TrainControl)
// 	tc.SetBlockDirection("a", "s")	//send command to Arduino
// 	tc.SetBlockSpeed("a", 0)
// 	tc.SetSwitch("a", "1")
//==============================================================================================================================
//available States==============================================================================================================
//==============================================================================================================================
//	speed := tc.Blocks['a'].Speed
//	dir := tc.Blocks['a'].Direction
//	Step := tc.Trains["N700"].Accelerate.Step
//	Step_time := tc.Trains["N700"].Accelerate.Time
//	Max_speed := tc.Trains["N700"].MaxSpeed
//	log.Println("",Step, Step_time, Max_speed)
//
//	if tc.Blocks['a'].Direction == traincontrol.Forward {}	// do action when ...
//	if tc.Blocks['a'].Direction == 'f' {}
//  if tc.Sensors[15].State == false	// do action when state reached
//  tc.Sensors[15].Await(false)  // hold program until state reached
//  tc.Sensors[15].CountTo(10)	// hold program and do action when state reached
// 	tc.GetSensorStates()	//request latest states from arduino. send command to arduino
// 	tc.GetBlockDirections()	//request latest states from arduino. send command to arduino
// 	tc.GetBlockSpeeds()		//request latest states from arduino. send command to arduino
// 	tc.GetSwitchStates()	//request latest states from arduino. send command to arduino
// 	tc.GetSignalStates()	//request latest states from arduino. send command to arduino
// traincontrol.Signal{}.ID		// get Stetes from Server
// traincontrol.Signal{}.State
// traincontrol.Signal{}.Color
// traincontrol.Sensor{}.ID
// traincontrol.Sensor{}.State
// traincontrol.Switch{}.ID
// traincontrol.Switch{}.State
// traincontrol.Block{}.ID
// traincontrol.Block{}.Speed
// traincontrol.Block{}.Direction
// traincontrol.Train{}.Accelerate
// traincontrol.Train{}.Accelerate.Step
// traincontrol.Train{}.Accelerate.Time
// traincontrol.Train{}.Brake.Step
// traincontrol.Train{}.Brake.Time
// traincontrol.Train{}.CrawlSpeed
// traincontrol.Train{}.MaxSpeed
// traincontrol.Train{}.Name
// traincontrol.Train{}.Block
// traincontrol.Track{}.Blocks
//	log.Println("Direction gesetzt")

//to website
// tc.PublishMessage(tc.Trains["N700"])
// // tc.PublishMessage(tc.Trains["beliebiger String"])

// tc.PublishMessage(struct {
// 	Speed int `json:"speed"`
// }{
// 	Speed: actualSpeed,
// }) //synchronize all websites with set state. in case sombody plaxed arpund during speedlock

// Track stops as soon as sensor 3 is reached
func Track(tc *traincontrol.TrainControl) {
	tc.Sensors[15].Await(false)
	tc.SetBlockSpeed("0", 0)
	log.Println("Ab hier wird ausgeloest;")
	Brake_Ramp(tc)
	tc.Sensors[3].Await(false)
	Stop(tc)
}

func JunctionTrack3(tc *traincontrol.TrainControl) {
	tc.SetSwitch("d", "1")
	tc.SetSwitch("e", "1")
	tc.SetSwitch("f", "0")
	tc.SetSwitch("a", "1")
	tc.SetSwitch("c", "0")
}

func Track3(tc *traincontrol.TrainControl) {
	tc.SetBlockSpeed("0", 0)
	JunctionTrack3(tc)
	tc.SetBlockSpeed("0", 0)
	tc.SetBlockDirection("c", "f")
	tc.SetBlockDirection("f", "f")
	tc.SetBlockDirection("g", "f")
	Accelerate_Ramp(tc)
	tc.Sensors[19].CountTo(false, 1)
	tc.SetBlockSpeed("0", 0)
	Brake_Ramp(tc)
	tc.Sensors[3].Await(false)
	Stop(tc)

}

// Tracktest1 sets the switches to the first track
func Tracktest1(tc *traincontrol.TrainControl) {
	tc.Sensors[3].CountTo(false, 1)
	Track(tc)
}

// Stop stops all blocks
func Stop(tc *traincontrol.TrainControl) {
	tc.SetBlockDirection("a", "s")
	tc.SetBlockDirection("b", "s")
	tc.SetBlockDirection("c", "s")
	tc.SetBlockDirection("d", "s")
	tc.SetBlockDirection("e", "s")
	tc.SetBlockDirection("f", "s")
	tc.SetBlockDirection("g", "s")
	tc.SetBlockDirection("h", "s")
	tc.SetBlockDirection("i", "s")
	tc.SetBlockDirection("j", "s")
}

func Speed10(tc *traincontrol.TrainControl) {
	//time.Sleep(100 * time.Millisecond)
	tc.SetBlockSpeed("v", 10)
	tc.SetBlockSpeed("e", 10)
	tc.SetBlockSpeed("f", 10)
	tc.SetBlockSpeed("c", 10)
}

func Ramp(tc *traincontrol.TrainControl) {
	Ramp_target := 15
	Ramp_actual := 99
	Ramp_time := 2000.0

	tc.SetBlockSpeed("0", 0)
	Ramp_steps := (math.Abs(float64(Ramp_actual - Ramp_target))) / 3
	Ramp_step_duration := Ramp_time / Ramp_steps
	Ramp_actual_speed := Ramp_actual

	for i := 0; i < int(Ramp_steps); i++ {
		Ramp_actual_speed = Ramp_actual_speed - 3
		tc.SetBlockSpeed("c", Ramp_actual_speed)
		time.Sleep(time.Duration(Ramp_step_duration) * time.Millisecond)
	}
}

func Brake_Ramp_obsolete(tc *traincontrol.TrainControl) {
	Ramp_target := 35
	Ramp_actual := 75
	Ramp_step := 15
	Step_time := time.Duration(70) //milliseconds

	for i := Ramp_actual; i >= Ramp_target; i = i - Ramp_step {
		tc.SetBlockSpeed("c", i)
		time.Sleep(Step_time * time.Millisecond)
	}
}

func Brake_Ramp(tc *traincontrol.TrainControl) {
	Ramp_target := 37
	Ramp_actual := 75
	Ramp_step := 3
	Step_time := time.Duration(500) //milliseconds

	start_flag := 0
	end_flag := 0
	interrupt_flag := 0
	start := time.Now()
	end := time.Now()
	distance := 0.51 //0.16
	threshold_speed := 160.0
	i := Ramp_actual

	for i >= Ramp_target {
		log.Println("Ramp started")
		//log.Println(tc.Sensors[7].State)
		if tc.Sensors[19].State == false && start_flag == 0 {
			start = time.Now()
			start_flag = 1
		}
		if tc.Sensors[15].State == false && end_flag == 0 {
			end = time.Now()
			speed := (distance / end.Sub(start).Seconds()) * 3.6 * 160
			if speed > threshold_speed {
				log.Println("Ramp interrupted. Crawl Speed set. Train too fast:", speed)
				i = int(Ramp_target / 2)
				interrupt_flag = 1
			}
			end_flag = 1
		}

		tc.SetBlockSpeed("c", i)
		time.Sleep(Step_time * time.Millisecond)
		i = i - Ramp_step
	}
	if interrupt_flag == 1 {
		time.Sleep(400 * time.Millisecond)
		tc.SetBlockSpeed("c", Ramp_target)
	}
	log.Println("Ramp done")
}

func Accelerate_Ramp(tc *traincontrol.TrainControl) {
	Ramp_target := 75
	Ramp_actual := 25
	Ramp_step := 1
	Step_time := 250 //milliseconds

	for i := Ramp_actual; i <= Ramp_target; i = i + Ramp_step {
		tc.SetBlockSpeed("c", i)
		tc.SetBlockSpeed("g", i)
		tc.SetBlockSpeed("f", i)
		time.Sleep(time.Duration(Step_time) * time.Millisecond)
	}
}

// Measure mueasures how long one round takes
func Measure(tc *traincontrol.TrainControl) {
	distance := 0.545

	tc.Sensors[17].Await(false)
	start := time.Now()
	tc.Sensors[18].Await(false)
	end := time.Now()

	duration := end.Sub(start)
	float_duration := duration.Seconds()

	speed := (distance / float_duration) * 3.6 * 160

	log.Println("Der Abschnitt dauert:", duration)
	log.Println("Der Zug faehrt:", speed)
}

func Track1(tc *traincontrol.TrainControl) {
	tc.SetSwitch("a", "1")
}

func DirectionWest(tc *traincontrol.TrainControl) {
	tc.SetBlockDirection("f", "f")
	tc.SetBlockDirection("g", "f")
}

func Track1DirectionWest(tc *traincontrol.TrainControl) {
	tc.SetBlockDirection("a", "f")
}

func Track2DirectionWest(tc *traincontrol.TrainControl) {
	tc.SetBlockDirection("b", "f")
}

func Track3DirectionWest(tc *traincontrol.TrainControl) {
	tc.SetBlockDirection("c", "f")
}

func Track4DirectionWest(tc *traincontrol.TrainControl) {
	tc.SetBlockDirection("d", "f")
}

func DirectionEast(tc *traincontrol.TrainControl) {
	tc.SetBlockDirection("f", "b")
	tc.SetBlockDirection("g", "b")
}

func Track1DirectionEast(tc *traincontrol.TrainControl) {
	tc.SetBlockDirection("a", "b")
}

func Track2DirectionEast(tc *traincontrol.TrainControl) {
	tc.SetBlockDirection("b", "b")
}

func Track3DirectionEast(tc *traincontrol.TrainControl) {
	tc.SetBlockDirection("c", "b")
}

func Track4DirectionEast(tc *traincontrol.TrainControl) {
	tc.SetBlockDirection("d", "b")
}

//====================================================
//func testfuncs(tc *traincontrol.TrainControl) {

//speed := tc.Blocks['a'].Speed
//dir := tc.Blocks['a'].Direction
//
//if tc.Blocks['a'].Direction == tc.Forward {}
//if tc.Blocks['a'].Direction == 'f' {}

//	log.Println("Direction gesetzt")
//	Step := tc.Trains["N700"].Accelerate.Step
//	Step_time := tc.Trains["N700"].Accelerate.Time
//	Max_speed := tc.Trains["N700"].MaxSpeed
//	log.Println("",Step, Step_time, Max_speed)

//}
