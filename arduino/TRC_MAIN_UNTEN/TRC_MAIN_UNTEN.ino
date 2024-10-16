//-------------------------------------//
//    TILO RAILWAY COMPANY             //
//    UNTERE EBENE                     //
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
Servo switch_machine1, switch_machine2, switch_machine3, switch_machine4, switch_machine5, switch_machine6; 

// Speed Array for quantization of 8Bit to 99 values (Throttle)
// linear Array 
// byte speedArray [] = {3,5,8,10,13,15,18,20,23,25,28,30,33,35,38,40,43,45,48,50,53,55,58,60,63,65,68,70,73,75,78,80,83,85,88,90,93,95,98,100,103,105,108,110,113,115,118,120,123,125,128,130,133,135,138,140,143,145,148,150,153,155,158,160,163,165,168,170,173,175,178,180,183,185,188,190,193,195,198,200,203,205,208,210,213,215,218,220,223,225,228,230,233,235,238,240,243,248,255};
// linear Array with offset
byte speedArray [] = {0,40,42,44,47,49,51,53,55,58,60,62,64,66,69,71,73,75,77,80,82,84,86,88,91,93,95,97,99,102,104,106,108,110,113,115,117,119,121,124,126,128,130,132,135,137,139,141,143,146,148,150,152,154,157,159,161,163,165,168,170,172,174,176,179,181,183,185,187,190,192,194,196,198,201,203,205,207,209,212,214,216,218,220,223,225,227,229,231,234,236,238,240,242,245,247,249,251,253,255};
// parabolic array with linear ramp in beginning (y=x²*0,026)
//byte speedArray [] = {1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,42,44,46,48,50,53,55,57,60,62,65,68,70,73,76,79,82,84,87,91,94,97,100,103,106,110,113,117,120,124,127,131,135,139,142,146,150,154,158,162,166,171,175,179,183,188,192,197,201,206,211,215,220,225,230,235,240,245,250,255};
// parabolic array with linear ramp in beginning and offset (y=x²*0.022+40)
// byte speedArray [] = {0,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,59,60,61,63,64,65,67,69,70,72,73,75,77,79,81,83,85,87,89,91,93,95,97,99,102,104,107,109,111,114,117,119,122,125,1127,130,133,136,139,142,145,148,151,154,157,160,164,167,170,174,177,181,184,188,192,195,199,203,207,210,214,218,222,226,230,234,239,243,247,251,255};

// Speed State to store Speed of each channel in between loops
byte PWMspeedState [] = {0,0,0,0,0,0,0,0,0,0};

// counter for loop runs - is used in sensor_states()
unsigned long current_time = 0, time_offset = 2000;  //unsigned long to prevent millis() overflow after 32 seconds
// int i = 0;

//Sensor stats
                                                              //unsigned long to prevent millis() overflow after 32 seconds
// counter for loop runs - is used in sensor_change()
unsigned long SENSOR_ACTIVATION[19], SENSOR_DURATION = 1600;  // Time Values to compare Sensor States between iterations and keep state stable
// counter for loop runs - is used in sensor_state_active()
unsigned long SENSOR_ACTIVATION2[19], SENSOR_DURATION2 = 3000;  // Time Values to compare Sensor States between iterations and keep state stable
                        // Set default Value
                        //High for no object in from of sensor, Sensor not released

