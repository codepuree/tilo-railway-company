package custom

import (
	"log"
	"math"
	"math/rand"
	"reflect"
	"strings"
	"time"

	"github.com/codepuree/tilo-railway-company/pkg/traincontrol"
)

//useful command for terminal: journalctl --unit=trc --follow -o short-precise
// sudo systemctl status trc
// sudo systemctl start trc
// sudo systemctl stop trc
// sudo systemctl restart trc

//root@train:~/Schreibtisch/trc/Tilo-Railway-Company# GOOS=linux GOARCH=arm GOARM=7 go build -a -tags netgo -ldflags '-w' -o ../trc ./cmd/

//=====================================================================================================================================================================
//======================================================================== I N I T ====================================================================================
//=====================================================================================================================================================================

// EmptyBlock valid for specific scenario. Definition of all available blocks in that scenario down at section  B L O C K S
var EmptyBlock = [4]string{"", "", "f", "g"}

// EmptySensors valid for specific scenario. Definition of all available blocks in that scenario down at section  S E N S O R S
var EmptySensors = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// EmptyDistances valid for specific scenario. Definition of all available blocks in that scenario down at section  S E N S O R S
var EmptyDistances = []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}

// ini init value for SensorTimes. Used for section  A U T O M A T I C
var ini = time.Date(1986, time.May, 15, 15, 0, 0, 0, time.UTC)

// SensorTimes valid for specific scenario. Definition of all available blocks in that scenario down at section  A U T O M A T I C
var SensorTimes = []time.Time{ini, ini, ini, ini, ini, ini, ini, ini, ini, ini, ini, ini, ini, ini, ini, ini, ini, ini, ini}

// global variables here
var actualBlocks [4]string = EmptyBlock
var targetBlocks [4]string = EmptyBlock
var sensorList []int = EmptySensors
var distanceList []float64 = EmptyDistances
var actualDirection string = "s"
var targetDirection string = "s"
var actualSpeed int = 0
var targetSpeed int = 0
var previousSpeed int = 0
var lastAccelerateTick time.Time = time.Unix(0, 0)

//flags for automatiion and program selection / behavior
var auto = 0              // will start automatic mode in control section
var manual = 0            // will enable the simple/direct mode
var doCircle = 0          // used in SetTrack and override individual branch selection for in- and outbound
var autoSleepTime = 1500  // sleeptime in Ms after each iteration in automatic mode
var autoBrake = 0         // used to activate autoBrake. reset at the end OR in case of acceleration (SpeedDiff < 10)
var autoBrakeReleased = 0 // used in autoBrake. flag used to mark action is running.
var autoBrakeAbsolute = 0 // autoBrake release because of block not clear.
var setSpeedFlag = 0      // enable disable setspeed until next track/event is set
// for Direction Selection
var randomDirection = 0     // Start Stop RandomDirectionFunction
var randomDirectionFlag = 1 // disable random Direction for one iteration
// for Round Selection
var randomRounds = 0      // Start Stop RandomRounds Function
var randomRoundsFlag = 0  // enable disable setTrack until next track/event is set
var maxRounds = 1         // amount of maximal rounds to be driven
var minRounds = 1         // amount of minimal rounds to be driven
var rounds = 0            // internal maxrounds (combine with random rounds)
var roundsCounter = 0     // Counter for actual driven Rounds.
var roundsCounterFlag = 0 // enable disable Random/OrderRounds until next track/event is set
// for Track Selection
var randomTrack = 0     // Start Stop RandomTrack Function
var randomTrackFlag = 1 // disable random Track for one iteration
var trackValue = -1.0   // valid value will be set after first SetTrack()
// for RoundRobin
var doRoundRobin = 0                                                                                       // Defines wether ControlRoundRobin should be run over common Control
var initialRoundRobin = 0                                                                                  // Defines if first Round RoundRobin is executed
var trackOffset = 0                                                                                        // Defines if track 1 is unused for display
var roundRobinRunning = 0                                                                                  // Defines if RoundRobin is actual in progress
var progressCurrentTrain = 0                                                                               // Defines Progress of Current Train (Train actual in display)
var progressNextTrain = 0                                                                                  // Defines Progress of Next Train (Train next round in display)
var waitCounter = 0                                                                                        // Counter to wait for some time after progress finished (multiplied by 50 milliseconds from cycletime)
var nextTrack = 0                                                                                          // Defines Track(index) of Next Train
var currentTrack = 0                                                                                       // Defines Track(index) of Current Train
var currentTrain traincontrol.Train                                                                        // Train variable for current Train
var nextTrain traincontrol.Train                                                                           // Train variable for next Train
var roundRobinTracks = [4]string{"a", "b", "c", "d"}                                                       // Array for usable Blocks
var roundRobinTargetSpeeds = [4]int{0, 0, 0, 0}                                                            // Array for Target Speeds during RoundRobin
var roundRobinActualSpeeds = [4]int{0, 0, 0, 0}                                                            // Array for Actual Speeds during RoundRobin
var lastAccelerateTicks = [4]time.Time{time.Unix(0, 0), time.Unix(0, 0), time.Unix(0, 0), time.Unix(0, 0)} // AccelerateTimeTick per Block

// variables used for velocity measurment
var timeResetFlag = 0
var sensorPerRound = 7 // number of sensors per full round (used to get proper distances for last sensor when end of list reached)
var speedAverageList = []int{0, 0, 0, 0}

//=====================================================================================================================================================================
//================================================================== C O N T R O L / M A I N ==========================================================================
//=====================================================================================================================================================================

// ControlRunner performs an arduino loop with controlCycleDuration
func ControlRunner(tc *traincontrol.TrainControl) {

	const controlCycleDuration = 50 * time.Millisecond
	var lastControlCycle = time.Unix(0, 0)

	for {
		var now = time.Now()
		var waitTime = now.Sub(lastControlCycle)
		var sleepTime = controlCycleDuration - waitTime
		if sleepTime < 1 {
			if isDriveable() {
				if doRoundRobin > 0 {
					ControlRoundRobin(tc, tc.Trains["RoundRobin"]) //for RoundRobin only one train definition is used
				} else {
					Control(tc, tc.GetActiveTrain())
					//Control(tc, tc.Trains["N700"])
					velocity(tc)
				}
				//velocity(tc)
			}
			lastControlCycle = now
		} else {
			time.Sleep(sleepTime)
		}
	}
}

