package perm

import "strings"

// Normalize is the line in its canonical form — what a field that edits it
// writes —, deciding the same:
//
//   - what follows a group goes into each alternative: "{a, b}:@get" is
//     "{a:@get, b:@get}";
//   - a group in an alternative of nothing else is its alternatives:
//     "{a, {b, c}}" is "{a, b, c}";
//   - the alternatives that begin by the same segment are one, a group after
//     it — "{posts:@list, posts:@get}" is "posts:{@list, @get}" —, when
//     moving them together decides the same: they and those between them
//     only allow or deny (no allow, no reset: their order then decides
//     nothing), or the segments of those between are others, of no globs;
//     at the root, the scopes: one item each;
//   - a group of one alternative is it: "admin:{posts:*}" is
//     "admin:posts:*", "admin:{deny posts:*}" "deny admin:posts:*";
//   - the prefix of every alternative of a group, before it:
//     "admin:{deny a, deny b}" is "deny admin:{a, b}";
//   - a group of none is nothing: no line, as ExpandLine.
func (a *ExprAlt) Normalize() *ExprAlt {
	return normAlt(cloneAlt(a), LineGrant, true)
}

// Canonical is the line's canonical form (Normalize): "" when it cannot be
// read (its error).
func Canonical(line string) (string, error) {
	e, err := ParseExpr(line)
	if err != nil {
		return "", err
	}
	return e.Normalize().String(), nil
}

func cloneAlt(a *ExprAlt) *ExprAlt {
	c := &ExprAlt{Pos: a.Pos, Kind: a.Kind, Nodes: make([]ExprNode, len(a.Nodes))}
	for i, n := range a.Nodes {
		c.Nodes[i] = cloneNode(n)
	}
	return c
}

func cloneNode(n ExprNode) ExprNode {
	switch n := n.(type) {
	case *ExprText:
		c := *n
		return &c
	case *ExprRecord:
		c := *n
		return &c
	case *ExprGroup:
		c := &ExprGroup{Pos: n.Pos, Alts: make([]*ExprAlt, len(n.Alts))}
		for i, a := range n.Alts {
			c.Alts[i] = cloneAlt(a)
		}
		return c
	}
	return n
}

func cloneNodes(nodes []ExprNode) []ExprNode {
	out := make([]ExprNode, len(nodes))
	for i, n := range nodes {
		out[i] = cloneNode(n)
	}
	return out
}

// effective is the kind a line of the alternative of kind k is, in a line of
// kind outer: its own, else the outer's.
func effective(k, outer LineKind) LineKind {
	if k != LineGrant {
		return k
	}
	return outer
}

// joinTexts joins the texts side by side.
func joinTexts(nodes []ExprNode) []ExprNode {
	var out []ExprNode
	for _, n := range nodes {
		if t, ok := n.(*ExprText); ok {
			if t.Text == "" {
				continue
			}
			if len(out) > 0 {
				if p, ok := out[len(out)-1].(*ExprText); ok {
					out[len(out)-1] = &ExprText{Pos: Pos{p.Start, t.End}, Text: p.Text + t.Text}
					continue
				}
			}
		}
		out = append(out, n)
	}
	return out
}

// normAlt is the alternative a, in a line of kind outer, normalized; root:
// the line itself (its group's alternatives, the scopes).
func normAlt(a *ExprAlt, outer LineKind, root bool) *ExprAlt {
	nodes := joinTexts(a.Nodes)
	// what follows the first group, into each of its alternatives
	for i, n := range nodes {
		g, ok := n.(*ExprGroup)
		if !ok || i == len(nodes)-1 {
			continue
		}
		suffix := nodes[i+1:]
		for _, c := range g.Alts {
			c.Nodes = append(c.Nodes, cloneNodes(suffix)...)
		}
		nodes = nodes[:i+1]
		break
	}
	a.Nodes = nodes
	if len(nodes) == 0 {
		return a
	}
	g, ok := nodes[len(nodes)-1].(*ExprGroup)
	if !ok {
		return a
	}
	kind := effective(a.Kind, outer)
	g.Alts = normAlts(g.Alts, kind, root && len(nodes) == 1)
	// the same prefix of every alternative: the alternative's — each line of
	// it of that prefix still
	if k := commonKind(g.Alts); k != LineGrant && len(g.Alts) > 1 {
		for _, c := range g.Alts {
			c.Kind = LineGrant
		}
		a.Kind = k
	}
	switch len(g.Alts) {
	case 0:
		// a group of none: no line
		a.Kind, a.Nodes = LineGrant, nil
	case 1:
		c := g.Alts[0]
		if c.Kind != LineGrant {
			a.Kind = c.Kind
		}
		a.Nodes = joinTexts(append(nodes[:len(nodes)-1], c.Nodes...))
	}
	return a
}

