package integration_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

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

// What a request asks is what the dump of the permissions says: every
// permission the events of the admin ask — the listing, the forms of a new
// record and of an edit, their saves, the detail, a section saved in place,
// the actions, the bulk actions, the pages, a nested model, a singleton, the
// fields, those of a model edited in place — is a
// resource of the dump (Builder.BuildPermissions), with that permission or
// action among those asked of it. A policy written from the dump, or from the
// documentation made of it, is then the one the request is checked against.

type PAPost struct {
	ID     uint
	Title  string
	Body   string
	Config PAConfig `gorm:"serializer:json"`
}

// PAConfig is a model edited in place, in a field of a post.
type PAConfig struct {
	ID   uint
	Note string
}

type PAComment struct {
	ID     uint
	PostID uint
	Text   string
}

type PASetting struct {
	ID   uint
	Name string
}

var paSeq int64

// paAsked are the permissions asked, as perm.Builder.OnAsk tells them.
type paAsked struct {
	mu    sync.Mutex
	asked []perm.Asked
}

func (a *paAsked) add(x perm.Asked) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.asked = append(a.asked, x)
}

func (a *paAsked) take() []perm.Asked {
	a.mu.Lock()
	defer a.mu.Unlock()
	r := a.asked
	a.asked = nil
	return r
}

func paApp(t *testing.T, policies ...*perm.PolicyBuilder) (*presets.Builder, http.Handler, *paAsked) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:permasked_%d?mode=memory&cache=shared", atomic.AddInt64(&paSeq, 1))), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&PAPost{}, &PAComment{}, &PASetting{}); err != nil {
		t.Fatal(err)
	}
	db.Create(&PAPost{ID: 1, Title: "old", Body: "b"})
	db.Create(&PAComment{ID: 1, PostID: 1, Text: "c"})
	db.Create(&PASetting{ID: 1, Name: "s"})

	asked := &paAsked{}
	p := presets.New(i18n.New()).URIPrefix("/admin")
	pb := perm.New().AllowAll()
	pb.CreatePolicies(policies...)
	p.Permission(pb.OnAsk(asked.add).SubjectsFunc(func(*http.Request) []string { return []string{"editor"} }))
	p.DataOperator(gorm2op.DataOperator(db))

	nop := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("ok")) })
	done := func(string, *web.EventContext) error { return nil }

	posts := p.Model(&PAPost{}, presets.ModelWithID("posts"))
	l := posts.Listing("ID", "Title")
	l.BulkAction("Archive").UpdateFunc(func([]string, *web.EventContext, *web.EventResponse) error { return nil })
	l.PagesRegistrator().AddHttpPage(presets.HttpPage("/export").Handler(nop).AutoPerm())
	config := presets.NewModelBuilder(p, &PAConfig{}, presets.ModelConfig().SetId("config"))
	ed := posts.Editing("Title", "Body", "Config")
	ed.Field("Config").AutoNested(config, &config.Editing("Note").FieldsBuilder)
	d := posts.Detailing("Main", "Body")
	d.Section("Main").Editing("Title")
	d.Action("Publish").UpdateFunc(done)
	d.PagesRegistrator().AddHttpPage(presets.HttpPage("/report").Handler(nop).AutoPerm().Methods("GET", "POST"))
	comments := presets.NewModelBuilder(p, &PAComment{}, presets.ModelConfig().SetId("comments"))
	comments.Editing("Text")
	posts.AddChild(comments)

	p.Model(&PASetting{}, presets.ModelConfig().SetSingleton(true).SetId("settings")).Editing("Name")
	p.MenuOrder(p.MenuGroup("content").Add(presets.ModelItem("posts")), presets.ModelItem("settings"))
	mux := http.NewServeMux()
	p.Build(mux)
	return p, mux, asked
}

// paRecord is a resource with its records any: "<1>" is "<*>".
var paRecord = regexp.MustCompile(`<[^>*]*>`)

