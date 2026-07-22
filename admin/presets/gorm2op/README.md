# gorm2op

The GORM implementation of the presets `DataOperator`: it turns the admin's
search / fetch / create / update / delete operations into GORM queries.

Two things live here:

- **`DataOperatorBuilder`** — the data operator itself, with a **callback
  pipeline** (pre/post hooks per operation) you use to customize how records are
  loaded and saved.
- **`SaveHasManyAssociation`** — a builder that reconciles a **has-many** field
  edited with the list editor: it creates new children, updates existing ones,
  and deletes removed ones, optionally writing an audit log for each change.

## Quick start

```go
import "github.com/go-rvq/rvq/admin/presets/gorm2op"

// wire the operator into a model
mb.DataOperator(gorm2op.DataOperator(db))
```

```go
// customize loading + save a has-many field
mb.UpdateDataOperator(func(do presets.DataOperator) presets.DataOperator {
    d := do.(*gorm2op.DataOperatorBuilder).
        WithReadCallbacks(func(cb *gorm2op.Callbacks[*gorm2op.DataOperatorBuilder]) {
            cb.Pre(func(state *gorm2op.CallbackState) error {
                state.DB = state.DB.Preload("Parcelas")
                return nil
            })
        })
    return gorm2op.SaveHasManyAssociation("Parcelas").Build(d)
})
```

## Documentation

Detailed docs live under [`docs/`](docs/):

- [Overview & architecture](docs/README.md)
- [Data operator](docs/data-operator.md) — operations, `Prepare`, transactions, custom create/update/delete/find.
- [Callbacks](docs/callbacks.md) — the pre/post pipeline, per-mode hooks, `CallbackState`, context callbacks.
- [Has-many associations](docs/has-many-associations.md) — `SaveHasManyAssociation`, the list-editor form protocol, deletes & audit.
- [API reference](docs/api-reference.md).

## Recent change — has-many reconciliation

`SaveHasManyAssociation` no longer relies on the global `ModifiedIndexesBuilder`.
Each submitted item now carries its own metadata (`__pos`, `__new`, `__deleted`,
`__present`), so the operator reconciles children from the form itself:

- **presence** — a `__present` marker distinguishes "the list was submitted but
  empty" (delete every child) from "the field was absent" (leave children as is);
- **deletes** — a removed persisted row posts only `{ID, __deleted:true}`, so its
  real fields are never overwritten and it stays revertible in the UI;
- **audit** — when the parent model is audited, each child create/update/delete
  is logged through an `AssociationAuditor` resolved from the context.

See [has-many-associations.md](docs/has-many-associations.md) for the full
protocol.
