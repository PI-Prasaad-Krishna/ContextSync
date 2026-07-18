Project: ContextSync (AI Memory Daemon)

🤖 Note to AI Agent (Antigravity/Cursor/Windsurf)

You are acting as a senior Go developer.

Write heavily commented (proper documentation), code explaining Go-specific concepts (goroutines, channels, pointers) enough to give idea to any engineer who will look at it.

Adhere strictly to the roadmap. Do not overcomplicate steps or skip ahead.

This file is your absolute source of truth for the project state and architecture.

🎯 Active Objective

Build a lightweight, zero-dependency local CLI tool in Go that maintains a dynamic .context.md memory bank for AI coding agents.

The Problem: AI IDEs currently perform expensive, slow "full-repository semantic scans" on every prompt, which wastes tokens and causes hallucination (the "Lost in the Middle" effect).
The Solution: A local background daemon that watches file saves and incrementally updates a single, highly dense markdown file. The AI agent reads this single file instead of the whole repo, dropping Time-To-First-Token (TTFT) and maximizing accuracy.

💻 Tech Stack

Language: Go (Golang) - chosen for zero end-user dependencies and cross-platform compilation into a single binary.

CLI Framework: Standard Go flag package.

File Watching: github.com/fsnotify/fsnotify for background OS-level file watching.

🏗️ Core Architecture & Lifecycle

Init (ctx init): Scaffolds a structured .context.md file in the user's project root.

Watch (ctx watch): A persistent background daemon watching the working directory.

Process & Sync: Intercepts file saves, debounces them, and (eventually) uses a fast LLM API to summarize the diff, appending it to .context.md.

Agent Consumption: The IDE agent reads the updated .context.md to instantly understand the project state.

🛑 Known Edge Cases to Solve (The "Gotchas")

As we build the file-watching logic, we must account for these specific traps:

The Infinite Loop: The watcher MUST explicitly ignore the .context.md file itself. If it watches the file it writes to, it will trigger an infinite loop of updates and crash the system.

The "Spam Save" Problem (Debouncing): Developers spam Ctrl+S. We must implement a debounce mechanism (e.g., waiting 2-3 seconds after the last save) before triggering the processing logic to prevent CPU/API overload.

Debouncer Batching: When multiple different files are changed within the debounce window (like a "Save All"), the debouncer must accumulate these unique file paths and hand them off as a single batch to the Sync Bridge, rather than dropping any changes.

API Modularity: The actual "smart summarization" will eventually be an external API call to a fast model (like Gemini 1.5 Flash). For the MVP, we just write a "dumb" local log to prove the pipes work.

File Filtering: Beyond ignoring `.context.md`, we must ignore `.git/` directories and ideally respect `.gitignore` to avoid processing noise from compiled binaries or third-party packages (like `node_modules`).

Graceful Shutdown: When the daemon is stopped (e.g., Ctrl+C), we must catch the OS interrupt signal (SIGINT) to cleanly stop the watcher and safely close any open file handles.

🗺️ Implementation Roadmap

Step 1: Scaffolding. Initialize go mod, set up main.go, and wire up the flag package to accept init and watch commands.

Step 2: Init Command. Write a function to generate a boilerplate .context.md file if it doesn't already exist.

Step 3: The Watcher (Dumb). Implement fsnotify to listen for file changes in the current directory and simply fmt.Println the changed file names.

Step 4: The Debouncer. Add a timer/channel system to debounce the fsnotify events so multiple rapid saves only trigger one event.

Step 5: The Sync Bridge. Connect the debounced watcher to actually write a dummy string into the .context.md file.

Step 6: Graceful Shutdown. Wire up `os/signal` to listen for interrupts and ensure the watcher and debouncer goroutines shut down cleanly.