// Control is run in a short interval
func Control(tc *traincontrol.TrainControl, train *traincontrol.Train) {
	if targetDirection != actualDirection {
		actualDirection = targetDirection
		setBlocksDirection(tc, actualBlocks, targetDirection)
		setSensorList(tc, targetBlocks, targetDirection)
		TimeReset(tc)
		log.Println("----------------sensorList: ", sensorList)
		log.Println("----------------distanceList: ", distanceList)
	}

	if (targetSpeed != actualSpeed) && (targetBlocks == actualBlocks) { // check for block change because tracsk could be changed in between
		adjustSpeed(tc, train, actualBlocks, targetSpeed)
	}

	if autoBrakeReleased == 0 && autoBrakeAbsolute == 0 && (canSwitch(tc, targetBlocks, actualBlocks) || allTrainsStopped()) { // switch blocks and tracks only if no brake procedure in progress
		if targetBlocks != actualBlocks {
			actualBlocks = targetBlocks
			setSwitches(tc, targetBlocks)
			setBlocksDirection(tc, targetBlocks, actualDirection)
			setBlocksSpeed(tc, train, targetBlocks, actualSpeed)
			setSensorList(tc, targetBlocks, actualDirection)
			resetInactiveBlocks(tc, targetBlocks)
			TimeReset(tc)
		}
	}

	// in case of manual send switch command as soon as it will be received without any double checks
	if manual == 1 {
		if targetBlocks != actualBlocks {
			actualBlocks = targetBlocks
			setSwitches(tc, targetBlocks)
			setBlocksDirection(tc, targetBlocks, actualDirection)
			setBlocksSpeed(tc, train, targetBlocks, actualSpeed)
			setSensorList(tc, targetBlocks, actualDirection)
			resetInactiveBlocks(tc, targetBlocks)
			TimeReset(tc)
		}
	}

	// Brake Control while block not clear
	// if block not clear Auto Brake will be set at entrance of tunnel. full stop at station end
	if tc.Sensors[sensorList[7]].State == false && blockClear(actualBlocks) == false && autoBrake == 0 {
		SetBrake(tc, 3)
		autoBrakeAbsolute = 1
	}

	if autoBrakeAbsolute == 1 && allTrainsStopped() && blockClear(targetBlocks) { // resume with previous speed in case of absolute brake since tracks weren't clear. To break procedure in station see below
		targetSpeed = previousSpeed
		autoBrakeAbsolute = 0
		autoBrakeReleased = 0

		tc.PublishMessage(struct {
			Speed int `json:"speed"`
		}{
			Speed: targetSpeed,
		}) //synchronize all websites with set state
	}

	if autoBrake == 1 && allTrainsStopped() { // Break condition to reset speed whil in station. set speed zero in case brake was release while train was stopped as it was blocked
		targetSpeed = 0
		previousSpeed = 0
		autoBrake = 0
		autoBrakeAbsolute = 0
		if auto == 1 {
			SetAuto(tc, 0)
		}
	}

	//Automation
	// autoBrake gets activated via sensor. Decrease speed down to crawlspeed. Stops completely when defined sensor is released.
	if autoBrake > 0 {
		// release brake process
		// sensorList[7] tunnel entry after open rail (start sensor)
		// sensorList[8] station entry, first sensor after switches but befor station entry
		// if position [7] still released but [8] already release, prevent start of brake procedure as remaining track not long enough. brake next round.
		if tc.Sensors[sensorList[7]].State == false && tc.Sensors[sensorList[8]].State == true && autoBrakeReleased == 0 {
			previousSpeed = targetSpeed
			targetSpeed = train.CrawlSpeed
			autoBrakeReleased = 1
			log.Println("----------------Braking. Speed set: ", targetSpeed)

			tc.PublishMessage(struct {
				Speed int `json:"speed"`
			}{
				Speed: targetSpeed,
			}) //synchronize all websites with set state
		}

		// While autoBrake, break speed ramp (at position Sensor [8]) in case of to slow braking (can be modified below)
		var speedMuliplicator = 3                                                                                                   // MODIFY THIS LINE TO CHANGE SPEED BREAK BEHAVIOR (RELEASE EARLIER/LATER)
		if tc.Sensors[sensorList[8]].State == false && autoBrakeReleased == 1 && actualSpeed > speedMuliplicator*train.CrawlSpeed { // Set CrawlSpeed hard since brake procedure took too long
			log.Println("----------------Override brake ramp. CrawlSpeed set, since actualSpeed was: ", actualSpeed)
			actualSpeed = train.CrawlSpeed

			setBlocksSpeed(tc, train, actualBlocks, actualSpeed) //override brake ramp

			tc.PublishMessage(struct {
				ActualSpeed int `json:"actualspeed"`
			}{
				ActualSpeed: actualSpeed,
			}) //synchronize all websites with set state
		}

		// brake process finished at position [10], end of platform
		if tc.Sensors[sensorList[10]].State == false && autoBrakeReleased == 1 {
			targetSpeed = 0
			actualSpeed = targetSpeed // set both values same level to not release brake ramp
			log.Println("----------------Stop now. Speed set: ", targetSpeed)

			tc.PublishMessage(struct {
				Speed int `json:"speed"`
			}{
				Speed: targetSpeed,
			}) //synchronize all websites with set state

			tc.PublishMessage(struct {
				ActualSpeed int `json:"actualspeed"`
			}{
				ActualSpeed: actualSpeed,
			}) //synchronize all websites with set state

			setBlocksSpeed(tc, train, actualBlocks, actualSpeed) //override brake ramp
			TimeReset(tc)
			autoBrake = 0 //reset autobrake
			autoBrakeReleased = 0
		}

		// break condition in case of acceleration while autobrake is running. twice brake.step because it can easily overshoot while braking
		if (actualSpeed-targetSpeed < -2*train.Brake.Step) && (targetBlocks == actualBlocks) { // check for block to prevent track change while braking
			autoBrake = 0
			autoBrakeReleased = 0
			log.Println("----------------AutoBrake Reset. SpeedDiff: ", actualSpeed-targetSpeed)
		}

		if autoBrake == 1 && auto == 1 { // reset automatic mode, in case user manually set velocity to zero
			SetAuto(tc, 0)
		}

	}

	if auto == 1 {
		// ================================================================================================================ Set (Random) Values
		if randomRounds == 1 && randomRoundsFlag == 0 {
			rounds = int(setRandomRounds(tc, minRounds, maxRounds))
			randomRoundsFlag = 1
		} else {
			rounds = maxRounds
		}

		if randomDirection == 1 && randomDirectionFlag == 0 {
			setRandomDirection(tc)
			randomDirectionFlag = 1
		}

		if randomTrack == 1 && randomTrackFlag == 0 {
			setRandomTrack(tc)
			randomTrackFlag = 1
		}
		if randomTrack == 0 && randomTrackFlag == 0 {
			setOrderTrack(tc, trackValue)
			randomTrackFlag = 1
		}

		// Start new round after Reset
		if setSpeedFlag == 0 {
			if previousSpeed != 0 {
				SetSpeed(tc, previousSpeed)
			} else {
				SetSpeed(tc, targetSpeed)
			}

			setSpeedFlag = 1
		}
		// ================================================================================================================ Count Rounds
		if tc.Sensors[sensorList[6]].State == false && roundsCounterFlag == 0 { // increase rounds counter each round
			roundsCounter++
			log.Println("----------------Actual Round: ", roundsCounter)
			roundsCounterFlag = 1
		}
		if tc.Sensors[sensorList[4]].State == false { // release rounds counter
			roundsCounterFlag = 0
		}

		// ================================================================================================================ Brake and Reset
		if roundsCounter >= rounds {
			SetBrake(tc, 2)

			if actualSpeed == 0 {
				time.Sleep(time.Duration(autoSleepTime) * time.Millisecond)
				resetAutoFlags(tc)
			}
		}

	}
}

