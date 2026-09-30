# perms

Reusable, resource-agnostic permission tooling for presets admins: a per-record
**permission manager** (a shortcut to the `x/perm` API, backed by
`perm.DefaultDBPolicy`) and a soft-delete **trash** gated by a `trash` listing
permission. Both work with any `*presets.ModelBuilder`.

## Permission manager

`InstallManager(mb, db)` adds a permissioned **Manage Permissions** detail
action. The dialog lists the record's grants and lets an authorized subject
grant/revoke, for one or more roles/users, permission over the record:

- **whole-record** verbs — view / edit / delete;
- **per-field** grants, per mode (detail / edit), including **nested** fields;
- **detail actions**;
- **auto-perm pages** (a page with its own developer-defined verifier keeps its
  own check and is not listed; an auto-perm page enters using its path as the
  permission segment, e.g. `…:carteiras:<id>:/relatorio`).

Grants are written as `perm.DefaultDBPolicy` rows keyed by a `ReferID` and the
resource computed by the presets permissioner (record resource + glob wildcard
for the field/action/page sub-resources), then the verifier is reloaded. The
action itself is guarded by the `managePermissions` verb.

```go
perms.InstallManager(mb, db)
```

### Enumeration helpers

The dialog is built from the resource's own structure; the same helpers are
public so you can build a custom UI:

| Helper | Returns |
|---|---|
| `FieldNodes(mb, mode)` | the mode's fields as a tree (recurses into nested). |
| `FieldResource(mb, id, path…)` | a field's exact permission resource. |
| `ActionNodes(mb)` | the detail actions (name + perm verb). |
| `PageNodes(mb)` | the auto-perm listing/detail pages. |
| `PageResource(mb, mode, id, path)` | a page's permission resource. |

### Low-level primitives

`RecordResource` / `ListResource`, and `Grant` / `Revoke` / `List` over
`perm.DefaultDBPolicy`, are exported for use outside the dialog.

## Trash

`SetupTrash(mb, db, opts)` turns a soft-delete listing into a trash-enabled one,
gated by the `trash` listing permission:

- a **Trash** filter tab after the model's own tabs — after an **All** one, the
  default, when it has none. It is added when the listing renders
  (`presets.FilterTabsWrapper`), so a `FilterTabsFunc` set before or after keeps
  it; the tab only shows for subjects allowed the trash verb;
- the unscoped trash query and a **restore** bulk action, both backend-gated;
- a record of the trash opens — its detail, its history, as the parent of its
  children —: fetched deleted or not, by who may see the trash;
- an optional `OnRestore` hook to audit restores.

```go
perms.SetupTrash(mb, db, perms.TrashOptions{
    TableName: models.Banco{}.TableName(),
    OnRestore: auditRestores, // optional
})
```

`AutoTrash(b, db)`, before the models, sets up the trash of **every**
soft-deleted model (a `gorm.DeletedAt` field, listed by gorm2op, not a
singleton) — the default of a soft-delete listing.

### Who deleted a record, and from where

A model that embeds `softdelete.Deletion` in place of its `DeletedAt` keeps its
deletion: `DeletedAt`, `DeletedByID` (a foreign key to the users —
`UsersTable`, set up by `MigrateDeleted` —, nil when the user is deleted) and
`DeletedOrigin` (`origin.Origin`, the columns `deleted_origin_*`: the address,
the browser and, when the application locates addresses — `origin.SetFunc` —,
the place). The trash fills them in when the admin deletes a record, and
clears them on restore; outside the admin, set them with
`softdelete.Columns(by, origin.Of(r))` (or `perms.DeletedBy(r)`).

In the Trash tab the listing shows, after its own columns
(`ListingBuilder.AppendTrailingFields`), **Deleted on**, **Deleted by** and
**Deleted from**; the last opens the `deleted_origin` action (of the record's
detail, in the admin's dialog): what is known of where from, and the place on
a map (`origin/originui`).

The detail of a deleted record opens with a warning (`DetailingBuilder.
AppendNoticeFunc`): it was deleted, when and — as the model keeps them — by
whom and from where.

## i18n

Messages ship in en-US and pt-BR and are registered automatically on the
builder when either feature is installed.
