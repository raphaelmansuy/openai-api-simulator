package nanochat

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// GetLatestLlamaTag fetches the latest llama.cpp release tag
func GetLatestLlamaTag() (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get("https://github.com/ggerganov/llama.cpp/releases/latest")
	if err != nil {
		return "", fmt.Errorf("failed to check latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound && resp.StatusCode != http.StatusMovedPermanently {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	location := resp.Header.Get("Location")
	if location == "" {
		return "", fmt.Errorf("no redirect location found")
	}

	// Extract tag from URL like: https://github.com/ggerganov/llama.cpp/releases/tag/b7083
	parts := strings.Split(location, "/")
	if len(parts) < 1 {
		return "", fmt.Errorf("invalid redirect URL: %s", location)
	}
	tag := parts[len(parts)-1]

	return tag, nil
}

// EnsureLlamaServer downloads and extracts llama-server if not present
func EnsureLlamaServer(cacheDir string, platform *PlatformInfo, version string) (string, error) {
	serverPath := filepath.Join(cacheDir, "llama-server")

	// Check if already downloaded
	if FileExists(serverPath) {
		return serverPath, nil
	}

	// Download llama.cpp binary
	zipName := strings.Replace(platform.LlamaCppZip, "{VERSION}", version, 1)
	zipURL := fmt.Sprintf("https://github.com/ggerganov/llama.cpp/releases/download/%s/%s", version, zipName)
	zipPath := filepath.Join(cacheDir, zipName)

	fmt.Printf("Downloading llama.cpp %s...\n", version)
	if err := DownloadWithProgress(zipURL, zipPath, "llama.cpp"); err != nil {
		return "", fmt.Errorf("failed to download llama.cpp: %w", err)
	}

	// Extract the binary
	fmt.Println("Extracting llama-server...")
	if err := extractLlamaServer(zipPath, cacheDir); err != nil {
		return "", fmt.Errorf("failed to extract llama-server: %w", err)
	}

	// Make executable
	if err := os.Chmod(serverPath, 0755); err != nil {
		return "", fmt.Errorf("failed to make llama-server executable: %w", err)
	}

	// Clean up zip file
	os.Remove(zipPath)

	return serverPath, nil
}

// extractLlamaServer extracts llama-server from the zip archive
func extractLlamaServer(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	// Look for llama-server or llama-server.exe
	for _, f := range r.File {
		baseName := filepath.Base(f.Name)
		if baseName == "llama-server" || baseName == "llama-server.exe" {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			destPath := filepath.Join(destDir, baseName)
			out, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer out.Close()

			_, err = io.Copy(out, rc)
			return err
		}
	}

	return fmt.Errorf("llama-server not found in archive")
}

// EnsureModel downloads the GGUF model if not present
func EnsureModel(cacheDir, modelName string) (string, error) {
	config, err := GetModelConfig(modelName)
	if err != nil {
		return "", err
	}

	modelPath := filepath.Join(cacheDir, config.HFFile)

	// Check if already downloaded
	if FileExists(modelPath) {
		return modelPath, nil
	}

	// Download model
	fmt.Printf("Downloading %s (%s)...\n", config.Description, config.Size)
	if err := DownloadWithProgress(config.URL, modelPath, config.Name); err != nil {
		return "", fmt.Errorf("failed to download model: %w", err)
	}

	return modelPath, nil
}

// StartLlamaServer starts the llama-server process
func StartLlamaServer(ctx context.Context, serverPath, modelPath string, platform *PlatformInfo, port int) (*exec.Cmd, error) {
	args := []string{
		"--host", "127.0.0.1",
		"--port", fmt.Sprintf("%d", port),
		"--model", modelPath,
		"--ctx-size", "8192",
		"--temp", "0.7",
		"--n-gpu-layers", fmt.Sprintf("%d", platform.GPULayers),
		"--threads", fmt.Sprintf("%d", 4), // Use 4 threads by default for small models
	}

	cmd := exec.CommandContext(ctx, serverPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start llama-server: %w", err)
	}

	return cmd, nil
}

// WaitForServer waits for the llama-server to be healthy
func WaitForServer(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: 2 * time.Second}

	for time.Now().Before(deadline) {
		resp, err := client.Get(url + "/v1/models")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("llama-server did not become healthy within %v", timeout)
}
