package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"time"
)

type ScanResult struct {
	hostname string
	port int
	isOpen bool
}

type ScanQueue struct {
	ports chan int
	results chan ScanResult
	hostname string
}

func scanner(queue ScanQueue) {
	// for each port in the supplied queue, attempt to connect
	// and return the result as a boolean
	for port := range queue.ports {
		address := fmt.Sprintf("%s:%d", queue.hostname, port)
		conn, err := net.Dial("tcp", address)

		if err != nil {
			queue.results <- ScanResult{ hostname: queue.hostname, port: port, isOpen: false }
			continue
		}
		conn.Close()
		queue.results <- ScanResult{ hostname: queue.hostname, port: port, isOpen: true }
	}
}

func main() {
	start := time.Now()

	fmt.Println("Spyglass Port Scanner v0.1")

	// flags
	bufferSize := flag.Int("buffer", 100, "The size of the buffer for the scanner to use")
	hostname := flag.String("hostname", "localhost", "The hostname to scan")
	startPort := flag.Int("start", 1, "The port to start scanning at")
	endPort := flag.Int("end", 1024, "The port to end scanning at")

	flag.Parse()

	if *startPort >= *endPort {
		err := fmt.Sprintf("Invalid start/end ports: %d - %d", *startPort, *endPort)
		panic(errors.New(err))
	}

	// worker channels
	ports := make(chan int, *bufferSize)
	results := make(chan ScanResult)

	// close the channels once scanning is complete
	defer close(ports)
	defer close(results)

	fmt.Println("Creating worker pool...")

	// create a pool of workers
	for i := 0; i < cap(ports); i++ {
		job := ScanQueue{ ports: ports, results: results, hostname: *hostname }
		go scanner(job)
	}

	fmt.Printf("Scanning %s...\n\n", *hostname)

	// send ports to be scanned
	go func() {
		for i := *startPort; i <= *endPort; i++ {
			ports <- i
		}
	}()

	// check port scanning results
	for i := *startPort; i <= *endPort; i++ {
		result := <- results

		if(result.isOpen) {
			fmt.Printf("Open port: %s:%d\n", result.hostname, result.port)
		}
	}

	elapsed := time.Since(start)
	fmt.Println("\nScan completed in", elapsed)
}