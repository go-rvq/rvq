package listeditor

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	. "github.com/go-rvq/rvq/admin/presets/integration"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/web/multipartestutils"
)

// The tests below drive the list editor through the exact browser flow and
// assert on the RENDERED HTML (er.Body), not only on the persisted rows. Each
// case documents the corresponding user interaction inline.
//
// Reactive-form recap (what a real browser posts on save). The bracket index is
// the row's stable __index (seeded at init), not its slice position:
//   - `Items.__present` = "1"                -> the list editor was initialized
//   - `Items[x].ID` / `Items[x].Label`       -> the row's hidden PK and real field
//   - `Items[x].__index` = x                 -> stable per-item id / form key
//   - `Items[x].__pos`                       -> display & submit order
//   - `Items[x].__deleted` = "true"          -> the row was removed client-side
//   - `Items[x].__new` = "true"              -> a browser-created row
//
// Removing a row does NOT hit the server; it only flips `__deleted` on the
// reactive form, which travels to the server on the next submit. So a save that
// fails validation must re-render every removed row as deleted and re-seed its
// `__deleted` flag — otherwise the removal is silently lost. The tests below use
// __index == slice position for readability; see index_test.go for the case
// where they differ (a sorted list).

// twoPersistedItems seeds a product (id=1) with two persisted children:
// #0 (ID=1, "A") and #1 (ID=2, "B").
func twoPersistedItems() (*presets.Builder, error) {
	db, err := NewDB()
	if err != nil {
		return nil, err
	}
	err = db.Create(&Product{ID: 1, Name: "P1", Items: []*Item{
		{ID: 1, ProductID: 1, Label: "A", Pos: 0},
		{ID: 2, ProductID: 1, Label: "B", Pos: 1},
	}}).Error
	if err != nil {
		return nil, err
	}
	return NewApp(db), nil
}

// bodyContainsAll asserts the rendered event-response body contains every
// fragment (order-independent), logging the body once on failure.
func bodyContainsAll(t *testing.T, body string, fragments ...string) {
	t.Helper()
	var missing []string
	for _, f := range fragments {
		if !strings.Contains(body, f) {
			missing = append(missing, f)
		}
	}
	if len(missing) > 0 {
		for _, m := range missing {
			t.Errorf("rendered body should contain:\n  %s", m)
		}
		t.Logf("---- rendered body ----\n%s\n-----------------------", body)
	}
}

// bodyContainsInOrder asserts the fragments appear in the body in the given
// order (each after the previous).
func bodyContainsInOrder(t *testing.T, body string, fragments ...string) {
	t.Helper()
	rest := body
	for _, f := range fragments {
		i := strings.Index(rest, f)
		if i < 0 {
			t.Errorf("rendered body should contain (in order):\n  %s", f)
			t.Logf("---- rendered body ----\n%s\n-----------------------", body)
			return
		}
		rest = rest[i+len(f):]
	}
}

func bodyContainsNone(t *testing.T, body string, fragments ...string) {
	t.Helper()
	for _, f := range fragments {
		if strings.Contains(body, f) {
			t.Errorf("rendered body should NOT contain:\n  %s", f)
			t.Logf("---- rendered body ----\n%s\n-----------------------", body)
			return
		}
	}
}

// TestSaveFail_PersistedItemRemoved_KeepsDeletedInUI reproduces:
//
//  1. open the edit form (two persisted items);
//  2. remove item #1 (client sets Items[1].__deleted = true);
//  3. clear the required Name;
//  4. SAVE.
//
// Expected UI after the (failed) save: the Name field shows the required error
// AND item #1 is still rendered as deleted — the component Setup re-seeds its
// __deleted flag so the removal survives the re-render.
func TestSaveFail_PersistedItemRemoved_KeepsDeletedInUI(t *testing.T) {
	app, err := twoPersistedItems()
	if err != nil {
		t.Fatal(err)
	}

	RunCase(t, TestCase{
		Name: "save fails, removed persisted item stays deleted",
		ReqFunc: func() *http.Request {
			return signed(t, app, NewMultipartBuilder().
				PageURL("/admin/products").
				EventFunc(actions.Update).
				Query(presets.ParamID, "1").
				AddField("Name", ""). // required -> validation fails
				AddField("Items.__present", "1").
				AddField("Items[0].ID", "1").
				AddField("Items[0].Label", "A").
				AddField("Items[0].__pos", "0").
				AddField("Items[1].ID", "2").
				AddField("Items[1].Label", "B").
				AddField("Items[1].__pos", "1").
				AddField("Items[1].__deleted", "true"). // removed in the browser
				BuildEventFuncRequest())
		},
		EventResponseMatch: func(t *testing.T, er *TestEventResponse) {
			bodyContainsAll(t, er.Body,
				// (1) the Name required error renders on the re-rendered form
				`v-model='form["Name"]'`,
				`:error-messages='["This field is required"]'`,
				// (2) item #1's deletion is re-seeded by the component Setup
				`const deletedKeys = ["Items[1]"]`,
				`deletedKeys && deletedKeys.forEach((k) => { form[k + ".__deleted"] = true })`,
				// (3) item #1 renders its "deleted" placeholder row
				`v-show='form["Items[1].__deleted"] && !form["Items[1].__purged"]'`,
			)
			// item #0 was NOT removed, so it must not be re-seeded as deleted
			bodyContainsNone(t, er.Body, `const deletedKeys = ["Items[0]`)
		},
	}, app)
}

