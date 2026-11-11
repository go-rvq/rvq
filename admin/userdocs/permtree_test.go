package userdocs

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	"golang.org/x/text/language"
)

// permTreeOf is the tree of the permissions of an admin of posts — in the
// group site, its column and field Title, a model edited in place (Meta), a
// section (Main, of Title), an action, a page of its record, nested
// comments — and a singleton settings, made as opts says.
func permTreeOf(t *testing.T, opts PermTreeOptions) (items []*PermTreeItem) {
	pb, ib := permApp(t)
	ib.EnsureLanguage(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		items = PermissionsTree(r.Context(), pb, opts)
	})).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/admin?lang=en", nil))
	return
}

// permApp is the admin of permTreeOf.
func permApp(t *testing.T) (*presets.Builder, *i18n.Builder) {
	ib := i18n.New().SupportLanguages(language.English)
	ConfigureMessages(ib)
	pb := presets.New(ib).URIPrefix("/admin")
	pb.Permission(perm.New())
	nop := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	posts := pb.Model(&treePost{}, presets.ModelWithID("posts"))
	posts.Listing("Title").Action("Import")
	meta := presets.NewModelBuilder(pb, &treeMeta{}, presets.ModelConfig().SetId("meta"))
	posts.Editing("Title", "Meta").Field("Meta").AutoNested(meta, &meta.Editing("Note").FieldsBuilder)
	d := posts.Detailing("Title")
	d.Action("Publish")
	d.Section("Main").Editing("Title")
	d.PagesRegistrator().AddHttpPage(presets.HttpPage("/report").Handler(nop).AutoPerm().Methods("GET", "POST"))
	posts.AddChild(presets.NewModelBuilder(pb, &treeComment{}, presets.ModelConfig().SetId("comments")))
	pb.Model(&treeSetting{}, presets.ModelConfig().SetSingleton(true).SetId("settings"))
	pb.MenuOrder(pb.MenuGroup("site").Add(presets.ModelItem("posts")), presets.ModelItem("settings"))
	pb.Build(http.NewServeMux())
	return pb, ib
}

// paths are the items of the tree by their paths of titles ("Admin/Posts/
// List"), the counts of the groups cut ("Fields (2)" is "Fields").
func paths(items []*PermTreeItem) map[string]*PermTreeItem {
	out := map[string]*PermTreeItem{}
	var walk func(items []*PermTreeItem, prefix string)
	walk = func(items []*PermTreeItem, prefix string) {
		for _, it := range items {
			title, _, _ := strings.Cut(it.Title, " (")
			p := prefix + title
			out[p] = it
			walk(it.Children, p+"/")
		}
	}
	walk(items, "")
	return out
}

