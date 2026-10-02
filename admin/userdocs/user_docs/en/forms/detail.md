# Detail: {%= doc.model %}

What a record shows when it is opened from the listing.

{%= admin.fields("detail") %}

## The menu of the record

The menu (**⋯**) at the top of the detail: what is nested in the record, what
may be done with it, and its pages. Each one is shown only to whoever may use
it; one available depending on the record says when.

{%= admin.menu() %}
