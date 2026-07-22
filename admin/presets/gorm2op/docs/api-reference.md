# API reference

A condensed map of the exported surface. See the source for full signatures.

## Operator

| Symbol | Purpose |
| --- | --- |
| `DataOperator(db *gorm.DB) *DataOperatorBuilder` | Construct the GORM data operator. |
| `(*DataOperatorBuilder) Search / Fetch / FetchTitle / Save / Create / Update / Delete` | The presets `DataOperator` operations. |
| `Clone()` / `CloneDataOperator()` | Copy the operator (fresh session). |
| `DB()` / `SetDB(db)` | Access/replace the base handle. |

### Base-behavior overrides

| Symbol | Type it sets |
| --- | --- |
| `SetPreparer` / `WrapPrepare` | `Preparer func(db, mode, obj, id, params, ctx) *gorm.DB` |
| `SetCreator` / `Creator()` | `Creator func(db, obj, ctx) error` |
| `SetUpdator` / `Updator()` | `Updator func(db, obj, id, ctx) error` |
| `SetDeleter` / `Deleter()` | `Deleter func(old func() error, db, obj, id, cascade, ctx) error` |
| `SetFinder` / `Finder()` | `Finder func(db, obj, ctx) (any, error)` |

## Modes

`Mode` bits: `Search`, `Create`, `Fetch`, `FetchTitle`, `Update`, `Delete`,
`DeletedRelated`; sets `Read`, `Write`; slice `Modes`. Methods: `Is`, `Has`,
`Split`.

## Callbacks

| Symbol | Purpose |
| --- | --- |
| `Callback` | `func(state *CallbackState) error`. |
| `Callbacks[T]` | A pre/post hook list; `Pre`, `Post`, `Merge`, `Clone`, `Build`. |
| `CallbacksRegistrator[T]` | Per-mode registries; `With{Search,Fetch,Create,Update,Delete,Read,Write,Mode,ModeSplit}Callbacks`, `ModeCallbacks`, `Dot`. |
| `CallbackSlice` | Ordered callbacks; `Execute(state)`. |
| `(*DataOperatorBuilder) GetCallbacks(mode, ctx)` | Merge builder + context callbacks for a mode. |
| `AddCallbacksToContext(ctx, cb…)` / `GetContextCallbacks(ctx)` | Inject/read per-request callbacks. |
| `NamedCallbacksRegistrator` / `NamedCallbacksRegistratorOf(mb)` | Named registries stored on a model builder. |

## CallbackState

| Field / method | Purpose |
| --- | --- |
| `Obj`, `Ctx` | record + event context. |
| `DB`, `SharedDB`, `CommonDB` | operation-scoped / transaction / plain handles. |
| `Set`, `Get`, `GetOk` | pass data between pre and post. |
| `Done(func() error)` | defer work to pipeline end. |

## Has-many associations

| Symbol | Purpose |
| --- | --- |
| `SaveHasManyAssociation(field string, fieldFormKey ...func(string, *CallbackState) string) *SaveHasManyAssociationBuilder` | Build the reconciler. |
| `(*SaveHasManyAssociationBuilder) Build(ob) *DataOperatorBuilder` | Install on create+update pipelines. |
| `Updator(func(db, r) error)` | Override how existing children are updated. |
| `Pre(...)` / `Post(...)` | Extra callbacks around the reconciliation. |

### Auditing

| Symbol | Purpose |
| --- | --- |
| `AssociationAuditor` | `LogCreated` / `LogUpdated` / `LogDeleted` — child audit hooks. |
| `ContextWithAssociationAuditor(ctx, a)` | Inject an auditor (done by the activity module). |
| `AssociationAuditorFromContext(ctx)` | Read it back (nil ⇒ no logging). |

## Schema helpers

| Symbol | Purpose |
| --- | --- |
| `Schema` / `Field` | Lightweight model schema over a gorm model. |
| `NewSchema(db, model)` | Build one. |

## Misc

| Symbol | Purpose |
| --- | --- |
| `SetTableNameClauseBuilder` | A gorm clause builder that rewrites the table name of a statement. |
| `CountingConfigKey` | DB config key marking a counting query. |
