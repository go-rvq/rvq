# 3. Editing

A change is saved only when it may be:

- **Another person saved first**: the form remembers the record as it was when
  it opened. If someone saved it in the meantime, saving is refused, saying
  who and when — so nobody's work is overwritten. Reload the form and make
  the change again.
- **A required or invalid field**: what is wrong is said under the field, and
  nothing is saved until it is fixed.
- **A deleted record** is not edited: see
  {%= admin.doc("guides/policies/02-trash").link %}.

What is never edited:

- **Records kept as history** — revisions, activity logs —: they are only
  looked at.
- **Settings** have a single record: they are edited, never created or
  deleted.
- A **published** page or post keeps showing on the site what was published:
  a change goes there only when it is published again
  ({%= admin.doc("actions/Publisher").link %}).