// By verb, without the groups of the menu: each verb with what it asks, down
// to the leaves — one permission each, by the unique name alone.
func TestPermissionsTreeByVerb(t *testing.T) {
	tree := paths(permTreeOf(t, PermTreeOptions{View: PermTreeByVerb}))
	for p, res := range map[string]string{
		"Admin/Tree Posts":                                               "admin:posts:",
		"Admin/Tree Posts/List":                                          "admin:posts:@list",
		"Admin/Tree Posts/List/Fields/Title":                             "admin:posts:#Title:@list",
		"Admin/Tree Posts/List/Actions/Import":                           "admin:posts:!import",
		"Admin/Tree Posts/Create/Fields/Title":                           "admin:posts:#Title:@create",
		"Admin/Tree Posts/Create/Models edited in place/Meta":            "admin:posts:&Meta:@create",
		"Admin/Tree Posts/See a record":                                  "admin:posts:<*>:@get",
		"Admin/Tree Posts/See a record/Fields/Title":                     "admin:posts:<*>:#Title:@get",
		"Admin/Tree Posts/See a record/Sections/Main":                    "admin:posts:<*>:$Main:@get",
		"Admin/Tree Posts/See a record/Sections/Main/Fields/Title":       "admin:posts:<*>:#Title:@get",
		"Admin/Tree Posts/See a record/Actions/Publish":                  "admin:posts:<*>:!publish",
		"Admin/Tree Posts/See a record/Pages/Permission Report":          "admin:posts:<*>:/report:",
		"Admin/Tree Posts/See a record/Pages/Permission Report/Open":     "admin:posts:<*>:/report:@get",
		"Admin/Tree Posts/See a record/Nested models/Tree Comments/List": "admin:posts:<*>:comments:@list",
		"Admin/Tree Posts/Edit/Fields/Title":                             "admin:posts:<*>:#Title:@edit",
		"Admin/Tree Posts/Edit/Models edited in place/Meta/Fields/Note":  "admin:posts:<*>:&Meta:#Note:@edit",
		"Admin/Tree Posts/Edit/Sections/Main/Fields/Title":               "admin:posts:<*>:#Title:@edit",
		"Admin/Tree Posts/Delete":                                        "admin:posts:<*>:@delete",
		// the singleton: its record's verbs
		"Admin/Tree Setting/See a record/Fields/Name": "admin:settings:#Name:@get",
	} {
		it := tree[p]
		if it == nil {
			var all []string
			for k := range tree {
				all = append(all, k)
			}
			t.Fatalf("no %s in the tree; it has:\n%s", p, strings.Join(all, "\n"))
		}
		if strings.Join(it.Resources, " ") != res {
			t.Errorf("%s: %v, want %s", p, it.Resources, res)
		}
	}
	// the kinds: of the parts and of the permissions asked of them; a
	// branch of the tree only, none
	for p, kind := range map[string]string{
		"Admin":                              PermKindScope,
		"Admin/Tree Posts":                   PermKindModel,
		"Admin/Tree Setting":                 PermKindSingleton,
		"Admin/Tree Posts/Edit":              PermKindVerb,
		"Admin/Tree Posts/Edit/Fields":       "",
		"Admin/Tree Posts/Edit/Fields/Title": PermKindField,
		"Admin/Tree Posts/Create/Models edited in place/Meta":        PermKindInline,
		"Admin/Tree Posts/See a record/Sections/Main":                PermKindSection,
		"Admin/Tree Posts/See a record/Actions/Publish":              PermKindAction,
		"Admin/Tree Posts/See a record/Pages/Permission Report":      PermKindPage,
		"Admin/Tree Posts/See a record/Pages/Permission Report/Open": PermKindVerb,
	} {
		if it := tree[p]; it == nil || it.Kind != kind {
			t.Errorf("%s: kind %+v, want %q", p, it, kind)
		}
	}
	// no groups of the menu: the models at the root
	if _, ok := tree["Admin/Site"]; ok {
		t.Error("the group site is in the tree")
	}
	// a leaf asks one permission, by the groups and by the unique name
	if c := tree["Admin/Tree Posts/Edit/Fields/Title"].Checks; len(c) != 1 ||
		c[0] != (PermCheck{Resource: "admin:site/:posts:<*>:#Title:@edit", Preferred: "admin:posts:<*>:#Title:@edit", Action: "@edit"}) {
		t.Errorf("the check of a leaf: %+v", c)
	}
}

// By kind, with the groups of the menu: a model by what it holds, each part
// with its verbs; the resources by the groups and by the unique name.
func TestPermissionsTreeByKind(t *testing.T) {
	tree := paths(permTreeOf(t, PermTreeOptions{View: PermTreeByKind, Groups: true}))
	for p, res := range map[string]string{
		"Admin/Site/Tree Posts":                                              "admin:site/:posts: admin:posts:",
		"Admin/Site/Tree Posts/Permissions/Edit":                             "admin:site/:posts:<*>:@edit admin:posts:<*>:@edit",
		"Admin/Site/Tree Posts/Fields/Title/Edit":                            "admin:site/:posts:<*>:#Title:@edit admin:posts:<*>:#Title:@edit",
		"Admin/Site/Tree Posts/Fields/Title/List":                            "admin:site/:posts:#Title:@list admin:posts:#Title:@list",
		"Admin/Site/Tree Posts/Sections/Main/Fields/Title/See a record":      "admin:site/:posts:<*>:#Title:@get admin:posts:<*>:#Title:@get",
		"Admin/Site/Tree Posts/Models edited in place/Meta/Fields/Note/Edit": "admin:site/:posts:<*>:&Meta:#Note:@edit admin:posts:<*>:&Meta:#Note:@edit",
	} {
		it := tree[p]
		if it == nil {
			t.Fatalf("no %s in the tree", p)
		}
		if strings.Join(it.Resources, " ") != res {
			t.Errorf("%s: %v, want %s", p, it.Resources, res)
		}
	}
}

