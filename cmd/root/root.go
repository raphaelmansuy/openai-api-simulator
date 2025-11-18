package root

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "openai-api-simulator",
	Short: "OpenAI API Simulator - A local OpenAI-compatible API server",
	Long: `OpenAI API Simulator is a lightweight, dependency-free OpenAI-compatible 
chat completion simulator that can run with fake responses or real local inference.

It provides:
  - Fake mode (default): Fast, predictable responses for testing
  - NanoChat mode: Real local inference with llama.cpp`,
	Version: "1.0.0",
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// GetRootCmd returns the root command for adding subcommands
func GetRootCmd() *cobra.Command {
	return rootCmd
}
