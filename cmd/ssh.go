package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
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

type ServerMetrics struct {
	Hostname     string
	CPU          ResourceUsage
	Memory       ResourceUsage
	Disk         ResourceUsage
	NetworkIO    ResourceUsage
	TimeStamp    time.Time
	ProcessCount int
	TopProcesses []ProcessStats
}

func CollectServerMetrics(servers []string, username, password string) ([]ServerMetrics, error) {
	var wg sync.WaitGroup
	metrics := make([]ServerMetrics, len(servers))
	metricsChan := make(chan ServerMetrics, len(servers))
	errorsChan := make(chan error, len(servers))

	for _, server := range servers {
		wg.Add(1)
		go func(hostname string) {
			defer wg.Done()

			config := &ssh.ClientConfig{
				User: username,
				Auth: []ssh.AuthMethod{
					ssh.Password(password),
				},
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
				Timeout:         10 * time.Second,
			}

			client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", hostname), config)
			if err != nil {
				errorsChan <- fmt.Errorf("failed to connect to %s: %v", hostname, err)
				return
			}
			defer client.Close()

			session, err := client.NewSession()
			if err != nil {
				errorsChan <- fmt.Errorf("failed to create session on %s: %v", hostname, err)
				return
			}
			defer session.Close()

			cmd := `
        # CPU Usage
        cpu_stats=$(top -bn1 | grep "Cpu(s)" | awk '{printf "{\\"total\\":100,\\"used\\":%s,\\"available\\":%s,\\"usedpct\\":%s}", $2, $8, $2}')
        
        # Memory Usage
        mem_stats=$(free -m | grep Mem | awk '{printf "{\\"total\\":%d,\\"used\\":%d,\\"available\\":%d,\\"usedpct\\":%0.2f}", $2, $3, $4, $3/$2 * 100}')
        
        # Disk Usage
        disk_stats=$(df -m / | awk 'NR==2 {printf "{\\"total\\":%d,\\"used\\":%d,\\"available\\":%d,\\"usedpct\\":%d}", $2, $3, $4, $5}')
        
        # Network IO (last minute average)
        network_stats=$(sar -n DEV 1 1 | grep eth0 | tail -1 | awk '{printf "{\\"total\\":%f,\\"used\\":%f,\\"available\\":%f,\\"usedpct\\":%f}", $5+$6, $5, $6, ($5+$6)/100}')
        
        # Process metrics (top 10 processes)
        processes=$(ps aux --sort=-%cpu | head -11 | tail -10 | awk '
            {
                printf "{\\"name\\":\\"%s\\",\\"pid\\":%d,\\"cpu\\":%f,\\"memory\\":%f,\\"thread_count\\":%d,\\"status\\":\\"%s\\",\\"used_pct\\":%f},"
            }' \
            'BEGIN {ORS=""}')

        echo "{
            \"cpu\":$cpu_stats,
            \"memory\":$mem_stats,
            \"disk\":$disk_stats,
            \"network\":$network_stats,
            \"process_count\":$(ps aux | wc -l),
            \"processes\":[$processes]
        }"
    `

			output, err := session.Output(cmd)
			if err != nil {
				errorsChan <- fmt.Errorf("failed to execute command on %s: %v", hostname, err)
				return
			}

			var data struct {
				CPU          ResourceUsage  `json:"cpu"`
				Memory       ResourceUsage  `json:"memory"`
				Disk         ResourceUsage  `json:"disk"`
				Network      ResourceUsage  `json:"network"`
				ProcessCount int            `json:"process_count"`
				Processes    []ProcessStats `json:"processes"`
			}

			if err := json.Unmarshal(output, &data); err != nil {
				errorsChan <- fmt.Errorf("failed to parse metrics from %s: %v", hostname, err)
				return
			}

			metricsChan <- ServerMetrics{
				Hostname:     hostname,
				CPU:          data.CPU,
				Memory:       data.Memory,
				Disk:         data.Disk,
				NetworkIO:    data.Network,
				ProcessCount: data.ProcessCount,
				TopProcesses: data.Processes,
				TimeStamp:    time.Now(),
			}
		}(server)
	}

	go func() {
		wg.Wait()
		close(metricsChan)
		close(errorsChan)
	}()

	// Collect results
	var errors []error
	for i := 0; i < len(servers); i++ {
		select {
		case metric := <-metricsChan:
			metrics[i] = metric
		case err := <-errorsChan:
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return metrics, fmt.Errorf("errors collecting metrics: %v", errors)
	}

	return metrics, nil
}

func ExecuteBackgroundCommandOnServer(host, username, password, command string) (int, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", host), config)
	if err != nil {
		return -1, err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return -1, err
	}
	defer session.Close()

	backgroundCmd := fmt.Sprintf("nohup %s > /dev/null 2>&1 &", command)
	err = session.Run(backgroundCmd)
	if err != nil {
		if exitError, ok := err.(*ssh.ExitError); ok {
			return exitError.ExitStatus(), nil
		}
		return -1, err
	}

	return 0, nil
}

func KillProcessOnServer(host, username, password, processName string) error {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", host), config)
	if err != nil {
		return fmt.Errorf("failed to dial: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	// Find and kill the process
	killCmd := fmt.Sprintf("pid=$(pgrep -f '%s') && [ ! -z \"$pid\" ] && kill -9 $pid", processName)
	output, err := session.CombinedOutput(killCmd)
	if err != nil && !strings.Contains(string(output), "not found") {
		return fmt.Errorf("failed to kill process: %v, output: %s", err, output)
	}

	return nil
}
