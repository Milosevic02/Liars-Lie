package controllers

import (
	"fmt"
	"net"
	"os"
)

// DeleteAgentLogs deletes the entire "logs" directory along with its contents.
func DeleteAgentLogs() error {
	logDir := "logs"

	// Pokušaj da obrišeš ceo "logs" direktorijum
	err := os.RemoveAll(logDir)
	if err != nil {
		fmt.Printf("Error deleting logs directory: %v\n", err)
		return err
	}

	return nil
}

// Checks if a port is available for binding
func IsPortAvailable(port string) bool {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

// CheckConfigExists verifies if the agent configuration file ("agents.config") exists
func CheckConfigExists() bool {
	_, err := os.Stat("agents.config")
	return !os.IsNotExist(err)
}
