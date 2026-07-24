package listeditor

import (
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"testing"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// CKChild is a has-many child with a COMPOSITE primary key: (OwnerID, Code). One
// key field is the foreign key, the other a natural code — so a new row has a
// non-zero Code but a zero OwnerID until the association assigns it.
type CKChild struct {
	OwnerID uint   `gorm:"primaryKey"`
	Code    string `gorm:"primaryKey"`
	Label   string
}

type CKOwner struct {
	ID       uint
	Name     string
	Children []*CKChild `gorm:"foreignKey:OwnerID"`
}

// ctxWithForm builds an EventContext whose multipart form carries the submitted
// list values (the same place UnmarshalForm reindexes by __pos and the list
// editor / persistence layer read per-item metadata from).
func ctxWithForm(vals url.Values) *web.EventContext {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.MultipartForm = &multipart.Form{Value: vals}
	return &web.EventContext{R: req}
}

// recordingAuditor captures the create/update/delete calls, keyed by Label, so a
// test can assert the reconciliation exercised the composite-key matching
// (snapshotExisting + pkMapKey + keptIDs).
type recordingAuditor struct {
	created []string
	updated [][2]string
	deleted []string
}

func (a *recordingAuditor) LogCreated(_ *gorm.DB, obj any) error {
	a.created = append(a.created, obj.(*CKChild).Label)
	return nil
}

func (a *recordingAuditor) LogUpdated(_ *gorm.DB, old, now any) error {
	a.updated = append(a.updated, [2]string{old.(*CKChild).Label, now.(*CKChild).Label})
	return nil
}

func (a *recordingAuditor) LogDeleted(_ *gorm.DB, obj any) error {
	a.deleted = append(a.deleted, obj.(*CKChild).Label)
	return nil
}

// TestCompositeKey_UpdateCreateDelete reconciles a composite-key has-many in one
// save: one child updated (both PK fields set), one created (__new, zero FK), one
// deleted — with an auditor so the composite-key snapshot/matching path runs.
func TestCompositeKey_UpdateCreateDelete(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&CKOwner{}, &CKChild{}); err != nil {
		t.Fatal(err)
	}
	owner := &CKOwner{ID: 1, Name: "o", Children: []*CKChild{
		{OwnerID: 1, Code: "A", Label: "alpha"},
		{OwnerID: 1, Code: "B", Label: "beta"},
	}}
	if err = db.Create(owner).Error; err != nil {
		t.Fatal(err)
	}

	// submitted: A -> "alpha2" (update), B deleted, one new "gamma" (Code=C, __new).
	vals := url.Values{
		"Children.__present":    {"1"},
		"Children[0].OwnerID":   {"1"},
		"Children[0].Code":      {"A"},
		"Children[0].__pos":     {"0"},
		"Children[1].OwnerID":   {"1"},
		"Children[1].Code":      {"B"},
		"Children[1].__pos":     {"1"},
		"Children[1].__deleted": {"true"},
		"Children[2].Code":      {"C"},
		"Children[2].__pos":     {"2"},
		"Children[2].__new":     {"true"},
	}
	obj := &CKOwner{ID: 1, Name: "o", Children: []*CKChild{
		{OwnerID: 1, Code: "A", Label: "alpha2"},
		{OwnerID: 1, Code: "B"},     // deleted: fields not submitted
		{Code: "C", Label: "gamma"}, // new: zero OwnerID (set by the association)
	}}

	aud := &recordingAuditor{}
	ctx := ctxWithForm(vals)
	ctx.R = ctx.R.WithContext(gorm2op.ContextWithAssociationAuditor(ctx.R.Context(), aud))

	op := gorm2op.SaveHasManyAssociation("Children").Build(
		gorm2op.DataOperator(db).SetUpdator(func(db *gorm.DB, obj interface{}, id model.ID, ctx *web.EventContext) error {
			return db.Model(obj).Where("id = ?", obj.(*CKOwner).ID).Updates(map[string]any{"name": obj.(*CKOwner).Name}).Error
		}),
	)
	if err = op.Update(obj, model.NewID(owner.ID), ctx); err != nil {
		t.Fatalf("update: %v", err)
	}

	// DB: A updated, B gone, C created — all under owner 1.
	var children []CKChild
	if err = db.Where("owner_id = ?", 1).Order("code").Find(&children).Error; err != nil {
		t.Fatal(err)
	}
	got := map[string]string{}
	for _, c := range children {
		got[c.Code] = c.Label
	}
	if len(children) != 2 || got["A"] != "alpha2" || got["C"] != "gamma" {
		t.Fatalf("children = %+v, want A=alpha2, C=gamma", children)
	}

	// Audit: the composite key matched the right rows.
	sort.Strings(aud.created)
	if len(aud.created) != 1 || aud.created[0] != "gamma" {
		t.Errorf("created = %v, want [gamma]", aud.created)
	}
	if len(aud.deleted) != 1 || aud.deleted[0] != "beta" {
		t.Errorf("deleted = %v, want [beta]", aud.deleted)
	}
	if len(aud.updated) != 1 || aud.updated[0] != [2]string{"alpha", "alpha2"} {
		t.Errorf("updated = %v, want [[alpha alpha2]]", aud.updated)
	}
}
