package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type LogEntry struct {
	Timestamp    time.Time `json:"timestamp"`
	Method       string    `json:"method"`
	Path         string    `json:"path"`
	StatusCode   int       `json:"status_code"`
	ResponseTime int64     `json:"response_time_ms"`
	UserAgent    string    `json:"user_agent,omitempty"`
	RemoteAddr   string    `json:"remote_addr,omitempty"`
	RequestID    string    `json:"request_id,omitempty"`
	Error        string    `json:"error,omitempty"`
	RequestBody  string    `json:"request_body,omitempty"`
	ResponseSize int       `json:"response_size,omitempty"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run view_logs.go [command] [options]")
		fmt.Println("Commands:")
		fmt.Println("  latest [n]     - Show latest n entries (default: 10)")
		fmt.Println("  errors         - Show only error entries")
		fmt.Println("  stats          - Show request statistics")
		fmt.Println("  filter [path]  - Filter by endpoint path")
		return
	}

	command := os.Args[1]
	
	// Find the latest log file
	logFile, err := findLatestLogFile()
	if err != nil {
		log.Fatal("Error finding log file:", err)
	}

	entries, err := readLogEntries(logFile)
	if err != nil {
		log.Fatal("Error reading log entries:", err)
	}

	switch command {
	case "latest":
		n := 10
		if len(os.Args) > 2 {
			fmt.Sscanf(os.Args[2], "%d", &n)
		}
		showLatest(entries, n)
	case "errors":
		showErrors(entries)
	case "stats":
		showStats(entries)
	case "filter":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a path to filter by")
			return
		}
		filterByPath(entries, os.Args[2])
	default:
		fmt.Println("Unknown command:", command)
	}
}

func findLatestLogFile() (string, error) {
	logsDir := "logs"
	files, err := filepath.Glob(filepath.Join(logsDir, "api_*.log"))
	if err != nil {
		return "", err
	}
	
	if len(files) == 0 {
		return "", fmt.Errorf("no log files found in %s", logsDir)
	}
	
	sort.Strings(files)
	return files[len(files)-1], nil
}

func readLogEntries(filename string) ([]LogEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		
		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue // Skip malformed entries
		}
		entries = append(entries, entry)
	}
	
	return entries, scanner.Err()
}

func showLatest(entries []LogEntry, n int) {
	fmt.Printf("Latest %d log entries:\n", n)
	fmt.Println(strings.Repeat("-", 80))
	
	start := len(entries) - n
	if start < 0 {
		start = 0
	}
	
	for i := start; i < len(entries); i++ {
		entry := entries[i]
		statusColor := getStatusColor(entry.StatusCode)
		fmt.Printf("%s [%s] %s %s%d%s (%dms)\n",
			entry.Timestamp.Format("2006-01-02 15:04:05"),
			entry.Method,
			entry.Path,
			statusColor,
			entry.StatusCode,
			"\033[0m",
			entry.ResponseTime,
		)
		if entry.Error != "" {
			fmt.Printf("  Error: %s\n", entry.Error)
		}
	}
}

func showErrors(entries []LogEntry) {
	fmt.Println("Error entries:")
	fmt.Println(strings.Repeat("-", 80))
	
	for _, entry := range entries {
		if entry.StatusCode >= 400 || entry.Error != "" {
			fmt.Printf("%s [%s] %s - %d\n",
				entry.Timestamp.Format("2006-01-02 15:04:05"),
				entry.Method,
				entry.Path,
				entry.StatusCode,
			)
			if entry.Error != "" {
				fmt.Printf("  Error: %s\n", entry.Error)
			}
		}
	}
}

func showStats(entries []LogEntry) {
	if len(entries) == 0 {
		fmt.Println("No log entries found")
		return
	}

	// Count by status code
	statusCounts := make(map[int]int)
	methodCounts := make(map[string]int)
	pathCounts := make(map[string]int)
	var totalResponseTime int64
	
	for _, entry := range entries {
		statusCounts[entry.StatusCode]++
		methodCounts[entry.Method]++
		pathCounts[entry.Path]++
		totalResponseTime += entry.ResponseTime
	}
	
	fmt.Printf("API Request Statistics (%d total requests)\n", len(entries))
	fmt.Println(strings.Repeat("=", 50))
	
	fmt.Println("\nStatus Code Distribution:")
	for status, count := range statusCounts {
		percentage := float64(count) / float64(len(entries)) * 100
		fmt.Printf("  %d: %d (%.1f%%)\n", status, count, percentage)
	}
	
	fmt.Println("\nMethod Distribution:")
	for method, count := range methodCounts {
		percentage := float64(count) / float64(len(entries)) * 100
		fmt.Printf("  %s: %d (%.1f%%)\n", method, count, percentage)
	}
	
	fmt.Println("\nTop Endpoints:")
	type pathCount struct {
		path  string
		count int
	}
	var paths []pathCount
	for path, count := range pathCounts {
		paths = append(paths, pathCount{path, count})
	}
	sort.Slice(paths, func(i, j int) bool {
		return paths[i].count > paths[j].count
	})
	
	for i, pc := range paths {
		if i >= 10 { // Show top 10
			break
		}
		percentage := float64(pc.count) / float64(len(entries)) * 100
		fmt.Printf("  %s: %d (%.1f%%)\n", pc.path, pc.count, percentage)
	}
	
	avgResponseTime := totalResponseTime / int64(len(entries))
	fmt.Printf("\nAverage Response Time: %dms\n", avgResponseTime)
}

func filterByPath(entries []LogEntry, pathFilter string) {
	fmt.Printf("Entries for path containing '%s':\n", pathFilter)
	fmt.Println(strings.Repeat("-", 80))
	
	for _, entry := range entries {
		if strings.Contains(entry.Path, pathFilter) {
			statusColor := getStatusColor(entry.StatusCode)
			fmt.Printf("%s [%s] %s %s%d%s (%dms)\n",
				entry.Timestamp.Format("2006-01-02 15:04:05"),
				entry.Method,
				entry.Path,
				statusColor,
				entry.StatusCode,
				"\033[0m",
				entry.ResponseTime,
			)
			if entry.Error != "" {
				fmt.Printf("  Error: %s\n", entry.Error)
			}
		}
	}
}

func getStatusColor(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "\033[32m" // Green
	case statusCode >= 300 && statusCode < 400:
		return "\033[33m" // Yellow
	case statusCode >= 400 && statusCode < 500:
		return "\033[31m" // Red
	case statusCode >= 500:
		return "\033[35m" // Magenta
	default:
		return "\033[37m" // White
	}
}