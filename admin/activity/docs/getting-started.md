# Getting started

## 1. Create the builder

```go
import "github.com/go-rvq/rvq/admin/activity"

ab := activity.New(db)
```

`New(db *gorm.DB, logModel ...ActivityLogInterface) *Builder`:

- **`db`** (required) is the global database. It is used whenever a specific
  DB is not supplied through the context (see below).
- **`logModel`** (optional) lets you store logs in your own table; it must
  implement `ActivityLogInterface`. When omitted, the built-in `ActivityLog`
  is used.

`New` **auto-migrates** the log model and installs a default permission policy
that denies `Create` / `Update` / `Delete` on the `activity_logs` resource (logs
are meant to be read-only in the admin).

## 2. Install into presets

`Builder` is a presets plugin — install it with `Use`:

```go
b := presets.New()
b.Use(ab)
```

This registers the i18n messages, applies the permission policy, and exposes the
log table as an admin resource (with a `mdi-book-edit` menu icon).

## 3. Register models

### Preset models (automatic logging)

Register a `*presets.ModelBuilder` and create / edit / delete performed through
the admin are logged automatically:

```go
postMB := b.Model(&Post{})
postMB.Use(ab)            // idiomatic: register as a per-model plugin
// or, equivalently:
ab.RegisterModel(postMB)
```

Registering a preset model also installs the **`activityLog`** detail action,
which opens a timeline of that record's history.

### Plain structs (manual logging)

```go
ab.RegisterModel(&Comment{})
ab.RegisterModels(&Comment{}, &Tag{}) // register several at once
```

Plain structs are not wired into any flow — you record their logs explicitly
with the `Add*Record` helpers (see [recording-logs.md](recording-logs.md)).

## 4. Wire creator, db and request info through the context

The recording API resolves three things from the `context.Context`:

| Value | Default context key | Set with | Configure key with |
| --- | --- | --- | --- |
| creator | `CreatorContextKey` | `ContextWithCreator(ctx, name)` | `ab.CreatorContextKey(key)` |
| db | `DBContextKey` | `ContextWithDB(ctx, db)` | `ab.DBContextKey(key)` |
| request origin | (internal) | `ContextWithRequestInfo(ctx, r)` | — |

A common setup points the creator key at the logged-in user so the current user
becomes the log's creator automatically:

```go
ab := activity.New(db).CreatorContextKey(login.UserKey)
```

The creator resolved from the context may be a `string` or a `CreatorInterface`
(`GetID() uuid.UUID`, `GetName() string`). A `CreatorInterface` also records the
user id on the log; a bare string records only the name. When nothing is found
the creator falls back to `"unknown"`.

> The automatic (preset) path already wraps the request with
> `ContextWithRequestInfo`, so IP / user agent are captured for you. For manual
> recording, wrap the context yourself if you want that origin persisted.
