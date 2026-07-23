# Components & directives

`plaidPlugin` (in `app.ts`) registers these globally, so server templates can use
them directly.

## GoPlaidPortal

A named region whose content the server can update or reload independently of the
rest of the page.

```html
<go-plaid-portal portal-name="detail">…</go-plaid-portal>
```

- Registers itself in `window.__goplaid.portals[name]` (anonymous portals get a
  generated name).
- `EventResponse.updatePortals` swaps a portal's template
  (`updatePortalTemplate(body, options)`); `reloadPortals` re-fetches it.
- A portal can carry its own `form` / `locals` scope and render options.

This is the primary mechanism for partial updates: an event returns only the
portals that changed.

## GoPlaidScope

Provides a local `locals` (and access to `form`) to its slot — the client side of
the Go `web.Scope`:

```html
<go-plaid-scope v-slot="{ locals }" :init='{ open: false }'>
  <v-btn @click="locals.open = !locals.open">toggle</v-btn>
</go-plaid-scope>
```

- `locals` is reactive and chains to the parent scope (`$parent`), so nested
  scopes can read outer locals.
- It also exposes `form`, letting scoped inputs bind into the shared form.

## UserComponent

A flexible wrapper that declares a **scope** and runs **setup** functions —
the client side of the Go `vue.UserComponent`. It powers form fields and the list
editor.

- `:scope="{ items: [...] }"` provides slot props (`v-slot="{ items, form }"`).
- `:setup="[fn]"` runs setup functions with `{ scope, $scope, ref, reactive,
  computed, inject, provide, watch, Vue, … }`, letting a field wire computed
  bindings against the ambient `form`.
- `assign` merges values into a target object on mount.

See the Go `web/vue` package and `admin/presets` field components for how these
are emitted.

## GoPlaidRunScript

Runs a script function on mount (and optional `setup` / `beforeUnmount` /
`unmounted` hooks), with a scope of `{ ...vue, view: window }`. Used to execute
the `runScript` returned in an `EventResponse`, or inline imperative snippets.

## Directives

| Directive | Purpose |
| --- | --- |
| `v-keep-scroll` | preserves an element's scroll position across re-renders. |
| `v-assign` | assigns/merges values into a reactive target on mount (`assign.ts`). |
| `GlobalEvents` | (from `vue-global-events`) attach window-level event listeners in templates. |
