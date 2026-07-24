// Package helper provides higher-level building blocks on top of the presets
// admin (github.com/go-rvq/rvq/admin/presets) for editing associations from a
// parent form: nested has-many / many-to-many editors, foreign-key selectors and
// inline child editing. It is used by HERMON-CMS and IPCD to wire related models
// into a parent's edit form with a minimum of boilerplate.
//
// # Files
//
//   - nestedslice.go — NestedSlice / NestedSliceBuilder: edits a parent's
//     has-many or many-to-many field as a list of related records, with a model
//     selector for many-to-many. For many-to-many it manages the join table
//     directly (link/unlink/filter). See the key-type notes below.
//
//   - model_select.go — ModelSelectorBuilder: a belongs-to (foreign-key) selector
//     that binds a `<Field>ID` column to a searchable single-select of the
//     foreign model.
//
//   - association.go — AssociationManagerBuilder: wires a data operator for an
//     association field.
//
//   - inline_edit.go — InlineEdit / InlineEditModel: edit a child record inline
//     within the parent form.
//
//   - fields.go — field/admin-tag helpers (AdminTag, KeyValueArray).
//
//   - messages.go, errors.go, url_params.go — i18n messages, error mapping and
//     URL-parameter helpers.
//
// # Primary-key types and composite keys
//
// The many-to-many machinery in nestedslice.go no longer assumes an integer `id`
// column: the related record's primary-key column is read from the schema and
// used as its own type, so integer, UUID and other single-column key types all
// work (see m2mInsertQuery / relatedPKColumn and TestM2MLinkQuery_NonIntegerKey).
//
// COMPOSITE primary keys are NOT yet supported by these helpers:
//
//   - NestedSlice rejects a related model with more than one primary field
//     (panics with a descriptive message); the join link/unlink/filter is built
//     for a single related key column.
//   - ModelSelectorBuilder binds a single `<Field>ID` foreign key and selects by a
//     single "ID" value, so it cannot represent a composite foreign key.
//
// Supporting composite keys here is a larger change (the raw join SQL and the
// single-`<Field>ID` selector would have to be generalized to every primary /
// foreign field). Note that the presets list editor itself
// (presets.NestedSlice + gorm2op.SaveHasManyAssociation) already handles
// composite primary keys for has-many reconciliation — these helper builders are
// a separate, selector-oriented layer.
package helper
