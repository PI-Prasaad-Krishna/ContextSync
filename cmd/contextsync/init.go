package main

import (
	"fmt"
	"os"
)

const defaultContextContent = `# Project Context

This file is automatically maintained by ContextSync. 
It contains a dense, continuously updated summary of recent file changes to provide context to AI coding agents.

## Recent Changes
`

// handleInit creates a boilerplate context file if it doesn't already exist.
func handleInit(cfg Config) error {
	// Check if file exists to avoid overwriting user data
	if _, err := os.Stat(cfg.OutFile); err == nil {
		fmt.Printf("File %s already exists. Skipping initialization.\n", cfg.OutFile)
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	// Create the file with standard rw-r--r-- permissions
	err := os.WriteFile(cfg.OutFile, []byte(defaultContextContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write %s: %w", cfg.OutFile, err)
	}

	fmt.Printf("Successfully initialized %s\n", cfg.OutFile)
	return nil
}
