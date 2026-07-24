# presets integration tests (bun)

End-to-end tests for the presets **EditForm/NewForm form-scope wrapper** and the
**list editor** (including four-level nesting), run with [bun](https://bun.sh).
They boot the Go test app (in-memory SQLite, seeded) as a real HTTP server and
drive its plaid event funcs over HTTP, asserting the exact response contract the
browser relies on — including the reactive `form` values (read from the v-assign
seeds each input binds to) at every nesting level.

## Layout

- `server/main.go` — a Go `main` that serves a fixture via `integration.ServeEnv`
  (prints `LISTENING http://<addr>`; `PORT=0` → ephemeral). The `APP` env var
  selects the fixture: `listeditor` (default) or `nested` (four levels).
- `server.ts` — boots the server (builds it once into `.tmp/`), bun-only.
- `helpers.ts` — call event funcs, parse plaid responses, follow the EditForm
  run-script, and read the reactive `form` from a rendered fragment
  (`formAssigns`).
- `editform.test.ts` — the flow from the Edit click: `presets_EditForm` returns
  the wrapper (closer scope + single `form` scope `{$parent: form}` + inner
  portal + run-script firing `presets_Edit`); following the run-script renders the
  real Edit dialog; an invalid save re-renders with the required error; a valid
  save tears the overlay down via `closer.show = false`.
- `nested.test.ts` — a four-level nested list editor (L0 → L1 → L2 → L3): opens
  with every level's `form` value from the DB; serializes the whole nested `form`
  on submit; a cleared required Name at the deepest level fails validation and the
  re-render preserves the other levels; an inclusion at a nested level is flagged
  `__new`; removing a deep row keeps it deleted on the validation re-render; and
  removing a first-level item persists.

The Go fixtures live in `admin/presets/integration` (`NewApp`/`NewSeededHandler`
and `NewNestedApp`/`NewNestedSeededHandler`), shared with the Go tests in
`admin/presets/tests/listeditor`.

## Run

```sh
cd js/integration_tests/admin/presets
bun test
```

Requires the Go toolchain on PATH (the server is built on first run). The binary
is cached at `<repo>/.tmp/presets_testserver`.

## Two kinds of test (both run under `bun test`)

- **HTTP contract** — `*.test.ts` (`editform`, `nested`, `composite`): drive the
  event funcs over HTTP and assert the plaid responses + the reactive `form`
  (read from the v-assign seeds). No DOM.
- **Mounted DOM** — `*.dom.test.ts` (`editform.dom`, `nested.dom`): mount the real
  corejs Root and prove the reactive flow from the click (closer → run-script →
  dialog), reading the mounted `form`. Enabled by `setup.ts` (bun preload):
  - a Vue-SFC loader (`@vue/compiler-sfc`) — bun has no real `.vue` loader;
  - happy-dom (`GlobalRegistrator`, same-origin disabled, ResizeObserver stub);
  - a light `vx-dialog` stub (the full vuetifyx is too heavy) so the fields
    render. The corejs `@/` alias resolves via `corejs/tsconfig.json` paths.

## Fixtures

Go fixtures in `admin/presets/integration`, selected by the server's `APP` env:
`""`/`listeditor` (one level), `nested` (four levels), `composite` (multi-level
COMPOSITE primary keys: Cart → Items PK(CartID,Sku) → Notes PK(CartID,Sku,Seq)).
