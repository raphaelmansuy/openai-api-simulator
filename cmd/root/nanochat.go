package root

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/openai/openai-api-simulator/internal/nanochat"
	"github.com/openai/openai-api-simulator/pkg/server"
	"github.com/openai/openai-api-simulator/pkg/streaming"
	"github.com/spf13/cobra"
)

var (
	nanochatModel     string
	nanochatModelPath string
	nanochatLlamaPort int
	nanochatPort      int
)

var nanochatCmd = &cobra.Command{
	Use:   "nanochat",
	Short: "One-command real local inference (auto-downloads llama.cpp + model)",
	Long: `Start the OpenAI API simulator with REAL local inference using llama.cpp.

This command will:
  1. Auto-detect your platform (macOS/Linux, Intel/ARM)
  2. Download llama.cpp binary if not cached
  3. Download a small GGUF model if not cached (default: nanochat ~316MB)
  4. Start llama.cpp server in the background
  5. Start the simulator proxy that routes to llama.cpp

First run takes ~30-90 seconds depending on your connection.
Subsequent runs are instant!

Example:
  openai-api-simulator nanochat                    # Use default nanochat model
  openai-api-simulator nanochat --model phi3.5-mini # Use larger Phi-3.5 model
  openai-api-simulator nanochat --model-path ./my-model.gguf # Use custom model`,
	RunE: runNanoChat,
}

func init() {
	nanochatCmd.Flags().StringVar(&nanochatModel, "model", nanochat.DefaultModel(), "Pre-defined model: nanochat (default), phi3.5-mini, qwen2.5-3b, gemma2-2b, tinyllama")
	nanochatCmd.Flags().StringVar(&nanochatModelPath, "model-path", "", "Custom local GGUF path (bypass download)")
	nanochatCmd.Flags().IntVar(&nanochatLlamaPort, "llama-port", 8081, "Internal llama.cpp port")
	nanochatCmd.Flags().IntVar(&nanochatPort, "port", 3080, "Public simulator port")

	rootCmd.AddCommand(nanochatCmd)
}

func runNanoChat(cmd *cobra.Command, args []string) error {
	// Step 1: Detect platform
	fmt.Print("[1/4] Detecting platform... ")
	platform, err := nanochat.DetectPlatform()
	if err != nil {
		fmt.Println("❌")
		return err
	}
	fmt.Printf("✓ %s\n", platform.Description)
	fmt.Printf("      GPU: %s\n", platform.GPUSupport)

	// Step 2: Ensure cache directory
	cacheDir := nanochat.GetCacheDir()
	if err := nanochat.EnsureCacheDir(cacheDir); err != nil {
		return err
	}

	// Step 3: Download llama.cpp if needed
	var serverPath string
	if nanochat.FileExists(cacheDir + "/llama-server") {
		fmt.Println("[2/4] llama.cpp binary found ✓")
		serverPath = cacheDir + "/llama-server"
	} else {
		fmt.Print("[2/4] Downloading llama.cpp... ")
		version, err := nanochat.GetLatestLlamaTag()
		if err != nil {
			fmt.Println("❌")
			return fmt.Errorf("failed to get latest llama.cpp version: %w", err)
		}
		fmt.Printf("(version %s)\n", version)

		serverPath, err = nanochat.EnsureLlamaServer(cacheDir, platform, version)
		if err != nil {
			return err
		}
		fmt.Println("      ✓ Downloaded and extracted")
	}

	// Step 4: Download model if needed
	var modelPath string
	if nanochatModelPath != "" {
		// Use custom model path
		if !nanochat.FileExists(nanochatModelPath) {
			return fmt.Errorf("custom model file not found: %s", nanochatModelPath)
		}
		modelPath = nanochatModelPath
		fmt.Printf("[3/4] Using custom model: %s ✓\n", nanochatModelPath)
	} else {
		// Use predefined model
		config, err := nanochat.GetModelConfig(nanochatModel)
		if err != nil {
			return err
		}

		cachedModelPath := cacheDir + "/" + config.HFFile
		if nanochat.FileExists(cachedModelPath) {
			fmt.Printf("[3/4] Model found: %s ✓\n", config.Name)
			modelPath = cachedModelPath
		} else {
			fmt.Printf("[3/4] Downloading %s (%s)...\n", config.Name, config.Size)
			modelPath, err = nanochat.EnsureModel(cacheDir, nanochatModel)
			if err != nil {
				return err
			}
			fmt.Println("      ✓ Downloaded")
		}
	}

	// Step 5: Start llama-server
	fmt.Printf("[4/4] Starting llama.cpp server on 127.0.0.1:%d...\n", nanochatLlamaPort)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	llamaCmd, err := nanochat.StartLlamaServer(ctx, serverPath, modelPath, platform, nanochatLlamaPort)
	if err != nil {
		return err
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n\nShutting down gracefully...")
		cancel()
		if llamaCmd.Process != nil {
			llamaCmd.Process.Kill()
		}
		os.Exit(0)
	}()

	// Wait for llama-server to be healthy
	fmt.Print("      Waiting for llama.cpp to be ready")
	llamaURL := fmt.Sprintf("http://127.0.0.1:%d", nanochatLlamaPort)
	if err := nanochat.WaitForServer(llamaURL, 60*time.Second); err != nil {
		fmt.Println(" ❌")
		cancel()
		return err
	}
	fmt.Println(" ✓")

	// Step 6: Start simulator proxy
	fmt.Printf("\n🚀 NanoChat ready!\n")
	fmt.Printf("   Simulator running on http://localhost:%d\n", nanochatPort)
	fmt.Printf("   Backend: llama.cpp on http://127.0.0.1:%d\n", nanochatLlamaPort)
	fmt.Printf("   Model: %s\n\n", nanochatModel)
	fmt.Printf("Test with:\n")
	fmt.Printf("  curl http://localhost:%d/v1/chat/completions -d '{\"model\":\"nanochat\",\"messages\":[{\"role\":\"user\",\"content\":\"Hello!\"}]}'\n\n", nanochatPort)

	// Start the simulator in proxy mode
	addr := fmt.Sprintf(":%d", nanochatPort)
	handler := server.NewRouterWithStreamDefaults(streaming.StreamOptions{})

	log.Printf("Starting simulator proxy on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		return fmt.Errorf("simulator proxy failed: %w", err)
	}

	return nil
}
