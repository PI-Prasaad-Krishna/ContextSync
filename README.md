# ContextSync

ContextSync is a local background daemon that maintains a dynamic .context.md memory bank for AI coding agents.

## Overview

AI IDEs currently perform expensive full-repository semantic scans on every prompt, which wastes tokens and causes the "Lost in the Middle" effect. ContextSync solves this by watching your project directory for file saves. It intercepts these saves, debounces them, and incrementally appends the list of changed files to a single, dense markdown file (.context.md). 

By pointing your AI agent to read this single file instead of the whole repository, it instantly knows exactly which files were modified, maximizing context efficiency.

## Features

- Zero Dependencies: Written in Go, it compiles to a single binary.
- Smart Watcher: Recursively monitors the working directory while dynamically respecting your `.gitignore` rules (using `go-gitignore`).
- Debouncer: Batches rapid file saves (e.g., "Save All") into a single event to prevent API/CPU overload.
- Sync Bridge: Appends timestamped logs of changed files to `.context.md` for AI consumption.

## Installation

Ensure you have Go installed, then build the binary:

```bash
go build
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

## License

This project is licensed under the MIT License. See the LICENSE file for details.
