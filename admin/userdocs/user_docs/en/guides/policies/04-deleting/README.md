# 4. Deleting

- A deletion is **confirmed** first: the dialog says what is deleted.
- **Records that depend on it**: the dialog lists them (**Show related
  items**), and the deletion may take them along (**Delete related items**).
- Where the listing has a trash, the record goes there and may be restored
  ({%= admin.doc("guides/policies/02-trash").link %}); elsewhere the deletion is
  for good.
- What may not be deleted is not offered: the settings, a record kept as
  history, a deleted record.

See {%= admin.doc("actions/Delete").link %}.
