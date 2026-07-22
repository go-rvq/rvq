# gorm2op — overview & architecture

`gorm2op` implements the presets `DataOperator` interface on top of GORM. The
admin (listing, editing, detailing) calls the operator; the operator runs the
matching GORM query.

## The operations

`DataOperatorBuilder` provides one method per admin operation:

| Method | Used by | GORM effect |
| --- | --- | --- |
| `Search(obj, params, ctx)` | listing | query + filters + pagination + count |
| `Fetch(obj, id, ctx)` | editing/detailing load | load one record by id |
| `FetchTitle(obj, id, ctx)` | breadcrumbs/labels | load the title of one record |
| `Save(obj, id, ctx)` | editing save | `Create` when id is zero, else `Update` |
| `Create(obj, ctx)` | creating | insert |
| `Update(obj, id, ctx)` | editing | update by id |
| `Delete(obj, id, cascade, ctx)` | listing delete | delete (optionally cascading) |

Create / update / delete run inside a **transaction** (`tx`), so callbacks and
the base write commit or roll back together.

## Two layers of customization

1. **Operator overrides** — replace the base behavior with a custom
   `Creator` / `Updator` / `Deleter` / `Finder` / `Preparer`. The `Preparer`
   shapes the `*gorm.DB` before every operation (e.g. add `Preload`, `Joins`,
   scopes).

2. **Callbacks** — pre/post hooks around each operation, without replacing the
   base write. This is the usual extension point: preload associations before a
   read, default a field before a create, reconcile a has-many field after a
   save. See [callbacks.md](callbacks.md).

`SaveHasManyAssociation` is built on the callback layer: it registers a pre
callback (detach the submitted slice) and a post callback (reconcile the
children) on both the create and update pipelines. See
[has-many-associations.md](has-many-associations.md).

## Files

| File | Contents |
| --- | --- |
| `operator.go` | `DataOperatorBuilder`, `Mode`, the operations, `CallbackState`. |
| `callback.go` | The callback pipeline: `Callbacks`, `CallbacksRegistrator`, context callbacks. |
| `assoc.go` | `SaveHasManyAssociation` — has-many reconciliation from the list-editor form. |
| `assoc_audit.go` | `AssociationAuditor` interface + context helpers (decoupled auditing). |
| `schema.go` | Lightweight `Schema` / `Field` helpers over a gorm model. |
| `clause.go` | `SetTableNameClauseBuilder` — override the table name in a query. |

## Reading the docs

1. [Data operator](data-operator.md)
2. [Callbacks](callbacks.md)
3. [Has-many associations](has-many-associations.md)
4. [API reference](api-reference.md)
