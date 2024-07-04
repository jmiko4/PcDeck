#include <Encoder.h>
#include <Bounce2.h>

unsigned long previousMillis = 0;  // Stores the last time the loop ran
const long interval = 10;          // Interval for the non-blocking delay in milliseconds

const int E1_CLK_PIN = 2;  // CLK (S1) (Clock) pin (encoder 1)
const int E1_DT_PIN = 3;   // DT (S2) (Data) pin (encoder 1)
const int E1_SW_PIN = 4;   // SW (Switch) pin (encoder 1)

const int PHOTORESISTOR_PIN1 = A3; // Photoresistor 1 pin
const int PHOTORESISTOR_PIN2 = A4; // Photoresistor 2 pin
const int CHANGE_THRESHOLD = 50;   // Threshold for detecting significant changes in light

Encoder encoder1(E1_DT_PIN, E1_CLK_PIN); // Create encoder 1 object
Bounce keyDebouncerE1 = Bounce(); // Create Bounce object for the switch

long encoder1Value = 0;
bool encoder1KeyState = HIGH;  // Default to HIGH (not pressed)
int photoresistor1Value = 0;
int photoresistor2Value = 0;

const int NUM_SLIDERS = 3;
const int NUM_MUTE_BUTTONS = 3;
const int analogInputs[NUM_SLIDERS] = {A0, A1, A2};
const int digitalInputs[NUM_MUTE_BUTTONS] = {22, 24, 26}; // Updated button pins
const int ledOutputs[NUM_MUTE_BUTTONS] = {23, 25, 27}; // LED pins
int analogSliderValues[NUM_SLIDERS];
int digitalButtonValues[NUM_MUTE_BUTTONS] = {HIGH, HIGH, HIGH}; // Initialize to HIGH because of INPUT_PULLUP
bool muteStates[NUM_MUTE_BUTTONS] = {false, false, false}; // Track mute states
String builtString;

void setup() {
  for (int i = 0; i < NUM_SLIDERS; i++) {
    pinMode(analogInputs[i], INPUT);
  }
  for (int i = 0; i < NUM_MUTE_BUTTONS; i++) {
    pinMode(digitalInputs[i], INPUT_PULLUP);
    pinMode(ledOutputs[i], OUTPUT);
    digitalWrite(ledOutputs[i], LOW); // Turn off LEDs initially
  }


  pinMode(E1_SW_PIN, INPUT_PULLUP);
  
  keyDebouncerE1.attach(E1_SW_PIN);
  keyDebouncerE1.interval(10);  // Debounce interval in milliseconds 


  Serial.begin(9600);
}

void loop() {

  keyDebouncerE1.update(); // Update the debouncer
  encoder1Value = encoder1.read() / 4;  // Read value and divide by 4 to get correct count
  encoder1KeyState = keyDebouncerE1.read();

  unsigned long currentMillis = millis();  // Get the current time
  // Check if the interval has passed
  if (currentMillis - previousMillis >= interval) {
    previousMillis = currentMillis;  // Save the current time
    
    updateSliderValues();
    updateMuteButtonValues();
    sendValues(); // Send combined data
  }


}

void updateSliderValues() {
  for (int i = 0; i < NUM_SLIDERS; i++) {
    if (!muteStates[i]) {
      analogSliderValues[i] = analogRead(analogInputs[i]);
    }

  }
}

void updateMuteButtonValues() {
  for (int i = 0; i < NUM_MUTE_BUTTONS; i++) {
    int currentButtonValue = digitalRead(digitalInputs[i]);
    if (currentButtonValue == LOW && digitalButtonValues[i] == HIGH) { // Button pressed
      muteStates[i] = !muteStates[i]; // Toggle mute state
      digitalWrite(ledOutputs[i], muteStates[i] ? HIGH : LOW); // Toggle LED
    }
    digitalButtonValues[i] = currentButtonValue;
  }
}

void sendValues() {
  builtString = String("");


  // Append slider values or 0 if muted
  for (int i = 0; i < NUM_SLIDERS; i++) {
    if (muteStates[i]) {
      builtString += String(0);
    } else {
      builtString += String(analogSliderValues[i]);
    }
    if (i < NUM_SLIDERS - 1) {
      builtString += "|";
    }
  }

  builtString += "$"; // Add delimiter between sliders and buttons

  // Append button values
  for (int i = 0; i < NUM_MUTE_BUTTONS; i++) {
    builtString += String(digitalButtonValues[i]);
    if (i < NUM_MUTE_BUTTONS - 1) {
      builtString += "|";
    }
  }


  builtString += "$0|0$"; // Add encoder 2 values. Currently it is not used.

  photoresistor1Value = analogRead(PHOTORESISTOR_PIN1);
  photoresistor2Value = analogRead(PHOTORESISTOR_PIN2);

  builtString += String(encoder1Value)+"|"+String(!encoder1KeyState)+"|"+String(photoresistor1Value)+"|"+String(photoresistor2Value);

  builtString += "$0|0|0|0|0|0"; // Add key states (not yet implemented, send 0 for now).


  Serial.println(builtString);
}
