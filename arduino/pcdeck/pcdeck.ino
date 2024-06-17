const int NUM_SLIDERS = 3;
const int NUM_MUTE_BUTTONS = 3;
const int analogInputs[NUM_SLIDERS] = {A0, A1, A2};
const int digitalInputs[NUM_MUTE_BUTTONS] = {26, 24, 22};
int analogSliderValues[NUM_SLIDERS];
int digitalButtonValues[NUM_MUTE_BUTTONS];
bool muteStates[NUM_MUTE_BUTTONS] = {false, false, false}; // Track mute states
String builtString;

void setup() { 
  for (int i = 0; i < NUM_SLIDERS; i++) {
    pinMode(analogInputs[i], INPUT);
  }
  for (int i = 0; i < NUM_MUTE_BUTTONS; i++) {
    pinMode(digitalInputs[i], INPUT_PULLUP);
  }

  Serial.begin(9600);
}

void loop() {
  updateSliderValues();
  updateMuteButtonValues();
  sendValues(); // Send combined data
  delay(10);
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

  Serial.println(builtString);
}
