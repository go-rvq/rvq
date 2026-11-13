package perm

import (
	"io"
	"strings"
)

// WriteContext is where a line is written (ExprNode.WriteCtx) and how: as
// written, one line, or indented — each alternative of a group on a line of
// its own, by how deep the group is —; its records with their labels or
// not.
type WriteContext struct {
	io.Writer
	// last is the last byte written
	last byte
	// Depth is how many groups the node written is in
	Depth int
	// Indent: each alternative of a group on a line of its own — the line
	// its separator —, Indent once a group it is in, the group after a space
	// ("admin: {"); "": the group on the line, its alternatives by ", "
	Indent string
	// NoLabels: the records without their labels ("<7>")
	NoLabels bool
	// sort: of FormatLine, its lines in their order first (WithSort)
	sort bool
}

// WriteOption is an option of a WriteContext.
type WriteOption func(ctx *WriteContext)

// WithIndent writes each alternative of a group on a line of its own,
// indented by indent once a group it is in.
func WithIndent(indent string) WriteOption {
	return func(ctx *WriteContext) { ctx.Indent = indent }
}

// WithoutLabels writes the records without their labels.
func WithoutLabels() WriteOption {
	return func(ctx *WriteContext) { ctx.NoLabels = true }
}

// NewWriteContext is a context writing to w, as opts say.
func NewWriteContext(w io.Writer, opts ...WriteOption) *WriteContext {
	ctx := &WriteContext{Writer: w}
	for _, o := range opts {
		o(ctx)
	}
	return ctx
}

// WriteString writes s.
func (ctx *WriteContext) WriteString(s string) {
	if s == "" {
		return
	}
	_, _ = io.WriteString(ctx.Writer, s)
	ctx.last = s[len(s)-1]
}

// newLine ends a line and indents the next by Depth.
func (ctx *WriteContext) newLine() {
	ctx.WriteString("\n" + strings.Repeat(ctx.Indent, ctx.Depth))
}

func (t *ExprText) WriteCtx(ctx *WriteContext) { ctx.WriteString(t.Text) }

func (r *ExprRecord) WriteCtx(ctx *WriteContext) {
	ctx.WriteString("<" + r.ID)
	if r.Label != "" && !ctx.NoLabels {
		ctx.WriteString("<" + EscapeLabel(r.Label) + ">")
	}
	ctx.WriteString(">")
}

func (g *ExprGroup) WriteCtx(ctx *WriteContext) {
	if ctx.Indent != "" && ctx.last == ':' {
		ctx.WriteString(" ")
	}
	ctx.WriteString("{")
	if ctx.Indent == "" {
		for i, a := range g.Alts {
			if i > 0 {
				ctx.WriteString(", ")
			}
			a.WriteCtx(ctx)
		}
		ctx.WriteString("}")
		return
	}
	// indented: a space after the ":" before it ("admin: {"), each
	// alternative a line — the line its separator, no ","
	ctx.Depth++
	for _, a := range g.Alts {
		ctx.newLine()
		a.WriteCtx(ctx)
	}
	ctx.Depth--
	ctx.newLine()
	ctx.WriteString("}")
}

// WriteCtx writes the alternative: its prefix, its parts.
func (a *ExprAlt) WriteCtx(ctx *WriteContext) {
	ctx.WriteString(a.Kind.Prefix())
	writeNodes(ctx, a.Nodes)
}

func writeNodes(ctx *WriteContext, nodes []ExprNode) {
	for _, n := range nodes {
		n.WriteCtx(ctx)
	}
}

// Format is the line as opts write it ("": as written).
func (a *ExprAlt) Format(opts ...WriteOption) string {
	var b strings.Builder
	a.WriteCtx(NewWriteContext(&b, opts...))
	return b.String()
}

// FormatLine is the line, canonical (Normalize), as opts write it — of
// WithIndent, to read; of WithSort, its lines in their order —; a line that cannot be read, as it is, and its
// error.
func FormatLine(line string, opts ...WriteOption) (string, error) {
	var probe WriteContext
	for _, o := range opts {
		o(&probe)
	}
	if probe.sort {
		sorted, err := SortValue(line)
		if err != nil {
			return line, err
		}
		line = sorted
	}
	e, err := ParseExpr(line)
	if err != nil {
		return line, err
	}
	return e.Normalize().Format(opts...), nil
}
