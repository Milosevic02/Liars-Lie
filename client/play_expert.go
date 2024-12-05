package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

func (c *Client) PlayExpert(numAgents int, liarRatio float64) int {
	valueCounts := make(map[int]int)

	// Randomly select agents to question directly
	selectedAgents := rand.Perm(len(c.AgentPorts))[:numAgents]

	for _, idx := range selectedAgents {
		port := c.AgentPorts[idx]

		// Prepare list of other agents' ports to send in request body
		otherPorts := []string{}
		for _, p := range c.AgentPorts {
			if p != port { // Exclude the current agent's port
				otherPorts = append(otherPorts, p)
			}
		}
		data, _ := json.Marshal(otherPorts) // Encode list of ports as JSON

		// Send HTTP POST request in expert mode with the list of other agent ports
		resp, err := http.Post("http://localhost:"+port+"/playexpert", "application/json", bytes.NewBuffer(data))
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		// Read the value sent by the agent
		scanner := bufio.NewScanner(resp.Body)
		if scanner.Scan() {
			val, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
			if err == nil {
				valueCounts[val]++
			} else {
				fmt.Println("Error reading value from agent on port", port, ":", err)
			}
		} else {
			fmt.Println("Error reading response from agent on port", port, ":", scanner.Err())
		}
	}

	// Determine the most common value among agent responses
	var networkValue, maxCount int
	for val, count := range valueCounts {
		if count > maxCount {
			maxCount = count
			networkValue = val
		}
	}
	return networkValue
}
