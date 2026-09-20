# RVQ documentation

RVQ unifies [rvq/web](https://github.com/go-rvq/web),
[rvq/x](https://github.com/go-rvq/x) and
[rvq/admin](https://github.com/go-rvq/admin) with local changes. This is the
top-level index of the in-tree documentation.

## Admin

- [Developing a package](../developer.md) — the convention for anything that
  installs models into a host: tables carry the package prefix, and so does the
  registration id.
- [Admin package reference](../admin/docs/reference.md) — index of the current
  per-package docs (presets, packages, features).
  - [presets](../admin/presets/README.md) · [presets/docs](../admin/presets/docs/README.md)
    ([Record IDs](../admin/presets/docs/record-id.md))
  - [packages/perms](../admin/packages/perms/README.md) — per-record permission
    manager and trash.
  - [activity](../admin/activity/README.md) — auditing and the activity log.

> The RVQ docsite under `admin/docs` (`README.md`, `docsrc/`) is an older
> upstream version, pending an update; prefer the per-package docs linked above.

## Modules overview

RVQ is a monorepo of three areas:

- **web** — the Go HTML builder and the Vue/Plaid event runtime.
- **x** — cross-cutting libraries (`i18n`, `perm`, `login`, UI kits, …).
- **admin** — the presets admin framework and its packages/features (see the
  [admin package reference](../admin/docs/reference.md)).
