package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

func processLine(line string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()
	// Logic: Identify failed transactions for alerting
	if strings.Contains(line, "FAILURE") {
		results <- line
	}
}

func main() {
	start := time.Now()
	file, err := os.Open("transactions.log")
	if err != nil {
		fmt.Println("Error opening file. Run generator.go first!")
		return
	}
	defer file.Close()

	var wg sync.WaitGroup
	results := make(chan string, 100000)
	scanner := bufio.NewScanner(file)

	fmt.Println("Starting distributed log processing...")

	for scanner.Scan() {
		wg.Add(1)
		go processLine(scanner.Text(), &wg, results) // Concurrency in action
	}

	// Closing results channel once all lines are processed
	go func() {
		wg.Wait()
		close(results)
	}()

	failureCount := 0
	for range results {
		failureCount++
	}

	duration := time.Since(start)
	fmt.Println("--------------------------------------------------")
	fmt.Printf("Processed 100,000 logs in: %v\n", duration)
	fmt.Printf("Critical Failures Found: %d\n", failureCount)
	fmt.Println("--------------------------------------------------")
}