// TestSaveFail_NewItemRemoved_VanishesFromUI covers the harder case:
//
//  1. open the edit form (two persisted items);
//  2. add a brand-new row (no ID, __new) at index 2;
//  3. remove that new row (Items[2].__deleted = true);
//  4. clear the required Name;
//  5. SAVE.
//
// Expected: a browser-created row that was removed has no persisted identity, so
// it VANISHES entirely — it is not rendered, not re-seeded, and is pruned from
// the reactive form (so it cannot leak into another form). The Name error still
// renders and the two persisted items are untouched.
func TestSaveFail_NewItemRemoved_VanishesFromUI(t *testing.T) {
	app, err := twoPersistedItems()
	if err != nil {
		t.Fatal(err)
	}

	RunCase(t, TestCase{
		Name: "save fails, removed new item vanishes",
		ReqFunc: func() *http.Request {
			return signed(t, app, NewMultipartBuilder().
				PageURL("/admin/products").
				EventFunc(actions.Update).
				Query(presets.ParamID, "1").
				AddField("Name", ""). // required -> validation fails
				AddField("Items.__present", "1").
				AddField("Items[0].ID", "1").
				AddField("Items[0].Label", "A").
				AddField("Items[0].__index", "0").
				AddField("Items[0].__pos", "0").
				AddField("Items[1].ID", "2").
				AddField("Items[1].Label", "B").
				AddField("Items[1].__index", "1").
				AddField("Items[1].__pos", "1").
				// a browser-created row (no ID, __new) that was then removed
				AddField("Items[2].Label", "C").
				AddField("Items[2].__index", "2").
				AddField("Items[2].__pos", "2").
				AddField("Items[2].__new", "true").
				AddField("Items[2].__deleted", "true").
				BuildEventFuncRequest())
		},
		EventResponseMatch: func(t *testing.T, er *TestEventResponse) {
			bodyContainsAll(t, er.Body,
				// the Name required error still renders
				`:error-messages='["This field is required"]'`,
			)
			// the removed new item is gone: not rendered, not re-seeded
			bodyContainsNone(t, er.Body,
				`form["Items[2].Label"]`,
				`Items[2].__deleted`,
				`const deletedKeys = [`, // nothing is deleted (the only removal vanished)
			)
			// the two persisted items are still rendered and not removed
			bodyContainsAll(t, er.Body,
				`form["Items[0].Label"]`,
				`form["Items[1].Label"]`,
			)
		},
	}, app)
}

