package main

import (
	"log"
	"net/http"
	"os"

	"github.com/MayurUbarhande0/Simple-cloud-storage/gateway/routes"
	"github.com/MayurUbarhande0/Simple-cloud-storage/gateway/server"
	"github.com/MayurUbarhande0/Simple-cloud-storage/gateway/statemanager"
)

var StateMgr *statemanager.Manager

func main() {
	// Initialize the manager instance and link it to your server package variable
	mgr, err := statemanager.NewManger("state.json")
	if err != nil {
		log.Fatalf("Failed to balance state records: %v", err)
	}
	server.StateMgr = mgr // Bind it to your server package variable!
	// 1. Sanity checks for Gateway requirements
	if os.Getenv("ENCRYPTION_KEY") == "" || os.Getenv("AUTH_KEY") == "" {
		log.Fatal("Gateway Error: Cryptographic environment variables are missing.")
	}
	if os.Getenv("STORAGE_IP") == "" {
		log.Fatal("Gateway Error: STORAGE_IP destination address is not configured.")
	}

	// 2. Resolve local HTTP server port
	httpPort := os.Getenv("GATEWAY_HTTP_PORT")
	if httpPort == "" {
		httpPort = ":8080" // Fallback default
	}

	// 3. Initialize your clean routing tree
	router := routes.NewRouter()

	log.Printf("🚀 Cloud Storage Gateway running on HTTP port %s...", httpPort)

	// 4. Fire up the HTTP Web Server
	if err := http.ListenAndServe(httpPort, router); err != nil {
		log.Fatalf("Gateway web server crashed: %v", err)
	}
}
