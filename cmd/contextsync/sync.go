package main

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// syncBatchToContext is the MVP "Sync Bridge".
// It takes a batch of changed file paths, calculates native incremental diffs, and appends it to .context.md.
func syncBatchToContext(cfg Config, cache *FileCache, files map[string]struct{}) error {
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
		diff := getFileDiff(cache, file)
		if diff == "" {
			continue // skip files with no meaningful text changes
		}
		if _, err := f.WriteString(diff); err != nil {
			f.Close()
			return fmt.Errorf("failed to write file entry: %w", err)
		}
	}
	f.Close()

	slog.Info("Successfully synced files", "count", len(files), "out", cfg.OutFile)

	// After appending, enforce the rolling window rotation
	return rotateContextFile(cfg)
}

func getFileDiff(cache *FileCache, file string) string {
	contentBytes, err := os.ReadFile(file)
	if err != nil {
		// Also remove from cache if deleted
		cache.mu.Lock()
		delete(cache.files, file)
		cache.mu.Unlock()
		return fmt.Sprintf("- **%s** (File deleted or unreadable)\n", file)
	}
	
	content := string(contentBytes)
	oldContent, exists := cache.Get(file)
	
	// Update cache immediately
	cache.Set(file, content)

	if !exists {
		// Fallback for new/untracked files we haven't cached yet
		truncated := content
		if len(truncated) > 3000 {
			truncated = truncated[:3000] + "\n...[Content truncated due to size]"
		}
		return fmt.Sprintf("- **%s** (Initial Cache):\n```\n%s\n```\n", file, truncated)
	}

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(oldContent, content, false)
	dmp.DiffCleanupSemantic(diffs)

	patches := dmp.PatchMake(oldContent, diffs)
	patchText := dmp.PatchToText(patches)

	if patchText == "" {
		return ""
	}

	// diffmatchpatch url-encodes strings (like %0A for newline). Decode them for readability.
	decodedPatch, err := url.QueryUnescape(patchText)
	if err == nil {
		patchText = decodedPatch
	}

	// Clean up patch text for better markdown formatting
	patchText = strings.TrimSpace(patchText)
	return fmt.Sprintf("- **%s**:\n```diff\n%s\n```\n", file, patchText)
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
		slog.Info("Rotating context file", "max_events", cfg.MaxEvents)
		
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
