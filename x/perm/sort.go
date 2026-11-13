package perm

import "strings"

// The order of the lines of a value (SortLines): the items of one level from
// the general to the particular, as the tree of the permissions — "*", the
// sets, the verbs (as the CRUD), the actions, the records, the fields, the
// sections, the records edited in place, the pages, the models and groups —,
// and their order deciding the same. Of grants and denies only, a value
// decides the same in any order — what its grants allow, but what its denies
// deny: a later grant takes nothing from a deny, a deny wins. An allow and a
// reset decide by what is before them: in a value of one (a line of a set
// that gives one, too), two lines change places only when they are of
// resources apart (independent) — then neither replaces, covers nor takes
// from the other —; any other two keep their order.

// verbOrder is the order of the verbs: as the CRUD, the trash after; the
// others after them, by name.
var verbOrder = map[string]int{"list": 0, "get": 1, "create": 2, "edit": 3, "delete": 4, "trash": 5, "restore": 6}

// The ranks of the segments of a level (segmentRank).
const (
	rankAll = iota
	rankSet
	rankVerb
	rankAction
	rankAllRecords
	rankRecord
	rankField
	rankSection
	rankInline
	rankPage
	rankOther
)

// segmentRank is the rank of a segment in its level, and what orders it in
// its rank.
func segmentRank(seg string) (rank int, sub int, name string) {
	s := strings.TrimSuffix(seg, ":")
	switch {
	case s == "*":
		return rankAll, 0, ""
	case strings.HasPrefix(s, SetPrefix):
		return rankSet, 0, s[1:]
	case strings.HasPrefix(s, "@"):
		v := s[1:]
		if i, ok := verbOrder[v]; ok {
			return rankVerb, i, v
		}
		return rankVerb, len(verbOrder), v
	case strings.HasPrefix(s, "!"):
		return rankAction, 0, s[1:]
	case s == "<*>":
		return rankAllRecords, 0, ""
	case strings.HasPrefix(s, "<"):
		return rankRecord, 0, StripLabels(s)
	case strings.HasPrefix(s, "#"):
		return rankField, 0, s[1:]
	case strings.HasPrefix(s, "$"):
		return rankSection, 0, s[1:]
	case strings.HasPrefix(s, "&"):
		return rankInline, 0, s[1:]
	case strings.HasPrefix(s, "/"):
		return rankPage, 0, s
	}
	return rankOther, 0, s
}

// compareSegments orders two segments of the same level.
func compareSegments(a, b string) int {
	ra, sa, na := segmentRank(a)
	rb, sb, nb := segmentRank(b)
	switch {
	case ra != rb:
		return ra - rb
	case sa != sb:
		return sa - sb
	case na != nb:
		return strings.Compare(na, nb)
	}
	return strings.Compare(a, b)
}

// ComparePatterns orders two patterns as SortLines does: by their first
// segment that differs (a pattern that is the beginning of the other,
// first).
func ComparePatterns(a, b string) int {
	sa, sb := splitTop(StripLabels(a)), splitTop(StripLabels(b))
	for i := 0; i < len(sa) && i < len(sb); i++ {
		if c := compareSegments(sa[i], sb[i]); c != 0 {
			return c
		}
	}
	return len(sa) - len(sb)
}

// globChars are what makes a segment a glob.
const globChars = "*?[{"

// sortKey is the pattern a line is ordered by, and the one told apart from
// the others: a line of a set, all under its resource ("admin:posts:*").
func sortKey(l string) (kind LineKind, order, span string) {
	kind, p := ParseLine(l)
	p = StripLabels(p)
	if res, _, ok := SetRef(p); ok {
		return kind, p, res + "*"
	}
	return kind, p, p
}

// segmentClass is what a segment may match: the segments of its sigil ("@"
// verbs, "!" actions, "%" sets, "<" records, "#" fields, "$" sections, "&"
// records edited in place, "/" pages; else models and groups), or any — a
// segment that begins by a glob ("*", "{a,b}").
func segmentClass(seg string) byte {
	if seg == "" {
		return 'o'
	}
	switch c := seg[0]; {
	case strings.IndexByte(globChars, c) >= 0:
		return '*'
	case strings.IndexByte("@!%<#$&/", c) >= 0:
		return c
	}
	return 'o'
}

// independent says the patterns a and b are of resources apart: at a
// segment, of sigils that match each other's never, or both of no glob and
// others.
func independent(a, b string) bool {
	sa, sb := splitTop(a), splitTop(b)
	for i := 0; i < len(sa) && i < len(sb); i++ {
		ca, cb := segmentClass(sa[i]), segmentClass(sb[i])
		if ca == '*' || cb == '*' {
			return false
		}
		if ca != cb {
			return true
		}
		ga, gb := strings.ContainsAny(sa[i], globChars), strings.ContainsAny(sb[i], globChars)
		if ga || gb {
			return false
		}
		if sa[i] != sb[i] {
			return true
		}
	}
	return false
}

// decidesByOrder says a line of lines decides by what is before it: an
// allow, a reset, a line of a set that gives one of them.
func decidesByOrder(lines []string) bool {
	for _, l := range lines {
		k, p := ParseLine(l)
		if k == LineAllow || k == LineReset {
			return true
		}
		if res, name, ok := SetRef(p); ok {
			given, _ := ExpandSet(res, name)
			for _, g := range given {
				if gk, _ := ParseLine(g); gk == LineAllow || gk == LineReset {
					return true
				}
			}
		}
	}
	return false
}

// SortLines are the lines in their order (the order of the lines of a
// value): by their patterns (ComparePatterns; a grant before another of the
// same); in a value that decides by its order (an allow, a reset), each time
// the first of those all the lines before of which they are not apart
// (independent) are out — any two that decide by their order kept in it.
// lines is not changed.
func SortLines(lines []string) []string {
	type item struct {
		line        string
		kind        LineKind
		order, span string
		done        bool
	}
	items := make([]*item, len(lines))
	for i, l := range lines {
		k, o, sp := sortKey(l)
		items[i] = &item{line: l, kind: k, order: o, span: sp}
	}
	less := func(a, b *item) bool {
		if c := ComparePatterns(a.order, b.order); c != 0 {
			return c < 0
		}
		return a.kind == LineGrant && b.kind != LineGrant
	}
	ordered := decidesByOrder(lines)
	out := make([]string, 0, len(lines))
	for len(out) < len(items) {
		var best *item
		for i, it := range items {
			if it.done {
				continue
			}
			free := true
			for _, b := range items[:i] {
				if ordered && !b.done && !independent(b.span, it.span) {
					free = false
					break
				}
			}
			if free && (best == nil || less(it, best)) {
				best = it
			}
		}
		best.done = true
		out = append(out, best.line)
	}
	return out
}

// SortValue is the value of its lines in their order (SortLines), canonical
// (JoinLines) — twice, the same; a value it cannot read, its error.
func SortValue(value string) (string, error) {
	// the canonical form groups what it may (Normalize), which may give
	// another order to sort: again, until it is the same
	for i := 0; i < 8; i++ {
		lines, err := LinesOf(value)
		if err != nil {
			return "", err
		}
		sorted, err := JoinLines(SortLines(lines))
		if err != nil || sorted == value {
			return sorted, err
		}
		value = sorted
	}
	return value, nil
}

// WithSort writes the line of its lines in their order (SortValue): of
// FormatLine.
func WithSort() WriteOption {
	return func(ctx *WriteContext) { ctx.sort = true }
}
