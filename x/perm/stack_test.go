package perm

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/ory/ladon"
)

// A line is the lines its groups of alternatives are, in their order; the
// prefix of an alternative takes the place of the line's; groups in groups
// too; the "!" of an action is the action's name.
func TestExpandLine(t *testing.T) {
	for _, c := range []struct {
		line string
		want []string
	}{
		{"admin:users:*", []string{"admin:users:*"}},
		{"admin:users:{list,new}", []string{"admin:users:list", "admin:users:new"}},
		{"deny admin:users:{list, new}", []string{"deny admin:users:list", "deny admin:users:new"}},
		{"admin:users:{*, deny @get}", []string{"admin:users:*", "deny admin:users:@get"}},
		{"admin:users:{*, deny {@get, @set}}", []string{"admin:users:*", "deny admin:users:@get", "deny admin:users:@set"}},
		{"deny admin:users:{*, allow @list}", []string{"deny admin:users:*", "allow admin:users:@list"}},
		{"admin:{users,pages}:{@list,@edit}", []string{"admin:users:@list", "admin:users:@edit", "admin:pages:@list", "admin:pages:@edit"}},
		{"admin:/site-files:{@get,!git}", []string{"admin:/site-files:@get", "admin:/site-files:!git"}},
		{"admin:users:{\n\tlist\n\tnew\n}", []string{"admin:users:list", "admin:users:new"}},
		{"admin:users:<5>:@edit", []string{"admin:users:<5>:@edit"}},
		{"   ", nil},
	} {
		if got := ExpandLine(c.line); strings.Join(got, "|") != strings.Join(c.want, "|") {
			t.Errorf("%q: %q, want %q", c.line, got, c.want)
		}
	}
}

// A prefix is its word and a space, before a resource.
func TestCheckLine(t *testing.T) {
	for _, l := range []string{"admin:*", "deny admin:users:*", "allow admin:users:@list", "reset admin:users:*", "", "admin:users:{*, deny @get}"} {
		if err := CheckLine(l); err != nil {
			t.Errorf("%q: %v", l, err)
		}
	}
	for _, l := range []string{"deny", "allow ", "reset", "deny   "} {
		if CheckLine(l) == nil {
			t.Errorf("%q: no error", l)
		}
	}
}

func situations(s *Stack) string {
	var out []string
	for _, l := range s.Lines {
		st := "active"
		switch {
		case l.ReplacedBy > 0:
			st = "by " + itoa(l.ReplacedBy)
		case l.Kind == LineReset:
			st = "reset"
		}
		for _, e := range l.ExceptBy {
			st += " -" + itoa(e)
		}
		out = append(out, itoa(l.Index)+" "+l.Text()+": "+st)
	}
	return strings.Join(out, "\n")
}

func itoa(i int) string { return strconv.Itoa(i) }

