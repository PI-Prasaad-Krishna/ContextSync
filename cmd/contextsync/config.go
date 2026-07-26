package main

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds the configuration parsed from CLI flags or JSON
type Config struct {
	DebounceDuration time.Duration `json:"-"`
	DebounceString   string        `json:"debounce"`
	MaxEvents        int           `json:"max_events"`
	OutFile          string        `json:"out"`
	Debug            bool          `json:"debug"`
	BaseContextFile  string        `json:"base_context_file"`
}

// loadConfig tries to read a .contextsync.json file and apply it to the config
func loadConfig(cfg *Config) {
	data, err := os.ReadFile(".contextsync.json")
	if err == nil {
		_ = json.Unmarshal(data, cfg)
		if cfg.DebounceString != "" {
			d, err := time.ParseDuration(cfg.DebounceString)
			if err == nil {
				cfg.DebounceDuration = d
			}
		}
	}
}
