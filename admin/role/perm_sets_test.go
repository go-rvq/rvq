package role

import (
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/go-rvq/rvq/x/perm"
)

// A set of the admin: saved, what it gives comes at the next stack; deleted,
// its lines ignored; a name of the code's, a kind of no resource, a
// definition that cannot be read: not saved.
func TestPermSetsOfTheAdmin(t *testing.T) {
	db, _, h := systemApp(t, nil)
	t.Cleanup(func() { perm.SetSetExpander(nil) })
	post := func(name, def string) {
		mb := multipartestutils.NewMultipartBuilder().PageURL("/admin/perm-sets").EventFunc(actions.Create).
			AddField("Name", name).AddField("Definition", def).
			AddField("Titles.__sent", "1").AddField("Titles.en", "Auditing")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, mb.BuildEventFuncRequest())
	}
	post("%auditors", "{model: {@list, <*>:@get}}")
	var s PermSet
	if err := db.First(&s, "name = ?", "auditors").Error; err != nil {
		t.Fatal("the set not saved: ", err)
	}
	if s.Titles["en"] != "Auditing" {
		t.Errorf("titles: %v", s.Titles)
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
	db.Model(&PermSet{}).Count(&n)
	if n != 1 {
		t.Errorf("sets saved: %d, want 1", n)
	}

	db.Delete(&s)
	Changed(nil)
	if _, ok := perm.ExpandSet("admin:roles:", "auditors"); ok {
		t.Error("a deleted set gives")
	}
	if st := (&perm.Stack{}).Add("role: R", "admin:roles:%auditors").Resolve(); len(st.Lines) != 1 || !st.Lines[0].Ignored {
		t.Error("the line of a deleted set not ignored")
	}
}
