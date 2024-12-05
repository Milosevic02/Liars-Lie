package agent

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func (a *Agent) handleGetValue(w http.ResponseWriter, r *http.Request) {
	a.LogInfo("Received request for value")
	responseValue := strconv.Itoa(a.Value)
	w.Write([]byte(responseValue))
	a.LogInfo(fmt.Sprintf("Sent value: %d", a.Value))
}

func (a *Agent) handleKill(w http.ResponseWriter, r *http.Request) {
	a.LogInfo("Received kill request")
	w.Write([]byte("Agent shutting down"))
	go a.shutdownServer() // Initiates shutdown in a separate goroutine
}

func (a *Agent) handlePlayExpert(w http.ResponseWriter, r *http.Request) {
	a.LogInfo("Received expert mode request")

	if !a.Truthful {
		a.LogInfo("Agent is not truthful, returning its own value")
		w.Write([]byte(strconv.Itoa(a.Value)))
		return
	}

	var otherAgentsPorts []string
	if err := json.NewDecoder(r.Body).Decode(&otherAgentsPorts); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Count occurrences of values from other agents
	valueCounts := make(map[int]int)
	for _, otherPort := range otherAgentsPorts {
		resp, err := http.Get("http://localhost:" + otherPort + "/value")
		if err != nil {
			continue // Skip if agent is unavailable
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		if scanner.Scan() {
			val, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
			if err == nil {
				valueCounts[val]++
			}
		}
	}

	// Determine the most frequent value (network value)
	var networkValue, maxCount int
	for val, count := range valueCounts {
		if count > maxCount {
			maxCount = count
			networkValue = val
		}
	}

	a.LogInfo(fmt.Sprintf("Responding with network value: %d", networkValue))
	w.Write([]byte(strconv.Itoa(networkValue)))
}

func (a *Agent) shutdownServer() {
	if err := a.server.Close(); err != nil {
		fmt.Println("Error shutting down server:", err)
		a.LogInfo("Failed to shut down server")
	} else {
		fmt.Printf("Agent %d on port %s shut down\n", a.ID, a.Port)
		a.LogInfo("Server shut down successfully")
	}
}

// handleTerminationSignals listens for OS signals to gracefully shut down the server
func (a *Agent) handleTerminationSignals() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	a.shutdownServer()
}
