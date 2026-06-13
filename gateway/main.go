package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	helper "github.com/MayurUbarhande0/Simple-cloud-storage/db"
	"github.com/MayurUbarhande0/Simple-cloud-storage/gateway/server"
	"github.com/MayurUbarhande0/Simple-cloud-storage/gateway/statemanager"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	if err := os.MkdirAll("/cloud", 0755); err != nil {
		log.Fatalf("Failed to create base storage directory: %v", err)
	}

	sqlDB, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}

	dbHelper := helper.NewDb(sqlDB)

	stateMgr, err := statemanager.NewManager("./state.json", dbHelper)
	if err != nil {
		log.Fatalf("Failed to initialize State Manager: %v", err)
	}

	server.StateMgr = stateMgr
	server.DbInstance = dbHelper

	mux := http.NewServeMux()
	mux.HandleFunc("/upload", server.Uploadfile)
	mux.HandleFunc("/download", server.Getfile)

	log.Printf("[Gateway] Storage server listening on port :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
