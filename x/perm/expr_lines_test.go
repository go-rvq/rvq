package perm

import (
	"errors"
	"slices"
	"testing"
)

func TestLines(t *testing.T) {
	e, err := ParseExpr("deny admin:posts:{<7<A, b>>:@get, allow @list}, site:*")
	if err != nil {
		t.Fatal(err)
	}
	got := e.Lines()
	want := []string{"deny admin:posts:<7<A, b>>:@get", "allow admin:posts:@list", "site:*"}
	if !slices.Equal(got, want) {
		t.Errorf("%q, want %q", got, want)
	}
}

func TestJoinSplit(t *testing.T) {
	lines := []string{"admin:posts:*", "deny admin:posts:@delete", "site:*", "admin:pages:@list", ""}
	v, err := JoinLines(lines)
	if err != nil {
		t.Fatal(err)
	}
	if want := "{admin:{posts:{*, deny @delete}, pages:@list}, site:*}"; v != want {
		t.Errorf("joined %q, want %q", v, want)
	}
	if !sameDecision(t, "{"+lines[0]+", "+lines[1]+", "+lines[2]+", "+lines[3]+"}", v) {
		t.Error("the joined decides otherwise")
	}
	split, err := SplitRoot(v)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"admin:{posts:{*, deny @delete}, pages:@list}", "site:*"}; !slices.Equal(split, want) {
		t.Errorf("split %q, want %q", split, want)
	}
	if s, _ := SplitRoot("deny {admin:a, site:b}"); !slices.Equal(s, []string{"deny admin:a", "deny site:b"}) {
		t.Errorf("split of a prefix %q", s)
	}
	if v, _ := JoinLines(nil); v != "" {
		t.Errorf("no line: %q", v)
	}
	if s, _ := SplitRoot(""); s != nil {
		t.Errorf("blank: %q", s)
	}
}

func TestUngroupPattern(t *testing.T) {
	for p, want := range map[string]string{
		// a resource: out of its groups
		"admin:content/:posts:*":                  "admin:posts:*",
		"admin:site/:/report:@get":                "admin:/report:@get",
		"admin:content/:posts:<*>:comments:@list": "admin:posts:<*>:comments:@list",
		"admin:content/:sub/:posts:<7<A:b>>:@get": "admin:posts:<7<A:b>>:@get",
		"admin:content/:posts:":                   "admin:posts:",
		// a group: as it is
		"admin:content/:*":      "admin:content/:*",
		"admin:content/:sub/:":  "admin:content/:sub/:",
		"admin:content/:sub/:*": "admin:content/:sub/:*",
		"admin:content/:":       "admin:content/:",
		// a glob, a verb after the group: as it is
		"admin:*/:posts:*":        "admin:*/:posts:*",
		"admin:content/:p*:@list": "admin:content/:p*:@list",
		"admin:content/:@list":    "admin:content/:@list",
		// no group
		"admin:posts:*": "admin:posts:*",
		"admin:*":       "admin:*",
	} {
		if got := UngroupPattern(p); got != want {
			t.Errorf("%q: %q, want %q", p, got, want)
		}
	}
}

func TestDropPlaceholders(t *testing.T) {
	got := DropPlaceholders([]string{"deny admin:builds:*", "deny admin:posts:#Body:", "deny admin:posts:#Body:*",
		"deny admin:", "admin:pages:", "site:<7<A>>:", "site:<7>:@get"})
	want := []string{"deny admin:builds:*", "deny admin:posts:#Body:*", "admin:pages:", "site:<7>:@get"}
	if !slices.Equal(got, want) {
		t.Errorf("%q, want %q", got, want)
	}
	lines, err := StoreLines("{admin:{deny builds:*, posts:{deny #Body:, deny #Body:*}}, deny admin:}")
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"deny admin:{builds:*, posts:#Body:*}"}; !slices.Equal(lines, want) {
		t.Errorf("stored %q, want %q", lines, want)
	}
	if lines, _ := StoreLines("admin:content/:posts:*, deny admin:"); !slices.Equal(lines, []string{"admin:posts:*"}) {
		t.Errorf("stored %q", lines)
	}
	// a placeholder alone decides nothing: not stored
	if lines, _ := StoreLines("admin:{users:, posts:<7<A>>:, content/:, roles:@list}"); !slices.Equal(lines, []string{"admin:roles:@list"}) {
		t.Errorf("stored %q", lines)
	}
	if lines, _ := StoreLines("deny admin:"); len(lines) != 0 {
		t.Errorf("stored %q", lines)
	}
	if got := DropIncomplete([]string{"admin:users:", "admin:users:*", "deny admin:", "site:<7<a:>>:@get"}); !slices.Equal(got, []string{"admin:users:*", "site:<7<a:>>:@get"}) {
		t.Errorf("%q", got)
	}
}

