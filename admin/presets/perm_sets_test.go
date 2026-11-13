package presets

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
)

type setPost struct {
	ID    uint
	Title string
}

type setSettings struct {
	ID   uint
	Name string
}

// setsApp is a model posts in the group content, a singleton settings, and a
// page /files with sets of its own, built: the expander of the process set.
func setsApp(t *testing.T) *Builder {
	t.Helper()
	b := New(i18n.New()).URIPrefix("/admin")
	b.Permission(perm.New())
	b.Model(&setPost{}, ModelWithID("posts"))
	b.Model(&setSettings{}, ModelWithID("settings")).Singleton(true)
	b.MenuGroup("content").Add(ModelItem("posts"))
	page := HttpPage("/files").Handler(nopHandler()).AutoPerm().PermActions("!edit", "!publish").
		PermSet("authors", "{@get, !edit}", nil, nil).
		PermSet("publishers", "{%authors, !publish}", func(context.Context) string { return "Publishing" }, nil)
	b.PagesRegistrator().AddHttpPage(page)
	b.Build()
	t.Cleanup(func() { perm.SetSetExpander(nil) })
	return b
}

func sorted(s []string) []string {
	out := append([]string{}, s...)
	sort.Strings(out)
	return out
}

// A set given to a resource: by its kind (a model, a singleton), by the
// resource's own (a page's, sets in sets), by a glob of resources; nothing
// for a kind of no block, a set of no such name.
func TestPermSetExpand(t *testing.T) {
	b := setsApp(t)
	exp := func(res, name string) []string {
		l, ok := perm.ExpandSet(res, name)
		if !ok {
			return nil
		}
		return sorted(l)
	}
	if got, want := exp("admin:posts:", "readers"), sorted([]string{"admin:posts:@list", "admin:posts:<*>:@get"}); !reflect.DeepEqual(got, want) {
		t.Errorf("readers of posts: %q, want %q", got, want)
	}
	// by the groups too: the lines given to what is written
	if got := exp("admin:content/:posts:", "readers"); len(got) != 2 || got[0] != "admin:content/:posts:<*>:@get" {
		t.Errorf("readers of posts by the groups: %q", got)
	}
	if got := exp("admin:settings:", "editors"); !reflect.DeepEqual(got, sorted([]string{"admin:settings:@get", "admin:settings:@edit"})) {
		t.Errorf("editors of the singleton: %q", got)
	}
	if got := exp("admin:posts:", "editors"); len(got) != 7 {
		t.Errorf("editors of posts (authors and removers): %q", got)
	}
	if got := exp("admin:/files:", "publishers"); !reflect.DeepEqual(got, sorted([]string{"admin:/files:@get", "admin:/files:!edit", "admin:/files:!publish"})) {
		t.Errorf("the page's own publishers: %q", got)
	}
	// a page and the set of its kind, readers
	if got := exp("admin:/files:", "readers"); !reflect.DeepEqual(got, []string{"admin:/files:@get"}) {
		t.Errorf("readers of the page: %q", got)
	}
	if got := exp("admin:posts:", "removers"); len(got) != 3 {
		t.Errorf("removers: %q", got)
	}
	if exp("admin:/files:", "removers") != nil || exp("admin:posts:", "nobody") != nil || exp("admin:nothing:", "readers") != nil {
		t.Error("a set that gives nothing there gave")
	}
	// a glob: each resource by its kind
	if got := exp("admin:*:", "readers"); len(got) < 3 {
		t.Errorf("readers of all: %q", got)
	}
	// the sets offered to a node: the page's own first, titled
	var files *PermNode
	b.permissions.Tree().Walk(func(_ []*PermNode, n *PermNode) {
		if n.Unique == "admin:/files:" || n.Name == "admin:/files:" {
			files = n
		}
	})
	offers := b.PermSetsOf(context.Background(), files)
	if len(offers) < 3 || offers[0].Name != "authors" || !offers[0].Own || offers[1].Title != "Publishing" {
		t.Fatalf("offers: %+v", offers)
	}
	for _, o := range offers {
		if o.Name == "readers" && (o.Own || o.Title != "Reading") {
			t.Errorf("readers offered: %+v", o)
		}
	}
}

// The sets of the admin: after the code's, by name; one of a name of the
// code's left out.
func TestPermSetSource(t *testing.T) {
	b := setsApp(t)
	b.PermSetSource(func() []*PermSet {
		return []*PermSet{{Name: "readers", Blocks: map[string]string{"model": "*"}}, {Name: "auditors", Blocks: map[string]string{"model": "@list"}}}
	})
	if s := b.PermSet("readers"); s == nil || !s.System {
		t.Error("the code's readers replaced")
	}
	if l, ok := perm.ExpandSet("admin:posts:", "auditors"); !ok || !reflect.DeepEqual(l, []string{"admin:posts:@list"}) {
		t.Errorf("auditors: %q %v", l, ok)
	}
}
