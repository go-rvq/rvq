package integration_test

import (
	"fmt"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SelNote struct {
	ID    uint
	Name  string
	Kind  string
	Value int
}

var selSeq int64

// The listing opened to choose several records: a column to mark them, the
// ones already chosen fixed, only the marked ones by a button (where they are
// unmarked), the button that gives them to who opened it; its condition and
// the filter it hides; no aggregates.
func TestListingSelectMany(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:selmany_%d?mode=memory&cache=shared", atomic.AddInt64(&selSeq, 1))), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&SelNote{}); err != nil {
		t.Fatal(err)
	}
	db.Create(&[]SelNote{
		{ID: 1, Name: "zz-one", Kind: "a", Value: 10},
		{ID: 2, Name: "zz-two", Kind: "a", Value: 20},
		{ID: 3, Name: "zz-three", Kind: "a", Value: 30},
		{ID: 4, Name: "zz-out", Kind: "b", Value: 40},
	})

	p := presets.New(i18n.New()).URIPrefix("/admin")
	p.Permission(perm.New().AllowAll())
	p.DataOperator(gorm2op.DataOperator(db))
	l := p.Model(&SelNote{}).URIName("notes").Listing("ID", "Name")
	l.FilterDataFunc(func(ctx *web.EventContext) vx.FilterData {
		return []*vx.FilterItem{{Key: "kind", Label: "Kind-filter", ItemType: vx.ItemTypeSelect, SQLCondition: `kind %s ?`,
			Options: []*vx.SelectItem{{Text: "A", Value: "a"}, {Text: "B", Value: "b"}}}}
	})
	l.AggregateBase("v", "sel_notes.value")
	l.Aggregate("total").Label("Total").Expr(presets.Sum("v"))
	// opened to choose: only kind a, its filter hidden
	l.SelectionConditionFunc(func(ctx *web.EventContext) []*presets.SQLCondition {
		return []*presets.SQLCondition{{Query: "kind = ?", Args: []any{"a"}}}
	}).SelectionHiddenFilters("kind")

	fire := func(event string, q url.Values) string {
		r := multipartestutils.NewMultipartBuilder().
			PageURL("/admin/notes?" + q.Encode()).
			EventFunc(event).
			BuildEventFuncRequest()
		w := httptest.NewRecorder()
		p.ServeHTTP(w, r)
		if w.Code != 200 {
			t.Fatalf("%s: %d %s", event, w.Code, w.Body.String())
		}
		return strings.NewReplacer(`<`, "<", `>`, ">", `&`, "&", `\"`, `"`).Replace(w.Body.String())
	}

	// the listing as always: all the kinds, the aggregates, the filter
	body := fire(actions.ReloadList, url.Values{})
	if !strings.Contains(body, "zz-out") || !strings.Contains(body, "rvq-listing-aggregates") {
		t.Fatal("the listing: all the records and the aggregates")
	}

	sel := url.Values{
		presets.SelectedEventParamName:   {"pickNotes"},
		presets.ParamSelectMany:          {"1"},
		presets.ParamSelectedIds:         {"1,2"},
		presets.ParamSelectedDisabledIds: {"3"},
	}
	body = fire(actions.OpenListingDialogForSelection, sel)
	if strings.Contains(body, "zz-out") || !strings.Contains(body, "zz-one") || !strings.Contains(body, "zz-three") {
		t.Fatal("opened to choose: only its condition's records (kind a)")
	}
	if strings.Contains(body, "rvq-listing-aggregates") {
		t.Fatal("opened to choose: no aggregates")
	}
	if strings.Contains(body, "Kind-filter") {
		t.Fatal("opened to choose: the filter it fixes is hidden")
	}
	if !strings.Contains(body, "item-selector") || !strings.Contains(body, "disabled") {
		t.Fatal("opened to choose several: the column to mark, the chosen one fixed")
	}
	if !strings.Contains(body, "rvq-listing-selection__add") || !strings.Contains(body, "pickNotes") ||
		!strings.Contains(body, "1,2") || !strings.Contains(body, "Selected (2)") {
		t.Fatal("opened to choose several: the toggle with the count and the add giving the marked")
	}

	// only the marked ones
	sel.Set(presets.ParamSelectedView, "1")
	body = fire(actions.UpdateListingDialog, sel)
	if !strings.Contains(body, "zz-one") || !strings.Contains(body, "zz-two") || strings.Contains(body, "zz-three") {
		t.Fatal("the marked only: the ones of selected_ids")
	}
	sel.Set(presets.ParamSelectedIds, "")
	body = fire(actions.UpdateListingDialog, sel)
	if strings.Contains(body, "zz-one") || strings.Contains(body, "zz-two") {
		t.Fatal("none marked: the table of the marked is empty")
	}
}
