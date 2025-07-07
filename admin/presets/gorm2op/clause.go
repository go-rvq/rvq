package gorm2op

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/hints"
)

type SetTableNameClauseBuilder struct {
	TableName string
}

func (h SetTableNameClauseBuilder) ModifyStatement(stmt *gorm.Statement) {
	for name, Clause := range stmt.Clauses {
		if Clause.AfterExpression == nil {
			Clause.AfterExpression = h
		} else if _, ok := Clause.AfterExpression.(SetTableNameClauseBuilder); ok {
		} else {
			Clause.AfterExpression = hints.Exprs{Clause.AfterExpression, h}
		}

		stmt.Clauses[name] = Clause
	}
}

func (h SetTableNameClauseBuilder) Build(builder clause.Builder) {
	builder.WriteString("---")
}
