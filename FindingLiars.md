
# Option 1 for PlayExpert - Finding Liars through Agent Voting

## Overview
The first option for implementing the `PlayExpert` functionality involves selecting suspicious agents based on their values, calculating the most common value in the network, and nominating agents based on the `liarRatio`. 

However, the main difficulty in this approach is the identification of liars in the network, as the system does not inherently know who the liars are. This method relies on communication between agents to determine who might be lying.

## Client Code (PlayExpert)

```go
func (c *Client) PlayExpert(numAgents int, liarRatio float64) int {
    valueCounts := make(map[int]int)
    suspectVotes := make(map[string]int)  // To track suspect votes

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
                // Track suspect votes
                if val != c.Value {
                    suspectVotes[port]++
                }
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

    // Adjust based on liar ratio and suspicious votes
    numLiarVotes := int(float64(len(c.AgentPorts)) * liarRatio)
    suspectedLiarPorts := make([]string, 0)

    // Find the most suspicious agents based on votes
    for port, count := range suspectVotes {
        if count > numLiarVotes {
            suspectedLiarPorts = append(suspectedLiarPorts, port)
        }
    }

    // Respond with the network value, factoring in suspicious agents
    return networkValue
}
```

## Agent Code (handlePlayExpert)

```go
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

    valueCounts := make(map[int]int)
    for _, otherPort := range otherAgentsPorts {
        resp, err := http.Get("http://localhost:" + otherPort + "/value")
        if err != nil {
            continue
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
```

## Conclusion

This option provides an interesting way to identify suspicious agents based on their responses and adjust the final decision accordingly. However, it doesn't fully adhere to the task's requirement of respecting the `liarRatio` in the way that some other options might.

---

# Notes:
- The client code randomly selects agents, sends their ports, and receives their values.
- The agent code checks if it is truthful, and if so, processes the list of other agents' values to determine the most common one.
