package nanochat

import (
	"strings"
	"testing"
)

func TestDefaultModel(t *testing.T) {
	defaultModel := DefaultModel()
	if defaultModel == "" {
		t.Error("DefaultModel() returned empty string")
	}

	// Default model should be "nanochat"
	if defaultModel != "nanochat" {
		t.Errorf("Expected default model to be 'nanochat', got '%s'", defaultModel)
	}

	// Default model should exist in predefined models
	_, err := GetModelConfig(defaultModel)
	if err != nil {
		t.Errorf("Default model '%s' not found in predefined models: %v", defaultModel, err)
	}
}

func TestGetModelConfig(t *testing.T) {
	tests := []struct {
		name      string
		modelName string
		wantErr   bool
	}{
		{
			name:      "valid model - nanochat",
			modelName: "nanochat",
			wantErr:   false,
		},
		{
			name:      "valid model - phi3.5-mini",
			modelName: "phi3.5-mini",
			wantErr:   false,
		},
		{
			name:      "valid model - qwen2.5-3b",
			modelName: "qwen2.5-3b",
			wantErr:   false,
		},
		{
			name:      "valid model - gemma2-2b",
			modelName: "gemma2-2b",
			wantErr:   false,
		},
		{
			name:      "valid model - tinyllama",
			modelName: "tinyllama",
			wantErr:   false,
		},
		{
			name:      "invalid model",
			modelName: "nonexistent-model",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := GetModelConfig(tt.modelName)
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify config fields are populated
			if config.Name == "" {
				t.Error("Model name is empty")
			}
			if config.HFRepo == "" {
				t.Error("HF repo is empty")
			}
			if config.HFFile == "" {
				t.Error("HF file is empty")
			}
			if config.URL == "" {
				t.Error("URL is empty")
			}
			if config.Size == "" {
				t.Error("Size is empty")
			}
			if config.Description == "" {
				t.Error("Description is empty")
			}

			// URL should start with https://huggingface.co
			if !strings.HasPrefix(config.URL, "https://huggingface.co") {
				t.Errorf("URL should start with https://huggingface.co, got: %s", config.URL)
			}

			// URL should contain the file name
			if !strings.Contains(config.URL, config.HFFile) {
				t.Errorf("URL should contain file name %s, got: %s", config.HFFile, config.URL)
			}
		})
	}
}

func TestPredefinedModels(t *testing.T) {
	if len(PredefinedModels) == 0 {
		t.Fatal("PredefinedModels is empty")
	}

	// Verify all predefined models have valid configurations
	for name, config := range PredefinedModels {
		t.Run(name, func(t *testing.T) {
			if config.Name == "" {
				t.Errorf("Model %s has empty name", name)
			}
			if config.URL == "" {
				t.Errorf("Model %s has empty URL", name)
			}
			if config.HFFile == "" {
				t.Errorf("Model %s has empty HFFile", name)
			}
			if config.Size == "" {
				t.Errorf("Model %s has empty Size", name)
			}
		})
	}
}

func TestNanoChatIsDefault(t *testing.T) {
	// Verify nanochat exists and is properly configured
	config, err := GetModelConfig("nanochat")
	if err != nil {
		t.Fatalf("nanochat model not found: %v", err)
	}

	// nanochat should be the smallest/fastest to download
	if !strings.Contains(config.Description, "561M") && !strings.Contains(config.Description, "561") {
		t.Logf("Warning: nanochat description doesn't mention 561M parameters: %s", config.Description)
	}

	// URL should point to sdobson/nanochat
	if !strings.Contains(config.URL, "sdobson/nanochat") {
		t.Errorf("Expected URL to contain sdobson/nanochat, got: %s", config.URL)
	}
}