// (default) States of all commands send and received by the arduino / (at startup)
String SENSOR_STATE[] = {"01hz","02hz","03hz","04hz","05hz","06hz","07hz","08hz","09hz","10hz","11hz","12hz","13hz","14hz","15hz","16hz","17hz","18hz","19hz"};
String SWITCH_STATE[] = {"ya0z","yb0z","yc0z","yd0z","ye0z","yf0z"};
String DIRECTION_STATE [] = {"adsz","bdsz","cdsz","ddsz","edsz","fdsz","gdsz","hdsz","idsz","jdsz"};
String SPEED_STATE [] = {"a00z","b00z","c00z","d00z","e00z","f00z","g00z","h00z","i00z","j00z"};
String SIGNAL_STATE [] = {"xa0z","xb0z"};

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////// I/O PINS ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Definition of Output Pins
#define CH0_ENA_PIN 6 // MOTOR-DRIVER L298 Channel0
#define CH0_IN1_PIN A12
#define CH0_IN2_PIN A13
#define CH0_IN3_PIN A14
#define CH0_IN4_PIN A15
#define CH0_ENB_PIN 7
#define CH1_ENA_PIN 9 // MOTOR-DRIVER L298 Channel1
#define CH1_IN1_PIN A8
#define CH1_IN2_PIN A9
#define CH1_IN3_PIN A10
#define CH1_IN4_PIN A11
#define CH1_ENB_PIN 8
#define CH2_ENA_PIN 11 // MOTOR-DRIVER L298 Channel2
#define CH2_IN1_PIN A4
#define CH2_IN2_PIN A5
#define CH2_IN3_PIN A6
#define CH2_IN4_PIN A7
#define CH2_ENB_PIN 10
#define CH3_ENA_PIN 2 // MOTOR-DRIVER L298 Channel3
#define CH3_IN1_PIN 20
#define CH3_IN2_PIN A1
#define CH3_IN3_PIN A2
#define CH3_IN4_PIN A3
#define CH3_ENB_PIN 12
#define CH4_ENA_PIN 43 // MOTOR-DRIVER L298 Channel4
#define CH4_IN1_PIN 50
#define CH4_IN2_PIN 49
#define CH4_IN3_PIN 48
#define CH4_IN4_PIN 47
#define CH4_ENB_PIN 44
#define switch1_pin 45 //Switch-Machines
#define switch2_pin 46
#define switch3_pin 13
#define switch4_pin 3
#define switch5_pin 4
#define switch6_pin 5
#define sig_grn_pin 16 //Signal Manipulation 
#define sig_red_pin 17


// Definition of Input Pins 	  
const byte S1_PIN = 24;  // Definition of IR Sensors (for position tracking)
const byte S2_PIN = 25;
const byte S3_PIN = 26;
const byte S4_PIN = 27;
const byte S5_PIN = 28;
const byte S6_PIN = 29;
const byte S7_PIN = 23;
const byte S8_PIN = 31;
const byte S9_PIN = 32;
const byte S10_PIN = 33;
const byte S11_PIN = 34;
const byte S12_PIN = 35;
const byte S13_PIN = 36;
const byte S14_PIN = 37;
const byte S15_PIN = 38;
const byte S16_PIN = 39;
const byte S17_PIN = 40;
const byte S18_PIN = 41;
const byte S19_PIN = 42;
 
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
  pinMode(CH3_ENA_PIN, OUTPUT);
  pinMode(CH3_IN1_PIN, OUTPUT);
  pinMode(CH3_IN2_PIN, OUTPUT);
  pinMode(CH3_IN3_PIN, OUTPUT);
  pinMode(CH3_IN4_PIN, OUTPUT);
  pinMode(CH3_ENB_PIN, OUTPUT);
  pinMode(CH4_ENA_PIN, OUTPUT);
  pinMode(CH4_IN1_PIN, OUTPUT);
  pinMode(CH4_IN2_PIN, OUTPUT);
  pinMode(CH4_IN3_PIN, OUTPUT);
  pinMode(CH4_IN4_PIN, OUTPUT);
  pinMode(CH4_ENB_PIN, OUTPUT);
  pinMode(sig_grn_pin, OUTPUT);
  pinMode(sig_red_pin, OUTPUT);
  
  pinMode(S1_PIN, INPUT);		// set input mode
  pinMode(S2_PIN, INPUT);
  pinMode(S3_PIN, INPUT);
  pinMode(S4_PIN, INPUT);
  pinMode(S5_PIN, INPUT);
  pinMode(S6_PIN, INPUT);
  pinMode(S7_PIN, INPUT);
  pinMode(S8_PIN, INPUT);
  pinMode(S9_PIN, INPUT);
  pinMode(S10_PIN, INPUT);
  pinMode(S11_PIN, INPUT);
  pinMode(S12_PIN, INPUT);
  pinMode(S13_PIN, INPUT);
  pinMode(S14_PIN, INPUT);
  pinMode(S15_PIN, INPUT);
  pinMode(S16_PIN, INPUT);
  pinMode(S17_PIN, INPUT);
  pinMode(S18_PIN, INPUT);
  pinMode(S19_PIN, INPUT);

