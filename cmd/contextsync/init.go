package main

import (
	"fmt"
	"os"
)

// handleInit creates a boilerplate context file if it doesn't already exist.
func handleInit(cfg Config) error {
	if _, err := os.Stat(cfg.OutFile); err == nil {
		return fmt.Errorf("context file %s already exists", cfg.OutFile)
	}

	boilerplate := getBoilerplate(cfg)

	err := os.WriteFile(cfg.OutFile, []byte(boilerplate), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("Successfully created %s\n", cfg.OutFile)
	return nil
}
