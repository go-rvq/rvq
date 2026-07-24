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
// These helpers do not assume an integer `id`: primary/foreign key columns are
// read from the schema and used as their own type, so integer, UUID and other
// key types all work — and COMPOSITE keys (two or more columns) are supported:
//
//   - NestedSlice reads related rows through a filter built from the relationship
//     references, and links/unlinks them through gorm's Association API, which
//     handles any key type and composite keys (see hasManyParentFilter /
//     m2mParentFilter / associationAppend / associationDeleter, and
//     TestM2MParentFilter_CompositeRelatedKey).
//   - ModelSelectorBuilder resolves every foreign-key column of the belongs-to
//     relationship (foreignKeyFieldsOf) and writes them all when a record is
//     selected, so a composite foreign key is persisted in full (see
//     TestForeignKeyFieldsOf_Composite).
//
// See README.md and docs/ in this package for a fuller guide.
package helper
