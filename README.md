# ContextSync

ContextSync is a local background daemon that maintains a dynamic .context.md memory bank for AI coding agents.

## Overview

AI IDEs currently perform expensive full-repository semantic scans on every prompt, which wastes tokens and causes the "Lost in the Middle" effect. ContextSync solves this by watching your project directory for file saves. It intercepts these saves, debounces them, and incrementally appends the list of changed files to a single, dense markdown file (.context.md). 

By pointing your AI agent to read this single file instead of the whole repository, it instantly knows exactly which files were modified, maximizing context efficiency.

## Features

- Zero Dependencies: Written in Go, it compiles to a single binary.
- Dynamic Watcher: Recursively monitors the working directory while dynamically attaching to newly created folders and respecting your `.gitignore` rules on the fly (using `go-gitignore`).
- Debouncer: Batches rapid file saves (e.g., "Save All") into a single event to prevent CPU overload.
- Super Context (Git Diffs): Executes `git diff` to extract exact code changes and appends beautifully formatted markdown diffs directly into `.context.md`.
- Rolling Window Rotation: Automatically prunes the oldest sync events from the context file when it reaches your configured limit, preventing amnesia wipes and token bloat!

## Installation

1. Clone the repository
2. Build the binary:
```bash
go build -o contextsync.exe ./cmd/contextsync
```

## Usage

Initialize a boilerplate context file in your project root:

```bash
./contextsync init
```

Start the file-watching daemon:
```bash
./contextsync.exe watch --debounce=2s --max-events=15 --out=.context.md
```

### CLI Flags
- `--debounce`: Time to wait after a file save before grouping them into a batch (default: `2s`)
- `--max-events`: The rolling window limit. Once this many events are in `.context.md`, the oldest events are pruned (default: `15`)
- `--out`: The name of the context file (default: `.context.md`)

## License

This project is licensed under the MIT License. See the LICENSE file for details.
