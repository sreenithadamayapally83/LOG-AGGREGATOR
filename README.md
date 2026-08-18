# Distributed Log Ingestion & Analytics Engine

A high-concurrency log processing pipeline simulating national-scale financial transaction systems. Built with Go and Python, this engine demonstrates low-latency file I/O, goroutine worker pools, real-time REST API monitoring, and cross-language database persistence.

## Architecture & Tech Stack

* **Producer-Consumer Engine (Go):** Utilizes `sync.WaitGroup`, channels, and mutex locks to safely process 1.7M+ logs/second using a decoupled worker pool architecture.
* **REST API (Go):** Exposes a real-time `net/http` endpoint (`/metrics`) for live monitoring of p50/p95/p99 latencies and failure counts.
* **Persistence Layer (PostgreSQL):** Aggregated metrics are permanently logged to a relational database using Go's `database/sql` package with built-in fault tolerance.
* **BI Consumer (Python):** A cross-language consumer utilizing `pandas` and `psycopg2` to query the SQL store and auto-generate executive `.xlsx` reports.

## Core Engineering Achievements
* **Unmatched Throughput:** Processes 100,000 transaction records in ~60ms (approx. 1.7 Million logs/sec) by heavily mitigating garbage collection overhead and leveraging `bufio`.
* **Concurrency Mastery:** Replaced anti-pattern 1:1 thread spawning with a strict 10-worker concurrent pool, preventing memory leaks on massive datasets.
* **Fault Tolerance:** Database connections gracefully degrade to terminal logging if the PostgreSQL instance goes offline, preventing complete system crashes.

## How to Run

**1. Generate Simulated Transaction Logs**
`go run ./generator/generator.go`

**2. Spin Up the Ingestion Engine & REST API**
`go get github.com/lib/pq`
`go run ./processor/processor.go`
*Engine runs continuously. View live metrics at `http://localhost:8080/metrics`.*

**3. Run the Cross-Language BI Consumer (Requires PostgreSQL)**
`pip install pandas psycopg2-binary openpyxl`
`python report_generator.py`
*Outputs timestamped executive reports to the `/reports` directory.*