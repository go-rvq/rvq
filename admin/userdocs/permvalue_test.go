package userdocs

import (
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
)

// rows are the items as "title|kind|prefix|edit", by their paths of titles
// ("Admin/Tree Posts/List"), the counts of the branches cut.
func rows(items []*PermTreeItem) map[string]string {
	out := map[string]string{}
	var walk func(items []*PermTreeItem, prefix string)
	walk = func(items []*PermTreeItem, prefix string) {
		for _, it := range items {
			title, _, _ := strings.Cut(it.Title, " (")
			p := prefix + title
			out[p] = it.Kind + "|" + it.Prefix + "|" + it.Edit + "|" + it.Active
			walk(it.Children, p+"/")
		}
	}
	walk(items, "")
	return out
}

func valueTreeOf(t *testing.T, value string, o PermValueOptions) (m map[string]string) {
	pb, ib := schemaApp(t)
	inRequest(t, ib, func(ctx *web.EventContext) {
		items, err := PermValueTree(ctx, pb, value, o)
		if err != nil {
			t.Fatal(err)
		}
		m = rows(items)
	})
	return
}

func TestPermValueTree(t *testing.T) {
	value := "admin:{posts:{@list, <1>:@edit}, settings:{allow *, deny @edit}, x:@y}"
	for _, groups := range []bool{false, true} {
		m := valueTreeOf(t, value, PermValueOptions{View: PermTreeByVerb, Groups: groups})
		posts := "Admin/Schema Posts"
		if groups {
			posts = "Admin/Site/Schema Posts"
		}
		for path, want := range map[string]string{
			// a line of its own: its prefix, the resource by its unique name
			posts + "/List": "verb|grant|admin:posts:@list|some", // its fields not
			// a record: a branch, what is asked of it under it
			posts + "/Record: First post":      "record||admin:posts:<1<First post>>:|some",
			posts + "/Record: First post/Edit": "verb|grant|admin:posts:<1<First post>>:@edit|some",
			"Admin/Tree Setting/Edit":          "verb|deny|admin:settings:@edit|some", // its fields allowed
			// what is under a branch: the branch's prefix
			"Admin/Tree Setting": "singleton|allow|admin:settings:|some",
			// a line of no item
			"Others/admin:x:@y": "other|grant|admin:x:@y|",
		} {
			if got, ok := m[path]; !ok || got != want {
				t.Errorf("groups %v: %s: %q, want %q", groups, path, got, want)
			}
		}
		// what the value does not touch: not shown
		if _, ok := m[posts+"/Create"]; ok {
			t.Errorf("groups %v: create shown", groups)
		}
		// the group given: with the groups, its item; without, a row of the
		// scope — what is in it allowed
		g := valueTreeOf(t, "admin:site/:*", PermValueOptions{View: PermTreeByVerb, Groups: groups})
		if got := g["Admin/Site"]; !strings.HasPrefix(got, "group|grant|admin:site/:|") {
			t.Errorf("groups %v: the group given %q (%v)", groups, got, g)
		}
		if _, ok := g["Admin/Site/Schema Posts/Create"]; ok != groups {
			t.Errorf("groups %v: what is in the group shown: %v", groups, ok)
		}
	}
	// the view by kind: the same lines, in its structure
	m := valueTreeOf(t, value, PermValueOptions{View: PermTreeByKind})
	if got := m["Admin/Schema Posts/Permissions/List"]; got != "verb|grant|admin:posts:@list|all" {
		t.Errorf("by kind: %q (%v)", got, m)
	}
}

// of a key: kept to what its user may, its denies its effect
func TestPermValueTreeBase(t *testing.T) {
	base := (&perm.Stack{}).Add("user", "admin:posts:*").Resolve()
	m := valueTreeOf(t, "deny admin:posts:@create", PermValueOptions{View: PermTreeByVerb,
		Base: func(c PermCheck) bool { return base.Allows(c.Resource, c.Preferred, c.Action) },
		Effect: func(lines []string) func(PermCheck) bool {
			st := (&perm.Stack{}).Add("key", append([]string{"*"}, lines...)...).Resolve()
			return func(c PermCheck) bool { return st.Allows(c.Resource, c.Preferred, c.Action) }
		}})
	// denied, its fields still allowed (the glob): some
	if got := m["Admin/Schema Posts/Create"]; got != "verb|deny|admin:posts:@create|some" {
		t.Errorf("the deny: %q", got)
	}
	if got := m["Admin/Schema Posts/List"]; got != "verb||admin:posts:@list|all" {
		t.Errorf("what the key still allows: %q", got)
	}
	if _, ok := m["Admin/Tree Setting"]; ok {
		t.Error("what the user may not, shown")
	}
}