// Set PWM frequency for D9 & D10 (for UNO)
// Set PWM frequency for D11 & D12 (for MEGA)
// Timer 1 divisor to 256 for PWM frequency of 122.55 Hz
  TCCR1B = TCCR1B & B11111000 | B00000100; 
//====================================================for MEGA only
//// Set PWM frequency for D9 & D10 //for MEGA
//TCCR2B = TCCR2B & B11111000 | B00000110;  // for PWM frequency of 122.55 Hz
//// Set PWM frequency for D2 & D3 & D5 //for MEGA
//TCCR4B = TCCR4B & B11111000 | B00000100;   // for PWM frequency of 122.55 Hz
//// Set PWM frequency for D6 & D7 & D8 //for MEGA
//TCCR4B = TCCR4B & B11111000 | B00000100;   // for PWM frequency of 122.55 Hz
//// Set PWM frequency for D44 & D45 & D46 //for MEGA
//TCCR5B = TCCR5B & B11111000 | B00000100;  // for PWM frequency of 122.55 Hz

// Set default direction to NEUTRAL
  digitalWrite(CH0_IN1_PIN, LOW);	//initial values
  digitalWrite(CH0_IN2_PIN, LOW); 
  digitalWrite(CH0_IN3_PIN, LOW);
  digitalWrite(CH0_IN4_PIN, LOW);
  digitalWrite(CH1_IN1_PIN, LOW);
  digitalWrite(CH1_IN2_PIN, LOW); 
  digitalWrite(CH1_IN3_PIN, LOW);
  digitalWrite(CH1_IN4_PIN, LOW);
  digitalWrite(CH2_IN1_PIN, LOW);
  digitalWrite(CH2_IN2_PIN, LOW); 
  digitalWrite(CH2_IN3_PIN, LOW);
  digitalWrite(CH2_IN4_PIN, LOW);
  digitalWrite(CH3_IN1_PIN, LOW);
  digitalWrite(CH3_IN2_PIN, LOW); 
  digitalWrite(CH3_IN3_PIN, LOW);
  digitalWrite(CH3_IN4_PIN, LOW);
  digitalWrite(CH4_IN1_PIN, LOW);
  digitalWrite(CH4_IN2_PIN, LOW); 
  digitalWrite(CH4_IN3_PIN, LOW);
  digitalWrite(CH4_IN4_PIN, LOW);

// Initializing Switches - default value
  switch_machine1.attach(switch1_pin); // set putput mode
  switch_machine1.write(45); // initial value for switch machines
  switch_machine2.attach(switch2_pin); 
  switch_machine2.write(135); //mirrored values as installed mirrored
  switch_machine3.attach(switch3_pin); 
  switch_machine3.write(45); 
  switch_machine4.attach(switch4_pin); 
  switch_machine4.write(135); 
  switch_machine5.attach(switch5_pin); 
  switch_machine5.write(135);
  switch_machine6.attach(switch6_pin); 
  switch_machine6.write(155); 

