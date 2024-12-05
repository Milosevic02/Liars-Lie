package main

import (
	"fmt"
	"lierslie/agent"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: go run agent.go <ID> <Value> <Truthful> <Port>")
		return
	}

	id, _ := strconv.Atoi(os.Args[1])
	value, _ := strconv.Atoi(os.Args[2])
	truthful, _ := strconv.ParseBool(os.Args[3])
	port := os.Args[4]

	agent := &agent.Agent{ID: id, Value: value, Truthful: truthful, Port: port}
	agent.StartAgent()
}
