package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/MayurUbarhande0/Simple-cloud-storage/gateway/helper"
	"github.com/MayurUbarhande0/Simple-cloud-storage/gateway/middleware"
	"github.com/MayurUbarhande0/Simple-cloud-storage/gateway/statemanager"
	"github.com/rs/xid"
)

var StateMgr *statemanager.Manager

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

	conn, err := net.Dial("tcp", os.Getenv("STORAGE_IP"))
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Internal server error")
		return
	}
	defer conn.Close()
	ENkey := os.Getenv("ENCRYPTION_KEY") //EN key to encrpyt auth key
	AUTH_KEY := os.Getenv("AUTH_KEY")    //gateway auth key to storage server

	EN_AUTH_KEY, err := helper.Encrypt([]byte(AUTH_KEY), []byte(ENkey))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Internal server error")
		return
	}
	tokenLen := byte(len(EN_AUTH_KEY))

	if _, err := conn.Write([]byte{tokenLen}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Internal server error")
		return
	}
	if _, err := conn.Write(EN_AUTH_KEY); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Internal server error")
		return
	}

	filesize := header.Size
	ServerStoragePath := fmt.Sprintf("/cloud/%s/%s/%s", userIDStr, path, file_id)
	handshake := fmt.Sprintf("%s|%s|%d\n", "0x01", ServerStoragePath, filesize)

	if _, err := conn.Write([]byte(handshake)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Internal server err")
		return
	}
	// send file contents to storage server
	if _, err := io.Copy(conn, file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Internal server err")
		return
	}
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"file_id": file_id, "status": "success"})

}
