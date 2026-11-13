package perm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gobwas/glob"
	"github.com/ory/ladon"
)

// The permissions of a user are a stack: lines in order — the policies of the
// code first, then the database's (what is granted a record, a share), then
// the ones of each of the user's roles, in the order of the user's roles —,
// each a resource (a glob pattern, as the policies') and maybe a prefix, its
// word and a space:
//
//	admin:pages:*            allows what it matches, unless a deny does
//	deny admin:users:*       denies what it matches
//	allow admin:users:*      the denies before it no longer count for what
//	                         it matches; it allows what it matches
//	reset admin:users:*      the lines before it no longer count for what it
//	                         matches; it allows nor denies anything
//
// A group of alternatives is lines, in its order: "admin:users:{list,new}" is
// "admin:users:list" and "admin:users:new"; an alternative may have its own
// prefix — "admin:users:{*, deny @get}" is "admin:users:*" and then
// "deny admin:users:@get" —, groups in groups too.
//
// A line covers another when its pattern matches the other's ("admin:*"
// covers "admin:users:list"; "admin:users:list" does not cover
// "admin:users:*"). Resolved (Resolve), a line covered by a later one is
// replaced (ReplacedBy): by a deny or an allow, a line of no prefix or a
// deny; by a reset, any; by a line of no prefix, an earlier one of no
// prefix. A wider line before an allow (a deny) or a reset (any) that it
// covers is not replaced: it no longer counts for what the later one matches
// alone (ExceptBy) — "deny admin:users:*" and then "allow
// admin:users:@list" deny the users but their list. What is left — the
// active lines, of no prefix (an allow is one) or deny, each with its
// exceptions — decides as the policies always did: a resource is allowed
// when an active line of no prefix matches it and no active deny does.

// Prefixes of a line of a stack: the word and a space.
const (
	PrefixDeny  = "deny "
	PrefixAllow = "allow "
	PrefixReset = "reset "
	// PrefixUnset is of an edit only (ExprAlt.Merge): never a line of a
	// stack
	PrefixUnset = "unset "
)

// LineKind is what a line of a stack does.
type LineKind int

const (
	// LineGrant allows (no prefix).
	LineGrant LineKind = iota
	// LineDeny denies.
	LineDeny
	// LineAllow replaces the denies before it it covers, and allows.
	LineAllow
	// LineReset replaces the lines before it it covers.
	LineReset
	// LineUnset is an edit, not a line of a stack: in a merge
	// (ExprAlt.Merge), it takes out the lines before it it covers; a stack
	// never takes it (CheckLine: ErrUnset)
	LineUnset
)

// LineKinds are the kinds of a prefix, in the order they are read.
var LineKinds = []LineKind{LineDeny, LineAllow, LineReset, LineUnset}

// ErrUnset is a line of a stack with an unset: an edit of a merge only.
var ErrUnset = errors.New(`"unset" takes out lines in a merge: it is no line of a stack`)

// Prefix is the prefix of the kind: "" for a grant.
func (k LineKind) Prefix() string {
	switch k {
	case LineDeny:
		return PrefixDeny
	case LineAllow:
		return PrefixAllow
	case LineReset:
		return PrefixReset
	case LineUnset:
		return PrefixUnset
	}
	return ""
}

// ErrLine is a line of no resource after its prefix, or of a prefix not
// followed by its space.
var ErrLine = errors.New(`a permission is a resource, maybe after a prefix — its word and a space: "deny admin:users:*"`)

// ParseLine is the kind and the pattern of a line: its prefix taken off.
func ParseLine(s string) (LineKind, string) {
	s = strings.TrimSpace(s)
	for _, k := range LineKinds {
		if p := k.Prefix(); strings.HasPrefix(s, p) {
			return k, strings.TrimSpace(s[len(p):])
		}
	}
	return LineGrant, s
}

// CheckLine says what is wrong with a line: nil when nothing is (a blank one
// is left out) — its prefix and its grammar (ParseExpr): what a field that
// writes lines takes.
func CheckLine(s string) error {
	if err := checkPrefix(s); err != nil {
		return err
	}
	if strings.TrimSpace(s) == "" {
		return nil
	}
	_, err := ParseExpr(s)
	return err
}

