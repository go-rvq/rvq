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
// A field declared `get name Type` — a getter — is SHOWN, never edited:
//
//	[]{get path str; content text}
//
// is a list of records whose path is read-only and whose content is a
// textarea. The getters come before the other fields (gad keeps an
// interface's properties apart from its fields). The form draws a getter
// read-only, and Schema.KeepReadOnly puts back the stored value over whatever
// a post said of it. A property with a setter, a `prop`, describes behaviour,
// and is not in the form.
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
// As anywhere in gad, the declarations may come in any order: the form may use
// an interface declared after it. An interface that contains itself — itself,
// or through others — describes a form without end, and is refused where the
// schema is read. Each occurrence is read on its own, so the words of `owner`
// are not the words of `creator`.
//
// Written in one piece the schema needs no name for anything, and needs none —
// `{…}`, `[]{…}`, `[]str`, or `interface []{…}` after the enums it uses: that
// one interface is the form. An interface without a name BESIDE the form is
// refused: nothing could refer to it.
//
// # Groups: fields side by side
//
// A group — `{ … }` written where a field goes, gad's sugar for the field
// `$N` of an anonymous class (interface) — is no field of the form: its
// fields are, in its place, the COLUMNS of one row, and a new group is a new
// row. What they hold is the record's own: the group adds no name to the
// value, nor to a field's path.
//
//	{
//	    { [width=2] name str; phone str; email? str }   // one row: 40% 30% 30%
//	    message text                                    // a row of its own
//	}
//
// A field's `[width=N]` is its share of the row in fifths (WidthUnits): 1 is
// 20%, and the fields with none share what is left, equally. The row is
// whole — the widths given are whole numbers from 1 to 5, add up to 5 when
// every field says one, and leave a share for the ones that do not —, or the
// schema is refused where it is read. Field.Row is the row a field is a
// column of, Field.Percent its share; Rows gives the fields by their rows.
// On a narrow screen the columns stack. A group inside a group is refused: a
// row has columns, not rows (gad nests them; a form does not). A group's `?`
// says nothing here: each field says whether it is required.
//
// A group's metadata says what else it may be (Field.Group): `[tabs] { … }`
// a tab each field — its label the tab's, its `[icon=…]` beside it —;
// `[steps] { … }` a step each, with back and next, `[steps, vertical]` down
// the page. Their fields take the whole width: no `width` is read.
//
// # Metadata: how the form is drawn
//
// The `[k=v, …]` block before the interface says how the form is DRAWN, never
// what the value is (Schema.Meta; a field's own block is Field.Meta):
//
//	[layout={name: "table", columns: [#label, #icon, #color, #disabled]}]
//	interface Form []{label str; icon str; color str; link? str; disabled? bool}
//
// `layout` is a layout's name, or `{name: …}` with what the layout is told.
// A layout is registered (RegisterLayout) with its Config — a gad class,
// written as code, whose fields are what it may be told, typed, with their
// defaults — and the ones there are by default:
//
//   - "form" (the default, the list): each record a card of its form, one
//     under the other. It is told nothing.
//   - "table": one row per record, one column per field.
//     `class Config { [fields=true, empty_as_all, sorted] columns []Field }`
//     — the fields it shows, in order (symbols name them: `#label` is the
//     string "label"); none, every field. A field left out keeps its value —
//     it is in the record, only not drawn.
//   - "grid": each record a card, `columns` of them side by side (on a narrow
//     screen they stack). `class Config { columns int = 4 }`, 1 to 12.
//
// A table and a grid are only of a list of records. The metadata a Config's
// field may carry: `fields=true` (it names fields of the list — a `Field`, or
// a list of them), `empty_as_all` (an empty list of fields is all of them),
// `sorted` (the order it names them in is the order drawn).
//
// What cannot be drawn is refused where the schema is read: a layout that does
// not exist, a table or a grid of something that is not a list of records, a
// config the layout's Config does not declare or of another type, a column
// that names no field, a grid out of range.
//
// A field's own `when` draws it only while the record holds what it says —
// `[when={name: "grid"}] columns? int` is there only while `name` is "grid".
//
// # Enums, and lists of values the schema cannot know
//
// A field may hold one of a CLOSED LIST of values, and is then a select. The
// list comes from one of three places:
//
//   - the field's own `[options=…]` (MetaOptions), its values with their
//     labels in the order written — a key-value array
//     `[options=(;opt1="option 1")] value str`, an array of pairs
//     `[options=[[1, "one"], [2, "two"]]] n int`, or of values alone
//     `[options=["a", "b"]]`. The VALUE goes into the record (as text, a
//     number for a number field), the label is only shown. It is a select
//     whatever the field's type — it wins over a registered component —, and
//     on a list of plain values (`tags []str`) each item is one.
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
// # Showing a value: the detail and the listing cell
//
// What is not edited is SHOWN, with the same schema and the same words:
//
//	value, _ := b.DecodeValue(stored)
//	b.DetailComponent(ctx, schema, value) // the detail page
//	b.ListComponent(ctx, schema, value)   // a listing cell
//
// Nothing is bound: the value itself is drawn. In a detail a record is its
// fields under their labels (an empty one as a dash), a list of plain values a
// bulleted list, a list of records each record with a line between them — and
// `[layout="table"]` a table of the columns it names. A listing cell keeps to
// what fits there: a record on one line, `Label: value · …`, without the empty
// fields; a list of plain values separated by commas; a list of records one
// line per record, showing the field a reader names it by (`label`, `title`,
// `name`, else the first text field); a list shows ListMaxItems and says how
// many more.
//
// Each type is shown by a component of its own, registered like the editing
// one — Builder.Display(name, f), read with Builder.DisplayFunc — which is
// given the value in Context.Data (not an expression in Context.Value) and
// Context.Compact in a cell. The defaults: text for str and the numbers,
// date, time and duration; `text` keeping its line breaks; `html` sanitized
// (only its text in a cell); `bool` a check or a dash; `color` a swatch; an
// enum, or a field EnumItemsFunc answers for, the LABEL of the value it holds.
// A type with no display shows as text: a read view never refuses a value.
//
// # A class as the form
//
// Builder.Class reads a gad class — options an application declares in code,
// as a layout does — as a form: a record of its fields, their defaults and
// metadata kept.
//
//		[label="Grade"]
//		class gridLayout { [label="Colunas", hint="de 1 a 6"] columns int = 3 }
//		class listLayout { comp? enum { index_list, compact } }
//		class PageOptions {
//		    postType? models.PostType          // a type only the application knows
//		    posts class {                      // a record of its own: a group
//		        layout? listLayout|gridLayout  // a choice of one class
//		    }
//		}
//
//	  - A field typed by a class — `a class { … }`, or a class by name — is a
//	    record of its own, drawn as a group.
//	  - A field typed by a union of classes is a CHOICE (Schema.Choice): a
//	    select of the classes — by the `label` of each class's metadata, or its
//	    name — and the form of the one chosen. Its value is an object of one
//	    key, the class chosen: `{gridLayout: {columns: 3}}`. In the schema it is
//	    a record with a field per class, so a path goes through the class:
//	    `posts.layout.gridLayout.columns`.
//	  - A type only the application knows is named by Builder.TypeOf, asked
//	    first for every field of a single type: it answers the component the
//	    field is drawn with (registered with Builder.Type).
//	  - `[label="…", hint="…", help="…"]` on a field names it where the
//	    FieldInfoFunc gives it no words.
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
