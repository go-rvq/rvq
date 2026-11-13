package perm

import (
	"math/rand"
	"strings"
	"testing"
)

// The order of a value: the items of a level from the general to the
// particular; a deny after a grant of the same; an allow and a reset where
// they are, the lines before them before; twice, the same.
func TestSortValue(t *testing.T) {
	for _, c := range []struct{ in, want string }{
		{"admin:posts:{deny <*>:@delete, !publish, @edit, %readers, @list}",
			"admin:posts:{%readers, @list, @edit, !publish, deny <*>:@delete}"},
		{"admin:posts:{#Title:@edit, /report:@get, <7>:@get, <*>:@get, $Main:@get, &Config:@get, *}",
			"admin:posts:{*, <*>:@get, <7>:@get, #Title:@edit, $Main:@get, &Config:@get, /report:@get}"},
		{"{site:*, admin:{users:*, posts:@list}}", "{admin:{posts:@list, users:*}, site:*}"},
		// of an allow: the lines of resources apart move, those of the same
		// keep their order
		{"admin:{users:@x, posts:@list, allow posts:@get, pages:*, deny ab:*}",
			"admin:{deny ab:*, pages:*, posts:{@list, allow @get}, users:@x}"},
		{"admin:{users:@list, reset users:*, posts:*}", "admin:{posts:*, users:{@list, reset *}}"},
		{"admin:{posts:@list, allow posts:*, posts:@get}", "admin:posts:{@list, allow *, @get}"},
	} {
		got, err := SortValue(c.in)
		if err != nil {
			t.Fatal(err)
		}
		if got != c.want {
			t.Errorf("%q:\n got %q\nwant %q", c.in, got, c.want)
		}
		if again, _ := SortValue(got); again != got {
			t.Errorf("%q: twice %q", got, again)
		}
	}
	if got, _ := FormatLine("admin:{users:*, posts:*}", WithSort()); got != "admin:{posts:*, users:*}" {
		t.Errorf("FormatLine WithSort: %q", got)
	}
	// a set that gives an allow: a barrier
	t.Cleanup(func() { SetSetExpander(nil) })
	SetSetExpander(func(resource, name string) ([]string, bool) {
		return []string{"allow " + resource + "@get"}, true
	})
	got := SortLines([]string{"admin:x:@list", "admin:x:%forcing", "admin:x:*", "admin:a:*"})
	if strings.Join(got, "|") != "admin:a:*|admin:x:@list|admin:x:%forcing|admin:x:*" {
		t.Errorf("a set of an allow: %q", got)
	}
}

// Sorted, a value decides the same: of random lines — grants, denies,
// allows, resets of resources and globs that cover each other —, what a
// stack of them allows, of each permission, before and after.
func TestSortLinesDecidesTheSame(t *testing.T) {
	r := rand.New(rand.NewSource(7))
	resources := []string{"admin:*", "admin:posts:<*>:*", "admin:posts:@*", "admin:{posts,users}:*", "admin:posts:*", "admin:posts:@list", "admin:posts:@get", "admin:posts:<*>:@edit",
		"admin:posts:<7>:@edit", "admin:posts:<*>:#Title:*", "admin:pages:*", "admin:pages:@list", "admin:{posts,pages}:@delete",
		"admin:users:*", "admin:users:@edit", "admin:*:@list"}
	kinds := []string{"", "", "", "deny ", "deny ", "allow ", "reset "}
	checks := []string{"admin:posts:@list", "admin:posts:@get", "admin:posts:<7>:@edit", "admin:posts:<8>:@edit",
		"admin:posts:<7>:#Title:@edit", "admin:pages:@list", "admin:pages:@delete", "admin:posts:@delete",
		"admin:users:@edit", "admin:users:@list", "admin:other:@list"}
	for n := 0; n < 20000; n++ {
		var lines []string
		for i := 0; i < 2+r.Intn(7); i++ {
			lines = append(lines, kinds[r.Intn(len(kinds))]+resources[r.Intn(len(resources))])
		}
		value, err := JoinLines(lines)
		if err != nil {
			t.Fatal(err)
		}
		before, _ := LinesOf(value)
		after := SortLines(before)
		a := (&Stack{}).Add("v", before...).Resolve()
		b := (&Stack{}).Add("v", after...).Resolve()
		for _, c := range checks {
			if a.Allows(c, "", "") != b.Allows(c, "", "") {
				t.Fatalf("%s decides another on %s:\n%q\n%q", value, c, before, after)
			}
		}
	}
}