func TestCleanLines(t *testing.T) {
	got, changed := CleanLines([]string{"admin:content/:posts:*", "deny admin:", "admin:{users:, roles:@list}", "site:*"})
	if want := []string{"admin:posts:*", "admin:roles:@list", "site:*"}; !changed || !slices.Equal(got, want) {
		t.Errorf("%q %v, want %q", got, changed, want)
	}
	// the trash's restore of before: the record's
	got, changed = CleanLines([]string{"admin:{posts:!restore, deny content/:pages:!restore}"})
	if want := []string{"admin:{posts:<*>:@restore, deny pages:<*>:@restore}"}; !changed || !slices.Equal(got, want) {
		t.Errorf("%q %v, want %q", got, changed, want)
	}
	// nothing to clean: as they are
	in := []string{"admin:{posts:*,pages:@list}", "admin:content/:*"}
	if got, changed := CleanLines(in); changed || !slices.Equal(got, in) {
		t.Errorf("%q %v", got, changed)
	}
}

func TestMerge(t *testing.T) {
	parse := func(s string) *ExprAlt {
		e, err := ParseExpr(s)
		if err != nil {
			t.Fatal(err)
		}
		return e
	}
	a := parse("admin:{posts:{@list, <7<A>>:@get}, pages:*}")
	got := a.Merge(parse("deny admin:pages:*"), parse("site:*, admin:users:@list"))
	if want := "{admin:{posts:{@list, <7<A>>:@get}, deny pages:*, users:@list}, site:*}"; got.String() != want {
		t.Errorf("%q, want %q", got.String(), want)
	}
	// a is not changed
	if a.String() != "admin:{posts:{@list, <7<A>>:@get}, pages:*}" {
		t.Errorf("a changed: %q", a.String())
	}
	// the later decides: its line after the others
	if got := parse("{admin:a:*, admin:b:*}").Merge(parse("deny admin:a:*")); !sameDecision(t, "{admin:b:*, deny admin:a:*}", got.String()) {
		t.Errorf("%q", got.String())
	}
	if got := (&ExprAlt{}).Merge(); got.String() != "" {
		t.Errorf("none: %q", got.String())
	}
}

func TestMergeUnset(t *testing.T) {
	parse := func(s string) *ExprAlt {
		e, err := ParseExpr(s)
		if err != nil {
			t.Fatal(err)
		}
		return e
	}
	a := parse("admin:{*, posts:{@list, deny <7<A>>:@get, *}, pages:@list}")
	got, idle := a.MergeReport(parse("unset admin:posts:*, admin:users:*, unset admin:x:*"))
	// the posts' lines out, whatever their prefix; admin:* (wider) kept
	if want := "admin:{*, pages:@list, users:*}"; got.String() != want {
		t.Errorf("%q, want %q", got.String(), want)
	}
	if !slices.Equal(idle, []string{"unset admin:x:*"}) {
		t.Errorf("idle %q", idle)
	}
	// on the lines before it only: one after kept
	if got := a.Merge(parse("{unset admin:pages:@list, admin:pages:@list}")); got.String() != "admin:{*, posts:{@list, deny <7<A>>:@get, *}, pages:@list}" {
		t.Errorf("%q", got.String())
	}
	// in a group: of the group's beginning
	if got := a.Merge(parse("admin:posts:{unset @list, unset <7>:*}")); got.String() != "admin:{*, posts:*, pages:@list}" {
		t.Errorf("%q", got.String())
	}
	// never a line of a stack
	for _, l := range []string{"unset admin:posts:*", "admin:{posts:*, unset pages:*}"} {
		if !errors.Is(CheckLine(l), ErrUnset) {
			t.Errorf("%q: %v", l, CheckLine(l))
		}
		if st := (&Stack{}).Add("t", l); len(st.Lines) != 0 {
			t.Errorf("%q in a stack: %v", l, st.Lines)
		}
	}
	// the word alone: a resource, as before
	if err := CheckLine("admin:x:!{unset, edit}"); err != nil {
		t.Error(err)
	}
}