// Resolved, a line covered by a later one that replaces its kind is replaced
// by it — the first —; the rest is active.
func TestResolve(t *testing.T) {
	for _, c := range []struct {
		name  string
		lines []string
		want  string
	}{
		{"a deny after: the grant it covers replaced",
			[]string{"admin:users:list", "deny admin:users:*"},
			"1 admin:users:list: by 2\n2 deny admin:users:*: active"},
		{"an allow after a deny it covers",
			[]string{"deny admin:users:list", "allow admin:users:*"},
			"1 deny admin:users:list: by 2\n2 allow admin:users:*: active"},
		{"an allow narrower than the deny: the deny but what it matches",
			[]string{"deny admin:users:*", "allow admin:users:list"},
			"1 deny admin:users:*: active -2\n2 allow admin:users:list: active"},
		{"a reset narrower: any line but what it matches",
			[]string{"admin:*", "deny admin:users:*", "reset admin:users:list"},
			"1 admin:*: active -3\n2 deny admin:users:*: active -3\n3 reset admin:users:list: reset"},
		{"replaced later: no exceptions",
			[]string{"deny admin:users:*", "allow admin:users:list", "reset admin:*"},
			"1 deny admin:users:*: by 3\n2 allow admin:users:list: by 3\n3 reset admin:*: reset"},
		{"a narrower grant after a deny: no exception",
			[]string{"deny admin:users:*", "admin:users:list"},
			"1 deny admin:users:*: active\n2 admin:users:list: active"},
		{"a grant after a deny: the deny stays",
			[]string{"deny admin:users:list", "admin:*"},
			"1 deny admin:users:list: active\n2 admin:*: active"},
		{"a grant after a grant it covers",
			[]string{"admin:users:list", "admin:*"},
			"1 admin:users:list: by 2\n2 admin:*: active"},
		{"a reset: the lines it covers replaced, one wider but what it matches",
			[]string{"admin:*", "admin:users:list", "deny admin:users:new", "reset admin:users:*"},
			"1 admin:*: active -4\n2 admin:users:list: by 4\n3 deny admin:users:new: by 4\n4 reset admin:users:*: reset"},
		{"replaced by the first",
			[]string{"admin:users:list", "deny admin:users:*", "reset admin:*"},
			"1 admin:users:list: by 2\n2 deny admin:users:*: by 3\n3 reset admin:*: reset"},
		{"a group broken",
			[]string{"admin:users:{*, deny @get}"},
			"1 admin:users:*: active\n2 deny admin:users:@get: active"},
	} {
		s := (&Stack{}).Add("role: R", c.lines...).Resolve()
		if got := situations(s); got != c.want {
			t.Errorf("%s:\n%s\nwant\n%s", c.name, got, c.want)
		}
	}
}

// The active lines are policies of the subject: a grant (an allow is one)
// allowed, a deny denied.
func TestStackPolicies(t *testing.T) {
	s := (&Stack{}).Add("role: R", "deny admin:users:list", "allow admin:users:*", "reset admin:pages:*").Resolve()
	ps := StackPolicies("user:1", s.Active())
	if len(ps) != 1 || ps[0].policy.Effect != Allowed || ps[0].policy.Resources[0] != "admin:users:*" ||
		ps[0].policy.Subjects[0] != "user:1" {
		t.Fatalf("policies: %+v", ps)
	}
}

type fakeSource map[string]*Stack

func (f fakeSource) StackOf(r *http.Request) (string, string, func() (*Stack, error)) {
	u := r.Header.Get("X-User")
	s, ok := f[u]
	if !ok {
		return "", "", nil
	}
	return "user:" + u, "1", func() (*Stack, error) {
		return &Stack{Lines: append([]*Line(nil), s.Lines...)}, nil
	}
}

