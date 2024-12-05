package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

var agentProcesses = make(map[int]*exec.Cmd)

// Starts the game by launching agents as HTTP servers and writing their ports to agents.config
func StartGame(value, max, numAgents int, liarRatio float64) {
	configFile, _ := os.Create("agents.config")
	defer configFile.Close()

	for i := 0; i < numAgents; i++ {
		truthful := i < int(float64(numAgents)*(1-liarRatio))
		agentValue := value
		if !truthful {
			agentValue = rand.Intn(max) + 1
			for agentValue == value {
				agentValue = rand.Intn(max) + 1
			}
		}
		port := strconv.Itoa(8000 + i)

		if !IsPortAvailable(port) {
			fmt.Printf("Port %s is already in use, try a different port.\n", port)
			continue
		}

		cmd := exec.Command("go", "run", "agent_starter/main.go", strconv.Itoa(i), strconv.Itoa(agentValue), strconv.FormatBool(truthful), port)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		go func(agentID int) {
			err := cmd.Run()
			if err != nil {
				fmt.Printf("Error running agent %d: %v\n", agentID, err)
			}
		}(i)

		agentProcesses[i] = cmd

		fmt.Fprintf(configFile, "%s\n", port)
	}

	fmt.Println("ready")
}

// Kills an agent by sending a /kill request
func KillAgentByID(id int, printMessage bool) {
	_, exists := agentProcesses[id]
	if !exists {
		fmt.Printf("Agent with ID %d not found.\n", id)
		return
	}
	port := strconv.Itoa(8000 + id)
	url := fmt.Sprintf("http://localhost:%s/kill", port)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("Error killing agent %d: %v\n", id, err)
	} else {
		if printMessage {
			fmt.Printf("Agent with ID %d successfully killed.\n", id)
		}
	}
	delete(agentProcesses, id)
}

// Stops all running agents by sending a /kill request to each, then removes logs and config files
func StopGame(exitProgram bool) {
	for id := range agentProcesses {
		KillAgentByID(id, exitProgram)
	}
	agentProcesses = make(map[int]*exec.Cmd)

	if err := os.Remove("agents.config"); err != nil && !os.IsNotExist(err) {
		fmt.Printf("Error removing config file: %v\n", err)
	}
	if err := DeleteAgentLogs(); err != nil {
		fmt.Printf("Error deleting log files: %v\n", err)
	}
	if exitProgram {
		fmt.Println("Exiting program...")
		os.Exit(0)
	}
}

// extendNetwork adds new agents to the game, appending their ports to agents.config
func ExtendNetwork(value, max, numAgents int, liarRatio float64) {
	configFile, err := os.OpenFile("agents.config", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening agents.config:", err)
		return
	}
	defer configFile.Close()

	initialAgentCount := len(agentProcesses)
	for i := initialAgentCount; i < initialAgentCount+numAgents; i++ {
		truthful := (i - initialAgentCount) < int(float64(numAgents)*(1-liarRatio))
		agentValue := value
		if !truthful {
			agentValue = rand.Intn(max) + 1
			for agentValue == value {
				agentValue = rand.Intn(max) + 1
			}
		}
		port := strconv.Itoa(8000 + i)

		if !IsPortAvailable(port) {
			fmt.Printf("Port %s is already in use, try a different port.\n", port)
			continue
		}

		cmd := exec.Command("go", "run", "agent_starter/main.go", strconv.Itoa(i), strconv.Itoa(agentValue), strconv.FormatBool(truthful), port)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		go func(agentID int) {
			err := cmd.Run()
			if err != nil {
				fmt.Printf("Error running agent %d: %v\n", agentID, err)
			}
		}(i)
		agentProcesses[i] = cmd

		fmt.Fprintf(configFile, "%s\n", port)
	}

	fmt.Printf("%d new agents have been added to the network.\n", numAgents)
}
