# Closer — the control object of an overlay, and the address bar

Every overlay (dialog, drawer, form host) is driven by one reactive object, the
**closer**. Turning `closer.show` on opens it, turning it off closes and destroys
it. Around that single flag there are two callback lists:

```js
closer.openCallbacks   // fn(closer) — run when show goes false → true
closer.closeCallbacks  // fn(closer) — run when show goes true → false
```

and the sugar that registers into them:

```js
const off = closer.onOpen((c) => …)   // returns the remover
closer.onClose((c) => …)
```

Rules, all in [`js/corejs/src/closer.ts`](../../../js/corejs/src/closer.ts):

- callbacks fire on the **transition** only — assigning the value it already has
  does nothing;
- the state the initializer carried (`{show: true}`) is the starting point and
  fires **nothing**: nobody is subscribed at creation time;
- a callback receives the closer itself;
- a callback that throws is logged and the others still run.

The trigger is an accessor on the OBJECT (`show` over `_show`), not a component
`watch`, because the same closer is adopted by several components — a watcher per
component would fire the callbacks once per adopter.

## Where a closer comes from

| Origin | How |
| --- | --- |
| the app root | `createCloser()` in `app.ts` |
| a scope | `go-plaid-scope`'s `:closer` prop — `web.ScopeBuilder.CloserInit`, `web.CloserScope` |
| a form host | a `user-component` scope variable: `$closer({show:false, …})` (see [form-host.md](form-host.md)) |
| adoption | a scope that receives an existing closer (`:closer="closer"`, `ParamCloserProvided`) or a portal that seeds it seeds THE SAME object |

`$closer` is a global property registered by `plaidPlugin`, so it resolves in any
template expression. Adoption never re-wraps: `asCloser` is idempotent, so an
adopted closer keeps exactly one copy of its callbacks.

## The address bar

An overlay showing a listing, a record's detail, its edit form or a create form
shows what a PAGE shows. So a closer that carries a `url` puts that clean page
address in the address bar while it is open, and puts back the previous one when
it closes:

```go
host.URL("/admin/articles/1")                                   // literal
host.URLExpr(`(closer) => "/admin/articles/" + closer.id`)      // per row
```

The expression form exists because a listing renders ONE host for every row: the
record is only known when a row opens it (`closer.id`), so the address is a
function of the closer.

Overlays stack, so the addresses unwind **LIFO**
([`closer-url.ts`](../../../js/corejs/src/closer-url.ts)):

```
/admin/articles                     listing page
  → row opens the DETAIL            /admin/articles/1
      → its EDIT                    /admin/articles/1/edit
      ← closed                      /admin/articles/1
  ← closed                          /admin/articles
```

Each entry records the address that was current when it opened; closing an entry
restores ITS address and drops everything above it (those overlays live inside
the one being closed, so they are gone anyway). Opening uses `pushState`, closing
`replaceState` — the overlay's address is not a place to come back to. The state
is `null` on purpose: plaid's popstate handler only reacts to a state it pushed
itself, so moving in this history must not reload the page.

**Back closes the top overlay.** A popstate with overlays open turns the innermost
one off (with the history writes suppressed — the browser has already moved), so
the address bar never shows something that is no longer on screen.

Wired out of the box, all pointing at routes that really exist:

| Overlay | Address |
| --- | --- |
| a row's DETAIL | `/{listing}/{id}` |
| a row's EDIT (or the row itself, with no detailing) | `/{listing}/{id}/edit` |
| NEW, from a listing | `/{listing}/new` |
| EDIT of a detailing | `/{listing}/{id}/edit` |

`/{listing}/new` is served by `EditingBuilder.GetCreatingPageFunc`
(`SetCreatingPageFunc` to replace it) — reloading or sharing the address opens
the same create form as a page.

## Go API

```go
// on a form host
host.URL(url)            // literal address
host.URLExpr(js)         // a function of the closer
host.OnOpen(js)          // appended to openCallbacks; the closer is `closer`
host.OnClose(js)

// on any scope that creates a closer
web.CloserScope(comp, true).
    CloserURL(`"/admin/articles/1"`).
    CloserOnOpen(`console.log("opened", closer)`).
    CloserOnClose(`presetsListing?.loader?.go()`)
```

## Testing it

- [`js/integration_tests/admin/presets/closer.test.ts`](../../../js/integration_tests/admin/presets/closer.test.ts) —
  the contract: transitions, idempotent adoption, error isolation, the LIFO
  address stack, the dynamic (per-row) address and Back.
- [`closer.dom.test.ts`](../../../js/integration_tests/admin/presets/closer.dom.test.ts) —
  the real runtime against the Go fixture: a row opens the detail, the detail
  opens its edit, and the addresses stack and unwind as the overlays do; the same
  addresses are then fetched to prove they serve a page.

> The runtime that a browser loads is the **built** bundle,
> `js/corejs/dist/index.js` (committed, embedded by `js/corejs.go`), while the
> bun suites import the corejs sources. After changing anything under
> `js/corejs/src`, run `cd js/corejs && bun run build-only` and commit the
> bundle. The build empties `dist/`, so restore the vendored files it does not
> produce (`vue.global.*.js`, `vue-i18n.js`) with `git checkout`.