// ControlRoundRobin is run in a short interval
func ControlRoundRobin(tc *traincontrol.TrainControl, train *traincontrol.Train) {

	//RoundRobin Initialization ===========================================================================================================

	//Decide between track 1 or 2 (Display Location)
	// if initialRoundRobin == 0 {
	// 	initialRoundRobin = 1
	// 	if tc.Blocks["a"].Train != nil { // check track 1 for current train
	// 		trackOffset = 0
	// 		currentTrack = trackOffset
	// 	} else if tc.Blocks["b"].Train != nil { // else check track 2 for current train
	// 		trackOffset = 1
	// 		currentTrack = trackOffset
	// 	} else {
	// 		doRoundRobin = 0
	// 		initialRoundRobin = 0
	// 		resetRoundRobin(tc)
	// 		tc.PublishMessage(struct {
	// 			DoRoundRobin bool `json:"doroundrobin"`
	// 		}{
	// 			DoRoundRobin: false,
	// 		}) //synchronize all websites with set state
	// 		return
	// 	}
	// }
	trackOffset = 0
	currentTrack = trackOffset
	// Define actual blocks and get Trains Definition
	if roundRobinRunning == 0 {
		//Define actual blocks
		nextTrack = roundRobinTrackSelection(tc, trackOffset, currentTrack)

		SetDirection(tc, "f")                           //drive always forward direction because of switch sensor distance
		SetTrack(tc, roundRobinTracks[trackOffset]+"w") // exit first the westbound track
		SetTrack(tc, roundRobinTracks[nextTrack]+"o")   // enter first the eastbound track
		actualDirection = targetDirection               // sctual direction hard coded (no direction change possible for RoundRobin)

		//Select train definition from both tracks
		//currentTrain = tc.Blocks[roundRobinTracks[trackOffset]].Train
		//nextTrain = tc.Blocks[roundRobinTracks[nextTrack]].Train

		// start current Train
		for i, _ := range roundRobinTargetSpeeds {
			if i != 0 { // exclude index 0 since this is block of next train
				roundRobinTargetSpeeds[i] = train.MaxSpeed // override trains definition since two trains active
			}
		}
		targetSpeed = 0
		actualSpeed = 0
		roundRobinRunning = 1
	}

	//RoundRobin is running ===============================================================================================================
	// ====================================================================================================================================
	if roundRobinRunning == 1 {
		if targetBlocks != actualBlocks {
			actualBlocks = targetBlocks
			setSwitches(tc, targetBlocks)
			setSensorList(tc, targetBlocks, actualDirection)
			setBlocksDirection(tc, targetBlocks, actualDirection)
			//setBlocksSpeed(tc, train, targetBlocks, actualSpeed)
			resetInactiveBlocks(tc, targetBlocks)
			TimeReset(tc)
		}

		if roundRobinActualSpeeds != roundRobinTargetSpeeds {
			for i, _ := range roundRobinTargetSpeeds {
				adjustSpeedPerBlock(tc, train, i)
			}
		}

		if tc.Sensors[sensorList[3]].State == false { // =================================================== Station Exit (Round start) 15%
			if progressCurrentTrain <= 15 { // Current Train on Display Track (will start first)
				updateCurrentTrain(tc, 15)
			}
			if progressCurrentTrain >= 50 && progressNextTrain < 15 {
				updateNextTrain(tc, 15) // Next train will follow when current train reached open track (sensor[5])
			}
		}

		if tc.Sensors[sensorList[4]].State == false { // =================================================== Station Block Exit 20%
			if progressCurrentTrain <= 20 {
				updateCurrentTrain(tc, 20)
			}
			if progressCurrentTrain >= 50 && progressNextTrain < 20 && progressNextTrain >= 15 {
				updateNextTrain(tc, 20)
			}
		}

		if tc.Sensors[sensorList[5]].State == false { // =================================================== Tunnel Exit 35%
			if progressCurrentTrain <= 35 {
				updateCurrentTrain(tc, 35)
			}
			if progressCurrentTrain >= 50 && progressNextTrain < 35 && progressNextTrain >= 20 {
				updateNextTrain(tc, 35)
			}
		}

		if tc.Sensors[sensorList[6]].State == false { // =================================================== open Track 50%
			if progressCurrentTrain <= 50 {
				updateCurrentTrain(tc, 50)
				// switch west outbound to next Train
				SetTrack(tc, roundRobinTracks[nextTrack]+"w") // exit first the westbound track
				setSwitches(tc, targetBlocks)
				// start next Train. position [1] represents next Track
				roundRobinTargetSpeeds[0] = train.MaxSpeed
			}
			if progressCurrentTrain >= 80 && progressNextTrain < 50 && progressNextTrain >= 35 {
				updateNextTrain(tc, 50)
				// brake to crawlspeed in case progressCurrentTrain not at least at 85%
				if progressCurrentTrain < 85 {
					// start next Train. position [3] represents open Track
					roundRobinTargetSpeeds[3] = train.CrawlSpeed
				}
			}
		}

		if tc.Sensors[sensorList[7]].State == false { // =================================================== Tunnel Entry 65%
			if progressCurrentTrain <= 65 {
				updateCurrentTrain(tc, 65)
			}
			if progressCurrentTrain >= 80 && progressNextTrain < 65 && progressNextTrain >= 50 {
				updateNextTrain(tc, 65)
				// if at least progressCurrentTrain reached 85%: switch inbound east to current train and target speed max
				// else: brake and wait until current train is at 85% before switching and ramping up again
				if progressCurrentTrain >= 100 {
					SetTrack(tc, roundRobinTracks[trackOffset]+"o")
					setSwitches(tc, targetBlocks)
					roundRobinTargetSpeeds[3] = train.MaxSpeed
					roundRobinTargetSpeeds[2] = train.MaxSpeed
				} else {
					roundRobinActualSpeeds[3] = 0 // open track // used actual speed instead of target to get rid of speed ramp
					roundRobinActualSpeeds[2] = 0 // switches
				}
			}
		}

		if tc.Sensors[sensorList[8]].State == false { // =================================================== Station Block Entry 80%
			if progressCurrentTrain <= 80 {
				updateCurrentTrain(tc, 80)
				roundRobinTargetSpeeds[1] = train.CrawlSpeed // corresponds with target track
			}
			if progressCurrentTrain >= 100 && progressNextTrain < 80 && progressNextTrain >= 65 {
				updateNextTrain(tc, 80)
				roundRobinTargetSpeeds[1] = train.CrawlSpeed // corresponds with target track
			}
		}

		if tc.Sensors[sensorList[9]].State == false { // =================================================== Station Entry 85%
			if progressCurrentTrain <= 85 {
				updateCurrentTrain(tc, 85)
			}
			if progressCurrentTrain >= 100 && progressNextTrain < 85 && progressNextTrain >= 80 {
				updateNextTrain(tc, 85)
			}
		}

		if tc.Sensors[sensorList[10]].State == false { // ================================================== Tunnel Exit (End of Round) 100%
			if progressCurrentTrain >= 85 {
				updateCurrentTrain(tc, 100)
				roundRobinActualSpeeds[1] = 0 // corresponds with target track
			}
			if progressCurrentTrain >= 100 && progressNextTrain < 100 && progressNextTrain >= 85 {
				updateNextTrain(tc, 100)
				roundRobinActualSpeeds[1] = 0 // corresponds with target track
			}
		}

		//write back trains definition to new tracks ======================================================= End Condition
		if progressCurrentTrain >= 100 && progressNextTrain >= 100 {
			waitCounter++
			// 20 waitCounter represents 1 second
			if waitCounter >= 100 { // wait for some seconds (waitcounter/10) before reset ================== EDIT THIS LINE TO IN/DECREASE THE WAITTIME AFTER ONE ROUND
				//tc.Blocks[roundRobinTracks[trackOffset]].Train = nextTrain
				//tc.Blocks[roundRobinTracks[nextTrack]].Train = currentTrain
				resetRoundRobin(tc)
			}
		}
	}
}

