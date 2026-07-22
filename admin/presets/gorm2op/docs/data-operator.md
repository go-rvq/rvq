# Data operator

## Construction & wiring

```go
op := gorm2op.DataOperator(db)   // db is a *gorm.DB
mb.DataOperator(op)              // attach to a preset model builder
```

The `*gorm.DB` you pass is the base handle for every operation. Anything you set
on it (a default scope, `Preload`, a session config) applies throughout — e.g.
enabling full association saves:

```go
gorm2op.DataOperator(db.Session(&gorm.Session{FullSaveAssociations: true}))
```

You can also call the operator directly (outside a model builder) when
implementing a custom `SearchFunc` / `SaveFunc`:

```go
return gorm2op.DataOperator(db.Preload("Addresses.Phones")).Fetch(obj, id, ctx)
```

## Modes

Every operation maps to a `Mode` bit: `Search`, `Create`, `Fetch`, `FetchTitle`,
`Update`, `Delete` (plus `DeletedRelated`). Convenience sets: `Read =
Search|Fetch`, `Write = Create|Update`. Modes drive both the `Preparer` and the
per-mode callback selection.

## Prepare (shaping the query)

Before each operation the operator runs the **`Preparer`** on the `*gorm.DB`.
The default is a no-op; wrap it to add preloads, joins or scopes for specific
modes:

```go
op.WrapPrepare(func(old gorm2op.Preparer) gorm2op.Preparer {
    return func(db *gorm.DB, mode gorm2op.Mode, obj any, id model.ID, params *presets.SearchParams, ctx *web.EventContext) *gorm.DB {
        db = old(db, mode, obj, id, params, ctx)
        if mode.Has(gorm2op.Read) {
            db = db.Preload("Marcadores")
        }
        return db
    }
})
```

## Custom base behavior

Replace what an operation does at its core (the callbacks still run around it):

| Setter | Signature | Replaces |
| --- | --- | --- |
| `SetCreator` | `func(db, obj, ctx) error` | the insert |
| `SetUpdator` | `func(db, obj, id, ctx) error` | the update |
| `SetDeleter` | `func(old func() error, db, obj, id, cascade, ctx) error` | the delete (call `old()` to run the default) |
| `SetFinder` | `func(db, obj, ctx) (result any, err error)` | the search query |
| `SetPreparer` / `WrapPrepare` | `Preparer` | the query shaping |

When no custom function is set, the operator uses GORM directly: `Create` →
`db.Create`, `Update` → `db.Select("*").Updates` by id, `Delete` → `db.Delete`
by id.

## Transactions

`Create`, `Update` and `Delete` wrap their work in `db.Transaction`. Inside a
callback, use the DB handles on `CallbackState` (`state.DB`, `state.SharedDB`)
so your writes join the same transaction — see [callbacks.md](callbacks.md).
