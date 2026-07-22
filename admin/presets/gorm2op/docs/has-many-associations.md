# Has-many associations

`SaveHasManyAssociation` reconciles a has-many field edited with the presets list
editor: it creates new children, updates existing ones, deletes removed ones, and
(optionally) writes an audit log per change.

```go
mb.UpdateDataOperator(func(do presets.DataOperator) presets.DataOperator {
    d := do.(*gorm2op.DataOperatorBuilder)
    return gorm2op.SaveHasManyAssociation("Parcelas").Build(d)
})
```

`Build` registers the reconciliation on both the **create** and **update**
pipelines. It returns the same `*DataOperatorBuilder`, so it composes with your
own callbacks (add them before calling `Build`).

## How it works

`Build` installs two callbacks:

- **pre** — records whether the list was present in the submitted form, then
  detaches the submitted slice from the object so the base create/update does not
  cascade-save it.
- **post** — reconciles the persisted children against the submitted list:
  1. items flagged `__deleted` (or absent from the list) are removed;
  2. items with a zero primary key (or flagged `__new`) are **created**;
  3. the remaining items are **updated**.

  A single `Association(...).Unscoped().Replace(kept)` applies the final set:
  it inserts the new children, keeps the updated ones, and deletes every
  persisted row that is no longer present.

## The list-editor form protocol

The list editor posts, for each row, the row's real fields plus per-item
metadata under `<field>[i].<meta>`. The operator reads this metadata straight
from the (reordered) form values, so it reflects the user's intent rather than
array position:

| Key | Meaning |
| --- | --- |
| `<field>.__present` | the list was in the form (even if empty). Present-and-empty ⇒ delete every child; absent ⇒ leave the children untouched. |
| `<field>[i].__pos` | the row's desired position. `UnmarshalForm` reorders the decoded slice (and reindexes these keys) by `__pos` before the operator runs. |
| `<field>[i].__new` | the row was created in the browser ⇒ create it. |
| `<field>[i].__deleted` | the row was removed. A **persisted** deleted row posts only `{ID, __deleted:true}` — its real fields are omitted, so nothing is overwritten and the UI can revert the deletion. |

These constants live in the presets package: `ListEditorPresentField`,
`ListEditorPositionField`, `ListEditorNewField`, `ListEditorDeletedField`.

> **Why presence matters.** Without the `__present` marker the operator cannot
> tell an intentionally emptied list (delete all children) from a form that
> simply did not include the field (leave them as is). The list editor always
> emits `__present`, and the operator also falls back to detecting any
> `<field>[i]` key for older forms.

## Auditing

When the model that **contains** the association is itself audited, each child
create/update/delete is logged. The operator stays decoupled from the activity
module via an interface resolved from the context:

```go
type AssociationAuditor interface {
    LogCreated(db *gorm.DB, obj any) error
    LogUpdated(db *gorm.DB, old, now any) error
    LogDeleted(db *gorm.DB, obj any) error
}
```

- The activity module injects an implementation into the request context while
  it saves an audited parent (`ContextWithAssociationAuditor`).
- The operator calls `AssociationAuditorFromContext`; when nothing is injected
  (**nil**), no logging happens and the extra snapshot reads are skipped.
- A child is logged only when its own type is registered with activity too;
  unregistered child types are silently skipped, so turning on parent audit
  never fails a save.

The auditor receives the transaction `*gorm.DB`, so log rows are written inside
the same transaction as the data change.

## Custom update & form key

- `Updator(func(db *gorm.DB, r any) error)` overrides how existing children are
  persisted (default: `db.Updates(r)`).
- `SaveHasManyAssociation(field, fieldFormKey…)` accepts an optional function to
  map the struct field name to the form key, when they differ.
- `Pre(...)` / `Post(...)` add extra callbacks that run around the association
  reconciliation.

## Behavior matrix

| Form state | Result |
| --- | --- |
| field absent | children left untouched |
| `__present`, no rows | all children deleted |
| row with zero id / `__new` | child created |
| row with id, real fields | child updated |
| row with id + `__deleted` | child deleted (fields untouched; revertible) |
