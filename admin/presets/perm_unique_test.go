package presets

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
)

type uniquePost struct {
	ID    uint
	Title string
}

type uniqueComment struct {
	ID   uint
	Body string
}

// uniqueApp is a model posts in the group content, with a nested model, and
// a page /import in the group tools, the policies given to the roles of the
// request.
func uniqueApp(roles []string, policies ...*perm.PolicyBuilder) (*Builder, *ModelBuilder, *ModelBuilder) {
	b := New(i18n.New()).URIPrefix("/admin")
	b.Permission(perm.New().Policies(policies...).SubjectsFunc(func(*http.Request) []string { return roles }))
	posts := b.Model(&uniquePost{}, ModelWithID("posts"))
	comments := b.Model(&uniqueComment{}, ModelWithID("comments"), ModelNotInMenu())
	posts.AddChild(comments)
	b.MenuGroup("content").Add(ModelItem("posts"))
	return b, posts, comments
}

// The resources of a model: by the chain of its groups, and by its unique
// name — its id, under the module (the groups end in "/").
func TestUniquePermResources(t *testing.T) {
	_, posts, comments := uniqueApp(nil)
	l := posts.Permissioner().ListVerifier()
	if l.Resource() != "admin:content/:posts:" || l.PreferredResource() != "admin:posts:" {
		t.Errorf("the listing: %q, %q", l.Resource(), l.PreferredResource())
	}
	if posts.UniquePermName() != "posts" {
		t.Errorf("the unique name %q", posts.UniquePermName())
	}
	// a nested model: under its parent, on both
	c := comments.Permissioner().ListVerifier()
	if c.Resource() != "admin:content/:posts:comments:" || c.PreferredResource() != "admin:posts:comments:" {
		t.Errorf("the nested model: %q, %q", c.Resource(), c.PreferredResource())
	}
	if p := (&HttpPageBuilder{path: "/import"}).UniquePermName(); p != "/import" {
		t.Errorf("a page: %q", p)
	}
}

// The policies by the unique name of a model decide before those of its
// groups; deny wins, of any role.
func TestUniquePermPrecedence(t *testing.T) {
	allow := func(sub, res string) *perm.PolicyBuilder {
		return perm.PolicyFor(sub).WhoAre(perm.Allowed).ToDo(perm.Anything).On(res)
	}
	deny := func(sub, res string) *perm.PolicyBuilder {
		return perm.PolicyFor(sub).WhoAre(perm.Denied).ToDo(perm.Anything).On(res)
	}
	const group, unique = "admin:content/:*", "admin:posts:*"
	for _, c := range []struct {
		name     string
		roles    []string
		policies []*perm.PolicyBuilder
		want     bool
	}{
		{"the group allows", []string{"a"}, []*perm.PolicyBuilder{allow("a", group)}, true},
		{"the unique name denies over the group", []string{"a"},
			[]*perm.PolicyBuilder{allow("a", group), deny("a", unique)}, false},
		{"the unique name allows over the group", []string{"a"},
			[]*perm.PolicyBuilder{deny("a", group), allow("a", unique)}, true},
		{"a role denies by the unique name, another allows the group", []string{"a", "b"},
			[]*perm.PolicyBuilder{deny("a", unique), allow("b", group)}, false},
	} {
		t.Run(c.name, func(t *testing.T) {
			_, posts, comments := uniqueApp(c.roles, c.policies...)
			r := httptest.NewRequest("GET", "/admin/content/posts", nil)
			if got := posts.Permissioner().ReqLister(r).Allowed(); got != c.want {
				t.Errorf("listing posts: %v, want %v", got, c.want)
			}
			// a nested model follows its parent's unique name
			if got := comments.Permissioner().Lister(r, ID{}).Allowed(); got != c.want {
				t.Errorf("listing the comments: %v, want %v", got, c.want)
			}
		})
	}
}

// A group is never taken for a model of its name: its part ends in "/".
func TestGroupPermPart(t *testing.T) {
	_, posts, _ := uniqueApp(nil)
	b := posts.Builder()
	b.Model(&uniqueComment{}, ModelWithID("content"))
	if GroupPermPart("content") != "content/" {
		t.Errorf("the group's part %q", GroupPermPart("content"))
	}
	if r := b.GetModelByID("content").Permissioner().ListVerifier().Resource(); r != "admin:content:" {
		t.Errorf("the model content: %q", r)
	}
	if ActionPerm("geo_ip:update") != "!geo_ip.update" {
		t.Errorf("a \":\" in an action %q", ActionPerm("geo_ip:update"))
	}
	if ActionPerm("Publish") != "!publish" || PermUpdate != "@edit" {
		t.Errorf("the verbs %q, %q", ActionPerm("Publish"), PermUpdate)
	}
}
