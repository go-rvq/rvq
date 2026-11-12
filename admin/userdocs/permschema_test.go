package userdocs

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	"golang.org/x/text/language"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type schemaPost struct {
	ID    uint
	Title string
}

type schemaComment struct {
	ID           uint
	SchemaPostID uint
	Text         string
}

func (p *schemaPost) String() string    { return p.Title }
func (c *schemaComment) String() string { return c.Text }

// schemaApp is an admin of posts — in the group site, with an action, a
// page of its record and nested comments — and a singleton settings, its
// records in a database: posts 1 "First post", 2 "Second post", the
// comment 5 "Nice" of the post 1.
func schemaApp(t *testing.T) (*presets.Builder, *i18n.Builder) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&schemaPost{}, &schemaComment{}); err != nil {
		t.Fatal(err)
	}
	db.Create(&schemaPost{ID: 1, Title: "First post"})
	db.Create(&schemaPost{ID: 2, Title: "Second post"})
	db.Create(&schemaComment{ID: 5, SchemaPostID: 1, Text: "Nice"})
	ib := i18n.New().SupportLanguages(language.English)
	ConfigureMessages(ib)
	pb := presets.New(ib).URIPrefix("/admin").DataOperator(gorm2op.DataOperator(db))
	pb.Permission(perm.New().AllowAll())
	nop := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	posts := pb.Model(&schemaPost{}, presets.ModelWithID("posts"))
	posts.Listing("Title").SearchColumns("title").Action("Import")
	posts.Editing("Title")
	d := posts.Detailing("Title")
	d.Action("Publish")
	d.PagesRegistrator().AddHttpPage(presets.HttpPage("/report").Handler(nop).AutoPerm().Methods("GET"))
	comments := presets.NewModelBuilder(pb, &schemaComment{}, presets.ModelConfig().SetId("comments"))
	posts.AddChild(comments)
	pb.Model(&treeSetting{}, presets.ModelConfig().SetSingleton(true).SetId("settings"))
	pb.MenuOrder(pb.MenuGroup("site").Add(presets.ModelItem("posts")), presets.ModelItem("settings"))
	pb.Build(http.NewServeMux())
	return pb, ib
}

// inRequest runs f in a request of the admin, in English.
func inRequest(t *testing.T, ib *i18n.Builder, f func(ctx *web.EventContext)) {
	ib.EnsureLanguage(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		f(&web.EventContext{R: r})
	})).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/admin?lang=en", nil))
}

func segments(r PermChoices) (out []string) {
	for _, c := range r.Choices {
		out = append(out, c.Segment)
	}
	return
}