// TestSaveOK_DeletedItemSkipsValidation reproduces the bug where a removed row's
// fields were still validated: item #1 is marked deleted and posts an EMPTY
// required Label. The save must succeed (the deleted row is skipped, not
// validated) and item #1 must be removed from the database.
func TestSaveOK_DeletedItemSkipsValidation(t *testing.T) {
	db, err := NewDB()
	if err != nil {
		t.Fatal(err)
	}
	if err = db.Create(&Product{ID: 1, Name: "P1", Items: []*Item{
		{ID: 1, ProductID: 1, Label: "A", Pos: 0},
		{ID: 2, ProductID: 1, Label: "B", Pos: 1},
	}}).Error; err != nil {
		t.Fatal(err)
	}
	app := NewApp(db)

	RunCase(t, TestCase{
		Name: "deleted item with empty required field does not fail validation",
		ReqFunc: func() *http.Request {
			return signed(t, app, NewMultipartBuilder().
				PageURL("/admin/products").
				EventFunc(actions.Update).
				Query(presets.ParamID, "1").
				Query(presets.ParamOverlay, string(actions.Dialog)).
				AddField("Name", "P1"). // valid
				AddField("Items.__present", "1").
				AddField("Items[0].ID", "1").
				AddField("Items[0].Label", "A").
				AddField("Items[0].__pos", "0").
				// item #1 removed, and its required Label posted EMPTY
				AddField("Items[1].ID", "2").
				AddField("Items[1].Label", ""). // empty required field
				AddField("Items[1].__pos", "1").
				AddField("Items[1].__deleted", "true").
				BuildEventFuncRequest())
		},
		// save succeeds -> the dialog closes; no validation error rendered
		ExpectRunScriptContainsInOrder: []string{"closer.show = false"},
		EventResponseMatch: func(t *testing.T, er *TestEventResponse) {
			if strings.Contains(er.Body, "This field is required") {
				t.Errorf("deleted item's empty required field should not be validated; body: %s", er.Body)
			}
			var items []Item
			if err := db.Order("id").Find(&items).Error; err != nil {
				t.Fatal(err)
			}
			if len(items) != 1 || items[0].ID != 1 {
				t.Errorf("expected only item #0 to remain, got %+v", items)
			}
		},
	}, app)
}

// TestSaveOK_DeletedItemRemovedFromDB validates the persistence side: a
// successful save (valid Name) with item #1 (ID=2) removed drops it from the
// database while keeping item #0 (ID=1).
func TestSaveOK_DeletedItemRemovedFromDB(t *testing.T) {
	db, err := NewDB()
	if err != nil {
		t.Fatal(err)
	}
	if err = db.Create(&Product{ID: 1, Name: "P1", Items: []*Item{
		{ID: 1, ProductID: 1, Label: "A", Pos: 0},
		{ID: 2, ProductID: 1, Label: "B", Pos: 1},
	}}).Error; err != nil {
		t.Fatal(err)
	}
	app := NewApp(db)

	RunCase(t, TestCase{
		Name: "successful save removes the deleted item",
		ReqFunc: func() *http.Request {
			return signed(t, app, NewMultipartBuilder().
				PageURL("/admin/products").
				EventFunc(actions.Update).
				Query(presets.ParamID, "1").
				Query(presets.ParamOverlay, string(actions.Dialog)). // edit runs in a dialog
				AddField("Name", "P1").                              // valid -> save succeeds
				AddField("Items.__present", "1").
				AddField("Items[0].ID", "1").
				AddField("Items[0].Label", "A").
				AddField("Items[0].__pos", "0").
				AddField("Items[1].ID", "2").
				AddField("Items[1].Label", "B").
				AddField("Items[1].__pos", "1").
				AddField("Items[1].__deleted", "true").
				BuildEventFuncRequest())
		},
		// the overlay is closed on success (closer.show = false)
		ExpectRunScriptContainsInOrder: []string{"closer.show = false"},
		EventResponseMatch: func(t *testing.T, er *TestEventResponse) {
			var items []Item
			if err := db.Order("id").Find(&items).Error; err != nil {
				t.Fatal(err)
			}
			if len(items) != 1 {
				t.Fatalf("expected 1 item after save, got %d: %+v", len(items), items)
			}
			if items[0].ID != 1 || items[0].Label != "A" {
				t.Errorf("expected only item #0 (ID=1, A) to remain, got %+v", items[0])
			}
		},
	}, app)
}

// signed is the request a browser would send: the form the test built, plus the
// record stamp the rendered form carried. Product has no UpdatedAt, so the
// stamp is the hash of the fields the form edits — and the record has to be
// read the way the form read it, through the model's own fetcher, or the nested
// items would hash differently (see presets.Builder.SignForm).
func signed(t *testing.T, app *presets.Builder, r *http.Request) *http.Request {
	t.Helper()

	mb := app.GetModel(&Product{})
	if mb == nil {
		t.Fatal("o modelo Product não está registrado")
	}

	mid, err := mb.ParseRecordID(r.URL.Query().Get(presets.ParamID))
	if err != nil {
		t.Fatal(err)
	}

	obj := mb.NewModel()
	if err = mb.Editing().Fetcher(obj, mid, &web.EventContext{R: r}); err != nil {
		t.Fatal(err)
	}

	out, err := app.SignForm(r, obj)
	if err != nil {
		t.Fatal(err)
	}
	return out
}
