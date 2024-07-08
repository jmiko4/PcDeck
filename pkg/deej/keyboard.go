package deej

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

// KeyboardController handles keyboard press events
type KeyboardController struct {
	logger *zap.SugaredLogger
}

// NewKeyboardController initializes a new KeyboardController instance
func NewKeyboardController(logger *zap.SugaredLogger) *KeyboardController {
	return &KeyboardController{
		logger: logger.Named("keyboard"),
	}
}

// HandleKeyboardInfo processes the keyboard information received from serial input
func (kc *KeyboardController) HandleKeyboardInfo(data string) error {
	kc.logger.Debugw("Received keyboard info", "data", data)

	// Split the data by |
	keyValues := strings.Split(data, "|")

	if len(keyValues) != 6 {
		return errors.New("keyboard: invalid data format")
	}

	// Map of key index to the corresponding key combination for Windows
	keyMap := map[int]string{
		0: "Ctrl+Shift+Esc",
		1: "discord",
		2: "Chatgpt",
		3: "F8",
		4: "Alt+Left",
		5: "F10",
	}

	// Detect the OS for platform-specific handling
	osType := detectOSType()

	// Iterate over key values and trigger key presses based on the mapping
	for idx, valueStr := range keyValues {
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return fmt.Errorf("keyboard: failed to parse key value %s", valueStr)
		}

		if value == 1 {
			keyCombination, ok := keyMap[idx]
			if !ok {
				return fmt.Errorf("keyboard: unknown key index %d", idx)
			}

			kc.logger.Infow("Sending key press", "key", keyCombination)

			// Send key press based on the OS type
			err := kc.sendKeyPress(osType, keyCombination)
			if err != nil {
				kc.logger.Errorw("Failed to send key press", "key", keyCombination, "error", err)
				return err
			}
		}
	}

	return nil
}

// detectOSType detects the current operating system type
func detectOSType() string {
	// For simplicity, assuming it's always Windows for this example
	return "windows"
}

// sendKeyPress sends a key press based on the OS type and key combination
func (kc *KeyboardController) sendKeyPress(osType, keyCombination string) error {
	var cmd *exec.Cmd

	switch osType {
	case "windows":
		// Windows-specific command to send key press
		switch keyCombination {
		case "Alt+Left":
			cmd = exec.Command("powershell", "-Command", "$wshell = New-Object -ComObject wscript.shell; $wshell.SendKeys('%{LEFT}')")
		case "Ctrl+Shift+Esc":
			cmd = exec.Command("powershell", "-Command", "$wshell = New-Object -ComObject wscript.shell; $wshell.SendKeys('^+{ESC}')")
		case "discord": 
            cmd = exec.Command("powershell", "-Command", `Start-Process "C:\Users\Gaming\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Discord Inc\Discord.lnk"`)
		case "Chatgpt":
			cmd = exec.Command("powershell", "-Command", `Start-Process "https://www.chatgpt.com"`)
		case "F8":
			cmd = exec.Command("powershell", "-Command", "$wshell = New-Object -ComObject wscript.shell; $wshell.SendKeys('{F8}')")
		case "F10":
			cmd = exec.Command("powershell", "-Command", "$wshell = New-Object -ComObject wscript.shell; $wshell.SendKeys('{F10}')")
		default:
			return fmt.Errorf("keyboard: unsupported key combination for Windows: %s", keyCombination)
		}
	default:
		return fmt.Errorf("keyboard: unsupported OS: %s", osType)
	}

	// Execute the command
	if cmd != nil {
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("keyboard: failed to execute command for key press: %w", err)
		}
	}

	return nil
}
