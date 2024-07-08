package deej

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

// SendSystemData sends the CPU load and time string over the serial connection every second
func (sio *SerialIO) SendSystemData() {
	// Ensure we have an active serial connection
	if !sio.connected {
		sio.logger.Warn("Not connected to serial, can't send data")
		return
	}

	// Send the CPU load and time string over the serial connection every second
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:

				// Format the time string
				// European
				// currentTime := time.Now().Format("06.01.02 15:04")
				// American
				currentTime := time.Now().Format("01/02/06 3:04 PM")

				// Get CPU load
				percentages, err := cpu.Percent(0, false)
				if err != nil {
					sio.logger.Warnw("Failed to get CPU load", "error", err)
					// Send the "Error" string if CPU load cannot be retrieved
					if _, err := sio.conn.Write([]byte("Error\r\n")); err != nil {
						sio.logger.Warnw("Failed to send error message over serial", "error", err)
					}
					continue
				}

				// Format the system data string
				systemDataString := fmt.Sprintf("%s|CPU: %2.0f%%\r\n", currentTime, percentages[0])

				// Send the CPU load and time string
				if _, err := sio.conn.Write([]byte(systemDataString)); err != nil {
					sio.logger.Warnw("Failed to send CPU load and time over serial", "error", err)
				}
			case <-sio.stopChannel:
				return // Stop sending data if requested to stop
			}
		}
	}()
}
