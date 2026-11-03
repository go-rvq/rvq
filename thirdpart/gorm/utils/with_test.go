package utils

import (
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestWithClausesRecursive(t *testing.T) {
	plain := WithClauses{{Name: "a", Query: "SELECT 1"}}
	if q, _ := plain.BuildQuery(); !strings.HasPrefix(q, "WITH a") {
		t.Fatalf("plain: %q", q)
	}
	rec := WithClauses{
		{Name: "a", Query: "SELECT ?", Args: []any{1}},
		{Name: "t", Fields: []string{"id"}, Recursive: true,
			Query: "SELECT id FROM x WHERE id = ? UNION SELECT x.id FROM x JOIN t ON x.parent_id = t.id", Args: []any{2}},
	}
	q, args := rec.BuildQuery()
	if !strings.HasPrefix(q, "WITH RECURSIVE a") || !strings.Contains(q, "t(id) AS (") {
		t.Fatalf("recursive: %q", q)
	}
	if len(args) != 2 || args[0] != 1 || args[1] != 2 {
		t.Fatalf("args %v", args)
	}
}

func TestWithDbClauseUpdateDelete(t *testing.T) {
	with := WithClauses{{Name: "old", Query: "SELECT id FROM nodes WHERE value < ?", Args: []any{10}}}
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	type node struct {
		ID    uint
		Value int
	}
	explain := func(db *gorm.DB) string {
		return db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
	}
	up := explain(db.Model(&node{}).Clauses(with.DbClause("UPDATE")).
		Where("id IN (SELECT id FROM old)").Update("value", 0))
	if !strings.HasPrefix(up, "WITH old AS (") || !strings.Contains(up, ") UPDATE `nodes` SET `value`=0 WHERE id IN (SELECT id FROM old)") {
		t.Fatalf("update: %s", up)
	}
	del := explain(db.Clauses(with.DbClause("DELETE")).Where("id IN (SELECT id FROM old)").Delete(&node{}))
	if !strings.HasPrefix(del, "WITH old AS (") || !strings.Contains(del, ") DELETE FROM `nodes` WHERE id IN (SELECT id FROM old)") {
		t.Fatalf("delete: %s", del)
	}
	if !strings.Contains(up, "value < 10") || !strings.Contains(del, "value < 10") {
		t.Fatalf("args not bound: %s / %s", up, del)
	}
}
