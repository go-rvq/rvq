# WITH clauses (common table expressions)

`with.go` puts a `WITH [RECURSIVE] … AS (…)` before a gorm query. gorm has a
`clause.With`, but it is an empty type (`type With struct{}`): it builds
nothing, and gorm's query callback has no `WITH` clause. These helpers fill
that gap.

## Types

| Type | What |
|------|------|
| `WithClause` | One table of the WITH: `Name[(Fields)] AS (Query)`, the `?` of `Query` bound to `Args`; `Recursive` when `Query` refers to `Name`. |
| `WithClauses` | The tables of one WITH, in order (a table may refer to those before it). `BuildQuery` gives the SQL and its values; `Append` / `AppendClause` add tables. |
| `WithDbClause` | A `WithClauses` as a gorm clause (`WithClauses.DbClause()`), for `db.Clauses(…)`. |

## How it works

`WithDbClause` is a `gorm.StatementModifier`. It does not add a clause of its
own. It becomes the `BeforeExpression` of the `SELECT` clause, so gorm writes
the WITH first and then the query: `WITH … SELECT … FROM …`.

- **The main query names the tables:** `db.Table("tree")`,
  `db.Joins("JOIN children c ON …")`, or a subquery.
- **Other statements:** `DbClause("UPDATE")` and `DbClause("DELETE")` put
  the WITH before that clause instead of the SELECT (`with_test.go`).
- **One WITH per query:** two `DbClause` on the same query are joined into a
  single WITH, the tables of the second after those of the first.
- **`RECURSIVE` for the whole WITH:** if any table is `Recursive`, the WITH is
  `WITH RECURSIVE`. That is the syntax of PostgreSQL, SQLite and MySQL 8.
- **Names are not quoted:** `Name` and `Fields` are written as they are.

## Examples

These are the examples of `example_test.go`; `go test` checks their SQL.

### A tree (recursive)

The subtree of a node: the node itself, then the children of each node at
every level.

```go
with := gormutils.WithClauses{{
	Name:      "tree",
	Recursive: true,
	Query: `SELECT n.* FROM nodes n WHERE n.id = ?
UNION
SELECT c.* FROM nodes c JOIN tree t ON c.parent_id = t.id`,
	Args: []any{rootID},
}}
var nodes []*Node
err := db.Clauses(with.DbClause()).Table("tree").Order("id").Find(&nodes).Error
```

```sql
WITH RECURSIVE tree AS (
	SELECT n.* FROM nodes n WHERE n.id = 1
	UNION
	SELECT c.* FROM nodes c JOIN tree t ON c.parent_id = t.id
) SELECT * FROM `tree` ORDER BY id
```

`UNION` (not `UNION ALL`) drops repeated rows, so a cycle in the data ends
the recursion. A model with a `gorm.DeletedAt` also gets
`tree.deleted_at IS NULL` from gorm. Filter the deleted rows inside `Query`
too, so the recursion does not descend from them.

### Several tables, joined

```go
var with gormutils.WithClauses
with.AppendClause("big", "SELECT id FROM nodes WHERE value > ?", []string{"id"}, 100)
with.AppendClause("children", "SELECT n.id FROM nodes n JOIN big b ON n.parent_id = b.id", nil)
err := db.Clauses(with.DbClause()).
	Joins("JOIN children c ON c.id = nodes.id").
	Find(&nodes).Error
```

```sql
WITH big(id) AS (
	SELECT id FROM nodes WHERE value > 100
)
, children AS (
	SELECT n.id FROM nodes n JOIN big b ON n.parent_id = b.id
) SELECT `nodes`.`id`,`nodes`.`parent_id`,`nodes`.`value` FROM `nodes` JOIN children c ON c.id = nodes.id
```

### In a data operator callback

The table restricts the records of a listing, and the presets query stays
the same:

```go
cb.Pre(func(state *gorm2op.CallbackState) error {
	with := gormutils.WithClauses{{
		Name:   "allowed",
		Fields: []string{"id"},
		Query:  "SELECT n.id FROM nodes n WHERE n.owner_id = ?",
		Args:   []any{ownerID},
	}}
	state.DB = state.DB.Clauses(with.DbClause()).
		Joins("INNER JOIN allowed a ON a.id = nodes.id")
	return nil
})
```

### Two clauses on one query

```go
a := gormutils.WithClauses{{Name: "a", Query: "SELECT 1 AS id"}}
b := gormutils.WithClauses{{Name: "b", Query: "SELECT id FROM a"}}
db.Clauses(a.DbClause(), b.DbClause()).Table("b").Find(&nodes)
```

```sql
WITH a AS (
	SELECT 1 AS id
)
, b AS (
	SELECT id FROM a
) SELECT * FROM `b`
```
