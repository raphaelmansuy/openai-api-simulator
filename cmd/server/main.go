package main

import (
	"os"

	"github.com/openai/openai-api-simulator/cmd/root"
)

func main() {
	// For backward compatibility, if run with old flags, run serve command
	if len(os.Args) > 1 && (os.Args[1] == "-port" || os.Args[1] == "--port" || 
		os.Args[1] == "-stream_delay_min_ms" || os.Args[1] == "--stream_delay_min_ms") {
		// Insert "serve" as first argument
		os.Args = append([]string{os.Args[0], "serve"}, os.Args[1:]...)
	} else if len(os.Args) == 1 {
		// If no arguments, default to serve command for backward compatibility
		os.Args = append(os.Args, "serve")
	}

	root.Execute()
}
