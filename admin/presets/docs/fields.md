# Fields

The fields layer decides which fields a resource shows, how each renders, and how
it is written back.

## FieldsBuilder & FieldBuilder

A `FieldsBuilder` holds the ordered `FieldBuilder`s and the layout. `Listing`,
`Editing`, `Creating` and `Detailing` each embed one.

```go
e.Field("Title").
    ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent { ... }).
    SetterFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) error { ... }).
    SetI18nLabel(func(ctx context.Context) string { return "..." })
```

Key `FieldsBuilder` calls: `Field(name)`, `Only(...)`, `Except(...)`,
`HiddenField(...)`, `FieldNames()`, `CurrentLayout()`. The layout can be flat
names, rows (`[]string`) or `*FieldsSection` groups — groups produce nested
column headers in tables/listing.

## Layout & sections

```go
e := m.Editing(&presets.FieldsSection{
    Rows: [][]string{
        {"FirstName", "LastName"},
        {"Email"},
    },
})
```

`FieldsSection{Name, Title, Rows}` groups fields; nested sections nest headers.
`FieldsLayout` is the `[]any` of `string` / `[]string` / `*FieldsSection`;
`FieldsLayout.Names()` flattens it to leaf names, and
`FieldBuilders.FieldTreeLayout(...)` builds the header tree.

## FieldContext

Passed to every component/setter. Important fields:

| Field | Meaning |
| --- | --- |
| `Obj` | the record |
| `Name` / `FormKey` | field name / full form key (`Parcelas[2].Valor`) |
| `Value()` / `StringValue()` | the current value |
| `Mode` | `FieldModeStack` (LIST / DETAIL / EDIT / NEW) |
| `Label` / `InputLabel()` | display label / label to put on the input |
| `Errors` | validation errors for the field |
| `ReadOnly`, `Required`, `Disabled` | render flags |
| `MustInput` | render only the bare input (no label/hint/container) |
| `Nested` | nested config (`NestedSlice` / `NestedStruct`) |

## Default component funcs

`field_component_funcs.go` provides the built-in components, wired by field type
in `FieldDefaults.builtInFieldTypes()`:

| Mode | Defaults |
| --- | --- |
| WRITE (edit/new) | `TextFieldComponentFunc`, `NumberComponentFunc`, `CheckboxComponentFunc`, `DateTimeComponentFunc`, … |
| DETAIL | `ReadonlyComponentFunc`, `CheckboxReadonlyComponentFunc`, `DateTimeReadonlyComponentFunc` |
| LIST | `TDStringComponentFunc`, `TDReadonlyBoolComponentFunc` (already `<td>`) |

Wrappers: `FieldComponentWrapper` (adds a label container),
`FormFieldComponentWrapper`, `FieldWithHint` (adds the hint).

## MustInput — bare inputs for table cells

Set `FieldContext.MustInput = true` to render **only the input**, with no label,
hint or surrounding container — for placing a field inside a table cell (the
column header carries the label). The default component funcs honour it:

- `InputLabel()` returns `""` when `MustInput` is set, so the vuetify `.Label(...)`
  and the custom select labels render nothing;
- `FieldWithHint` skips the hint;
- `FieldComponentWrapper` skips the label container.

This is what the list-editor table renderer uses for its cells — see
[list-editor.md](list-editor.md).

## Nested fields

A field can nest another fields set:

- `NestedStruct(mb, fb)` — an embedded struct.
- `NestedSlice(mb, fb)` — a has-many slice, edited with the list editor.

See [list-editor.md](list-editor.md).
