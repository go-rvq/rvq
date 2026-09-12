package vuetifyx

import (
	"encoding/json"

	h "github.com/go-rvq/htmlgo"
)

// VXCodeBuilder renders the <vx-code> component: a readonly code viewer with
// Prism syntax highlighting, numbered lines and optional per-line diff marks.
type VXCodeBuilder struct {
	tag *h.HTMLTagBuilder
}

// VXCode builds a code viewer for the given source. Set the language (Prism
// grammar name, e.g. "json", "gad", "markup") with Language; an unknown or empty
// language renders as escaped plain text with line numbers.
func VXCode(code string) *VXCodeBuilder {
	return &VXCodeBuilder{tag: h.Tag("vx-code").Attr("code", code)}
}

func (b *VXCodeBuilder) Language(v string) *VXCodeBuilder {
	b.tag.Attr("language", v)
	return b
}

// MarkKind sets how MarkLines are highlighted: "add" (success) or "del" (error).
func (b *VXCodeBuilder) MarkKind(v string) *VXCodeBuilder {
	b.tag.Attr("mark-kind", v)
	return b
}

// MarkLines highlights the given 1-based line numbers as diff marks.
func (b *VXCodeBuilder) MarkLines(lines []int) *VXCodeBuilder {
	j, _ := json.Marshal(lines)
	b.tag.Attr(":mark-lines", string(j))
	return b
}

func (b *VXCodeBuilder) Attr(vs ...interface{}) *VXCodeBuilder {
	b.tag.Attr(vs...)
	return b
}

func (b *VXCodeBuilder) Write(ctx *h.Context) error {
	return b.tag.Write(ctx)
}
