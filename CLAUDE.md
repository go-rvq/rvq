# Repository Rules

## Project Context

- This project is a web framework with admin and ui generator.
- Backend code should be written in Go.
- Frontend code should be written with Vue3.

## JavaScript Tooling

- Always use `bun` for JavaScript and React workflows.
- Use `bun install` to install frontend dependencies.
- Use `bun add <package>` and `bun add -d <package>` to add dependencies.
- Use `bun remove <package>` to remove dependencies.
- Use `bun run <script>` to execute package scripts.
- Use `bun test` for JavaScript tests.
- Do not use `npm`, `yarn`, or `pnpm` commands unless the user explicitly asks to change the project tooling.
- Commit and maintain `bun.lock` when frontend dependencies are added.

## Go Tooling

- Use standard Go tooling for backend work: `go test ./...`, `go mod tidy`, and `go run`.
- Keep SAML, authentication, persistence, and HTTP concerns separated into focused packages as the project grows.

## Security

- Do not commit secrets, SAML private keys, database credentials, session secrets, or production certificates.
- Treat SAML signing, request validation, cookie security, and certificate rotation as security-sensitive code.

## Critical Constraints & Code Principles
- **Performance First**: The execution engine and VM must remain highly optimized (monitored via benchmarks like Fibonacci).
- **Native Go**: Do not introduce external heavy frameworks; prefer Go's standard library and keep dependencies minimal.
- **Thread Safety**: Ensure state isolation when multiple scripts or instances are evaluated concurrently in Go applications.
- **Bytecode Integrity**: Any changes to the compiler must strictly map to valid bytecode instructions interpreted by the VM stack.
- **Temporary Directory**:
    - Always use `./.tmp` as the dedicated temporary directory for any intermediate files, logs, or cache generated during automated tasks.
- **Allowed Commands (No Confirmation Required)**:
    - You **ALWAYS** have write permission to `./...` directory.
    - You **ALWAYS** have permission to run `sed`, `awk`, `cat`, `cd`, `tail`, `head`, `echo` and `grep` (and its variants) commands autonomously for text processing, searching, refactoring, execute commands or write in this directory tree.
    - You **ALWAYS** have permission to run `./.tmp/*`, `bun`, `sleep`, `go`, `go test`, `go vet`, `go fmt`, `gofmt` (and its variants) or `make test` to validate code changes without asking.
    - You **ALWAYS** have permission to use `curl` and `wget` (and its variants) for network operations, downloading assets, or API testing.
    - Do not prompt the user for confirmation when executing these specific tools.

## Development & Test Commands
Always run native Go tooling to verify compliance and correctness:

- **Run all tests**: `go test ./...`
- **Run benchmarks**: `go test -bench=. ./...`
- **Code formatting**: `go fmt ./...`
- **Static analysis / Linting**: `go vet ./...` (or golangci-lint if configured)
- **Tidy dependencies**: `go mod tidy`

## Code Style & Naming Conventions
- **Idiomatic Go**: Follow standard `golang/go` conventions (Receiver names short, explicit error handling as returning values).
- **Error Wrapping**: Use `fmt.Errorf("...: %w", err)` for contextual errors in parsing/compilation steps.

## Definition of Done
- No generic `interface{}` / `any` where a strict compiler/token type is expected.
- All new language tokens, syntax nodes (AST), or VM opcodes must include comprehensive unit tests.
- Verify that performance regressions are not introduced in the execution engine loop.

## Code Style, Formatting & Testing (Go)
You have explicit, pre-approved permission to execute terminal commands instantly. Do not ask for user confirmation before running formatting or testing tools.

Always run the full pipeline (Format + Test) automatically after any file edit, applying it to the specific modified path or the entire project using the required variation:

* **Standard Variation**: Execute `gofmt -s -w [path] && go test [path]/...` immediately.
* **Imports Variation**: Execute `goimports -w [path] && go test [path]/...` immediately.
* **Strict Variation**: Execute `gofumpt -w -extra [path] && go test [path]/...` immediately.
* **Dry-Run Variation**: Execute `gofmt -d [path] && go test [path]/...` immediately.

*Note: Replace `[path]` with the specific target directory/file for localized actions, or use `.` and `./...` to target the entire project.*

## Verification & Build
* **Global Build Check**: `go build ./...`