// checkPrefix says what is wrong with the prefix of a line: what a stack
// leaves out (Stack.Add) — a line of an odd grammar, stored before the
// grammar was, decides as it always did.
func checkPrefix(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	for _, k := range LineKinds {
		w := strings.TrimSpace(k.Prefix())
		if s == w {
			return ErrLine
		}
	}
	k, p := ParseLine(s)
	if p == "" {
		return ErrLine
	}
	for _, l := range ExpandLine(s) {
		lk, e := ParseLine(l)
		if lk == LineUnset {
			return ErrUnset
		}
		if e == "" && k != LineGrant {
			return ErrLine
		}
	}
	return nil
}

// ExpandLine is the lines a line is, its groups of alternatives broken in
// their order: "deny admin:users:{list, new}" is "deny admin:users:list" and
// "deny admin:users:new"; the prefix of an alternative takes the place of
// the line's. Blank lines are none. The labels of its records are none of it
// (StripLabels); its alternatives at the root, by commas or lines, a group:
// "admin:a, site:b" is "admin:a" and "site:b".
func ExpandLine(s string) []string {
	kind, pattern := ParseLine(StripLabels(s))
	pattern = rootList(pattern)
	if pattern == "" {
		return nil
	}
	var out []string
	for _, e := range expandPattern(pattern) {
		k := kind
		if e.kind != nil {
			k = *e.kind
		}
		out = append(out, k.Prefix()+e.pattern)
	}
	return out
}

type expanded struct {
	kind    *LineKind
	pattern string
}

// expandPattern breaks the first group of alternatives of p — "{…}" — and,
// recursively, the rest: each alternative in the place of the group, its
// prefix, when it has one, the line's.
func expandPattern(p string) []expanded {
	start, end := firstGroup(p)
	if start < 0 {
		return []expanded{{pattern: p}}
	}
	prefix, suffix := p[:start], p[end+1:]
	var out []expanded
	for _, alt := range splitAlternatives(p[start : end+1]) {
		var kind *LineKind
		if k, a := ParseLine(alt); k != LineGrant {
			kind, alt = &k, a
		}
		for _, e := range expandPattern(prefix + alt + suffix) {
			if e.kind == nil {
				e.kind = kind
			}
			out = append(out, e)
		}
	}
	return out
}

// firstGroup is where the first group of alternatives of p begins and ends,
// "{…}" (its outer braces); -1 when none.
func firstGroup(p string) (start, end int) {
	depth := 0
	start = -1
	for i := 0; i < len(p); i++ {
		switch p[i] {
		case '\\':
			i++
		case '{':
			if depth == 0 {
				start = i
			}
			depth++
		case '}':
			if depth == 0 {
				continue
			}
			depth--
			if depth == 0 {
				return start, i
			}
		}
	}
	return -1, -1
}

// Line is a line of a stack: its place (Index, from 1), what it does, its
// pattern, where it comes from (Origin: "role: Editor", "code: …"), and,
// resolved, the line that replaced it (0: none).
type Line struct {
	Index   int
	Kind    LineKind
	Pattern string
	Origin  string
	// Actions, of a policy of the code that names them: what the request's
	// action must be one of (nil: any — the action ends the resource).
	Actions []string
	// Conditions, of a policy of the code: the line counts only when the
	// request fulfills them.
	Conditions ladon.Conditions
	// ReplacedBy is the index of the line that replaced it (Resolve).
	ReplacedBy int
	// ExceptBy are the indexes of the later lines, narrower, that it no
	// longer counts for (Resolve): an allow after a deny, a reset after any
	// line, each covered by it.
	ExceptBy []int
	// Except are the patterns of ExceptBy: what it no longer matches.
	Except []string
	// Via is the set the line comes of ("readers": a line "…:%readers"
	// expanded, Stack.Add); "" none.
	Via string
	// Ignored, a line of a set that gives nothing there (no such set,
	// nothing for its kind of resource): kept, to be shown, but it counts
	// for nothing.
	Ignored bool

	g  glob.Glob
	eg []glob.Glob
}

// Text is the line as written: its prefix, its pattern and, of a policy of
// the code that names its actions, them.
func (l *Line) Text() string {
	t := l.Kind.Prefix() + l.Pattern
	if anyAction(l.Actions) {
		return t
	}
	return t + " [" + strings.Join(l.Actions, ", ") + "]"
}

// Active says whether the line counts: not replaced, not a reset, not a set
// that gives nothing (Ignored).
func (l *Line) Active() bool { return l.ReplacedBy == 0 && l.Kind != LineReset && !l.Ignored }

