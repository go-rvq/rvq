# package helper

Higher-level building blocks on top of the presets admin
(`github.com/go-rvq/rvq/admin/presets`) for editing **associations** from a
parent form.

## What's here

| File | Purpose |
|------|---------|
| `nestedslice.go` | `NestedSlice` / `NestedSliceBuilder` — edit a parent's **has-many** or **many-to-many** field as a list of related records (with a model selector for m2m). |
| `model_select.go` | `ModelSelectorBuilder` — a **belongs-to** (foreign-key) single/multi select bound to the owner's foreign-key column(s). |
| `association.go` | `AssociationManagerBuilder` — wires a data operator for an association field. |
| `inline_edit.go` | `InlineEdit` / `InlineEditModel` — edit a child record inline in the parent form. |
| `fields.go` | field / admin-tag helpers (`AdminTag`, `KeyValueArray`). |
| `messages.go`, `errors.go`, `url_params.go` | i18n messages, error mapping, URL-parameter helpers. |

## Primary keys: any type, composite supported

None of these helpers assume an integer column named `id`. Key columns are read
from the schema and used as their own type, so **integer, UUID and other key
types** all work, and **composite keys (two or more columns)** are supported.

### `NestedSlice` (has-many / many-to-many)

- Related rows are **listed** through a filter built from the relationship
  references — `hasManyParentFilter` (`<fk> = ? [AND ...]`) and `m2mParentFilter`
  (an `EXISTS(...)` correlated subquery over the join table). Every reference
  participates, so composite foreign/related keys are matched in full and the
  bound values come from the parent id (`parentFilterArgs`).
- Related rows are **linked / unlinked** through gorm's Association API
  (`associationAppend` → `Association(field).Append`, `associationDeleter` →
  `Association(field).Delete`). gorm manages the join table for any key type,
  composite keys and duplicate links. Unlinking removes only the join row, never
  the related record.
- A legacy single-key link query (`m2mInsertQuery`) is still exposed on
  `NestedSliceBuilderInfo.LinkInsertQuery` for callers that read it; the actual
  link/unlink no longer depends on it.

### `ModelSelectorBuilder` (belongs-to)

- The owner's **foreign-key columns** are resolved from the relationship
  (`foreignKeyFieldsOf`) in related-primary-key order, falling back to the
  `<Field>ID` convention for a single key or a non-gorm operator.
- When a record is selected, **all** foreign-key columns are written
  (`model.ID.Related(schema, fkFields...).SetTo(owner)`); clearing the selection
  zeroes them all. So a composite foreign key is persisted in full.
- Selected values are exchanged as `model.ID` strings (composite keys join their
  parts with `_`; see `presets.ParseRecordID`).

## Related

The presets **list editor** itself (`presets.NestedSlice` +
`gorm2op.SaveHasManyAssociation`) reconciles has-many children (create / update /
delete) and also supports composite primary keys — a separate, lower layer from
these selector-oriented builders. See `admin/presets/listeditor_partition.go` and
`admin/presets/tests/listeditor/`.

## Tests

- `nestedslice_test.go` — non-integer (UUID-like) m2m link; composite related-key
  m2m filter + link/unlink via the Association API.
- `model_select_test.go` — composite belongs-to foreign-key resolution and the
  `model.ID.Related` mapping onto the owner columns.

See `docs/keys.md` for the key-handling rationale and a worked example.
