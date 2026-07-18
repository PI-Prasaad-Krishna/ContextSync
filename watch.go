package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
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

				// File Filtering: Ignore .git directory and files inside it
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
					fmt.Printf("\n--- Processing batch of %d files ---\n", len(changedFiles))
					for file := range changedFiles {
						fmt.Println("-", file)
					}
					fmt.Println("------------------------------------")

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

	// Watch the current directory "."
	err = watcher.Add(".")
	if err != nil {
		return fmt.Errorf("failed to add directory to watcher: %w", err)
	}

	fmt.Println("ContextSync daemon is now watching the current directory...")

	// Block the main thread indefinitely so the program doesn't exit.
	<-make(chan struct{})

	return nil
}
