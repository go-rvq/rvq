package integration_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"

	h "github.com/go-rvq/htmlgo"
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

// A detail page's sections are edited in place (Edit, Save; for a list
// section Create, Edit, Save, Delete an element). They change the record as
// the edit form does, so they ask what it asks — updating the record, by the
// model's route, and each field's own permission (#NAME) — and, besides, the
// section's (the section $NAME: get to see it, update to edit it).

type SPArticle struct {
	ID    uint
	Title string
	Body  string
	Tags  []*SPTag `gorm:"serializer:json"`
}

type SPTag struct {
	Name string
}

var spSeq int64

// sectionPermApp is the app with one article; canUpdate says whether the
// policy lets anybody update records.
func sectionPermApp(t *testing.T, canUpdate bool) (http.Handler, *gorm.DB) {
	t.Helper()
	var extra []*perm.PolicyBuilder
	if !canUpdate {
		extra = append(extra, perm.PolicyFor(perm.Anybody).WhoAre(perm.Denied).ToDo(presets.PermUpdate).On(perm.Anything))
	}
	return spApp(t, extra...)
}

// spApp is the app with one article, under a policy that allows everything
// but what extra denies.
func spApp(t *testing.T, extra ...*perm.PolicyBuilder) (http.Handler, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:sectionperm_%d?mode=memory&cache=shared", atomic.AddInt64(&spSeq, 1))), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&SPArticle{}); err != nil {
		t.Fatal(err)
	}
	if err = db.Create(&SPArticle{ID: 1, Title: "old", Body: "b", Tags: []*SPTag{{Name: "t1"}, {Name: "t2"}, {Name: "t3"}}}).Error; err != nil {
		t.Fatal(err)
	}

	policies := append([]*perm.PolicyBuilder{
		perm.PolicyFor(perm.Anybody).WhoAre(perm.Allowed).ToDo(perm.Anything).On(perm.Anything),
	}, extra...)

	p := presets.New(i18n.New()).URIPrefix("/admin")
	p.Permission(perm.New().Policies(policies...).SubjectsFunc(func(*http.Request) []string { return []string{"editor"} }))
	p.DataOperator(gorm2op.DataOperator(db))

	mb := p.Model(&SPArticle{}).URIName("articles")
	mb.Listing("ID", "Title")
	dp := mb.Detailing("Main", "Tags")
	dp.Section("Main").Editing("Title", "Body")
	dp.Section("Tags").IsList(&SPTag{}).Editing("Name").
		ElementShowComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return h.Span(field.Obj.(*SPTag).Name)
		}).
		ElementEditComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return h.Input("").Attr(web.VField(field.FormKey+".Name", field.Obj.(*SPTag).Name)...)
		})
	mb.Editing("Title", "Body")
	return p, db
}

func spFire(t *testing.T, h http.Handler, event string, query [][2]string, form [][2]string) *httptest.ResponseRecorder {
	t.Helper()
	b := multipartestutils.NewMultipartBuilder().PageURL("/admin/articles").EventFunc(event).Query(presets.ParamID, "1")
	for _, q := range query {
		b = b.Query(q[0], q[1])
	}
	for _, f := range form {
		b = b.AddField(f[0], f[1])
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, b.BuildEventFuncRequest())
	return w
}

func spArticle(t *testing.T, db *gorm.DB) SPArticle {
	t.Helper()
	var a SPArticle
	if err := db.First(&a, 1).Error; err != nil {
		t.Fatal(err)
	}
	return a
}

