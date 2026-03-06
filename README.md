# LOG-AGGREGATOR

# High-Performance Transaction Log Aggregator (Go)
 ## Overview
This project is a high-concurrency engine designed to simulate and process massive volumes of payment transaction logs.
As an Electronics and Communication Engineering (ECE) student, I developed this to demonstrate how hardware-level logic and concurrent programming can ensure the reliability and performance of national-scale financial systems.

## Key Technical Achievements
Unmatched Throughput: Capable of processing 100,000 transaction records in ~59ms (approx. 1.7 Million logs/sec).

Concurrency mastery:
Leveraged Go’s Goroutines and sync.WaitGroup to achieve parallel execution, drastically reducing latency compared to traditional single-threaded software.

Scalable Architecture: Designed with a decoupled structure (Generator and Processor) to reflect real-world Distributed Systems patterns.

## Features
Real-time Simulation: Generates 100k simulated UPI-style transaction logs with randomized statuses (SUCCESS, FAILURE, PENDING) and latencies.

Critical Health Monitoring: Automatically identifies and aggregates "FAILURE" states to simulate real-time alerting for system maintenance.

Memory Optimization: Engineered to handle large datasets with minimal memory overhead using Go’s efficient bufio scanner.

## Tech Stack
Language: Go (Golang)

Core Concepts: Concurrency, Parallelism, Low-Latency File I/O, System Reliability

## How to Run
Clone the repository:

Bash
git clone https://github.com/sreenithadamayapally83/Log-Aggregator.git
Generate Transaction Logs:

Bash
go run ./generator/generator.go
Analyze & Process Logs:

Bash
go run ./processor/processor.go
