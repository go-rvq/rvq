package utils

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/hints"
)

type WithClause struct {
	Name   string
	Fields []string
	Query  string
	Args   []any
}

func (c *WithClause) WriteTo(b *strings.Builder) {
	b.WriteString(c.Name)
	if len(c.Fields) > 0 {
		b.WriteString("(")
		b.WriteString(strings.Join(c.Fields, ","))
		b.WriteString(")")
	}
	b.WriteString(" AS (\n\t")
	lines := strings.Split(c.Query, "\n")
	for _, line := range lines {
		b.WriteString("\t")
		b.WriteString(line)
		b.WriteString("\n")
	}
	b.WriteString("\n)")
}

func (c *WithClause) String() string {
	var b strings.Builder
	c.WriteTo(&b)
	return b.String()
}

type WithClauses []*WithClause

func (c WithClauses) BuildQuery() (q string, args []any) {
	var (
		b strings.Builder
	)

	b.WriteString("WITH ")

	for _, wc := range c[:len(c)-1] {
		wc.WriteTo(&b)
		b.WriteString("\n, ")
		args = append(args, wc.Args...)
	}

	last := c[len(c)-1]
	last.WriteTo(&b)
	args = append(args, last.Args...)
	return b.String(), args
}

func (c *WithClauses) Append(wc ...*WithClause) {
	*c = append(*c, wc...)
}

func (c *WithClauses) AppendClause(name string, query string, fields []string, args ...any) {
	c.Append(&WithClause{Name: name, Fields: fields, Query: query, Args: args})
}

func (c WithClauses) DbClause(clause ...string) *WithDbClause {
	if len(clause) == 0 {
		clause = append(clause, "SELECT")
	}
	return &WithDbClause{clauses: clause, items: c}
}

type WithDbClause struct {
	items   WithClauses
	clauses []string
}

func (c WithDbClause) ModifyStatement(stmt *gorm.Statement) {
	for _, name := range c.clauses {
		name = strings.ToUpper(name)
		clause := stmt.Clauses[name]
		if clause.BeforeExpression == nil {
			clause.BeforeExpression = c
		} else if old, ok := clause.BeforeExpression.(WithDbClause); ok {
			clause.BeforeExpression = WithDbClause{items: append(old.items, c.items...)}
		} else {
			clause.BeforeExpression = hints.Exprs{clause.BeforeExpression, c}
		}

		stmt.Clauses[name] = clause
	}
}

func (c WithDbClause) Build(builder clause.Builder) {
	q, args := c.items.BuildQuery()
	gorm.Expr(q, args...).Build(builder)
}
