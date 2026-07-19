# presets docs

Reference notes for the `presets` admin package.

- [Record IDs](record-id.md) — how a model's primary key is encoded into a
  record id, parsed back (`ParseRecordID` / `ParseRecordIDTo`), the supported
  key types (basic kinds, `uuid.UUID`, a `Parse` method, `sql.Scanner`) and the
  `IdParserFallback` extension point for custom id types.
