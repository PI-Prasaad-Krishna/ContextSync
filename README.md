# ContextSync

ContextSync is a local background daemon that maintains a dynamic .context.md memory bank for AI coding agents.

## Overview

AI IDEs currently perform expensive full-repository semantic scans on every prompt, which wastes tokens and causes the "Lost in the Middle" effect. ContextSync solves this by watching your project directory for file saves. It intercepts these saves, debounces them, and incrementally appends the list of changed files to a single, dense markdown file (.context.md). 

By pointing your AI agent to read this single file instead of the whole repository, it instantly knows exactly which files were modified, maximizing context efficiency.

## Features

- Zero Dependencies: Written in Go, it compiles to a single binary. No need for `git` or external libraries.
- Dynamic Watcher: Recursively monitors the working directory while dynamically attaching to newly created folders and respecting your `.gitignore` rules on the fly.
- Debouncer: Batches rapid file saves (e.g., "Save All") into a single event to prevent CPU overload.
- Incremental Native Diffs: Natively calculates string diffs in memory to output perfectly accurate, token-efficient patch logs of only the lines you changed!
- Rolling Window Rotation: Automatically prunes the oldest sync events from the context file when it reaches your configured limit, preventing token bloat.
- JSON Configuration: Configure all behaviors in a lightweight `.contextsync.json` file to avoid typing long commands.

## Installation

**Option 1: Download Pre-compiled Binary (Recommended)**
Head over to the [Releases](https://github.com/) page and download the native binary for your operating system (Windows, macOS, or Linux).

**Option 2: Build from Source**
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
./contextsync watch
```

### Configuration
You can customize behavior using a `.contextsync.json` file in your project root, or via CLI flags (which take precedence):
```json
{
  "debounce": "2s",
  "max_events": 15,
  "out": ".context.md",
  "debug": false
}
```
*CLI equivalents: `--debounce`, `--max-events`, `--out`, `--debug`*
## License

This project is licensed under the MIT License. See the LICENSE file for details.
