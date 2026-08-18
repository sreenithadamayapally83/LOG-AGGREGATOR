package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type MetricsStore struct {
	mu           sync.Mutex
	TotalLogs    int     `json:"total_logs"`
	FailureCount int     `json:"failure_count"`
	P50Latency   float64 `json:"p50_latency_ms"`
	P95Latency   float64 `json:"p95_latency_ms"`
	P99Latency   float64 `json:"p99_latency_ms"`
}

var globalMetrics MetricsStore
var allLatencies []int

func worker(jobs <-chan string, latenciesChan chan<- int, failuresChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for line := range jobs {
		parts := strings.Split(line, " | ")
		if len(parts) >= 3 {
			status := strings.TrimPrefix(parts[1], "STATUS:")
			latencyStr := strings.TrimPrefix(parts[2], "LATENCY:")
			latencyStr = strings.TrimSuffix(latencyStr, "ms")

			latency, err := strconv.Atoi(latencyStr)
			if err == nil {
				latenciesChan <- latency
			}

			if status == "REJECTED" || status == "FAILURE" {
				failuresChan <- 1
			}
		}
	}
}

func startRESTAPI() {
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		globalMetrics.mu.Lock()
		defer globalMetrics.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(globalMetrics)
	})
	fmt.Println("--> REST API live at http://localhost:8080/metrics")
	http.ListenAndServe(":8080", nil)
}

func persistToPostgres() {
	// IMPORTANT: Update this connection string with your local Postgres credentials
	connStr := "user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("--> [WARN] Could not connect to PostgreSQL: %v\n", err)
		return
	}
	defer db.Close()

	// Ensure our table exists
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS ingestion_metrics (
		id SERIAL PRIMARY KEY,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		total_logs INT,
		failure_count INT,
		p50_latency FLOAT,
		p95_latency FLOAT,
		p99_latency FLOAT
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		fmt.Printf("--> [WARN] Could not create table: %v\n", err)
		return
	}

	// Insert the aggregated metrics
	insertQuery := `
	INSERT INTO ingestion_metrics (total_logs, failure_count, p50_latency, p95_latency, p99_latency)
	VALUES ($1, $2, $3, $4, $5)`

	_, err = db.Exec(insertQuery, globalMetrics.TotalLogs, globalMetrics.FailureCount,
		globalMetrics.P50Latency, globalMetrics.P95Latency, globalMetrics.P99Latency)

	if err != nil {
		fmt.Printf("--> [WARN] Failed to insert metrics: %v\n", err)
	} else {
		fmt.Println("--> [SUCCESS] Metrics successfully persisted to PostgreSQL.")
	}
}

func main() {
	go startRESTAPI()

	start := time.Now()
	file, err := os.Open("transactions.log")
	if err != nil {
		fmt.Println("Error opening file. Run generator.go first!")
		return
	}
	defer file.Close()

	jobs := make(chan string, 100000)
	latenciesChan := make(chan int, 100000)
	failuresChan := make(chan int, 100000)
	var wg sync.WaitGroup

	numWorkers := 10
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(jobs, latenciesChan, failuresChan, &wg)
	}

	go func() {
		for lat := range latenciesChan {
			allLatencies = append(allLatencies, lat)
		}
	}()

	go func() {
		for range failuresChan {
			globalMetrics.mu.Lock()
			globalMetrics.FailureCount++
			globalMetrics.mu.Unlock()
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		jobs <- scanner.Text()
		globalMetrics.mu.Lock()
		globalMetrics.TotalLogs++
		globalMetrics.mu.Unlock()
	}
	close(jobs)

	wg.Wait()
	close(latenciesChan)
	close(failuresChan)

	sort.Ints(allLatencies)
	n := len(allLatencies)
	if n > 0 {
		globalMetrics.mu.Lock()
		globalMetrics.P50Latency = float64(allLatencies[int(float64(n)*0.50)])
		globalMetrics.P95Latency = float64(allLatencies[int(float64(n)*0.95)])
		globalMetrics.P99Latency = float64(allLatencies[int(float64(n)*0.99)])
		globalMetrics.mu.Unlock()
	}

	duration := time.Since(start)

	fmt.Println("--------------------------------------------------")
	fmt.Printf("Processed %d logs in: %v\n", globalMetrics.TotalLogs, duration)
	fmt.Printf("Critical Failures Found: %d\n", globalMetrics.FailureCount)
	fmt.Printf("Latency Profile -> p50: %.0fms | p95: %.0fms | p99: %.0fms\n",
		globalMetrics.P50Latency, globalMetrics.P95Latency, globalMetrics.P99Latency)
	fmt.Println("--------------------------------------------------")

	// Persist to Database before blocking
	persistToPostgres()

	fmt.Println("Engine running. Press Ctrl+C to terminate.")
	select {}
}
