package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	// Replace this import path with the exact module name inside your go.mod file
	"github.com/MayurUbarhande0/Simple-cloud-storage/gateway/auth"
)

// ///
// /THIS IS AN AI GENRATE CODE
// //
func main() {
	gatewayURL := "http://127.0.0.1:8080/upload"
	testFileName := "test_up.txt"

	// 1. Generate a valid, signed JWT for testing
	// We pass a mock user ID and an expiration duration.
	// Make sure your local auth.GenerateToken or equivalent signature matches this format!
	testUserID := "user_test_99"
	tokenString, err := auth.GenerateToken(testUserID)
	if err != nil {
		log.Fatalf("Failed to generate test authorization token: %v", err)
	}

	// 2. Create local test file data
	log.Println("Creating local test file...")
	err = os.WriteFile(testFileName, []byte("Cloud architecture verification test! Secure stream pipe cleared."), 0644)
	if err != nil {
		log.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(testFileName)

	file, err := os.Open(testFileName)
	if err != nil {
		log.Fatalf("Failed to open test file: %v", err)
	}
	defer file.Close()

	// 3. Construct a multipart form buffer payload
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("path", "documents/production/test")
	_ = writer.WriteField("filename", "verified_notes.txt")

	part, err := writer.CreateFormFile("file", testFileName)
	if err != nil {
		log.Fatalf("Failed to establish form file field entry: %v", err)
	}
	_, _ = io.Copy(part, file)
	writer.Close()

	// 4. Construct the HTTP Post request
	req, err := http.NewRequest("POST", gatewayURL, body)
	if err != nil {
		log.Fatalf("Request initialization crash: %v", err)
	}

	// Apply headers
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Attach the newly signed token into the Authorization Header
	// (Ensure the prefix matches your middleware style: e.g., "Bearer " or raw string)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	// 5. Send it down the wire straight to the Gateway
	log.Printf("Firing live authenticated HTTP POST request to Gateway at %s...", gatewayURL)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Network transaction failed: %v", err)
	}
	defer resp.Body.Close()

	// 6. Read the response logs
	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("Gateway Status Response Code: %d", resp.StatusCode)
	log.Printf("Gateway Core JSON Body payload: %s", string(respBody))
}
