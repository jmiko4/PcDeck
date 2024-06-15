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
	currentBrightness   int
	prevEncoderValue    int // Previous absolute encoder value
	lastButtonPress     int // Last button press state (0 or 1)
}

// NewBrightnessController initializes a new BrightnessController instance
func NewBrightnessController() *BrightnessController {
	return &BrightnessController{
		automaticBrightness: false,
		photoresistorLeft:   0,
		photoresistorRight:  0,
		currentBrightness:   90, // Default brightness (90%)
		prevEncoderValue:    0,  // Initialize previous encoder value to 0
		lastButtonPress:     0,  // Initialize last button press state to 0 (assuming no button press)
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

	// Adjust brightness based on encoder change
	bc.currentBrightness += encoderChange * 5

	// Ensure brightness is within 0% and 100%
	if bc.currentBrightness < 0 {
		bc.currentBrightness = 0
	} else if bc.currentBrightness > 100 {
		bc.currentBrightness = 100
	}

	fmt.Printf("Manual brightness adjusted to %d%%\n", bc.currentBrightness)
	// Implement logic to update brightness value in the system
	bc.updateSystemBrightness()
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

// updateSystemBrightness updates the brightness value in the system
func (bc *BrightnessController) updateSystemBrightness() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("brightnessctl", "set", fmt.Sprintf("%d%%", bc.currentBrightness))
	case "windows":
		// Windows-specific command or script to adjust brightness
		// For example, using PowerShell:
		cmd = exec.Command("powershell", "-Command", fmt.Sprintf("(Get-WmiObject -Namespace root/WMI -Class WmiMonitorBrightnessMethods).WmiSetBrightness(1,%d)", bc.currentBrightness))
	case "darwin":
		cmd = exec.Command("brightness", fmt.Sprintf("%.2f", float64(bc.currentBrightness)/100.0))
	default:
		fmt.Printf("Unsupported OS: %s\n", runtime.GOOS)
		return
	}

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to update system brightness: %v\n", err)
	}
}
