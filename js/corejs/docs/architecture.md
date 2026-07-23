# Architecture

corejs is a small runtime around a **server-driven** rendering model: the Go
server sends HTML (with Vue directives); corejs compiles it into a live Vue
component and re-renders regions as the server responds to events.

## Bootstrap

`main.ts` mounts the app:

1. reads the initial HTML from `#app`'s `innerHTML`;
2. `createWebApp(template)` builds the Vue app around the `Root` component;
3. runs every registration in `window.__goplaidVueComponentRegisters` (extra
   components contributed by the server / plugins);
4. `app.mount('#app')`.

## The `Root` component (`app.ts`)

`Root` sets up the shared reactive state and `provide`s it to the whole tree:

| provide | Purpose |
| --- | --- |
| `form` | the reactive form object serialized on submit (see [forms.md](forms.md)) |
| `locals` | per-scope reactive UI state |
| `closer` | dialog/drawer close state |
| `vars` | app-wide vars (e.g. `__notification`) |
| `plaid` | a factory returning a fresh [Plaid](plaid.md) `Builder` wired to the root |
| `isFetching` | a ref toggled around requests (`fetchStart` / `fetchEnd` events) |
| `updateRootTemplate` | replaces the whole page component from a template string |

It also wires `popstate` to Plaid so browser back/forward re-fetches
(`plaid().onpopstate`).

`plaidPlugin` registers the global components and directives:
`GoPlaidScope`, `GoPlaidPortal`, `GoPlaidRunScript`, `UserComponent`,
`GlobalEvents`, and the `v-keep-scroll` / `v-assign` directives (see
[components.md](components.md)).

## Server-driven templates

`componentByTemplate(template, scope, ...)` (in `component-by-template.ts`)
compiles an HTML string into a Vue component whose render context includes the
shared `form` and `locals`. This is how the server "renders": it returns HTML,
and corejs turns it into a reactive component — for the whole page
(`updateRootTemplate`) or for a single [portal](components.md#go-plaid-portal).

## Request lifecycle

```
user action ──► plaid().eventFunc("evt").go()
                     │  POST FormData(form) to the event URL
                     ▼
              server EventResponse
                     │  applied by go():
                     ├─ updatePortals / reloadPortals   (swap regions)
                     ├─ pushState                       (URL + history)
                     ├─ runScript                       (imperative JS)
                     ├─ redirectURL / pageTitle / reload
                     └─ body → updateRootTemplate       (replace page)
```

See [plaid.md](plaid.md) for the full builder API and response handling.
