# API reference

A condensed map of the exported surface. See the source for full signatures.

## Constructor & install

| Symbol | Purpose |
| --- | --- |
| `New(db *gorm.DB, logModel ...ActivityLogInterface) *Builder` | Create the module; auto-migrates the log table; sets the read-only perm policy. |
| `(*Builder) Install(pb *presets.Builder) error` | Presets plugin install (called by `pb.Use(ab)`). |
| `(*Builder) ModelInstall(pb *presets.Builder, m *presets.ModelBuilder) error` | Per-model plugin install (called by `mb.Use(ab)`); registers `m`. |

## Builder — configuration

| Symbol | Purpose |
| --- | --- |
| `CreatorContextKey(key) *Builder` | Context key for the creator. |
| `DBContextKey(key) *Builder` | Context key for the request/transaction DB. |
| `TabHeading(func(ActivityLogInterface) string) *Builder` | Heading format in the activity tab. |
| `PermPolicy(*perm.PolicyBuilder) *Builder` | Replace the default permission policy. |
| `WrapLogModelInstall(...) *Builder` | Wrap how the log resource is installed. |

## Builder — model registration

| Symbol | Purpose |
| --- | --- |
| `RegisterModel(m interface{}) *ModelBuilder` | Register a preset model (auto-logging) or a plain struct (manual). |
| `RegisterModels(models ...interface{}) *Builder` | Register several. |
| `GetModelBuilder(v) (*ModelBuilder, bool)` | Look up a registered builder. |
| `MustGetModelBuilder(v) *ModelBuilder` | Same; panics if missing. |
| `GetModelBuilders() []*ModelBuilder` | All registered builders. |

## Builder — recording (context-based)

| Symbol | Purpose |
| --- | --- |
| `AddRecords(action string, ctx context.Context, vs ...interface{}) error` | Record for each value; `Edit` diffs against the DB row. |
| `AddCustomizedRecord(action string, diff bool, ctx, obj) error` | Record an arbitrary action, optionally with a diff. |
| `AddEditRecordWithOldAndContext(ctx, old, now) error` | Edit log from a held snapshot, resolving creator/db from ctx. |

## Builder — recording (explicit creator + db)

| Symbol | Purpose |
| --- | --- |
| `AddCreateRecord(creator, v, db) error` | Create log. |
| `AddEditRecord(creator, now, db) error` | Edit log; fetches old from db. |
| `AddEditRecordWithOld(creator, old, now, db) error` | Edit log from a held snapshot. |
| `AddDeleteRecord(creator, v, db) error` | Delete log. |
| `AddViewRecord(creator, v, db) error` | View log. |
| `AddSaveRecord(creator, now, db) error` | Create if no old row, else edit. |

The same recording methods exist on `*ModelBuilder` (use those to disambiguate
when several builders share a Go type), plus `AddEditRecordWithOldCtx(ctx, old,
now)`.

## Builder — querying

| Symbol | Purpose |
| --- | --- |
| `GetActivityLogs(m, db) []*ActivityLog` | Logs for a record (default model). |
| `GetCustomizeActivityLogs(m, db) interface{}` | Logs as the configured slice type. |
| `NewLogModelData() interface{}` / `NewLogModelSlice() interface{}` | Fresh log value / slice. |

## ModelBuilder — configuration

| Symbol | Purpose |
| --- | --- |
| `AddKeys(...)` / `Keys(...)` | Identity keys. |
| `AddIgnoredFields(...)` / `SetIgnoredFields(...)` | Diff exclusions. |
| `AddTypeHanders(v, TypeHandler)` | Custom diffing for a type. |
| `LinkFunc(func(interface{}) string)` | Link stored on logs. |
| `SkipCreate()` / `SkipUpdate()` / `SkipDelete()` | Opt out of auto-logging (preset models). |
| `EnableActivityInfoTab()` | Activity tab on the editing page. |
| `KeysValue(v) string` | The joined key value for a record. |
| `Diff(old, now) ([]Diff, error)` | Compute a diff. |
| `RecordLogs(id string) ([]ActivityLogInterface, error)` | Logs for a record id. |

## Context helpers

| Symbol | Purpose |
| --- | --- |
| `ContextWithCreator(ctx, name string) context.Context` | Store the creator name. |
| `ContextWithDB(ctx, db *gorm.DB) context.Context` | Store the request/transaction DB. |
| `ContextWithRequestInfo(ctx, r *http.Request) context.Context` | Store request origin (IP / user agent). |
| `RequestInfoFromContext(ctx) *RequestInfo` | Read it back. |
| `RequestInfoFromRequest(r) *RequestInfo` | Build request info from a request. |

## Constants & vars

| Symbol | Value / purpose |
| --- | --- |
| `ActivityCreate` / `ActivityEdit` / `ActivityDelete` / `ActivityView` | Action strings. |
| `LogViewAction` | `"activityLog"` — detail action name. |
| `I18nActivityKey` | i18n module key. |
| `DefaultIgnoredFields` | `ID`, `UpdatedAt`, `DeletedAt`, `CreatedAt`. |
| `DefaultTypeHandles` | Built-in handlers (`time.Time`, `MediaBox`). |

## Types

`Builder`, `ModelBuilder`, `ActivityLog`, `ActivityLogInterface`,
`CreatorInterface`, `RequestInfo`, `RequestInfoSetter`, `Diff`, `DiffBuilder`,
`TypeHandler`.