// normAlts are the alternatives of a group, in a line of kind outer,
// normalized: each, the groups of nothing else made theirs, those of the same
// first segment one.
func normAlts(alts []*ExprAlt, outer LineKind, root bool) []*ExprAlt {
	var flat []*ExprAlt
	for _, c := range alts {
		c = normAlt(c, outer, false)
		if len(c.Nodes) == 0 {
			continue
		}
		if g, ok := c.Nodes[0].(*ExprGroup); ok && len(c.Nodes) == 1 {
			// a group of nothing else: its alternatives, of its kind when
			// they have none
			for _, gc := range g.Alts {
				if gc.Kind == LineGrant {
					gc.Kind = c.Kind
				}
				flat = append(flat, gc)
			}
			continue
		}
		flat = append(flat, c)
	}
	return mergeAlts(flat, outer, root)
}

// head is the first segment of the alternative a, up to its first ":" (with
// it) — texts and records, no group —, and what follows it; ok false when it
// has none (no ":", a group first) or nothing follows it.
func head(a *ExprAlt) (key string, h, rest []ExprNode, ok bool) {
	var b strings.Builder
	for i, n := range a.Nodes {
		switch n := n.(type) {
		case *ExprGroup:
			return "", nil, nil, false
		case *ExprRecord:
			b.WriteString("<" + n.ID + ">")
		case *ExprText:
			j := segmentEnd(n.Text)
			if j < 0 {
				b.WriteString(n.Text)
				continue
			}
			b.WriteString(n.Text[:j+1])
			h = append(cloneNodes(a.Nodes[:i]), &ExprText{Pos: Pos{n.Start, n.Start + j + 1}, Text: n.Text[:j+1]})
			if j+1 < len(n.Text) {
				rest = append(rest, &ExprText{Pos: Pos{n.Start + j + 1, n.End}, Text: n.Text[j+1:]})
			}
			rest = append(rest, a.Nodes[i+1:]...)
			if len(rest) == 0 {
				return "", nil, nil, false
			}
			return b.String(), h, rest, true
		}
	}
	return "", nil, nil, false
}

// segmentEnd is where the first ":" of a text is, out of brackets and
// escapes; -1 when none.
func segmentEnd(s string) int {
	depth := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\\':
			i++
		case '[':
			depth++
		case ']':
			if depth > 0 {
				depth--
			}
		case ':':
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// literal says the segment has no glob: what it matches is itself.
func literal(key string) bool {
	return key != "" && !strings.ContainsAny(key, "*?[{")
}

// mergeAlts makes the alternatives of the same first segment one, a group
// after it, when moving them together decides the same (Normalize); at the
// root, always (the scopes).
func mergeAlts(alts []*ExprAlt, outer LineKind, root bool) []*ExprAlt {
	type entry struct {
		alt  *ExprAlt
		key  string
		h    []ExprNode
		rest [][]ExprNode
		kind []LineKind
		ok   bool
	}
	var out []*entry
	for _, c := range alts {
		key, h, rest, ok := head(c)
		if ok {
			for i := len(out) - 1; i >= 0; i-- {
				e := out[i]
				if !e.ok || e.key != key {
					continue
				}
				// those after it, before c: c goes before them
				safe := root
				if !safe {
					safe = orderFree(c, outer)
					for _, o := range out[i+1:] {
						safe = safe && orderFree(o.alt, outer)
						for j, r := range o.rest[1:] {
							safe = safe && orderFree(&ExprAlt{Kind: o.kind[j+1], Nodes: r}, outer)
						}
					}
				}
				if !safe {
					safe = literal(key)
					for _, o := range out[i+1:] {
						if !o.ok || !literal(o.key) {
							safe = false
							break
						}
					}
				}
				if safe {
					e.rest = append(e.rest, rest)
					e.kind = append(e.kind, c.Kind)
					c = nil
				}
				break
			}
		}
		if c != nil {
			out = append(out, &entry{alt: c, key: key, h: h, rest: [][]ExprNode{rest}, kind: []LineKind{c.Kind}, ok: ok})
		}
	}
	res := make([]*ExprAlt, 0, len(out))
	for _, e := range out {
		if len(e.rest) == 1 {
			res = append(res, e.alt)
			continue
		}
		g := &ExprGroup{Pos: e.alt.Pos}
		for i, r := range e.rest {
			g.Alts = append(g.Alts, &ExprAlt{Pos: e.alt.Pos, Kind: e.kind[i], Nodes: r})
		}
		m := &ExprAlt{Pos: e.alt.Pos, Nodes: append(e.h, g)}
		res = append(res, normAlt(m, outer, false))
	}
	return res
}

// orderFree says the lines of the alternative a, in a line of kind outer,
// only allow or deny — no allow, no reset —: where they are among others
// decides nothing.
func orderFree(a *ExprAlt, outer LineKind) bool {
	k := effective(a.Kind, outer)
	if k != LineGrant && k != LineDeny {
		return false
	}
	for _, n := range a.Nodes {
		if g, ok := n.(*ExprGroup); ok {
			for _, c := range g.Alts {
				if !orderFree(c, k) {
					return false
				}
			}
		}
	}
	return true
}

// commonKind is the prefix every alternative has (LineGrant: none, or not
// the same).
func commonKind(alts []*ExprAlt) LineKind {
	if len(alts) == 0 {
		return LineGrant
	}
	k := alts[0].Kind
	for _, c := range alts[1:] {
		if c.Kind != k {
			return LineGrant
		}
	}
	return k
}
