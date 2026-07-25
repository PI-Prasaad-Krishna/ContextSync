package main

import (
	"os"
	"strings"
	"testing"
)

func TestRotateContextFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "contextsync-test-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	cfg := Config{
		MaxEvents: 2,
		OutFile:   tmpFile.Name(),
	}

	// 3 events, but max is 2
	initialContent := "# Boilerplate\n\n### Sync Event: 1\nEvent 1\n### Sync Event: 2\nEvent 2\n### Sync Event: 3\nEvent 3\n"
	os.WriteFile(tmpFile.Name(), []byte(initialContent), 0644)

	err = rotateContextFile(cfg)
	if err != nil {
		t.Fatal(err)
	}

	result, _ := os.ReadFile(tmpFile.Name())
	resultStr := string(result)

	if strings.Contains(resultStr, "Event 1") {
		t.Errorf("Expected Event 1 to be pruned, got: %s", resultStr)
	}
	if !strings.Contains(resultStr, "Event 2") || !strings.Contains(resultStr, "Event 3") {
		t.Errorf("Expected Events 2 and 3 to remain, got: %s", resultStr)
	}
}
