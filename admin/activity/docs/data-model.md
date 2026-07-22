# Data model

## `ActivityLog`

The default log row (table `activity_logs`, auto-migrated by `New`):

| Field | Type | Notes |
| --- | --- | --- |
| `ID` | `uint` | primary key |
| `UserID` | `uuid.UUID` | set when the creator is a `CreatorInterface`; indexed |
| `CreatedAt` | `time.Time` | when the action happened |
| `Creator` | `string` | creator name (`"unknown"` when none resolved) |
| `Action` | `string` | `Create` / `Edit` / `Delete` / `View`, or a custom action |
| `ModelKeys` | `string` | primary-key values of the record, joined by `:`; indexed |
| `ModelName` | `string` | Go type name of the record; indexed |
| `ModelLabel` | `string` | preset resource URI name, or `-` |
| `ModelLink` | `string` | optional link to the record (see `LinkFunc`) |
| `ModelDiffs` | `string` | JSON array of `Diff` (edits only) |
| `IP` | `string` | request origin IP; indexed |
| `UserAgent` | `string` | request user agent |

`ModelDiffs` holds a JSON array of:

```go
type Diff struct {
    Field string
    Old   string
    Now   string
}
```

## Interfaces

### `ActivityLogInterface`

Every log model must implement this (getter/setter per column):
`Set/GetCreatedAt`, `Set/GetUserID`, `Set/GetCreator`, `Set/GetAction`,
`Set/GetModelKeys`, `Set/GetModelName`, `Set/GetModelLabel`, `Set/GetModelLink`,
`Set/GetModelDiffs`.

### `RequestInfoSetter` (optional)

```go
type RequestInfoSetter interface {
    SetIP(string)
    SetUserAgent(string)
}
```

When the log model implements it **and** the context carried request info (via
`ContextWithRequestInfo`), the IP and user agent are persisted. `ActivityLog`
implements it.

### `CreatorInterface` (optional, for the creator value)

```go
type CreatorInterface interface {
    GetID() uuid.UUID
    GetName() string
}
```

A creator that implements this records both the user id (`UserID`) and name
(`Creator`); a plain `string` creator records only the name.

## Custom log table

Pass your own model to `New` to store logs elsewhere or with extra columns:

```go
type AuditLog struct {
    activity.ActivityLog          // embed to inherit the interface + columns
    TenantID uuid.UUID `gorm:"index"`
}

ab := activity.New(db, &AuditLog{})
```

Embedding `ActivityLog` is the easy path: you inherit the full
`ActivityLogInterface` (and `RequestInfoSetter`) and only add what you need.
Otherwise, implement `ActivityLogInterface` yourself.

## Retrieving logs

- `ab.GetActivityLogs(record, db) []*ActivityLog` — logs for a record (only when
  the default `ActivityLog` model is used).
- `ab.GetCustomizeActivityLogs(record, db) interface{}` — same, but returns the
  configured (possibly custom) slice type.
- `mb.RecordLogs(id string)` — logs for a record id, as `[]ActivityLogInterface`.

Both `GetActivityLogs`/`GetCustomizeActivityLogs` match on `model_name` +
`model_keys` and fall back to the global DB when `db` is nil.
