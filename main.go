package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	watchCmd := flag.NewFlagSet("watch", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Expected 'init' or 'watch' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		initCmd.Parse(os.Args[2:])
		if err := handleInit(); err != nil {
			fmt.Printf("Error initializing context: %v\n", err)
			os.Exit(1)
		}
	case "watch":
		watchCmd.Parse(os.Args[2:])
		if err := handleWatch(); err != nil {
			fmt.Printf("Daemon error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Println("Expected 'init' or 'watch' subcommands")
		os.Exit(1)
	}
}
