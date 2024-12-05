package agent

import (
	"fmt"
	"net/http"
)

type Agent struct {
	ID       int
	Value    int
	Truthful bool
	Port     string
	server   *http.Server // Holds the HTTP server instance for this agent
}

// StartAgent initializes the HTTP server for the agent and registers endpoints
func (a *Agent) StartAgent() {
	mux := http.NewServeMux()
	mux.HandleFunc("/value", a.handleGetValue)
	mux.HandleFunc("/kill", a.handleKill)
	mux.HandleFunc("/playexpert", a.handlePlayExpert)

	a.server = &http.Server{
		Addr:    ":" + a.Port,
		Handler: mux,
	}

	a.LogInfo("Agent started successfully as HTTP server")
	go a.handleTerminationSignals() // Graceful shutdown handler

	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("Error starting HTTP server:", err)
		a.LogInfo("Failed to start HTTP server")
	}
}
