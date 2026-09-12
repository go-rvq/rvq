package admin

import (
	"net/http/httptest"
	"testing"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// histDoc is a minimal versioned model for the integration tests.
type histDoc struct {
	ID    uint `gorm:"primarykey"`
	Title string
	Body  string
}

func setupHistory(t *testing.T) (*gorm.DB, *presets.ModelBuilder, *ModelHistory) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&histDoc{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	b := presets.New(i18n.New()).URIPrefix("/admin")
	b.DataOperator(gorm2op.DataOperator(db))
	mb := b.Model(&histDoc{}).URIName("docs")
	mb.Editing("Title", "Body")

	mh := New(db).Model(mb).HTMLFields("Body").Build()
	return db, mb, mh
}

func save(t *testing.T, mb *presets.ModelBuilder, obj *histDoc, id model.ID) {
	t.Helper()
	ctx := &web.EventContext{R: httptest.NewRequest("POST", "/admin/docs", nil)}
	if err := mb.Editing().Saver(obj, id, ctx); err != nil {
		t.Fatalf("save: %v", err)
	}
}

// TestCreateRecordsRevision: a create save records the first revision (no parent).
func TestCreateRecordsRevision(t *testing.T) {
	_, mb, mh := setupHistory(t)

	obj := &histDoc{Title: "A", Body: "<p>one</p>"}
	save(t, mb, obj, model.ID{})

	revs, err := mh.Chain(mb.MustRecordID(obj).String())
	if err != nil {
		t.Fatalf("chain: %v", err)
	}
	if len(revs) != 1 {
		t.Fatalf("after create: %d revisions, want 1", len(revs))
	}
	if len(revs[0].Parent) != 0 {
		t.Fatalf("first revision should have no parent, got %x", revs[0].Parent)
	}
}

// TestUpdateChainsRevision: an update save records a second revision whose
// parent is the first, and re-saving unchanged fields records nothing (dedup).
func TestUpdateChainsRevision(t *testing.T) {
	_, mb, mh := setupHistory(t)

	obj := &histDoc{Title: "A", Body: "<p>one</p>"}
	save(t, mb, obj, model.ID{})
	key := mb.MustRecordID(obj).String()

	obj.Title = "B"
	save(t, mb, obj, mb.MustRecordID(obj))

	revs, err := mh.Chain(key)
	if err != nil {
		t.Fatalf("chain: %v", err)
	}
	if len(revs) != 2 {
		t.Fatalf("after update: %d revisions, want 2", len(revs))
	}
	// revs is oldest→newest; the newest's parent is the first's hash.
	if string(revs[1].Parent) != string(revs[0].Hash) {
		t.Fatalf("revision 2 parent = %x, want %x", revs[1].Parent, revs[0].Hash)
	}

	// Saving again with no change to the versioned fields must be a no-op.
	save(t, mb, obj, mb.MustRecordID(obj))
	revs, _ = mh.Chain(key)
	if len(revs) != 2 {
		t.Fatalf("unchanged re-save added a revision: %d, want 2", len(revs))
	}
}

// TestChangedFieldsAndFieldHistory: the first revision lists every versioned
// field; a later one lists only the fields that actually changed, and
// FieldHistory returns exactly the revisions that touched a given field.
func TestChangedFieldsAndFieldHistory(t *testing.T) {
	_, mb, mh := setupHistory(t)

	obj := &histDoc{Title: "A", Body: "<p>one</p>"}
	save(t, mb, obj, model.ID{})
	key := mb.MustRecordID(obj).String()

	// Change only Title.
	obj.Title = "B"
	save(t, mb, obj, mb.MustRecordID(obj))

	revs, _ := mh.Chain(key)
	if len(revs) != 2 {
		t.Fatalf("revisions = %d, want 2", len(revs))
	}
	// First revision: both fields are "changed" (initial state).
	first := revs[0].ChangedFields.Data
	if len(first) != 2 {
		t.Fatalf("first ChangedFields = %v, want [Title Body]", first)
	}
	// Second revision: only Title.
	second := revs[1].ChangedFields.Data
	if len(second) != 1 || second[0] != "Title" {
		t.Fatalf("second ChangedFields = %v, want [Title]", second)
	}

	titleHist, err := mh.FieldHistory(key, "Title")
	if err != nil {
		t.Fatalf("FieldHistory(Title): %v", err)
	}
	if len(titleHist) != 2 {
		t.Fatalf("Title history = %d, want 2", len(titleHist))
	}
	bodyHist, err := mh.FieldHistory(key, "Body")
	if err != nil {
		t.Fatalf("FieldHistory(Body): %v", err)
	}
	if len(bodyHist) != 1 {
		t.Fatalf("Body history = %d, want 1 (only the first touched Body)", len(bodyHist))
	}
}

// TestBackfillFillsLegacyRows: revisions saved before the ChangedFields column
// (simulated by nulling it) are filled by Backfill, computing each from its
// parent, and it is idempotent.
func TestBackfillFillsLegacyRows(t *testing.T) {
	db, mb, mh := setupHistory(t)

	obj := &histDoc{Title: "A", Body: "<p>one</p>"}
	save(t, mb, obj, model.ID{})
	key := mb.MustRecordID(obj).String()
	obj.Title = "B"
	save(t, mb, obj, mb.MustRecordID(obj))

	// Simulate legacy rows: clear the column on every revision.
	if err := db.Table(mh.Table()).Where("record_key = ?", key).
		Update("changed_fields", nil).Error; err != nil {
		t.Fatalf("null out: %v", err)
	}

	if err := mh.Backfill(); err != nil {
		t.Fatalf("backfill: %v", err)
	}

	revs, _ := mh.Chain(key)
	if len(revs[0].ChangedFields.Data) != 2 {
		t.Fatalf("first ChangedFields = %v, want both fields", revs[0].ChangedFields.Data)
	}
	if len(revs[1].ChangedFields.Data) != 1 || revs[1].ChangedFields.Data[0] != "Title" {
		t.Fatalf("second ChangedFields = %v, want [Title]", revs[1].ChangedFields.Data)
	}

	// Idempotent: a second run changes nothing and errors nothing.
	if err := mh.Backfill(); err != nil {
		t.Fatalf("backfill (2nd): %v", err)
	}
}

// TestLatestHashReadsBytea exercises the exact query that regressed: reading the
// newest hash back through Hash's sql.Scanner (a gorm .Scan would misread it).
func TestLatestHashReadsBytea(t *testing.T) {
	_, mb, mh := setupHistory(t)

	obj := &histDoc{Title: "A", Body: "x"}
	save(t, mb, obj, model.ID{})
	key := mb.MustRecordID(obj).String()

	got, err := mh.latestHash(key)
	if err != nil {
		t.Fatalf("latestHash: %v", err)
	}
	revs, _ := mh.Chain(key)
	if string(got) != string(revs[0].Hash) {
		t.Fatalf("latestHash = %x, want %x", got, revs[0].Hash)
	}

	// No revisions for an unknown record → nil, no error.
	none, err := mh.latestHash("999")
	if err != nil {
		t.Fatalf("latestHash(unknown): %v", err)
	}
	if none != nil {
		t.Fatalf("latestHash(unknown) = %x, want nil", none)
	}
}
