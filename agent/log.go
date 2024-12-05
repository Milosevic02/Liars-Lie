package agent

import (
	"fmt"
	"os"
	"time"
)

// LogInfo logs information about the agent in the "logs" folder.
func (a *Agent) LogInfo(message string) {
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, 0755)
		if err != nil {
			fmt.Println("Error creating logs directory:", err)
			return
		}
	}

	filename := fmt.Sprintf("%s/agent_%d.log", logDir, a.ID)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error creating log file for agent", a.ID, ":", err)
		return
	}
	defer file.Close()

	logMessage := fmt.Sprintf("%s | Port: %s | Value: %d | Truthful: %v | %s\n",
		time.Now().Format("2006-01-02 15:04:05"), a.Port, a.Value, a.Truthful, message)

	if _, err := file.WriteString(logMessage); err != nil {
		fmt.Println("Error writing to log file:", err)
	}
}
