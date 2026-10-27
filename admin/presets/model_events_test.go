package presets_test

import (
	"net/http/httptest"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type eventsParent struct {
	ID uint
}

type eventsChild struct {
	ID             uint
	EventsParentID uint
}

// An event of a child model runs with the IDs of its parents, of the path —
// the hrefs it makes have them, not {parent_0_id} —, even when it is
// registered before the model is made a child.
func TestChildEventHasParentsID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&eventsParent{}, &eventsChild{}); err != nil {
		t.Fatal(err)
	}
	db.Create(&eventsParent{ID: 7})
	b := presets.New(i18n.New()).URIPrefix("/admin").DataOperator(gorm2op.DataOperator(db))
	parent := b.Model(&eventsParent{}).URIName("parents")
	child := presets.NewModelBuilder(b, &eventsChild{}).URIName("children")
	var href string
	child.RegisterEventFunc("child_href", func(ctx *web.EventContext) (r web.EventResponse, err error) {
		href = child.Info().ListingHrefCtx(ctx)
		return
	})
	parent.AddChild(child)
	b.Build()
	w := httptest.NewRecorder()
	b.ServeHTTP(w, httptest.NewRequest("POST", "/admin/parents/7/children?__execute_event__=child_href", nil))
	if href != "/admin/parents/7/children" {
		t.Errorf("the href of the event of a child: %q (status %d)", href, w.Code)
	}
}
