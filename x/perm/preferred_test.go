package perm_test

import (
	"net/http/httptest"
	"testing"

	"github.com/go-rvq/rvq/x/perm"
)

// The resource of a model by its unique name (Prefer) decides before the
// one of its ancestors: a deny of any subject by the unique name denies, an
// allow allows; when none of its policies matches, the ancestors decide —
// as before, an allow of any subject allowing.
func TestPreferredResource(t *testing.T) {
	const (
		group  = "presets:content/:*" // the ancestors: the group content
		unique = "presets:posts:*"    // the unique name of the model posts
	)
	allow := func(sub, res string) *perm.PolicyBuilder {
		return perm.PolicyFor(sub).WhoAre(perm.Allowed).ToDo(perm.Anything).On(res)
	}
	deny := func(sub, res string) *perm.PolicyBuilder {
		return perm.PolicyFor(sub).WhoAre(perm.Denied).ToDo(perm.Anything).On(res)
	}
	for _, c := range []struct {
		name     string
		subjects []string
		policies []*perm.PolicyBuilder
		want     bool
	}{
		{"no policy", []string{"a"}, nil, false},
		{"the group allows", []string{"a"}, []*perm.PolicyBuilder{allow("a", group)}, true},
		{"the group denies", []string{"a"}, []*perm.PolicyBuilder{deny("a", group)}, false},
		{"the unique name allows, the group denies", []string{"a"},
			[]*perm.PolicyBuilder{deny("a", group), allow("a", unique)}, true},
		{"the unique name denies, the group allows", []string{"a"},
			[]*perm.PolicyBuilder{allow("a", group), deny("a", unique)}, false},
		{"the unique name allows and denies: deny wins", []string{"a"},
			[]*perm.PolicyBuilder{allow("a", unique), deny("a", unique)}, false},
		{"two roles: one denies by the unique name, the other allows the group", []string{"a", "b"},
			[]*perm.PolicyBuilder{deny("a", unique), allow("b", group)}, false},
		{"two roles: one denies by the unique name, the other allows by it", []string{"a", "b"},
			[]*perm.PolicyBuilder{deny("a", unique), allow("b", unique)}, false},
		{"two roles: one allows by the unique name, the other denies the group", []string{"a", "b"},
			[]*perm.PolicyBuilder{allow("a", unique), deny("b", group)}, true},
		{"two roles: none by the unique name, one allows the group", []string{"a", "b"},
			[]*perm.PolicyBuilder{allow("b", group)}, true},
		{"another model's unique name does not count", []string{"a"},
			[]*perm.PolicyBuilder{allow("a", group), deny("a", "presets:pages:*")}, true},
	} {
		t.Run(c.name, func(t *testing.T) {
			p := perm.New().Policies(c.policies...).SubjectsFunc(sf(c.subjects...))
			v := perm.NewVerifier("presets", p).Spawn()
			base := v.ResourceParts()
			v.On("content/", "posts")
			v.Prefer(append(base, "posts")...)
			v.On("12") // a record: on both
			v = v.Do("presets:update").WithReq(httptest.NewRequest("GET", "/", nil))
			if got := v.Allowed(); got != c.want {
				t.Errorf("allowed %v, want %v (resources %s, %s)", got, c.want, v.Resource(), v.PreferredResource())
			}
		})
	}
}

// The parts added after Prefer go to both resources; a verifier spawned
// keeps them.
func TestPreferredResourceParts(t *testing.T) {
	v := perm.NewVerifier("presets", perm.New()).Spawn()
	base := v.ResourceParts()
	v.On("site/", "seo/", "seo_config").Prefer(append(base, "seo_config")...)
	v.On("7", "#Title")
	if v.Resource() != "presets:site/:seo/:seo_config:7:#Title:" || v.PreferredResource() != "presets:seo_config:7:#Title:" {
		t.Errorf("resources %q, %q", v.Resource(), v.PreferredResource())
	}
	if s := v.Spawn(); s.PreferredResource() != v.PreferredResource() {
		t.Errorf("spawned: %q", s.PreferredResource())
	}
	if v.RemoveOn(1); v.PreferredResource() != "presets:seo_config:7:" {
		t.Errorf("removed: %q", v.PreferredResource())
	}
	if (perm.NewVerifier("presets", perm.New()).Spawn()).PreferredResource() != "" {
		t.Error("a resource by a unique name nobody set")
	}
}
