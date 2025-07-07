package utils

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RawColumn struct {
	Table, Name, Query string
}

func (r *RawColumn) ModifyStatement(stmt *gorm.Statement) {
	clause := stmt.Clauses["SELECT"]
	clause.AfterNameExpression = r
	stmt.Clauses["SELECT"] = clause
}

func (r *RawColumn) Build(builder clause.Builder) {
	stmt, _ := builder.(*gorm.Statement)
	if stmt == nil {
		return
	}
	selClause := stmt.Clauses["SELECT"]
	expr := stmt.Clauses["SELECT"].Expression.(clause.Select)
	if len(expr.Columns) == 0 {
		builder.WriteString(r.Query + " AS " + r.Name)
	} else {
		var found bool
		for i, col := range expr.Columns {
			if col.Table == r.Table && col.Name == r.Name {
				col.Raw = true
				col.Table = ""
				col.Name = r.Query + " as " + r.Name
				stmt.Clauses["SELECT"].Expression.(clause.Select).Columns[i] = col
				found = true
				break
			}
		}
		if !found {
			builder.WriteString(r.Query + " AS " + r.Name + ", ")
		}
	}

	selClause.Expression = expr
	stmt.Clauses["SELECT"] = selClause
}

func SetRawColumn(table, name, query string) *RawColumn {
	return &RawColumn{Table: table, Name: name, Query: query}
}
