package perm_test

import (
	"net/http/httptest"
	"testing"

	"github.com/go-rvq/rvq/x/perm"
)

// A request restricted (an access key's) is allowed only what its user's
// subjects allow AND the restriction allows: the key's policies never give
// more; a deny of the user stays; the unique name decides first on both.
func TestRestriction(t *testing.T) {
	const (
		group  = "presets:content/:posts:"
		unique = "presets:posts:"
		key    = "key:k1"
	)
	allow := func(sub, res string) *perm.PolicyBuilder {
		return perm.PolicyFor(sub).WhoAre(perm.Allowed).ToDo(perm.Anything).On(res)
	}
	deny := func(sub, res string) *perm.PolicyBuilder {
		return perm.PolicyFor(sub).WhoAre(perm.Denied).ToDo(perm.Anything).On(res)
	}
	for _, c := range []struct {
		name     string
		policies []*perm.PolicyBuilder
		want     bool
	}{
		{"the user may, the key allows", []*perm.PolicyBuilder{allow("a", group+"*"), allow(key, group+"*")}, true},
		{"the user may, the key says nothing", []*perm.PolicyBuilder{allow("a", group+"*")}, false},
		{"the user may not, the key allows", []*perm.PolicyBuilder{allow(key, "presets:*")}, false},
		{"the user denied, the key allows", []*perm.PolicyBuilder{allow("a", "presets:*"), deny("a", unique+"*"),
			allow(key, "presets:*")}, false},
		{"the user denied by the group, the key allows", []*perm.PolicyBuilder{allow("a", group+"*"),
			deny("a", group+"presets:update"), allow(key, "presets:*")}, false},
		{"the key: all the user may", []*perm.PolicyBuilder{allow("a", group+"*"), allow(key, "presets:*")}, true},
		{"the key by the unique name", []*perm.PolicyBuilder{allow("a", group+"*"), allow(key, unique+"*")}, true},
		{"the key allows another action only", []*perm.PolicyBuilder{allow("a", group+"*"), allow(key, group+"*:presets:list")}, false},
		{"the key of another", []*perm.PolicyBuilder{allow("a", group+"*"), allow("key:other", "presets:*")}, false},
	} {
		t.Run(c.name, func(t *testing.T) {
			p := perm.New().Policies(c.policies...).SubjectsFunc(sf("a"))
			v := perm.NewVerifier("presets", p).Spawn()
			base := v.ResourceParts()
			v.On("content/", "posts")
			v.Prefer(append(base, "posts")...)
			r := httptest.NewRequest("GET", "/", nil)
			r = r.WithContext(perm.WithRestriction(r.Context(), key))
			v = v.Do("presets:update").WithReq(r)
			if got := v.Allowed(); got != c.want {
				t.Errorf("allowed %v, want %v", got, c.want)
			}
		})
	}

	// with no restriction: the user's alone
	p := perm.New().Policies(allow("a", group+"*")).SubjectsFunc(sf("a"))
	v := perm.NewVerifier("presets", p).Spawn().On("content/", "posts").Do("presets:update").
		WithReq(httptest.NewRequest("GET", "/", nil))
	if !v.Allowed() {
		t.Error("no restriction: denied")
	}
	if perm.RestrictionOf(httptest.NewRequest("GET", "/", nil).Context()) != nil {
		t.Error("a restriction out of nowhere")
	}
}
