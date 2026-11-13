package role

import (
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/go-rvq/rvq/x/login"
	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// A set of the admin: saved in a language, what it gives comes at the next
// stack; its name and definition, the set's in every language; deleted, its
// lines ignored; a name of the code's, a kind of no resource, a definition
// that cannot be read: not saved.
func TestPermSetsOfTheAdmin(t *testing.T) {
	db, b, h := systemApp(t, nil)
	t.Cleanup(func() { perm.SetSetExpander(nil) })
	post := func(name, def string) {
		mb := multipartestutils.NewMultipartBuilder().PageURL("/admin/perm-sets").EventFunc(actions.Create).
			AddField("Name", name).AddField("Definition", def).AddField("Title", "Auditing")
		h.ServeHTTP(httptest.NewRecorder(), mb.BuildEventFuncRequest())
	}
	post("%auditors", "{model: {@list, <*>:@get}}")
	var s PermSet
	if err := db.First(&s, "name = ?", "auditors").Error; err != nil {
		t.Fatal("the set not saved: ", err)
	}
	if s.Title != "Auditing" || s.LocaleCode == "" || s.System {
		t.Errorf("the set: %+v", s)
	}
	if l, ok := perm.ExpandSet("admin:roles:", "auditors"); !ok || !reflect.DeepEqual(l, []string{"admin:roles:@list", "admin:roles:<*>:@get"}) {
		t.Errorf("auditors of the roles: %q %v", l, ok)
	}
	st := (&perm.Stack{}).Add("role: R", "admin:roles:%auditors").Resolve()
	if !st.Allows("admin:roles:@list", "", "") || st.Allows("admin:roles:<1>:@edit", "", "") {
		t.Error("what the set gives")
	}

	for _, bad := range [][2]string{{"readers", "{model: @list}"}, {"auditors", "{model: @list}"}, {"Bad", "{model: @list}"},
		{"kinds", "{models: @list}"}, {"defs", "model: @list"}} {
		post(bad[0], bad[1])
	}
	var n int64
	db.Model(&PermSet{}).Where("NOT system").Count(&n)
	if n != 1 {
		t.Errorf("sets saved: %d, want 1", n)
	}

	// in another language: of its title translated, the same set
	if got, err := b.LocalizePermSets(s.LocaleCode, "pt-BR", func(texts []string) ([]string, error) {
		out := make([]string, len(texts))
		for i, t := range texts {
			out[i] = strings.ToUpper(t)
		}
		return out, nil
	}); err != nil || got == 0 {
		t.Fatalf("localized: %d %v", got, err)
	}
	var pt PermSet
	if err := db.First(&pt, "id = ? AND locale_code = ?", s.ID, "pt-BR").Error; err != nil || pt.Title != "AUDITING" || pt.Definition != s.Definition {
		t.Errorf("the set in pt-BR: %+v %v", pt, err)
	}

	db.Where("id = ?", s.ID).Delete(&PermSet{})
	Changed(nil)
	if _, ok := perm.ExpandSet("admin:roles:", "auditors"); ok {
		t.Error("a deleted set gives")
	}
	if st := (&perm.Stack{}).Add("role: R", "admin:roles:%auditors").Resolve(); len(st.Lines) != 1 || !st.Lines[0].Ignored {
		t.Error("the line of a deleted set not ignored")
	}
}

// The sets of the code: registered in each language of the admin, of the
// system, of the code's definition; their name and definition not changed by
// a form, their title is; not deleted; the title translated is the one
// offered.
func TestPermSetsOfTheSystem(t *testing.T) {
	db, b, h := systemApp(t, func(db *gorm.DB) {
		db.Exec(`CREATE TABLE users (id uuid PRIMARY KEY, name text)`)
	})
	ana := uuid.New()
	db.Exec(`INSERT INTO users (id, name) VALUES (?, ?)`, ana, "Ana")
	t.Cleanup(func() { perm.SetSetExpander(nil) })
	b.pb.PermSets() // the first read registers them
	var readers []PermSet
	db.Where("name = ?", "readers").Find(&readers)
	if len(readers) == 0 || !readers[0].System || readers[0].Definition != "{model: {@list, <*>:@get}, page: @get, singleton: @get}" {
		t.Fatalf("readers: %+v", readers)
	}
	r := readers[0]
	if r.Title == "" {
		t.Errorf("readers of no title: %+v", r)
	}
	edit := multipartestutils.NewMultipartBuilder().PageURL("/admin/perm-sets").EventFunc(actions.Update).
		Query("id", r.PrimarySlug()).
		AddField("Name", "hackers").AddField("Definition", "{model: *}").AddField("Title", "Leitura geral")
	w := httptest.NewRecorder()
	req := edit.BuildEventFuncRequest()
	ctx := login.WithTrustedUser(web.WithSkipFormSign(req.Context()), &struct {
		ID   uuid.UUID
		Name string
	}{ana, "Ana"})
	h.ServeHTTP(w, req.WithContext(ctx))
	var got PermSet
	db.First(&got, "id = ? AND locale_code = ?", r.ID, r.LocaleCode)
	if got.Name != "readers" || got.Definition != r.Definition || !got.System {
		t.Errorf("a set of the system changed: %+v", got)
	}
	if got.Title != "Leitura geral" || got.UpdatedByID != model.UUID(ana) {
		t.Errorf("its title not saved: %q, by %v", got.Title, got.UpdatedByID)
	}
	// who changed it: the user, by the join
	var joined PermSet
	if err := db.Joins("CreatedBy").Joins("UpdatedBy").First(&joined, `"perm_sets".id = ? AND "perm_sets".locale_code = ?`, r.ID, r.LocaleCode).Error; err != nil ||
		joined.UpdatedBy.String() != "Ana" || joined.CreatedBy != nil {
		t.Errorf("the users of the set: %+v %+v %v", joined.CreatedBy, joined.UpdatedBy, err)
	}
	if s := b.pb.PermSet("readers"); s == nil || !s.System || s.TitleOf(b.pb.I18n().LanguageContext(t.Context(), r.LocaleCode)) != "Leitura geral" {
		t.Errorf("the title offered: %+v", s)
	}
	del := multipartestutils.NewMultipartBuilder().PageURL("/admin/perm-sets").EventFunc(actions.DoDelete).
		Query("id", r.PrimarySlug())
	h.ServeHTTP(httptest.NewRecorder(), del.BuildEventFuncRequest())
	var n int64
	db.Model(&PermSet{}).Where("name = ?", "readers").Count(&n)
	if n != int64(len(readers)) {
		t.Errorf("a set of the system deleted: %d", n)
	}
}
