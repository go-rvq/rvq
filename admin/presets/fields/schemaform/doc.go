// Package schemaform builds an edit form out of a gad interface.
//
// A value edited as a form is described by a SCHEMA — the gad interface its
// content satisfies, written without the keyword:
//
//	[]{label str; icon str; href}
//
// which is `interface []{label str; icon str; href}`: a list of records of
// three fields. The form the schema describes is what a person edits; the value
// itself is stored the way it always was — YAML, for a LocaleMessage — and only
// the editor changes.
//
// # The shapes a schema takes
//
//	{title str; count int}           a record: one form
//	[]{label str; href}              a list of records: the array sorter
//	[]str                            a list of PLAIN VALUES: one input per item
//	[][]str                          a list of those
//	{sub: {x str}}                   a field that is a form of its own
//	{sub: []{x str}}                 a field that is a list of forms
//	{tags []str}                     a field that is a list of plain values
//	{name?: {x str}}                 the nested field may be nil or absent
//
// `[]str` is the short form of `[]<str>`; `sub: {…}` the short form of
// `sub interface {…}`. A list of plain values has no Fields and one Item: the
// item IS the value, so it binds by index and carries no label of its own — the
// list already has one.
//
// # The schema is a gad program
//
// A schema is not parsed into a form: it RUNS. It is compiled and run in a gad
// VM of its own (Builder.Parse), and the interface named `Form` is what the run
// returns; the form is then read off that interface through gad's own
// reflection — its fields and their resolved types, its array depth
// (`@depth`), the element of a list of values (`@elem`), the `[k=v, …]`
// metadata of the interface and of each field (`@meta`). So a schema may use
// whatever gad can declare, and is checked by gad the way any gad code is: a
// type name nothing declared is an `unresolved reference`, reported where the
// schema is read — not later, in the form.
//
// The types a field may use are the builder's: every type registered on it is
// in scope while the schema runs, so `color` or `html` resolve although gad has
// no such types — and `time`, which in gad is the time NAMESPACE, is the
// builder's time type here. A name gad already has as a type (`str`, `int`,
// `bool`) stays gad's.
//
// The run is bounded (RunTimeout), and what a schema read into is cached by the
// source and the builder's types: a schema is read on every draw and every
// save, and it is the same declaration each time.
//
// # What a field's type means
//
// The type NAMES THE COMPONENT that edits the field: `Builder.Type(name, comp)`
// registers one, and a schema may only use the types the builder knows. A field
// the schema left untyped is `str`, plain text — an omission, not a type of its
// own, which is why a dynamic list of values may take it over (see below).
//
// A new Builder knows: str (and untyped), text (textarea), html (rich text),
// int, uint, float, decimal, bool (switch), color, date, time, duration, and
// FormType — the form itself, one level down, which is what makes the schema
// recursive.
//
// # A schema written in parts
//
// A named interface is a declaration another may use BY ITS NAME, and the one
// named `Form` (FormName) is the form:
//
//	interface User { name, id }
//	interface Form { owner User; creator User }
//
// As anywhere in gad, a name is in scope after its declaration, so what the
// form uses is declared before it — which is also why an interface that
// contains itself cannot be written at all. Each occurrence is read on its own,
// so the words of `owner` are not the words of `creator`.
//
// Written in one piece the schema needs no name for anything, and needs none —
// `{…}`, `[]{…}`, `[]str`, or `interface []{…}` after the enums it uses: that
// one interface is the form. An interface without a name BESIDE the form is
// refused: nothing could refer to it.
//
// # Metadata: how the form is drawn
//
// The `[k=v, …]` block before the interface says how the form is DRAWN, never
// what the value is (Schema.Meta; a field's own block is Field.Meta):
//
//		[layout="table", columns=[#label, #icon, #color, #disabled]]
//		interface Form []{label str; icon str; color str; link? str; disabled? bool}
//
//	  - `layout` — how a list of records is drawn: "form" (the default: each
//	    record a form, one under the other) or "table" (one row per record, one
//	    column per field). A table is only of a list of records.
//	  - `columns` — the fields a table shows, in the order it shows them. Symbols
//	    name them (`#label` is the string "label"). A field left out keeps its
//	    value — it is in the record, only not drawn. Without it, every field, in
//	    the schema's order.
//
// What cannot be drawn is refused where the schema is read: a layout that does
// not exist, a table of something that is not a list of records, a column that
// names no field.
//
// # Enums, and lists of values the schema cannot know
//
// A field may hold one of a CLOSED LIST of values, and is then a select. The
// list comes from one of two places:
//
//   - an `enum` the schema declares, beside it or inline —
//     `enum Perm { Read, Write }` / `{perm enum { Read, Write }}`. What goes
//     into the record is the member's NAME (`Read`), never the number behind
//     it, because the value is written out and read back by name.
//   - `Builder.EnumItems`, for a list that is not fixed and so cannot be
//     written in a schema at all: the locales in the database, the categories
//     of this account. It answers by the field's path, with the request at
//     hand, and may fail — a list that should have been there and could not be
//     fetched is drawn as the failure it is, and refuses the save rather than
//     letting through a value nobody checked.
//
// Both are resolved in ONE place (Builder.itemsOf), so what the select offers
// is exactly what a posted value is checked against. Precedence, when a field
// could be drawn in more than one way: a component registered for the field's
// type wins — that is what registering a type is for, an enum's own name
// included — except for a field the schema left untyped, where the list wins,
// since `str` is what a field has when the schema said nothing about it.
//
// # Paths: how a field is named
//
// A field is asked about by its PATH — the names from the root down. A list
// adds no name of its own, so a field of the records in `links` is
// `links.href`, not `links.*.href`, and the item of a list of plain values is
// the list's own path. The same path names a field for its words
// (Builder.FieldInfo), for its values (Builder.EnumInfo, Builder.EnumItems) and
// for Schema.FieldAt.
//
// # The words around a field
//
// A schema says nothing about what a field is CALLED. Builder.FieldInfo answers
// that, by path: the label (which falls back to the field's name, humanized),
// the hint under the input, and the help behind a `?` at its right.
// Builder.EnumInfo does the same for the values of a declared enum.
//
// # Drawing, and reading back
//
// Builder.ComponentFunc(schema) is a presets field component: it binds the
// value under the field's own form key (`form["<FormKey>"]`), a record by name
// and a list through the array sorter — the list iterating in the sorter's
// DEFAULT slot (the sorter draws rows of its own only while sorting), with a
// button per row to remove it and one at the end to add another.
//
// In a table each column's header carries the field's label, with its hint
// behind it, and a cell carries only the input, on one line: the component is
// given Context.Compact, and spreads Context.CompactAttrs on its input.
//
// What the browser posts is flat, by the two rules the binding follows: an
// array indexes and a record names, so a list of records arrives as
// `key[0].label` and a list of values as `key[0]`. Schema.Decode reads that
// back — typed by the schema, records keeping the schema's order (Record) —
// and Builder.DecodeForm reads AND checks it, which is the one that has the
// request and so the values a field may hold.
//
// Builder.Codec says how the decoded value is carried outside the form:
// EncodeValue/DecodeValue, YAML unless an application says otherwise
// (JSONEncode/JSONDecode are here for a value stored as JSON).
//
// # A whole use, end to end
//
//	schema, err := schemaform.Parse(`[]{label str; icon str; href}`)
//	if err != nil {
//	    return err
//	}
//
//	b := schemaform.New().
//	    FieldInfo(func(ctx *web.EventContext, path string) schemaform.FieldInfo {
//	        w := words[path] // whatever the application calls this field
//	        return schemaform.FieldInfo{Label: w.Label, Hint: w.Hint}
//	    }).
//	    EnumItems(func(ctx *web.EventContext, path string, f *schemaform.Field) ([]schemaform.EnumItem, error) {
//	        if path != "icon" {
//	            return nil, nil
//	        }
//	        return iconsOf(ctx) // a list only the request knows
//	    })
//
//	// drawing: the form over the value the record holds
//	mb.Editing().Field("Links").ComponentFunc(b.ComponentFunc(schema))
//	initial, err := b.DecodeValue(record.Links) // the stored text -> the form's value
//
//	// reading back, on save
//	value, err := b.DecodeForm(ctx, schema, ctx.R.Form, "Links")
//	if err != nil {
//	    return err // a value no select offered, or a list that could not be fetched
//	}
//	record.Links, err = b.EncodeValue(value)
package schemaform
