# presets docs

Reference notes for the `presets` admin package. See the package
[README](../README.md) for the high-level overview.

- [Getting started](getting-started.md) — the `Builder`, registering models with
  `Model`, installing plugins with `Use`, and wiring the data operator.
- [Models](models.md) — `ModelBuilder` and its `Listing` / `Editing` /
  `Creating` / `Detailing` builders.
- [Fields](fields.md) — `FieldsBuilder`, `FieldBuilder`, `FieldContext`, the
  default component funcs, and `MustInput` (bare inputs for table cells).
- [List editor & tables](list-editor.md) — nested slices, the client-side list
  editor (`__pos`/`__new`/`__deleted`), and the table renderer
  (`ListEditorTableBuilder`).
- [Form host](form-host.md) — how a page opens and destroys its edit/create
  overlays with a single variable (`VAR.show = true|false`): the NEW / row
  DETAIL-EDIT / detailing EDIT flows, the dialog & drawer closer, the re-renders
  on validation errors and list-editor add/remove, and driving the form from your
  own buttons.
- [Post-save refresh](post-save-refresh.md) — what is refreshed after a
  successful save (`onSaveCallbacks`): the listing behind a NEW, the detail and
  the listing that opened it after an EDIT, the detail PAGE and its `<title>`
  without a reload, and the singleton cases.
- [Record IDs](record-id.md) — how a model's primary key is encoded into a
  record id, parsed back (`ParseRecordID` / `ParseRecordIDTo`), the supported
  key types (basic kinds, `uuid.UUID`, a `Parse` method, `sql.Scanner`) and the
  `IdParserFallback` extension point for custom id types.
