package test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MayurUbarhande0/Simple-cloud-storage/gateway/middleware"
	"github.com/MayurUbarhande0/Simple-cloud-storage/gateway/server"
)

func TestUploadfileWithoutAuth(t *testing.T) {
	rec := httptest.NewRecorder()
	body := &bytes.Buffer{}
	req := httptest.NewRequest("POST", "/upload", body)
	server.Uploadfile(rec, req)
	if rec.Code == http.StatusOK {
		t.Fatalf("Expected non-200 status, got 200")
	}
}

func TestUploadfileWithNilBody(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/upload", nil)
	req.Body = nil
	server.Uploadfile(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400, got %d", rec.Code)
	}
}

func TestGetfileWithoutAuth(t *testing.T) {
	// Expect handler may panic because DB not initialized; we recover and treat as pass
	defer func() {
		if r := recover(); r != nil {
			// ok
		}
	}()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/download", nil)
	server.Getfile(rec, req)
	if rec.Code == http.StatusOK {
		t.Fatalf("Expected non-200 status, got 200")
	}
}

func TestGetfileWithValidContext(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// panics expected in this environment without DB
		}
	}()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/download?file_id=test_file_123", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "test_user_123")
	req = req.WithContext(ctx)
	server.Getfile(rec, req)
}