// Set Signal Manipulation Pins to low (inactive)
  digitalWrite(sig_grn_pin, LOW);
  digitalWrite(sig_red_pin, LOW);
}

 
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////// LOOP ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
void loop(){
 if (stringComplete ) {						//if string complete was set, last entered character is "z", this loop will be executed
 
   // LOOP FUNCTIONS 
    if ((inputString.charAt(0) =='r') && (inputString.charAt(1) =='s') && (inputString.charAt(2) =='t')) {    //if "rst" reset arduino
      resetFunc(); //call reset
    }
    
    if (inputString.charAt(0) =='y') {		//if first character is "y", it is interpreted as junction (Y) command and will execute junction()
      if (inputString.charAt(1) =='a') junctions(switch_machine1,45,135,0); //character position 1: "a-f" will define turnout 1-6
      if (inputString.charAt(1) =='b') junctions(switch_machine2,135,45,1); //character position 2: 0 or 1 will define direction of switch (0 for straight, outer raduis, 1 for branch change, inner direction)
      if (inputString.charAt(1) =='c') junctions(switch_machine3,45,135,2); //character position 3: "z" will close the command
      if (inputString.charAt(1) =='d') junctions(switch_machine4,135,45,3); //last number/value corresponds with indexnumber to store values/states in arrays
      if (inputString.charAt(1) =='e') junctions(switch_machine5,135,45,4);
      if (inputString.charAt(1) =='f') junctions(switch_machine6,155,25,5);
    }

    if (inputString.charAt(0) == 'x') { //if first character is "x", it is interpreted as signal manipulation attempt for independet arduino nano.
      if (inputString.charAt(1) =='a') signal_manipulation(sig_grn_pin, "a",0);             // stay green
      if (inputString.charAt(1) =='b') signal_manipulation(sig_red_pin, "b",1);             // stay red
    }

    if (inputString.charAt(0) == 'w') { //if first charachter is "w", states will be send.
      if ((inputString.charAt(1) =='s') && (inputString.charAt(2) =='e')) sensor_state_ondemand(); //Sensor States.
      if ((inputString.charAt(1) =='s') && (inputString.charAt(2) =='w')) junction_state_ondemand(); //Junction States.
      if ((inputString.charAt(1) =='b') && (inputString.charAt(2) =='l')) direction_state_ondemand(); //Block States / Direction per Block.
      if ((inputString.charAt(1) =='s') && (inputString.charAt(2) =='p')) speed_state_ondemand(); //Speed States per Block.   
      if ((inputString.charAt(1) =='s') && (inputString.charAt(2) =='i')) signal_state_ondemand(); //Signal States. 
    }

    if (inputString.charAt(0) =='a') speedControl(CH0_ENA_PIN, CH0_IN1_PIN, CH0_IN2_PIN, PWMspeedState[0], 0); //if first character is "a-j" it is interpreted as motor driver command, 25 available
    if (inputString.charAt(0) =='b') speedControl(CH0_ENB_PIN, CH0_IN3_PIN, CH0_IN4_PIN, PWMspeedState[1], 1); //valid commands are p.e. adfz, adbz, adsz, a33z
    if (inputString.charAt(0) =='c') speedControl(CH1_ENA_PIN, CH1_IN1_PIN, CH1_IN2_PIN, PWMspeedState[2], 2); //character position 0: a-j, defines channel 1-5, subchannel A or B, overall 25 possible
    if (inputString.charAt(0) =='d') speedControl(CH1_ENB_PIN, CH1_IN3_PIN, CH1_IN4_PIN, PWMspeedState[3], 3); //character position 1: "d", specify direction,
    if (inputString.charAt(0) =='e') speedControl(CH2_ENA_PIN, CH2_IN1_PIN, CH2_IN2_PIN, PWMspeedState[4], 4); //                       followed by character position 2: "b" for backward, "f" for forward, or "s" for stop
    if (inputString.charAt(0) =='f') speedControl(CH2_ENB_PIN, CH2_IN3_PIN, CH2_IN4_PIN, PWMspeedState[5], 5); //character position 1 and 2: numeric number 0-99 defines drivespeed (roughly 255/99) -> see SpeedArray
    if (inputString.charAt(0) =='g') speedControl(CH3_ENA_PIN, CH3_IN1_PIN, CH3_IN2_PIN, PWMspeedState[6], 6); //character position 3: "z" will close the command     
    if (inputString.charAt(0) =='h') speedControl(CH3_ENB_PIN, CH3_IN3_PIN, CH3_IN4_PIN, PWMspeedState[7], 7); //last number/value corresponds with indexnumber to store values/states in arrays
    if (inputString.charAt(0) =='i') speedControl_noPWM(CH4_ENA_PIN, CH4_IN1_PIN, CH4_IN2_PIN, 8); //PWM not possible (digital pin              // no PWM -> no Speedstate 
    if (inputString.charAt(0) =='j') speedControl_noPWM(CH4_ENB_PIN, CH4_IN3_PIN, CH4_IN4_PIN, 9); //PWM should be possible on this pin (couldn't verify during test -> noPWM instead        // only full throttle 00z corresponds with 0, everything else with full throttle

    
  
   // RESET INPUT STRING
    inputString = "";						//mandatory reset of input string to get new commands
    stringComplete = false;
  }


  // LOGIC - is be done every loop
  
  SENSOR_STATE[0] = sensor_change(S1_PIN,SENSOR_STATE[0],1);  // regular update of all  U P D A T E D  sensors
  SENSOR_STATE[1] = sensor_change(S2_PIN,SENSOR_STATE[1],2);  // Sensors will only get updated when change
  SENSOR_STATE[2] = sensor_change(S3_PIN,SENSOR_STATE[2],3);  // will keep value until time expired, see SENSOR_DURATION
  SENSOR_STATE[3] = sensor_change(S4_PIN,SENSOR_STATE[3],4);
  SENSOR_STATE[4] = sensor_change(S5_PIN,SENSOR_STATE[4],5);
  SENSOR_STATE[5] = sensor_change(S6_PIN,SENSOR_STATE[5],6);
  SENSOR_STATE[6] = sensor_change(S7_PIN,SENSOR_STATE[6],7);
  SENSOR_STATE[7] = sensor_change(S8_PIN,SENSOR_STATE[7],8);
  SENSOR_STATE[8] = sensor_change(S9_PIN,SENSOR_STATE[8],9);
  SENSOR_STATE[9] = sensor_change(S10_PIN,SENSOR_STATE[9],10);
  SENSOR_STATE[10] = sensor_change(S11_PIN,SENSOR_STATE[10],11);
  SENSOR_STATE[11] = sensor_change(S12_PIN,SENSOR_STATE[11],12);
  SENSOR_STATE[12] = sensor_change(S13_PIN,SENSOR_STATE[12],13);
  SENSOR_STATE[13] = sensor_change(S14_PIN,SENSOR_STATE[13],14);
  SENSOR_STATE[14] = sensor_change(S15_PIN,SENSOR_STATE[14],15);
  SENSOR_STATE[15] = sensor_change(S16_PIN,SENSOR_STATE[15],16);
  SENSOR_STATE[16] = sensor_change(S17_PIN,SENSOR_STATE[16],17);
  SENSOR_STATE[17] = sensor_change(S18_PIN,SENSOR_STATE[17],18);
  SENSOR_STATE[18] = sensor_change(S19_PIN,SENSOR_STATE[18],19);

  sensor_state_active(S1_PIN,SENSOR_STATE[0],1);  // regular update of all  A C T I V E  sensors
  sensor_state_active(S2_PIN,SENSOR_STATE[1],2);  // will repeat output every SENSOR_DURATION2
  sensor_state_active(S3_PIN,SENSOR_STATE[2],3);
  sensor_state_active(S4_PIN,SENSOR_STATE[3],4);
  sensor_state_active(S5_PIN,SENSOR_STATE[4],5);
  sensor_state_active(S6_PIN,SENSOR_STATE[5],6);
  sensor_state_active(S7_PIN,SENSOR_STATE[6],7);
  sensor_state_active(S8_PIN,SENSOR_STATE[7],8);
  sensor_state_active(S9_PIN,SENSOR_STATE[8],9);
  sensor_state_active(S10_PIN,SENSOR_STATE[9],10);
  sensor_state_active(S11_PIN,SENSOR_STATE[10],11);
  sensor_state_active(S12_PIN,SENSOR_STATE[11],12);
  sensor_state_active(S13_PIN,SENSOR_STATE[12],13);
  sensor_state_active(S14_PIN,SENSOR_STATE[13],14);
  sensor_state_active(S15_PIN,SENSOR_STATE[14],15);
  sensor_state_active(S16_PIN,SENSOR_STATE[15],16);
  sensor_state_active(S17_PIN,SENSOR_STATE[16],17);
  sensor_state_active(S18_PIN,SENSOR_STATE[17],18);
  sensor_state_active(S19_PIN,SENSOR_STATE[18],19);

  //regular_sensor_state();                                          //regular update of all sensors  A F T E R   S P E C I F I C   T I M E 
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


 
