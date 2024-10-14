package main

import (
	_ "embed" // Required for embedding files
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

//go:embed help.txt
var helpText string // Variable to hold the contents of the embedded help.txt file

// Struct to hold host and port details
type HostPort struct {
	Host   string
	Port   string
	Status string
}

// Function to check if a given host and port are reachable
func checkHostPort(hostPort HostPort) bool {
	address := hostPort.Host + ":" + hostPort.Port
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// Function to check multiple hosts and ports until timeout is reached or all are available
func waitForMultipleHosts(hostPorts []HostPort, timeout, retryInterval int) bool {
	start := time.Now()
	ready := make(map[HostPort]bool)

	// Initialize all hosts as not ready
	for _, hp := range hostPorts {
		ready[hp] = false
	}

	for {
		allReady := true

		// Iterate over all host/port pairs and check their availability
		for i, hp := range hostPorts {
			if ready[hp] {
				continue
			}

			log.WithFields(log.Fields{
				"host": hp.Host,
				"port": hp.Port,
			}).Debug("Checking availability")

			if checkHostPort(hp) {
				log.WithFields(log.Fields{
					"host": hp.Host,
					"port": hp.Port,
				}).Info("Host is available")
				hostPorts[i].Status = "available"
				ready[hp] = true
			} else {
				allReady = false
				hostPorts[i].Status = "unavailable"
				log.WithFields(log.Fields{
					"host": hp.Host,
					"port": hp.Port,
				}).Debug("Host is not available yet")
			}
		}

		// Break the loop if all hosts are ready
		if allReady {
			log.Info("All hosts are available!")
			return true
		}

		// Check if the timeout is reached
		if int(time.Since(start).Seconds()) >= timeout {
			log.Warn("Timeout reached. Not all hosts became available.")
			return false
		}

		time.Sleep(time.Duration(retryInterval) * time.Second)
	}
}

// Function to parse the input host:port pairs and return a slice of HostPort structs
func parseHosts(input string) []HostPort {
	var hostPorts []HostPort
	pairs := strings.Split(input, ",")

	for _, pair := range pairs {
		hp := strings.Split(pair, ":")
		if len(hp) != 2 {
			log.WithFields(log.Fields{
				"input": pair,
			}).Error("Invalid format for host:port")
			os.Exit(1)
		}
		hostPorts = append(hostPorts, HostPort{Host: hp[0], Port: hp[1], Status: "unknown"})
	}

	return hostPorts
}

func main() {
	// Define flags
	timeout := flag.Int("timeout", 15, "Timeout in seconds")
	retryInterval := flag.Int("retry-interval", 1, "Retry interval in seconds")
	quiet := flag.Bool("quiet", false, "Enable quiet mode")
	debug := flag.Bool("debug", false, "Enable debug mode")
	format := flag.String("format", "text", "Set log format (options: 'text' or 'json')")
	help := flag.Bool("help", false, "Show help")
	flag.Parse()

	// Show help if requested
	if *help {
		fmt.Println(helpText) // Print the embedded help text
		os.Exit(0)
	}

	// Set log level based on quiet or debug mode
	if *quiet {
		log.SetLevel(log.WarnLevel)
	} else if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// Set log format
	if *format == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{})
	}

	// Determine the source of host:port pairs (either from positional args or the WAIT environment variable)
	var hosts string
	if len(flag.Args()) > 0 {
		hosts = flag.Args()[0] // Use positional argument if provided
	} else {
		hosts = os.Getenv("WAIT") // Otherwise, fallback to the WAIT environment variable
	}

	if hosts == "" {
		fmt.Println(helpText) // Print help if no arguments or WAIT variable is provided
		os.Exit(1)
	}

	// Parse the host:port pairs from the chosen source
	hostPorts := parseHosts(hosts)

	// Wait for all hosts and handle the exit status based on availability
	if !waitForMultipleHosts(hostPorts, *timeout, *retryInterval) {
		os.Exit(1) // Exit with status 1 if some hosts are still unavailable
	}
}
