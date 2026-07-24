package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	var cfg Config

	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	initCmd.StringVar(&cfg.OutFile, "out", ".context.md", "Output file for context")

	watchCmd := flag.NewFlagSet("watch", flag.ExitOnError)
	watchCmd.DurationVar(&cfg.DebounceDuration, "debounce", 2*time.Second, "Debounce duration for rapid saves")
	watchCmd.IntVar(&cfg.MaxEvents, "max-events", 15, "Maximum number of sync events to keep in the context file")
	watchCmd.StringVar(&cfg.OutFile, "out", ".context.md", "Output file for context")

	if len(os.Args) < 2 {
		fmt.Println("Expected 'init' or 'watch' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		initCmd.Parse(os.Args[2:])
		if err := handleInit(cfg); err != nil {
			fmt.Printf("Error initializing context: %v\n", err)
			os.Exit(1)
		}
	case "watch":
		watchCmd.Parse(os.Args[2:])
		if err := handleWatch(cfg); err != nil {
			fmt.Printf("Daemon error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Println("Expected 'init' or 'watch' subcommands")
		os.Exit(1)
	}
}
