# Plaid

**Plaid** is the client ↔ server bridge: a fluent `Builder` (`builder.ts`) that
POSTs an *event function* to the server and applies the returned `EventResponse`.
Get one from the injected factory:

```ts
const plaid = inject('plaid') as () => Builder
plaid().eventFunc('doSomething').query('id', '42').go()
```

The Go side generates these calls (`web.Plaid()...Go()`), so most templates never
write Plaid by hand.

## Building a request

| Method | Effect |
| --- | --- |
| `eventFunc(id)` | the event function to invoke (`reload()` sets `__reload__`). |
| `url(v)` | POST to a specific URL (default: current window URL). |
| `query(k, v)` / `queries(obj)` | query params; `mergeQuery(true)` / `clearMergeQuery(keys)` control merging. |
| `location(loc)` / `stringQuery(s)` | set the full location / raw query string. |
| `form(v)` | the reactive form object serialized into the POST body. |
| `formData(v)` | extra data merged into a separate object, also serialized into the body. |
| `fieldValue(name, v)` | set a single form field. |
| `scope(v)` / `locals(v)` / `vars(v)` / `closer(v)` | pass scope state through. |
| `pushState(true)` / `pushStateURL(url)` | push a history entry / URL. |
| `method(m)` | HTTP method (default POST). |
| `skipFiles(true)` | drop file entries and send the body url-encoded (see below). |
| `parent(index, value)` | nested parent record ids. |
| `noCache()` | bypass caching. |
| `preFetch(fn)` / `postFetch(fn)` | hooks around the request. |
| `run(script)` | attach a client script to run. |

Terminals:

- `go()` — run the request and apply the response (the usual call).
- `fetch()` — just perform the request, returning the `EventResponse`.
- `buildFetchURL()` / `buildPushStateArgs()` — inspect what would be sent.
- `parseUrl(url)` — populate the builder from a URL (event + queries).

## The POST body

On a POST, Plaid serializes both `form` and `formData` into a single `FormData`
via `objectToFormData` (see [forms.md](forms.md)):

- **default** — the body is that `FormData` (multipart/form-data; file uploads
  are included).
- **`skipFiles(true)`** — file entries (non-string values) are removed and the
  body is re-encoded as `URLSearchParams` (application/x-www-form-urlencoded),
  i.e. a plain form post with no multipart overhead. Use it for events that never
  upload a file.

## Applying the `EventResponse`

`go()` handles each field of the response (`types.ts`):

| Field | Handling |
| --- | --- |
| `pageTitle` | sets `document.title`. |
| `redirectURL` | `document.location.replace`. |
| `reloadPortals` | re-fetch each named portal. |
| `updatePortals` | swap each named portal's template (`defer` waits for a not-yet-mounted portal). |
| `pushState` | push the history entry and reload the page location. |
| `runScript` | run the returned script. |
| `body` | replace the root template (unless `loadPortalBody`). |
| `reload` | reload the current view. |

`isFetching` is toggled via the `fetchStart` / `fetchEnd` window events around
the request, and `popstate` routes back/forward through `onpopstate`.
