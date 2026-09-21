package listeditor

import (
	"net/http"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	. "github.com/go-rvq/rvq/admin/presets/integration"
	. "github.com/go-rvq/rvq/web/multipartestutils"
)

// TestAddRow_FlagsNewItemAndPersists covers the add-row event end to end:
//
//  1. two persisted items are submitted; the add event appends a fresh row;
//  2. the re-rendered Setup must flag the appended row as __new (so the browser
//     posts __new on the next submit) — the flag has to be set on the same values
//     the render reads (the reindexed multipart values), not ctx.R.Form;
//  3. simulating that next save (the appended row posts __new with a zero ID)
//     creates the row instead of failing with "WHERE conditions required".
func TestAddRow_FlagsNewItemAndPersists(t *testing.T) {
	app, err := twoPersistedItems()
	if err != nil {
		t.Fatal(err)
	}

	// (1)+(2) the add event flags the appended row (slice index 2) as new.
	RunCase(t, TestCase{
		Name: "add row flags the appended item as __new",
		ReqFunc: func() *http.Request {
			return signed(t, app, NewMultipartBuilder().
				PageURL("/admin/products").
				EventFunc(actions.AddRowEvent).
				Query(presets.ParamID, "1").
				Query(presets.ParamAddRowFormKey, "Items").
				AddField("Items.__present", "1").
				AddField("Items[0].ID", "1").
				AddField("Items[0].Label", "A").
				AddField("Items[0].__index", "0").
				AddField("Items[0].__pos", "0").
				AddField("Items[1].ID", "2").
				AddField("Items[1].Label", "B").
				AddField("Items[1].__index", "1").
				AddField("Items[1].__pos", "1").
				BuildEventFuncRequest())
		},
		EventResponseMatch: func(t *testing.T, er *TestEventResponse) {
			bodyContainsAll(t, er.Body,
				`const newKeys = ["Items[2]"]`,
				`newKeys && newKeys.forEach((k) => { form[k + ".__new"] = true })`,
			)
		},
	}, app)

	// (3) the browser then saves; the appended row posts __new with ID=0 and must
	// be created (classification by __new, not by the zero primary key).
	db, err := NewDB()
	if err != nil {
		t.Fatal(err)
	}
	if err = db.Create(&Product{ID: 1, Name: "P1", Items: []*Item{
		{ID: 1, ProductID: 1, Label: "A", Pos: 0},
	}}).Error; err != nil {
		t.Fatal(err)
	}
	app2 := NewApp(db)

	RunCase(t, TestCase{
		Name: "save with appended __new row (zero ID) creates it",
		ReqFunc: func() *http.Request {
			return signed(t, app2, NewMultipartBuilder().
				PageURL("/admin/products").
				EventFunc(actions.Update).
				Query(presets.ParamID, "1").
				Query(presets.ParamOverlay, string(actions.Dialog)).
				AddField("Name", "P1").
				AddField("Items.__present", "1").
				AddField("Items[0].ID", "1").
				AddField("Items[0].Label", "A").
				AddField("Items[0].__index", "0").
				AddField("Items[0].__pos", "0").
				// appended row: zero ID + __new (what the browser posts after add)
				AddField("Items[1].ID", "0").
				AddField("Items[1].Label", "NEW").
				AddField("Items[1].__index", "1").
				AddField("Items[1].__pos", "1").
				AddField("Items[1].__new", "true").
				BuildEventFuncRequest())
		},
		ExpectRunScriptContainsInOrder: []string{"closer.show = false"},
		EventResponseMatch: func(t *testing.T, er *TestEventResponse) {
			if er.Body != "" && containsAny(er.Body, "WHERE conditions required") {
				t.Fatalf("save failed: %s", er.Body)
			}
			var items []Item
			if err := db.Order("id").Find(&items).Error; err != nil {
				t.Fatal(err)
			}
			if len(items) != 2 {
				t.Fatalf("expected 2 items after save, got %d: %+v", len(items), items)
			}
		},
	}, app2)
}

// TestSaveOK_ZeroPKItemWithoutNewDoesNotCrash reproduces the reported crash: a
// row posted with a zero ID and NO __new (a client that omitted the flag) reached
// the update path and gorm raised "WHERE conditions required". The save must now
// succeed — a row without a primary key is left for Replace to create.
func TestSaveOK_ZeroPKItemWithoutNewDoesNotCrash(t *testing.T) {
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
		Name: "zero-PK row without __new does not crash the save",
		ReqFunc: func() *http.Request {
			return signed(t, app, NewMultipartBuilder().
				PageURL("/admin/products").
				EventFunc(actions.Update).
				Query(presets.ParamID, "1").
				Query(presets.ParamOverlay, string(actions.Dialog)).
				AddField("Name", "P1").
				AddField("Items.__present", "1").
				AddField("Items[0].ID", "1").
				AddField("Items[0].Label", "A").
				AddField("Items[0].__index", "0").
				AddField("Items[0].__pos", "0").
				// a new row with a zero ID and NO __new (the reported request)
				AddField("Items[1].ID", "0").
				AddField("Items[1].Label", "NEW").
				AddField("Items[1].__index", "1").
				AddField("Items[1].__pos", "1").
				BuildEventFuncRequest())
		},
		EventResponseMatch: func(t *testing.T, er *TestEventResponse) {
			if containsAny(er.Body, "WHERE conditions required") {
				t.Fatalf("save crashed: %s", er.Body)
			}
			var items []Item
			if err := db.Order("id").Find(&items).Error; err != nil {
				t.Fatal(err)
			}
			if len(items) != 2 {
				t.Fatalf("expected 2 items (A kept, NEW created via Replace), got %d: %+v", len(items), items)
			}
		},
	}, app)
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		for i := 0; i+len(sub) <= len(s); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
	}
	return false
}
