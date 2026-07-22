package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	ignore "github.com/sabhiram/go-gitignore"
)

const debounceDuration = 2 * time.Second

// handleWatch starts a background daemon that monitors the directory for changes.
func handleWatch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	// Ensure the watcher is closed when the function exits
	defer watcher.Close()

	// Compile .gitignore rules. If the file doesn't exist, we fall back to an empty parser.
	ignoreParser, err := ignore.CompileIgnoreFile(".gitignore")
	if err != nil {
		ignoreParser = ignore.CompileIgnoreLines()
	}

	// The watcher emits events on channels. We spin up a goroutine
	// to process these events concurrently so we don't block the main thread.
	go func() {
		// Track unique files changed within the debounce window
		changedFiles := make(map[string]struct{})

		// Create a timer for debouncing but stop it immediately.
		// We'll reset it every time a file changes.
		timer := time.NewTimer(debounceDuration)
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}

		for {
			// select blocks until one of its cases can proceed.
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// Step 3 edge case: The Infinite Loop.
				// Explicitly ignore events for our own context file.
				if event.Name == contextFileName || event.Name == ".\\"+contextFileName {
					continue
				}

				// Clean the path (remove leading .\ if present) for the gitignore parser
				cleanPath := strings.TrimPrefix(event.Name, ".\\")
				cleanPath = strings.TrimPrefix(cleanPath, "./")
				cleanPath = filepath.ToSlash(cleanPath)

				// If it's a directory, append a slash so rules like `ignore_me/` work
				if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
					cleanPath += "/"
				}

				// File Filtering: Ignore based on .gitignore rules
				if ignoreParser.MatchesPath(cleanPath) {
					continue
				}

				// Always ignore .git directory explicitly to avoid noise
				if strings.Contains(event.Name, ".git") {
					continue
				}

				if event.Op&fsnotify.Write != 0 || event.Op&fsnotify.Create != 0 {
					// Add file to our batch (map automatically handles uniqueness)
					changedFiles[event.Name] = struct{}{}

					// Reset the debounce timer. If rapid saves happen, this keeps pushing the timer back.
					timer.Reset(debounceDuration)
				}

			case <-timer.C:
				// The timer fired because 2 seconds passed since the last save.
				// We now process the entire accumulated batch safely.
				if len(changedFiles) > 0 {
					if err := syncBatchToContext(changedFiles); err != nil {
						log.Printf("Error syncing batch: %v\n", err)
					}

					// Clear the map for the next batch
					changedFiles = make(map[string]struct{})
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Watcher error:", err)
			}
		}
	}()

	// Recursively watch the current directory and all subdirectories, respecting .gitignore
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Clean path for gitignore parser
		cleanPath := filepath.ToSlash(path)

		// If it's a directory and matches our .gitignore, skip it entirely
		if info.IsDir() && ignoreParser.MatchesPath(cleanPath) {
			return filepath.SkipDir
		}

		// Always skip .git to avoid noise
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		// Add valid directories to the watcher
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to add directories to watcher: %w", err)
	}

	fmt.Println("ContextSync daemon is now watching the current directory...")

	// Step 6: Graceful Shutdown
	// Create a channel to listen for OS interrupt signals (like Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Block the main thread until we receive a signal
	sig := <-sigChan
	fmt.Printf("\nReceived signal (%v). Shutting down ContextSync gracefully...\n", sig)

	// Returning here triggers the `defer watcher.Close()` at the top of the function.
	// Closing the watcher closes its Events and Errors channels, which safely terminates
	// our background goroutine.
	return nil
}
