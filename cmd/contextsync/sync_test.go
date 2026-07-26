package main

import (
	"os"
	"strings"
	"testing"
)

func TestSyncBatchToContext(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "contextsync-test-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	baseFile, err := os.CreateTemp("", "contextsync-base-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(baseFile.Name())

	os.WriteFile(baseFile.Name(), []byte("BASE ARCHITECTURE"), 0644)

	cfg := Config{
		MaxEvents:       2,
		OutFile:         tmpFile.Name(),
		BaseContextFile: baseFile.Name(),
	}

	// 3 previous events in the file
	initialContent := "BOILERPLATE\n\n### Sync Event: 1\nEvent 1\n### Sync Event: 2\nEvent 2\n### Sync Event: 3\nEvent 3\n"
	os.WriteFile(tmpFile.Name(), []byte(initialContent), 0644)

	cache := NewFileCache()
	
	// Create a dummy file to simulate a new save
	dummyFile, _ := os.CreateTemp("", "dummy-*.txt")
	os.WriteFile(dummyFile.Name(), []byte("new changes"), 0644)
	defer os.Remove(dummyFile.Name())

	files := map[string]struct{}{
		dummyFile.Name(): {},
	}

	err = syncBatchToContext(cfg, cache, files)
	if err != nil {
		t.Fatal(err)
	}

	result, _ := os.ReadFile(tmpFile.Name())
	resultStr := string(result)

	// Verify the base architecture is dynamically prepended!
	if !strings.Contains(resultStr, "BASE ARCHITECTURE") {
		t.Errorf("Expected base context to be dynamically prepended")
	}

	// We had 3 events, and added 1 new one. Max is 2. 
	// Event 1 and Event 2 should be deleted.
	if strings.Contains(resultStr, "Event 1") || strings.Contains(resultStr, "Event 2") {
		t.Errorf("Expected Event 1 and 2 to be pruned, got:\n%s", resultStr)
	}
	
	if !strings.Contains(resultStr, "Event 3") {
		t.Errorf("Expected Event 3 to remain")
	}
}
