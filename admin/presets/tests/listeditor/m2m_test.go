package listeditor

import (
	"net/url"
	"sort"
	"testing"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type M2MTag struct {
	ID   uint
	Name string
}

type M2MPost struct {
	ID   uint
	Name string
	Tags []*M2MTag `gorm:"many2many:m2m_post_tags;"`
}

func postTags(t *testing.T, db *gorm.DB, postID uint) []string {
	t.Helper()
	var p M2MPost
	if err := db.Preload("Tags").First(&p, postID).Error; err != nil {
		t.Fatal(err)
	}
	names := make([]string, 0, len(p.Tags))
	for _, tg := range p.Tags {
		names = append(names, tg.Name)
	}
	sort.Strings(names)
	return names
}

// TestManyToMany_AssociateDisassociate drives SaveHasManyAssociation over a
// many-to-many relation: an existing tag is associated, another disassociated
// (removed), keeping the tag rows themselves intact. Association classification
// still uses the __new/__deleted flags; Replace manages the join table.
func TestManyToMany_AssociateDisassociate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&M2MPost{}, &M2MTag{}); err != nil {
		t.Fatal(err)
	}

	tags := []*M2MTag{{ID: 1, Name: "go"}, {ID: 2, Name: "rust"}, {ID: 3, Name: "js"}}
	if err = db.Create(&tags).Error; err != nil {
		t.Fatal(err)
	}
	// post initially tagged go(1) + rust(2)
	post := &M2MPost{ID: 1, Name: "p", Tags: []*M2MTag{tags[0], tags[1]}}
	if err = db.Create(post).Error; err != nil {
		t.Fatal(err)
	}

	// submitted: keep go(1), remove rust(2), add js(3). Existing tags carry only
	// their ID (a selector), so they are not flagged __new.
	vals := url.Values{
		"Tags.__present":    {"1"},
		"Tags[0].ID":        {"1"},
		"Tags[0].__pos":     {"0"},
		"Tags[1].ID":        {"2"},
		"Tags[1].__pos":     {"1"},
		"Tags[1].__deleted": {"true"},
		"Tags[2].ID":        {"3"},
		"Tags[2].__pos":     {"2"},
	}
	obj := &M2MPost{ID: 1, Name: "p", Tags: []*M2MTag{
		{ID: 1}, {ID: 2}, {ID: 3},
	}}

	ctx := ctxWithForm(vals)
	op := gorm2op.SaveHasManyAssociation("Tags").Build(
		gorm2op.DataOperator(db).SetUpdator(func(db *gorm.DB, obj interface{}, id model.ID, ctx *web.EventContext) error {
			return db.Model(obj).Where("id = ?", obj.(*M2MPost).ID).Updates(map[string]any{"name": obj.(*M2MPost).Name}).Error
		}),
	)
	if err = op.Update(obj, model.NewID(post.ID), ctx); err != nil {
		t.Fatalf("update: %v", err)
	}

	// the post is now tagged go + js (rust disassociated)
	if got, want := postTags(t, db, 1), []string{"go", "js"}; !equalStr(got, want) {
		t.Fatalf("post tags = %v, want %v", got, want)
	}

	// all three tag rows still exist with their names intact (a selected tag is
	// associated, never overwritten with empty fields, and a removed tag is only
	// disassociated, not deleted)
	var allTags []M2MTag
	if err = db.Order("id").Find(&allTags).Error; err != nil {
		t.Fatal(err)
	}
	if len(allTags) != 3 {
		t.Fatalf("expected 3 tags to remain, got %d: %+v", len(allTags), allTags)
	}
	names := map[uint]string{}
	for _, tg := range allTags {
		names[tg.ID] = tg.Name
	}
	if names[1] != "go" || names[2] != "rust" || names[3] != "js" {
		t.Errorf("tag names were altered: %+v", allTags)
	}
}

func equalStr(a, b []string) bool {
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
