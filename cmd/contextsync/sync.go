package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// syncBatchToContext is the MVP "Sync Bridge".
// It takes a batch of changed file paths, extracts the git diff, and appends it to .context.md.
func syncBatchToContext(cfg Config, files map[string]struct{}) error {
	if len(files) == 0 {
		return nil
	}

	// Open the context file in append mode. Create it if it somehow doesn't exist.
	f, err := os.OpenFile(cfg.OutFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open context file for sync: %w", err)
	}

	timestamp := time.Now().Format(time.RFC1123)
	header := fmt.Sprintf("\n### Sync Event: %s\n", timestamp)
	if _, err := f.WriteString(header); err != nil {
		f.Close()
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write each file diff
	for file := range files {
		diff := getFileDiff(file)
		if _, err := f.WriteString(diff); err != nil {
			f.Close()
			return fmt.Errorf("failed to write file entry: %w", err)
		}
	}
	f.Close()

	fmt.Printf("Successfully synced %d files to %s\n", len(files), cfg.OutFile)

	// After appending, enforce the rolling window rotation
	return rotateContextFile(cfg)
}

func getFileDiff(file string) string {
	cmd := exec.Command("git", "diff", "HEAD", "--", file)
	var out bytes.Buffer
	cmd.Stdout = &out
	_ = cmd.Run() // Ignore errors (e.g., if it's untracked or no git repo)

	diff := out.String()
	if diff == "" {
		// Fallback for new/untracked files
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("- %s (File deleted or unreadable)\n", file).Error()
		}
		truncated := string(content)
		if len(truncated) > 3000 {
			truncated = truncated[:3000] + "\n...[Content truncated due to size]"
		}
		return fmt.Sprintf("- **%s** (Untracked/New File):\n```\n%s\n```\n", file, truncated)
	}
	return fmt.Sprintf("- **%s**:\n```diff\n%s\n```\n", file, diff)
}

// rotateContextFile ensures the context file doesn't exceed the configured MaxEvents
func rotateContextFile(cfg Config) error {
	content, err := os.ReadFile(cfg.OutFile)
	if err != nil {
		return err
	}

	strContent := string(content)
	eventDelimiter := "\n### Sync Event:"
	events := strings.Split(strContent, eventDelimiter)

	// The first split chunk is everything before the first event (our boilerplate header)
	// So len(events) - 1 is the number of actual sync events.
	if len(events)-1 > cfg.MaxEvents {
		fmt.Printf("Rotating context file (exceeded %d max events)\n", cfg.MaxEvents)
		
		newContent := events[0]
		// Slice off the oldest events, keeping only the last MaxEvents
		keepEvents := events[len(events)-cfg.MaxEvents:]
		for _, event := range keepEvents {
			newContent += eventDelimiter + event
		}

		err = os.WriteFile(cfg.OutFile, []byte(newContent), 0644)
		if err != nil {
			return fmt.Errorf("failed to rotate file: %w", err)
		}
	}
	return nil
}
