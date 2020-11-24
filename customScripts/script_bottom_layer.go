package custom

import (
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/codepuree/tilo-railway-company/pkg/traincontrol"
)

// EmptyBlock valid for specific scenario. Definition of all available blocks in that scenario down at section  B L O C K S
var EmptyBlock = [4]string{"", "", "f", "g"}

// EmptySensors valid for specific scenario. Definition of all available blocks in that scenario down at section  S E N S O R S
var EmptySensors = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// EmptyDistances valid for specific scenario. Definition of all available blocks in that scenario down at section  S E N S O R S
var EmptyDistances = []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}

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

// PrintAll is just a function to print status of all values
func PrintAll(tc *traincontrol.TrainControl) {
	log.Println("----------------PRINT ALL------------------- ", actualBlocks)
	log.Println("----------------actualBlocks: ", actualBlocks)
	log.Println("----------------targetBlocks: ", targetBlocks)
	log.Println("----------------sensorList: ", sensorList)
	log.Println("----------------distanceList: ", distanceList)
	log.Println("----------------actualDirection: ", actualDirection)
	log.Println("----------------targetDirection: ", targetDirection)
	log.Println("----------------actualSpeed: ", actualSpeed)
	log.Println("----------------targetSpeed: ", targetSpeed)
	log.Println("----------------previousSpeed: ", previousSpeed)
	log.Println("----------------lastAccelerateTick: ", lastAccelerateTick)

}

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
				// Control(tc, tc.GetActiveTrain())
				Control(tc, tc.Trains["N700"])
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
	}

	if targetSpeed != actualSpeed {
		adjustSpeed(tc, train, actualBlocks, targetSpeed)
	}

	if targetBlocks != actualBlocks {
		actualBlocks = targetBlocks
		setSwitches(tc, targetBlocks)
		setBlocksDirection(tc, targetBlocks, actualDirection)
		setBlocksSpeed(tc, targetBlocks, actualSpeed)
		setSensorList(tc, targetBlocks, actualDirection)
		resetInactiveBlocks(tc, targetBlocks)
	}
}

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
		setBlocksSpeed(tc, actualBlocks, actualSpeed)
	}
}

// SetDirection sets the direction
func SetDirection(tc *traincontrol.TrainControl, d string) {
	targetDirection = d
	tc.PublishMessage(struct {
		Direction string `json:"direction"`
	}{
		Direction: targetDirection,
	})
	log.Println("----------------Direction set: ", targetDirection)
}

// SetSpeed sets the speed
func SetSpeed(tc *traincontrol.TrainControl, s int) {
	targetSpeed = s
	//previousSpeed = actualSpeed  //need to be set only in case of sensor activation. do later in sensor block/actions
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
		targetBlocks[0] = block + switchLocation
	case "w":
		targetBlocks[1] = block + switchLocation
	default:
		targetBlocks[0] = block + "o"
		targetBlocks[1] = block + "w"
	}
	log.Println("----------------setTrack: Blocks set: ", targetBlocks)
	tc.PublishMessage(struct {
		Blocks [4]string `json:"blocks"`
	}{
		Blocks: targetBlocks,
	})
	//synchronize all websites with set state
}

// isDriveable checks wheather a train can drive
func isDriveable() bool {
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

//=====================================================================================================================================================================
//======================================================================== B L O C K S ================================================================================
//=====================================================================================================================================================================

const allBlocks = "abcdfg"

// setBlocksDirection sets the direction for all blocks
func setBlocksDirection(tc *traincontrol.TrainControl, blocks [4]string, direction string) {
	for _, block := range blocks {
		direction2Arduino(tc, getBlock(block), direction)
	}
}

// setBlocksSpeed sets the speed for all blocks
func setBlocksSpeed(tc *traincontrol.TrainControl, blocks [4]string, speed int) {
	for _, block := range blocks {
		speed2Arduino(tc, getBlock(block), speed)
	}
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

	// start sensorList for first letter of defined block. add all sensors but skip last sensor to sensorList
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
}

// getNextSensor provides information (ID and distance) of following sensor
func getNextSensor(tc *traincontrol.TrainControl) {

}

// getPreviousDistance provides distnace to last sensor
func getPreviousDistance(tc *traincontrol.TrainControl, int id) float64{
	for i := 0; i < len(distanceList); i++ {
		if sensorList[i] == id {
			if i == 0 {
				return distanceList[len(distanceList)]
			} else {
				return distanceList[i-1]
			}
		}
	}
	return 0
}

// getNextDistance provides distnace to last sensor
func getNextDistance(tc *traincontrol.TrainControl, int id) float64{
	for i := 0; i < len(distanceList); i++ {
		if sensorList[i] == id {
			return distanceList[i]
		}
	}
	return 0
}

//=====================================================================================================================================================================
//======================================================================= A R D U I N O ===============================================================================
//=====================================================================================================================================================================

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
	tc.SetBlockDirection("a", "s")
	tc.SetBlockDirection("b", "s")
	tc.SetBlockDirection("c", "s")
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

//=====================================================================================================================================================================
//========================================================================== M I S C ==================================================================================
//=====================================================================================================================================================================

func reverseInt(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

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