// Denies says whether the line, active, denies.
func (l *Line) Denies() bool { return l.Kind == LineDeny }

func (l *Line) glob() glob.Glob {
	if l.g == nil {
		g, err := glob.Compile(l.Pattern)
		if err != nil {
			// a pattern that does not compile matches only itself
			g = glob.MustCompile(glob.QuoteMeta(l.Pattern))
		}
		l.g = g
	}
	return l.g
}

func anyAction(actions []string) bool {
	if len(actions) == 0 {
		return true
	}
	for _, a := range actions {
		if a == Anything {
			return true
		}
	}
	return false
}

// covers says whether l covers o: its pattern matches the other's, its
// actions are any or the other's too, and neither has conditions.
func (l *Line) covers(o *Line) bool {
	if len(l.Conditions) > 0 || len(o.Conditions) > 0 {
		return false
	}
	if !anyAction(l.Actions) {
		if anyAction(o.Actions) {
			return false
		}
		for _, a := range o.Actions {
			if !containsString(l.Actions, a) {
				return false
			}
		}
	}
	return l.glob().Match(o.Pattern)
}

func containsString(l []string, s string) bool {
	for _, e := range l {
		if e == s {
			return true
		}
	}
	return false
}

// replaces says whether l, later, replaces o, earlier, when it covers it.
func (l *Line) replaces(o *Line) bool {
	switch l.Kind {
	case LineDeny, LineAllow:
		return o.Kind != LineReset
	case LineReset:
		return true
	default: // a grant: an earlier grant (an allow is one)
		return o.Kind == LineGrant || o.Kind == LineAllow
	}
}

// excepts says whether l, later and narrower, takes from o what it matches:
// an allow from a deny, a reset from any line.
func (l *Line) excepts(o *Line) bool {
	switch l.Kind {
	case LineAllow:
		return o.Kind == LineDeny
	case LineReset:
		return o.Kind != LineReset
	}
	return false
}

// Stack is the lines of the permissions of a user, in order.
type Stack struct {
	Lines []*Line
}

// Add adds the lines of origin, each expanded (ExpandLine), the blank ones
// and those of a wrong prefix left out (a grammar odd decides as it always
// did: CheckLine is for what is written now). A line of a set ("…:%readers")
// is the lines the set gives there (ExpandSet), each of it (Via), its kind
// the line's unless its own says another (SetLineKind); a set that gives
// nothing there, the line kept, ignored.
func (s *Stack) Add(origin string, lines ...string) *Stack {
	for _, l := range lines {
		if checkPrefix(l) != nil {
			continue
		}
		for _, e := range ExpandLine(l) {
			k, p := ParseLine(e)
			if res, name, ok := SetRef(p); ok {
				got, ok := ExpandSet(res, name)
				if !ok {
					s.Lines = append(s.Lines, &Line{Index: len(s.Lines) + 1, Kind: k, Pattern: p, Origin: origin,
						Via: name, Ignored: true})
					continue
				}
				for _, g := range got {
					for _, ge := range ExpandLine(g) {
						gk, gp := ParseLine(ge)
						s.Lines = append(s.Lines, &Line{Index: len(s.Lines) + 1, Kind: SetLineKind(k, gk), Pattern: gp,
							Origin: origin, Via: name})
					}
				}
				continue
			}
			s.Lines = append(s.Lines, &Line{Index: len(s.Lines) + 1, Kind: k, Pattern: p, Origin: origin})
		}
	}
	return s
}

// AddLine adds a line made already (a policy of the code's).
func (s *Stack) AddLine(l *Line) *Stack {
	l.Index = len(s.Lines) + 1
	s.Lines = append(s.Lines, l)
	return s
}

// Resolve sets the line that replaced each: the first later line that covers
// it and replaces its kind.
func (s *Stack) Resolve() *Stack {
	for i, l := range s.Lines {
		l.ReplacedBy, l.ExceptBy, l.Except = 0, nil, nil
		if l.Ignored {
			continue
		}
		for _, later := range s.Lines[i+1:] {
			if later.Ignored {
				continue
			}
			if later.replaces(l) && later.covers(l) {
				l.ReplacedBy = later.Index
				l.ExceptBy, l.Except = nil, nil
				break
			}
			if later.excepts(l) && l.covers(later) {
				l.ExceptBy = append(l.ExceptBy, later.Index)
				l.Except = append(l.Except, later.Pattern)
			}
		}
	}
	return s
}

