package basics

import (
	"github.com/go-rvq/rvq/admin/docs/docsrc/generated"
	. "github.com/theplant/docgo"
	"github.com/theplant/docgo/ch"
)

var History = Doc(
	Markdown(`
The **history** package (~github.com/go-rvq/rvq/admin/packages/history~) adds
git-like revision history to any model. On every save it records a **revision**
of the versioned fields — a content hash (so identical states dedupe), the
author, a parent→child chain, a publish tag and a per-revision access counter.
It provides, through the admin UI:

* a **comparator** — a field-by-field diff of two revisions (or a revision vs the
  current record), with Prism syntax highlighting, line-by-line changes and the
  changed content highlighted **within each line** (GoLand/IntelliJ style);
* **navigation by field** (a field's own revision timeline); and
* **revert** — whole record, named fields, one field's whole content, or only the
  **selected hunks** of a field (git ~checkout -p~), including for singleton
  models.

## Activate history on a model

Activate it per model after ~history.Configure(pb, db)~ (once per admin):

~~~go
history.Configure(pb, db)
history.New(db).Model(mb).Build()            // versions the EDIT fields
history.New(db).Model(mb).Fields("Body").Build()  // only these fields
~~~

~HTMLFields~ marks a field as HTML (block-level diff / partial revert);
~WholeFields~ marks structured fields (Cover, relations) as whole-only;
~ExtraFields~ versions a field that is not an inline form field.

## Diff and partial revert of a structured (JSON) field, as YAML

A field whose value is a JSON map is not a string, so by default it is
whole-only. Two hooks adapt it to the text-based comparator and partial revert:

* ~FieldDiffHandler(field, fn)~ overrides how a field is diffed — here it renders
  the JSON payload as **YAML** with ~PrismDiffComponents("yaml", …)~.
* ~FieldContentHandler(field, language, codec)~ registers a **content codec**
  (~ToText~ / ~Apply~). It is what lets a structured field take part in partial
  revert: the field is diffed and edited as text (YAML), and the patched text is
  parsed back into the map on apply. The changed spans become clickable hunks in
  the comparator (check + super-highlight), and only the selected ones are
  reverted.

The example below activates history on a singleton config model whose whole
configuration is one JSON map, and wires both hooks so the comparator and the
partial revert operate on it as YAML:
`),
	ch.Code(generated.HistoryYAMLConfigExample).Language("go"),
	Markdown(`
The same pattern fits any settings/config model: version the JSON payload, show
it as YAML in the comparator, and let the user revert just the YAML hunks they
pick. For a plain string field (text or HTML) no codec is needed — partial revert
is available out of the box.
`),
).Title("History").Slug("basics/history")
