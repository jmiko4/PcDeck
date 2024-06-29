package deej

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// BrightnessController handles PC brightness control
type BrightnessController struct {
	automaticBrightness bool
	photoresistorLeft   int
	photoresistorRight  int
	prevEncoderValue    int // Previous absolute encoder value
	lastButtonPress     int // Last button press state (0 or 1)
}

// NewBrightnessController initializes a new BrightnessController instance
func NewBrightnessController() *BrightnessController {
	return &BrightnessController{
		automaticBrightness: false,
		photoresistorLeft:   0,
		photoresistorRight:  0,
		prevEncoderValue:    0, // Initialize previous encoder value to 0
		lastButtonPress:     0, // Initialize last button press state to 0 (assuming no button press)
	}
}

// HandleBrightnessInfo handles brightness information received over serial
func (bc *BrightnessController) HandleBrightnessInfo(info string) {
	// Split the info string by |
	fields := strings.Split(info, "|")

	if len(fields) != 4 {
		// Invalid format, ignore
		return
	}

	// Parse the fields
	encoderValue, _ := strconv.Atoi(fields[0])
	buttonPress, _ := strconv.Atoi(fields[1])
	photoresistorLeft, _ := strconv.Atoi(fields[2])
	photoresistorRight, _ := strconv.Atoi(fields[3])

	// Check if encoder value changed
	if encoderValue != bc.prevEncoderValue {
		// Adjust brightness based on encoder value change
		bc.adjustBrightness(encoderValue)
		// Update previous encoder value
		bc.prevEncoderValue = encoderValue
	}

	// Toggle automatic brightness control based on button press
	if buttonPress == 1 {
		// Only toggle if lastButtonPress was 0 (indicating no consecutive 1s)
		if bc.lastButtonPress == 0 {
			bc.toggleAutomaticBrightness()
		}
		bc.lastButtonPress = 1
	} else {
		bc.lastButtonPress = 0
	}

	// Update photoresistor values
	bc.photoresistorLeft = photoresistorLeft
	bc.photoresistorRight = photoresistorRight
}

func (bc *BrightnessController) adjustBrightness(encoderValue int) {
	// Compute brightness adjustment based on encoder value change
	encoderChange := encoderValue - bc.prevEncoderValue

	if encoderChange > 0 {
		for i := 0; i < encoderChange; i++ {
			bc.sendKeyPress("Alt+PgUp")
		}
	} else {
		for i := 0; i < -encoderChange; i++ {
			bc.sendKeyPress("Alt+PgDown")
		}
	}
}

func (bc *BrightnessController) toggleAutomaticBrightness() {
	bc.automaticBrightness = !bc.automaticBrightness
	if bc.automaticBrightness {
		fmt.Println("Automatic brightness control enabled")
		// Implement logic to enable automatic brightness control
	} else {
		fmt.Println("Automatic brightness control disabled")
		// Implement logic to disable automatic brightness control
	}
}

func (bc *BrightnessController) sendKeyPress(key string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdotool", "key", key)
	case "windows":
		// Windows-specific command or script to send key press
		if key == "Alt+PgUp" {
			cmd = exec.Command("powershell", "-Command", "$wshell = New-Object -ComObject wscript.shell; $wshell.SendKeys('%{PGUP}')")
		} else {
			cmd = exec.Command("powershell", "-Command", "$wshell = New-Object -ComObject wscript.shell; $wshell.SendKeys('%{PGDN}')")
		}
	case "darwin":
		// macOS-specific command or script to send key press
		if key == "Alt+PgUp" {
			cmd = exec.Command("osascript", "-e", "tell application \"System Events\" to key code 116 using {option down}") // Option+Page Up
		} else {
			cmd = exec.Command("osascript", "-e", "tell application \"System Events\" to key code 121 using {option down}") // Option+Page Down
		}
	default:
		fmt.Printf("Unsupported OS: %s\n", runtime.GOOS)
		return
	}

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to send key press: %v\n", err)
	} else {
		fmt.Printf("Generated key press: %s\n", key)
	}
}
