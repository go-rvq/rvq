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
	relatedPKCol = relatedPKColumn(rel)
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
