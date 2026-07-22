# Callbacks

Callbacks are pre/post hooks that run around an operation's base action, without
replacing it. They are the primary way to customize loading and saving.

## The pipeline

For a given operation the executed slice is:

```
pre callbacks…  →  base action (Create/Update/Delete/…)  →  post callbacks…
```

A `Callback` is `func(state *CallbackState) error`. Returning an error aborts the
pipeline (and rolls back the surrounding transaction for writes).

`CallbackState.Done(f)` registers a deferred function that runs after the
pipeline finishes (in registration order), useful for cleanup or restoring
temporarily-detached values.

## Registering callbacks

A `DataOperatorBuilder` embeds a `CallbacksRegistrator`. Register per mode:

```go
d := gorm2op.DataOperator(db)

d.WithReadCallbacks(func(cb *gorm2op.Callbacks[*gorm2op.DataOperatorBuilder]) {
    cb.Pre(func(state *gorm2op.CallbackState) error {
        state.DB = state.DB.Joins("Moeda").Preload("Parcelas")
        return nil
    })
})

d.WithCreateCallbacks(func(cb *gorm2op.Callbacks[*gorm2op.DataOperatorBuilder]) {
    cb.Pre(func(state *gorm2op.CallbackState) error {
        m := state.Obj.(*Movimentacao)
        if m.MoedaID == uuid.Nil {
            m.MoedaID = defaultMoedaID
        }
        return nil
    })
})
```

Selectors: `WithSearchCallbacks`, `WithFetchCallbacks`, `WithCreateCallbacks`,
`WithUpdateCallbacks`, `WithDeleteCallbacks`, plus the sets `WithReadCallbacks`
(search+fetch), `WithWriteCallbacks` (create+update), and the generic
`WithModeCallbacks(mode, …)` / `WithModeSplitCallbacks(mode, …)`.

Each selector hands you a `*Callbacks`, on which `Pre(...)` / `Post(...)` append
hooks.

## CallbackState

Passed to every callback:

| Field | Purpose |
| --- | --- |
| `Obj` | the record being operated on |
| `Ctx` | the `*web.EventContext` (request, form values, queries) |
| `DB` | the prepared, operation-scoped `*gorm.DB` (already `Model`-bound, id-filtered for update/delete) |
| `SharedDB` | a transaction-bound handle with a fresh statement — use it for side queries/writes (associations, extra rows) so they join the same transaction |
| `CommonDB` | a plain session on the same connection |
| `Set` / `Get` / `GetOk` | stash values to pass from a pre to a post callback |
| `Done(func() error)` | defer work until the pipeline ends |

Rule of thumb: mutate `state.DB` in a **read** pre callback to shape the query;
use `state.SharedDB` in a **write** post callback to persist related data inside
the same transaction.

## Context callbacks

Callbacks can also be injected per request through the context — they are merged
with the operator's own for the running mode:

```go
ctx.R = ctx.R.WithContext(
    gorm2op.AddCallbacksToContext(ctx.R.Context(), extraRegistrator),
)
```

`GetCallbacks(mode, ctx)` is what the operator uses internally to merge the
builder's callbacks with any found on the context. `NamedCallbacksRegistratorOf`
stores named registrators on a model builder for reuse.