// Active: a leaf by its check, a branch by its children — all, some, none.
func TestPermissionsTreeActive(t *testing.T) {
	items := permTreeOf(t, PermTreeOptions{View: PermTreeByVerb})
	s := (&perm.Stack{}).Add("role", "admin:posts:*", "deny admin:posts:<*>:@edit").Resolve()
	SetActive(items, func(c PermCheck) bool { return s.Allows(c.Resource, c.Preferred, c.Action) })
	tree := paths(items)
	for p, want := range map[string]string{
		"Admin/Tree Posts/List": PermActiveAll,
		"Admin/Tree Posts/Edit": PermActiveSome, // denied, its fields not
		"Admin/Tree Posts":      PermActiveSome,
		"Admin/Tree Setting":    PermActiveNone,
		"Admin":                 PermActiveSome,
		"Admin/Tree Posts/See a record/Fields/Title": PermActiveAll,
	} {
		if it := tree[p]; it == nil || it.Active != want {
			t.Errorf("%s: %+v, want %s", p, it, want)
		}
	}
}

// The view has no tree in it: the table, its controls, and the event that
// answers the tree of the view shown, with the query given; the column
// Active only when asked.
func TestPermissionsTreeView(t *testing.T) {
	_, ib := permApp(t)
	var out, active string
	ib.EnsureLanguage(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		out = h.MustString(PermissionsTreeView(r.Context(), PermTreeViewOptions{URL: "/admin/docs", Event: PermissionsTreeEvent}), r.Context())
		active = h.MustString(PermissionsTreeView(r.Context(), PermTreeViewOptions{URL: "/admin/users", Event: "e",
			Query: map[string]string{"id": "7"}, Active: true}), r.Context())
	})).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/admin?lang=en", nil))
	for _, want := range []string{"vx-tree-data-table", "st.items", PermissionsTreeEvent, "/admin/docs", "data-permissions-tree-kinds",
		"data-permissions-tree-views", "data-permissions-tree-groups", "Permission"} {
		if !strings.Contains(out, want) {
			t.Errorf("no %q in the view", want)
		}
	}
	if strings.Contains(out, "admin:posts") || strings.Contains(out, `"key":"active"`) {
		t.Error("the tree, or the column Active, in the view of the documentation")
	}
	if !strings.Contains(active, `.query("id", "7")`) {
		t.Errorf("the query of the event:\n%s", active)
	}
	if !strings.Contains(active, `"key":"active"`) {
		t.Error("no column Active")
	}
}

// The event answers the tree of the view asked, with the groups or not;
// with allows, each row its Active.
func TestPermissionsTreeEventFunc(t *testing.T) {
	pb, ib := permApp(t)
	call := func(query string, allows func(*web.EventContext) (func(PermCheck) bool, error)) (items []*PermTreeItem) {
		ib.EnsureLanguage(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			res, err := PermissionsTreeEventFunc(pb, allows)(&web.EventContext{R: r})
			if err != nil {
				t.Fatal(err)
			}
			items = res.Data.([]*PermTreeItem)
		})).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/admin?lang=en&"+query, nil))
		return
	}
	if _, ok := paths(call("view=kind&groups=1", nil))["Admin/Site/Tree Posts/Fields/Title"]; !ok {
		t.Error("by kind, with the groups")
	}
	tree := paths(call("view=verb", func(*web.EventContext) (func(PermCheck) bool, error) {
		return func(c PermCheck) bool { return strings.HasPrefix(c.Preferred, "admin:posts:") }, nil
	}))
	if it := tree["Admin/Tree Posts/List/Fields/Title"]; it == nil || it.Active != PermActiveAll {
		t.Errorf("by verb, active: %+v", it)
	}
	if it := tree["Admin/Tree Setting"]; it == nil || it.Active != PermActiveNone {
		t.Errorf("the singleton not allowed: %+v", it)
	}
}

// In a document, the view is a block of HTML the markdown keeps as it is:
// none of its tags escaped (an attribute as @click is no inline HTML).
func TestPermissionsTreeInMarkdown(t *testing.T) {
	pb, ib := permApp(t)
	b := &Builder{p: pb, mb: pb.Model(&treeComment{}, presets.ModelWithID("docs"))}
	var out bytes.Buffer
	ib.EnsureLanguage(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		md := "# The tree\n\nSome text.\n" + (&renderer{b: b, ctx: &web.EventContext{R: r}}).permissionsTree() + "\nAfter it.\n"
		if err := newMarkdown(func(d string) string { return d }, codeLabels{}, false).Convert([]byte(md), &out); err != nil {
			t.Fatal(err)
		}
	})).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/admin/docs?lang=en", nil))
	html := out.String()
	if strings.Contains(html, "&lt;") || !strings.Contains(html, "<v-btn") || !strings.Contains(html, "<p>After it.</p>") {
		t.Errorf("the view in the markdown:\n%s", html)
	}
}
