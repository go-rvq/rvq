package perm

import (
	"strings"

	"github.com/gobwas/glob"
)

// A line and its lines: a field that writes permissions edits the lines a
// line is (Lines), as a list — one put, one taken out, a prefix changed —,
// and makes the line of them again (JoinLines), canonical; it stores the
// items of its root, one a scope (SplitRoot). What is given to a resource
// is stored by its unique name, never by the groups of the menu it is in
// (CleanLine): a resource moved to another group keeps its
// permissions.

// Lines are the lines the line a is, its groups broken in their order (as
// ExpandLine), the labels of its records kept: "deny admin:posts:{<7<A>>:@get,
// allow @list}" is "deny admin:posts:<7<A>>:@get" and "allow
// admin:posts:@list".
func (a *ExprAlt) Lines() (out []string) {
	for _, l := range expandAlt(a.Kind, a.Nodes) {
		out = append(out, l.kind.Prefix()+NodesString(l.nodes))
	}
	return
}

type expandedNodes struct {
	kind  LineKind
	nodes []ExprNode
}

// expandAlt breaks the first group of nodes — each alternative in its place,
// of its prefix when it has one — and, recursively, the rest.
func expandAlt(kind LineKind, nodes []ExprNode) (out []expandedNodes) {
	for i, n := range nodes {
		g, ok := n.(*ExprGroup)
		if !ok {
			continue
		}
		for _, c := range g.Alts {
			k := kind
			if c.Kind != LineGrant {
				k = c.Kind
			}
			rest := append(append(append([]ExprNode{}, nodes[:i]...), c.Nodes...), nodes[i+1:]...)
			out = append(out, expandAlt(k, rest)...)
		}
		return
	}
	if len(nodes) > 0 {
		out = append(out, expandedNodes{kind, nodes})
	}
	return
}

// Merge is the line of a and others merged, canonical (Normalize): their
// lines (Lines) in their order, a line of a resource an earlier one has
// taking its place — the earlier out, the later after, as a later line of a
// stack decides —; an unset ("unset admin:posts:*") taking out the lines
// before it its pattern covers, whatever their prefix, and gone; the labels
// of the records kept. a is not changed.
func (a *ExprAlt) Merge(other ...*ExprAlt) *ExprAlt {
	m, _ := a.MergeReport(other...)
	return m
}

// MergeReport is Merge, and the unsets that took nothing out (idle): a line
// wider than the pattern ("admin:*" of "unset admin:posts:*") is not taken
// out — a deny is what narrows it.
func (a *ExprAlt) MergeReport(other ...*ExprAlt) (merged *ExprAlt, idle []string) {
	type line struct{ text, key string }
	var out []line
	for _, e := range append([]*ExprAlt{a}, other...) {
		if e == nil {
			continue
		}
		for _, l := range e.Lines() {
			k, p := ParseLine(l)
			key := StripLabels(p)
			kept := out[:0]
			if k == LineUnset {
				g, err := glob.Compile(key)
				for _, o := range out {
					if o.key == key || err == nil && g.Match(o.key) {
						continue
					}
					kept = append(kept, o)
				}
				if len(kept) == len(out) {
					idle = append(idle, l)
				}
				out = kept
				continue
			}
			for _, o := range out {
				if o.key != key {
					kept = append(kept, o)
				}
			}
			out = append(kept, line{l, key})
		}
	}
	g := &ExprGroup{}
	for _, o := range out {
		if e, err := ParseExpr(o.text); err == nil {
			g.Alts = append(g.Alts, e)
		}
	}
	root := &ExprAlt{}
	if len(g.Alts) > 0 {
		root.Nodes = []ExprNode{g}
	}
	return root.Normalize(), idle
}

// NodesString is the nodes as written, the labels of their records with
// them.
func NodesString(nodes []ExprNode) string {
	var b strings.Builder
	writeNodes(NewWriteContext(&b), nodes)
	return b.String()
}

// LinesOf are the lines of the line (Lines), canonical first; a line that
// cannot be read, its error.
func LinesOf(line string) ([]string, error) {
	if strings.TrimSpace(line) == "" {
		return nil, nil
	}
	e, err := ParseExpr(line)
	if err != nil {
		return nil, err
	}
	return e.Normalize().Lines(), nil
}

// JoinLines is the line of the lines, canonical: "{l1, l2, …}" normalized —
// the items of its root the scopes, one each. No line, "".
func JoinLines(lines []string) (string, error) {
	var ls []string
	for _, l := range lines {
		if l = strings.TrimSpace(l); l != "" {
			ls = append(ls, l)
		}
	}
	switch len(ls) {
	case 0:
		return "", nil
	case 1:
		return Canonical(ls[0])
	}
	return Canonical("{" + strings.Join(ls, ", ") + "}")
}

// SplitRoot are the items of the root of the line, canonical, a line each —
// one a scope: what a stack stores. A blank line, none.
func SplitRoot(line string) ([]string, error) {
	if strings.TrimSpace(line) == "" {
		return nil, nil
	}
	e, err := ParseExpr(line)
	if err != nil {
		return nil, err
	}
	e = e.Normalize()
	if len(e.Nodes) == 0 {
		return nil, nil
	}
	if g, ok := e.Nodes[0].(*ExprGroup); ok && len(e.Nodes) == 1 {
		var out []string
		for _, c := range g.Alts {
			k := c.Kind
			if k == LineGrant {
				k = e.Kind
			}
			out = append(out, k.Prefix()+NodesString(c.Nodes))
		}
		return out, nil
	}
	return []string{e.String()}, nil
}

