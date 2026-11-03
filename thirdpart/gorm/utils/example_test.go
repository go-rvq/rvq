package utils_test

import (
	"fmt"

	gormutils "github.com/go-rvq/rvq/thirdpart/gorm/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Node is a record of a tree: its parent, a node of the same table.
type Node struct {
	ID       uint
	ParentID *uint
	Value    int
}

// dryRun is a database that builds the SQL of a query without running it.
func dryRun() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{DryRun: true})
	if err != nil {
		panic(err)
	}
	return db
}

// sql is the SQL of the query db ran, its values in place.
func sql(db *gorm.DB) string {
	return db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
}

// A recursive WITH: the subtree of the node 1 — it and, every level, the
// children of each.
func ExampleWithClauses_recursive() {
	with := gormutils.WithClauses{{
		Name:      "tree",
		Recursive: true,
		Query: `SELECT n.* FROM nodes n WHERE n.id = ?
UNION
SELECT c.* FROM nodes c JOIN tree t ON c.parent_id = t.id`,
		Args: []any{1},
	}}
	var nodes []*Node
	fmt.Println(sql(dryRun().Clauses(with.DbClause()).Table("tree").Order("id").Find(&nodes)))
	// Output:
	// WITH RECURSIVE tree AS (
	// 	SELECT n.* FROM nodes n WHERE n.id = 1
	// 	UNION
	// 	SELECT c.* FROM nodes c JOIN tree t ON c.parent_id = t.id
	// ) SELECT * FROM `tree` ORDER BY id
}

// Several clauses, one referring to the one before, joined by the main query.
func ExampleWithClauses_join() {
	var with gormutils.WithClauses
	with.AppendClause("big", "SELECT id FROM nodes WHERE value > ?", []string{"id"}, 100)
	with.AppendClause("children", "SELECT n.id FROM nodes n JOIN big b ON n.parent_id = b.id", nil)
	var nodes []*Node
	fmt.Println(sql(dryRun().Clauses(with.DbClause()).
		Joins("JOIN children c ON c.id = nodes.id").Find(&nodes)))
	// Output:
	// WITH big(id) AS (
	// 	SELECT id FROM nodes WHERE value > 100
	// )
	// , children AS (
	// 	SELECT n.id FROM nodes n JOIN big b ON n.parent_id = b.id
	// ) SELECT `nodes`.`id`,`nodes`.`parent_id`,`nodes`.`value` FROM `nodes` JOIN children c ON c.id = nodes.id
}

// Two DbClause on one query make a single WITH.
func ExampleWithDbClause_ModifyStatement() {
	a := gormutils.WithClauses{{Name: "a", Query: "SELECT 1 AS id"}}
	b := gormutils.WithClauses{{Name: "b", Query: "SELECT id FROM a"}}
	var nodes []*Node
	fmt.Println(sql(dryRun().Clauses(a.DbClause(), b.DbClause()).Table("b").Find(&nodes)))
	// Output:
	// WITH a AS (
	// 	SELECT 1 AS id
	// )
	// , b AS (
	// 	SELECT id FROM a
	// ) SELECT * FROM `b`
}

// NodeDouble is a Node with a column computed by the query.
type NodeDouble struct {
	Node
	Double int
}

// A column of the SELECT made of an expression: it takes the place of the
// column of the model, selected by Select or by the model (QueryFields, a
// Joins); none selected, it goes before the *.
func ExampleSetRawColumn() {
	double := func() *gormutils.RawColumn { return gormutils.SetRawColumn("nodes", "double", "value * 2") }
	var nodes []*NodeDouble
	fmt.Println(sql(dryRun().Clauses(double()).Table("nodes").Select("id", "double").Find(&nodes)))
	fmt.Println(sql(dryRun().Clauses(double()).Table("nodes").Joins("JOIN x ON x.id = nodes.id").Find(&nodes)))
	fmt.Println(sql(dryRun().Clauses(double()).Table("nodes").Find(&nodes)))
	// Output:
	// SELECT  `id`,value * 2 AS double FROM `nodes`
	// SELECT  `nodes`.`id`,`nodes`.`parent_id`,`nodes`.`value`,value * 2 AS double FROM `nodes` JOIN x ON x.id = nodes.id
	// SELECT value * 2 AS double, * FROM `nodes`
}
