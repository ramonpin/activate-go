# Agent Guidelines for activate-go

## Build & Test Commands

- Build: `go build` or `go build -v -o activate-go`
- Run tests: `go test -v`
- Run single test: `go test -v -run TestName` (e.g.,
  `go test -v -run TestIsVenv`)
- Format: `go fmt ./...`
- Lint: `go vet ./...`

## Code Style

- **Imports**: Standard library first, then third-party (grouped by blank line)
- **Naming**: camelCase for private, PascalCase for public; descriptive names
- **Functions**: Small, focused functions; verb-based names (e.g., `isVenv`,
  `findVenvsInDir`)
- **Error Handling**: Always check errors; use `fmt.Fprintf(os.Stderr, ...)` for
  errors; exit with `os.Exit(1)` on failure
- **Comments**: Spanish is acceptable (codebase uses Spanish comments); document
  exported functions
- **Types**: Use Go's type inference where clear; explicit types for function
  signatures
- **Formatting**: Use `go fmt`; tabs for indentation

## Architecture

- Uses Bubble Tea for TUI; model-update-view pattern
- Main logic: search current/parent dirs for venvs (with `bin/activate` file)
- Outputs shell command to stderr for UX, activation command to stdout for eval
