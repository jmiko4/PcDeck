// include libraries for display
#include <Wire.h>
#include <Adafruit_GFX.h>
#include <Adafruit_SSD1306.h>
#include <Encoder.h>
#include <Bounce2.h>

// define constants for display pins
#define SCREEN_WIDTH 128
#define SCREEN_HEIGHT 64
#define OLED_MOSI   49
#define OLED_CLK   50
#define OLED_CS    51
#define OLED_DC    52
#define OLED_RESET 53
Adafruit_SSD1306 display(SCREEN_WIDTH, SCREEN_HEIGHT,
  OLED_MOSI, OLED_CLK, OLED_DC, OLED_RESET, OLED_CS);

#define LINE_1_START 0
#define LINE_2_START 16
#define LINE_3_START 32
#define LINE_4_START 48

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

const int NUM_KEYS = 6;
const int keyPins[NUM_KEYS] = {14, 15, 16, 17, 18, 19}; // Key button pins
int keyValues[NUM_KEYS] = {HIGH, HIGH, HIGH, HIGH, HIGH, HIGH}; // Initialize to HIGH because of INPUT_PULLUP

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

  for (int i = 0; i < NUM_KEYS; i++) {
    pinMode(keyPins[i], INPUT_PULLUP);
  }

  pinMode(E1_SW_PIN, INPUT_PULLUP);
  
  keyDebouncerE1.attach(E1_SW_PIN);
  keyDebouncerE1.interval(10);  // Debounce interval in milliseconds 

  Serial.begin(9600);

  initializeDisplay();

}

void loop() {
  if (Serial.available() > 0) {
    // Read the incoming string
    String receivedData = Serial.readStringUntil('\n'); // Read until newline character
    
    // Parse the received string
    parseString(receivedData);
  }
  keyDebouncerE1.update(); // Update the debouncer
  encoder1Value = encoder1.read() / 4;  // Read value and divide by 4 to get correct count
  encoder1KeyState = keyDebouncerE1.read();

  unsigned long currentMillis = millis();  // Get the current time
  // Check if the interval has passed
  if (currentMillis - previousMillis >= interval) {
    previousMillis = currentMillis;  // Save the current time
    
    updateSliderValues();
    updateMuteButtonValues();
    updateKeyValues();
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

void updateKeyValues() {
  for (int i = 0; i < NUM_KEYS; i++) {
    keyValues[i] = digitalRead(keyPins[i]);
  }
}

void sendValues() {
  builtString = "";

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

  builtString += String(encoder1Value) + "|" + String(!encoder1KeyState) + "|" + String(photoresistor1Value) + "|" + String(photoresistor2Value);

  builtString += "$"; // Add delimiter before key states

  // Append key states
  for (int i = 0; i < NUM_KEYS; i++) {
    builtString += String(keyValues[i] == LOW ? 1 : 0); // Convert HIGH/LOW to 1/0
    if (i < NUM_KEYS - 1) {
      builtString += "|";
    }
  }

  Serial.println(builtString);
}

void printSliderValues() {
  for (int i = 0; i < NUM_SLIDERS; i++) {
    String printedString = String("Slider #") + String(i + 1) + String(": ") + String(analogSliderValues[i]) + String(" mV");
    Serial.write(printedString.c_str());

    if (i < NUM_SLIDERS - 1) {
      Serial.write(" | ");
    } else {
      Serial.write("\n");
    }
  }
}



void initializeDisplay(){
// Initialize the display
  if (!display.begin(SSD1306_SWITCHCAPVCC, 0x3C)) { // Address 0x3C for 128x64
    Serial.println(F("SSD1306 allocation failed"));
    for (;;); // Don't proceed, loop forever
  }
}

String parseString(String data) {
  // Variables to hold parsed values
  String value1;
  String value2;

  // Find the first comma
  int separationIndex1 = data.indexOf('|');

  // Extract the command
  if (separationIndex1 != -1) {

    value1 = data.substring(0, separationIndex1);

    value2 = data.substring(separationIndex1 + 1);

    updateDisplay(value1, value2);
      // You can now use the command and values as needed
    } 
    return value1, value2;
  } 

  void updateDisplay(String dateAndTime, String cpuLoad){
// Set text size and color
  
  
  display.clearDisplay();

  display.setTextSize(1); // Adjust as needed for your font
  display.setTextColor(SSD1306_WHITE);

  display.setCursor(0, LINE_1_START);
  display.println(dateAndTime);

  
  display.setCursor(0, LINE_3_START);
  display.println(cpuLoad);
  

  // Display the text
  display.display();
}
