package utils

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// CKItem is a many-to-many related model with a COMPOSITE primary key
// (OrgID, Code) of mixed types.
type CKItem struct {
	OrgID uint   `gorm:"primaryKey"`
	Code  string `gorm:"primaryKey"`
	Name  string
}

type CKParent struct {
	ID    uint
	Items []*CKItem `gorm:"many2many:ck_parent_items;"`
}

// HMParent/HMChild model a has-many with a COMPOSITE foreign key.
type HMParent struct {
	OrgID    uint      `gorm:"primaryKey"`
	Code     string    `gorm:"primaryKey"`
	Children []HMChild `gorm:"foreignKey:ParentOrgID,ParentCode;references:OrgID,Code"`
}

type HMChild struct {
	ID          uint
	ParentOrgID uint
	ParentCode  string
	Name        string
}

func newDB(t *testing.T, models ...any) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(models...); err != nil {
		t.Fatal(err)
	}
	return db
}

// TestM2MParentFilter_CompositeRelatedKey proves the m2m parent filter correlates
// every related key column (composite) with no fixed SQL cast, and that link /
// unlink via gorm's Association API keep the related rows intact.
func TestM2MParentFilter_CompositeRelatedKey(t *testing.T) {
	db := newDB(t, &CKParent{}, &CKItem{})
	rel := db.Model(&CKParent{}).Association("Items").Relationship

	filterSQL, ownerFields := M2MParentFilter(rel)
	if len(ownerFields) != 1 || ownerFields[0] != "ID" {
		t.Fatalf("owner filter fields = %v, want [ID]", ownerFields)
	}
	if strings.Contains(strings.ToUpper(filterSQL), "BIGINT") {
		t.Fatalf("filter must not cast to BIGINT:\n%s", filterSQL)
	}
	for _, col := range []string{"org_id", "code"} {
		if !strings.Contains(filterSQL, col) {
			t.Fatalf("filter must correlate related key column %q:\n%s", col, filterSQL)
		}
	}

	items := []*CKItem{
		{OrgID: 1, Code: "a", Name: "A"},
		{OrgID: 1, Code: "b", Name: "B"},
		{OrgID: 2, Code: "a", Name: "A2"},
	}
	if err := db.Create(&items).Error; err != nil {
		t.Fatal(err)
	}
	parent := &CKParent{ID: 1}
	if err := db.Create(parent).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Model(parent).Association("Items").Append([]*CKItem{items[0], items[2]}); err != nil {
		t.Fatalf("append: %v", err)
	}

	linkedNames := func() []string {
		var got []CKItem
		if err := db.Table(rel.FieldSchema.Table).Where(filterSQL, 1).Order("name").Find(&got).Error; err != nil {
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
	if err := db.Model(parent).Association("Items").Delete(&CKItem{OrgID: 1, Code: "a"}); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if got := linkedNames(); fmt.Sprint(got) != fmt.Sprint([]string{"A2"}) {
		t.Fatalf("after unlink linked = %v, want [A2]", got)
	}
	var count int64
	db.Model(&CKItem{}).Count(&count)
	if count != 3 {
		t.Fatalf("expected 3 items to remain, got %d", count)
	}
}

// TestHasManyParentFilter_Composite covers a has-many child filtered by a
// composite foreign key.
func TestHasManyParentFilter_Composite(t *testing.T) {
	db := newDB(t, &HMParent{}, &HMChild{})
	rel := db.Model(&HMParent{}).Association("Children").Relationship

	sql, ownerFields := HasManyParentFilter(rel)
	if fmt.Sprint(ownerFields) != fmt.Sprint([]string{"OrgID", "Code"}) {
		t.Fatalf("owner fields = %v, want [OrgID Code]", ownerFields)
	}
	for _, col := range []string{"parent_org_id", "parent_code"} {
		if !strings.Contains(sql, col) {
			t.Fatalf("filter must reference fk column %q: %s", col, sql)
		}
	}
	if strings.Count(sql, "?") != 2 {
		t.Fatalf("filter must have two placeholders (composite fk): %s", sql)
	}
}

// TestForeignKeyFields_Composite proves belongs-to foreign-key resolution returns
// every FK column in related-primary-key order, and that mapping a selected id
// onto those columns via model.ID.Related writes them all.
func TestForeignKeyFields_Composite(t *testing.T) {
	type Target struct {
		OrgID uint   `gorm:"primaryKey"`
		Code  string `gorm:"primaryKey"`
		Name  string
	}
	type Owner struct {
		ID          uint
		TargetOrgID uint
		TargetCode  string
		Target      Target `gorm:"foreignKey:TargetOrgID,TargetCode;references:OrgID,Code"`
	}
	db := newDB(t, &Target{}, &Owner{})

	got := ForeignKeyFields(db, &Owner{}, "Target")
	if want := []string{"TargetOrgID", "TargetCode"}; fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("ForeignKeyFields = %v, want %v", got, want)
	}

	id := model.ID{
		Fields: model.Fields{model.SingleField("OrgID"), model.SingleField("Code")},
		Values: []any{uint(7), "abc"},
	}
	owner := &Owner{}
	id.Related(ownerSchema{}, got...).SetTo(owner)
	if owner.TargetOrgID != 7 || owner.TargetCode != "abc" {
		t.Fatalf("owner FK not set from composite id: %+v", owner)
	}
}

func TestPKHelpers(t *testing.T) {
	type Row struct {
		A uint
		B string
	}
	if !PKAllZero(reflect.ValueOf(Row{}), []string{"A", "B"}) {
		t.Error("PKAllZero should be true for a zero row")
	}
	if PKAllZero(reflect.ValueOf(Row{A: 1}), []string{"A", "B"}) {
		t.Error("PKAllZero should be false when any key is set")
	}
	k1 := PKMapKey(reflect.ValueOf(Row{A: 3, B: "x"}), []string{"A", "B"})
	k2 := PKMapKey(reflect.ValueOf(Row{A: 3, B: "x"}), []string{"A", "B"})
	if k1 != k2 {
		t.Errorf("PKMapKey not stable: %v != %v", k1, k2)
	}
	if k1 == PKMapKey(reflect.ValueOf(Row{A: 33, B: "x"}), []string{"A", "B"}) {
		t.Error("PKMapKey must distinguish different composite keys")
	}
	if NormalizeKey(reflect.ValueOf(uint(5))) != int64(5) {
		t.Error("NormalizeKey should normalize uint to int64")
	}
}

// ownerSchema is a minimal model.Schema exposing owner fields for
// model.ID.Related (which only needs FieldsByName).
type ownerSchema struct{}

func (ownerSchema) Model() any                  { return nil }
func (ownerSchema) Table() string               { return "owners" }
func (ownerSchema) QuotedTable() string         { return "owners" }
func (ownerSchema) Fields() model.Fields        { return nil }
func (ownerSchema) PrimaryFields() model.Fields { return model.Fields{model.SingleField("ID")} }
func (ownerSchema) FieldByName(name string) model.Field {
	return model.SingleField(name)
}
func (ownerSchema) FieldsByName(names ...string) model.Fields {
	out := make(model.Fields, len(names))
	for i, n := range names {
		out[i] = model.SingleField(n)
	}
	return out
}
