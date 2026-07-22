# Model configuration

Everything below is configured on the `*ModelBuilder` returned by
`RegisterModel` (or fetched with `GetModelBuilder` / `MustGetModelBuilder`).
Methods are chainable.

```go
mb := ab.RegisterModel(postMB)
mb.AddKeys("ID", "Version").
   AddIgnoredFields("UpdatedAt").
   EnableActivityInfoTab()
```

## Keys (record identity)

The keys identify *which record* a log belongs to; their values are joined into
`ModelKeys`. They default to the model's gorm primary keys.

- `AddKeys(keys ...string)` — append keys (deduplicated).
- `Keys(keys ...string)` — replace the key set.

Use `AddKeys` when a record's identity is a composite (e.g. `ID` + `Version`).

## Ignored fields (diff noise)

Ignored fields are excluded from edit diffs. They default to the model's primary
keys; the package-level `DefaultIgnoredFields` (`ID`, `UpdatedAt`, `DeletedAt`,
`CreatedAt`) documents the fields you typically also want to skip.

- `AddIgnoredFields(fields ...string)` — append to the current set.
- `SetIgnoredFields(fields ...string)` — replace the set.

```go
mb.AddIgnoredFields("UpdatedAt", "SearchIndex")
```

## Type handlers (custom diffing)

A `TypeHandler` controls how a field type is diffed:

```go
type TypeHandler func(old, now interface{}, prefixField string) []Diff
```

Register one per type:

```go
mb.AddTypeHanders(time.Time{}, func(old, now interface{}, prefixField string) []Diff {
    o := old.(time.Time).Format(time.RFC3339)
    n := now.(time.Time).Format(time.RFC3339)
    if o == n {
        return nil
    }
    return []Diff{{Field: prefixField, Old: o, Now: n}}
})
```

The package ships `DefaultTypeHandles` for `time.Time` and
`media_library.MediaBox`.

## Skipping automatic logging

For registered **preset** models, opt out of individual auto-logged actions
(no effect on plain structs):

- `SkipCreate()`
- `SkipUpdate()`
- `SkipDelete()`

## Record link

`LinkFunc(func(obj interface{}) string)` sets the URL stored on each log
(`ModelLink`), letting the log view link back to the modified record.

## Editing-page activity tab

`EnableActivityInfoTab()` adds an "Activities" tab to the model's editing page
showing that record's log timeline. (A registered preset model also always gets
the standalone **`activityLog`** detail action — see below.)

## Builder-level options

Set on the `Builder`, not per model:

- `CreatorContextKey(key)` — where to read the creator from the context
  (e.g. `login.UserKey`).
- `DBContextKey(key)` — where to read the request-scoped/transaction DB from the
  context.
- `TabHeading(func(ActivityLogInterface) string)` — customize the heading shown
  for each log entry in the activity tab.
- `PermPolicy(*perm.PolicyBuilder)` — replace the default (read-only) permission
  policy for the log resource.
- `RegisterModels(models ...interface{})` — register several models at once.
