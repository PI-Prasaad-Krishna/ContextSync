package main

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
)

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
		for {
			// select blocks until one of its cases can proceed.
			// This is how we listen to multiple channels in Go.
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

				// For this MVP step, we just print the file that changed.
				// We filter by Write or Create operations to reduce noise.
				if event.Op&fsnotify.Write != 0 || event.Op&fsnotify.Create != 0 {
					fmt.Println("File changed:", event.Name)
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
	// In Step 6, we'll replace this with a proper os.Signal listener for graceful shutdown.
	<-make(chan struct{})

	return nil
}
