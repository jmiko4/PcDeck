// include libraries for display
#include <Wire.h>
#include <Adafruit_GFX.h>
#include <Adafruit_SSD1306.h>

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


const int NUM_SLIDERS = 3;
const int analogInputs[NUM_SLIDERS] = {A0, A1, A2};

int analogSliderValues[NUM_SLIDERS];

void setup() { 
  for (int i = 0; i < NUM_SLIDERS; i++) {
    pinMode(analogInputs[i], INPUT);
  }

  Serial.begin(9600);

  initializeDisplay();

}

void loop() {
  updateSliderValues();
  sendSliderValues(); // Actually send data (all the time)
  // printSliderValues(); // For debug
  delay(10);

  if (Serial.available() > 0) {
    // Read the incoming string
    String receivedData = Serial.readStringUntil('\n'); // Read until newline character
    
    // Parse the received string
    parseString(receivedData);
  }
}

void updateSliderValues() {
  for (int i = 0; i < NUM_SLIDERS; i++) {
     analogSliderValues[i] = analogRead(analogInputs[i]);
  }
}

void sendSliderValues() {
  String builtString = String("");

  for (int i = 0; i < NUM_SLIDERS; i++) {
    builtString += String((int)analogSliderValues[i]);

    if (i < NUM_SLIDERS - 1) {
      builtString += String("|");
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