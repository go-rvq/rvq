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
| `select_column.go` | `RawColumn` / `SetRawColumn` — inject a raw SQL expression as a named SELECT column. |
| `with.go` | `WithClause` / `WithClauses` / `WithDbClause` — build and attach SQL `WITH` (CTE) clauses to a gorm statement. |
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
