package deej

import (
	"fmt"
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
}

// NewBrightnessController initializes a new BrightnessController instance
func NewBrightnessController() *BrightnessController {
	return &BrightnessController{
		automaticBrightness: false,
		photoresistorLeft:   0,
		photoresistorRight:  0,
		currentBrightness:   80, // Default brightness (50%)
		prevEncoderValue:    0,  // Initialize previous encoder value to 0
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
		bc.toggleAutomaticBrightness()
	}

	// Update photoresistor values
	bc.photoresistorLeft = photoresistorLeft
	bc.photoresistorRight = photoresistorRight
}

func (bc *BrightnessController) adjustBrightness(encoderValue int) {
	// Compute brightness adjustment based on encoder value change
	encoderChange := encoderValue - bc.prevEncoderValue

	// Adjust brightness based on encoder change
	bc.currentBrightness += encoderChange * 10

	// Ensure brightness is within 0% and 100%
	if bc.currentBrightness < 0 {
		bc.currentBrightness = 0
	} else if bc.currentBrightness > 100 {
		bc.currentBrightness = 100
	}

	fmt.Printf("Manual brightness adjusted to %d%%\n", bc.currentBrightness)
	// Implement logic to update brightness value in the system
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
