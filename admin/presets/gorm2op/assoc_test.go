package gorm2op

import (
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strconv"
	"testing"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/web"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type assocChild struct {
	ID       uint `gorm:"primarykey"`
	ParentID uint
	Name     string
}

type assocParent struct {
	ID       uint `gorm:"primarykey"`
	Title    string
	Children []*assocChild `gorm:"foreignKey:ParentID"`
}

// newAssocOperator builds a data operator whose base update targets the parent
// by its own id. The default update path relies on a schema-bound model.ID that
// the presets layer produces at runtime; here we bypass it so the tests can
// focus on the association reconciliation done by SaveHasManyAssociation.
func newAssocOperator(db *gorm.DB) *DataOperatorBuilder {
	return DataOperator(db).SetUpdator(func(db *gorm.DB, obj interface{}, id model.ID, ctx *web.EventContext) error {
		return db.Model(obj).Where("id = ?", obj.(*assocParent).ID).Updates(obj).Error
	})
}

func newAssocDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&assocParent{}, &assocChild{}); err != nil {
		t.Fatal(err)
	}
	return db
}

// ctxWithForm builds an EventContext whose multipart form carries the submitted
// list values (the same place UnmarshalForm reindexes by __pos).
func ctxWithForm(vals url.Values) *web.EventContext {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.MultipartForm = &multipart.Form{Value: vals}
	return &web.EventContext{R: req}
}

func childrenOf(t *testing.T, db *gorm.DB, parentID uint) []*assocChild {
	t.Helper()
	var out []*assocChild
	if err := db.Where("parent_id = ?", parentID).Order("id").Find(&out).Error; err != nil {
		t.Fatal(err)
	}
	return out
}

// fakeAuditor records the calls made by SaveHasManyAssociation.
type fakeAuditor struct {
	created []string
	updated [][2]string // old name, new name
	deleted []string
}

func (a *fakeAuditor) LogCreated(_ *gorm.DB, obj any) error {
	a.created = append(a.created, obj.(*assocChild).Name)
	return nil
}

func (a *fakeAuditor) LogUpdated(_ *gorm.DB, old, now any) error {
	a.updated = append(a.updated, [2]string{old.(*assocChild).Name, now.(*assocChild).Name})
	return nil
}

func (a *fakeAuditor) LogDeleted(_ *gorm.DB, obj any) error {
	a.deleted = append(a.deleted, obj.(*assocChild).Name)
	return nil
}

// TestSaveHasManyAssociation_UpdateCreateDelete covers the mixed case: one child
// updated, one created, one deleted — all in a single submit, with auditing.
func TestSaveHasManyAssociation_UpdateCreateDelete(t *testing.T) {
	db := newAssocDB(t)
	parent := &assocParent{Title: "p", Children: []*assocChild{
		{Name: "a"},
		{Name: "b"},
	}}
	if err := db.Create(parent).Error; err != nil {
		t.Fatal(err)
	}
	c1, c2 := parent.Children[0].ID, parent.Children[1].ID

	// submitted: c1 -> "a2" (edit), c2 deleted (posts only ID+__deleted), one new
	// "c" flagged __new (creation is classified by __new, not by a zero PK).
	vals := url.Values{
		"Children.__present":    {"1"},
		"Children[0].ID":        {itoa(c1)},
		"Children[1].ID":        {itoa(c2)},
		"Children[1].__deleted": {"true"},
		"Children[2].__new":     {"true"},
	}
	obj := &assocParent{ID: parent.ID, Title: "p", Children: []*assocChild{
		{ID: c1, Name: "a2"},
		{ID: c2}, // deleted: real fields are not submitted
		{Name: "c"},
	}}

	aud := &fakeAuditor{}
	ctx := ctxWithForm(vals)
	ctx.R = ctx.R.WithContext(ContextWithAssociationAuditor(ctx.R.Context(), aud))

	d := SaveHasManyAssociation("Children").Build(newAssocOperator(db))
	if err := d.Update(obj, model.NewID(parent.ID), ctx); err != nil {
		t.Fatal(err)
	}

	got := names(childrenOf(t, db, parent.ID))
	if want := []string{"a2", "c"}; !equal(got, want) {
		t.Fatalf("children = %v, want %v", got, want)
	}

	if want := []string{"c"}; !equal(aud.created, want) {
		t.Errorf("created = %v, want %v", aud.created, want)
	}
	if want := []string{"b"}; !equal(aud.deleted, want) {
		t.Errorf("deleted = %v, want %v", aud.deleted, want)
	}
	if len(aud.updated) != 1 || aud.updated[0] != [2]string{"a", "a2"} {
		t.Errorf("updated = %v, want [[a a2]]", aud.updated)
	}
}

