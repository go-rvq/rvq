package helper

import (
	"fmt"
	"testing"

	"github.com/go-rvq/rvq/admin/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SelTarget has a COMPOSITE primary key (OrgID, Code) of mixed types.
type SelTarget struct {
	OrgID uint   `gorm:"primaryKey"`
	Code  string `gorm:"primaryKey"`
	Name  string
}

// SelOwner belongs-to SelTarget through a COMPOSITE foreign key.
type SelOwner struct {
	ID          uint
	TargetOrgID uint
	TargetCode  string
	Target      SelTarget `gorm:"foreignKey:TargetOrgID,TargetCode;references:OrgID,Code"`
}

// TestForeignKeyFieldsOf_Composite proves the belongs-to foreign-key resolution
// returns every FK column in related-primary-key order, so a composite foreign
// key is written in full by the selector.
func TestForeignKeyFieldsOf_Composite(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&SelTarget{}, &SelOwner{}); err != nil {
		t.Fatal(err)
	}

	got := foreignKeyFieldsOf(db, &SelOwner{}, "Target")
	if want := []string{"TargetOrgID", "TargetCode"}; fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("foreignKeyFieldsOf = %v, want %v", got, want)
	}

	// the resolved FK fields, in related-PK order, must map a selected target's
	// composite id onto the owner via model.ID.Related (what the setter does).
	s := selOwnerSchema{}
	id := model.ID{
		Fields: model.Fields{model.SingleField("OrgID"), model.SingleField("Code")},
		Values: []any{uint(7), "abc"},
	}
	owner := &SelOwner{}
	id.Related(s, got...).SetTo(owner)
	if owner.TargetOrgID != 7 || owner.TargetCode != "abc" {
		t.Fatalf("owner FK not set from composite id: %+v", owner)
	}
}

// selOwnerSchema is a minimal model.Schema exposing the owner FK fields for
// model.ID.Related (which only needs FieldsByName).
type selOwnerSchema struct{}

func (selOwnerSchema) Model() any                  { return &SelOwner{} }
func (selOwnerSchema) Table() string               { return "sel_owners" }
func (selOwnerSchema) QuotedTable() string         { return "sel_owners" }
func (selOwnerSchema) Fields() model.Fields        { return nil }
func (selOwnerSchema) PrimaryFields() model.Fields { return model.Fields{model.SingleField("ID")} }
func (selOwnerSchema) FieldByName(name string) model.Field {
	return model.SingleField(name)
}
func (s selOwnerSchema) FieldsByName(names ...string) model.Fields {
	out := make(model.Fields, len(names))
	for i, n := range names {
		out[i] = model.SingleField(n)
	}
	return out
}
