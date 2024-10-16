//-------------------------------------//
//    TILO RAILWAY COMPANY             //
//    Streckensignal                   //
//-------------------------------------//
//    V.2.2020                         //
//-------------------------------------//


/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////// VARIABLES ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
unsigned long SENSOR_ACTIVATION[2], SENSOR_DURATION = 200, SIGNAL_DURATION = 7000;  // Time Values to compare Sensor States between iterations and keep state stable
String SENSOR_STATE[] = {"high","high"};
bool red_flag = false;
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////// I/O PINS ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Definition of Input Pins     
const byte S1_PIN = 10;  // Definition of IR Sensors (for Signal trigger)
const byte S2_PIN = 11;
const byte sig_grn_pin = 4; // Definition of Signal Override from Master Arduino (Signal Manipulation)
const byte sig_red_pin = 5;

// Definition of Output Pins
const byte SIGNAL_RED = 9; // LED PINS
const byte SIGNAL_YEL = 8;
const byte SIGNAL_GRN = 7;
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////// SETUP ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
void setup() {

// Initializing Motor-Driver
  pinMode(SIGNAL_RED, OUTPUT);  // set output mode
  pinMode(SIGNAL_YEL, OUTPUT);
  pinMode(SIGNAL_GRN, OUTPUT);

  pinMode(S1_PIN, INPUT);    // set input mode
  pinMode(S2_PIN, INPUT);
  pinMode(sig_grn_pin, INPUT);
  pinMode(sig_red_pin, INPUT);


// Set default Signal to GREEN
  digitalWrite(SIGNAL_RED, LOW);  //initial values
  digitalWrite(SIGNAL_YEL, LOW); 
  digitalWrite(SIGNAL_GRN, HIGH);
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////// LOOP ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
void loop() {

  int red, grn;
  grn = digitalRead(sig_grn_pin);
  red = digitalRead(sig_red_pin);
  
  if (grn == HIGH) {                                                  // force green if set to high. no change will be accepted (will win over red)
    digitalWrite(SIGNAL_GRN, HIGH);
    digitalWrite(SIGNAL_RED, LOW);
  }
  else if (red == HIGH){                                              // force red if set to high. no change will be accepted
    digitalWrite(SIGNAL_GRN, LOW);
    digitalWrite(SIGNAL_RED, HIGH);
    red_flag = true;
  }
  else if (red == LOW && red_flag == true) {                          // return to initial values after override with red was set to low again or reset function after change. 
                                                                      // last function that will be executed so it is guarantueed it wll set back to green light, a bit fishy but working
    digitalWrite(SIGNAL_RED, LOW);                                    //initial values
    digitalWrite(SIGNAL_GRN, HIGH);
    red_flag = false;    
  }
  else {
    SENSOR_STATE[0] = sensor_change(S1_PIN,SENSOR_STATE[0],1,red);    // Check Sensor States and set back to green light after forced red was resetted
    SENSOR_STATE[1] = sensor_change(S2_PIN,SENSOR_STATE[1],2,red);
  }

  red, grn = 0;                                                       // reset forced state, otherwise no normal change will accepted in future
}

String sensor_change(const byte PIN, String STATE,int SENSOR_INT,int red){
//SENSOR 
  int sread = digitalRead(PIN);                                       // temp variable sread to prevent to call so often digitalRead
  String TEMP_STATE = STATE;                                          // create temp variables for Sensor Number to convert it in String and for State to do comparisons between iterations                                                                   
                                                                      // Monoflop Part - debouncer part
  if (sread == LOW){                                                  // Monoflop Part - debouncer part
    SENSOR_ACTIVATION[SENSOR_INT-1] = millis();                       // if Sensor is activated, the actual time is stored to ACTIVATION
    TEMP_STATE = "low";                                               // the state is set to active
  }
  
  if (sread == HIGH && (millis() > SENSOR_ACTIVATION[SENSOR_INT-1]+SENSOR_DURATION)){     // if the sensor not active and the time window (Sensor Duration = 2 seconds or so) elapsed
    TEMP_STATE = "high";                                              // the state is set to inactive
  }
    
  if (TEMP_STATE == "high" && STATE == "low") {  // will send inactive if submitted state and new state (as well as state after time window) was changed from active
    STATE = "high";
            delay(SIGNAL_DURATION);
    digitalWrite(SIGNAL_RED, LOW);  //initial values
    digitalWrite(SIGNAL_GRN, HIGH);
    }
  if (TEMP_STATE == "low" && STATE == "high") {  // will send active if submitted state and new state (as well as state after time window) was changed from inactive
    STATE = "low";
    digitalWrite(SIGNAL_GRN, LOW);
    digitalWrite(SIGNAL_RED, HIGH);
    }

    return STATE;                                                     // return new state to have it available for next iteration
}    