// A user's request is decided by the active lines of the user's stack, as
// its subject alone: the policies of the code (of the subject "*", of a
// role) no longer decide it — they are lines of the stack, first. A request
// of no user is decided by the policies, as before.
func TestStacksDecide(t *testing.T) {
	b := New().Policies(
		PolicyFor(Anybody).WhoAre(Denied).ToDo(Anything).On("admin:secret:*"),
		PolicyFor("editor").WhoAre(Allowed).ToDo(Anything).On("admin:pages:*"),
		PolicyFor(Anybody).WhoAre(Allowed).ToDo(Anything).On("admin:public:*"),
	)
	build := func(roles []string, lines ...string) *Stack {
		s := &Stack{}
		for _, l := range b.CodeLines(append([]string{Anybody}, roles...)) {
			s.AddLine(l)
		}
		return s.Add("role: R", lines...)
	}
	src := fakeSource{
		"ana": build([]string{"editor"}, "admin:users:*", "deny admin:users:@edit", "allow admin:secret:*"),
		"bob": build(nil, "admin:*"),
		"cid": build(nil, "admin:*", "deny admin:users:*", "allow admin:users:@list", "reset admin:pages:@get"),
	}
	b.Stacks(src)
	asks := func(user, resource string) bool {
		parts := strings.Split(resource, ":")
		r := httptest.NewRequest("GET", "/", nil)
		if user != "" {
			r.Header.Set("X-User", user)
		}
		return NewVerifier("admin", b).Spawn().Do(parts[len(parts)-1]).On(parts[1:len(parts)-1]...).WithReq(r).IsAllowed() == nil
	}
	for _, c := range []struct {
		user, resource string
		want           bool
	}{
		{"ana", "admin:users:@list", true},
		{"ana", "admin:users:@edit", false},
		{"ana", "admin:pages:@list", true},  // the code's, for editor
		{"ana", "admin:secret:@get", true},  // the code's deny, replaced by her allow
		{"ana", "admin:public:@get", true},  // the code's, for anybody
		{"ana", "admin:other:@get", false},  // nothing allows
		{"bob", "admin:pages:@list", true},  // his admin:*
		{"bob", "admin:secret:@get", false}, // the code's deny: admin:* comes after, a grant
		{"cid", "admin:users:@list", true},  // the allow, under a wider deny
		{"cid", "admin:users:@edit", false}, // the deny, but the list
		{"cid", "admin:pages:@list", true},  // admin:*
		{"cid", "admin:pages:@get", false},  // admin:*, but what the reset matches
		{"", "admin:public:@get", true},     // no user: the policies, as before
		{"", "admin:secret:@get", false},
	} {
		if got := asks(c.user, c.resource); got != c.want {
			t.Errorf("%s %s: %v, want %v", c.user, c.resource, got, c.want)
		}
	}
	if !b.Registered("user:ana") || b.Registered("user:carl") {
		t.Error("the stacks registered")
	}
	// a policy of the code made: the stacks made again, with it
	b.CreatePolicies(PolicyFor(Anybody).WhoAre(Denied).ToDo(Anything).On("admin:users:@list"))
	src["ana"] = build([]string{"editor"}, "admin:users:*", "deny admin:users:@edit", "allow admin:secret:*")
	if asks("ana", "admin:users:@list") {
		// admin:users:* after the code's deny is a grant: the deny stays
		t.Error("the code's new deny: not in the stack")
	}
}

// The conditions of a policy of the code are the line's: it counts only when
// the request fulfills them.
func TestStackConditions(t *testing.T) {
	b := New().Policies(
		PolicyFor(Anybody).WhoAre(Allowed).ToDo(Anything).On("admin:roles:*").Given(Conditions{
			"editor": &ladon.BooleanCondition{BooleanValue: true},
		}),
	)
	b.ContextFunc(func(r *http.Request, _ []interface{}) Context {
		return Context{"editor": r.Header.Get("X-Editor") == "1"}
	})
	s := &Stack{}
	for _, l := range b.CodeLines([]string{Anybody}) {
		s.AddLine(l)
	}
	b.Stacks(fakeSource{"ana": s})
	asks := func(editor string) bool {
		r := httptest.NewRequest("GET", "/", nil)
		r.Header.Set("X-User", "ana")
		r.Header.Set("X-Editor", editor)
		return NewVerifier("admin", b).Spawn().Do("@list").On("roles").WithReq(r).IsAllowed() == nil
	}
	if !asks("1") || asks("0") {
		t.Errorf("the condition: %v %v", asks("1"), asks("0"))
	}
}

// An access key restricts by the stack of its subject (WithRestriction): what
// the user may, but its denies.
func TestStackRestriction(t *testing.T) {
	b := New()
	b.Stacks(fakeSource{"ana": (&Stack{}).Add("role: R", "admin:*")})
	if err := b.Register("key:1", "1", func() (*Stack, error) {
		return (&Stack{}).Add("key", "*", "deny admin:users:*"), nil
	}); err != nil {
		t.Fatal(err)
	}
	asks := func(resource ...string) bool {
		r := httptest.NewRequest("GET", "/", nil)
		r.Header.Set("X-User", "ana")
		r = r.WithContext(WithRestriction(r.Context(), "key:1"))
		return NewVerifier("admin", b).Spawn().Do("@list").On(resource...).WithReq(r).IsAllowed() == nil
	}
	if !asks("pages") || asks("users") {
		t.Errorf("the key's deny: pages %v users %v", asks("pages"), asks("users"))
	}
}