// sectionEvents are the in-place section events, each with what it leaves
// when it runs (checked on an editor who may update) — so the denied case
// proves the refusal, not a request that never did anything.
var sectionEvents = []struct {
	name, event string
	query, form [][2]string
	// did reports whether the event took effect, from the record and the
	// response body
	did func(a SPArticle, body string) bool
}{
	{
		name: "edit a section", event: actions.DoEditDetailingField,
		query: [][2]string{{presets.SectionFieldName, "Main"}},
		did:   func(_ SPArticle, body string) bool { return strings.Contains(body, `Main.Title`) },
	},
	{
		name: "save a section", event: actions.DoSaveDetailingField,
		query: [][2]string{{presets.SectionFieldName, "Main"}},
		form:  [][2]string{{"Main.Title", "new"}, {"Main.Body", "b"}},
		did:   func(a SPArticle, _ string) bool { return a.Title == "new" },
	},
	{
		name: "edit a list element", event: actions.DoEditDetailingListField,
		query: [][2]string{{presets.SectionFieldName, "Tags"}, {"detailListFieldEditBtn_Tags", "0"}},
		did:   func(_ SPArticle, body string) bool { return strings.Contains(body, `Tags[0].Name`) },
	},
	{
		// the third element: its form key is Tags[2] (it was Tags[10], the
		// index written in binary, and the form wrote element 10)
		name: "save the third list element", event: actions.DoSaveDetailingListField,
		query: [][2]string{{presets.SectionFieldName, "Tags"}, {"detailListFieldSaveBtn_Tags", "2"}},
		form:  [][2]string{{"Tags[2].Name", "t9"}},
		did:   func(a SPArticle, _ string) bool { return spTagNames(a) == "t1,t2,t9" },
	},
	{
		name: "save a list element", event: actions.DoSaveDetailingListField,
		query: [][2]string{{presets.SectionFieldName, "Tags"}, {"detailListFieldSaveBtn_Tags", "0"}},
		form:  [][2]string{{"Tags[0].Name", "t9"}},
		did:   func(a SPArticle, _ string) bool { return a.Tags[0].Name == "t9" },
	},
	{
		name: "add a list element", event: actions.DoCreateDetailingListField,
		query: [][2]string{{presets.SectionFieldName, "Tags"}},
		did:   func(a SPArticle, _ string) bool { return len(a.Tags) == 4 },
	},
	{
		name: "delete a list element", event: actions.DoDeleteDetailingListField,
		query: [][2]string{{presets.SectionFieldName, "Tags"}, {"detailListFieldDeleteBtn_Tags", "0"}},
		did:   func(a SPArticle, _ string) bool { return len(a.Tags) == 2 },
	},
}

func TestSectionEventsAskUpdatePermission(t *testing.T) {
	for _, e := range sectionEvents {
		t.Run(e.name, func(t *testing.T) {
			h, db := sectionPermApp(t, true)
			w := spFire(t, h, e.event, e.query, e.form)
			if !e.did(spArticle(t, db), w.Body.String()) {
				t.Fatalf("allowed: the event did nothing (%d) %.300s", w.Code, w.Body.String())
			}

			h, db = sectionPermApp(t, false)
			w = spFire(t, h, e.event, e.query, e.form)
			if e.did(spArticle(t, db), w.Body.String()) {
				t.Errorf("denied: the event ran without the update permission")
			}
			if a := spArticle(t, db); a.Title != "old" || spTagNames(a) != "t1,t2,t3" {
				t.Errorf("denied: the record changed: %+v", a)
			}
		})
	}

	// the same policy refuses the edit form's save: the sections now ask
	// what it asks
	t.Run("the edit form, for comparison", func(t *testing.T) {
		h, db := sectionPermApp(t, false)
		spFire(t, h, actions.Update, nil, [][2]string{{"Title", "new"}, {"Body", "b"}})
		if got := spArticle(t, db).Title; got != "old" {
			t.Errorf("the edit form saved without the update permission: %q", got)
		}
	})
}

func spTagNames(a SPArticle) string {
	var names []string
	for _, t := range a.Tags {
		names = append(names, t.Name)
	}
	return strings.Join(names, ",")
}

// A field's write permission is the same for the edit form and for a section
// of the detail page: the field is the record's ("#Title"), whichever of them
// writes it. Denied, the field is kept and shown read only; the others save.
func TestFieldWritePermissionIsTheSameInFormAndSection(t *testing.T) {
	denyTitle := perm.PolicyFor(perm.Anybody).WhoAre(perm.Denied).ToDo(presets.PermUpdate).On("*#Title:")

	t.Run("section", func(t *testing.T) {
		h, db := spApp(t, denyTitle)
		spFire(t, h, actions.DoSaveDetailingField,
			[][2]string{{presets.SectionFieldName, "Main"}},
			[][2]string{{"Main.Title", "new"}, {"Main.Body", "nb"}})
		if a := spArticle(t, db); a.Title != "old" || a.Body != "nb" {
			t.Errorf("section save: title %q body %q, want old and nb", a.Title, a.Body)
		}

		// the section's form shows it read only, the others editable
		body := spFire(t, h, actions.DoEditDetailingField, [][2]string{{presets.SectionFieldName, "Main"}}, nil).Body.String()
		for field, disabled := range map[string]string{"Title": "true", "Body": "false"} {
			re := regexp.MustCompile(`form\[\\"Main\.` + field + `\\"\][^>]*:disabled='` + disabled + `'`)
			if !re.MatchString(body) {
				t.Errorf("the section's %s is not disabled=%s", field, disabled)
			}
		}
	})

	t.Run("edit form", func(t *testing.T) {
		h, db := spApp(t, denyTitle)
		w := spFire(t, h, actions.Edit, nil, nil)
		stamp := stampRe.FindStringSubmatch(w.Body.String())
		if stamp == nil {
			t.Fatalf("the edit form carries no %s", presets.RecordStampFormKey)
		}
		spFire(t, h, actions.Update, nil, [][2]string{{"Title", "new"}, {"Body", "nb"}, {presets.RecordStampFormKey, stamp[1]}})
		if a := spArticle(t, db); a.Title != "old" || a.Body != "nb" {
			t.Errorf("form save: title %q body %q, want old and nb", a.Title, a.Body)
		}
	})
}

