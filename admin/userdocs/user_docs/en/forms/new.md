# New record: {%= doc.model %}

The form of a new record. Open it with **New** at the top of the listing; the
record is made when the form is saved, and its detail opens.

{%= admin.fields("new") %}

A field may be required, or checked when the form is saved: what is wrong is
said under the field, and nothing is saved until it is fixed. A field *shown
depending on the record* appears once the fields it depends on are filled.