// roundRobinTrackSelection is repsonsible to set next track for RoundRobin
func roundRobinTrackSelection(tc *traincontrol.TrainControl, trackOffset int, currentTrack int) int {
	if currentTrack < 3 {
		nextTrack = currentTrack + 1 // increment tracks between offset (display train) and maximum number of tracks
	} else {
		nextTrack = trackOffset + 1 // start always on track above display train
	}
	return nextTrack
}

// resetRoundRobin set all RoundRobin related/important variables back to default. Default is used to initialize and/or start new round
func resetRoundRobin(tc *traincontrol.TrainControl) {
	roundRobinRunning = 0
	waitCounter = 0
	currentTrack = nextTrack
	roundRobinTargetSpeeds = [4]int{0, 0, 0, 0}
	progressCurrentTrain = -1 // set to -1 to for proper update
	progressNextTrain = -1
	updateCurrentTrain(tc, 0)
	updateNextTrain(tc, 0)
}

// PrintAll is just a function to print status of all values
func PrintAll(tc *traincontrol.TrainControl) {
	log.Println("----------------PRINT ALL------------------- ", actualBlocks)
	log.Println("----------------actualBlocks: ", actualBlocks)
	log.Println("----------------targetBlocks: ", targetBlocks)
	log.Println("----------------sensorStates: ")
	for _, id := range tc.Sensors {
		{
			log.Println("----------------Sensor: ", id, tc.Sensors[id.ID].State)
		}
	}
	log.Println("----------------sensorList: ", sensorList)
	log.Println("----------------distanceList: ", distanceList)
	log.Println("----------------actualDirection: ", actualDirection)
	log.Println("----------------targetDirection: ", targetDirection)
	log.Println("----------------actualSpeed: ", actualSpeed)
	log.Println("----------------targetSpeed: ", targetSpeed)
	log.Println("----------------previousSpeed: ", previousSpeed)
	log.Println("----------------lastAccelerateTick: ", lastAccelerateTick)
	log.Println("----------------isDriveable: ", isDriveable())
	log.Println("----------------autoBrake: ", autoBrake)
	log.Println("----------------autoBrakeRelease: ", autoBrakeReleased)
	log.Println("----------------autoBrakeAbsolute: ", autoBrakeAbsolute)
	log.Println("----------------setSpeedFlag: ", setSpeedFlag)

	log.Println("--------------------------------------------- ")
	log.Println("--------------------------------------------- ")
	log.Println("----------------ROUND ROBIN ----------------- ")
	log.Println("--------------------------------------------- ")
	log.Println("--------------------------------------------- ")
	log.Println("----------------doRoundRobin: ", doRoundRobin)
	log.Println("----------------initialRoundRobin: ", initialRoundRobin)
	log.Println("----------------trackOffset: ", trackOffset)
	log.Println("----------------roundRobinRunning: ", roundRobinRunning)
	log.Println("----------------progressCurrentTrain: ", progressCurrentTrain)
	log.Println("----------------seprogressNextTrain: ", progressNextTrain)
	log.Println("----------------waitCounter: ", waitCounter)
	log.Println("----------------nextTrack: ", nextTrack)
	log.Println("----------------currentTrain: ", currentTrain)
	log.Println("----------------nextTrain: ", nextTrain)
	log.Println("----------------roundRobinTrackSelection: ", roundRobinTracks)
	log.Println("----------------roundRobinTargetSpeeds: ", roundRobinTargetSpeeds)
	log.Println("----------------roundRobinActualSpeeds: ", roundRobinActualSpeeds)
}

// adjustSpeed for Speed regulation during normal (all) modes (gets blocks with 4 strings)
func adjustSpeed(tc *traincontrol.TrainControl, train *traincontrol.Train, actualBlocks [4]string, targetSpeed int) {
	now := time.Now()
	tickDuration := train.Accelerate.Time * time.Millisecond
	if now.Sub(lastAccelerateTick) > tickDuration {
		lastAccelerateTick = now
		speedDiff := actualSpeed - targetSpeed
		inc := 0
		if speedDiff > 0 {
			// deccelerate
			inc = -train.Brake.Step
		}
		if speedDiff < 0 {
			//accelerate
			inc = train.Accelerate.Step
		}
		actualSpeed += inc
		if manual == 1 {
			actualSpeed = targetSpeed
		}

		setBlocksSpeed(tc, train, actualBlocks, actualSpeed)
		tc.PublishMessage(struct {
			ActualSpeed int `json:"actualspeed"`
		}{
			ActualSpeed: actualSpeed,
		}) //synchronize all websites with set state
	}
}

// adjustSpeedPerBlock for Speed regulation during RoundRobin only (gets only one block)
func adjustSpeedPerBlock(tc *traincontrol.TrainControl, train *traincontrol.Train, blockPosition int) {
	now := time.Now()
	tickDuration := train.Accelerate.Time * time.Millisecond
	if now.Sub(lastAccelerateTicks[blockPosition]) > tickDuration {
		lastAccelerateTicks[blockPosition] = now
		speedDiff := roundRobinActualSpeeds[blockPosition] - roundRobinTargetSpeeds[blockPosition]
		inc := 0
		if speedDiff > 0 {
			// deccelerate
			inc = -train.Brake.Step
		}
		if speedDiff < 0 {
			//accelerate
			inc = train.Accelerate.Step
		}
		roundRobinActualSpeeds[blockPosition] += inc

		setBlockSpeed(tc, train, actualBlocks[blockPosition], roundRobinActualSpeeds[blockPosition])
	}
}

