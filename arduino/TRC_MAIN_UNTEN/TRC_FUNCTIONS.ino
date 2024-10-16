//-------------------------------------//
//    TILO RAILWAY COMPANY             //
//-------------------------------------//
//    TRC_FUNCTIONS.ino                //
//    FUNCTION UNIT                    //
//    PWM 122 Hz                       //
//    V.3.2020                         //
//    based on ArduinoRailwayProject   //
//-------------------------------------//


////////////////////////////////////////////////////////////////////////////////////////////////////////////////// JUNCTIONS ////////////////////////////////////////////////////////////////////////////////////////////////
void junctions(Servo switch_machine,int MIN_ANGLE, int MAX_ANGLE, int state_index) {
  char c = inputString.charAt(1);
// Switch
    if (inputString.charAt(2) =='0') { // Branch direction
      if (inputString.charAt(3)=='z') {
        Serial.print((String)"y"+c+"0z");
        switch_machine.write(MIN_ANGLE);
        SWITCH_STATE[state_index]=((String)"y"+c+"0z");
      }
    }
    if (inputString.charAt(2) =='1') { // Throw direction;
      if (inputString.charAt(3)=='z') {
        Serial.print((String)"y"+c+"1z");
        switch_machine.write(MAX_ANGLE); 
        SWITCH_STATE[state_index]=((String)"y"+c+"1z");
      }
    } 
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////// SPEED CONTROL ////////////////////////////////////////////////////////////////////////////////////////////////
byte speedControl(const byte ENABL, const byte IN1, const byte IN2, int speedTrain, int state_index) {
  int one = 0, ten = 0, flag1 = 0, flag2 = 0;
  // Speed Selection
  if (inputString.charAt(1)=='0') {ten = 0;  flag1 = 1;}
  if (inputString.charAt(1)=='1') {ten = 10; flag1 = 1;}
  if (inputString.charAt(1)=='2') {ten = 20; flag1 = 1;}
  if (inputString.charAt(1)=='3') {ten = 30; flag1 = 1;}
  if (inputString.charAt(1)=='4') {ten = 40; flag1 = 1;}
  if (inputString.charAt(1)=='5') {ten = 50; flag1 = 1;}
  if (inputString.charAt(1)=='6') {ten = 60; flag1 = 1;}
  if (inputString.charAt(1)=='7') {ten = 70; flag1 = 1;}
  if (inputString.charAt(1)=='8') {ten = 80; flag1 = 1;}
  if (inputString.charAt(1)=='9') {ten = 90; flag1 = 1;}

  if (inputString.charAt(2)=='0') {one = 0; flag2 = 1;}
  if (inputString.charAt(2)=='1') {one = 1; flag2 = 1;}
  if (inputString.charAt(2)=='2') {one = 2; flag2 = 1;}
  if (inputString.charAt(2)=='3') {one = 3; flag2 = 1;}
  if (inputString.charAt(2)=='4') {one = 4; flag2 = 1;}
  if (inputString.charAt(2)=='5') {one = 5; flag2 = 1;}
  if (inputString.charAt(2)=='6') {one = 6; flag2 = 1;}
  if (inputString.charAt(2)=='7') {one = 7; flag2 = 1;}
  if (inputString.charAt(2)=='8') {one = 8; flag2 = 1;}
  if (inputString.charAt(2)=='9') {one = 9; flag2 = 1;}


  if (inputString.charAt(1)!='d'){ // execution only in case of speed control, skip this block in case of direction control
    if (inputString.charAt(2)=='0' && inputString.charAt(1)=='0') {
      speedTrain = 0;
      if (inputString.charAt(3)=='z') Serial.print(inputString);        //send back inpiutString for execution command success
      SPEED_STATE[state_index] = inputString;
      }
    else {
      speedTrain = speedArray[ten+one-1];
      if ((flag1 == 1) && (flag2 == 1) && inputString.charAt(3)=='z') Serial.print(inputString);        //send back inpiutString for execution command success
      SPEED_STATE[state_index] = inputString;
      }
    analogWrite(ENABL, speedTrain); // Throttle
    return speedTrain;
  }
    
  // Direction and Stop in1,in2,enable
  if (inputString.charAt(1) =='d') {
    if (inputString.charAt(2) =='f') { // (f) Forward
      digitalWrite(IN1, HIGH);
      digitalWrite(IN2, LOW);
      if (inputString.charAt(3)=='z') Serial.print(inputString);        //send back inpiutString for execution command success 
      DIRECTION_STATE[state_index] = inputString;
    }
    if (inputString.charAt(2) =='b') { // (b) Backward
      digitalWrite(IN1, LOW);
      digitalWrite(IN2, HIGH);
      if (inputString.charAt(3)=='z') Serial.print(inputString);        //send back inpiutString for execution command success
      DIRECTION_STATE[state_index] = inputString;
    }
    if (inputString.charAt(2) =='s') { // (s) Stop button
      speedTrain = 0;                  // Motor free spinning, after some time it will be blocked and set signal correctly
      analogWrite(ENABL, speedTrain); // Throttle
      //delay(1000);
      digitalWrite(IN1, LOW);
      digitalWrite(IN2, LOW);
      if (inputString.charAt(3)=='z') Serial.print(inputString);        //send back inpiutString for execution command success
      DIRECTION_STATE[state_index] = inputString;
    } 
    if (inputString.charAt(2) =='x') { // (x) emergency stop - MOTOR BLOCKED
      speedTrain = 0;
      analogWrite(ENABL, speedTrain); // Throttle
      digitalWrite(IN1, LOW);
      digitalWrite(IN2, LOW);
      if (inputString.charAt(3)=='z') Serial.print(inputString);        //send back inpiutString for execution command success
      DIRECTION_STATE[state_index] = inputString;
    }
  }

  return speedTrain;
}  




byte speedControl_noPWM(const byte ENABL, const byte IN1, const byte IN2, int state_index) {
  int one = 0, ten = 0, flag1 = 0, flag2 = 0;
  // Speed Selection for string validation only
  if (inputString.charAt(1)=='0') {ten = 0;  flag1 = 1;}
  if (inputString.charAt(1)=='1') {ten = 10; flag1 = 1;}
  if (inputString.charAt(1)=='2') {ten = 20; flag1 = 1;}
  if (inputString.charAt(1)=='3') {ten = 30; flag1 = 1;}
  if (inputString.charAt(1)=='4') {ten = 40; flag1 = 1;}
  if (inputString.charAt(1)=='5') {ten = 50; flag1 = 1;}
  if (inputString.charAt(1)=='6') {ten = 60; flag1 = 1;}
  if (inputString.charAt(1)=='7') {ten = 70; flag1 = 1;}
  if (inputString.charAt(1)=='8') {ten = 80; flag1 = 1;}
  if (inputString.charAt(1)=='9') {ten = 90; flag1 = 1;}

  if (inputString.charAt(2)=='0') {one = 0; flag2 = 1;}
  if (inputString.charAt(2)=='1') {one = 1; flag2 = 1;}
  if (inputString.charAt(2)=='2') {one = 2; flag2 = 1;}
  if (inputString.charAt(2)=='3') {one = 3; flag2 = 1;}
  if (inputString.charAt(2)=='4') {one = 4; flag2 = 1;}
  if (inputString.charAt(2)=='5') {one = 5; flag2 = 1;}
  if (inputString.charAt(2)=='6') {one = 6; flag2 = 1;}
  if (inputString.charAt(2)=='7') {one = 7; flag2 = 1;}
  if (inputString.charAt(2)=='8') {one = 8; flag2 = 1;}
  if (inputString.charAt(2)=='9') {one = 9; flag2 = 1;}

  if (inputString.charAt(1)!='d'){ // execution only in case of speed control, skip this block in case of direction control
    if (inputString.charAt(2)=='0' && inputString.charAt(1)=='0') {
      digitalWrite(ENABL, LOW);
      if (inputString.charAt(3)=='z') Serial.print(inputString);        //send back inpiutString for execution command success
      SPEED_STATE[state_index] = inputString;
      }
    else {
      digitalWrite(ENABL, HIGH);
      if ((flag1 == 1) && (flag2 == 1) && inputString.charAt(3)=='z') Serial.print(inputString);        //send back inpiutString for execution command success
      SPEED_STATE[state_index] = inputString;
      }
  }
    
  // Direction and Stop in1,in2,enable
  if (inputString.charAt(1) =='d') {
    if (inputString.charAt(2) =='f') { // (f) Forward
      digitalWrite(IN1, HIGH);
      digitalWrite(IN2, LOW); 
      if (inputString.charAt(3)=='z') Serial.print(inputString);        //send back inpiutString for execution command success
      DIRECTION_STATE[state_index] = inputString;
    }
    if (inputString.charAt(2) =='b') { // (b) Backward
      digitalWrite(IN1, LOW);
      digitalWrite(IN2, HIGH); 
      if (inputString.charAt(3)=='z') Serial.print(inputString);        //send back inpiutString for execution command success
      DIRECTION_STATE[state_index] = inputString;
    }
    if (inputString.charAt(2) =='s') { // (s) Stop button
                                       // Motor free spinning, after some time it will be blocked and set signal correctly
      digitalWrite(ENABL, LOW);
      digitalWrite(IN1, LOW);
      digitalWrite(IN2, LOW);
      Serial.print(inputString);        //send back inpiutString for execution command success
      DIRECTION_STATE[state_index] = inputString;
    } 
    if (inputString.charAt(2) =='x') { // (x) emergency stop - MOTOR BLOCKED
      digitalWrite(ENABL, LOW);
      digitalWrite(IN1, LOW);
      digitalWrite(IN2, LOW);
      if (inputString.charAt(3)=='z') Serial.print(inputString);        //send back inpiutString for execution command success
      DIRECTION_STATE[state_index] = inputString;
    }
  }
}  


////////////////////////////////////////////////////////////////////////////////////////////////////////////////// SENSORS ////////////////////////////////////////////////////////////////////////////////////////////////
String sensor_change(const byte PIN, String STATE,int SENSOR_INT){
//SENSOR 
  int sread = digitalRead(PIN);                                       // temp variable sread to prevent to call so often digitalRead
  String SENSOR_NR = "", TEMP_STATE = STATE;                          // create temp variables for Sensor Number to convert it in String and for State to do comparisons between iterations

  if (SENSOR_INT < 10){                                               // Sensor Number (int) converted to string
    SENSOR_NR = "0"+((String)SENSOR_INT);                             // if below 10, a leading 0 is added to match the naming convention
    }
  else{
    SENSOR_NR = (String)SENSOR_INT;                                   // if higher 10, no conversion needed as already 2 numbers
  }
                                                                      
                                                                      // Monoflop Part - debouncer part
  if (sread == LOW){                                                  // Monoflop Part - debouncer part
    SENSOR_ACTIVATION[SENSOR_INT-1] = millis();                       // if Sensor is activated, the actual time is stored to ACTIVATION
    TEMP_STATE = SENSOR_NR+"lz";                                   // the state is set to active
  }
  if (sread == HIGH && (millis() > SENSOR_ACTIVATION[SENSOR_INT-1]+SENSOR_DURATION)){     // if the sensor not active and the time window (Sensor Duration = 2 seconds or so) elapsed
    TEMP_STATE = SENSOR_NR+"hz";                                   // the state is set to inactive
  }
    
  if (TEMP_STATE == SENSOR_NR+"hz" &&
  STATE == SENSOR_NR+"lz") {  // will send inactive if submitted state and new state (as well as state after time window) was changed from active
    STATE = SENSOR_NR+"hz";
    Serial.print(SENSOR_NR+"hz");
    }
  if (TEMP_STATE == SENSOR_NR+"lz" && STATE == SENSOR_NR+"hz") {  // will send active if submitted state and new state (as well as state after time window) was changed from inactive
    STATE = SENSOR_NR+"lz";
    Serial.print(SENSOR_NR+"lz");
    }

    return STATE;                                                     // return new state to have it available for next iteration
}    

void sensor_state_active(const byte PIN, String STATE,int SENSOR_INT){
  String SENSOR_NR = "";
    if (SENSOR_INT < 10){                                               // Sensor Number (int) converted to string
      SENSOR_NR = "0"+((String)SENSOR_INT);                             // if below 10, a leading 0 is added to match the naming convention
    }
    else{
      SENSOR_NR = (String)SENSOR_INT;                                   // if higher 10, no conversion needed as already 2 numbers
    }
    
    if (STATE == SENSOR_NR+"lz" && (millis() > SENSOR_ACTIVATION2[SENSOR_INT-1]+SENSOR_DURATION2)){
      Serial.print(STATE);
      SENSOR_ACTIVATION2[SENSOR_INT-1] = millis();
    }
}

void regular_sensor_state(){
                                                          //Backup solution - Print eyery x seconds the last Sensor State    
  if (millis() > current_time + time_offset){             // will be executed after x loops -- Backup Solution
    current_time = millis();
    Serial.print(SENSOR_STATE[0]); 
    Serial.print(SENSOR_STATE[1]);
    Serial.print(SENSOR_STATE[2]);
    Serial.print(SENSOR_STATE[3]);
    Serial.print(SENSOR_STATE[4]);
    Serial.print(SENSOR_STATE[5]);
    Serial.print(SENSOR_STATE[6]);
    Serial.print(SENSOR_STATE[7]);
    Serial.print(SENSOR_STATE[8]);
    Serial.print(SENSOR_STATE[9]);
    Serial.print(SENSOR_STATE[10]);
    Serial.print(SENSOR_STATE[11]);
    Serial.print(SENSOR_STATE[12]);
    Serial.print(SENSOR_STATE[13]);
    Serial.print(SENSOR_STATE[14]);
    Serial.print(SENSOR_STATE[15]);
    Serial.print(SENSOR_STATE[16]);
    Serial.print(SENSOR_STATE[17]);
    Serial.print(SENSOR_STATE[18]);
  }         
}




////////////////////////////////////////////////////////////////////////////////////////////////////////////////// ACTUAL STATES ON DEMAND ////////////////////////////////////////////////////////////////////////////////////////////



void sensor_state_ondemand(){
                                                          //Print eyery last Sensor State on demand
        Serial.print(SENSOR_STATE[0]); 
        Serial.print(SENSOR_STATE[1]);
        Serial.print(SENSOR_STATE[2]);
        Serial.print(SENSOR_STATE[3]);
        Serial.print(SENSOR_STATE[4]);
        Serial.print(SENSOR_STATE[5]);
        Serial.print(SENSOR_STATE[6]);
        Serial.print(SENSOR_STATE[7]);
        Serial.print(SENSOR_STATE[8]);
        Serial.print(SENSOR_STATE[9]);
        Serial.print(SENSOR_STATE[10]);
        Serial.print(SENSOR_STATE[11]);
        Serial.print(SENSOR_STATE[12]);
        Serial.print(SENSOR_STATE[13]);
        Serial.print(SENSOR_STATE[14]);
        Serial.print(SENSOR_STATE[15]);
        Serial.print(SENSOR_STATE[16]);
        Serial.print(SENSOR_STATE[17]);
        Serial.print(SENSOR_STATE[18]);
  }         

void junction_state_ondemand(){
                                                          //Print eyery junction State on demand
        Serial.print(SWITCH_STATE[0]);    
        Serial.print(SWITCH_STATE[1]);  
        Serial.print(SWITCH_STATE[2]);  
        Serial.print(SWITCH_STATE[3]);  
        Serial.print(SWITCH_STATE[4]);  
        Serial.print(SWITCH_STATE[5]);                                                       
} 

void direction_state_ondemand(){
                                                          //Print eyery direction State on demand
        Serial.print(DIRECTION_STATE[0]);                                
        Serial.print(DIRECTION_STATE[1]);
        Serial.print(DIRECTION_STATE[2]);
        Serial.print(DIRECTION_STATE[3]);
        Serial.print(DIRECTION_STATE[4]);
        Serial.print(DIRECTION_STATE[5]);
        Serial.print(DIRECTION_STATE[6]);
        Serial.print(DIRECTION_STATE[7]);
        Serial.print(DIRECTION_STATE[8]);
        Serial.print(DIRECTION_STATE[9]);
}

void speed_state_ondemand(){
                                                          //Print eyery speed State on demand
        Serial.print(SPEED_STATE[0]);                                             
        Serial.print(SPEED_STATE[1]);
        Serial.print(SPEED_STATE[2]);
        Serial.print(SPEED_STATE[3]);
        Serial.print(SPEED_STATE[4]);
        Serial.print(SPEED_STATE[5]);
        Serial.print(SPEED_STATE[6]);
        Serial.print(SPEED_STATE[7]);
        Serial.print(SPEED_STATE[8]);
        Serial.print(SPEED_STATE[9]);
}

void signal_state_ondemand(){
                                                          //Print eyery signal State on demand
        Serial.print(SIGNAL_STATE[0]);                                              
        Serial.print(SIGNAL_STATE[1]);
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////// SIGNAL MANIPULATION ////////////////////////////////////////////////////////////////////////////////////////////////

void signal_manipulation(const byte PIN,String channel,int state_index){
    String sread = "";
    int state = 0;
    
    if (inputString.charAt(2)=='0') digitalWrite(PIN, LOW); 
    if (inputString.charAt(2)=='1') digitalWrite(PIN, HIGH);

    state = digitalRead(PIN);
    if (state == HIGH){sread = "1";}
    if (state == LOW){sread = "0";}
    Serial.print("x"+channel+state+"z");        //send back inpiutString for execution command success
    SIGNAL_STATE[state_index] = "x"+channel+state+"z";
}
