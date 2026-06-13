package routes

import (
	"net/http"

	"github.com/MayurUbarhande0/Simple-cloud-storage/gateway/middleware"
	"github.com/MayurUbarhande0/Simple-cloud-storage/gateway/server"
)

// NewRouter registers paths and wires up handlers with middleware
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// 1. Create a handler chain that forces authentication checks
	// This ensures r.Context().Value(middleware.UserIDKey) is never nil inside Uploadfile!
	uploadHandler := http.HandlerFunc(server.Uploadfile)
	authenticatedUpload := middleware.Authhandler(uploadHandler) // Wrap it up!

	// 2. Register the upload endpoint route pattern
	mux.Handle("/upload", authenticatedUpload)

	// 3. Wrap and register the download endpoint with middleware
	downloadHandler := http.HandlerFunc(server.Getfile)
	authenticatedDownload := middleware.Authhandler(downloadHandler)
	mux.Handle("/download", authenticatedDownload)

	return mux
}
