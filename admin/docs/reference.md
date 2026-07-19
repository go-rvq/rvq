# RVQ admin — package reference

Index of the current, in-tree documentation for the `admin` packages of RVQ.

> Note: the QOR5 docsite in this folder (`README.md`, `docsrc/`) is an older
> upstream version and is pending an update. The links below point to the
> up-to-date per-package docs that ship with the source.

## Core

- [presets](../presets/README.md) — the admin builder (models, listing,
  detailing, editing, actions, permissions).
  - [presets/docs](../presets/docs/README.md)
    - [Record IDs](../presets/docs/record-id.md) — id encoding/parsing,
      supported primary-key types and the `IdParserFallback` extension point.

## Packages

- [packages/perms](../packages/perms/README.md) — reusable per-record
  permission manager and soft-delete trash.
- [packages/db-tools](../packages/db-tools/README.md) — database backup/restore
  tools.
- [packages/fs-tools](../packages/fs-tools/README.md) — filesystem tools.
- [packages/mail-sender](../packages/mail-sender/README.md) — mail sending.
- [packages/helper](../packages/helper/README.md) — model selectors and shared
  admin helpers.

## Features

- [activity](../activity/README.md) — activity log / auditing (creator, time,
  request origin IP & browser, per-field diffs, per-record log view).
- [seo](../seo/README.md) — SEO metadata.
- [slug](../slug/README.md) — slug fields.
- [worker](../worker/README.md) — background jobs.
- [media/base](../media/base/README.md), [media/vips](../media/vips/README.md) —
  media library.
