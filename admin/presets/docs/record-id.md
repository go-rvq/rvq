# Record IDs

Every model registered in a presets admin has a **record id**: the string form
of its primary key(s), used in URLs, permission resources, actions and nested
routes. This document explains how ids are encoded, parsed and how to support
custom primary-key types.

## The `ID` type

`presets.ID` (an alias of `model.ID`) holds the primary-key field(s) and their
values for one record:

```go
type ID struct {
    Schema Schema
    Fields model.Fields
    Values []any
}
```

Key methods:

| Method | Purpose |
|---|---|
| `String()` | Slug form: the values joined by `_` (empty when zero). |
| `StringValues()` | The values as strings. |
| `IsZero()` | True when every value is the zero value. |
| `SetTo(obj)` | Writes the values back into `obj`'s primary-key fields. |
| `GetValue(field)` | The raw value of a named field. |

A single-key record with id `7` has `String() == "7"`; a composite key
`{Page, "en"}` becomes `"7_en"`.

## Parsing an id from a string

`ParseRecordID(schema, s)` turns a slug back into an `ID`, converting each part
to the primary-key field's Go type. The `ModelBuilder` wrappers are the usual
entry points:

```go
id, err := mb.ParseRecordID(s)   // -> presets.ID
err := mb.ParseRecordIDTo(&obj, s) // parse and SetTo(&obj) in one step
```

`ParseRecordIDTo` is handy in actions and services where you have the record id
string and want the primary key filled into a fresh model value:

```go
var m models.Movimentacao
if err := mb.ParseRecordIDTo(&m, id); err != nil {
    return err
}
// m.ID is now set
```

## Supported primary-key types

`ParseRecordID` resolves each field by kind, in this order:

1. **Basic kinds** — `string`, the signed/unsigned integer widths.
2. **`uuid.UUID`** — parsed from its canonical string form.
3. **A `Parse(string) (T, error)` method** — any type (value or pointer
   receiver) exposing such a method is parsed with it. This is the extension
   point for custom id types that ship their own parser.
4. **`sql.Scanner`** — the value is `Scan`-ed from the string; the scanned
   value (not the pointer) is stored so `ID.SetTo` can assign it without a
   reflection panic.
5. **`IdParserFallback`** — a globally registered parser (see below), tried
   last.

If none matches, parsing fails with an `Unsupported type` error.

### Custom id types

Prefer a `Parse` method on the type:

```go
type Code [4]byte

func (Code) Parse(s string) (Code, error) { /* ... */ }
```

When you cannot add a method (e.g. a third-party type), register a fallback
parser once at startup:

```go
presets.IdParserFallback[reflect.TypeOf(SomeID{})] =
    func(s string) (any, error) { return parseSomeID(s) }
```

The registered function returns the concrete value (not a pointer); it is used
only when the built-in resolution above does not apply.

## Why the pointer matters

`ID.SetTo` assigns each parsed value into the model's field with
`reflect.Value.Convert`. A parser must therefore yield the **value** type
(`uuid.UUID`), never a pointer (`*uuid.UUID`) — converting `*uuid.UUID` to
`uuid.UUID` panics. The built-in `sql.Scanner` path dereferences the scanned
pointer for exactly this reason; custom parsers should do the same.
