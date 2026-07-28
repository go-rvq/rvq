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
- **UI integration tests live in `js/integration_tests/`** (e.g.
  `js/integration_tests/admin/presets`). Run them from the suite directory —
  that is where `bunfig.toml` sets the `preload` the DOM tests need:
  `cd js/integration_tests/admin/presets && bun test`. They boot the Go fixture
  server (in-memory SQLite) and drive it over HTTP; `*.dom.test.ts` also mount
  the real Vue app with happy-dom.
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

## Before committing

- **Run `make check` and only commit when it passes.** It is the gate: `gofmt -s`
  on the healthy trees, `go build` of the packages that compile, `go test` of the
  packages whose tests pass, and the bun UI integration suites. Long-broken
  upstream packages (RVQ docs/examples, pagebuilder, seo, media, libvips, and
  the `web` / `x/perm` tests) are deliberately out, so a failure means a
  regression you introduced.

## Development & Test Commands
Always run native Go tooling to verify compliance and correctness:

- **Full gate (use this)**: `make check`
- **Run all tests**: `go test ./...` (inclui pacotes quebrados de longa data)
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
* When changing `admin/packages/helper` or `admin/presets/gorm2op`, also build any
  downstream applications that consume them.

## Admin associations & list editor
- **Never assume an integer `id` or a single-column key.** Record keys can be any
  type (int, UUID, string) and **composite** (two or more columns). Read keys from
  the schema/relationship, not from hardcoded `id`/`<Field>ID`.
  - Use `model.ID` (composite-capable: `Values`/`Fields`, `SetTo`, `Related`,
    `String`) and `presets.ParseRecordID` (composite string parts joined by `_`).
  - For gorm associations, drive link/unlink through the **Association API**
    (`db.Model(parent).Association(field).Append/Delete`) and build filters from
    the relationship `References` — never with `id::BIGINT`-style raw SQL.
- **List editor per-item metadata** (flat reactive-form keys, see
  `admin/presets/listeditor.go`): `__present` (field submitted), `__index` (stable
  per-item form-key bracket, not the slice index), `__pos` (order), `__new`
  (browser-created), `__deleted`/`__purged` (client-side removal). Classify items
  by flags, never by primary key; a removed `__new` row vanishes.
  - `presets.PartitionListEditorItems(values, fieldFormKey, items)` groups the
    submitted slice (Deleted/New/Others/Kept, `__pos` order) or returns `nil` when
    not initialized. `gorm2op.SaveHasManyAssociation` uses it and supports
    composite PKs.
- Reference docs: `admin/packages/helper/{doc.go,README.md,docs/keys.md}`.

## Tasks
- **New tasks:** `todo.md`
- **Done tasks:** `todo-done.md`
- **When a task is DONE**, move it to `todo-done.md`.