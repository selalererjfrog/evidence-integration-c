package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// writeToFile writes data to a file
func writeToFile(filename string, data []byte) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}

	return os.WriteFile(filename, data, 0644)
}
