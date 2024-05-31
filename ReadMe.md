# PcDeck
### Justin Mikolajcik, Elizabeth Hunter, Vladislovas Karalius
## Project Scope

The goal of this project is to build a streamdeck type device that allows you to quickly change app specific volumes (discord, game volume, spotify), use custom macro keys to do frequent tasks (mute microphone, open discord, toggle night mode), automatically or manually change brightness, and add a display to easily monitor PC metrics (CPU temperature, CPU usage, GPU usage, Time, Etc). All of these features would be incredibly useful for those who play games or PC enthusiasts.

### Main controller
The main controller will be implemented using Arduino Mega 2560 Rev3 microcontroller board based on the ATmega2560. 
It has the following features relevant to this project:
54 digital pins
16 analog pins
SPI communication on pins 50 (MISO), 51 (MOSI), 52 (SCK), 53 (SS)
USB connectivity and serial communication over a virtual port
Operating voltage 5V, USB can be used as the power supply
DC Current per I/O Pin 20 mA

### Peripherals
OLED  display (single color), qty. 1
Resolution 128x64
Size 2.4”
SSD1309 driver chip
SPI communication

Slide potentiometers for volume control, qty. 3
Size 90x20 mm
Voltage 5V
Linear type
Resistance 10k

Rotary encoders KY-040, qty. 2
Voltage 5V
Pulse number 20
Knob with a switch
Knob cap from aluminum alloy
Knob cap diameter 15 mm, height 16.5 mm

Photoresistors (photocells), qty, 2 + resistors 1 kΩ, qty. 2
Voltage 5V
Resistance in darkness 50 kΩ
Resistance in bright light 500 kΩ
Connected together with 1 kΩ fixed resistor


### Other hardware
Momentary Push Button Switch SPST AC250V/3A AC125V/6A Mini Off(ON) ON 5 Colour with Pre-soldered Wires R13-507-5C-X., qty 6
Small buttons, qty. 3
LEDs, qty. 3
Resistors for LEDs 330 Ω, qty. 3
Plastic or cardboard case, qty. 1
Breadboard, wires




### User interface
The user interface will contain the following elements:
OLED display 128x64
3 slide potentiometers
2 rotary encoders with knob switches
6 large buttons
3 small buttons
3 LEDs

## Connectivity
Communication with a PC. We will use serial communication only. The board will not directly function as an USB keyboard. So, a significant part of this project will be writing some PC program that would accept the communication from the microcontroller and perform certain actions: 
Control PC volume
Increase/decrease monitor brightness
Generate keyboard scancodes
Send information about PC stats to the microcontroller

## Lessons to Learn
Lessons to learn include: 
Programming an Arduino in python or go
Connecting and displaying to a display
Modifying windows settings with Arduino inputs
Gathering information typically displayed in task manager
Working with various analog and digital hardware elements, such as SPI displays, potentiometers, photoresistors, rotary encoders
Using various team programs such as Monday.com, LucidChart/draw.io/Miro, Github 

## Roles and Responsibilities
### Team Leader
#### Justin Mikolajcik
Set up project in Monday.com and invite team
Manage project tasks in Monday.com
Submits group assignments
Establishes meeting schedule (1-2 times per week)
### Hardware Lead
#### Vladislovas Karalius
Final decision maker on hardware selection
Hardware block diagram owner
Gathers necessary hardware (purchase or loan)
### Software Lead
#### Elizabeth Hunter
Final decision maker on software architecture
Software block diagram owner


## General Requirements
System size shall be no larger than 8”x6”.
The system shall  weight no more than 250g
Coding the PC part: Python or Go
Coding the microcontroller part: Arduino IDE
The following libraries will be used in the Arduino IDE:
“Encoder.h” for handling encoder input
“Bounce2.h” for buttons debouncing


## Interface Requirements
The system shall communicate with the PC over a virtual serial port, at 9600 8N1 settings.
The system shall provide 5 analog pins for slide potentiometers and photoresistors.
The system shall provide 14 digital pins for buttons and LEDs.
The system shall provide 4 digital pins with interrupts for rotary encoders.
The system shall provide 4 pins for SPI communication.
The system shall use a 2.4” 128x64 SPI display for information received from the PC.
The system shall provide 3 slide potentiometers for volume control (1 for master volume and 2 for separate apps).
The system shall provide 3 small buttons for muting (1 for master volume and 2 for separate apps).
The system shall provide 3 LEDs to display the muting state.
The system shall provide 6 large buttons for keyboard macros.
The system shall provide 1 rotary encoder for brightness control.
The system shall provide 1 rotary encoder for mic volume control. The knob switch shall be used for mic muting.



## Functional Requirements
Features include:
Change app specific volume with volume sliders
Mute app volume with button
Multi-function larger buttons able to execute specific commands
Automatically or manually adjust brightness via rotary encoder or photoresistors
Display PC metrics on the included display
Two photoresistors shall be used for automatic brightness control. They will be placed some distance apart and their inputs will be averaged. When the change of detected brightness level reaches a certain threshold, the microcontroller will send a command to the PC to increase or decrease the monitor brightness.
It will be possible to adjust the base brightness level using a rotary encoder. 
A knob switch shall be used to disable the brightness control and set the monitor brightness to a predetermined level.