// splitTop are the segments of a pattern, each with its ":" at its end (but
// the last): the colons inside a record ("<…<…>>"), brackets or an escape
// are their own.
func splitTop(p string) (out []string) {
	depth, brackets, from := 0, 0, 0
	for i := 0; i < len(p); i++ {
		switch p[i] {
		case '\\':
			i++
		case '<':
			depth++
		case '>':
			if depth > 0 {
				depth--
			}
		case '[':
			brackets++
		case ']':
			if brackets > 0 {
				brackets--
			}
		case ':':
			if depth == 0 && brackets == 0 {
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

// isGroupSegment says the segment is a group of the menu, literal:
// "content/:".
func isGroupSegment(s string) bool {
	return strings.HasSuffix(s, "/:") && literal(s)
}

// UngroupPattern is the pattern of one line (no group of alternatives), what
// it gives to a resource by its unique name: the groups of the menu after
// its scope left out when a resource — a segment not a group, of no glob —
// follows them ("admin:content/:posts:*" is "admin:posts:*",
// "admin:site/:/report:@get" "admin:/report:@get"); given to a group
// ("admin:content/:*", "admin:content/:sub/:"), as it is.
func UngroupPattern(p string) string {
	segs := splitTop(p)
	if len(segs) < 3 {
		return p
	}
	i := 1
	for i < len(segs) && isGroupSegment(segs[i]) {
		i++
	}
	if i == 1 || i >= len(segs) {
		return p
	}
	next := strings.TrimSuffix(segs[i], ":")
	// a resource: a model, a page — not a verb, an action, a field, a record,
	// a set (given to the group, as written)
	if next == "" || !literal(next) || strings.ContainsAny(next[:1], "<@!#$&%") {
		return p
	}
	return segs[0] + strings.Join(segs[i:], "")
}

// DropPlaceholders are the lines without the placeholders that others are
// under: a line ending by ":" ("deny admin:", "admin:posts:#Body:") decides
// nothing — it only holds the place of what will follow it —, and once a
// line is under it ("admin:posts:#Body:*"), it is out. A placeholder alone
// is kept: what follows it is still missing.
func DropPlaceholders(lines []string) []string {
	keys := make([]string, len(lines))
	for i, l := range lines {
		_, p := ParseLine(l)
		keys[i] = StripLabels(p)
	}
	var out []string
lines:
	for i, l := range lines {
		if strings.HasSuffix(keys[i], ":") {
			for j, k := range keys {
				if j != i && len(k) > len(keys[i]) && strings.HasPrefix(k, keys[i]) {
					continue lines
				}
			}
		}
		out = append(out, l)
	}
	return out
}

// DropIncomplete are the lines that decide something: none of the
// placeholders ("admin:users:", "deny admin:"), which name a part but not
// what of it ("admin:users:*" is all of it).
func DropIncomplete(lines []string) (out []string) {
	for _, l := range lines {
		if _, p := ParseLine(l); !strings.HasSuffix(StripLabels(p), ":") {
			out = append(out, l)
		}
	}
	return
}

// StoreLines are the lines of the value a field of permissions stores: its
// lines (LinesOf) out of the groups of the menu (UngroupPattern), without any
// placeholder — a line ending by ":" ("admin:users:") names a part but not
// what of it, and decides nothing (DropIncomplete) —, in their order
// (SortValue), the items of their root, one a scope (SplitRoot). A value it cannot read, its error.
func StoreLines(value string) ([]string, error) {
	ls, err := LinesOf(value)
	if err != nil {
		return nil, err
	}
	for i, l := range ls {
		k, p := ParseLine(l)
		ls[i] = k.Prefix() + UngroupPattern(p)
	}
	v, err := JoinLines(DropIncomplete(ls))
	if err != nil {
		return nil, err
	}
	// in their order (SortValue)
	if v, err = SortValue(v); err != nil {
		return nil, err
	}
	return SplitRoot(v)
}

// CleanLine is the line as stored, and whether it changed: what it gives to
// resources by their unique names (UngroupPattern), the permissions renamed
// (the trash's "!restore" is the record's "<*>:@restore"), none of its
// placeholders (DropIncomplete) — "" when nothing is left —; a line that
// needs no change, or cannot be read, as it is.
func CleanLine(line string) (string, bool) {
	ls, err := LinesOf(line)
	if err != nil || len(ls) == 0 {
		return line, false
	}
	changed := false
	var out []string
	for _, l := range ls {
		k, p := ParseLine(l)
		if u := renameVerbs(UngroupPattern(p)); u != p {
			p, l, changed = u, k.Prefix()+u, true
		}
		if strings.HasSuffix(StripLabels(p), ":") {
			changed = true
			continue
		}
		out = append(out, l)
	}
	if !changed {
		return line, false
	}
	v, err := JoinLines(out)
	if err != nil {
		return line, false
	}
	return v, true
}

// renameVerbs is the pattern of one line with the permissions renamed:
// the restore of the trash of before — an action of the listing,
// "admin:posts:!restore" — is the record's "admin:posts:<*>:@restore".
func renameVerbs(p string) string {
	if strings.HasSuffix(p, ":!restore") {
		return strings.TrimSuffix(p, "!restore") + "<*>:@restore"
	}
	return p
}

// CleanLines are lines stored (a stack's) cleaned (CleanLine) — those left
// empty out —, and whether any changed.
func CleanLines(lines []string) ([]string, bool) {
	out := make([]string, 0, len(lines))
	changed := false
	for _, l := range lines {
		c, ch := CleanLine(l)
		changed = changed || ch
		if strings.TrimSpace(c) != "" {
			out = append(out, c)
		}
	}
	return out, changed
}
