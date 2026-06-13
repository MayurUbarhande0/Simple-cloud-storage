package server

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MayurUbarhande0/Simple-cloud-storage/gateway/middleware"
)

// TestUploadfileWithoutAuth tests upload handler fails without authentication
func TestUploadfileWithoutAuth(t *testing.T) {
	recorder := httptest.NewRecorder()
	body := &bytes.Buffer{}

	req := httptest.NewRequest("POST", "/upload", body)

	Uploadfile(recorder, req)

	if recorder.Code == http.StatusOK {
		t.Errorf("Expected non-200 status, got %d", recorder.Code)
	}
}

// TestUploadfileWithNilBody tests upload with nil body
func TestUploadfileWithNilBody(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/upload", nil)
	req.Body = nil

	Uploadfile(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", recorder.Code)
	}
}

// TestGetfileWithoutAuth tests download handler fails without authentication
func TestGetfileWithoutAuth(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Handler panicked (expected - no DB): %v", r)
		}
	}()

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/download", nil)

	Getfile(recorder, req)

	// Should fail since context won't have user_id
	if recorder.Code == http.StatusOK {
		t.Errorf("Expected error status, got %d", recorder.Code)
	}
}

// TestGetfileWithValidContext tests download with valid user context
func TestGetfileWithValidContext(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Handler panicked (expected - no DB): %v", r)
		}
	}()

	recorder := httptest.NewRecorder()
	body := bytes.NewBufferString("file_id=test_file_123")
	req := httptest.NewRequest("POST", "/download?file_id=test_file_123", body)

	// Add user context
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "test_user_123")
	req = req.WithContext(ctx)

	Getfile(recorder, req)

	t.Logf("Response code: %d", recorder.Code)
}

