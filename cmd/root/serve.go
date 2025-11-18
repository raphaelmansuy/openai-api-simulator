package root

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/openai/openai-api-simulator/pkg/server"
	"github.com/openai/openai-api-simulator/pkg/streaming"
	"github.com/spf13/cobra"
)

var (
	port         int
	delayMin     int
	delayMax     int
	tokensPerSec float64
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the OpenAI API simulator server (fake mode)",
	Long: `Start the OpenAI API simulator server with fake/simulated responses.
This mode generates coherent random text without requiring any LLM.
Perfect for testing, development, and CI/CD pipelines.`,
	RunE: runServe,
}

func init() {
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the simulator HTTP server on")
	serveCmd.Flags().IntVar(&delayMin, "stream-delay-min-ms", 0, "Default min per-chunk delay (ms) to simulate jitter")
	serveCmd.Flags().IntVar(&delayMax, "stream-delay-max-ms", 0, "Default max per-chunk delay (ms) to simulate jitter")
	serveCmd.Flags().Float64Var(&tokensPerSec, "stream-tokens-per-second", 0, "Default token emission rate for streaming chunks; 0 disables throttling")

	rootCmd.AddCommand(serveCmd)
}

func runServe(cmd *cobra.Command, args []string) error {
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting OpenAI API Simulator on %s", addr)

	defaults := streaming.StreamOptions{}
	if delayMin > 0 {
		defaults.DelayMin = time.Duration(delayMin) * time.Millisecond
	}
	if delayMax > 0 {
		defaults.DelayMax = time.Duration(delayMax) * time.Millisecond
	}
	if tokensPerSec > 0 {
		defaults.TokensPerSecond = tokensPerSec
	}

	handler := server.NewRouterWithStreamDefaults(defaults)
	if err := http.ListenAndServe(addr, handler); err != nil {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}
