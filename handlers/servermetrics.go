package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testcaseautomation/cmd"
)

type ProcessStats struct {
	Name        string
	PID         int
	CPU         float64
	Memory      float64
	ThreadCount int
	Status      string
	UsedPct     float64
}

type ResourceUsage struct {
	Total     float64
	Used      float64
	Available float64
	UsedPct   float64
}

// Sample data for demonstration
var processStats = []ProcessStats{
	{Name: "Process1", PID: 1234, CPU: 10.5, Memory: 2048, ThreadCount: 5, Status: "Running", UsedPct: 50.0},
	{Name: "Process2", PID: 5678, CPU: 20.0, Memory: 4096, ThreadCount: 10, Status: "Sleeping", UsedPct: 75.0},
}

var resourceUsage = ResourceUsage{Total: 100.0, Used: 75.0, Available: 25.0, UsedPct: 75.0}

func GetProcessStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(processStats)
}

func GetResourceUsage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resourceUsage)
}

func ExecuteBackgroundCommandOnServer(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	command := r.URL.Query().Get("command")

	if host == "" || username == "" || password == "" || command == "" {
		http.Error(w, "Missing required query parameters", http.StatusBadRequest)
		return
	}

	exitCode, err := cmd.ExecuteBackgroundCommandOnServer(host, username, password, command)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing command: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"exitCode": exitCode,
		"message":  "Command executed successfully",
	})
}

func CollectMetrics(w http.ResponseWriter, r *http.Request) {
	// Assuming servers, metricsChan, and errorsChan are defined elsewhere in the code
	servers := []string{"server1", "server2"} // Example servers
	metrics, err := cmd.CollectServerMetrics(servers, "username", "password")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error collecting metrics: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
