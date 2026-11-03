package utils

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// RawColumn is a column of the SELECT made of an SQL expression: Query AS
// Name. db.Clauses(SetRawColumn(table, name, query)) replaces the column
// table.name the query selects by it, or, not selected (or a SELECT of every
// column), adds it before the others. Query is written as it is (Args are
// not bound).
type RawColumn struct {
	Table, Name, Query string
	Args               []any
}

// ModifyStatement puts r in the SELECT of stmt (its AfterNameExpression).
func (r *RawColumn) ModifyStatement(stmt *gorm.Statement) {
	clause := stmt.Clauses["SELECT"]
	clause.AfterNameExpression = r
	stmt.Clauses["SELECT"] = clause
}

// Build replaces the column Table.Name of the SELECT — selected by the model
// (QueryFields, a Joins), or by Select("name") or Select("table.name") —
// by Query AS Name; not selected, writes Query AS Name before the columns
// (before the *, when none).
func (r *RawColumn) Build(builder clause.Builder) {
	stmt, _ := builder.(*gorm.Statement)
	if stmt == nil {
		return
	}
	if sel, ok := stmt.Clauses["SELECT"].Expression.(clause.Select); ok {
		for i, col := range sel.Columns {
			if r.is(col) {
				// the slice of the clause: the SELECT built next writes it
				sel.Columns[i] = clause.Column{Name: r.Query + " AS " + r.Name, Raw: true}
				return
			}
		}
	}
	builder.WriteString(r.Query + " AS " + r.Name + ",")
}

// is says col is the column of r: Table.Name, Name, or the raw "Table.Name".
func (r *RawColumn) is(col clause.Column) bool {
	if col.Raw {
		return col.Name == r.Name || col.Name == r.Table+"."+r.Name
	}
	return col.Name == r.Name && (col.Table == "" || col.Table == r.Table || col.Table == clause.CurrentTable)
}

// SetRawColumn is the RawColumn table.name = query.
func SetRawColumn(table, name, query string) *RawColumn {
	return &RawColumn{Table: table, Name: name, Query: query}
}
