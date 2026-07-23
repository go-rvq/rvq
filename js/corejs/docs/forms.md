# Forms

## The reactive `form`

`Root` provides a single reactive `form` object (`provide('form', form)`). Field
components bind into it by **flat key**: an input for `Parcelas[0].Valor` writes
`form["Parcelas[0].Valor"]`. On submit, Plaid serializes the whole object into a
`FormData` and POSTs it — so the DOM inputs are the source of truth and the server
receives a conventional form payload.

`locals` is the sibling reactive object for transient UI state (dialog open flags,
etc.), scoped via [GoPlaidScope](components.md#go-plaid-scope).

## Serialization — `objectToFormData`

`objectToFormData(obj, formData, parentKey?)` (in `utils.ts`) walks `obj` and
writes entries into `formData`:

- nested objects/arrays recurse, building keys like `Parcelas[0].Valor`;
- a `$parent` key is skipped (it only links scopes, it is not data);
- **function values are called** and their result is serialized under the key —
  this lets a component compute what to submit for a key lazily, e.g. the list
  editor posts per-item metadata for a slice field;
- `File` / `Blob` values are set as-is (multipart upload).

`setFormValue(form, name, val)` writes a single value (handling arrays and
add/remove `ValueOp`s). `Plaid.skipFiles(true)` drops file entries and sends the
body url-encoded (`URLSearchParams`) instead of multipart — see
[plaid.md](plaid.md#the-post-body).

## Field binding pattern

The Go field components (via `web/vue`) typically:

1. seed the initial value: `form["<key>"] = <value>`;
2. bind the input with `v-model` on a computed that reads/writes `form["<key>"]`;

so editing the input updates `form`, and `objectToFormData` serializes it on the
next Plaid request. Because `form` is a single flat object, list/nested editors
compose by writing their own `form["field[i].sub"]` keys.

## History — `buildPushState`

`buildPushState(eventFuncId, url)` builds the `pushState` args from an event id
and URL; Plaid uses it (`buildPushStateArgs` / `runPushState`) to keep the browser
URL and history in sync with server navigation, and restores state on `popstate`.
