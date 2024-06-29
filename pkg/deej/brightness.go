package deej

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// BrightnessController handles PC brightness control
type BrightnessController struct {
	automaticBrightness          int // Use 0 for disabled, 1 for enabled
	photoresistorLeft            int
	photoresistorRight           int
	avgPhotoresistor             int           // Average of photoresistor values
	prevAvgPhotoresistor         int           // Previous average of photoresistor values
	prevEncoderValue             int           // Previous absolute encoder value
	lastButtonPress              int           // Last button press state (0 or 1)
	photoresistorChangeThreshold int           // Threshold for photoresistor change to trigger brightness adjustment
	lastBrightnessChangeTime     time.Time     // Last time brightness was changed
	brightnessChangeDelay        time.Duration // Delay between brightness changes
}

// NewBrightnessController initializes a new BrightnessController instance
func NewBrightnessController() *BrightnessController {
	return &BrightnessController{
		automaticBrightness:          0, // Initialize to 0 (disabled)
		photoresistorLeft:            0,
		photoresistorRight:           0,
		avgPhotoresistor:             0,
		prevAvgPhotoresistor:         0,
		prevEncoderValue:             0,               // Initialize previous encoder value to 0
		lastButtonPress:              0,               // Initialize last button press state to 0 (assuming no button press)
		photoresistorChangeThreshold: 50,              // Default threshold for photoresistor change
		lastBrightnessChangeTime:     time.Now(),      // Initialize last change time
		brightnessChangeDelay:        1 * time.Second, // 1 second delay between brightness changes
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

	// Calculate average photoresistor value
	bc.avgPhotoresistor = (photoresistorLeft + photoresistorRight) / 2

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

	if bc.automaticBrightness == 1 {
		bc.adjustBrightnessAutomatically()
	}
}

func (bc *BrightnessController) adjustBrightness(encoderValue int) {
	// Compute brightness adjustment based on encoder value change
	encoderChange := encoderValue - bc.prevEncoderValue

	if encoderChange > 0 {
		for i := 0; i < encoderChange; i++ {
			bc.sendKeyPress("Alt+PgUp")
		}
	} else {
		for i := 0; i > encoderChange; i-- {
			bc.sendKeyPress("Alt+PgDown")
		}
	}
}

func (bc *BrightnessController) toggleAutomaticBrightness() {
	if bc.automaticBrightness == 0 {
		bc.automaticBrightness = 1
		fmt.Println("Automatic brightness control enabled")
	} else {
		bc.automaticBrightness = 0
		fmt.Println("Automatic brightness control disabled")
	}
}

func (bc *BrightnessController) adjustBrightnessAutomatically() {
	// Compare current average with previous average and check if change exceeds threshold
	if abs(bc.avgPhotoresistor-bc.prevAvgPhotoresistor) >= bc.photoresistorChangeThreshold {
		// If the delay period has passed since the last adjustment, adjust brightness
		if time.Since(bc.lastBrightnessChangeTime) >= bc.brightnessChangeDelay {
			// If current average is greater than previous, increase brightness
			if bc.avgPhotoresistor > bc.prevAvgPhotoresistor {
				bc.sendKeyPress("Alt+PgUp")
			} else { // Decrease brightness
				bc.sendKeyPress("Alt+PgDown")
			}

			// Update previous average and last brightness change time
			bc.prevAvgPhotoresistor = bc.avgPhotoresistor
			bc.lastBrightnessChangeTime = time.Now()
		}
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

// Helper function to get absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