// Active are the lines that count, in order (Resolve first).
func (s *Stack) Active() []*Line {
	var out []*Line
	for _, l := range s.Lines {
		if l.Active() {
			out = append(out, l)
		}
	}
	return out
}

// match says whether the line, active, counts for the resource asked with
// action: its pattern matches it, its actions take the action, none of its
// exceptions matches it. A line of conditions does not count: what the
// request will be is not known.
func (l *Line) match(resource, action string) bool {
	if len(l.Conditions) > 0 || !l.glob().Match(resource) {
		return false
	}
	if !anyAction(l.Actions) && !containsString(l.Actions, action) {
		return false
	}
	if l.eg == nil && len(l.Except) > 0 {
		l.eg = NewExceptCondition(l.Except...).globs
	}
	for _, g := range l.eg {
		if g.Match(resource) {
			return false
		}
	}
	return true
}

// allows is the decision of the active lines for the resource asked with
// action: matched, some line counts for it; allowed, a line of no prefix
// counts and no deny does (denied, a deny counts).
func (s *Stack) allows(active []*Line, resource, action string) (allowed, denied bool) {
	for _, l := range active {
		if !l.match(resource, action) {
			continue
		}
		if l.Denies() {
			return false, true
		}
		allowed = true
	}
	return
}

// Allows says whether the stack, resolved, allows the permission asked —
// the resource with its verb at the end, as a Verifier asks it
// ("admin:content/:posts:@list") —, with action (the verb; the lines of the
// code may name theirs), and by the unique name preferred ("" none) first,
// as the Verifier decides: there a deny denies and an allow allows; when no
// line counts for it, the resource decides. A line of conditions does not
// count.
func (s *Stack) Allows(resource, preferred, action string) bool {
	active := s.Active()
	if preferred != "" {
		if allowed, denied := s.allows(active, preferred, action); allowed || denied {
			return allowed
		}
	}
	allowed, _ := s.allows(active, resource, action)
	return allowed
}

// CodeLines are the lines of the policies of the code — those not loaded from
// the database — whose subjects match one of subjects, in their order: an
// allow of no prefix, a deny "deny ", a line a resource; their origin
// "code: <their subjects>".
func (b *Builder) CodeLines(subjects []string) []*Line {
	return b.policyLines(subjects, false, "code")
}

// DBLines are the lines of the policies of the database (DBPolicy) whose
// subjects match one of subjects — what is granted a record, a share —, as
// CodeLines; their origin "database: <their subjects>".
func (b *Builder) DBLines(subjects []string) []*Line {
	return b.policyLines(subjects, true, "database")
}

func (b *Builder) policyLines(subjects []string, db bool, from string) []*Line {
	b.m.Lock()
	defer b.m.Unlock()
	var out []*Line
	for _, p := range b.policies {
		if b.dbIDs[p.GetID()] != db {
			continue
		}
		pol := p.policy
		if !subjectsMatch(pol.Subjects, subjects) {
			continue
		}
		kind := LineGrant
		if pol.Effect == Denied {
			kind = LineDeny
		}
		origin := fmt.Sprintf("%s: %s", from, strings.Join(pol.Subjects, ", "))
		for _, res := range pol.Resources {
			for _, e := range ExpandLine(res) {
				k, pat := ParseLine(e)
				if k == LineGrant {
					k = kind
				}
				out = append(out, &Line{Kind: k, Pattern: pat, Origin: origin, Actions: pol.Actions,
					Conditions: pol.Conditions})
			}
		}
	}
	return out
}

func subjectsMatch(patterns, subjects []string) bool {
	for _, p := range patterns {
		g, err := glob.Compile(p)
		if err != nil {
			continue
		}
		for _, s := range subjects {
			if g.Match(s) {
				return true
			}
		}
	}
	return false
}

// splitAlternatives are the alternatives of a group, "{a, b}", trimmed — of
// a text that is no group, the text itself —: the commas of the groups
// inside an alternative are its own.
func splitAlternatives(s string) (out []string) {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "{") || !strings.HasSuffix(s, "}") || !wholeGroup(s) {
		return []string{s}
	}
	s = s[1 : len(s)-1]
	depth, from := 0, 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\\':
			i++
		case '{', '[', '<':
			depth++
		case '}', ']', '>':
			if depth > 0 {
				depth--
			}
		case ',', '\n':
			if depth == 0 {
				if a := strings.TrimSpace(s[from:i]); a != "" {
					out = append(out, a)
				}
				from = i + 1
			}
		}
	}
	if a := strings.TrimSpace(s[from:]); a != "" {
		out = append(out, a)
	}
	return
}

