package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/armon/go-socks5"
)

const (
	appName    = "littlesocks"
	appVersion = "dev"
)

func main() {
	// Command line flags
	addr := flag.String("addr", "127.0.0.1:1080", "Address to bind the SOCKS5 server (host:port)")
	version := flag.Bool("version", false, "Print version information")
	flag.Parse()

	// Handle version flag
	if *version {
		fmt.Printf("%s version %s\n", appName, appVersion)
		os.Exit(0)
	}

	// Create a SOCKS5 server with minimal configuration
	conf := &socks5.Config{}
	server, err := socks5.New(conf)
	if err != nil {
		log.Fatalf("Failed to create SOCKS5 server: %v", err)
	}

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		log.Printf("Starting SOCKS5 proxy server on %s", *addr)
		if err := server.ListenAndServe("tcp", *addr); err != nil {
			errChan <- err
		}
	}()

	// Wait for shutdown signal or error
	select {
	case err := <-errChan:
		log.Fatalf("Server error: %v", err)
	case sig := <-sigChan:
		log.Printf("Received signal %v, shutting down gracefully...", sig)
	}
}