// TestSaveHasManyAssociation_PresentButEmpty deletes every child when the list
// was in the form but empty.
func TestSaveHasManyAssociation_PresentButEmpty(t *testing.T) {
	db := newAssocDB(t)
	parent := &assocParent{Title: "p", Children: []*assocChild{{Name: "a"}, {Name: "b"}}}
	if err := db.Create(parent).Error; err != nil {
		t.Fatal(err)
	}

	vals := url.Values{"Children.__present": {"1"}}
	obj := &assocParent{ID: parent.ID, Title: "p"} // no children

	aud := &fakeAuditor{}
	ctx := ctxWithForm(vals)
	ctx.R = ctx.R.WithContext(ContextWithAssociationAuditor(ctx.R.Context(), aud))

	d := SaveHasManyAssociation("Children").Build(newAssocOperator(db))
	if err := d.Update(obj, model.NewID(parent.ID), ctx); err != nil {
		t.Fatal(err)
	}

	if got := childrenOf(t, db, parent.ID); len(got) != 0 {
		t.Fatalf("children = %v, want none", names(got))
	}
	sort.Strings(aud.deleted)
	if want := []string{"a", "b"}; !equal(aud.deleted, want) {
		t.Errorf("deleted = %v, want %v", aud.deleted, want)
	}
}

// TestSaveHasManyAssociation_AbsentLeavesUntouched keeps the children when the
// field was not part of the form at all.
func TestSaveHasManyAssociation_AbsentLeavesUntouched(t *testing.T) {
	db := newAssocDB(t)
	parent := &assocParent{Title: "p", Children: []*assocChild{{Name: "a"}, {Name: "b"}}}
	if err := db.Create(parent).Error; err != nil {
		t.Fatal(err)
	}

	obj := &assocParent{ID: parent.ID, Title: "p2"} // no children, no form marker

	d := SaveHasManyAssociation("Children").Build(newAssocOperator(db))
	if err := d.Update(obj, model.NewID(parent.ID), ctxWithForm(url.Values{})); err != nil {
		t.Fatal(err)
	}

	if got := names(childrenOf(t, db, parent.ID)); !equal(got, []string{"a", "b"}) {
		t.Fatalf("children = %v, want [a b] untouched", got)
	}
}

// TestSaveHasManyAssociation_CreateParentWithChildren covers the create flow.
func TestSaveHasManyAssociation_CreateParentWithChildren(t *testing.T) {
	db := newAssocDB(t)

	vals := url.Values{
		"Children.__present": {"1"},
		"Children[0].ID":     {"0"},
		"Children[0].__new":  {"true"},
		"Children[1].ID":     {"0"},
		"Children[1].__new":  {"true"},
	}
	obj := &assocParent{Title: "p", Children: []*assocChild{{Name: "x"}, {Name: "y"}}}

	d := SaveHasManyAssociation("Children").Build(newAssocOperator(db))
	if err := d.Save(obj, model.ID{}, ctxWithForm(vals)); err != nil {
		t.Fatal(err)
	}
	if obj.ID == 0 {
		t.Fatal("parent id not set")
	}
	if got := names(childrenOf(t, db, obj.ID)); !equal(got, []string{"x", "y"}) {
		t.Fatalf("children = %v, want [x y]", got)
	}
}

func names(cs []*assocChild) []string {
	out := make([]string, len(cs))
	for i, c := range cs {
		out[i] = c.Name
	}
	return out
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func itoa(v uint) string {
	return strconv.FormatUint(uint64(v), 10)
}
