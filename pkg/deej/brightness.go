package deej

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

// Import the COM libraries (for Windows)
var (
	ole32              = syscall.NewLazyDLL("ole32.dll")
	procCoInitializeEx = ole32.NewProc("CoInitializeEx")
	procCoUninitialize = ole32.NewProc("CoUninitialize")
)

// CoInitializeEx initializes the COM library for use by the calling thread
func CoInitializeEx(pvReserved uintptr, dwCoInit uint32) (err error) {
	hr, _, _ := procCoInitializeEx.Call(
		pvReserved,
		uintptr(dwCoInit))
	if hr != 0 {
		err = fmt.Errorf("CoInitializeEx failed: hr = 0x%x", hr)
	}
	return
}

// CoUninitialize uninitializes the COM library
func CoUninitialize() {
	procCoUninitialize.Call()
}

const (
	COINIT_APARTMENTTHREADED = 0x2 // Single-threaded apartment
	COINIT_MULTITHREADED     = 0x0 // Multi-threaded apartment
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
		currentBrightness:   90, // Default brightness (90%)
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
		// Initialize COM for this thread
		err := CoInitializeEx(0, COINIT_APARTMENTTHREADED)
		if err != nil {
			fmt.Printf("Failed to initialize COM: %v\n", err)
			return
		}
		defer CoUninitialize()

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
