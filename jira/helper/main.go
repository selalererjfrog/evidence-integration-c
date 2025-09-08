package main

import (
	"fmt"
	"os"
)

func main() {
	// Parse command line flags
	flags, args := ParseFlags()

	// Handle help flags
	if flags.Help || flags.HelpLong {
		DisplayUsage()
		return
	}

	// Load configuration
	config, err := LoadConfig(flags, args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		DisplayUsage()
		os.Exit(1)
	}

	// Determine and execute the appropriate mode
	if err := determineExecutionMode(flags, args, config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
