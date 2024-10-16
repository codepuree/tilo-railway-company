//-------------------------------------//
//    TILO RAILWAY COMPANY             //
//    OBERE EBENE                      //
//-------------------------------------//
//    TRC_MAIN.ino                     //
//    MAIN UNIT                        //
//    PWM 122 Hz                       //
//    V.2.2020                         //
//    based on ArduinoRailwayProject   //
//-------------------------------------//

void(* resetFunc) (void) = 0;

#include <SoftwareSerial.h>
#include <Wire.h>
#include <Servo.h> 
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////// VARIABLES ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

  
bool stringComplete = false; //semaphor for string parsing
String inputString = ""; //empty string for serial input

//creates Servos with name „switch_machine*“
Servo switch_machine1, switch_machine2; 

// Speed Array for quantization of 8Bit to 99 values (Throttle)
byte speedArray [] = {3,5,8,10,13,15,18,20,23,25,28,30,33,35,38,40,43,45,48,50,53,55,58,60,63,65,68,70,73,75,78,80,83,85,88,90,93,95,98,100,103,105,108,110,113,115,118,120,123,125,128,130,133,135,138,140,143,145,148,150,153,155,158,160,163,165,168,170,173,175,178,180,183,185,188,190,193,195,198,200,203,205,208,210,213,215,218,220,223,225,228,230,233,235,238,240,243,248,255};
// Speed State to store Speed of each channel in between loops
byte speedState [] = {0,0,0,0,0,0};

// counter for loop runs - is used in sensor_states()
unsigned long current_time = 0, time_offset = 2000;  //unsigned long to prevent millis() overflow after 32 seconds
int i = 0;

//Sensor stats
unsigned long SENSOR_ACTIVATION[9], SENSOR_DURATION = 2000;  // Time Values to compare Sensor States between iterations and keep state stable
                        // Set default Value
                        //High for no object in from of sensor, Sensor not released
String SENSOR_STATE[] = {"01hz","02hz","03hz","04hz","05hz","06hz","07hz","08hz","09hz","10hz","11hz","12hz","13hz","14hz","15hz","16hz","17hz","18hz","19hz","20hz","21hz","22hz","23hz","24hz","25hz","26hz","27hz","28hz"};


/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////// I/O PINS ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Definition of Output Pins
#define CH0_ENA_PIN 9 // MOTOR-DRIVER L298 Channel0
#define CH0_IN1_PIN A8
#define CH0_IN2_PIN A9
#define CH0_IN3_PIN A10
#define CH0_IN4_PIN A11
#define CH0_ENB_PIN 8
#define CH1_ENA_PIN 11 // MOTOR-DRIVER L298 Channel1
#define CH1_IN1_PIN A4
#define CH1_IN2_PIN A5
#define CH1_IN3_PIN A6
#define CH1_IN4_PIN A7
#define CH1_ENB_PIN 10
#define CH2_ENA_PIN 13 // MOTOR-DRIVER L298 Channel2
#define CH2_IN1_PIN A0
#define CH2_IN2_PIN A1
#define CH2_IN3_PIN A2
#define CH2_IN4_PIN A3
#define CH2_ENB_PIN 12
#define switch1_pin 7 //Switch-Machines
#define switch2_pin 6



// Definition of Input Pins 	  
const byte S1_PIN = 24;  // Definition of IR Sensors (for position tracking)
const byte S2_PIN = 25;
const byte S3_PIN = 26;
const byte S4_PIN = 27;
const byte S5_PIN = 28;
const byte S6_PIN = 29;
const byte S7_PIN = 30;
const byte S8_PIN = 31;
const byte S9_PIN = 32;

 
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////// SETUP ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

void setup()

