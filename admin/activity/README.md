# activity

Audit logging for the `go-rvq/rvq` admin. The package records **who** did **what**
to **which record** and **when**, persisting a row per action (create / edit /
delete / view) into an `activity_logs` table, including a field-level diff for
edits and the request origin (IP / user agent).

It plugs into `presets` in two ways:

- **Automatically** — register a `*presets.ModelBuilder` and create / edit /
  delete performed through the admin are logged for you.
- **Manually** — for plain structs, or records saved outside the preset flow
  (nested associations, batch jobs, custom actions), call the `Add*Record`
  helpers yourself.

## Quick start

```go
import "github.com/go-rvq/rvq/admin/activity"

ab := activity.New(db)          // creates & auto-migrates the activity_logs table
b.Use(ab)                        // install into the presets.Builder (i18n + perm + admin resource)

ab.RegisterModel(postModelBuilder)          // auto-log create/edit/delete
ab.RegisterModel(&Comment{})                // plain struct: record manually
```

```go
// manual recording (creator, db and request origin resolved from context)
ctx := activity.ContextWithCreator(
    activity.ContextWithRequestInfo(r.Context(), r), "alice")
ab.AddRecords(activity.ActivityCreate, ctx, comment)
```

## Documentation

Detailed docs live under [`docs/`](docs/):

- [Overview & architecture](docs/README.md)
- [Getting started](docs/getting-started.md) — install, register models, wire creator/db/request-info
- [Recording logs](docs/recording-logs.md) — automatic vs. manual, context helpers, nested records
- [Model configuration](docs/model-configuration.md) — keys, ignored fields, type handlers, skip, tabs, links
- [Data model](docs/data-model.md) — `ActivityLog`, the interfaces, and custom log tables
- [API reference](docs/api-reference.md) — every exported symbol at a glance

## Actions

| Constant | Value | Meaning |
| --- | --- | --- |
| `activity.ActivityCreate` | `"Create"` | a record was created |
| `activity.ActivityEdit` | `"Edit"` | a record was edited (stores a field diff) |
| `activity.ActivityDelete` | `"Delete"` | a record was deleted |
| `activity.ActivityView` | `"View"` | a record was viewed |
