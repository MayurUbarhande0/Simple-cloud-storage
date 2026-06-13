package utils

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Compress reads the incoming payload, expands the home directory path,
// creates the necessary folders, and writes compressed gzip blocks directly to disk.
func Compress(reader io.Reader, storagePath string, fileSize int64) error {
	// 1. Resolve the absolute path (Fixing the "~/"" issue)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Converts "cloud/user_id/path/file_id" into "/home/username/cloud/user_id/path/file_id"
	absolutePath := filepath.Join(homeDir, storagePath)

	// 2. Ensure all parent directories exist on the server disk
	// os.MkdirAll does nothing if they already exist
	parentDir := filepath.Dir(absolutePath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("failed to create storage directories: %w", err)
	}

	// 3. Create the physical target file on disk
	outputFile, err := os.Create(absolutePath)
	if err != nil {
		return fmt.Errorf("failed to create target file: %w", err)
	}
	defer outputFile.Close() // Explicitly close when streaming is done

	// 4. Wrap the file writer in a Gzip Writer
	gzipWriter := gzip.NewWriter(outputFile)
	defer gzipWriter.Close() // Flushes remaining gzip bits on return

	// 5. Read EXACTLY the file size out of the stream using LimitReader
	// This ensures we don't accidentally read data from a subsequent action
	limitedStream := io.LimitReader(reader, fileSize)

	// Stream the raw network bytes into the gzip encoder, straight to disk
	_, err = io.Copy(gzipWriter, limitedStream)
	if err != nil {
		return fmt.Errorf("failed to stream compressed blocks: %w", err)
	}

	return nil
}

// GetFileAbsolutePath converts your database storage path directly to an OS readable path
func GetFileAbsolutePath(storagePath string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// No walking required! Instant lookup since the unique file_id is at the end of the path
	return filepath.Join(homeDir, storagePath), nil
}
