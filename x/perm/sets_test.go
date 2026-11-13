package perm

import (
	"reflect"
	"testing"
)

// A line of a set: its resource and its name; not a set, no.
func TestSetRef(t *testing.T) {
	for _, c := range []struct {
		p, res, name string
		ok           bool
	}{
		{"admin:posts:%readers", "admin:posts:", "readers", true},
		{"admin:/site-files:%publishers", "admin:/site-files:", "publishers", true},
		{"admin:posts:<*>:%readers", "admin:posts:<*>:", "readers", true},
		{"admin:posts:@list", "", "", false},
		{"admin:posts:%Readers", "", "", false},
		{"%readers", "", "", false},
	} {
		res, name, ok := SetRef(c.p)
		if res != c.res || name != c.name || ok != c.ok {
			t.Errorf("%q: %q %q %v", c.p, res, name, ok)
		}
	}
	// given to a group of the menu: kept on the group
	if got := UngroupPattern("admin:content/:%readers"); got != "admin:content/:%readers" {
		t.Errorf("ungroup: %q", got)
	}
}

// The definition of a set: its blocks by kind, each checked.
func TestSetBlocks(t *testing.T) {
	b, err := SetBlocks("{model: {@list, <*>:@get}, singleton: @get}")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(b, map[string]string{"model": "{@list, <*>:@get}", "singleton": "@get"}) {
		t.Errorf("blocks: %#v", b)
	}
	for _, bad := range []string{"model: @get", "{}", "{model}", "{Model: @get}", "{model: @get, model: @list}", "{model: unset @get}"} {
		if _, err := SetBlocks(bad); err == nil {
			t.Errorf("%q: no error", bad)
		}
	}
	lines, err := SetBlockLines("admin:posts:", "{@list, deny <*>:#Secret:*, %more}", func(name string) []string {
		return []string{"admin:posts:!" + name}
	})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"admin:posts:@list", "deny admin:posts:<*>:#Secret:*", "admin:posts:!more"}; !reflect.DeepEqual(lines, want) {
		t.Errorf("lines: %q, want %q", lines, want)
	}
}

// A stack of a line of a set: the lines the set gives there, of it (Via), the
// line's kind unless the set's says another; a set that gives nothing, the
// line ignored — counting for nothing, replacing nothing.
func TestStackSets(t *testing.T) {
	t.Cleanup(func() { SetSetExpander(nil) })
	SetSetExpander(func(resource, name string) ([]string, bool) {
		if name != "readers" {
			return nil, false
		}
		return []string{resource + "@list", resource + "<*>:@get", "deny " + resource + "<*>:#Secret:*"}, true
	})
	s := (&Stack{}).Add("role: R", "admin:posts:%readers", "admin:pages:%nobody").Resolve()
	var got []string
	for _, l := range s.Lines {
		got = append(got, l.Kind.Prefix()+l.Pattern+"|"+l.Via)
	}
	want := []string{"admin:posts:@list|readers", "admin:posts:<*>:@get|readers", "deny admin:posts:<*>:#Secret:*|readers",
		"admin:pages:%nobody|nobody"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("lines: %q, want %q", got, want)
	}
	if !s.Lines[3].Ignored || s.Lines[3].Active() {
		t.Error("a set of nothing there counts")
	}
	if !s.Allows("admin:posts:@list", "", "") || !s.Allows("admin:posts:<7>:@get", "", "") || s.Allows("admin:posts:<7>:@edit", "", "") {
		t.Error("what the set gives")
	}
	if s.Allows("admin:posts:<7>:#Secret:@get", "", "") {
		t.Error("the set's deny")
	}
	// a deny of a set: every line it gives denies
	d := (&Stack{}).Add("role: R", "admin:*", "deny admin:posts:%readers").Resolve()
	if d.Allows("admin:posts:@list", "", "") || !d.Allows("admin:pages:@list", "", "") {
		t.Error("a deny of a set")
	}
	// no expander: ignored
	SetSetExpander(nil)
	n := (&Stack{}).Add("role: R", "admin:posts:%readers").Resolve()
	if len(n.Lines) != 1 || !n.Lines[0].Ignored || n.Allows("admin:posts:@list", "", "") {
		t.Error("without an expander")
	}
}
