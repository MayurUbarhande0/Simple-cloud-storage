package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/MayurUbarhande0/Simple-cloud-storage/gateway/auth"
)

// TestAuthhandlerMissingToken tests middleware rejects requests without token
func TestAuthhandlerMissingToken(t *testing.T) {
	handler := Authhandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/test", nil)

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401, got %d", recorder.Code)
	}
}

// TestAuthhandlerInvalidFormat tests middleware rejects invalid token format
func TestAuthhandlerInvalidFormat(t *testing.T) {
	handler := Authhandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Authorization", "InvalidToken")

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for invalid format, got %d", recorder.Code)
	}
}

// TestAuthhandlerValidToken tests middleware accepts valid Bearer token
func TestAuthhandlerValidToken(t *testing.T) {
	os.Setenv("SECRET_KEY", "test-secret-key-123")
	
	token, err := auth.GenerateToken("test_user_456")
	if err != nil {
		t.Fatalf("Failed to generate test token: %v", err)
	}

	nextHandlerCalled := false
	handler := Authhandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextHandlerCalled = true
		userID := r.Context().Value(UserIDKey)
		if userID == nil {
			t.Error("Expected user_id in context")
		}
		w.WriteHeader(http.StatusOK)
	}))

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", recorder.Code)
	}

	if !nextHandlerCalled {
		t.Error("Next handler was not called")
	}
}

