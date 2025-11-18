package nanochat

import "fmt"

// ModelConfig defines a GGUF model configuration
type ModelConfig struct {
	Name        string
	HFRepo      string
	HFFile      string
	URL         string
	Size        string
	Description string
}

// PredefinedModels contains the available model presets
var PredefinedModels = map[string]ModelConfig{
	"nanochat": {
		Name:        "nanochat",
		HFRepo:      "sdobson/nanochat",
		HFFile:      "nanochat-Q4_K_M.gguf",
		URL:         "https://huggingface.co/sdobson/nanochat/resolve/main/nanochat-Q4_K_M.gguf",
		Size:        "~316 MB",
		Description: "NanoChat - Community-trained 561M-parameter model optimized for chat",
	},
	"phi3.5-mini": {
		Name:        "phi3.5-mini",
		HFRepo:      "bartowski/Phi-3.5-mini-instruct-GGUF",
		HFFile:      "Phi-3.5-mini-instruct-Q4_K_M.gguf",
		URL:         "https://huggingface.co/bartowski/Phi-3.5-mini-instruct-GGUF/resolve/main/Phi-3.5-mini-instruct-Q4_K_M.gguf",
		Size:        "~2.4 GB",
		Description: "Phi-3.5-mini - Microsoft's 3.8B-parameter model with strong reasoning",
	},
	"qwen2.5-3b": {
		Name:        "qwen2.5-3b",
		HFRepo:      "Qwen/Qwen2.5-3B-Instruct-GGUF",
		HFFile:      "qwen2.5-3b-instruct-q4_k_m.gguf",
		URL:         "https://huggingface.co/Qwen/Qwen2.5-3B-Instruct-GGUF/resolve/main/qwen2.5-3b-instruct-q4_k_m.gguf",
		Size:        "~1.9 GB",
		Description: "Qwen2.5-3B - Alibaba's multilingual 3B model",
	},
	"gemma2-2b": {
		Name:        "gemma2-2b",
		HFRepo:      "bartowski/gemma-2-2b-it-GGUF",
		HFFile:      "gemma-2-2b-it-Q4_K_M.gguf",
		URL:         "https://huggingface.co/bartowski/gemma-2-2b-it-GGUF/resolve/main/gemma-2-2b-it-Q4_K_M.gguf",
		Size:        "~1.6 GB",
		Description: "Gemma2-2B - Google's efficient 2B instruction-tuned model",
	},
	"tinyllama": {
		Name:        "tinyllama",
		HFRepo:      "TheBloke/TinyLlama-1.1B-Chat-v1.0-GGUF",
		HFFile:      "tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf",
		URL:         "https://huggingface.co/TheBloke/TinyLlama-1.1B-Chat-v1.0-GGUF/resolve/main/tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf",
		Size:        "~669 MB",
		Description: "TinyLlama - Compact 1.1B chat model",
	},
}

// DefaultModel returns the default model name
func DefaultModel() string {
	return "nanochat"
}

// GetModelConfig returns the configuration for a model name
func GetModelConfig(name string) (ModelConfig, error) {
	config, ok := PredefinedModels[name]
	if !ok {
		return ModelConfig{}, fmt.Errorf("unknown model: %s (available: nanochat, phi3.5-mini, qwen2.5-3b, gemma2-2b, tinyllama)", name)
	}
	return config, nil
}
