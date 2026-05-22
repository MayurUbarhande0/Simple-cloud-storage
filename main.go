package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	// Replace this import path with the exact module name defined inside your go.mod file
	"github.com/MayurUbarhande0/Simple-cloud-storage/protocol"
)

func main() {
	// 1. Structural Sanity Check: Ensure required cryptographic secrets are active
	if os.Getenv("ENCRYPTION_KEY") == "" || os.Getenv("AUTH_KEY") == "" {
		log.Fatal("Critical Boot Failure: ENCRYPTION_KEY and AUTH_KEY environment variables must be configured.")
	}

	// 2. Resolve target socket network address
	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":8081" // Default fallback port if none specified
	}

	// 3. Initialize your custom protocol server engine instance
	server := protocol.NewServer(listenAddr)

	// 4. Spin up your AcceptLoop in a background thread to prevent blocking our signal hook
	go func() {
		log.Printf("🚀 Storage engine node initializing network socket listener on %s...", listenAddr)
		if err := server.Start(); err != nil {
			log.Fatalf("Fatal network runtime engine error: %v", err)
		}
	}()

	// 5. Graceful Intercept: Block main thread execution until system kill signals hit the application
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)

	<-shutdownSignal // Halts right here until Ctrl+C (SIGINT) or SIGTERM is registered

	log.Println("🛑 Termination signal recognized. Completing active background transfer sockets and spinning down gracefully.")
}
