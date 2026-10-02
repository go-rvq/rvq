# Edit: {%= doc.model %}

The form that changes a record. Open it with the pencil (**Edit**) of its
detail, or from the menu of its row (**⋯**) in the listing; what is changed is
kept when the form is saved.

{%= admin.fields("edit") %}

A field may be required, or checked when the form is saved: what is wrong is
said under the field, and nothing is saved until it is fixed. A field *shown
depending on the record* appears only for the records it applies to.
