package copier_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/go-rvq/rvq/admin/copier"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// CopyRec is a record in a language: its key the id and the language, as a
// localized record's.
type CopyRec struct {
	ID         uint   `gorm:"primaryKey;autoIncrement:false"`
	LocaleCode string `gorm:"primaryKey"`
	Title      string
	Note       string
}

func (r *CopyRec) PrimarySlug() string { return fmt.Sprintf("%d_%s", r.ID, r.LocaleCode) }

func (r *CopyRec) PrimaryColumnValuesBySlug(slug string) map[string]string {
	id, code, _ := strings.Cut(slug, "_")
	return map[string]string{"id": id, "locale_code": code}
}

var seq int64

type asked struct {
	mu sync.Mutex
	rs []string
}

func (a *asked) add(x perm.Asked) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.rs = append(a.rs, x.Resource)
}

// count is how many of the permissions asked end with suffix.
func (a *asked) count(suffix string) (n int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, r := range a.rs {
		if strings.HasSuffix(r, suffix) {
			n++
		}
	}
	return
}

func (a *asked) reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.rs = nil
}

type app struct {
	db    *gorm.DB
	h     http.Handler
	asked *asked
	list  string
}

// copyApp is a model of CopyRec with the copier: a copy takes, with the
// switch translations, the versions of its records in the languages picked;
// a record copied in no language before gets the next id.
func copyApp(t *testing.T, policies ...*perm.PolicyBuilder) *app {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:copier_%d?mode=memory&cache=shared", atomic.AddInt64(&seq, 1))), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&CopyRec{}); err != nil {
		t.Fatal(err)
	}
	for _, r := range []*CopyRec{{1, "en", "One", "n1"}, {1, "pt", "Um", "n1pt"}, {2, "en", "Two", "n2"}, {3, "en", "Three", "n3"}} {
		db.Create(r)
	}

	a := &app{db: db, asked: &asked{}}
	p := presets.New(i18n.New()).URIPrefix("/admin")
	pb := perm.New().AllowAll()
	pb.CreatePolicies(policies...)
	p.Permission(pb.OnAsk(a.asked.add).SubjectsFunc(func(*http.Request) []string { return []string{"editor"} }))
	p.DataOperator(gorm2op.DataOperator(db))

	m := p.Model(&CopyRec{}, presets.ModelWithID("recs"))
	m.Listing("ID", "LocaleCode", "Title")
	m.Detailing("Title", "Note")
	m.Use(copier.New(db).
		Switches(&copier.Switch{Name: "translations", Choices: func(_ *web.EventContext, mids model.IDSlice) (cs []copier.Choice) {
			var codes []string
			db.Model(&CopyRec{}).Distinct("locale_code").Where("id IN ?", ids(mids)).Order("locale_code").Pluck("locale_code", &codes)
			for _, c := range codes {
				cs = append(cs, copier.Choice{Value: c, Label: c})
			}
			return
		}}).
		Expand(func(_ *web.EventContext, mids model.IDSlice, st *copier.State) (*copier.Expansion, error) {
			exp := &copier.Expansion{IDs: append(model.IDSlice{}, mids...), Added: map[string]int{}}
			if !st.On["translations"] {
				return exp, nil
			}
			have := map[string]bool{}
			for _, mid := range mids {
				have[mid.String()] = true
			}
			var recs []*CopyRec
			db.Where("id IN ? AND locale_code IN ?", ids(mids), st.Picked["translations"]).Order("id, locale_code").Find(&recs)
			for _, r := range recs {
				mid := model.ID{Fields: mids[0].Fields, Values: []any{r.ID, r.LocaleCode}}
				if !have[mid.String()] {
					have[mid.String()] = true
					exp.IDs = append(exp.IDs, mid)
					exp.Added["translations"]++
				}
			}
			return exp, nil
		}).
		Prepare(func(s *copier.Session, c *copier.Copy) error {
			to := c.To.(*CopyRec)
			if to.ID == 0 {
				var max uint
				s.DB.Model(&CopyRec{}).Select("max(id)").Scan(&max)
				to.ID = max + 1
			}
			return nil
		}))
	mux := http.NewServeMux()
	p.Build(mux)
	a.h = mux
	a.list = m.Info().ListingHref()
	return a
}

func ids(mids model.IDSlice) (r []any) {
	for _, m := range mids {
		r = append(r, m.GetValue("ID"))
	}
	return
}

