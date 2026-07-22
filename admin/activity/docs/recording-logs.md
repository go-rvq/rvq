# Recording logs

There are two ways logs get written: automatically for registered preset models,
and manually through the `Add*Record` API.

## Automatic (preset models)

When a `*presets.ModelBuilder` is registered, the module wraps its editing,
creating and deleting flow:

- **Create** — on the creating builder's `CreateFunc`, and on `SaveFunc` when
  the id is zero (or no existing row is found). Records `ActivityCreate`.
- **Edit** — on `SaveFunc` when the id is non-zero and a previous row exists.
  Records `ActivityEdit` with a field diff of *old → now*.
- **Delete** — on the listing's `DeleteFunc`. The old row is fetched first (so
  its keys/label are known) and `ActivityDelete` is recorded.

No call-site changes are required. Use `SkipCreate` / `SkipUpdate` /
`SkipDelete` on the model builder to opt out of any of these (see
[model-configuration.md](model-configuration.md)).

> An edit log is only persisted when the diff is non-empty — an update that
> changes nothing (after ignored fields) writes no row.

## Manual — context-based

These resolve creator, db and request info from the context (see
[getting-started.md](getting-started.md#4-wire-creator-db-and-request-info-through-the-context)).
Prefer them when you already have a request/context on hand.

```go
ctx := activity.ContextWithCreator(
    activity.ContextWithRequestInfo(r.Context(), r), "alice")

ab.AddRecords(activity.ActivityCreate, ctx, obj)      // Create / Delete / View: stored as-is
ab.AddRecords(activity.ActivityEdit, ctx, obj)        // Edit: fetches old from db, stores the diff
```

`AddRecords(action, ctx, vs...)` dispatches each value to the `ModelBuilder`
registered for its type. For `ActivityEdit` it loads the current row from the
database (by primary key) and diffs it against `obj` — so call it **before**
saving your changes, while the DB still holds the old values, or use the
`WithOld` variants below.

Other context-based helpers:

- `ab.AddCustomizedRecord(action, diff, ctx, obj)` — record an arbitrary action
  string (e.g. `"Approved"`, `"Paid"`). With `diff=true` it fetches the old row
  and stores a diff; with `diff=false` it stores the action with no diff.
- `mb.AddEditRecordWithOldCtx(ctx, old, now)` — edit log from an *old* snapshot
  you already hold, avoiding a re-fetch.

## Manual — explicit creator + db

When you don't have a wired context (background jobs, tests), pass creator and
`*gorm.DB` explicitly. `creator` may be a `string` or a `CreatorInterface`.

```go
ab.AddCreateRecord(creator, obj, db)
ab.AddEditRecord(creator, obj, db)             // fetches old from db, diffs
ab.AddEditRecordWithOld(creator, old, obj, db) // diff a snapshot you already hold
ab.AddDeleteRecord(creator, obj, db)
ab.AddViewRecord(creator, obj, db)
ab.AddSaveRecord(creator, obj, db)             // create if no old row exists, else edit
```

## Multiple builders for the same type

`ab.Add*` dispatch by Go type. If one struct type is registered under several
`ModelBuilder`s, disambiguate by recording on the specific one:

```go
mb := ab.MustGetModelBuilder(postPresetModel)
mb.AddRecords(activity.ActivityEdit, ctx, post)
```

## Nested records (has-many children, custom flows)

The automatic path logs the **aggregate** model that the admin saves. Children
persisted through an association (or any record saved outside the preset save
flow) are **not** logged unless you record them yourself.

Register the child model and record each affected child as you create / update /
delete it, threading the same context so the creator and request origin match
the parent action:

```go
ab.RegisterModel(&OrderItem{})

// inside your save/association logic, with `ctx` carrying creator + request info:
ab.AddRecords(activity.ActivityCreate, ctx, newItem)
ab.AddEditRecordWithOldCtx(ctx, oldItem, updatedItem)
ab.AddRecords(activity.ActivityDelete, ctx, deletedItem)
```

Pass a transaction DB through `ContextWithDB(ctx, tx)` so the log rows are
written inside the same transaction as the data change.
