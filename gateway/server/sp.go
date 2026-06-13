package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	utils "github.com/MayurUbarhande0/Simple-cloud-storage/IO"
	helper2 "github.com/MayurUbarhande0/Simple-cloud-storage/db"
	"github.com/MayurUbarhande0/Simple-cloud-storage/gateway/middleware"
	"github.com/MayurUbarhande0/Simple-cloud-storage/gateway/statemanager"
	"github.com/rs/xid"
)

var StateMgr *statemanager.Manager
var DbInstance *helper2.Db

func Uploadfile(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid request")
		return
	}
	// ensure body check returns
	User_id := r.Context().Value(middleware.UserIDKey)

	// convert user id to string for paths
	userIDStr := fmt.Sprint(User_id)
	path := r.FormValue("path")
	filename := r.FormValue("filename")

	file_id := xid.New().String()
	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	filesize := header.Size
	// store a relative path (no leading slash) so IO utilities will resolve
	// it under the user's home directory (see IO.Compress/GetFileAbsolutePath)
	ServerStoragePath := fmt.Sprintf("cloud/%s/%s/%s", userIDStr, path, file_id)

	metadata := &statemanager.Fileheader{
		File_id:     file_id,
		Filename:    filename,
		Uploaded:    time.Now(),
		StoragePath: ServerStoragePath,
		Size:        int(filesize),
		LastUpdated: time.Now(),
	}
	if err = StateMgr.Addfile(userIDStr, *metadata); err != nil {
		log.Printf("[ERROR] Failed to add file to database: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Database error: %v", err)})
		return
	}
	log.Printf("[SUCCESS] File metadata saved to database - UserID: %s, FileID: %s", userIDStr, file_id)
	if err := utils.Compress(file, ServerStoragePath, filesize); err != nil {
		log.Printf("[ERROR] Failed to compress/save file: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("File storage error: %v", err)})
		return
	}
	log.Printf("[SUCCESS] File compressed and stored at: %s", ServerStoragePath)

	w.WriteHeader(http.StatusOK)
	log.Printf("[UPLOAD_COMPLETE] File uploaded successfully - FileID: %s, Size: %d bytes", file_id, filesize)
	json.NewEncoder(w).Encode(map[string]string{"file_id": file_id, "status": "success"})

}

func Getfile(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("invalid request")
		return

	}
	User_id := r.Context().Value(middleware.UserIDKey)

	UserIDstr := fmt.Sprint(User_id)
	file_id := r.FormValue("file_id")

	log.Printf("[DOWNLOAD_REQUEST] UserID: %s, FileID: %s", UserIDstr, file_id)

	strg, err := DbInstance.GetPath(UserIDstr, file_id)
	if err != nil {
		log.Printf("[ERROR] Failed to get file path from database: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("File not found: %v", err)})
		return
	}

	// Convert stored storage path to an OS absolute path before serving
	absPath, pErr := utils.GetFileAbsolutePath(strg)
	if pErr != nil {
		log.Printf("[ERROR] Failed to resolve absolute file path: %v", pErr)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Path resolution error: %v", pErr)})
		return
	}
	log.Printf("[SUCCESS] Retrieved file path: %s (absolute: %s)", strg, absPath)
	http.ServeFile(w, r, absPath)
	log.Printf("[DOWNLOAD_COMPLETE] File served successfully - FileID: %s", file_id)

}