// SetBrake set flag to brake
func SetBrake(tc *traincontrol.TrainControl, s int) {
	autoBrake = s
	// 1: used for manual mode (overrides and resets auto mode)
	// 2: used to brake in auto mode
	// 3: used to brake in absolute mode (block not clear)

	if s == 1 {
		auto = 0 // reset Automatic mode in case of manual switch OFF
		resetAutoFlags(tc)
	}

	if manual == 1 { // in manual mode ignos brake procedure, just stop
		autoBrake = 0
		targetSpeed = 0
	}
}

// SetDirection sets the direction
func SetDirection(tc *traincontrol.TrainControl, d string) {
	if targetDirection != d && autoBrakeAbsolute == 1 {
		targetSpeed = 0
		autoBrake = 0
		autoBrakeAbsolute = 0
	}
	targetDirection = d

	tc.PublishMessage(struct {
		Direction string `json:"direction"`
	}{
		Direction: targetDirection,
	})
	log.Println("----------------Direction set: ", targetDirection)
}

// SetActualSpeed sets the speed
func SetActualSpeed(tc *traincontrol.TrainControl, s int) {
	if manual == 1 {
		SetSpeed(tc, s)
	}
}

// SetSpeed sets the speed
func SetSpeed(tc *traincontrol.TrainControl, s int) {
	if s > 99 {
		s = 99
	}
	if autoBrakeAbsolute == 1 { // set previous speed to s indstead of targetspeed because brake active because block not clear
		previousSpeed = s // will be set to target speed after brake process done
	} else {
		targetSpeed = s
	}

	log.Println("----------------Speed set: ", s)

	tc.PublishMessage(struct {
		Speed int `json:"speed"`
	}{
		Speed: s,
	}) //synchronize all websites with set state
}

// SetTrack sets the track
func SetTrack(tc *traincontrol.TrainControl, track string) {
	var switchLocation = string(getSwitchLocation(track))
	var block = string(getBlock(track))
	switch switchLocation {
	case "o":
		if doCircle == 1 { // for doCircle set east and westbound to same branch. no independent tracks allowed
			targetBlocks[0] = block + "o"
			targetBlocks[1] = block + "w"
		} else {
			targetBlocks[0] = block + switchLocation
		}
	case "w":
		if doCircle == 1 {
			targetBlocks[0] = block + "o"
			targetBlocks[1] = block + "w"
		} else {
			targetBlocks[1] = block + switchLocation
		}
	default:
		targetBlocks[0] = block + "o"
		targetBlocks[1] = block + "w"
	}

	log.Println("----------------setTrack: Blocks set: ", targetBlocks)

	if track[0] == 'a' { // set trackValue.. used for setRandomTrack and setOrderTrack..  set here to avoid use same track twice after initialisation
		trackValue = 0.0
	} else if track[0] == 'b' {
		trackValue = 0.25
	} else if track[0] == 'c' {
		trackValue = 0.5
	} else if track[0] == 'd' {
		trackValue = 0.75
	}
	//randomTrackFlag = 1 // disable random Track for one iteration

	SendMapVisuals(tc, targetBlocks, targetDirection)
}

// isDriveable checks wheather a train can drive
func isDriveable() bool {
	if doRoundRobin == 1 {
		return true
	}
	if targetDirection == "s" {
		return false
	}
	for _, block := range targetBlocks {
		if block == "" {
			return false
		}
	}
	for _, block := range targetBlocks {
		if block == "+" || block == "-" {
			return false
		}
	}
	return true
}

// allTrainsStopped will stop actual train at stop location befor program continues.
func allTrainsStopped() bool {
	if actualSpeed > 0 {
		return false
	}
	return true
}

// canSwitch will return if allowed to send switches depending from train position (sensor released)
func canSwitch(tc *traincontrol.TrainControl, targetBlocks [4]string, actualBlocks [4]string) bool {
	//f = westbound = default
	targetInbound := getBlock(targetBlocks[0])
	targetOutbound := getBlock(targetBlocks[1])
	actualInbound := getBlock(actualBlocks[0])
	actualOutbound := getBlock(actualBlocks[1])
	if actualDirection == "b" { //b= eastbound
		targetInbound = getBlock(targetBlocks[1])
		targetOutbound = getBlock(targetBlocks[0])
		actualInbound = getBlock(actualBlocks[1])
		actualOutbound = getBlock(actualBlocks[0])
	}

	//[7] Tunnel Entry, [5] Tunnel Exit
	//allow setSwitches only at specific locations of track (most far away) to operate secure
	if (targetInbound != actualInbound) && tc.Sensors[sensorList[5]].State == false {
		return true
	}

	if (targetOutbound != actualOutbound) && tc.Sensors[sensorList[7]].State == false {
		return true
	}
	return false
}

//=====================================================================================================================================================================
//========================================================================== M E N U ==================================================================================
//=====================================================================================================================================================================

// SwitchDoCircle Menufunction for named funtion (selection of described mode)
func SwitchDoCircle(tc *traincontrol.TrainControl, b int) {
	if b == 0 {
		doCircle = 0
	} else {
		doCircle = 1
		if isDriveable() {
			SetTrack(tc, actualBlocks[0]) // synchronize track and switches with first actualBlock for doCircle Mode
		}
	}
}

// SetSimpleMode Menufunction for named funtion (selection of described mode)
func SetSimpleMode(tc *traincontrol.TrainControl, b int) {
	if b == 0 {
		manual = 0
	} else {
		manual = 1
		auto = 0 // make sure automatic is off
	}
}

// SetAuto Menufunction for named funtion (selection of described mode)
func SetAuto(tc *traincontrol.TrainControl, b int) {
	if b == 1 {
		// switch auto on, but for visuals only (no effect for Round Robin)
		if isDriveable() && !allTrainsStopped() && doRoundRobin == 0 {
			SetBrake(tc, 1)
		}
		auto = 1
		manual = 0 // make sure manual mode is off
		if doCircle == 0 && doRoundRobin == 0 {
			SwitchDoCircle(tc, 1)
		}
	} else {
		if isDriveable() && !allTrainsStopped() && doRoundRobin == 0 {
			SetBrake(tc, 1)
		}
		auto = 0
		resetAutoFlags(tc)
	}
}

// SwitchRandomDirection Menufunction for named funtion (selection of described mode)
func SwitchRandomDirection(tc *traincontrol.TrainControl, b int) {
	if b == 1 {
		randomDirection = 1
	} else {
		randomDirection = 0
	}
}

// SwitchRandomTrack Menufunction for named funtion (selection of described mode)
func SwitchRandomTrack(tc *traincontrol.TrainControl, b int) {
	if b == 1 {
		randomTrack = 1 // overrides ordered track selection
	} else {
		randomTrack = 0 // if automatic active always ordered track is choosen
	}
}

// SwitchRandomRounds Menufunction for named funtion (selection of described mode)
func SwitchRandomRounds(tc *traincontrol.TrainControl, b int) {
	if b == 1 {
		randomRounds = 1
	} else {
		randomRounds = 0
	}
}

// SetMaxRoundsInt Menufunction for named funtion (selection of described mode)
func SetMaxRoundsInt(tc *traincontrol.TrainControl, i int) {
	maxRounds = i
}

