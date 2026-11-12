package perm

import (
	"fmt"
	"strings"
)

// The grammar of a line of a stack, as a tree (ParseExpr): what a field that
// edits a line (vuetifyx's vx-permission-input) reads and writes.
//
//	line   = [prefix] nodes | alt {("," | "\n") alt}   several at the root: a group
//	alt    = [prefix] nodes
//	prefix = ("deny" | "allow" | "reset" | "unset") " "     the word alone: a resource
//	nodes  = {text | record | group}
//	record = "<" id [ "<" label ">" ] ">"           <7>, <7<First post>>
//	group  = "{" alt {("," | "\n") alt} "}"          after a ":", spaces before it: "admin: {"
//
// A record's label only shows the record: the decision never sees it
// (StripLabels) — "<7<First post>>" decides as "<7>". In a label, "<", ">" and
// "\" are escaped by "\"; what is between "<" and ">" is the record's, so the
// ":", "," and "{}" of a title break nothing.
//
// Several alternatives at the root, by commas or lines, are a group:
// "admin:content/:*, site:*" is "{admin:content/:*, site:*}". The items of the
// root are the scopes ("admin", "site"…): one each — Normalize makes
// "{admin:a, admin:b}" "admin:{a, b}".

// Pos is where a part of a line is: from Start to End (bytes).
type Pos struct {
	Start, End int
}

// ExprNode is a part of an alternative: an *ExprText, an *ExprRecord, an
// *ExprGroup.
type ExprNode interface {
	Position() Pos
	// WriteCtx writes the node as ctx says (WriteContext)
	WriteCtx(ctx *WriteContext)
}

// ExprText is text of a resource: its segments, ":" ending each, and globs.
type ExprText struct {
	Pos
	Text string
}

// ExprRecord is a record, "<7>" — its label, "<7<First post>>", what shows
// it only.
type ExprRecord struct {
	Pos
	ID, Label string
}

// ExprGroup is a group of alternatives, "{a, deny b}": each a line in the
// group's place, in its order.
type ExprGroup struct {
	Pos
	Alts []*ExprAlt
}

// ExprAlt is a line, or an alternative of a group: its prefix (Kind; LineGrant
// when it has none — of a group's alternative: the line's) and its parts.
type ExprAlt struct {
	Pos
	Kind  LineKind
	Nodes []ExprNode
}

// Position is where the part is.
func (p Pos) Position() Pos { return p }

// ExprError is what is wrong in a line and where (Pos: the byte).
type ExprError struct {
	Pos int
	Msg string
}

func (e *ExprError) Error() string { return fmt.Sprintf("at %d: %s", e.Pos+1, e.Msg) }

// The messages of ExprError.
const (
	ErrExprGroupOpen    = "a group not closed: a } is missing"
	ErrExprGroupClose   = "a } of no group"
	ErrExprRecordOpen   = "a record not closed: a > is missing"
	ErrExprRecordClose  = "a > of no record"
	ErrExprRecordID     = "a record of no id"
	ErrExprPrefixAlone  = "a prefix of no resource"
	ErrExprEscapeAtEnd  = "a \\ at the end"
	ErrExprLabelAtEnd   = "a label not closed: a > is missing"
	ErrExprBracketsOpen = "a [ not closed: a ] is missing"
)

// ParseExpr is the line as a tree. A blank line is an alternative of no
// parts.
func ParseExpr(line string) (*ExprAlt, error) {
	p := &exprParser{s: line}
	root, err := p.alts(-1)
	if err != nil {
		return nil, err
	}
	if len(root) == 1 {
		return root[0], nil
	}
	a := &ExprAlt{Pos: Pos{0, len(line)}}
	if len(root) > 1 {
		a.Nodes = []ExprNode{&ExprGroup{Pos: Pos{0, len(line)}, Alts: root}}
	}
	return a, nil
}

type exprParser struct {
	s string
	i int
}

// alts are the alternatives up to the end of the group that begins at group
// (its "}", read; -1: none) or of the line, empty ones left out.
func (p *exprParser) alts(group int) (out []*ExprAlt, err error) {
	for {
		a, end, err := p.alt(group)
		if err != nil {
			return nil, err
		}
		if a != nil {
			out = append(out, a)
		}
		switch end {
		case 0: // the end of the line
			return out, nil
		case '}':
			return out, nil
		}
	}
}

