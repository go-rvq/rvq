package perm

import (
	"errors"
	"reflect"
	"slices"
	"testing"
)

func TestParseExprRoundTrip(t *testing.T) {
	for _, line := range []string{
		"admin:posts:*",
		"deny admin:users:*",
		"admin:{posts:{@list, <7<Primeiro post>>:{@get, @edit}}, deny pages:*}",
		"admin:posts:<7>:@get",
		"admin:posts:<*>:#Title:@edit",
		"{admin:content/:*, site:*}",
		`admin:posts:<7<a \<b\> c\\d>>:@get`,
		"admin:users:{list, new}:@x",
		"admin:p[ab,c]:*",
		// the words of the prefixes alone: resources (actions)
		"admin:/site-files:!{edit, reset, deny, allow}",
		"admin:reset:*",
	} {
		e, err := ParseExpr(line)
		if err != nil {
			t.Fatalf("%q: %v", line, err)
		}
		if got := e.String(); got != line {
			t.Errorf("%q: written %q", line, got)
		}
	}
}

func TestParseExprRootList(t *testing.T) {
	for line, want := range map[string]string{
		"admin:content/:*, site:*":   "{admin:content/:*, site:*}",
		"admin:content/:*\nsite:*\n": "{admin:content/:*, site:*}",
		"admin:a,\n deny site:b":     "{admin:a, deny site:b}",
		"admin:a\n\n":                "admin:a",
		"":                           "",
	} {
		e, err := ParseExpr(line)
		if err != nil {
			t.Fatalf("%q: %v", line, err)
		}
		if got := e.String(); got != want {
			t.Errorf("%q: %q, want %q", line, got, want)
		}
	}
	got := ExpandLine("admin:content/:*\ndeny site:*")
	if want := []string{"admin:content/:*", "deny site:*"}; !reflect.DeepEqual(got, want) {
		t.Errorf("expanded %q, want %q", got, want)
	}
}

func TestParseExprTree(t *testing.T) {
	e, err := ParseExpr("deny admin:{posts:<7<T>>:@get, allow pages:*}")
	if err != nil {
		t.Fatal(err)
	}
	if e.Kind != LineDeny || len(e.Nodes) != 2 {
		t.Fatalf("%#v", e)
	}
	g := e.Nodes[1].(*ExprGroup)
	if len(g.Alts) != 2 || g.Alts[1].Kind != LineAllow {
		t.Fatalf("%#v", g)
	}
	r := g.Alts[0].Nodes[1].(*ExprRecord)
	if r.ID != "7" || r.Label != "T" {
		t.Fatalf("%#v", r)
	}
	if p := r.Position(); p.Start != 18 || p.End != 24 {
		t.Errorf("record at %v", p)
	}
	var recs []string
	e.Records(func(path []ExprNode, r *ExprRecord) { recs = append(recs, NodesText(path)+"|"+r.ID) })
	if want := []string{"admin:posts:|7"}; !slices.Equal(recs, want) {
		t.Errorf("records %q, want %q", recs, want)
	}
}