func TestPermChoicesOf(t *testing.T) {
	pb, ib := schemaApp(t)
	inRequest(t, ib, func(ctx *web.EventContext) {
		for path, want := range map[string][]string{
			// the root: the scopes
			"": {"admin:"},
			// a scope, no group: its models and pages by their unique names
			"admin:": {"posts:", "settings:", "*"},
			// a group: what is in it, by the groups
			"admin:site/:": {"posts:", "*"},
			// a model: its verbs and actions, its fields, all its records
			"admin:posts:":       {"@list", "@create", "!import", "#Title:", "<*>:", "*"},
			"admin:site/:posts:": {"@list", "@create", "!import", "#Title:", "<*>:", "*"},
			// a record (any): what is asked of it, its nested model
			"admin:posts:<7>:":             {"@get", "@edit", "@delete", "@delete_with_related", "!publish", "#Title:", "/report:", "comments:", "*"},
			"admin:posts:<7<First post>>:": {"@get", "@edit", "@delete", "@delete_with_related", "!publish", "#Title:", "/report:", "comments:", "*"},
			"admin:posts:<1>:comments:":    {"@list", "@create", "#Text:", "<*>:", "*"},
			"admin:posts:<1>:/report:":     {"@get", "*"},
			// a verb: nothing follows it
			"admin:nothing:": nil,
		} {
			r, err := PermChoicesOf(ctx.R.Context(), pb, path, false, nil)
			if err != nil {
				t.Fatalf("%q: %v", path, err)
			}
			got := segments(r)
			for _, w := range want {
				if !slices.Contains(got, w) {
					t.Errorf("%q: %q missing in %q", path, w, got)
				}
			}
			if want == nil && len(got) > 0 {
				t.Errorf("%q: %q, want none", path, got)
			}
		}
		r, _ := PermChoicesOf(ctx.R.Context(), pb, "admin:posts:", false, nil)
		if !r.Records {
			t.Error("the records of the model posts not offered")
		}
		for _, c := range r.Choices {
			switch c.Segment {
			case "<*>:":
				if c.Kind != PermKindRecord || c.Title != Messages_en_US.PermInputAllRecords || !c.Children {
					t.Errorf("all the records: %+v", c)
				}
			case "@list":
				if c.Kind != PermKindVerb || c.Title != Messages_en_US.PermList || c.Children {
					t.Errorf("@list: %+v", c)
				}
			case "!import":
				if c.Kind != PermKindAction {
					t.Errorf("!import: %+v", c)
				}
			}
		}
		if r, _ := PermChoicesOf(ctx.R.Context(), pb, "admin:posts:<1>:", false, nil); r.Records {
			t.Error("the records of a record offered")
		}
		if r, _ := PermChoicesOf(ctx.R.Context(), pb, "admin:", false, nil); slices.Contains(segments(r), "site/:") {
			t.Errorf("a group without the groups: %q", segments(r))
		}
		// the groups: a resource in its group only, written out of it
		r, _ = PermChoicesOf(ctx.R.Context(), pb, "admin:", true, nil)
		if got := segments(r); !slices.Contains(got, "site/:") || slices.Contains(got, "posts:") || !slices.Contains(got, "settings:") {
			t.Errorf("the scope by the menu: %q", got)
		}
		r, _ = PermChoicesOf(ctx.R.Context(), pb, "admin:site/:", true, nil)
		for _, c := range r.Choices {
			want := map[string]string{"posts:": "admin:posts:", "*": "admin:site/:*"}[c.Segment]
			if want != "" && c.Resource != want {
				t.Errorf("%q writes %q, want %q", c.Segment, c.Resource, want)
			}
		}
		if r, _ := PermChoicesOf(ctx.R.Context(), pb, "admin:posts:<7<A>>:", false, nil); len(r.Choices) == 0 || r.Choices[0].Resource != "admin:posts:<7<A>>:"+r.Choices[0].Segment {
			t.Errorf("a record's part writes %+v", r.Choices[0])
		}
		root, _ := PermChoicesOf(ctx.R.Context(), pb, "", false, nil)
		if c := root.Choices[0]; c.Kind != PermKindScope || c.Title != Messages_en_US.PermScopeAdmin {
			t.Errorf("the scope: %+v", c)
		}
	})
}

// a stack's: only the parts of which it allows something
func TestPermChoicesOfActive(t *testing.T) {
	pb, ib := schemaApp(t)
	st := (&perm.Stack{}).Add("t", "admin:posts:{@list, <*>:@get}").Resolve()
	allows := func(c PermCheck) bool { return st.Allows(c.Resource, c.Preferred, c.Action) }
	inRequest(t, ib, func(ctx *web.EventContext) {
		for path, want := range map[string][]string{
			"admin:":       {"posts:", "*"},
			"admin:posts:": {"@list", "<*>:", "*"},
			// the glob of the stack: <*> and * span the colons, the nested
			// comments too
			"admin:posts:<*>:": {"@get", "comments:", "*"},
		} {
			r, err := PermChoicesOf(ctx.R.Context(), pb, path, false, allows)
			if err != nil {
				t.Fatal(err)
			}
			got := segments(r)
			if strings.Join(got, " ") != strings.Join(want, " ") {
				t.Errorf("%q: %q, want %q", path, got, want)
			}
		}
		r, _ := PermChoicesOf(ctx.R.Context(), pb, "admin:posts:", false, allows)
		for _, c := range r.Choices {
			if c.Segment == "*" && c.Active != PermActiveSome {
				t.Errorf("*: %+v", c)
			}
			if c.Segment == "@list" && c.Active != PermActiveAll {
				t.Errorf("@list: %+v", c)
			}
			if c.Segment == "<*>:" && c.Active != PermActiveSome {
				t.Errorf("<*>: %+v", c)
			}
		}
	})
}

