package perm

import (
	"context"
	"reflect"
	"testing"

	"github.com/lib/pq"
	"github.com/ory/ladon"
)

// The resources of a policy of the database, separated by commas: the commas
// of an alternative ({a,b}) and of a record (<…>) are theirs.
func TestSplitResources(t *testing.T) {
	for in, want := range map[string][]string{
		"admin:posts:*":                         {"admin:posts:*"},
		"admin:posts:*, admin:pages:*":          {"admin:posts:*", "admin:pages:*"},
		"admin:{pages,posts}:*@{list,get}":      {"admin:{pages,posts}:*@{list,get}"},
		"admin:{a,b}:*,admin:c:<x,y>:@get,":     {"admin:{a,b}:*", "admin:c:<x,y>:@get"},
		`admin:\{a,b`:                           {`admin:\{a`, "b"},
		"admin:/site-files:<{static,img}/*>:!*": {"admin:/site-files:<{static,img}/*>:!*"},
		"":                                      nil,
	} {
		if got := SplitResources(in); !reflect.DeepEqual(got, want) {
			t.Errorf("%q: %q, want %q", in, got, want)
		}
	}
}

// A policy of the database with alternatives allows what they name.
func TestDBPolicyAlternatives(t *testing.T) {
	p := DefaultDBPolicy{Subject: "Editor", Effect: Allowed, Actions: pq.StringArray{"*"},
		Resources: pq.StringArray{"admin:{pages,posts}:*@{list,get}", "admin:menus:*"}}
	b := New()
	b.UpdateOrCreatePolicies(p.ToPolicy())
	for res, allowed := range map[string]bool{
		"admin:posts:@list":       true,
		"admin:pages:<1>:@get":    true,
		"admin:menus:<1>:@edit":   true,
		"admin:posts:<1>:@delete": false,
		"admin:users:@list":       false,
	} {
		err := b.ladon.IsAllowed(context.TODO(), &ladon.Request{Subject: "Editor", Resource: res, Action: "x"})
		if (err == nil) != allowed {
			t.Errorf("%s: %v", res, err)
		}
	}
}

// A policy of no actions is saved as one of any: its resources say all.
func TestDBPolicyNoActions(t *testing.T) {
	for _, in := range []pq.StringArray{nil, {}, {""}, {" "}} {
		p := &DefaultDBPolicy{Actions: in, Resources: pq.StringArray{"admin:x:*"}}
		_ = p.BeforeSave(nil)
		if len(p.Actions) != 1 || p.Actions[0] != "*" {
			t.Errorf("%q: %q", in, p.Actions)
		}
	}
	p := &DefaultDBPolicy{Actions: pq.StringArray{"@get", " "}}
	_ = p.BeforeSave(nil)
	if len(p.Actions) != 1 || p.Actions[0] != "@get" {
		t.Errorf("kept: %q", p.Actions)
	}
}
