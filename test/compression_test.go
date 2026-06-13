package test

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"

	utils "github.com/MayurUbarhande0/Simple-cloud-storage/IO"
)

// TestCompressionRoundTrip verifies that IO.Compress writes valid gzip data
// and that the compressed payload round-trips back to the original bytes.
func TestCompressionRoundTrip(t *testing.T) {
	oldHome := os.Getenv("HOME")
	tempHome := t.TempDir()
	if err := os.Setenv("HOME", tempHome); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}
	defer os.Setenv("HOME", oldHome)

	original := bytes.Repeat([]byte("The quick brown fox jumps over the lazy dog. "), 50)
	storagePath := "cloud/test-user/docs/sample.dat"

	if err := utils.Compress(bytes.NewReader(original), storagePath, int64(len(original))); err != nil {
		t.Fatalf("Compress returned error: %v", err)
	}

	absPath := filepath.Join(tempHome, storagePath)
	compressed, err := os.ReadFile(absPath)
	if err != nil {
		t.Fatalf("failed to read compressed file: %v", err)
	}

	if len(compressed) < 2 || compressed[0] != 0x1f || compressed[1] != 0x8b {
		t.Fatalf("file is not gzip data; first bytes=% x", compressed[:min(2, len(compressed))])
	}

	zr, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		t.Fatalf("failed to open gzip stream: %v", err)
	}
	defer zr.Close()

	decoded, err := io.ReadAll(zr)
	if err != nil {
		t.Fatalf("failed to decompress: %v", err)
	}

	if !bytes.Equal(decoded, original) {
		t.Fatalf("decompressed output does not match original")
	}

	t.Logf("original=%d bytes compressed=%d bytes path=%s", len(original), len(compressed), absPath)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

