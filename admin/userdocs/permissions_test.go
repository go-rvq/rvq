package userdocs

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	"golang.org/x/text/language"
)

type treePost struct {
	ID    uint
	Title string
	Meta  treeMeta
}

// treeMeta is a model edited in place, in a field of a post.
type treeMeta struct {
	ID   uint
	Note string
}

type treeComment struct {
	ID   uint
	Text string
}

type treeSetting struct {
	ID   uint
	Name string
}

// Every item of the tree of the permissions has a label and a description in
// the language of the request: the scope, the groups, a model, its record,
// fields, section, pages, a nested model, a singleton, a page of the admin.
func TestPermissionsTreeLocalized(t *testing.T) {
	ib := i18n.New().SupportLanguages(language.English, language.BrazilianPortuguese)
	ConfigureMessages(ib)
	pb := presets.New(ib).URIPrefix("/admin")
	pb.Permission(perm.New())
	nop := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	posts := pb.Model(&treePost{}, presets.ModelWithID("posts"))
	posts.Listing("Title").PagesRegistrator().AddHttpPage(presets.HttpPage("/export").Handler(nop).AutoPerm())
	meta := presets.NewModelBuilder(pb, &treeMeta{}, presets.ModelConfig().SetId("meta"))
	posts.Editing("Title", "Meta").Field("Meta").AutoNested(meta, &meta.Editing("Note").FieldsBuilder)
	d := posts.Detailing("Title")
	d.Section("Main").Editing("Title")
	d.PagesRegistrator().AddHttpPage(presets.HttpPage("/report").Handler(nop).AutoPerm())
	posts.AddChild(presets.NewModelBuilder(pb, &treeComment{}, presets.ModelConfig().SetId("comments")))
	pb.Model(&treeSetting{}, presets.ModelConfig().SetSingleton(true).SetId("settings"))
	pb.PagesRegistrator().AddHttpPage(presets.HttpPage("/tools").Handler(nop).AutoPerm().MenuGroup("site"))
	// descriptions of their own: they win over the one by the kind
	own := func(s string) func(context.Context) string { return func(context.Context) string { return s } }
	posts.DescriptionFunc(own("All the posts of the site"))
	pb.MenuGroup("site").DescriptionFunc(own("The site"))
	pb.MenuOrder(pb.MenuGroup("site").Add(presets.ModelItem("posts")), presets.ModelItem("settings"))
	pb.Build(http.NewServeMux())

	tree := func(lang string) (items []*PermTreeItem) {
		ib.EnsureLanguage(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			items = PermissionsTree(r.Context(), pb, PermTreeOptions{View: PermTreeByKind, Groups: true})
		})).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/admin/docs?lang="+lang, nil))
		SetSubtitles(items)
		return
	}

	descs, subs, titles := map[string]string{}, map[string]string{}, map[string]string{}
	var walk func(items []*PermTreeItem, msgs *Messages)
	walk = func(items []*PermTreeItem, msgs *Messages) {
		for _, it := range items {
			if it.Title == "" || it.Subtitle == "" {
				t.Errorf("%s: no label or description: %+v", it.Value, it)
			}
			desc, _, _ := strings.Cut(it.Subtitle, " · ")
			descs[it.Value], subs[it.Value], titles[it.Value] = desc, it.Subtitle, it.Title
			walk(it.Children, msgs)
		}
	}
	pt := tree("pt-BR")
	walk(pt, Messages_pt_BR)
	if pt[0].Title != Messages_pt_BR.PermScopeAdmin {
		t.Errorf("the scope %q", pt[0].Title)
	}
	for res, prefix := range map[string]string{
		"admin:site/:":       "The site",
		"admin:site/:posts:": "All the posts of the site",
		// inside the model, in groups — its record's too
		"admin:site/:posts:(verbs)":    "O que se pode fazer com",
		"admin:site/:posts:(fields)":   "Os campos de",
		"admin:site/:posts:(sections)": "As seções do detalhe de",
		"admin:site/:posts:(pages)":    "As páginas de",
		"admin:site/:posts:(models)":   "Os modelos dentro de",
		"admin:site/:posts:<*>:@edit":  "admin:site/:posts:<*>:@edit",
		// a model edited in place: its structure inside the field
		"admin:site/:posts:&Meta:":              "O campo",
		"admin:site/:posts:(inlines)":           "Os registros editados no lugar dentro de",
		"admin:site/:posts:&Meta:(fields)":      "Os campos de",
		"admin:site/:posts:&Meta:#Note:":        "O campo",
		"admin:site/:posts:#Title:":             "O campo",
		"admin:site/:posts:<*>:$Main:":          "A seção",
		"admin:site/:posts:/export:":            "A página",
		"admin:site/:posts:<*>:/report:":        "A página",
		"admin:site/:posts:<*>:comments:":       "Os registros de",
		"admin:site/:posts:<*>:comments:#Text:": "O campo",
		"admin:settings:":                       "O registro único de",
		"admin:site/:/tools:":                   "A página",
	} {
		if d, ok := descs[res]; !ok {
			t.Errorf("%s: not in the tree", res)
		} else if !strings.HasPrefix(d, prefix) {
			t.Errorf("%s: description %q, want %q…", res, d, prefix)
		}
	}
	// a page's permissions: one by HTTP method it registers, localized
	if titles["admin:site/:/tools:@get"] != "Abrir (GET)" || titles["admin:site/:/tools:@post"] != "Enviar (POST)" {
		t.Errorf("the permissions of a page by its methods: %q, %q", titles["admin:site/:/tools:@get"], titles["admin:site/:/tools:@post"])
	}
	// a leaf is one permission: its resources by the groups and by the
	// unique name, each with its verb
	if sub := subs["admin:site/:posts:<*>:@edit"]; sub != "admin:site/:posts:<*>:@edit · admin:posts:<*>:@edit" {
		t.Errorf("the leaf @edit: %q", sub)
	}
	// the section: its verbs, then the fields of the record it holds
	if _, ok := titles["admin:site/:posts:<*>:$Main:(fields)"]; !ok {
		t.Error("the fields of the section Main")
	}
	en := tree("en")
	walk(en, Messages_en_US)
	if en[0].Title != Messages_en_US.PermScopeAdmin || !strings.HasPrefix(descs["admin:settings:"], "The single record of") {
		t.Errorf("in English: %q, %q", en[0].Title, descs["admin:settings:"])
	}
}
