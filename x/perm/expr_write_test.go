package perm

import (
	"strings"
	"testing"
)

func TestFormatLine(t *testing.T) {
	got, err := FormatLine("admin:{posts:{@list, <7<a, {b}>>:{@get, @edit}}, deny pages:*}, site:*", WithIndent("  "))
	if err != nil {
		t.Fatal(err)
	}
	want := strings.Join([]string{
		"{",
		"  admin: {",
		"    posts: {",
		"      @list",
		"      <7<a, {b}>>: {",
		"        @get",
		"        @edit",
		"      }",
		"    }",
		"    deny pages:*",
		"  }",
		"  site:*",
		"}",
	}, "\n")
	if got != want {
		t.Errorf("indented:\n%s\nwant:\n%s", got, want)
	}
	// read back: the same line
	if again, _ := Canonical(got); again != "{admin:{posts:{@list, <7<a, {b}>>:{@get, @edit}}, deny pages:*}, site:*}" {
		t.Errorf("read back %q", again)
	}
	// no group, no labels, as written
	if got, _ := FormatLine("deny admin:posts:<7<A>>:@get", WithIndent("\t"), WithoutLabels()); got != "deny admin:posts:<7>:@get" {
		t.Errorf("%q", got)
	}
	// a space before a group, after its ":": none of it
	if got, _ := FormatLine("admin: {a:x,\n b:y}  "); got != "admin:{a:x, b:y}" {
		t.Errorf("one line: %q", got)
	}
	if got, err := FormatLine("admin:{a"); err == nil || got != "admin:{a" {
		t.Errorf("unread: %q %v", got, err)
	}
}

// the space after a ":" before a group is no part of the resource: the same
// decision
func TestSpaceBeforeGroup(t *testing.T) {
	for _, line := range []string{"admin: {posts: {@list, @get}}", "admin:\t{posts:{@list,@get}}"} {
		got, err := Canonical(line)
		if err != nil {
			t.Fatal(err)
		}
		if got != "admin:posts:{@list, @get}" {
			t.Errorf("%q: %q", line, got)
		}
	}
	// a space that is not after a ":": the resource's, as before
	if got, _ := Canonical("admin:a b{c, d}"); got != "admin:a b{c, d}" {
		t.Errorf("%q", got)
	}
}
