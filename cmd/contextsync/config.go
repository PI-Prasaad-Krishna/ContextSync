package main

import "time"

// Config holds the configuration parsed from CLI flags
type Config struct {
	DebounceDuration time.Duration
	MaxEvents        int
	OutFile          string
}
