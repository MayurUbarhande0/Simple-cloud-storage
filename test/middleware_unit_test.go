package test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/MayurUbarhande0/Simple-cloud-storage/gateway/auth"
	"github.com/MayurUbarhande0/Simple-cloud-storage/gateway/middleware"
)

func TestAuthhandlerMissingToken(t *testing.T) {
	handler := middleware.Authhandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/test", nil)

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401, got %d", recorder.Code)
	}
}

func TestAuthhandlerInvalidFormat(t *testing.T) {
	handler := middleware.Authhandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Authorization", "InvalidToken")

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 for invalid format, got %d", recorder.Code)
	}
}

func TestAuthhandlerValidToken(t *testing.T) {
	os.Setenv("SECRET_KEY", "test-secret-key-123")
	token, err := auth.GenerateToken("test_user_456")
	if err != nil {
		t.Fatalf("Failed to generate test token: %v", err)
	}

	nextCalled := false
	handler := middleware.Authhandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		if r.Context().Value(middleware.UserIDKey) == nil {
			t.Error("user id missing in context")
		}
		w.WriteHeader(http.StatusOK)
	}))

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", recorder.Code)
	}
	if !nextCalled {
		t.Fatal("next handler not called")
	}
}

