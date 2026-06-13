package server

import (
	"encoding/json"
	"fmt"
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
	ServerStoragePath := fmt.Sprintf("/cloud/%s/%s/%s", userIDStr, path, file_id)

	metadata := &statemanager.Fileheader{
		File_id:     file_id,
		Filename:    filename,
		Uploaded:    time.Now(),
		StoragePath: path,
		Size:        int(filesize),
		LastUpdated: time.Now(),
	}
	if err = StateMgr.Addfile(userIDStr, *metadata); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Internal Server Error")
		return
	}
	if err := utils.Compress(file, ServerStoragePath, filesize); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
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

	strg, err := DbInstance.GetPath(UserIDstr, file_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Internal Server Error")
		return

	}
	http.ServeFile(w, r, strg)

}