// wholeGroup says whether s, "{…}", is one group: its first brace closed by
// its last.
func wholeGroup(s string) bool {
	depth := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\\':
			i++
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 && i != len(s)-1 {
				return false
			}
		}
	}
	return depth == 0
}

// PolicyLines are the lines a policy of before is, a resource each: "deny "
// before them when it denies; of actions other than any, each appended to
// the resource — the action ends it ("admin:posts:*" of @get is
// "admin:posts:*@get"). What the roles and the keys of before had is made
// lines so (once, by their packages).
func PolicyLines(effect string, actions, resources []string) (out []string) {
	prefix := ""
	if effect == Denied {
		prefix = PrefixDeny
	}
	for _, res := range resources {
		for _, r := range SplitResources(res) {
			if anyAction(actions) {
				out = append(out, prefix+r)
				continue
			}
			for _, a := range actions {
				out = append(out, prefix+r+a)
			}
		}
	}
	return
}

// Tree is lines — the active ones of a stack, whose order does not matter —
// as one, indented: their common beginnings once, by their segments (":"),
// what follows each as a group of alternatives, a line each, its prefix
// before it — "admin:users:list" and "admin:users:new" are
//
//	admin:users:{
//		list
//		new
//	}
//
// recursively. It is a line itself (a new line separates alternatives as a
// comma does): ExpandLine gives the lines back.
func Tree(lines []string) string {
	root := &treeNode{}
	for _, l := range lines {
		k, p := ParseLine(l)
		if p == "" {
			continue
		}
		root.add(splitSegments(p), k.Prefix())
	}
	var b strings.Builder
	for i, c := range root.children {
		if i > 0 {
			b.WriteString("\n")
		}
		c.write(&b, "", "")
	}
	return b.String()
}

type treeNode struct {
	segment  string
	children []*treeNode
	// leaf: a line ends here, of its prefix
	leaf   bool
	prefix string
}

func (n *treeNode) child(seg string) *treeNode {
	for _, c := range n.children {
		if c.segment == seg {
			return c
		}
	}
	c := &treeNode{segment: seg}
	n.children = append(n.children, c)
	return c
}

func (n *treeNode) add(segs []string, prefix string) {
	if len(segs) == 0 {
		return
	}
	c := n.child(segs[0])
	if len(segs) == 1 {
		c.leaf, c.prefix = true, prefix
		return
	}
	c.add(segs[1:], prefix)
}

// write writes the node — after what came before it, of the indent — and
// what follows it: one way, its segments joined; more, a group.
func (n *treeNode) write(b *strings.Builder, indent, prefix string) {
	if n.leaf && len(n.children) == 0 {
		b.WriteString(indent + prefix + n.prefix + n.segment)
		return
	}
	if !n.leaf && len(n.children) == 1 {
		n.children[0].write(b, indent, prefix+n.segment)
		return
	}
	b.WriteString(indent + prefix + n.segment + "{\n")
	in := indent + "\t"
	first := true
	if n.leaf {
		// a line ends here, others go on: the empty alternative
		b.WriteString(in + n.prefix)
		first = false
	}
	for _, c := range n.children {
		if !first {
			b.WriteString("\n")
		}
		first = false
		c.write(b, in, "")
	}
	b.WriteString("\n" + indent + "}")
}

// splitSegments are the segments of a pattern, each with its ":" at its end
// (but the last): "admin:users:@list" is "admin:", "users:", "@list". The
// colons inside a group or a record ("<…>") are their own.
func splitSegments(p string) (out []string) {
	depth, from := 0, 0
	for i := 0; i < len(p); i++ {
		switch p[i] {
		case '\\':
			i++
		case '{', '<', '[':
			depth++
		case '}', '>', ']':
			if depth > 0 {
				depth--
			}
		case ':':
			if depth == 0 {
				out = append(out, p[from:i+1])
				from = i + 1
			}
		}
	}
	if from < len(p) {
		out = append(out, p[from:])
	}
	return
}
