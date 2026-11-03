package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type capNode struct {
	ID       uint
	ParentID *uint
	Value    int
	Children []*capNode `gorm:"foreignKey:ParentID"`
	found    bool
}

func (n *capNode) AfterFind(*gorm.DB) error { n.found = true; return nil }

// The page, the count and the totals of the whole filter in one query: the
// count by a window function, the totals by a subquery of the WITH in a JSON
// column; the model keeps its hook and its preload.
func TestCaptureColumns(t *testing.T) {
	t.Run("sqlite", func(t *testing.T) {
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}
		testCaptureColumns(t, db, "json_object")
	})
	// RVQ_TEST_POSTGRES: a DSN of a database to test against (a schema of its
	// own, dropped after)
	t.Run("postgres", func(t *testing.T) {
		dsn := os.Getenv("RVQ_TEST_POSTGRES")
		if dsn == "" {
			t.Skip("RVQ_TEST_POSTGRES not set")
		}
		db, err := gorm.Open(postgres.Open(dsn+" search_path=rvq_capture_test"), &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}
		if err := db.Exec("DROP SCHEMA IF EXISTS rvq_capture_test CASCADE; CREATE SCHEMA rvq_capture_test").Error; err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { db.Exec("DROP SCHEMA IF EXISTS rvq_capture_test CASCADE") })
		testCaptureColumns(t, db, "json_build_object")
	})
}

func testCaptureColumns(t *testing.T, db *gorm.DB, jsonObject string) {
	var err error
	if err := RegisterColumnCapture(db); err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&capNode{}); err != nil {
		t.Fatal(err)
	}
	// 6 roots (-30, -20, -10, 10, 20, 30), a child of the first
	for i, v := range []int{-30, -20, -10, 10, 20, 30} {
		db.Create(&capNode{ID: uint(i + 1), Value: v})
	}
	one := uint(1)
	db.Create(&capNode{ID: 7, ParentID: &one, Value: 5})

	with := WithClauses{{Name: "base", Query: "SELECT id, value FROM cap_nodes WHERE parent_id IS NULL"}}
	q, capt := CaptureColumns(db.Model(&capNode{}).Clauses(with.DbClause()))
	capt.Names = []string{"__count", "__aggs"}
	var page []*capNode
	err = q.Select("cap_nodes.*",
		"COUNT(*) OVER () AS __count",
		`(SELECT `+jsonObject+`('total', SUM(value),
			'negative', SUM(CASE WHEN value < 0 THEN value END),
			'positive', SUM(CASE WHEN value > 0 THEN value END)) FROM base) AS __aggs`).
		Where("cap_nodes.id IN (SELECT id FROM base)").
		Preload("Children").
		Order("cap_nodes.id").Limit(2).
		Find(&page).Error
	if err != nil {
		t.Fatal(err)
	}

	if len(page) != 2 || page[0].ID != 1 || page[1].ID != 2 {
		t.Fatalf("the page: %v", page)
	}
	if !page[0].found || !page[1].found {
		t.Fatal("AfterFind did not run")
	}
	if len(page[0].Children) != 1 || page[0].Children[0].ID != 7 {
		t.Fatalf("the preload: %v", page[0].Children)
	}
	if fmt.Sprint(capt.Values["__count"]) != "6" {
		t.Fatalf("__count %#v, want 6 (the whole filter, not the page)", capt.Values["__count"])
	}
	var aggs map[string]int
	if err := json.Unmarshal([]byte(fmt.Sprint(capt.Values["__aggs"])), &aggs); err != nil {
		t.Fatalf("__aggs %#v: %v", capt.Values["__aggs"], err)
	}
	if aggs["total"] != 0 || aggs["negative"] != -60 || aggs["positive"] != 60 {
		t.Fatalf("aggs %v", aggs)
	}
	t.Logf("count %v, aggs %v", capt.Values["__count"], aggs)

	// a query that asks for no capture: gorm's own
	var all []*capNode
	if err := db.Find(&all).Error; err != nil || len(all) != 7 {
		t.Fatalf("plain query: %d, %v", len(all), err)
	}
}