// SetMinRoundsInt Menufunction for named funtion (selection of described mode)
func SetMinRoundsInt(tc *traincontrol.TrainControl, i int) {
	minRounds = i
}

// SwitchRoundRobin Menufunction for named funtion (selection of described mode)
func SwitchRoundRobin(tc *traincontrol.TrainControl, b int) {
	if b == 1 {
		doRoundRobin = 1
		SetAuto(tc, 1)
		tc.PublishMessage(struct {
			DoRoundRobin bool `json:"doroundrobin"`
		}{
			DoRoundRobin: true,
		}) //synchronize all websites with set state
	} else {
		doRoundRobin = 0
		initialRoundRobin = 0
		resetRoundRobin(tc)
		SetAuto(tc, 0)
		tc.PublishMessage(struct {
			DoRoundRobin bool `json:"doroundrobin"`
		}{
			DoRoundRobin: false,
		}) //synchronize all websites with set state
	}
}

// SwitchLightTunnel will turn on the service light in tunnel and back of station wall (lower level)
func SwitchLightTunnel(tc *traincontrol.TrainControl, b int) {
	if b > 0 {
		tc.SetBlockDirection("j", "b")
		tc.SetBlockSpeed("j", b)
	} else {
		tc.SetBlockDirection("j", "s")
		tc.SetBlockSpeed("j", 0)
	}
}

//=====================================================================================================================================================================
//======================================================================== B L O C K S ================================================================================
//=====================================================================================================================================================================

const allBlocks = "abcdfg"

// setBlocksDirection sets the direction for all blocks
func setBlocksDirection(tc *traincontrol.TrainControl, blocks [4]string, direction string) {
	for _, block := range blocks {
		direction2Arduino(tc, getBlock(block), direction)
	}
	tc.PublishMessage(struct {
		Direction string `json:"direction"`
	}{
		Direction: direction,
	}) //synchronize all websites with set state
	SendMapVisuals(tc, blocks, direction)
}

// setBlocksSpeed sets the speed for all blocks
func setBlocksSpeed(tc *traincontrol.TrainControl, train *traincontrol.Train, blocks [4]string, speed int) {
	for _, block := range blocks {
		// recalculate real speed based on MaxSpeed per train defined in Config Json
		speed2Arduino(tc, getBlock(block), int(math.Round(float64(speed)*(float64(train.MaxSpeed)/100))))
	}
}

// setBlockSpeed sets the speed for one block (used in RoundRobin)
func setBlockSpeed(tc *traincontrol.TrainControl, train *traincontrol.Train, block string, speed int) {
	// recalculate real speed based on MaxSpeed per train defined in Config Json
	speed2Arduino(tc, getBlock(block), int(math.Round(float64(speed)*(float64(train.MaxSpeed)/100))))
}

// setSwitches sets all switches for inbound and outbound direction
func setSwitches(tc *traincontrol.TrainControl, blocks [4]string) {
	switches2Arduino(tc, blocks[0])
	switches2Arduino(tc, blocks[1])
}

// getBlock return block letter for direction and speed (a,b,c,d,f,g and so on)
func getBlock(block string) byte {
	if len(block) > 0 {
		return block[0]
	}
	return '+'
}

// getSwitchLocation return "o" or "w" to decide which switches should operate
func getSwitchLocation(block string) byte {
	if len(block) > 1 {
		return block[1]
	}
	return '-'
}

// getInactiveBlocks returns a string of inactive blocks, based on allBlocks defined in the begin of  B L O C K S  section
func getInactiveBlocks(blocks [4]string) string {
	var currentBlocks = allBlocks
	for _, block := range blocks {
		currentBlocks = strings.Replace(currentBlocks, string(getBlock(block)), "", 1)
	}
	return currentBlocks
}

// resetInactiveBlocks set all inactive blocks to speed = 0 and direction 's'
func resetInactiveBlocks(tc *traincontrol.TrainControl, blocks [4]string) {
	var inactiveBlocks = getInactiveBlocks(blocks)
	for _, block := range inactiveBlocks {
		partialResetBlock2Arduino(tc, byte(block))
	}
}

//blockClear returns if train can safely leave the station
func blockClear(blocks [4]string) bool {
	eastbound := getBlock(blocks[0])
	westbound := getBlock(blocks[1])

	if eastbound == westbound {
		return true
	}
	return false
}

//=====================================================================================================================================================================
//===================================================================== A U T O M A T I C =============================================================================
//=====================================================================================================================================================================

//setRandomRounds sets a random amount of rounds between minRounds and maxRounds
func setRandomRounds(tc *traincontrol.TrainControl, minRounds int, maxRounds int) float64 {
	diff := absolute(maxRounds, minRounds)
	return rand.Float64()*float64(diff) + float64(minRounds)
}

// setRandomDirection sets a random Direction
func setRandomDirection(tc *traincontrol.TrainControl) {
	r := rand.Float64()
	log.Println("----------------DirectionValue now: ", r)

	if r > 0.5 {
		SetDirection(tc, "f")
		log.Println("----------------Forward because DirectionValue was: ", r)
	} else {
		SetDirection(tc, "b")
		log.Println("----------------Backward because DirectionValue was: ", r)
	}
}

// setOrderTrack sets a ordered Track
func setOrderTrack(tc *traincontrol.TrainControl, value float64) {
	if value == 0.0 {
		SetTrack(tc, "bo")
	} else if value == 0.25 {
		SetTrack(tc, "co")
	} else if value == 0.5 {
		SetTrack(tc, "do")
	} else {
		SetTrack(tc, "ao")
	}
	log.Println("----------------setTrack in order: Blocks set: ", targetBlocks)
}

// setRandomTrack sets a random Track
func setRandomTrack(tc *traincontrol.TrainControl) {
	r := rand.Float64()
	r = round(r, 0.25)
	log.Println("----------------trackValue now: ", r)
	for r == trackValue || r == 1 { // recalculate to make sure to not drive same track twice in a row
		r = rand.Float64()
		r = round(r, 0.25)
		log.Println("----------------Recalculation. trackValue now: ", r)
	}
	trackValue = r // store value to exclude old track from new selection (next round)

	if r == 0 {
		SetTrack(tc, "ao")
	} else if r == 0.25 {
		SetTrack(tc, "bo")
	} else if r == 0.5 {
		SetTrack(tc, "co")
	} else {
		SetTrack(tc, "do")
	}

	log.Println("----------------setTrack randomly: Blocks set: ", targetBlocks)
}

