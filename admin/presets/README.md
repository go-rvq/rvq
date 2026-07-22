# presets

`presets` is the core of the `go-rvq/rvq` admin: it turns a Go model into a full
CRUD resource — **listing**, **editing**, **creating** and **detailing** — driven
by a declarative **fields** layer, with pluggable **data operators**,
**permissions** and **i18n**.

## Quick start

```go
import (
    "github.com/go-rvq/rvq/admin/presets"
    "github.com/go-rvq/rvq/admin/presets/gorm2op"
    "github.com/go-rvq/rvq/x/i18n"
)

b := presets.New(i18n.New())             // the admin Builder
b.DataOperator(gorm2op.DataOperator(db)) // default persistence (GORM)

post := b.Model(&Post{})                 // register a resource
post.Listing("Title", "Status")          // list columns
post.Editing("Title", "Body", "Status")  // edit form fields
post.Detailing("Title", "Body")          // detail view
```

Mount `b` as an `http.Handler`.

## Concepts

| Type | Role |
| --- | --- |
| `Builder` | The admin. Holds models, i18n, permissions, the default data operator; created with `New(i18n.Builder)`. |
| `ModelBuilder` | One resource. Exposes `Listing` / `Editing` / `Creating` / `Detailing`. |
| `ListingBuilder` | The index list (read-only table) — columns, filters, bulk actions, pagination. |
| `EditingBuilder` / `CreatingBuilder` | The edit / create form — fields, validators, save hooks. |
| `DetailingBuilder` | The read-only detail view — fields, actions, tabs. |
| `FieldsBuilder` / `FieldBuilder` | The fields layer: which fields, their layout, component and setter. |
| `FieldContext` | Per-render field state (obj, form key, mode, label, errors…). |
| `DataOperator` | Persistence abstraction (search/fetch/save/delete); `gorm2op` is the GORM impl. |
| `Plugin` | A unit installed with `b.Use(...)` (e.g. `activity`, `login`, packages). |

## Documentation

Detailed docs live under [`docs/`](docs/):

- [Getting started](docs/getting-started.md) — the Builder, registering models, plugins, data operator.
- [Models](docs/models.md) — Listing / Editing / Creating / Detailing.
- [Fields](docs/fields.md) — FieldsBuilder, FieldBuilder, FieldContext, the default component funcs, and `MustInput`.
- [List editor & tables](docs/list-editor.md) — nested slices, the client-side list editor, and the table renderer.
- [Record IDs](docs/record-id.md) — how primary keys are encoded/parsed.

## Feature map

```
Builder (New)
  Model
    Listing        — columns, Filter Menu/Tabs, Bulk Actions (UpdateForm/UpdateFunc), pagination
    Editing        — fields, validators, SaveFunc; multiple forms by name
    Creating       — the create form (derived from Editing)
    Detailing      — fields, Actions (UpdateForm/UpdateFunc), tabs
  Fields           — FieldsBuilder / FieldBuilder / FieldContext, nested (slice/struct)
  DataOperator     — gorm2op (search/fetch/save/delete + callbacks)
  Use(Plugin...)   — activity, login helpers, packages
```
