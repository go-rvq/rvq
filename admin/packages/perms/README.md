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

- an **All / Trash** filter-tab pair (the Trash tab only shows for subjects
  allowed the trash verb);
- the unscoped trash query and a **restore** bulk action, both backend-gated;
- an optional `OnRestore` hook to audit restores.

```go
perms.SetupTrash(mb, db, perms.TrashOptions{
    TableName: models.Banco{}.TableName(),
    OnRestore: auditRestores, // optional
})
```

## i18n

Messages ship in en-US and pt-BR and are registered automatically on the
builder when either feature is installed.
