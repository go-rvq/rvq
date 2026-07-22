# Models

A `*ModelBuilder` (from `b.Model(&T{})`) exposes four sub-builders. Each takes an
optional list of field names to seed its fields layout.

## Listing

The index list — a read-only table with search, filters, pagination and bulk
actions. (It is the model's INDEX LIST; it is always read-only and distinct from
the editable table renderer of the list editor.)

```go
l := m.Listing("Title", "Status", "CreatedAt").PerPage(20)
l.SearchColumns("title")
l.FilterDataFunc(func(ctx *web.EventContext) vx.FilterData { ... })   // Filter Menu
l.BulkAction("approve").UpdateFunc(func(ids []string, ctx *web.EventContext) error { ... })
```

Highlights: `SearchFunc` / `WrapSearchFunc`, `OrderableFields`,
`SelectableColumns`, `CellWrapperFunc`, filter menu/tabs, bulk actions
(`UpdateForm` / `UpdateFunc`). Columns are built from the fields layout (with
nested/grouped headers when the layout has sections).

## Editing & Creating

The edit form. `Creating` is derived from `Editing` (the create form).

```go
e := m.Editing("Title", "Body", "Status")
e.Field("Title").ComponentFunc(...).SetterFunc(...)
e.Validators.AppendFunc(func(obj any, mode presets.FieldModeStack, ctx *web.EventContext) web.ValidationErrors { ... })
e.WrapSaveFunc(func(in presets.SaveFunc) presets.SaveFunc { ... })

c := e.CreatingBuilder()      // the create variant
c.WrapCreateFunc(...)
```

- Multiple named forms: `m.Editing(name, ...)` lets a model have more than one
  edit form (e.g. a change-password sub-form).
- Save flow: `SaveFunc` / `WrapSaveFunc`; the create flow uses `CreateFunc`.
- Fields can be `HiddenField(...)`, laid out in rows/sections, and validated.

## Detailing

The read-only detail view — fields, actions and tabs.

```go
d := m.Detailing("Title", "Body", "Author")
d.Field("Author").ComponentFunc(...)
d.Action("publish").OnClick(func(ctx *web.EventContext, id string, obj any) string { ... })
d.AppendTabsPanelFunc(func(obj any, ctx *web.EventContext) (tab, content h.HTMLComponent) { ... })
```

- Actions: `UpdateForm` / `UpdateFunc`, `OnClick`, enabled predicates, i18n labels.
- Tabs: e.g. the `activity` module adds an "Activities" tab.

## Restrictions

Each builder has restriction hooks that gate operations per request/object:
`CreatingRestriction`, `EditingRestriction`, `DetailingRestriction`,
`DeletingRestriction`, `DeletingWithRelatedRestriction`.

```go
m.DeletingRestriction.ObjHandler(presets.OkObjHandlerFuncT(
    func(u *User, ctx *web.EventContext) (ok, handled bool) { ... }))
```

## Fields

All four builders embed a `FieldsBuilder`. Field configuration (component, setter,
labels, nesting) is shared — see [fields.md](fields.md).