func TestEditValue(t *testing.T) {
	for _, c := range []struct {
		value, op string
		resources []string
		prefix    string
		want      string
	}{
		{"admin:posts:@list", PermEditSet, []string{"admin:posts:@list"}, PermPrefixDeny, "deny admin:posts:@list"},
		{"admin:posts:{@list, @get}", PermEditSet, []string{"admin:posts:@list"}, "", "admin:posts:@get"},
		{"admin:posts:@list", PermEditSet, []string{"admin:pages:*"}, PermPrefixDeny, "admin:{deny pages:*, posts:@list}"},
		// what is under a branch: the branch alone out
		{"admin:{posts:@list, pages:}", PermEditSet, []string{"admin:pages:*"}, PermPrefixDeny, "admin:{deny pages:*, posts:@list}"},
		{"admin:{posts:{@list, <1<A>>:@get}, pages:*}", PermEditRemove, []string{"admin:posts:"}, "", "admin:pages:*"},
		// added in a group: out of it
		{"admin:pages:*", PermEditAdd, []string{"admin:site/:posts:", "admin:site/:*"}, "", "admin:{pages:*, posts:, site/:*}"},
		{"", PermEditAdd, []string{"admin:"}, PermPrefixDeny, "deny admin:"},
		// a placeholder a line is under: out
		{"deny admin:", PermEditAdd, []string{"admin:posts:#Body:*"}, PermPrefixDeny, "deny admin:posts:#Body:*"},
		{"admin:{deny builds:*, deny posts:#Body:}", PermEditSet, []string{"admin:posts:#Body:*"}, PermPrefixDeny,
			"deny admin:{builds:*, posts:#Body:*}"},
		// a text merged: its lines after, one of the same resource in the
		// place of the value's
		{"admin:{posts:@list, pages:*}", PermEditMerge, []string{"deny admin:pages:*", "admin:content/:users:*"}, "",
			"admin:{deny pages:*, posts:@list, users:*}"},
		// the value stored inside a group: out of it on any edit
		{"admin:site/:posts:*", PermEditAdd, []string{"admin:pages:@list"}, "", "admin:{pages:@list, posts:*}"},
	} {
		got, err := EditValue(c.value, c.op, c.resources, c.prefix)
		if err != nil {
			t.Fatalf("%+v: %v", c, err)
		}
		if got != c.want {
			t.Errorf("%s %q of %q: %q, want %q", c.op, c.resources, c.value, got, c.want)
		}
	}
}

// only what is defined: the lines' items and what holds them
func TestPermValueTreeOnly(t *testing.T) {
	m := valueTreeOf(t, "admin:posts:*", PermValueOptions{View: PermTreeByVerb, Only: true})
	if got := m["Admin/Schema Posts"]; !strings.HasPrefix(got, "model|grant|admin:posts:|") {
		t.Errorf("the line's item: %q", got)
	}
	if _, ok := m["Admin/Schema Posts/List"]; ok {
		t.Errorf("what it allows, not defined, shown: %v", m)
	}
	if _, ok := valueTreeOf(t, "admin:posts:*", PermValueOptions{View: PermTreeByVerb})["Admin/Schema Posts/List"]; !ok {
		t.Error("all it allows not shown")
	}
}

func TestMergeText(t *testing.T) {
	// any form: indented, by lines, ": {"
	got, _, err := mergeText("admin:posts:@list", "admin: {\n  content/:posts: {\n    @get\n  }\n  deny pages:*\n}", nil)
	if err != nil {
		t.Fatal(err)
	}
	if want := "admin:{deny pages:*, posts:{@list, @get}}"; got != want {
		t.Errorf("%q, want %q", got, want)
	}
	if _, _, err := mergeText("admin:posts:@list", "admin:{posts", nil); err == nil {
		t.Error("a text it cannot read taken")
	}
	// of a key: denies only
	if _, _, err := mergeText("", "admin:posts:*", []string{PermPrefixDeny}); err == nil {
		t.Error("a grant in a key taken")
	}
	if got, _, err := mergeText("", "deny admin:posts:*", []string{PermPrefixDeny}); err != nil || got != "deny admin:posts:*" {
		t.Errorf("%q %v", got, err)
	}
	// an unset: the lines it covers out — of a key too —; one of nothing, said
	got, idle, err := mergeText("deny admin:{posts:{@edit, @delete}, pages:*}", "unset admin:posts:*\nunset admin:x:*",
		[]string{PermPrefixDeny})
	if err != nil || got != "deny admin:pages:*" || len(idle) != 1 || idle[0] != "unset admin:x:*" {
		t.Errorf("%q %q %v", got, idle, err)
	}
}

// the help: a ? in the slot at the end of the bar of the field
func TestPermissionInputHelp(t *testing.T) {
	_, ib := schemaApp(t)
	inRequest(t, ib, func(ctx *web.EventContext) {
		html := h.MustString(PermissionInput(ctx, PermissionInputOptions{URL: "/admin/roles", Event: "e",
			HelpURL: "/docs/perm", HelpTitle: "Help"}), ctx.R.Context())
		for _, want := range []string{`<template #bar-append>`, `href='/docs/perm'`, `data-permissions-help`} {
			if !strings.Contains(html, want) {
				t.Errorf("%q not in %s", want, html)
			}
		}
		if html := h.MustString(PermissionInput(ctx, PermissionInputOptions{}), ctx.R.Context()); strings.Contains(html, "bar-append") {
			t.Errorf("a help of none: %s", html)
		}
	})
}