// alt is an alternative, up to "," or a new line (end: the byte), "}" of
// its group (that begins at group; -1: none) or the end of the line (end
// 0); nil when empty.
func (p *exprParser) alt(group int) (a *ExprAlt, end byte, err error) {
	s := p.s
	for p.i < len(s) && (s[p.i] == ' ' || s[p.i] == '\t' || s[p.i] == '\r') {
		p.i++
	}
	a = &ExprAlt{Pos: Pos{Start: p.i}}
	for _, k := range LineKinds {
		w := strings.TrimSpace(k.Prefix())
		if strings.HasPrefix(s[p.i:], w) {
			rest := s[p.i+len(w):]
			if strings.HasPrefix(rest, " ") || strings.HasPrefix(rest, "\t") {
				a.Kind = k
				p.i += len(w)
				for p.i < len(s) && (s[p.i] == ' ' || s[p.i] == '\t') {
					p.i++
				}
				break
			}
			// the word alone is no prefix: a resource ("!{edit, reset}")
		}
	}
	var text strings.Builder
	textStart := -1
	flush := func(at int) {
		if textStart >= 0 {
			a.Nodes = append(a.Nodes, &ExprText{Pos: Pos{textStart, at}, Text: text.String()})
			text.Reset()
			textStart = -1
		}
	}
	addText := func(t string) {
		if textStart < 0 {
			textStart = p.i
		}
		text.WriteString(t)
	}
	finish := func(at int) (*ExprAlt, error) {
		flush(at)
		// the spaces at its end are none of it
		if n := len(a.Nodes); n > 0 {
			if t, ok := a.Nodes[n-1].(*ExprText); ok {
				t.Text = strings.TrimRight(t.Text, " \t\r")
				t.End = t.Start + len(t.Text)
				if t.Text == "" {
					a.Nodes = a.Nodes[:n-1]
				}
			}
		}
		a.End = at
		if len(a.Nodes) == 0 {
			if a.Kind != LineGrant {
				return nil, &ExprError{Pos: a.Start, Msg: ErrExprPrefixAlone}
			}
			return nil, nil
		}
		return a, nil
	}
	for p.i < len(s) {
		c := s[p.i]
		switch c {
		case '\\':
			if p.i+1 >= len(s) {
				return nil, 0, &ExprError{Pos: p.i, Msg: ErrExprEscapeAtEnd}
			}
			addText(s[p.i : p.i+2])
			p.i += 2
		case '[':
			j := strings.IndexByte(s[p.i:], ']')
			if j < 0 {
				return nil, 0, &ExprError{Pos: p.i, Msg: ErrExprBracketsOpen}
			}
			addText(s[p.i : p.i+j+1])
			p.i += j + 1
		case '{':
			// the spaces between a ":" and its group are none of it: "admin: {"
			if t := strings.TrimRight(text.String(), " \t"); t != text.String() && strings.HasSuffix(t, ":") {
				text.Reset()
				text.WriteString(t)
			}
			flush(p.i)
			start := p.i
			p.i++
			alts, err := p.alts(start)
			if err != nil {
				return nil, 0, err
			}
			a.Nodes = append(a.Nodes, &ExprGroup{Pos: Pos{start, p.i}, Alts: alts})
		case '}':
			if group < 0 {
				return nil, 0, &ExprError{Pos: p.i, Msg: ErrExprGroupClose}
			}
			at := p.i
			p.i++
			r, err := finish(at)
			return r, '}', err
		case ',', '\n':
			at := p.i
			p.i++
			r, err := finish(at)
			return r, c, err
		case '<':
			flush(p.i)
			r, err := p.record()
			if err != nil {
				return nil, 0, err
			}
			a.Nodes = append(a.Nodes, r)
		case '>':
			return nil, 0, &ExprError{Pos: p.i, Msg: ErrExprRecordClose}
		default:
			addText(s[p.i : p.i+1])
			p.i++
		}
	}
	if group >= 0 {
		// its "}" never came
		return nil, 0, &ExprError{Pos: group, Msg: ErrExprGroupOpen}
	}
	r, err := finish(p.i)
	return r, 0, err
}

