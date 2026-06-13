package routes

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/MayurUbarhande0/Simple-cloud-storage/gateway/auth"
)

// TestNewRouterCreatesValidMux tests that router creates valid mux
func TestNewRouterCreatesValidMux(t *testing.T) {
	mux := NewRouter()

	if mux == nil {
		t.Fatal("NewRouter returned nil")
	}
}

// TestUploadRouteRequiresAuth tests /upload endpoint requires authentication
func TestUploadRouteRequiresAuth(t *testing.T) {
	mux := NewRouter()

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/upload", nil)

	mux.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for /upload without auth, got %d", recorder.Code)
	}
}

// TestDownloadRouteRequiresAuth tests /download endpoint requires authentication
func TestDownloadRouteRequiresAuth(t *testing.T) {
	mux := NewRouter()

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/download", nil)

	mux.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for /download without auth, got %d", recorder.Code)
	}
}

// TestUploadRouteWithValidToken tests /upload accepts valid token
func TestUploadRouteWithValidToken(t *testing.T) {
	os.Setenv("SECRET_KEY", "test-secret-key-123")

	token, err := auth.GenerateToken("test_user")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	mux := NewRouter()

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/upload", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	mux.ServeHTTP(recorder, req)

	if recorder.Code == http.StatusUnauthorized {
		t.Errorf("Expected non-401 status, got %d (auth middleware failed)", recorder.Code)
	}
}

// TestDownloadRouteWithValidToken tests /download accepts valid token
func TestDownloadRouteWithValidToken(t *testing.T) {
	os.Setenv("SECRET_KEY", "test-secret-key-123")

	token, err := auth.GenerateToken("test_user")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	mux := NewRouter()

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/download", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// This test just verifies auth passes - actual download requires DB setup
	// The handler will fail due to nil DbInstance, but auth should succeed
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Handler panicked as expected (DB not initialized): %v", r)
		}
	}()

	mux.ServeHTTP(recorder, req)

	// If we get here without auth error, middleware worked
	if recorder.Code == http.StatusUnauthorized {
		t.Errorf("Expected non-401 status, got %d (auth middleware failed)", recorder.Code)
	}
}