{
// Initializing COMM
  Wire.begin();
  Serial.begin(9600);
  inputString.reserve(4);

// Initializing Motor-Driver
  pinMode(CH0_ENA_PIN, OUTPUT);	// set output mode
  pinMode(CH0_IN1_PIN, OUTPUT);
  pinMode(CH0_IN2_PIN, OUTPUT);
  pinMode(CH0_IN3_PIN, OUTPUT);
  pinMode(CH0_IN4_PIN, OUTPUT);
  pinMode(CH0_ENB_PIN, OUTPUT);
  pinMode(CH1_ENA_PIN, OUTPUT);
  pinMode(CH1_IN1_PIN, OUTPUT);
  pinMode(CH1_IN2_PIN, OUTPUT);
  pinMode(CH1_IN3_PIN, OUTPUT);
  pinMode(CH1_IN4_PIN, OUTPUT);
  pinMode(CH1_ENB_PIN, OUTPUT);
  pinMode(CH2_ENA_PIN, OUTPUT);
  pinMode(CH2_IN1_PIN, OUTPUT);
  pinMode(CH2_IN2_PIN, OUTPUT);
  pinMode(CH2_IN3_PIN, OUTPUT);
  pinMode(CH2_IN4_PIN, OUTPUT);
  pinMode(CH2_ENB_PIN, OUTPUT);

  
  pinMode(S1_PIN, INPUT);		// set input mode
  pinMode(S2_PIN, INPUT);
  pinMode(S3_PIN, INPUT);
  pinMode(S4_PIN, INPUT);
  pinMode(S5_PIN, INPUT);
  pinMode(S6_PIN, INPUT);
  pinMode(S7_PIN, INPUT);
  pinMode(S8_PIN, INPUT);
  pinMode(S9_PIN, INPUT);


// Set PWM frequency for D9 & D10
// Timer 1 divisor to 256 for PWM frequency of 122.55 Hz
  TCCR1B = TCCR1B & B11111000 | B00000100;   

// Set default direction to NEUTRAL
  digitalWrite(CH0_IN1_PIN, HIGH);	//initial values
  digitalWrite(CH0_IN2_PIN, HIGH); 
  digitalWrite(CH0_IN3_PIN, HIGH);
  digitalWrite(CH0_IN4_PIN, HIGH);
  digitalWrite(CH1_IN1_PIN, HIGH);
  digitalWrite(CH1_IN2_PIN, HIGH); 
  digitalWrite(CH1_IN3_PIN, HIGH);
  digitalWrite(CH1_IN4_PIN, HIGH);
  digitalWrite(CH2_IN1_PIN, HIGH);
  digitalWrite(CH2_IN2_PIN, HIGH); 
  digitalWrite(CH2_IN3_PIN, HIGH);
  digitalWrite(CH2_IN4_PIN, HIGH);

// Initializing Switches - default value
  switch_machine1.attach(switch1_pin); // set putput mode
  switch_machine1.write(45); // initial value for switch machines
  switch_machine2.attach(switch2_pin); 
  switch_machine2.write(135); //mirrored values as installed mirrored
}

 
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////// LOOP ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
void loop(){

 if (stringComplete) {						//if string complete was set, last entered character is "z", this loop will be executed
 
   // LOOP FUNCTIONS 
    if ((inputString.charAt(0) =='r') && (inputString.charAt(1) =='s') && (inputString.charAt(2) =='t')) {    //if "rst" reset arduino
      resetFunc(); //call reset
    }
    
    if (inputString.charAt(0) =='y') {		//if first character is "y", it is interpreted as junction (Y) command and will execute junction()
      if (inputString.charAt(1) =='g') junctions(switch_machine1,45,135); //character position 1: "g-h" will define turnout 7-8
      if (inputString.charAt(1) =='h') junctions(switch_machine2,135,45); //character position 2: 0 or 1 will define direction of switch (0 for straight, outer raduis, 1 for branch change, inner direction)
                                                                          //character position 3: "z" will close the command
    }

    if (inputString.charAt(0) =='k') speedControl(CH0_ENA_PIN, CH0_IN1_PIN, CH0_IN2_PIN, speedState[0]); //if first character is "k-p" it is interpreted as motor driver command, 25 available
    if (inputString.charAt(0) =='l') speedControl(CH0_ENB_PIN, CH0_IN3_PIN, CH0_IN4_PIN, speedState[1]); //valid commands are p.e. kdfz, kdbz, kdsz, k33z
    if (inputString.charAt(0) =='m') speedControl(CH1_ENA_PIN, CH1_IN1_PIN, CH1_IN2_PIN, speedState[2]); //character position 0: k-p, defines channel 0-2, subchannel A or B, overall 25 possible
    if (inputString.charAt(0) =='n') speedControl(CH1_ENB_PIN, CH1_IN3_PIN, CH1_IN4_PIN, speedState[3]); //character position 1: "d", specify direction,
    if (inputString.charAt(0) =='o') speedControl(CH2_ENA_PIN, CH2_IN1_PIN, CH2_IN2_PIN, speedState[4]); //followed by character position 2: "b" for backward, "f" for forward, or "s" for stop
    if (inputString.charAt(0) =='p') speedControl(CH2_ENB_PIN, CH2_IN3_PIN, CH2_IN4_PIN, speedState[5]); //character position 1 and 2: numeric number 0-99 defines drivespeed (roughly 255/99) -> see SpeedArray
                                                                                                         //character position 3: "z" will close the command     

    
  
   // RESET INPUT STRING
    inputString = "";						//mandatory reset of input string to get new commands
    stringComplete = false;
  }


  // LOGIC - is be done every loop

  SENSOR_STATE[0] = sensor_change(S1_PIN,SENSOR_STATE[19],20);  // Sensors will only get updated when change
  SENSOR_STATE[1] = sensor_change(S2_PIN,SENSOR_STATE[20],21);  // will keep value until time expired, see SENSOR_DURATION
  SENSOR_STATE[2] = sensor_change(S3_PIN,SENSOR_STATE[21],22);
  SENSOR_STATE[3] = sensor_change(S4_PIN,SENSOR_STATE[22],23);
  SENSOR_STATE[4] = sensor_change(S5_PIN,SENSOR_STATE[23],24);
  SENSOR_STATE[5] = sensor_change(S6_PIN,SENSOR_STATE[24],25);
  SENSOR_STATE[6] = sensor_change(S7_PIN,SENSOR_STATE[25],26);
  SENSOR_STATE[7] = sensor_change(S8_PIN,SENSOR_STATE[26],27);
  SENSOR_STATE[8] = sensor_change(S9_PIN,SENSOR_STATE[27],28);

  //sensor_state();                                          //regulary update of all sensors
}



/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////// FUNCTIONS ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
void serialEvent() {						//used to get serial input all the time
  if (Serial.available()) {
    char inChar = (char)Serial.read();
    inputString += inChar;
    if (inChar == 'z') {
      stringComplete = true;
    }
  }
}


 