func TestPermRecordsOf(t *testing.T) {
	pb, ib := schemaApp(t)
	inRequest(t, ib, func(ctx *web.EventContext) {
		for _, c := range []struct {
			path, q string
			want    []presets.PermRecord
		}{
			{"admin:posts:", "first", []presets.PermRecord{{ID: "1", Label: "First post"}}},
			{"admin:site/:posts:", "post", []presets.PermRecord{{ID: "1", Label: "First post"}, {ID: "2", Label: "Second post"}}},
			// by the id
			{"admin:posts:", "2", []presets.PermRecord{{ID: "2", Label: "Second post"}}},
			// of a nested model: of its parent's record
			{"admin:posts:<1>:comments:", "", []presets.PermRecord{{ID: "5", Label: "Nice"}}},
			// no model's listing
			{"admin:posts:<1>:", "", nil},
		} {
			got, err := PermRecordsOf(ctx, pb, c.path, c.q)
			if err != nil {
				t.Fatalf("%q %q: %v", c.path, c.q, err)
			}
			if !slices.Equal(got, c.want) {
				t.Errorf("%q %q: %v, want %v", c.path, c.q, got, c.want)
			}
		}
	})
}

func TestResolveLabels(t *testing.T) {
	pb, ib := schemaApp(t)
	inRequest(t, ib, func(ctx *web.EventContext) {
		got, errs := ResolveLabels(ctx, pb,
			"admin:{posts:{<1>:@get, <2<an old title>>:@edit, <9<gone>>:@get, <*>:@list}, site/:posts:<1>:comments:<5>:@get}")
		want := "admin:{posts:{<1<First post>>:@get, <2<Second post>>:@edit, <9<gone>>:@get, <*>:@list}, site/:posts:<1<First post>>:comments:<5<Nice>>:@get}"
		if got != want {
			t.Errorf("resolved %q\nwant %q", got, want)
		}
		if len(errs) != 1 || errs[0].ID != "9" {
			t.Errorf("errors %+v", errs)
		}
		// labels decide nothing
		if perm.StripLabels(got) != perm.StripLabels(want) {
			t.Error("the line changed")
		}
	})
}

func TestDescribeLine(t *testing.T) {
	pb, ib := schemaApp(t)
	inRequest(t, ib, func(ctx *web.EventContext) {
		line, known := DescribeLine(ctx, pb, "admin:posts:{@list, <1>:@get}, admin:posts:<9<gone>>:@edit, admin:site/:*")
		if want := "admin:{posts:{@list, <1<First post>>:@get, <9<gone>>:@edit}, site/:*}"; line != want {
			t.Errorf("line %q\nwant %q", line, want)
		}
		for path, want := range map[string]PermKnown{
			"admin:":               {Title: Messages_en_US.PermScopeAdmin, Kind: PermKindScope},
			"admin:posts:@list":    {Title: Messages_en_US.PermList, Kind: PermKindVerb},
			"admin:posts:<1>:":     {Title: "First post", Kind: PermKindRecord},
			"admin:posts:<1>:@get": {Title: Messages_en_US.PermGet, Kind: PermKindVerb},
			"admin:posts:<9>:":     {Title: "9", Kind: PermKindRecord, Error: true},
			"admin:site/:":         {Kind: PermKindGroup},
			"admin:site/:*":        {Title: Messages_en_US.PermInputEverything, Kind: PermKindAll},
		} {
			got, ok := known[path]
			if !ok {
				t.Errorf("%q not known (%v)", path, known)
				continue
			}
			if got.Kind != want.Kind || got.Error != want.Error || (want.Title != "" && got.Title != want.Title) {
				t.Errorf("%q: %+v, want %+v", path, got, want)
			}
		}
		// a line it cannot read: as it is
		if line, known := DescribeLine(ctx, pb, "admin:{posts"); line != "admin:{posts" || len(known) != 0 {
			t.Errorf("%q %v", line, known)
		}
	})
}
