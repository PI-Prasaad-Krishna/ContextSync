package main

import (
	"fmt"
	"os"
)

const contextFileName = ".context.md"

const defaultContextContent = `# Project Context

This file is automatically maintained by ContextSync. 
It contains a dense, continuously updated summary of recent file changes to provide context to AI coding agents.

## Recent Changes
`

// handleInit creates a boilerplate .context.md file if it doesn't already exist.
func handleInit() error {
	// Check if file exists to avoid overwriting user data
	if _, err := os.Stat(contextFileName); err == nil {
		fmt.Printf("File %s already exists. Skipping initialization.\n", contextFileName)
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	// Create the file with standard rw-r--r-- permissions
	err := os.WriteFile(contextFileName, []byte(defaultContextContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write %s: %w", contextFileName, err)
	}

	fmt.Printf("Successfully initialized %s\n", contextFileName)
	return nil
}
