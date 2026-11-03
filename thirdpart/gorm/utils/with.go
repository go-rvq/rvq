package utils

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/hints"
)

// WithClause is a common table expression of a WITH: Name[(Fields)] AS
// (Query), the placeholders (?) of Query bound to Args. Name and Fields are
// written as they are (not quoted).
type WithClause struct {
	// Name is the name of the table the query defines, the one the main query
	// (and, recursive, Query itself) refers to.
	Name string
	// Fields name the columns of the table; none: those of Query.
	Fields []string
	// Query is the SELECT of the table; a recursive one: the initial SELECT,
	// UNION [ALL], the SELECT that joins Name.
	Query string
	// Args are the values of the placeholders of Query, in order.
	Args []any
	// DB is the query as a query of gorm (its values its own), in place of
	// Query and Args: a query built by gorm — the filtered query of a
	// listing, say — goes into the WITH as it is.
	DB *gorm.DB
	// Recursive says Query refers to Name (a tree, say): any of the clauses
	// recursive makes the WITH a WITH RECURSIVE (PostgreSQL, SQLite, MySQL 8).
	Recursive bool
}

// WriteTo writes c to b: Name[(Fields)] AS (Query), Query indented; with
// DB, a placeholder (?) in place of the query, DB its value (BuildQuery).
func (c *WithClause) WriteTo(b *strings.Builder) {
	b.WriteString(c.Name)
	if len(c.Fields) > 0 {
		b.WriteString("(")
		b.WriteString(strings.Join(c.Fields, ","))
		b.WriteString(")")
	}
	b.WriteString(" AS (\n")
	if c.DB != nil {
		b.WriteString("\t?\n)")
		return
	}
	for _, line := range strings.Split(strings.TrimSpace(c.Query), "\n") {
		b.WriteString("\t")
		b.WriteString(line)
		b.WriteString("\n")
	}
	b.WriteString(")")
}

// values are the values of the placeholders of c: Args, or DB alone.
func (c *WithClause) values() []any {
	if c.DB != nil {
		return []any{c.DB}
	}
	return c.Args
}

// String is c as WriteTo writes it.
func (c *WithClause) String() string {
	var b strings.Builder
	c.WriteTo(&b)
	return b.String()
}

// WithClauses are the common table expressions of one WITH, in order: a
// clause may refer to those before it.
type WithClauses []*WithClause

// BuildQuery is the WITH of c — WITH [RECURSIVE] a AS (…), b AS (…) — and the
// values of its placeholders, those of each clause in order (the DB of a
// clause, its value: gorm writes it as its query). None: "" and no values.
func (c WithClauses) BuildQuery() (q string, args []any) {
	if len(c) == 0 {
		return "", nil
	}
	var (
		b strings.Builder
	)

	b.WriteString("WITH ")
	for _, wc := range c {
		if wc.Recursive {
			b.WriteString("RECURSIVE ")
			break
		}
	}

	for _, wc := range c[:len(c)-1] {
		wc.WriteTo(&b)
		b.WriteString("\n, ")
		args = append(args, wc.values()...)
	}

	last := c[len(c)-1]
	last.WriteTo(&b)
	args = append(args, last.values()...)
	return b.String(), args
}

// Append adds the clauses wc to c.
func (c *WithClauses) Append(wc ...*WithClause) {
	*c = append(*c, wc...)
}

// AppendClause adds the clause name[(fields)] AS (query), its placeholders
// bound to args.
func (c *WithClauses) AppendClause(name string, query string, fields []string, args ...any) {
	c.Append(&WithClause{Name: name, Fields: fields, Query: query, Args: args})
}

// DbClause is c as a gorm clause: db.Clauses(c.DbClause()) puts the WITH of c
// before the SELECT of the query (clause names the gorm clauses to put it
// before instead: "UPDATE", "DELETE"). The main query refers to the
// tables by their names: db.Table(name), db.Joins("JOIN name …"), a
// subquery.
func (c WithClauses) DbClause(clause ...string) *WithDbClause {
	if len(clause) == 0 {
		clause = append(clause, "SELECT")
	}
	return &WithDbClause{clauses: clause, items: c}
}

// WithDbClause is a WITH as a gorm clause (a gorm.StatementModifier): its
// clauses before the clauses of the statement it names. See
// WithClauses.DbClause.
type WithDbClause struct {
	items   WithClauses
	clauses []string
}

// ModifyStatement puts c before each clause of c of stmt (its
// BeforeExpression). A WITH already there is joined into one — a query has a
// single WITH —, the clauses of c after its own; any other expression there,
// kept before c.
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

// Build writes the WITH of c and binds its values. None: nothing.
func (c WithDbClause) Build(builder clause.Builder) {
	q, args := c.items.BuildQuery()
	if q == "" {
		return
	}
	gorm.Expr(q, args...).Build(builder)
}
