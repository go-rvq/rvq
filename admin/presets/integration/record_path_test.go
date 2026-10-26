package integration_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	. "github.com/go-rvq/rvq/admin/presets/integration"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/multipartestutils"
)

type rpSetting struct {
	ID   uint
	Name string
}

type rpNote struct {
	ID   string `gorm:"primarykey"`
	Text string
}

// A record path is the resource of the permissions by the unique names
// (the groups left out), with the whole id of each record: it goes to text
// and back — the models and ids, root first, the page, the action — and to
// the admin's address: the detail under its parents, the listing, the
// singleton's page, a page, and an action's, its page's with
// "/action:<name>", which opens it once loaded. A path that names nothing
// there is an error.
func TestRecordPath(t *testing.T) {
	db, err := NewNestedModelsDB()
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&rpNote{}, &rpSetting{}); err != nil {
		t.Fatal(err)
	}
	p := NewNestedModelsApp(db)
	var archived string
	noop := func(*web.EventContext) (r web.PageResponse, err error) { return }
	settings := p.Model(&rpSetting{}, presets.ModelConfig().SetSingleton(true).SetId("settings"))
	settings.Editing("Name")
	settings.Detailing("Name").Action("SendTest").UpdateFunc(func(string, *web.EventContext) error { return nil })
	notes := p.Model(&rpNote{}, presets.ModelConfig().SetId("notes"))
	notes.Listing("ID", "Text").Action("ImportAll").UpdateFunc(func(string, *web.EventContext) error { return nil })
	notes.Detailing().Action("Fail").UpdateFunc(func(id string, ctx *web.EventContext) error {
		ctx.Flash = &presets.FlashMessage{Text: "it failed", Color: "error"}
		return nil
	})
	notes.Detailing("Text").Action("Archive").UpdateFunc(func(id string, ctx *web.EventContext) error {
		if ctx.R.FormValue("Reason") == "" {
			ve := &web.ValidationErrors{}
			ve.GlobalError("the reason is required")
			return ve
		}
		archived = id
		return nil
	}).ComponentFunc(func(id string, ctx *web.EventContext) (h.HTMLComponent, error) {
		return h.Input("Reason").Attr("data-archive-of", id), nil
	})
	notes.Detailing().PagesRegistrator().New(presets.HttpPage("/preview")).Layout(noop)
	notes.Listing().PagesRegistrator().New(presets.HttpPage("/import")).Layout(noop)
	files := p.PagesRegistrator().New(presets.HttpPage("/files")).Layout(noop)
	files.Action("Publish").UpdateFunc(func(string, *web.EventContext) error { return nil })
	files.Build()
	p.Build()
	mux := p // served as the applications do (hermon: mux.Handle("/", presets))

	surveys := p.GetModel(&NMSurvey{})
	places := p.GetModel(&NMPlace{})
	products := p.GetModel(&NMProduct{})
	id := func(mb *presets.ModelBuilder, v string) model.ID {
		t.Helper()
		r, err := mb.ParseRecordID(v)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}
	s, pl, pr := surveys.UniquePermName(), places.UniquePermName(), products.UniquePermName()
	productPath := "admin:" + s + ":<1>:" + pl + ":<2>:" + pr + ":<3>:"

	cases := []struct {
		name, made, path, href string
	}{
		{"listing", surveys.Info().ListingPath(nil), "admin:" + s + ":", "/admin/surveys"},
		{"record", surveys.Info().RecordPath(nil, id(surveys, "1")), "admin:" + s + ":<1>:", "/admin/surveys/1"},
		{"child of child", products.Info().RecordPath(model.IDSlice{id(surveys, "1"), id(places, "2")}, id(products, "3")),
			productPath, "/admin/surveys/1/places/2/products/3"},
		{"listing of a child", places.Info().ListingPath(model.IDSlice{id(surveys, "1")}),
			"admin:" + s + ":<1>:" + pl + ":", "/admin/surveys/1/places"},
		{"singleton", settings.Info().RecordPath(nil, model.ID{}), "admin:settings:", "/admin/settings"},
		{"an id with : and /", notes.Info().RecordPath(nil, id(notes, "a:b/c")), "admin:notes:<a:b%2Fc>:", "/admin/notes/a:b%2Fc"},
		{"page of a record", presets.PagePath(notes.Info().RecordPath(nil, id(notes, "n1")), "/preview"),
			"admin:notes:<n1>:/preview:", "/admin/notes/n1/preview"},
		{"page of the listing", presets.PagePath(notes.Info().ListingPath(nil), "/import"),
			"admin:notes:/import:", "/admin/notes/import"},
		{"page of the admin", p.PagePath("/files"), "admin:/files:", "/admin/files"},
		{"action of the singleton", presets.ActionPath(settings.Info().RecordPath(nil, model.ID{}), "SendTest"),
			"admin:settings:!send_test", "/admin/settings/action:SendTest"},
		{"action of the listing", presets.ActionPath(notes.Info().ListingPath(nil), "ImportAll"),
			"admin:notes:!import_all", "/admin/notes/action:ImportAll"},
		{"action of a page", presets.ActionPath(p.PagePath("/files"), "Publish"),
			"admin:/files:!publish", "/admin/files/action:Publish"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.made != c.path {
				t.Fatalf("made %q, want %q", c.made, c.path)
			}
			href, err := p.RecordPathHref(c.path)
			if err != nil {
				t.Fatal(err)
			}
			if href != c.href {
				t.Errorf("href %q, want %q", href, c.href)
			}
		})
	}

	target, err := p.ParseRecordPath(productPath)
	if err != nil {
		t.Fatal(err)
	}
	if n := len(target.Models); n != 3 || target.Models[2].Model != products || target.Models[0].ID.String() != "1" ||
		target.Models[2].ID.String() != "3" || target.Listing || target.Page != nil || target.Action != nil {
		t.Errorf("parsed: %+v", target)
	}

	// the href opens the record; the action's is the page of the action: its
	// form in the body, its button doing it — done, back to the record, with
	// the message; with errors, the form again in the page
	if err = db.Create(&NMSurvey{ID: 1, Name: "S1", Places: []*NMPlace{{ID: 2, Name: "P2",
		Products: []*NMProduct{{ID: 3, Name: "X"}}}}}).Error; err != nil {
		t.Fatal(err)
	}
	db.Create(&rpSetting{ID: 1, Name: "S"})
	get := func(url string) string {
		t.Helper()
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, httptest.NewRequest("GET", url, nil))
		if w.Code != http.StatusOK {
			t.Errorf("%s: %d", url, w.Code)
		}
		return w.Body.String()
	}
	get("/admin/surveys/1/places/2/products/3")
	db.Create(&rpNote{ID: "n1", Text: "N"})
	for url, want := range map[string]string{
		"/admin/settings/action:SendTest":  "SendTest",
		"/admin/notes/n1/action:Archive":   "data-archive-of",
		"/admin/notes/action:ImportAll":    "ImportAll",
		"/admin/files/action:Publish":      "Publish",
		"/admin/notes/n1/action:Archive?x": presets.ActionPagePortalName,
	} {
		page := get(url)
		if !strings.Contains(page, want) || !strings.Contains(page, presets.ActionPagePortalName) {
			t.Errorf("%s: the page of the action has no %q:\n%s", url, want, page)
		}
	}
	if strings.Contains(get("/admin/notes/n1"), presets.ActionPagePortalName) {
		t.Error("the detail alone is the page of an action")
	}
	if page := get("/admin/notes/n1/action:Nothing"); !strings.Contains(page, "404") ||
		strings.Contains(page, presets.ActionPagePortalName) {
		t.Error("the page of no action is not the one of not found")
	}

	// what its button sends
	do := func(action, reason string) string {
		t.Helper()
		b := multipartestutils.NewMultipartBuilder().PageURL("/admin/notes/n1").EventFunc(actions.DoAction).
			Query(presets.ParamAction, action).Query(presets.ParamID, "n1").Query(presets.ParamActionPage, "1")
		if reason != "" {
			b = b.AddField("Reason", reason)
		}
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, b.BuildEventFuncRequest())
		if w.Code != http.StatusOK {
			t.Fatalf("do: %d %s", w.Code, w.Body.String())
		}
		return w.Body.String()
	}
	if got := do("Fail", ""); !strings.Contains(got, `"pushState":null`) || !strings.Contains(got, "it failed") {
		t.Errorf("an error message: the page of the action stays, the message shown: %s", got)
	}
	if got := do("Archive", ""); archived != "" || !strings.Contains(got, presets.ActionPagePortalName) ||
		!strings.Contains(got, "the reason is required") || !strings.Contains(got, `"pushState":null`) {
		t.Errorf("with errors: archived %q, %s", archived, got)
	}
	if got := do("Archive", "old"); archived != "n1" || !strings.Contains(got, `"url":"/admin/notes/n1"`) {
		t.Errorf("done: archived %q, %s", archived, got)
	}

	for _, bad := range []string{"", "admin:", "admin:nothing:", "admin:nothing:<1>:", "admin:" + s + ":<1>:nothing:",
		"admin:" + s + ":<x>:", "admin:settings:<1>:", "admin:" + pl + ":", "admin:" + s, "posts:" + s + ":",
		"admin:notes:<n1>:/nothing:", "admin:/nothing:", "admin:settings:!nothing", "admin:notes:<n1>:!import_all"} {
		if _, err := p.RecordPathHref(bad); !errors.Is(err, presets.ErrRecordPath) {
			t.Errorf("%q: %v, want ErrRecordPath", bad, err)
		}
	}
}
