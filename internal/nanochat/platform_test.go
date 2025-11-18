package nanochat

import (
	"runtime"
	"strings"
	"testing"
)

func TestDetectPlatform(t *testing.T) {
	platform, err := DetectPlatform()
	if err != nil {
		t.Fatalf("DetectPlatform() failed: %v", err)
	}

	if platform == nil {
		t.Fatal("DetectPlatform() returned nil platform")
	}

	// Check that OS and Arch are set correctly
	if platform.OS != runtime.GOOS {
		t.Errorf("Expected OS %s, got %s", runtime.GOOS, platform.OS)
	}

	if platform.Arch != runtime.GOARCH {
		t.Errorf("Expected Arch %s, got %s", runtime.GOARCH, platform.Arch)
	}

	// Check that BinaryName is set
	if platform.BinaryName == "" {
		t.Error("BinaryName is empty")
	}

	// Check that Description is set
	if platform.Description == "" {
		t.Error("Description is empty")
	}

	// Check that LlamaCppZip is set and contains {VERSION}
	if platform.LlamaCppZip == "" {
		t.Error("LlamaCppZip is empty")
	}
	if !strings.Contains(platform.LlamaCppZip, "{VERSION}") {
		t.Error("LlamaCppZip should contain {VERSION} placeholder")
	}

	// Check GPU settings are reasonable
	if platform.GPULayers < 0 {
		t.Errorf("GPU layers should be >= 0, got %d", platform.GPULayers)
	}

	// Apple Silicon should have Metal GPU offload
	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		if platform.GPULayers == 0 {
			t.Error("Expected Apple Silicon to have GPU offload enabled")
		}
		if !strings.Contains(platform.GPUSupport, "Metal") {
			t.Errorf("Expected Metal GPU support on Apple Silicon, got: %s", platform.GPUSupport)
		}
	}
}

func TestDetectPlatform_SupportedPlatforms(t *testing.T) {
	// This test verifies current platform is supported
	platform, err := DetectPlatform()
	if err != nil {
		// This is expected on unsupported platforms
		if !strings.Contains(err.Error(), "unsupported platform") {
			t.Errorf("Unexpected error: %v", err)
		}
		t.Skipf("Platform %s/%s is not supported (expected)", runtime.GOOS, runtime.GOARCH)
	}

	if platform == nil {
		t.Fatal("Platform is nil but no error returned")
	}
}
