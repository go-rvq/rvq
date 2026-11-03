# package utils (thirdpart/gorm/utils)

Small, reusable helpers on top of [gorm](https://gorm.io) — the **single source
of truth** for reading gorm schema/relationship metadata (primary/foreign keys,
composite keys, any key type) and for a few query-building conveniences. Consumed
by `admin/presets/gorm2op`, `admin/packages/helper` and application code.

Import path: `github.com/go-rvq/rvq/thirdpart/gorm/utils`.

## Files

| File | What |
|------|------|
| `relationship.go` | Key & relationship helpers: parent filters, foreign-key resolution, primary-key matching. Composite- and any-type-aware. |
| `select_column.go` | `RawColumn` / `SetRawColumn` — a SELECT column made of a raw SQL expression. |
| `with.go` | `WithClause` / `WithClauses` / `WithDbClause` — a `WITH [RECURSIVE]` (common table expressions) before a gorm query: the `clause.With` gorm lacks. |
| `capture.go` | `RegisterColumnCapture` / `CaptureColumns` — read the columns a model has no field for (gorm discards them). |
| `association.go` | `AssociationDB` / `Instance` — `//go:linkname` accessors into gorm internals (`(*Association).buildCondition`, `(*DB).getInstance`). |

## Keys & relationships (`relationship.go`)

Nothing assumes an integer column named `id`. Columns are read from the schema /
relationship and used as their own type, and **composite keys (two or more
columns) are supported**.

- `ForeignKeyFields(db, model, field) []string` — the owner's foreign-key field
  names of a belongs-to `field`, in related-primary-key order.
- `HasManyParentFilter(rel) (sql, ownerFields)` — `WHERE <fk> = ? [AND ...]` to
  filter a has-many child by its parent, plus the owner PK field names (bind
  order).
- `M2MParentFilter(rel) (sql, ownerFields)` — an `EXISTS(...)` correlated
  subquery over the join table, matching related rows linked to a parent.
- `ParentFilterArgs(id, fieldNames) []any` — extract the parent key values (in
  `fieldNames` order) from a `FieldValuer` (e.g. `model.ID`) to bind the filter.
- `PrimaryFieldNames(schema) []string` — every primary-key field name.
- `PKAllZero(item, pkNames) bool` — the row has no usable key (all PK fields
  zero) → treat as new/keyless.
- `PKMapKey(item, pkNames) any` — an unambiguous, comparable map key from the
  (possibly composite) primary key.
- `NormalizeKey(value) any` — normalize an int/uint/string key for map use.

To **link/unlink** related records, prefer gorm's Association API
(`db.Model(parent).Association(field).Append/Delete`) rather than raw SQL — it
handles any key type, composite keys and duplicate links natively. These helpers
provide only the read/filter side that the Association API does not.

See `docs/keys.md` for the rationale and worked examples, and
`relationship_test.go` for composite has-many / many-to-many / belongs-to tests.

## WITH clauses (`with.go`)

gorm's `clause.With` is an empty type. `WithClauses` builds the WITH, and
`DbClause()` puts it before the SELECT (or, named, before the UPDATE or
DELETE) of a query:

```go
with := gormutils.WithClauses{{
	Name:      "tree",
	Recursive: true, // any recursive table: WITH RECURSIVE
	Query: `SELECT n.* FROM nodes n WHERE n.id = ?
UNION
SELECT c.* FROM nodes c JOIN tree t ON c.parent_id = t.id`,
	Args: []any{rootID},
}}
err := db.Clauses(with.DbClause()).Table("tree").Find(&nodes).Error
// WITH RECURSIVE tree AS (…) SELECT * FROM `tree`
```

The main query refers to the tables by name (`Table`, `Joins`, a subquery),
and two `DbClause` on one query make a single WITH. A table can also be a query of gorm (`WithClause.DB`), with its own values:

```go
with := gormutils.WithClauses{{Name: "big", DB: db.Model(&Node{}).Select("id", "value").Where("value > ?", 100)}}
db.Clauses(with.DbClause()).Table("big").Find(&nodes)
// WITH big AS (SELECT `id`,`value` FROM `nodes` WHERE value > 100) SELECT * FROM `big`
``` More examples, and how it
works: `docs/with.md`.

## Capturing columns (`capture.go`)

When gorm scans into a model, it reads a column the model has no field for
into a local variable and discards it. There is no option or hook to get it
back (gorm v1.31.2). `RegisterColumnCapture(db)` puts a query callback in
place of `gorm:query`, once at the start of the app. A query marked with
`CaptureColumns` gets its rows wrapped, and the values of the named columns of
the first row are kept. A query that is not marked runs gorm's own callback.

```go
q, capt := gormutils.CaptureColumns(db.Model(&Node{}), "__count")
q.Select("nodes.*", "COUNT(*) OVER () AS __count").Limit(10).Find(&page)
capt.Values["__count"] // the count of the whole filter
```

- **Hooks and preloads:** the model keeps its hooks (`AfterFind`) and its
  preloads. The mark is only for that query (`InstanceSet`), not for the
  preloads.
- **Values:** a value comes as the driver gives it. Bytes (a JSON in
  PostgreSQL) come as a copied string.
- **No row, nothing captured.**

## Raw SELECT columns (`select_column.go`)

`SetRawColumn(table, name, query)` selects `query AS name` in place of the
column `table.name` of the model. That column can be selected by `Select`,
or by the model itself (`QueryFields`, a `Joins`). When nothing is selected,
the expression goes before the `*`:

```go
db.Clauses(gormutils.SetRawColumn("nodes", "double", "value * 2")).
	Table("nodes").Select("id", "double").Find(&rows)
// SELECT `id`,value * 2 AS double FROM `nodes`
```

The query is written as it is: `Args` are not bound.

## Examples

`example_test.go` has runnable examples of the WITH clauses and the raw
columns. `go test` checks the SQL they build, using a dry-run SQLite
database (`go doc -all` shows them).
