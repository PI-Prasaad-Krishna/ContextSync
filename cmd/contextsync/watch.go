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

// handleWatch starts a background daemon that monitors the directory for changes.
func handleWatch(cfg Config) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Close()

	ignoreParser, err := ignore.CompileIgnoreFile(".gitignore")
	if err != nil {
		ignoreParser = ignore.CompileIgnoreLines()
	}

	go func() {
		changedFiles := make(map[string]struct{})

		timer := time.NewTimer(cfg.DebounceDuration)
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// Ignore our own context file
				if event.Name == cfg.OutFile || event.Name == ".\\"+cfg.OutFile {
					continue
				}

				cleanPath := strings.TrimPrefix(event.Name, ".\\")
				cleanPath = strings.TrimPrefix(cleanPath, "./")
				cleanPath = filepath.ToSlash(cleanPath)

				info, err := os.Stat(event.Name)
				isDir := err == nil && info.IsDir()
				
				if isDir {
					cleanPath += "/"
				}

				if ignoreParser.MatchesPath(cleanPath) {
					continue
				}

				if strings.Contains(event.Name, ".git") {
					continue
				}

				if event.Op&fsnotify.Create != 0 && isDir {
					// Dynamic Directory Watching: newly created directories
					err = filepath.Walk(event.Name, func(path string, walkInfo os.FileInfo, walkErr error) error {
						if walkErr != nil {
							return walkErr
						}
						wClean := filepath.ToSlash(path)
						if walkInfo.IsDir() {
							if ignoreParser.MatchesPath(wClean+"/") || walkInfo.Name() == ".git" {
								return filepath.SkipDir
							}
							return watcher.Add(path)
						}
						return nil
					})
					if err != nil {
						log.Println("Error dynamically adding directory:", err)
					}
				}

				if event.Op&fsnotify.Write != 0 || event.Op&fsnotify.Create != 0 {
					if !isDir { // Only batch actual files
						changedFiles[event.Name] = struct{}{}
						timer.Reset(cfg.DebounceDuration)
					}
				}

			case <-timer.C:
				if len(changedFiles) > 0 {
					if err := syncBatchToContext(cfg, changedFiles); err != nil {
						log.Printf("Error syncing batch: %v\n", err)
					}
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

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		cleanPath := filepath.ToSlash(path)
		if info.IsDir() && ignoreParser.MatchesPath(cleanPath+"/") {
			return filepath.SkipDir
		}
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to add directories to watcher: %w", err)
	}

	fmt.Println("ContextSync daemon is now watching the current directory...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	fmt.Printf("\nReceived signal (%v). Shutting down ContextSync gracefully...\n", sig)

	return nil
}
