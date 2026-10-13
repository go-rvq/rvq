package presets

import (
	"reflect"
	"testing"

	"github.com/go-rvq/rvq/admin/softdelete"
	"github.com/go-rvq/rvq/x/i18n"
)

// The fields of the record's times and of its deletion — who deleted it, from
// where (softdelete.Deletion) — are the system's: in the detail only, audited,
// never in a form a user writes; any other field as it was.
func TestDefaultModelSetupSystemFields(t *testing.T) {
	for _, name := range []string{"CreatedAt", "UpdatedAt", "DeletedAt", "DeletedByID", "DeletedOrigin"} {
		f := &FieldBuilder{}
		f.name = name
		defaultModelSetup{}.InitField(EDIT, f)
		if f.Mode() != DETAIL || !f.Audited() {
			t.Errorf("%s: mode %v audited %v, want the detail's, audited", name, f.Mode(), f.Audited())
		}
	}
	f := &FieldBuilder{}
	f.name = "Title"
	defaultModelSetup{}.InitField(EDIT, f)
	if f.Mode() == DETAIL || f.Audited() {
		t.Errorf("Title: mode %v audited %v", f.Mode(), f.Audited())
	}
}

type deletedRecord struct {
	ID    uint
	Title string
	softdelete.Deletion
}

// A model's form has none of the fields of its deletion; a listing that asks
// for them by name — the trash's columns — has them.
func TestDeletionFieldsOutOfTheForm(t *testing.T) {
	mb := New(i18n.New()).URIPrefix("/admin").Model(&deletedRecord{})
	ed := mb.Editing()
	for _, name := range []string{"DeletedAt", "DeletedByID", "DeletedOrigin"} {
		if ed.GetField(name) != nil {
			t.Errorf("%s in the form", name)
		}
		// of the model's field (its type), not a virtual one (a string)
		if f := mb.Listing().Field(name); f == nil || f.rt == nil || f.rt.Kind() == reflect.String {
			t.Errorf("%s: not a column of the listing", name)
		}
	}
	if ed.GetField("Title") == nil {
		t.Error("Title not in the form")
	}
}
