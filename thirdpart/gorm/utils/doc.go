// Package utils holds small helpers on top of gorm, used by the presets data
// operator (admin/presets/gorm2op), admin/packages/helper and applications:
//
//   - keys and relationships (relationship.go): the filters of a child by its
//     parent and the primary keys, of any type and composite (docs/keys.md);
//   - common table expressions (with.go): WithClause, WithClauses and
//     WithDbClause put a WITH [RECURSIVE] before a query — the clause.With
//     gorm lacks (an empty type there) — (docs/with.md);
//   - a raw SELECT column (select_column.go): RawColumn, SetRawColumn;
//   - gorm internals (association.go): AssociationDB, Instance (by
//     go:linkname) and AssociationManager.
package utils
