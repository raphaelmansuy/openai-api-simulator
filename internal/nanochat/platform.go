package nanochat

import (
	"fmt"
	"runtime"
)

// PlatformInfo contains platform-specific information
type PlatformInfo struct {
	OS           string
	Arch         string
	Description  string
	GPULayers    int
	GPUSupport   string
	BinaryName   string
	LlamaCppZip  string
}

// DetectPlatform detects the current OS and architecture
func DetectPlatform() (*PlatformInfo, error) {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	var platform PlatformInfo
	platform.OS = goos
	platform.Arch = goarch
	platform.BinaryName = "llama-server"

	switch {
	case goos == "darwin" && goarch == "arm64":
		platform.Description = "macOS arm64 (Apple Silicon)"
		platform.GPULayers = 999 // Full Metal GPU offload
		platform.GPUSupport = "Metal (auto-enabled)"
		platform.LlamaCppZip = "llama-{VERSION}-bin-macos-arm64.zip"
	case goos == "darwin" && goarch == "amd64":
		platform.Description = "macOS x64 (Intel)"
		platform.GPULayers = 0 // CPU only
		platform.GPUSupport = "CPU only"
		platform.LlamaCppZip = "llama-{VERSION}-bin-macos-x64.zip"
	case goos == "linux" && goarch == "amd64":
		platform.Description = "Linux x86_64"
		platform.GPULayers = 0 // CPU (Vulkan possible but not in official prebuilt)
		platform.GPUSupport = "CPU"
		platform.LlamaCppZip = "llama-{VERSION}-bin-ubuntu-x64.zip"
	case goos == "linux" && goarch == "arm64":
		platform.Description = "Linux arm64"
		platform.GPULayers = 0 // CPU
		platform.GPUSupport = "CPU"
		platform.LlamaCppZip = "llama-{VERSION}-bin-ubuntu-arm64.zip"
	default:
		return nil, fmt.Errorf("unsupported platform: %s/%s (only macOS arm64/x64 and Linux x86_64/arm64 are supported)", goos, goarch)
	}

	return &platform, nil
}
