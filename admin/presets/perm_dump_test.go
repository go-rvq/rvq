package presets

import (
	"net/http"
	"slices"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	"golang.org/x/text/language"
)

type dumpPost struct {
	ID    uint
	Title string
	Body  string
}

type dumpComment struct {
	ID   uint
	Text string
}

type dumpSetting struct {
	ID   uint
	Name string
}

type dumpHidden struct {
	ID   uint
	Code string
}

// dumpApp has a model posts in the group content inside site, with bulk,
// listing, item and detail actions, a section, a page of its listing and one
// of its record, a nested model comments; a singleton settings at the root;
// a model out of the menu; a page of the admin in the group site.
func dumpApp(t *testing.T) *Builder {
	b := New(i18n.New().SupportLanguages(language.English)).URIPrefix("/admin")
	b.Permission(perm.New())
	posts := b.Model(&dumpPost{}, ModelWithID("posts"))
	l := posts.Listing("Title")
	l.BulkAction("Archive")
	l.Action("Import")
	l.ItemAction("Duplicate")
	l.PagesRegistrator().AddHttpPage(HttpPage("/export").Handler(nopHandler()))
	posts.Editing("Title", "Body")
	d := posts.Detailing("Title", "Body")
	d.Action("Publish")
	d.Section("Main")
	d.PagesRegistrator().AddHttpPage(HttpPage("/report").Handler(nopHandler()))
	comments := NewModelBuilder(b, &dumpComment{}, ModelConfig().SetId("comments"))
	posts.AddChild(comments)

	b.Model(&dumpSetting{}, ModelConfig().SetSingleton(true).SetId("settings"))
	b.Model(&dumpHidden{}, ModelWithID("hidden"), ModelNotInMenu())
	b.PagesRegistrator().AddHttpPage(HttpPage("/report").Handler(nopHandler()).AutoPerm().MenuGroup("site"))

	b.MenuOrder(b.MenuGroup("site").Add(b.MenuGroup("content").Add(ModelItem("posts"))), ModelItem("settings"))
	b.Build(http.NewServeMux())
	return b
}

// node is the node of the tree of resource name, nil when there is none.
func node(n *PermNode, name string) *PermNode {
	if n.Name == name {
		return n
	}
	for _, c := range n.Children {
		if f := node(c, name); f != nil {
			return f
		}
	}
	return nil
}

func actionNames(n *PermNode) (r []string) {
	for _, a := range n.Actions {
		r = append(r, a.Name)
	}
	return
}

// The dump maps, recursively, every part with a permission: the groups,
// nested; a model's listing and its record, each by its groups and by its
// unique name; the permissions and the actions asked of each; the fields,
// the sections, the pages; the nested models; a singleton's record; the
// models and the pages out of the menu.
func TestPermissionsDump(t *testing.T) {
	tree := dumpApp(t).BuildPermissions().Tree()
	want := func(name, unique string, actions ...string) *PermNode {
		t.Helper()
		n := node(tree, name)
		if n == nil {
			var all []string
			for _, e := range tree.Zip() {
				all = append(all, e.Resource)
			}
			t.Fatalf("no %s in the dump:\n%s", name, strings.Join(all, "\n"))
		}
		if n.Unique != unique {
			t.Errorf("%s: unique %q, want %q", name, n.Unique, unique)
		}
		got := actionNames(n)
		for _, a := range actions {
			if !slices.Contains(got, a) {
				t.Errorf("%s: no %s among %v", name, a, got)
			}
		}
		return n
	}

	// the groups, nested
	site := want("presets:site/:", "")
	if node(site, "presets:site/:content/:") == nil {
		t.Error("the group content is not inside site")
	}
	// the listing, the record
	want("presets:site/:content/:posts:", "presets:posts:", PermList, PermCreate, "!archive", "!import")
	record := want("presets:site/:content/:posts:<*>:", "presets:posts:<*>:",
		PermGet, PermUpdate, PermDelete, PermDeleteWithRelated, "!publish", "!duplicate")
	if slices.Contains(actionNames(record), PermList) {
		t.Error("the record is asked @list")
	}
	// the fields, the section, the pages
	want("presets:site/:content/:posts:<*>:#Title:", "presets:posts:<*>:#Title:", PermGet, PermUpdate, PermCreate)
	want("presets:site/:content/:posts:<*>:$Main:", "presets:posts:<*>:$Main:", PermGet, PermUpdate)
	want("presets:site/:content/:posts:/export:", "presets:posts:/export:")
	want("presets:site/:content/:posts:<*>:/report:", "presets:posts:<*>:/report:")
	// the nested model, under any record, recursively the same
	want("presets:site/:content/:posts:<*>:comments:", "presets:posts:<*>:comments:", PermList)
	want("presets:site/:content/:posts:<*>:comments:<*>:", "presets:posts:<*>:comments:<*>:", PermGet, PermUpdate)
	// the singleton: its record is its node
	settings := want("presets:settings:", "presets:settings:", PermGet, PermUpdate)
	if slices.Contains(actionNames(settings), PermList) || node(settings, "presets:settings:<*>:") != nil {
		t.Error("the singleton has a listing")
	}
	// out of the menu: the model, at the root; the page of the admin
	want("presets:hidden:", "presets:hidden:", PermList)
	want("presets:site/:/report:", "presets:/report:")

	// the list holds every node, with its unique name
	var uniques int
	for _, e := range tree.Zip() {
		if e.Unique != "" {
			uniques++
		}
	}
	if uniques < 10 {
		t.Errorf("%d entries with a unique name", uniques)
	}
}
