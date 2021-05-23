package main

import (
	"fmt"
	"net"
)

type ScanResult struct {
	hostname string
	port int
	isOpen bool
}

type ScanJob struct {
	ports chan int
	results chan ScanResult
	hostname string
}

func scanner(job ScanJob) {
	// for each port in the supplied job, attempt to connect
	// and return the result as a boolean
	for port := range job.ports {
		address := fmt.Sprintf("%s:%d", job.hostname, port)
		conn, err := net.Dial("tcp", address)

		if err != nil {
			job.results <- ScanResult{ hostname: job.hostname, port: port, isOpen: false }
			continue
		}
		conn.Close()
		job.results <- ScanResult{ hostname: job.hostname, port: port, isOpen: true }
	}
}

func main() {
	fmt.Println("Spyglass Port Scanner v0.1")

	ports := make(chan int, 100)
	results := make(chan ScanResult)
	hostname := "localhost"

	// close the channels once scanning is complete
	defer close(ports)
	defer close(results)

	fmt.Println("Creating worker pool...")

	// create a pool of workers
	for i := 0; i < cap(ports); i++ {
		job := ScanJob{ ports: ports, results: results, hostname: hostname }
		go scanner(job)
	}

	fmt.Printf("Scanning %s...\n\n", hostname)

	// send ports to be scanned
	go func() {
		for i := 1; i <= 10000; i++ {
			ports <- i
		}
	}()

	// check port scanning results
	for i := 1; i <= 10000; i++ {
		result := <- results

		if(result.isOpen) {
			fmt.Printf("Open port: %s:%d\n", result.hostname, result.port)
		}
	}
}