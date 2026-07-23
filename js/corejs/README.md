# corejs

The Vue 3 + TypeScript front-end runtime of the `go-rvq/rvq` admin. The Go server
renders HTML templates; **corejs** compiles them into live Vue components and
provides the client ↔ server bridge (**Plaid**), so most interactions are
server-driven without hand-written JavaScript.

## What it does

- **Bootstraps the SPA** — mounts a `Root` component on `#app`, seeds it with the
  server's initial template, and registers any extra components exposed by the
  server (`window.__goplaidVueComponentRegisters`).
- **Compiles server templates** — `componentByTemplate` turns an HTML string sent
  by the server into a Vue component, injecting the shared reactive `form` and
  `locals`.
- **Plaid** — a fluent request builder (`plaid()`) that POSTs an *event function*
  to the server, receives an `EventResponse`, and applies it: update/reload
  named **portals**, `pushState` the URL, run a returned script, or replace the
  root template.
- **Reactive state & components** — provides `form` / `locals` / `closer` /
  `vars`, and registers `GoPlaidScope`, `GoPlaidPortal`, `UserComponent`,
  `GoPlaidRunScript`, plus the `v-keep-scroll` / `v-assign` directives.

## Documentation

Detailed docs live under [`docs/`](docs/):

- [Architecture](docs/architecture.md) — bootstrap, the `Root` providers, and the server-driven template model.
- [Plaid](docs/plaid.md) — the `plaid()` request builder and how `EventResponse` is applied.
- [Components & directives](docs/components.md) — GoPlaidScope, GoPlaidPortal, UserComponent, GoPlaidRunScript.
- [Forms](docs/forms.md) — the reactive `form`/`locals`, field binding, and `objectToFormData`.

## Build

The build is driven by [Bun](https://bun.sh) (see `Makefile`):

```sh
make build          # bun install && bun run build && cp assets/* dist
# or directly:
bun run build       # type-check (vue-tsc) + vite build -> dist/index.js
```

`dist/index.js` is the bundle embedded/served by the Go side.

## Development

```sh
bun install         # (or pnpm install)
bun run dev         # vite dev server with HMR
bun run test:unit   # vitest
bun run type-check  # vue-tsc --build
bun run lint        # eslint --fix
bun run format      # prettier
```

### IDE

VSCode + [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar)
(disable Vetur). `.vue` type-checking uses `vue-tsc` instead of `tsc`.
