package integration_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type FTNote struct {
	ID   uint
	Name string
	Kind string
}

var ftSeq int64

// A listing opens in its Default tab: with no tab picked by the query, the
// default tab's filters apply and the tab is the selected one.
func TestFilterTabDefault(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:filtertab_%d?mode=memory&cache=shared", atomic.AddInt64(&ftSeq, 1))), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&FTNote{}); err != nil {
		t.Fatal(err)
	}
	db.Create(&[]FTNote{{Name: "zz-text-note", Kind: "text"}, {Name: "zz-system-note", Kind: "system"}})

	p := presets.New(i18n.New()).URIPrefix("/admin")
	p.Permission(perm.New().AllowAll())
	p.DataOperator(gorm2op.DataOperator(db))
	l := p.Model(&FTNote{}).URIName("notes").Listing("ID", "Name")
	l.FilterDataFunc(func(ctx *web.EventContext) vx.FilterData {
		return []*vx.FilterItem{{Key: "kind", ItemType: vx.ItemTypeSelect, SQLCondition: `kind %s ?`, Invisible: true}}
	})
	l.FilterTabsFunc(func(ctx *web.EventContext) []*presets.FilterTab {
		return []*presets.FilterTab{
			{ID: "all", Label: "All"},
			{ID: "texts", Label: "Texts", Query: url.Values{"kind": {"text"}}, Default: true},
			{ID: "system", Label: "System", Query: url.Values{"kind": {"system"}}},
		}
	})

	get := func(query string) string {
		w := httptest.NewRecorder()
		p.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/admin/notes"+query, nil))
		if w.Code != http.StatusOK {
			t.Fatalf("GET %s: %d", query, w.Code)
		}
		return w.Body.String()
	}

	body := get("")
	if !strings.Contains(body, "zz-text-note") || strings.Contains(body, "zz-system-note") {
		t.Error("no tab picked: the listing should be the default tab's (texts only)")
	}
	// the tabs' model value is the default tab's index (1)
	if !strings.Contains(body, `:model-value=&#39;1&#39; class=&#39;mb-2&#39;`) {
		t.Error("no tab picked: the default tab should be the selected one")
	}

	body = get("?active_filter_tab=system&f_kind=system")
	if strings.Contains(body, "zz-text-note") || !strings.Contains(body, "zz-system-note") {
		t.Error("another tab picked: its filter, not the default's")
	}

	body = get("?active_filter_tab=all")
	if !strings.Contains(body, "zz-text-note") || !strings.Contains(body, "zz-system-note") {
		t.Error("the tab with no filter picked: everything")
	}
}