func (a *app) event(t *testing.T, ev string, query, form [][2]string) string {
	t.Helper()
	b := multipartestutils.NewMultipartBuilder().PageURL(a.list).EventFunc(ev)
	for _, q := range query {
		b = b.Query(q[0], q[1])
	}
	for _, f := range form {
		b = b.AddField(f[0], f[1])
	}
	w := httptest.NewRecorder()
	a.h.ServeHTTP(w, b.BuildEventFuncRequest())
	return w.Body.String()
}

func (a *app) titles(t *testing.T) map[string]string {
	t.Helper()
	var recs []*CopyRec
	a.db.Order("id, locale_code").Find(&recs)
	m := map[string]string{}
	for _, r := range recs {
		m[strconv.Itoa(int(r.ID))+"_"+r.LocaleCode] = r.Title
	}
	return m
}

// A record copied is a new one — its title the one of a copy, its other
// fields its own —, the original as it was.
func TestCopyARecord(t *testing.T) {
	a := copyApp(t)
	a.event(t, copier.EventDo, [][2]string{{presets.ParamID, "2_en"}}, nil)
	got := a.titles(t)
	if got["2_en"] != "Two" || got["4_en"] != "Copy of Two" || len(got) != 5 {
		t.Fatalf("the records: %v", got)
	}
	var c CopyRec
	a.db.First(&c, "id = 4")
	if c.Note != "n2" {
		t.Errorf("the note of the copy: %q", c.Note)
	}
	// the title of the dialog
	a.event(t, copier.EventDo, [][2]string{{presets.ParamID, "3_en"}}, [][2]string{{"CopyTitle", "Mine"}})
	if got := a.titles(t); got["5_en"] != "Mine" {
		t.Errorf("the title chosen: %v", got)
	}
}

// A bulk copy of a record in two languages, with its translations, and of
// another: MODEL:@create asked once, <ID>:@get once by id — the versions of
// a record one —; the versions of a record copied take one new id.
func TestCopyAsksEachOnce(t *testing.T) {
	a := copyApp(t)
	a.event(t, actions.DoBulkAction, [][2]string{
		{presets.ParamBulkActionName, "Copy"}, {presets.ParamSelectedIds, "1_en,2_en"},
	}, [][2]string{{"CopySwitches.translations", "true"}, {"CopyChoices.translations", "en"}, {"CopyChoices.translations", "pt"}})
	if n := a.asked.count(":@create"); n != 1 {
		t.Errorf("@create asked %d times", n)
	}
	if n1, n2 := a.asked.count(":<1>:@get"), a.asked.count(":<2>:@get"); n1 != 1 || n2 != 1 {
		t.Errorf("@get asked: <1> %d, <2> %d", n1, n2)
	}
	got := a.titles(t)
	if got["4_en"] != "Copy of One" || got["4_pt"] != "Um" || got["5_en"] != "Copy of Two" || len(got) != 7 {
		t.Errorf("the records: %v", got)
	}
}

// A record it may not read: nothing copied, saying which.
func TestCopyDenied(t *testing.T) {
	a := copyApp(t, perm.PolicyFor(perm.Anybody).WhoAre(perm.Denied).ToDo(perm.Anything).On("*:<2>:@get"))
	body := a.event(t, actions.DoBulkAction, [][2]string{
		{presets.ParamBulkActionName, "Copy"}, {presets.ParamSelectedIds, "1_en,2_en"},
	}, nil)
	if got := a.titles(t); len(got) != 4 {
		t.Errorf("copied, denied: %v", got)
	}
	if !strings.Contains(body, "No permission to read") {
		t.Errorf("the denial not said: %.300s", body)
	}
}

// The items of the menus ask MODEL:@create once a menu, whatever its rows,
// and nothing of the records.
func TestCopyMenuAsksOnce(t *testing.T) {
	a := copyApp(t)
	get := func() {
		w := httptest.NewRecorder()
		a.h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, a.list, nil))
	}
	get()
	four := a.asked.count(":@create")
	if four == 0 {
		t.Fatal("the listing asked no @create: the menu not built")
	}
	t.Logf("@create asked %d times by a listing", four)
	a.db.Create(&CopyRec{9, "en", "Nine", ""})
	a.db.Create(&CopyRec{10, "en", "Ten", ""})
	a.asked.reset()
	get()
	if six := a.asked.count(":@create"); six != four {
		t.Errorf("@create asked %d times of 4 rows, %d of 6", four, six)
	}
	if n := a.asked.count(":@get"); n != 0 {
		// the listing asks the read of no record of its own: the copy does not
		t.Logf("@get asked %d times by the listing", n)
	}
}