// velocity is the main function to measure alls velocitys for the sensors
func velocity(tc *traincontrol.TrainControl) {
	for _, id := range sensorList { //for each id in sensorlist
		if tc.Sensors[id] == nil {
			continue
		}
		if tc.Sensors[id].State == false && SensorTimes[id-1] == ini {
			distance := getPreviousDistance(tc, id)
			end := time.Now()
			start := getPreviousTime(tc, id)
			SensorTimes[id-1] = end
			speed := getVelocity(tc, start, end, distance)
			lastID := getPreviousSensor(tc, id)

			if speed > 0 && speed < 999 && //Smoothing Attempt 1: large Outlier Detection
				id != sensorList[2] && // Smoothing Attempt 2: ignore measurments from short sections as inaccurate
				id != sensorList[4] &&
				id != sensorList[9] &&
				id != sensorList[11] { // Publish Speed (dirty and cheap attempt)
				speed = averageSpeedExcludeOutliers(tc, speed) // Smoothing Attempt 3: floating average of last value (Exclude largest Outlier in case of false measurement)
				log.Println("----------------Velocity between Sensor ", id, " and sensor ", lastID, ": ", speed, " km/h")
				tc.PublishMessage(struct {
					Velocity int `json:"velocity"`
				}{
					Velocity: speed,
				})
			}

			if actualSpeed == 0 { // Publish 0 Speed
				log.Println("----------------Velocity is now: ", 0, " km/h")
				tc.PublishMessage(struct {
					Velocity int `json:"velocity"`
				}{
					Velocity: 0,
				})
			}
		}

		if tc.Sensors[sensorList[len(sensorList)-1]].State == true && timeResetFlag == 1 { // enable reset
			timeResetFlag = 0
		}

		if tc.Sensors[sensorList[len(sensorList)-1]].State == false && timeResetFlag == 0 { // reset timelist and enable all sensors for next measurement
			TimeReset(tc)
			timeResetFlag = 1
		}
	}
}

//getVelocity measures velocity between sensprs
func getVelocity(tc *traincontrol.TrainControl, start time.Time, end time.Time, distance float64) int {
	duration := (end.Sub(start)).Seconds()
	speed := (distance / duration) * 3.6 * 160 // calculate float velocity in n scale (1:160) in km/h
	//speedRoundedUp := int(math.Ceil(speed/10)) * 10 // round up to next ten
	speedRoundedDown := int(speed/5) * 5 // round down to next fifth
	return int(speedRoundedDown)
}

// averageSpeed calculates mean of last values (depends on length of average list)
func averageSpeed(tc *traincontrol.TrainControl, s int) int {
	avg := 0
	for i := 1; i < len(speedAverageList); i++ {
		avg = avg + speedAverageList[i]
		speedAverageList[i-1] = speedAverageList[i]
	}
	speedAverageList[len(speedAverageList)-1] = s
	avg = (avg + s) / len(speedAverageList)
	return int(avg)
}

// averageSpeedExcludeOutliers same as Aver Speed but exlude largest Outlier
func averageSpeedExcludeOutliers(tc *traincontrol.TrainControl, s int) int {
	avg := 0
	largest := s
	for i := 1; i < len(speedAverageList); i++ {
		avg = avg + speedAverageList[i]
		if largest < speedAverageList[i] {
			largest = speedAverageList[i]
		}
		speedAverageList[i-1] = speedAverageList[i]
	}
	speedAverageList[len(speedAverageList)-1] = s

	avg = (avg + s - largest) / (len(speedAverageList) - 1)
	return int(avg)
}

// resetAutoFlags reset all Flags for Automatic Mode.
func resetAutoFlags(tc *traincontrol.TrainControl) {
	randomDirectionFlag = 0
	randomTrackFlag = 0
	randomRoundsFlag = 0
	roundsCounterFlag = 0
	setSpeedFlag = 0
	roundsCounter = 0
	log.Println("----------------Reset Flags for Automatic Mode ")
}

//=====================================================================================================================================================================
//======================================================================= S E N S O R S ===============================================================================
//=====================================================================================================================================================================

// getSensors receive list of Sensors, revere order if needed and cut it to length
func getSensors(tc *traincontrol.TrainControl, block string, direction string) []int {
	sensors := tc.Blocks[string(getBlock(block))].Sensors
	if direction == "b" {
		sensors = sensors[1:]
	} else {
		sensors = sensors[:len(sensors)-1]
	}
	return sensors
}

// getDistances receive list of Sensors, revere order if needed and cut it to length
func getDistances(tc *traincontrol.TrainControl, block string, direction string) []float64 {
	distances := tc.Blocks[string(getBlock(block))].Distances
	return distances
}

// setSensorList creates depending on actual blocks the sensor list with defined order referring distances
func setSensorList(tc *traincontrol.TrainControl, blocks [4]string, direction string) {

	// start sensorList for first letter of defined block. add all sensors but skip last sensor to sensorList (done in getSensors)
	// start distanceList for first letter of defined block. add all distances in correct order
	if direction == "b" { // Backward Direction (0,3,1 -> east, middle, west))
		sensorPart := getSensors(tc, blocks[1], direction)
		sensorPart = append(sensorPart, getSensors(tc, blocks[3], direction)...) // append defined block to list (after reverse 0,3,1 -> east, middle, west)
		sensorPart = append(sensorPart, getSensors(tc, blocks[0], direction)...)
		sensorList = sensorPart
		reverseAny(sensorList)

		distancePart := getDistances(tc, blocks[1], direction)
		distancePart = append(distancePart, getDistances(tc, blocks[3], direction)...)
		distancePart = append(distancePart, getDistances(tc, blocks[0], direction)...)
		distanceList = distancePart
		reverseAny(distanceList)
	} else { // Forward Direction (1,3,0 -> west, middle, east)
		sensorPart := getSensors(tc, blocks[1], direction)
		sensorPart = append(sensorPart, getSensors(tc, blocks[3], direction)...) // append defined block to list (1,3,0 -> west, middle, east)
		sensorPart = append(sensorPart, getSensors(tc, blocks[0], direction)...)
		sensorList = sensorPart

		distancePart := getDistances(tc, blocks[1], direction)
		distancePart = append(distancePart, getDistances(tc, blocks[3], direction)...)
		distancePart = append(distancePart, getDistances(tc, blocks[0], direction)...)
		distanceList = distancePart
	}

	log.Println("----------------sensorList: ", sensorList)
	log.Println("----------------distanceList: ", distanceList)
}

// getPreviousSensor provides distance to last sensor
func getPreviousSensor(tc *traincontrol.TrainControl, id int) int {
	for i := 0; i < len(sensorList); i++ {
		if sensorList[i] == id {
			if i == 0 {
				return sensorList[sensorPerRound-1]
			}
			return sensorList[i-1]

		}
	}
	return 0
}

// getPreviousDistance provides distance to last sensor
func getPreviousDistance(tc *traincontrol.TrainControl, id int) float64 {
	for i := 0; i < len(distanceList); i++ {
		if sensorList[i] == id {
			if i == 0 {
				return distanceList[sensorPerRound-1]
			}
			return distanceList[i-1]
		}
	}
	return 0
}

// getPreviousTime provides time of last sensor activation
func getPreviousTime(tc *traincontrol.TrainControl, id int) time.Time {
	lastID := -1
	for i := 0; i < len(sensorList); i++ {
		if sensorList[i] == id {
			if i == 0 {
				lastID = sensorList[len(sensorList)-1]
			} else {
				lastID = sensorList[i-1]
			}
		}
	}

	return SensorTimes[lastID-1]
}