// record reads "<id>" or "<id<label>>", at p.i ("<").
func (p *exprParser) record() (*ExprRecord, error) {
	s := p.s
	start := p.i
	p.i++
	var id strings.Builder
	for {
		if p.i >= len(s) {
			return nil, &ExprError{Pos: start, Msg: ErrExprRecordOpen}
		}
		c := s[p.i]
		switch c {
		case '\\':
			if p.i+1 >= len(s) {
				return nil, &ExprError{Pos: p.i, Msg: ErrExprEscapeAtEnd}
			}
			id.WriteString(s[p.i : p.i+2])
			p.i += 2
			continue
		case '>':
			p.i++
			if id.Len() == 0 {
				return nil, &ExprError{Pos: start, Msg: ErrExprRecordID}
			}
			return &ExprRecord{Pos: Pos{start, p.i}, ID: id.String()}, nil
		case '<':
			if id.Len() == 0 {
				return nil, &ExprError{Pos: start, Msg: ErrExprRecordID}
			}
			p.i++
			var label strings.Builder
			for {
				if p.i >= len(s) {
					return nil, &ExprError{Pos: start, Msg: ErrExprLabelAtEnd}
				}
				c := s[p.i]
				if c == '\\' {
					if p.i+1 >= len(s) {
						return nil, &ExprError{Pos: p.i, Msg: ErrExprEscapeAtEnd}
					}
					label.WriteByte(s[p.i+1])
					p.i += 2
					continue
				}
				if c == '>' {
					p.i++
					break
				}
				label.WriteByte(c)
				p.i++
			}
			if p.i >= len(s) || s[p.i] != '>' {
				return nil, &ExprError{Pos: start, Msg: ErrExprRecordOpen}
			}
			p.i++
			return &ExprRecord{Pos: Pos{start, p.i}, ID: id.String(), Label: label.String()}, nil
		default:
			id.WriteByte(c)
			p.i++
		}
	}
}

// EscapeLabel is a label of a record as written: its "<", ">" and "\"
// escaped.
func EscapeLabel(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `<`, `\<`, `>`, `\>`)
	return r.Replace(s)
}

// String is the line as written, its labels with it; the canonical line is
// Normalize's.
func (a *ExprAlt) String() string {
	return a.Format()
}

// Records calls f with each record of the line, in its order: what is
// under a record (its path, the record's model) is before it.
func (a *ExprAlt) Records(f func(path []ExprNode, r *ExprRecord)) {
	var walk func(path []ExprNode, a *ExprAlt)
	walk = func(path []ExprNode, a *ExprAlt) {
		for i, n := range a.Nodes {
			switch n := n.(type) {
			case *ExprRecord:
				f(append(append([]ExprNode{}, path...), a.Nodes[:i]...), n)
			case *ExprGroup:
				for _, c := range n.Alts {
					walk(append(append([]ExprNode{}, path...), a.Nodes[:i]...), c)
				}
			}
		}
	}
	walk(nil, a)
}

// NodesText is the nodes as written, the records without their labels: a
// resource, or what begins one ("admin:posts:<7>:").
func NodesText(nodes []ExprNode) string {
	var b strings.Builder
	writeNodes(NewWriteContext(&b, WithoutLabels()), nodes)
	return b.String()
}

// StripLabels is the line without the labels of its records: "<7<First
// post>>" is "<7>". A line it cannot read is as it is.
func StripLabels(line string) string {
	if !strings.Contains(line, "<") {
		return line
	}
	var b strings.Builder
	s := line
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '\\':
			if i+1 < len(s) {
				b.WriteString(s[i : i+2])
				i++
				continue
			}
			b.WriteByte(c)
		case '<':
			// the id, up to its ">" or its label's "<"
			j := i + 1
			for j < len(s) && s[j] != '>' && s[j] != '<' {
				if s[j] == '\\' {
					j++
				}
				j++
			}
			if j >= len(s) {
				b.WriteString(s[i:])
				return b.String()
			}
			b.WriteString(s[i:j])
			if s[j] == '>' {
				b.WriteByte('>')
				i = j
				continue
			}
			// a label: up to its ">", then the record's
			k := j + 1
			for k < len(s) && s[k] != '>' {
				if s[k] == '\\' {
					k++
				}
				k++
			}
			if k+1 >= len(s) || s[k+1] != '>' {
				b.WriteString(s[j:])
				return b.String()
			}
			b.WriteByte('>')
			i = k + 1
		default:
			b.WriteByte(c)
		}
	}
	return b.String()
}

// rootList is the line with its alternatives at the root, by commas or
// lines, as a group: "a, deny b" is "{a, deny b}". One is as it is.
func rootList(pattern string) string {
	depth := 0
	for i := 0; i < len(pattern); i++ {
		switch pattern[i] {
		case '\\':
			i++
		case '{', '<', '[':
			depth++
		case '}', '>', ']':
			if depth > 0 {
				depth--
			}
		case ',', '\n':
			if depth == 0 {
				return "{" + pattern + "}"
			}
		}
	}
	return pattern
}
