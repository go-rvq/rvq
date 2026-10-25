package integration_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GSNote struct {
	ID   uint
	Name string
}

type GSHidden struct {
	ID   uint
	Name string
}

type GSOff struct {
	ID   uint
	Name string
}

// The search of the top bar (presets_GlobalSearch): the records whose words
// match, by model, with the link of each; a model out of the menu is not
// looked in; under 2 letters, nothing; and the bar is in the page.
func TestGlobalSearch(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:globalsearch_%d?mode=memory&cache=shared", atomic.AddInt64(&ftSeq, 1))), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&GSNote{}, &GSHidden{}, &GSOff{}); err != nil {
		t.Fatal(err)
	}
	db.Create(&[]GSNote{{Name: "alpha zebra"}, {Name: "beta"}})
	db.Create(&[]GSHidden{{Name: "zebra hidden"}})
	db.Create(&[]GSOff{{Name: "zebra off"}})

	p := presets.New(i18n.New()).URIPrefix("/admin")
	p.Permission(perm.New().AllowAll())
	p.DataOperator(gorm2op.DataOperator(db))
	p.Model(&GSNote{}).URIName("gs-notes").Listing("ID", "Name")
	p.Model(&GSHidden{}).URIName("gs-hidden").InMenu(false).Listing("ID", "Name")
	p.Model(&GSOff{}).URIName("gs-off").GlobalSearch(false).Listing("ID", "Name")

	search := func(q string) (groups []struct {
		Title string `json:"title"`
		Items []struct {
			Title string `json:"title"`
			Href  string `json:"href"`
		} `json:"items"`
	}) {
		t.Helper()
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/admin/gs-notes?__execute_event__="+presets.GlobalSearchEvent+"&q="+q, nil)
		p.ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Fatalf("search %q: %d", q, w.Code)
		}
		var resp struct {
			Data json.RawMessage `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("search %q: %v %s", q, err, w.Body.String())
		}
		if err := json.Unmarshal(resp.Data, &groups); err != nil {
			t.Fatalf("search %q, its data: %v %s", q, err, resp.Data)
		}
		return
	}

	groups := search("zebra")
	if len(groups) != 1 || len(groups[0].Items) != 1 {
		t.Fatalf("zebra: one model (out of the menu, GlobalSearch(false): not looked in), one record: %+v", groups)
	}
	if it := groups[0].Items[0]; !strings.HasPrefix(it.Href, "/admin/gs-notes/") || it.Title == "" {
		t.Errorf("the record: %+v", it)
	}
	if groups := search("z"); len(groups) != 0 {
		t.Errorf("one letter: nothing, got %+v", groups)
	}
	if groups := search("nothing-here"); len(groups) != 0 {
		t.Errorf("no match: nothing, got %+v", groups)
	}

	// the page: the bar of the search, Ctrl K, the menu's items for it
	w := httptest.NewRecorder()
	p.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/admin/gs-notes", nil))
	body := w.Body.String()
	for _, want := range []string{"rvq-global-search__field", "vars.gs.show()", "/admin/gs-notes"} {
		if !strings.Contains(body, want) {
			t.Errorf("the page has no %q", want)
		}
	}
}
