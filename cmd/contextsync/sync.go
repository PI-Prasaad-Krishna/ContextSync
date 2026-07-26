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

// getBoilerplate generates the static top half of the context file.
func getBoilerplate(cfg Config) string {
	base := "# Project Context\n\nThis file is automatically maintained by ContextSync. \nIt contains a dense, continuously updated summary of recent file changes to provide context to AI coding agents.\n"
	
	if cfg.BaseContextFile != "" {
		content, err := os.ReadFile(cfg.BaseContextFile)
		if err == nil {
			base += "\n## Architecture & Rules\n" + string(content) + "\n"
		} else {
			base += fmt.Sprintf("\n> [!WARNING]\n> Could not read base_context_file: %s\n", cfg.BaseContextFile)
		}
	}
	
	base += "\n## Recent Changes\n"
	return base
}

// syncBatchToContext is the MVP "Sync Bridge".
// It takes a batch of changed file paths, calculates native incremental diffs, and appends it to .context.md.
func syncBatchToContext(cfg Config, cache *FileCache, files map[string]struct{}) error {
	if len(files) == 0 {
		return nil
	}

	timestamp := time.Now().Format(time.RFC1123)
	newEventBody := fmt.Sprintf("%s\n", timestamp)
	
	hasChanges := false
	for file := range files {
		diff := getFileDiff(cache, file)
		if diff == "" {
			continue // skip files with no meaningful text changes
		}
		hasChanges = true
		newEventBody += diff
	}

	if !hasChanges {
		return nil
	}

	// Extract existing events
	var events []string
	contentBytes, err := os.ReadFile(cfg.OutFile)
	if err == nil {
		parts := strings.Split(string(contentBytes), "\n### Sync Event: ")
		if len(parts) > 1 {
			events = parts[1:]
		}
	}

	// Append the new event
	events = append(events, newEventBody)

	// Enforce rolling window
	if len(events) > cfg.MaxEvents {
		slog.Info("Rotating context file", "max_events", cfg.MaxEvents)
		events = events[len(events)-cfg.MaxEvents:]
	}

	// Generate fresh boilerplate and stitch together
	finalOutput := getBoilerplate(cfg)
	for _, event := range events {
		finalOutput += "\n### Sync Event: " + event
	}

	err = os.WriteFile(cfg.OutFile, []byte(finalOutput), 0644)
	if err == nil {
		slog.Info("Successfully synced files", "count", len(files), "out", cfg.OutFile)
	}
	return err
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
