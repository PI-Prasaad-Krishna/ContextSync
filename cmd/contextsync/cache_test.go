package main

import (
	"testing"
)

func TestFileCache(t *testing.T) {
	cache := NewFileCache()

	_, exists := cache.Get("test.txt")
	if exists {
		t.Errorf("Expected new cache to be empty")
	}

	cache.Set("test.txt", "hello world")
	content, exists := cache.Get("test.txt")
	if !exists || content != "hello world" {
		t.Errorf("Expected to retrieve cached content")
	}

	// Update the cache
	cache.Set("test.txt", "hello universe")
	content, _ = cache.Get("test.txt")
	if content != "hello universe" {
		t.Errorf("Expected cache to update")
	}
}
