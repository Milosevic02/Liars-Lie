package client

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Client struct {
	AgentPorts []string
}

func (c *Client) ReadConfig() {
	file, err := os.Open("agents.config")
	if err != nil {
		fmt.Println("Error reading agents.config:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		port := strings.TrimSpace(scanner.Text())
		c.AgentPorts = append(c.AgentPorts, port)
	}
}
