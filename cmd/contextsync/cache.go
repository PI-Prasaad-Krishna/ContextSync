package main

import (
	"sync"
)

// FileCache stores the last known state of files in memory.
// This allows us to calculate incremental diffs natively without relying on Git.
type FileCache struct {
	mu    sync.RWMutex
	files map[string]string
}

// NewFileCache initializes a new thread-safe file cache.
func NewFileCache() *FileCache {
	return &FileCache{
		files: make(map[string]string),
	}
}

// Get retrieves the last known content of a file.
func (c *FileCache) Get(path string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	content, exists := c.files[path]
	return content, exists
}

// Set updates the known content of a file.
func (c *FileCache) Set(path, content string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.files[path] = content
}
