# List editor & tables

A has-many field (a slice of structs) is edited with the **list editor**: an
add/remove/reorder UI over the nested fields.

## Nested slice

```go
e.Field("Parcelas").
    Nested(presets.NestedSlice(itemModel, &itemFields).
        SetDisplayFieldInSorter("Nome"))
```

`NestedSlice(mb, fb)` renders each item with the nested fields of `fb`. Options:
`SetDisplayFieldInSorter`, `SetAddListItemRowEvent`, `SetRemoveListItemRowEvent`,
`SetSortListItemsEvent`, and `SetComponentBuilder` (see below).

## Client-side edit model

Editing happens on the reactive form; the server reconciles from the submitted
values. Each item posts, besides its real fields, per-item metadata under
`<field>[i].<meta>`:

| Key | Meaning |
| --- | --- |
| `<field>.__present` | the list was submitted (even if empty). Present + empty ⇒ delete all children; absent ⇒ leave untouched. |
| `<field>[i].__pos` | the row's desired position (the decoder reorders by it). |
| `<field>[i].__new` | the row was created in the browser. |
| `<field>[i].__deleted` | the row was removed. A persisted row posts only `{ID, __deleted:true}` so its fields are not overwritten and the deletion is revertible. |

Deletion is toggled on the reactive form client-side (no round-trip): the item's
inputs stay in the DOM (hidden) so it can be reverted, and the server data
operator (see `gorm2op.SaveHasManyAssociation`) creates / updates / deletes
children accordingly.

## Component builders

The list editor delegates rendering to a `ListEditorComponentBuilder`:

```go
type ListEditorComponentBuilder interface {
    Container(ctx *ListEditorContainerContext, items h.HTMLComponents) h.HTMLComponent
    Item(ctx *ListEditorItemContext) h.HTMLComponent
    DeletedItem(ctx *ListEditorItemContext) h.HTMLComponent
}
```

- `Item` renders an editable item; `DeletedItem` renders the removed placeholder
  (with a revert button), shown reactively in its place so order is preserved.
- Each builder applies its own `v-show` toggle via `ctx.DeletedCond()` /
  `ctx.DeleteExpr()` / `ctx.RevertExpr()`.

The default is `ListEditorListBuilder` — one card per item.

## Table renderer

`ListEditorTableBuilder` renders the items as a **table**: one column per field
(with nested/grouped headers built from the fields layout) and, in edit mode, a
trailing **Actions** column with the delete/revert control plus custom actions.
`Item` / `DeletedItem` return a `<tr>`; cells render with `FieldContext.MustInput`
so only the bare input shows (the label is in the header — see
[fields.md](fields.md#mustinput--bare-inputs-for-table-cells)).

```go
e.Field("Parcelas").
    Nested(presets.NestedSlice(itemModel, &itemFields).
        SetComponentBuilder(
            presets.NewListEditorTableBuilder().
                // extra per-item actions (before the built-in delete):
                Action(func(ctx *presets.ListEditorItemContext) h.HTMLComponent {
                    return myRowButton(ctx)
                }).
                // per-mode column layouts (optional; fall back to the current layout):
                ViewLayoutFunc(func(fb *presets.FieldsBuilder) presets.FieldsLayout { ... }).
                EditLayoutFunc(func(fb *presets.FieldsBuilder) presets.FieldsLayout { ... }).
                NewLayoutFunc(func(fb *presets.FieldsBuilder) presets.FieldsLayout { ... }),
        ))
```

Configuration:

| Method | Purpose |
| --- | --- |
| `Action(...)` | Extra per-item action components in the actions column. |
| `ActionsLabel(s)` | Override the actions column header. |
| `Density(s)` | Table density. |
| `ViewLayoutFunc` / `EditLayoutFunc` / `NewLayoutFunc` | Columns layout per mode (read-only / editing / creating). |

Both `Item`'s hidden fields (e.g. the primary key) are kept invisibly in the row
so it stays submittable — including the `{ID, __deleted}` payload when removed.
