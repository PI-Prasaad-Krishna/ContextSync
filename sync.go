package main

import (
	"fmt"
	"os"
	"time"
)

// syncBatchToContext is the MVP "Sync Bridge".
// It takes a batch of changed file paths and appends a simple log to .context.md.
// This gives AI agents raw context about what files they should inspect.
func syncBatchToContext(files map[string]struct{}) error {
	if len(files) == 0 {
		return nil
	}

	// Open the context file in append mode. Create it if it somehow doesn't exist.
	f, err := os.OpenFile(contextFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open context file for sync: %w", err)
	}
	defer f.Close()

	// Write the sync header with a timestamp
	timestamp := time.Now().Format(time.RFC1123)
	header := fmt.Sprintf("\n### Sync Event: %s\nFiles changed:\n", timestamp)
	if _, err := f.WriteString(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write each file that changed
	for file := range files {
		line := fmt.Sprintf("- %s\n", file)
		if _, err := f.WriteString(line); err != nil {
			return fmt.Errorf("failed to write file entry: %w", err)
		}
	}

	fmt.Printf("Successfully synced %d files to %s\n", len(files), contextFileName)
	return nil
}