func TestParseExprErrors(t *testing.T) {
	for line, want := range map[string]struct {
		pos int
		msg string
	}{
		"admin:{posts:*":          {6, ErrExprGroupOpen},
		"admin:{a, {b}":           {6, ErrExprGroupOpen},
		"admin:posts:*}":          {13, ErrExprGroupClose},
		"admin:posts:<7":          {12, ErrExprRecordOpen},
		"admin:posts:<7<T>":       {12, ErrExprRecordOpen},
		"admin:posts:<7<T":        {12, ErrExprLabelAtEnd},
		"admin:posts:7>":          {13, ErrExprRecordClose},
		"admin:posts:<>":          {12, ErrExprRecordID},
		"admin:{posts, deny }":    {14, ErrExprPrefixAlone},
		"deny ":                   {0, ErrExprPrefixAlone},
		`admin:posts:\`:           {12, ErrExprEscapeAtEnd},
		"admin:p[ab:*":            {7, ErrExprBracketsOpen},
		"admin:{a, b}, {c}}, d:e": {17, ErrExprGroupClose},
	} {
		_, err := ParseExpr(line)
		var ee *ExprError
		if !errors.As(err, &ee) {
			t.Errorf("%q: %v, want an ExprError", line, err)
			continue
		}
		if ee.Pos != want.pos || ee.Msg != want.msg {
			t.Errorf("%q: at %d %q, want at %d %q", line, ee.Pos, ee.Msg, want.pos, want.msg)
		}
		if CheckLine(line) == nil {
			t.Errorf("%q: CheckLine took it", line)
		}
	}
}

func TestStripLabels(t *testing.T) {
	for line, want := range map[string]string{
		"admin:posts:<7<Primeiro post>>:@get":   "admin:posts:<7>:@get",
		"admin:posts:<7>:@get":                  "admin:posts:<7>:@get",
		`admin:posts:<7<a \<b\>, {c}: d>>:@get`: "admin:posts:<7>:@get",
		"admin:{posts:<1<A>>, pages:<2<B>>}":    "admin:{posts:<1>, pages:<2>}",
		"admin:posts:<7<unclosed":               "admin:posts:<7<unclosed",
		"admin:*":                               "admin:*",
	} {
		if got := StripLabels(line); got != want {
			t.Errorf("%q: %q, want %q", line, got, want)
		}
	}
}

// the decision is the same with and without labels
func TestLabelsDecideNothing(t *testing.T) {
	with := (&Stack{}).Add("t", "admin:posts:<7<Primeiro post, {x}: y>>:{@get, @edit}").Resolve()
	without := (&Stack{}).Add("t", "admin:posts:<7>:{@get, @edit}").Resolve()
	for _, res := range []string{"admin:posts:<7>:@get", "admin:posts:<7>:@edit", "admin:posts:<8>:@get", "admin:posts:<7>:@delete"} {
		if a, b := with.Allows(res, "", ""), without.Allows(res, "", ""); a != b {
			t.Errorf("%s: %v with the label, %v without", res, a, b)
		}
	}
	if !with.Allows("admin:posts:<7>:@get", "", "") {
		t.Error("the record of the label not allowed")
	}
}

func TestNormalize(t *testing.T) {
	for line, want := range map[string]string{
		// a group of one
		"admin:{posts:*}":            "admin:posts:*",
		"admin:{deny posts:*}":       "deny admin:posts:*",
		"allow admin:{deny posts:*}": "deny admin:posts:*",
		"admin:{{posts:*}}":          "admin:posts:*",
		"admin:{}":                   "",
		"admin:{a:{}, b:x}":          "admin:b:x",
		"deny admin:{}":              "",
		// what follows a group, into each alternative
		"admin:{posts, pages}:@list": "admin:{posts:@list, pages:@list}",
		// a group of nothing else
		"admin:{a:x, {b:y, deny c:z}}":       "admin:{a:x, b:y, deny c:z}",
		"admin:{a:x, deny {b:y, allow c:z}}": "admin:{a:x, deny b:y, allow c:z}",
		// the same segment: one
		"admin:{posts:@list, posts:@get}":               "admin:posts:{@list, @get}",
		"admin:{posts:@list, pages:*, deny posts:@get}": "admin:{posts:{@list, deny @get}, pages:*}",
		"admin:{posts:<7<A>>:@get, posts:<7>:@edit}":    "admin:posts:<7<A>>:{@get, @edit}",
		// an allow between: not moved, unless the segments are others, of no
		// globs
		"admin:{posts:@list, allow pages:*, posts:@get}": "admin:{posts:{@list, @get}, allow pages:*}",
		"admin:{posts:@list, allow *:@get, posts:@get}":  "admin:{posts:@list, allow *:@get, posts:@get}",
		"admin:{*:@list, allow posts:x, *:@get}":         "admin:{*:@list, allow posts:x, *:@get}",
		// the root: its items the scopes, one each
		"{admin:a, site:b, admin:c}":                "{admin:{a, c}, site:b}",
		"admin:a, admin:b":                          "admin:{a, b}",
		"admin:a\nallow site:*\nadmin:b":            "{admin:{a, b}, allow site:*}",
		"{admin:posts:*}":                           "admin:posts:*",
		"deny {admin:a, admin:b}":                   "deny admin:{a, b}",
		"{admin:{posts:@list}, admin:{posts:@get}}": "admin:posts:{@list, @get}",
	} {
		got, err := Canonical(line)
		if err != nil {
			t.Fatalf("%q: %v", line, err)
		}
		if got != want {
			t.Errorf("%q: %q, want %q", line, got, want)
		}
		// the same lines, in an order that decides the same
		if !sameDecision(t, line, got) {
			t.Errorf("%q and %q decide otherwise", line, got)
		}
		// canonical: once is all
		if again, _ := Canonical(got); again != got {
			t.Errorf("%q: again %q", got, again)
		}
	}
}

// sameDecision says the lines a and b, alone in a stack, decide the same of
// the resources their lines name and of some others.
func sameDecision(t *testing.T, a, b string) bool {
	t.Helper()
	sa := (&Stack{}).Add("t", a).Resolve()
	sb := (&Stack{}).Add("t", b).Resolve()
	var res []string
	for _, l := range append(ExpandLine(a), ExpandLine(b)...) {
		_, p := ParseLine(l)
		res = append(res, p, p+":x", "admin:posts:@get", "admin:pages:@get", "admin:x:@get", "site:b", "admin:c")
	}
	for _, r := range res {
		if sa.Allows(r, "", "") != sb.Allows(r, "", "") {
			t.Logf("%s: %v / %v", r, sa.Allows(r, "", ""), sb.Allows(r, "", ""))
			return false
		}
	}
	return true
}

// a word of a prefix alone is a resource: the line stays in the stack
func TestPrefixWordAlone(t *testing.T) {
	line := "admin:/site-files:!{create,edit,reset,share}"
	if err := CheckLine(line); err != nil {
		t.Fatalf("%q: %v", line, err)
	}
	st := (&Stack{}).Add("t", line).Resolve()
	for _, r := range []string{"admin:/site-files:!reset", "admin:/site-files:!edit"} {
		if !st.Allows(r, "", "") {
			t.Errorf("%s not allowed", r)
		}
	}
	if CheckLine("deny") == nil || CheckLine("deny ") == nil {
		t.Error("a prefix alone taken")
	}
}

// a line of an odd grammar, stored before: in the stack as it always was
func TestStackKeepsOddLines(t *testing.T) {
	line := "admin:posts:*}"
	if CheckLine(line) == nil {
		t.Fatal("taken by CheckLine")
	}
	if st := (&Stack{}).Add("t", line).Resolve(); len(st.Lines) != 1 {
		t.Errorf("lines %v", st.Lines)
	}
}
