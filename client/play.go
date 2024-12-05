package client

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (c *Client) Play() int {
	valueCounts := make(map[int]int)

	for _, port := range c.AgentPorts {

		// Send HTTP GET request to the agent's /value endpoint
		resp, err := http.Get("http://localhost:" + port + "/value")
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

	// Find the most common value among all agents
	var networkValue, maxCount int
	for val, count := range valueCounts {
		if count > maxCount {
			maxCount = count
			networkValue = val
		}
	}
	return networkValue
}
