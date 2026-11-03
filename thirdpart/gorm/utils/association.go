package utils

import (
	"fmt"
	_ "unsafe"

	. "gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// AssociationDB is the *gorm.DB of the records of assoc: its condition —
// gorm's (*Association).buildCondition, linked by name.
//
//go:linkname AssociationDB gorm.io/gorm.(*Association).buildCondition
func AssociationDB(assoc *Association) *DB

// Instance is the statement instance of db — gorm's (*DB).getInstance,
// linked by name.
//
//go:linkname Instance gorm.io/gorm.(*DB).getInstance
func Instance(*DB) *DB

// AssociationManager holds the SQL that filters, links and unlinks the
// records of a many-to-many association through its join table. Its queries
// assume integer ids (BIGINT casts, a related table's "id" column): to link
// and unlink records of any key, prefer gorm's Association API (README).
type AssociationManager struct {
	assoc *Association
	relatedTable,
	filterQuery, insertQuery string
	deleteQuery, linkQuery string
}

// NewAssociationManager is the manager of assoc: its queries made for a
// many-to-many relationship (none for the others).
func NewAssociationManager(assoc *Association) *AssociationManager {
	a := &AssociationManager{assoc: assoc}
	relation := assoc.Relationship
	a.relatedTable = relation.FieldSchema.Table

	switch relation.Type {
	case schema.Many2Many:
		joinTable := relation.JoinTable
		a.filterQuery = fmt.Sprintf("EXISTS (SELECT 1 FROM %s rel WHERE rel.%s = ? AND rel.%s = %s.id)",
			joinTable.Name, joinTable.DBNames[0], joinTable.DBNames[1], a.relatedTable)
		a.insertQuery = fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES (?, ?)",
			joinTable.Table, joinTable.DBNames[0], joinTable.DBNames[1])
		a.deleteQuery = fmt.Sprintf("DELETE FROM %s WHERE %s = ? AND %s = ?",
			joinTable.Table, joinTable.DBNames[0], joinTable.DBNames[1])
		a.linkQuery = fmt.Sprintf(m2mInsertQuery, joinTable.Table, joinTable.DBNames[0], joinTable.DBNames[1])
	}

	return a
}

const m2mInsertQuery = `with 
data as (
	select ?::BIGINT as f_id, unnest('{$ID$}'::BIGINT[]) as p_id
) 
, data_ok as (
	select f_id::BIGINT, p_ID::BIGINT from data d where not exists (
	select 1 
	from %s fp 
	where fp.%s = f_id and fp.%s = p_id)
)
insert into %[1]s (%[2]s, %[3]s) select f_id::BIGINT, p_id::BIGINT from data_ok;
`
