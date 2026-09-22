package admin

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupNoRevert versions both fields of histDoc but keeps Body out of every
// revert: it is the field the application owns.
func setupNoRevert(t *testing.T) (*gorm.DB, *presets.ModelBuilder, *ModelHistory) {
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

	mh := New(db).Model(mb).Fields("Title", "Body").NoRevertFields("Body").Build()
	return db, mb, mh
}

func revertCtx() *web.EventContext {
	return &web.EventContext{R: httptest.NewRequest("POST", "/admin/docs", nil)}
}

// A revert of the whole record restores every other field and leaves the
// non-revertible one with the value it has now.
func TestRevertRecordSkipsNoRevertField(t *testing.T) {
	_, mb, mh := setupNoRevert(t)

	obj := &histDoc{Title: "A", Body: "um"}
	save(t, mb, obj, model.ID{})
	key := mb.MustRecordID(obj).String()

	obj.Title, obj.Body = "B", "dois"
	save(t, mb, obj, mb.MustRecordID(obj))

	revs, err := mh.Chain(key)
	if err != nil {
		t.Fatal(err)
	}
	if err := mh.RevertRecord(obj, revs[0].Hash, revertCtx()); err != nil {
		t.Fatal(err)
	}

	if obj.Title != "A" {
		t.Errorf("Title = %q, want a volta para A", obj.Title)
	}
	if obj.Body != "dois" {
		t.Errorf("Body = %q, want o valor de agora (não é reversível)", obj.Body)
	}

	// and what was saved is what the record holds
	var stored histDoc
	if err := mh.db.First(&stored, obj.ID).Error; err != nil {
		t.Fatal(err)
	}
	if stored.Title != "A" || stored.Body != "dois" {
		t.Errorf("gravado = %+v", stored)
	}
}

// Naming the field alone is refused, and nothing is written.
func TestRevertFieldRefusesNoRevertField(t *testing.T) {
	_, mb, mh := setupNoRevert(t)

	obj := &histDoc{Title: "A", Body: "um"}
	save(t, mb, obj, model.ID{})
	key := mb.MustRecordID(obj).String()

	obj.Body = "dois"
	save(t, mb, obj, mb.MustRecordID(obj))

	revs, _ := mh.Chain(key)
	err := mh.RevertField(obj, revs[0].Hash, "Body", revertCtx())
	if err == nil {
		t.Fatal("um field não reversível foi revertido")
	}
	if !strings.Contains(err.Error(), "not revertible") {
		t.Errorf("err = %v", err)
	}
	if obj.Body != "dois" {
		t.Errorf("Body = %q, want intocado", obj.Body)
	}
}

// The partial paths refuse it too — a field that is not revertible accepts no
// revert at all, whole or in hunks.
func TestPartialRevertRefusesNoRevertField(t *testing.T) {
	_, mb, mh := setupNoRevert(t)

	obj := &histDoc{Title: "A", Body: "um"}
	save(t, mb, obj, model.ID{})

	if mh.AcceptsPartial("Body") {
		t.Error("um field não reversível não aceita revert parcial")
	}
	if err := mh.RevertFieldContent(obj, "Body", "outro", revertCtx()); err == nil {
		t.Error("RevertFieldContent devia recusar")
	}
	if err := mh.applyFieldText(obj, "Body", "outro", revertCtx()); err == nil {
		t.Error("applyFieldText devia recusar")
	}
	if obj.Body != "um" {
		t.Errorf("Body = %q, want intocado", obj.Body)
	}

	// the revertible field goes on accepting both
	if !mh.AcceptsPartial("Title") {
		t.Error("o field reversível parou de aceitar revert parcial")
	}
	if !mh.Revertible("Title") || mh.Revertible("Body") {
		t.Error("Revertible não responde o que foi marcado")
	}
}

// Asking for only non-revertible fields does nothing — and records no revision,
// because nothing changed.
func TestRevertFieldsWithOnlyNoRevertFieldsDoesNothing(t *testing.T) {
	_, mb, mh := setupNoRevert(t)

	obj := &histDoc{Title: "A", Body: "um"}
	save(t, mb, obj, model.ID{})
	key := mb.MustRecordID(obj).String()

	obj.Body = "dois"
	save(t, mb, obj, mb.MustRecordID(obj))

	revs, _ := mh.Chain(key)
	before := len(revs)

	if err := mh.RevertFields(obj, revs[0].Hash, []string{"Body"}, revertCtx()); err != nil {
		t.Fatal(err)
	}
	if obj.Body != "dois" {
		t.Errorf("Body = %q, want intocado", obj.Body)
	}

	after, _ := mh.Chain(key)
	if len(after) != before {
		t.Errorf("revisões = %d, want %d: nada mudou", len(after), before)
	}
}
