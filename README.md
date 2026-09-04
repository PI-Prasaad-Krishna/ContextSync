# ContextSync

ContextSync is a local background daemon that maintains a dynamic .context.md memory bank for AI coding agents.

## Overview

AI IDEs currently perform expensive full-repository semantic scans on every prompt, which wastes tokens and causes the "Lost in the Middle" effect. ContextSync solves this by watching your project directory for file saves. It intercepts these saves, debounces them, and incrementally appends the list of changed files to a single, dense markdown file (.context.md). 

By pointing your AI agent to read this single file instead of the whole repository, it instantly knows exactly which files were modified, maximizing context efficiency.

## Architecture

[#architecture](#architecture)

ContextSync runs as a background daemon with four stages:

![Architecture Flow](architecture-flow.svg)

1. **Watcher** — recursively monitors the project directory for file save events, dynamically attaching to newly created subfolders and respecting `.gitignore` rules.
2. **Debouncer** — batches rapid successive saves (e.g. "Save All") into a single event within a configurable window (default 2s), preventing redundant processing.
3. **Diff Engine** — computes incremental string diffs in memory against the last known state of each changed file, producing only the changed lines rather than full file contents.
4. **Context Writer** — appends the diff as a timestamped entry to `.context.md`, prepending the static `Architecture.md` base context (if configured) so every write carries both long-term project understanding and the most recent changes. A rolling window rotation prunes the oldest entries once `max_events` is reached, keeping the file bounded.

Flow: `File Save → Watcher → Debouncer → Diff Engine → Context Writer → .context.md`

An AI coding agent is pointed at `.context.md` instead of the full repository, so every prompt carries only what actually changed.

## Features

- Zero Dependencies: Written in Go, it compiles to a single binary. No need for `git` or external libraries.
- Dynamic Watcher: Recursively monitors the working directory while dynamically attaching to newly created folders and respecting your `.gitignore` rules on the fly.
- Debouncer: Batches rapid file saves (e.g., "Save All") into a single event to prevent CPU overload.
- Incremental Native Diffs: Natively calculates string diffs in memory to output perfectly accurate, token-efficient patch logs of only the lines you changed!
- Rolling Window Rotation: Automatically prunes the oldest sync events from the context file when it reaches your configured limit, preventing token bloat.
- JSON Configuration: Configure all behaviors in a lightweight `.contextsync.json` file to avoid typing long commands.

## The Cost Problem (Measured Benchmark)

[#the-cost-problem-measured-benchmark](#the-cost-problem-measured-benchmark)

*Measured on the ContextSync repo itself, 26 files, 3-file editing session:*

- **Full-repo AI scan:** 9,875 tokens per prompt — processing takes longer and runs a higher risk of "Lost in the Middle" hallucination from 23 untouched files.
- **With ContextSync:** a session touching 3 files produces a `.context.md` diff of only 1,014 tokens.

**Result:** An **89.73%** reduction in tokens sent per prompt.

This is the gap ContextSync is designed to close: pointing agents at exactly what changed, not the whole repository.

## Best Practices (Hybrid Memory)

Because ContextSync only records your most *recent* file diffs, the AI will know exactly what you just did, but it might lose sight of the big picture. 

To solve this, ContextSync features a native **Base Context** pipeline. 
1. **Long-Term Memory:** Write a static `Architecture.md` file that explains your overall project goals, structs, and patterns.
2. **Configuration:** Add `"base_context_file": "Architecture.md"` to your `.contextsync.json`.
3. **The Magic:** On every file save, ContextSync dynamically copies your entire `Architecture.md` and glues it to the very top of `.context.md`, right above the rolling diffs.

You just hand the AI `.context.md`, and it instantly has both a high-level understanding of the whole project AND a perfect timeline of your recent code edits!

## Installation

**Option 1: Download Pre-compiled Binary (Recommended)**
Head over to the [Releases](https://github.com/PI-Prasaad-Krishna/ContextSync/releases) page and download the native binary for your operating system (Windows, macOS, or Linux).

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
  "base_context_file": "Architecture.md",
  "debug": false
}
```
*CLI equivalents: `--debounce`, `--max-events`, `--out`, `--base-context`, `--debug`*
## License

This project is licensed under the MIT License. See the LICENSE file for details.
