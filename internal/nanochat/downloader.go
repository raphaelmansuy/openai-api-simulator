package nanochat

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

// DownloadWithProgress downloads a file with a progress bar
func DownloadWithProgress(url, destPath, description string) error {
	// Create temporary file
	tmpPath := destPath + ".tmp"

	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download %s: HTTP %d", url, resp.StatusCode)
	}

	// Create destination file
	out, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", tmpPath, err)
	}
	defer out.Close()

	// Create progress bar
	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		description,
	)

	// Copy with progress
	_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
	if err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to download file: %w", err)
	}

	// Rename temp file to final destination
	if err := os.Rename(tmpPath, destPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

// EnsureCacheDir creates the cache directory if it doesn't exist
func EnsureCacheDir(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}
	return nil
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// GetCacheDir returns the cache directory path
func GetCacheDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "/tmp"
	}
	return filepath.Join(home, ".cache", "openai-api-simulator")
}