func TestAskedPermissionsAreTheDump(t *testing.T) {
	p, h, asked := paApp(t)
	posts := p.GetModel(&PAPost{})
	list := posts.Info().ListingHref()

	event := func(url, ev string, query [][2]string, form [][2]string) {
		t.Helper()
		b := multipartestutils.NewMultipartBuilder().PageURL(url).EventFunc(ev)
		for _, q := range query {
			b = b.Query(q[0], q[1])
		}
		for _, f := range form {
			b = b.AddField(f[0], f[1])
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, b.BuildEventFuncRequest())
		if w.Code != http.StatusOK {
			t.Errorf("%s %s: %d %.200s", url, ev, w.Code, w.Body.String())
		}
	}
	send := func(method, url string, status int) {
		t.Helper()
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(method, url, nil))
		if w.Code != status {
			t.Errorf("%s %s: %d, want %d %.200s", method, url, w.Code, status, w.Body.String())
		}
	}
	get := func(url string) { t.Helper(); send(http.MethodGet, url, http.StatusOK) }
	id := [][2]string{{presets.ParamID, "1"}}
	post := [][2]string{{"Title", "t"}, {"Body", "b"}, {"Config.Note", "n"}}

	type step struct {
		name string
		do   func()
		// want are permissions the step must ask, records any ("<*>")
		want []string
	}
	steps := []step{
		{"the listing", func() { get(list) }, []string{"admin:content/:posts:@list", "admin:content/:posts:#Title:@list"}},
		{"the form of a new record", func() { event(list, actions.New, nil, nil) },
			[]string{"admin:content/:posts:@create"}},
		// a new record's fields: of the model, no record yet
		{"creating", func() { event(list, actions.Create, nil, post) }, []string{
			"admin:content/:posts:@create",
			"admin:content/:posts:#Title:@create",
			"admin:content/:posts:&Config:@create",
			"admin:content/:posts:&Config:#Note:@create",
		}},
		{"the edit form", func() { event(list, actions.Edit, id, nil) }, []string{"admin:content/:posts:<*>:@edit"}},
		{"updating", func() { event(list, actions.Update, id, post) }, []string{
			"admin:content/:posts:<*>:@edit",
			"admin:content/:posts:<*>:#Body:@edit",
			"admin:content/:posts:<*>:&Config:@edit",
			"admin:content/:posts:<*>:&Config:#Note:@edit",
		}},
		{"the detail", func() { get(list + "/1") }, []string{"admin:content/:posts:<*>:@get"}},
		{"a section saved in place", func() {
			event(list+"/1", actions.DoSaveDetailingField, append(id, [2]string{presets.SectionFieldName, "Main"}), [][2]string{{"Main.Title", "s"}})
		}, []string{"admin:content/:posts:<*>:$Main:@edit", "admin:content/:posts:<*>:#Title:@get"}},
		{"an action of the detail", func() {
			event(list+"/1", actions.DoAction, append(id, [2]string{presets.ParamAction, "Publish"}), nil)
		}, []string{"admin:content/:posts:<*>:!publish"}},
		{"a bulk action", func() {
			event(list, actions.DoBulkAction, [][2]string{{presets.ParamBulkActionName, "Archive"}, {presets.ParamSelectedIds, "1"}}, nil)
		}, []string{"admin:content/:posts:!archive"}},
		{"a page of the listing", func() { get(list + "/export") }, []string{"admin:content/:posts:/export:@get"}},
		{"a page of the record", func() { get(list + "/1/report") }, []string{"admin:content/:posts:<*>:/report:@get"}},
		// a page asks the permission of the HTTP method; one it does not
		// register is not answered
		{"a page of the record, by POST", func() { send(http.MethodPost, list+"/1/report", http.StatusOK) },
			[]string{"admin:content/:posts:<*>:/report:@post"}},
		{"a page of the record, by a method it has not", func() { send(http.MethodPut, list+"/1/report", http.StatusMethodNotAllowed) }, nil},
		{"a page of the listing, by any method", func() { send(http.MethodPatch, list+"/export", http.StatusOK) },
			[]string{"admin:content/:posts:/export:@patch"}},
		{"the nested model", func() { get(list + "/1/comments") }, []string{"admin:content/:posts:<*>:comments:@list"}},
		{"the singleton", func() { get(p.GetModel(&PASetting{}).Info().ListingHref()) }, []string{"admin:settings:@get"}},
		{"deleting", func() { event(list, actions.DoDelete, id, nil) }, []string{"admin:content/:posts:<*>:@delete"}},
	}

	// the dump: each resource, by the groups and by the unique name, with
	// what is asked of it
	dump := map[string]map[string]bool{}
	for _, e := range p.BuildPermissions().Tree().Zip() {
		for _, r := range []string{e.Resource, e.Unique} {
			if r == "" {
				continue
			}
			if dump[r] == nil {
				dump[r] = map[string]bool{}
			}
			for _, a := range e.Actions {
				dump[r][a.Name] = true
			}
		}
	}

	for _, s := range steps {
		t.Run(s.name, func(t *testing.T) {
			asked.take()
			s.do()
			got := map[string]bool{}
			for _, a := range asked.take() {
				for _, full := range []string{a.Resource, a.Preferred} {
					if full == "" {
						continue
					}
					full = paRecord.ReplaceAllString(full, "<*>")
					got[full] = true
					res := strings.TrimSuffix(full, a.Action)
					acts, ok := dump[res]
					switch {
					case !ok:
						t.Errorf("asked %s: %s is not in the dump", full, res)
					case a.Action != "" && len(acts) > 0 && !acts[a.Action]:
						t.Errorf("asked %s: the dump has not %s among %v", full, a.Action, keys(acts))
					}
				}
			}
			for _, w := range s.want {
				if !got[w] {
					t.Errorf("not asked %s; asked:\n%s", w, strings.Join(keys(got), "\n"))
				}
			}
		})
	}
}

func keys(m map[string]bool) (r []string) {
	for k := range m {
		r = append(r, k)
	}
	sort.Strings(r)
	return
}

// A policy on the permission of a method of a page decides that method alone.
func TestPagePermissionByMethod(t *testing.T) {
	p, h, _ := paApp(t, perm.PolicyFor(perm.Anybody).WhoAre(perm.Denied).ToDo(perm.Anything).On("admin:posts:<*>:/report:@post"))
	report := p.GetModel(&PAPost{}).Info().ListingHref() + "/1/report"
	for method, want := range map[string]int{http.MethodGet: http.StatusOK, http.MethodPost: http.StatusForbidden} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(method, report, nil))
		if w.Code != want {
			t.Errorf("%s: %d, want %d", method, w.Code, want)
		}
	}
}
