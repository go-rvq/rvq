package listeditor

import (
	"net/http"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	. "github.com/go-rvq/rvq/admin/presets/integration"
	. "github.com/go-rvq/rvq/web/multipartestutils"
)

// TestIndexStable_FormKeysFollowIndexNotPosition submits two items reordered by
// __pos (B shown first, A second) while their stable __index stays A=0 / B=1. On
// the re-render each item must keep its form key from its __index — B under
// Items[1], A under Items[0] — even though B is displayed first. This is what
// lets a row keep its identity across a sort instead of inheriting the metadata
// of whatever now sits at its old array position.
func TestIndexStable_FormKeysFollowIndexNotPosition(t *testing.T) {
	app, err := twoPersistedItems()
	if err != nil {
		t.Fatal(err)
	}

	RunCase(t, TestCase{
		Name: "form keys follow __index, not slice position",
		ReqFunc: func() *http.Request {
			return signed(t, app, NewMultipartBuilder().
				PageURL("/admin/products").
				EventFunc(actions.Update).
				Query(presets.ParamID, "1").
				AddField("Name", ""). // fail validation -> re-render
				AddField("Items.__present", "1").
				// slice pos 0 = A, stable __index 0, shown 2nd (__pos 1)
				AddField("Items[0].ID", "1").
				AddField("Items[0].Label", "A").
				AddField("Items[0].__index", "0").
				AddField("Items[0].__pos", "1").
				// slice pos 1 = B, stable __index 1, shown 1st (__pos 0)
				AddField("Items[1].ID", "2").
				AddField("Items[1].Label", "B").
				AddField("Items[1].__index", "1").
				AddField("Items[1].__pos", "0").
				BuildEventFuncRequest())
		},
		EventResponseMatch: func(t *testing.T, er *TestEventResponse) {
			bodyContainsAll(t, er.Body,
				// each item keeps the form key derived from its __index
				`v-model='form["Items[1].Label"]' v-assign='[form, {"Items[1].Label": "B"}]`,
				`v-model='form["Items[0].Label"]' v-assign='[form, {"Items[0].Label": "A"}]`,
			)
			// B (Items[1]) is rendered before A (Items[0]) — display follows __pos,
			// so the stable-key item Items[1] appears first in the Setup seed list.
			bodyContainsInOrder(t, er.Body,
				`const keys = ["Items[1]`,
				`"Items[0]"`,
			)
		},
	}, app)
}

// TestSaveOK_NewItemCreated saves (successfully) an existing item plus a
// browser-created row flagged __new with no ID. The new row must be persisted —
// classification uses __new, never the (zero) primary key.
func TestSaveOK_NewItemCreated(t *testing.T) {
	db, err := NewDB()
	if err != nil {
		t.Fatal(err)
	}
	if err = db.Create(&Product{ID: 1, Name: "P1", Items: []*Item{
		{ID: 1, ProductID: 1, Label: "A", Pos: 0},
	}}).Error; err != nil {
		t.Fatal(err)
	}
	app := NewApp(db)

	RunCase(t, TestCase{
		Name: "new item (flagged __new, no ID) is created",
		ReqFunc: func() *http.Request {
			return signed(t, app, NewMultipartBuilder().
				PageURL("/admin/products").
				EventFunc(actions.Update).
				Query(presets.ParamID, "1").
				Query(presets.ParamOverlay, string(actions.Dialog)).
				AddField("Name", "P1"). // valid -> save succeeds
				AddField("Items.__present", "1").
				AddField("Items[0].ID", "1").
				AddField("Items[0].Label", "A").
				AddField("Items[0].__index", "0").
				AddField("Items[0].__pos", "0").
				// a browser-created row: no ID, flagged __new
				AddField("Items[1].Label", "NEW").
				AddField("Items[1].__index", "1").
				AddField("Items[1].__pos", "1").
				AddField("Items[1].__new", "true").
				BuildEventFuncRequest())
		},
		ExpectRunScriptContainsInOrder: []string{"closer.show = false"},
		EventResponseMatch: func(t *testing.T, er *TestEventResponse) {
			var items []Item
			if err := db.Order("id").Find(&items).Error; err != nil {
				t.Fatal(err)
			}
			if len(items) != 2 {
				t.Fatalf("expected 2 items after save, got %d: %+v", len(items), items)
			}
			got := map[string]bool{}
			for _, it := range items {
				got[it.Label] = true
			}
			if !got["A"] || !got["NEW"] {
				t.Errorf("expected items A and NEW, got %+v", items)
			}
		},
	}, app)
}
