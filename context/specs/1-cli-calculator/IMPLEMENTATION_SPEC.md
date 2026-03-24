# Implementation Spec: CLI Calculator with Expression Parsing and History

## Overview

Build a Go CLI calculator that evaluates mathematical expressions from stdin, supporting basic arithmetic operators (+, -, *, /), parentheses with correct operator precedence, and a session history feature. This provides a lightweight terminal-based calculator alternative to tools like `bc`.

## Source Issue

[Issue #1 — Build CLI calculator with expression parsing and history](../../..)

## Current State

Greenfield project — no Go code exists yet. The repository contains only the baf workflow scaffolding (context/, .claude/, .github/). A `go.mod` file and all source code must be created from scratch.

## Design

### Architecture

The calculator is a single-binary CLI tool with three internal packages:

```
cmd/calc/main.go          ← entry point, CLI argument routing
internal/parser/parser.go ← expression tokenizer + recursive descent parser
internal/history/history.go ← session history (file-backed)
```

### Expression Evaluation

Use a **recursive descent parser** implementing standard math precedence:

1. **Tokenizer** — splits input string into tokens: numbers (float64), operators (+, -, *, /), parentheses
2. **Parser** — recursive descent with two precedence levels:
   - `parseExpr`: handles `+` and `-` (lowest precedence)
   - `parseTerm`: handles `*` and `/`
   - `parseFactor`: handles numbers, unary minus, and parenthesized sub-expressions

This avoids external dependencies and keeps the implementation simple.

### CLI Interface

```bash
calc "2 + 3 * 4"          # → 14
calc "(2 + 3) * 4"         # → 20
calc history               # → last 100 expressions + results
calc clear                 # → clears history
calc --version             # → calc v0.1.0
```

- Single expression argument (quoted to prevent shell expansion of `*`)
- Reserved subcommands: `history`, `clear`
- Flags: `--version` prints version string (`calc v<version>`) and exits
- Version is set via a `Version` variable in `cmd/calc/main.go` (can be overridden at build time with `-ldflags`)
- Exit code 0 on success, 1 on error

### History Storage

- File-backed: `~/.calc_history` (one JSON line per entry)
- Each entry: `{"expr": "2 + 3", "result": 14, "time": "2026-03-24T..."}`
- `history` command shows last 100 entries, most recent last
- `clear` command truncates the file
- History file created on first write if it doesn't exist
- No locking needed — single-user CLI tool

### Error Handling

| Error | Message |
|-------|---------|
| Division by zero | `error: division by zero` |
| Malformed expression | `error: unexpected token '<token>' at position <N>` |
| Unmatched parentheses | `error: unmatched parenthesis` |
| Empty input | `error: no expression provided` |
| No args | Print usage to stderr |

All errors go to stderr with exit code 1. Results go to stdout.

## Configuration

No configuration files or environment variables. The history file location (`~/.calc_history`) is hardcoded but extracted as a constant for testability.

## File Plan

| File | Action | Purpose |
|------|--------|---------|
| `go.mod` | Create | Go module definition (`github.com/brandon/calc` or similar) |
| `cmd/calc/main.go` | Create | CLI entry point — argument parsing, subcommand routing |
| `internal/parser/parser.go` | Create | Tokenizer + recursive descent parser |
| `internal/parser/parser_test.go` | Create | Unit tests for parsing and evaluation |
| `internal/history/history.go` | Create | History read/write/clear operations |
| `internal/history/history_test.go` | Create | Unit tests for history (using temp files) |

## Implementation Order

1. **Initialize Go module** — `go mod init`, create directory structure
2. **Implement tokenizer** — `Tokenize(input string) ([]Token, error)` — split input into number/operator/paren tokens
3. **Implement parser** — `Eval(input string) (float64, error)` — recursive descent parser using tokenizer output; handles precedence, parentheses, unary minus
4. **Add parser tests** — cover: basic arithmetic, precedence, parentheses, nested parens, division by zero, malformed input, empty input
5. **Implement history** — `Store`, `Load`, `Clear` functions operating on `~/.calc_history`
6. **Add history tests** — use temp directory for isolation
7. **Wire up CLI** — `cmd/calc/main.go` routes args to eval or history subcommands, prints results/errors
8. **End-to-end manual test** — `go build ./cmd/calc && ./calc "2 + 3 * 4"`

## Testing

### Unit tests

- **Parser:** arithmetic correctness (all 4 operators), operator precedence, parentheses (including nested), division by zero error, malformed expression errors, whitespace handling, negative numbers, decimal numbers
- **History:** store and load entries, load returns last 100 max, clear removes all entries, handles missing file gracefully

### Integration

- Build the binary and test CLI argument routing
- Verify stdout/stderr separation
- Verify exit codes

### Running

```bash
go test ./...
```

## Not In Scope

- Variables or named constants
- Mathematical functions (sin, cos, sqrt, etc.)
- Graphing or visualization
- GUI
- Floating point precision beyond float64
- Multi-user or concurrent access to history
- Configuration files or environment variables
