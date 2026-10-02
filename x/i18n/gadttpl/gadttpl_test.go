package gadttpl

import (
	"strings"
	"testing"

	"github.com/gad-lang/gad"
)

func TestRender(t *testing.T) {
	tpl := Template(`{ if valid begin }{= enabled }, keep for {= join_and(", ", " and ", days, weeks, other) }{ else }<span>Do not keep</span>{ end }.`)
	for _, tc := range []struct {
		data gad.Dict
		want string
	}{
		{gad.Dict{"valid": gad.True, "enabled": gad.Str("On"), "days": gad.Str("2 days"), "weeks": gad.Str(""), "other": gad.Str("1 year"), "join_and": JoinAnd},
			"On, keep for 2 days and 1 year."},
		{gad.Dict{"valid": gad.False, "enabled": gad.Str("On"), "days": gad.Str(""), "weeks": gad.Str(""), "other": gad.Str(""), "join_and": JoinAnd},
			"<span>Do not keep</span>."},
	} {
		got, err := tpl.Render(tc.data)
		if err != nil || got != tc.want {
			t.Errorf("got %q, %v; want %q", got, err, tc.want)
		}
	}
}

func TestJoinAnd(t *testing.T) {
	tpl := Template(`{= join_and(", ", " and ", a, b, c, d) }`)
	for _, tc := range []struct{ a, b, c, d, want string }{
		{"", "", "", "", ""},
		{"x", "", "", "", "x"},
		{"x", "", "y", "", "x and y"},
		{"x", "y", "z", "w", "x, y, z and w"},
	} {
		got, err := tpl.Render(gad.Dict{"a": gad.Str(tc.a), "b": gad.Str(tc.b), "c": gad.Str(tc.c), "d": gad.Str(tc.d), "join_and": JoinAnd})
		if err != nil || got != tc.want {
			t.Errorf("%v: got %q, %v", tc, got, err)
		}
	}
}

func TestEmptyAndError(t *testing.T) {
	if got, err := Template("  ").Render(nil); got != "" || err != nil {
		t.Errorf("empty: %q, %v", got, err)
	}
	if got := Template(`{= nope( }`).HTML(nil); !strings.Contains(string(got), "text-error") {
		t.Errorf("error: %q", got)
	}
}