// sectionHTML is the section's portal on the rendered detail page, "" when
// the page has none.
func sectionHTML(body, name string) string {
	i := strings.Index(body, "DetailFieldPortal_"+name)
	if i < 0 {
		return ""
	}
	body = body[i:]
	if j := strings.Index(body, `/go-plaid-portal`); j >= 0 {
		body = body[:j]
	}
	return body
}

var sectionEditBtn = regexp.MustCompile(`isHovering(\\u0026|&amp;|&){2}(true|false)`)

func detailPage(t *testing.T, h http.Handler) string {
	t.Helper()
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/admin/articles/1", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("detail page: %d", w.Code)
	}
	return w.Body.String()
}

// A section has a permission of its own, $NAME: seeing it (get) and editing
// it in place (update). Its fields keep theirs (#NAME): denying a section
// does not deny its fields to the edit form.
func TestSectionPermission(t *testing.T) {
	t.Run("may see and edit", func(t *testing.T) {
		h, _ := spApp(t)
		main := sectionHTML(detailPage(t, h), "Main")
		if m := sectionEditBtn.FindStringSubmatch(main); m == nil || m[2] != "true" {
			t.Errorf("the section should show its edit button: %.300s", main)
		}
	})

	t.Run("may see, not edit", func(t *testing.T) {
		h, db := spApp(t, perm.PolicyFor(perm.Anybody).WhoAre(perm.Denied).ToDo(presets.PermUpdate).On("*$Main:"))

		main := sectionHTML(detailPage(t, h), "Main")
		if !strings.Contains(main, "v-card") {
			t.Fatalf("the section is not shown: %.300s", main)
		}
		if m := sectionEditBtn.FindStringSubmatch(main); m == nil || m[2] != "false" {
			t.Errorf("the section shows its edit button: %.300s", main)
		}

		spFire(t, h, actions.DoSaveDetailingField,
			[][2]string{{presets.SectionFieldName, "Main"}},
			[][2]string{{"Main.Title", "new"}, {"Main.Body", "nb"}})
		if a := spArticle(t, db); a.Title != "old" || a.Body != "b" {
			t.Errorf("the section was saved: %+v", a)
		}

		// the fields are not the section's: the edit form writes them
		w := spFire(t, h, actions.Edit, nil, nil)
		stamp := stampRe.FindStringSubmatch(w.Body.String())
		if stamp == nil {
			t.Fatalf("the edit form carries no %s", presets.RecordStampFormKey)
		}
		spFire(t, h, actions.Update, nil, [][2]string{{"Title", "new"}, {"Body", "b"}, {presets.RecordStampFormKey, stamp[1]}})
		if got := spArticle(t, db).Title; got != "new" {
			t.Errorf("denying the section denied its field to the edit form: title %q", got)
		}
	})

	t.Run("may not see", func(t *testing.T) {
		h, _ := spApp(t, perm.PolicyFor(perm.Anybody).WhoAre(perm.Denied).ToDo(presets.PermGet).On("*$Main:"))
		body := detailPage(t, h)
		if main := sectionHTML(body, "Main"); strings.Contains(main, "v-card") {
			t.Errorf("a section that may not be seen is shown: %.300s", main)
		}
		if tags := sectionHTML(body, "Tags"); !strings.Contains(tags, "t1") {
			t.Errorf("the other section is not shown: %.300s", tags)
		}
	})

	t.Run("a section the request makes up", func(t *testing.T) {
		h, db := spApp(t)
		w := spFire(t, h, actions.DoSaveDetailingField,
			[][2]string{{presets.SectionFieldName, "Title"}},
			[][2]string{{"Title.Title", "new"}})
		if got := spArticle(t, db).Title; got != "old" {
			t.Errorf("an unknown section saved: title %q (%.200s)", got, w.Body.String())
		}
	})
}