// TimeReset reset SensorTime to init values
func TimeReset(tc *traincontrol.TrainControl) {
	for i := 0; i < len(SensorTimes); i++ {
		SensorTimes[i] = ini
		log.Println("----------------TIME RESET at Sensor: ", sensorList[len(sensorList)-1])
	}
}

//=====================================================================================================================================================================
//======================================================================= A R D U I N O ===============================================================================
//=====================================================================================================================================================================

// resetArduino resets Arduino
func resetArduino(tc *traincontrol.TrainControl) {
	log.Println("----------------RESET ARDUINO BOTTOM LAYER (MAIN LINE)")
	//rstz
	tc.ResetArduino()
}

// direction2Arduino sets direction for block
func direction2Arduino(tc *traincontrol.TrainControl, block byte, direction string) {
	tc.SetBlockDirection(string(block), direction)
}

// speed2Arduino sets the speed of an arduino with a byte
func speed2Arduino(tc *traincontrol.TrainControl, block byte, speed int) {
	tc.SetBlockSpeed(string(block), speed)
}

// partialResetBlock2Arduino resets a block
func partialResetBlock2Arduino(tc *traincontrol.TrainControl, block byte) {
	speed2Arduino(tc, block, 0)
	direction2Arduino(tc, block, "s")
	time.Sleep(100 * time.Millisecond)
	log.Println("----------------reset block", string(block))
}

// switches2Arduino alters junctions
func switches2Arduino(tc *traincontrol.TrainControl, block string) {
	if block == "aw" {
		log.Println("----------------Send Switches for Track 1 west in-/outbound to Arduino")
		tc.SetSwitch("e", "0")
		tc.SetSwitch("f", "0")
	}
	if block == "ao" {
		log.Println("----------------Send Switches for Track 1 to east in-/outbound Arduino")
		tc.SetSwitch("a", "0")
		tc.SetSwitch("b", "0")
	}
	if block == "bw" {
		log.Println("----------------Send Switches for Track 2 west in-/outbound to Arduino")
		tc.SetSwitch("d", "0")
		tc.SetSwitch("e", "1")
		tc.SetSwitch("f", "0")
	}
	if block == "bo" {
		log.Println("----------------Send Switches for Track 2 east in-/outbound to Arduino")
		tc.SetSwitch("a", "0")
		tc.SetSwitch("b", "1")
	}
	if block == "cw" {
		log.Println("----------------Send Switches for Track 3 west in-/outbound to Arduino")
		tc.SetSwitch("d", "1")
		tc.SetSwitch("e", "1")
		tc.SetSwitch("f", "0")
	}
	if block == "co" {
		log.Println("----------------Send Switches for Track 3 east in-/outbound to Arduino")
		tc.SetSwitch("a", "1")
		tc.SetSwitch("c", "0")
	}
	if block == "dw" {
		log.Println("----------------Send Switches for Track 4 west in-/outbound to Arduino")
		tc.SetSwitch("f", "1")
	}
	if block == "do" {
		log.Println("----------------Send Switches for Track 4 east in-/outbound to Arduino")
		tc.SetSwitch("a", "1")
		tc.SetSwitch("c", "1")
	}
}

// EmergencyStop2Arduino stops all tracks
func EmergencyStop2Arduino(tc *traincontrol.TrainControl) {
	log.Println("----------------Emergency Stop released")
	var history = 1
	if manual == 0 {
		history = 0
		manual = 1
	}

	//SetActualSpeed(tc, 0)
	resetArduino(tc)
	actualDirection = "s"
	targetDirection = "s"
	actualSpeed = 0
	targetSpeed = 0
	previousSpeed = 0
	actualBlocks = EmptyBlock
	targetBlocks = EmptyBlock
	sensorList = EmptySensors
	distanceList = EmptyDistances

	log.Println("----------------Velocity is now: ", 0, " km/h")
	tc.PublishMessage(struct {
		Velocity int `json:"velocity"`
	}{
		Velocity: 0,
	})

	if history == 0 { // set manual to 0 only if it was 0 before
		manual = 0
	}

	tc.PublishMessage(struct {
		Speed int `json:"speed"`
	}{
		Speed: 0,
	}) //synchronize all websites with set state

	tc.PublishMessage(struct {
		ActualSpeed int `json:"actualspeed"`
	}{
		ActualSpeed: 0,
	}) //synchronize all websites with set state
}

//=====================================================================================================================================================================
//========================================================================== M I S C ==================================================================================
//=====================================================================================================================================================================

func reverse(s []interface{}) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func reverseAny(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}

}

func round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

// absolute returns the absolute value of x-y.
func absolute(x int, y int) int {
	if x-y < 0 {
		return -(x - y)
	}
	return x - y
}

// updateCurrentTrain publish message to all listeners and update variable
func updateCurrentTrain(tc *traincontrol.TrainControl, progress int) {
	if progressCurrentTrain < progress {
		log.Println("----------------progressCurrentTrain is now: ", progress)
		tc.PublishMessage(struct {
			ProgressCurrentTrain int `json:"progressCurrentTrain"`
		}{
			ProgressCurrentTrain: progress,
		}) //synchronize all websites with set state

		sendInvalidVelocity(tc)
	}
	progressCurrentTrain = progress
}

// updateNextTrain publish message to all listeners and update variable
func updateNextTrain(tc *traincontrol.TrainControl, progress int) {
	if progressNextTrain < progress {
		log.Println("----------------progressNextTrain is now: ", progress)
		tc.PublishMessage(struct {
			ProgressNextTrain int `json:"progressNextTrain"`
		}{
			ProgressNextTrain: progress,
		}) //synchronize all websites with set state

		sendInvalidVelocity(tc)
	}
	progressNextTrain = progress
}

func sendInvalidVelocity(tc *traincontrol.TrainControl) {
	log.Println("----------------Velocity invalid since RoundRodin running")
	tc.PublishMessage(struct {
		Velocity string `json:"velocity"`
	}{
		Velocity: "",
	})
}

// SendMapVisuals sned messages to all websites to synchronize map visualization. Track, Signals...
func SendMapVisuals(tc *traincontrol.TrainControl, blocks [4]string, direction string) {
	var tracks []string
	var signals []string

	for _, block := range blocks {
		var switchLocation = string(getSwitchLocation(block))
		var track = string(getBlock(block))

		switch switchLocation {
		case "o":
			tracks = append(tracks, block, track+"m")
			if direction == "b" {
				signals = append(signals, block)
			}
		case "w":
			tracks = append(tracks, block, track+"m")
			if direction == "f" {
				signals = append(signals, block)
			}
		}
	}

	if len(tracks) > 0 {
		tracks = append(tracks, "g")
	}

	tc.PublishMessage(struct {
		whiteTrack  []string `json:"whiteTracks"`
		greenSignal []string `json:"greenSignals"`
	}{
		whiteTrack:  tracks,
		greenSignal: signals,
	}) //synchronize all websites with set state
}
