# activity — overview & architecture

`activity` is the audit-log module for the `go-rvq/rvq` admin. Every audited
action produces one row in the `activity_logs` table describing:

- **who** — the creator (a name string, or a `CreatorInterface` with id + name);
- **what** — the action (`Create` / `Edit` / `Delete` / `View`);
- **which** — the model name, its label, and its primary-key values;
- **when** — a timestamp;
- **details** — a JSON field-level diff (for edits) and the request origin
  (IP / user agent).

## Building blocks

| Type | Role |
| --- | --- |
| `Builder` | The module. Holds the DB, the log model, the registered models, the permission policy and the context keys. Created with `New`. Installed into a `presets.Builder` with `Use`. |
| `ModelBuilder` | Per-model configuration and the entry point for recording logs. Obtained from `RegisterModel` (or `b.Model(&X{}).Use(ab)`). |
| `ActivityLog` | The default log row. Any type implementing `ActivityLogInterface` can replace it. |
| `Diff` / `DiffBuilder` / `TypeHandler` | Compute the field-level diff stored on edit logs. |
| `RequestInfo` | Request origin (IP, user agent) carried through the context. |

## How it fits together

```
presets.Builder
   └── Use(activity.Builder)          // i18n + perm policy + "activity_logs" admin resource
          ├── RegisterModel(postMB)   // *presets.ModelBuilder → auto create/edit/delete logging
          │      └── installs LogViewAction detail action ("activityLog")
          └── RegisterModel(&Plain{}) // plain struct → record manually with Add*Record
```

For a registered **preset** model the module wraps the editing / creating /
deleting flow (`WrapSaveFunc`, `WrapCreateFunc`, `WrapDeleteFunc`), so no
call-site changes are needed. For plain structs — or records saved outside the
preset flow, such as **has-many children** or **custom actions** — you call the
recording API yourself (see [recording-logs.md](recording-logs.md)).

## Reading the docs

1. [Getting started](getting-started.md) — install, register, wire creator / db / request info.
2. [Recording logs](recording-logs.md) — the automatic path, the manual API, and nested records.
3. [Model configuration](model-configuration.md) — keys, ignored fields, type handlers, skips, tabs, links.
4. [Data model](data-model.md) — the `ActivityLog` schema, the interfaces, and custom log tables.
5. [API reference](api-reference.md) — every exported symbol.
