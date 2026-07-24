package helper

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// LinkTag has a STRING primary key (standing in for a UUID or any non-integer
// key), to prove the m2m link query is not tied to an integer `id`.
type LinkTag struct {
	Code string `gorm:"primaryKey"`
	Name string
}

type LinkPost struct {
	ID   uint
	Name string
	Tags []*LinkTag `gorm:"many2many:link_post_tags;"`
}

// buildLinkQuery mirrors Build()'s many-to-many link query construction so the
// test exercises the exact template + column resolution the builder uses.
func buildLinkQuery(db *gorm.DB) (query, relatedPKCol string, ownerCol, relatedCol string) {
	rel := db.Model(&LinkPost{}).Association("Tags").Relationship
	jt := rel.JoinTable
	relatedPKCol = rel.FieldSchema.PrimaryFields[0].DBName
	query = fmt.Sprintf(m2mInsertQuery, rel.FieldSchema.Table, jt.Table, jt.DBNames[0], jt.DBNames[1], relatedPKCol)
	return query, relatedPKCol, jt.DBNames[0], jt.DBNames[1]
}

func TestM2MLinkQuery_NonIntegerKey(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&LinkPost{}, &LinkTag{}); err != nil {
		t.Fatal(err)
	}

	query, relatedPKCol, ownerCol, relatedCol := buildLinkQuery(db)

	// the query must not cast the key to a fixed SQL type and must key off the
	// schema's real primary-key column (here "code"), not a hardcoded "id".
	if strings.Contains(strings.ToUpper(query), "BIGINT") {
		t.Fatalf("link query must not cast to BIGINT:\n%s", query)
	}
	if relatedPKCol != "code" {
		t.Fatalf("related PK column = %q, want \"code\"", relatedPKCol)
	}

	// seed tags with string keys and a post
	if err = db.Create(&[]*LinkTag{{Code: "a", Name: "A"}, {Code: "b", Name: "B"}, {Code: "c", Name: "C"}}).Error; err != nil {
		t.Fatal(err)
	}
	if err = db.Create(&LinkPost{ID: 1, Name: "p"}).Error; err != nil {
		t.Fatal(err)
	}

	// link a + b
	if err = db.Exec(query, 1, []any{"a", "b"}).Error; err != nil {
		t.Fatalf("link a,b: %v", err)
	}
	// link b + c again — b must be de-duplicated, c added
	if err = db.Exec(query, 1, []any{"b", "c"}).Error; err != nil {
		t.Fatalf("link b,c: %v", err)
	}

	var codes []string
	if err = db.Table("link_post_tags").
		Where(ownerCol+" = ?", 1).
		Order(relatedCol).
		Pluck(relatedCol, &codes).Error; err != nil {
		t.Fatal(err)
	}
	sort.Strings(codes)
	if fmt.Sprint(codes) != fmt.Sprint([]string{"a", "b", "c"}) {
		t.Fatalf("linked codes = %v, want [a b c] (deduplicated)", codes)
	}
}

// CKItem is a many-to-many related model with a COMPOSITE primary key
// (OrgID, Code), of mixed types — proving the join filter and link/unlink work
// with more than one related key field.
type CKItem struct {
	OrgID uint   `gorm:"primaryKey"`
	Code  string `gorm:"primaryKey"`
	Name  string
}

type CKParent struct {
	ID    uint
	Items []*CKItem `gorm:"many2many:ck_parent_items;"`
}

func TestM2MParentFilter_CompositeRelatedKey(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&CKParent{}, &CKItem{}); err != nil {
		t.Fatal(err)
	}

	rel := db.Model(&CKParent{}).Association("Items").Relationship

	// the filter builder must produce one owner placeholder (parent id) and
	// correlate BOTH related key columns (org_id, code) — no fixed SQL cast.
	filterSQL, ownerFields := m2mParentFilter(rel)
	if len(ownerFields) != 1 || ownerFields[0] != "ID" {
		t.Fatalf("owner filter fields = %v, want [ID]", ownerFields)
	}
	if strings.Contains(strings.ToUpper(filterSQL), "BIGINT") {
		t.Fatalf("filter must not cast to BIGINT:\n%s", filterSQL)
	}
	for _, col := range []string{"org_id", "code"} {
		if !strings.Contains(filterSQL, "ck_item_"+col) && !strings.Contains(filterSQL, col) {
			t.Fatalf("filter must correlate related key column %q:\n%s", col, filterSQL)
		}
	}

	// seed items + parent, link two of three via the Association API (exactly what
	// associationAppend does)
	items := []*CKItem{
		{OrgID: 1, Code: "a", Name: "A"},
		{OrgID: 1, Code: "b", Name: "B"},
		{OrgID: 2, Code: "a", Name: "A2"},
	}
	if err = db.Create(&items).Error; err != nil {
		t.Fatal(err)
	}
	parent := &CKParent{ID: 1}
	if err = db.Create(parent).Error; err != nil {
		t.Fatal(err)
	}
	if err = db.Model(parent).Association("Items").Append([]*CKItem{items[0], items[2]}); err != nil {
		t.Fatalf("append: %v", err)
	}

	linkedNames := func() []string {
		var got []CKItem
		if err := db.Table(rel.FieldSchema.Table).
			Where(filterSQL, 1).Order("name").Find(&got).Error; err != nil {
			t.Fatal(err)
		}
		names := make([]string, len(got))
		for i, g := range got {
			names[i] = g.Name
		}
		return names
	}

	if got := linkedNames(); fmt.Sprint(got) != fmt.Sprint([]string{"A", "A2"}) {
		t.Fatalf("linked = %v, want [A A2]", got)
	}

	// unlink (1,a) via the Association API (exactly what associationDeleter does)
	if err = db.Model(parent).Association("Items").Delete(&CKItem{OrgID: 1, Code: "a"}); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if got := linkedNames(); fmt.Sprint(got) != fmt.Sprint([]string{"A2"}) {
		t.Fatalf("after unlink linked = %v, want [A2]", got)
	}

	// the related rows themselves are untouched (only the link was removed)
	var count int64
	db.Model(&CKItem{}).Count(&count)
	if count != 3 {
		t.Fatalf("expected 3 items to remain, got %d", count)
	}
}
