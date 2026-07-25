package main

import (
	"flag"
	"log/slog"
	"os"
	"time"
)

func main() {
	var cfg Config
	
	// Defaults
	cfg.DebounceDuration = 2 * time.Second
	cfg.MaxEvents = 15
	cfg.OutFile = ".context.md"
	cfg.Debug = false

	// Attempt to load from JSON
	loadConfig(&cfg)

	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	initCmd.StringVar(&cfg.OutFile, "out", cfg.OutFile, "Output file for context")

	watchCmd := flag.NewFlagSet("watch", flag.ExitOnError)
	watchCmd.DurationVar(&cfg.DebounceDuration, "debounce", cfg.DebounceDuration, "Debounce duration for rapid saves")
	watchCmd.IntVar(&cfg.MaxEvents, "max-events", cfg.MaxEvents, "Maximum number of sync events to keep in the context file")
	watchCmd.StringVar(&cfg.OutFile, "out", cfg.OutFile, "Output file for context")
	watchCmd.BoolVar(&cfg.Debug, "debug", cfg.Debug, "Enable verbose debug logging")

	if len(os.Args) < 2 {
		slog.Error("Expected 'init' or 'watch' subcommands")
		os.Exit(1)
	}

	// Initialize Logger
	logLevel := slog.LevelInfo
	if cfg.Debug {
		logLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))

	switch os.Args[1] {
	case "init":
		initCmd.Parse(os.Args[2:])
		if err := handleInit(cfg); err != nil {
			slog.Error("Error initializing context", "err", err)
			os.Exit(1)
		}
	case "watch":
		watchCmd.Parse(os.Args[2:])
		cache := NewFileCache()
		if err := handleWatch(cfg, cache); err != nil {
			slog.Error("Daemon error", "err", err)
			os.Exit(1)
		}
	default:
		slog.Error("Expected 'init' or 'watch' subcommands")
		os.Exit(1)
	}
}