// A policy of before is lines: a resource each, "deny " before a deny's,
// the actions it names appended.
func TestPolicyLines(t *testing.T) {
	for _, c := range []struct {
		effect             string
		actions, resources []string
		want               string
	}{
		{Allowed, []string{"*"}, []string{"admin:pages:*", "admin:{a,b}:*@{list,get}"}, "admin:pages:*|admin:{a,b}:*@{list,get}"},
		{Denied, nil, []string{"admin:users:*"}, "deny admin:users:*"},
		{Allowed, []string{"@get", "@list"}, []string{"admin:posts:*"}, "admin:posts:*@get|admin:posts:*@list"},
		{Allowed, []string{"*"}, []string{"admin:x:*, admin:y:*"}, "admin:x:*|admin:y:*"},
	} {
		if got := strings.Join(PolicyLines(c.effect, c.actions, c.resources), "|"); got != c.want {
			t.Errorf("%v %v %v: %s, want %s", c.effect, c.actions, c.resources, got, c.want)
		}
	}
}

// The active lines as a tree: their beginnings once, what follows each a
// group, recursively, a deny's prefix before its leaf; expanded, the lines
// again.
func TestTree(t *testing.T) {
	lines := []string{"admin:users:list", "admin:users:new", "deny admin:users:@get", "admin:pages:*", "site:x"}
	got := Tree(lines)
	want := "admin:{\n\tusers:{\n\t\tlist\n\t\tnew\n\t\tdeny @get\n\t}\n\tpages:*\n}\nsite:x"
	if got != want {
		t.Fatalf("the tree:\n%s\nwant\n%s", got, want)
	}
	var back []string
	for _, part := range strings.Split(got, "\nsite:") {
		back = append(back, ExpandLine(part)...)
	}
	if strings.Join(back, "|") != "admin:users:list|admin:users:new|deny admin:users:@get|admin:pages:*|x" {
		t.Errorf("expanded back: %q", back)
	}
	if Tree([]string{"admin:users:list"}) != "admin:users:list" {
		t.Errorf("one line: %q", Tree([]string{"admin:users:list"}))
	}
}

// The policies of the database are lines after the code's (DBLines), of the
// origin "database: …"; a policy changed — of the database too —, the stacks
// are made again.
func TestStackDBLines(t *testing.T) {
	b := New().Policies(
		PolicyFor("Editor").WhoAre(Allowed).ToDo(Anything).On("admin:pages:*"),
		PolicyFor("user:1").WhoAre(Allowed).ToDo(Anything).On("admin:posts:<7>:*").ID("db-1"),
	)
	b.dbIDs = map[string]bool{"db-1": true}
	text := func(lines []*Line) (out []string) {
		for _, l := range lines {
			out = append(out, l.Text()+" ("+l.Origin+")")
		}
		return
	}
	subjects := []string{"Editor", "user:1"}
	if got := strings.Join(text(b.CodeLines(subjects)), "|"); got != "admin:pages:* (code: Editor)" {
		t.Errorf("the code's: %s", got)
	}
	if got := strings.Join(text(b.DBLines(subjects)), "|"); got != "admin:posts:<7>:* (database: user:1)" {
		t.Errorf("the database's: %s", got)
	}
	builds := 0
	build := func() (*Stack, error) { builds++; return &Stack{}, nil }
	_ = b.Register("user:1", "1", build)
	_ = b.Register("user:1", "1", build)
	b.CreatePolicies(PolicyFor("user:1").WhoAre(Denied).ToDo(Anything).On("admin:posts:*").ID("db-2"))
	_ = b.Register("user:1", "1", build)
	if builds != 2 {
		t.Errorf("made %d times, want 2", builds)
	}
}
