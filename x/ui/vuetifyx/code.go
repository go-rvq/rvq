package vuetifyx

import (
	"encoding/json"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
)

// VXCodeBuilder renders the <vx-code> component: a readonly code viewer with
// Prism syntax highlighting — the admin's one, for a document's code, a diff,
// a record's data. By default it is highlighted with Prism's plugins: line
// numbers, matching braces, a toolbar with the language and a copy button, and
// (Diff) a git diff whose lines are highlighted in the file's language. With
// diff marks (MarkKind, MarkLines, ChangeRanges) it is a row per line instead,
// its marked lines tinted.
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

// ChangeRanges highlights, within a line, the exact changed character spans
// (GoLand/IntelliJ-style). The map is keyed by 1-based line number; each value is
// a list of [start,end,hunk) triples — start/end are rune offsets in that line
// and hunk is the change region's index (-1 when there is no hunk grouping, e.g.
// a plain compare). Used with MarkLines/MarkKind so a changed line is tinted and
// its changed content is emphasized; the hunk index lets the `changed` slot make
// each span interactive (e.g. a partial-revert toggle).
func (b *VXCodeBuilder) ChangeRanges(ranges map[int][][3]int) *VXCodeBuilder {
	if len(ranges) == 0 {
		return b
	}
	j, _ := json.Marshal(ranges)
	b.tag.Attr(":change-ranges", string(j))
	return b
}

// ChangedSlot fills the `changed` scoped slot: each changed span is rendered by
// these children instead of the default highlight mark. The slot scope exposes
// { text, line, start, end, hunk, kind, added } — added is true on the NEW side
// (an insertion) and false on the OLD side (a deletion) — so a caller can render,
// say, a clickable partial-revert toggle bound to the hunk index.
func (b *VXCodeBuilder) ChangedSlot(children ...h.HTMLComponent) *VXCodeBuilder {
	b.tag.Children(web.Slot(children...).Name("changed").Scope("{ text, line, start, end, hunk, kind, added }"))
	return b
}

// Diff says the code is a diff (git's), its lines highlighted in the Language
// (Prism's diff-highlight): diff-gadx, diff-yaml; a plain diff when the
// language is none or unknown.
func (b *VXCodeBuilder) Diff(v bool) *VXCodeBuilder {
	b.tag.Attr(":diff", v)
	return b
}

// LineNumbers numbers the lines (default on).
func (b *VXCodeBuilder) LineNumbers(v bool) *VXCodeBuilder {
	b.tag.Attr(":line-numbers", v)
	return b
}

// MatchBraces highlights a brace and its pair (default on).
func (b *VXCodeBuilder) MatchBraces(v bool) *VXCodeBuilder {
	b.tag.Attr(":match-braces", v)
	return b
}

// Label is what the toolbar says of the code: by default its language's name.
func (b *VXCodeBuilder) Label(v string) *VXCodeBuilder {
	b.tag.Attr("label", v)
	return b
}

// CopyLabels are the words of the copy button: its own, once it copied, and
// when it could not.
func (b *VXCodeBuilder) CopyLabels(copy, copied, failed string) *VXCodeBuilder {
	b.tag.Attr("copy-text", copy, "copied-text", copied, "copy-error-text", failed)
	return b
}

func (b *VXCodeBuilder) Attr(vs ...interface{}) *VXCodeBuilder {
	b.tag.Attr(vs...)
	return b
}

func (b *VXCodeBuilder) Write(ctx *h.Context) error {
	return b.tag.Write(ctx)
}

// VXDiffHunksBuilder renders the <vx-diff-hunks> component: a clickable
// partial-revert diff. It shows an instruction, then a current | revision split
// whose change regions (marked with data-h="<hunk index>") the user clicks to
// select (super-highlight + a check). The selected hunk indexes are the model
// value (bind with v-model to a locals array the revert button reads).
type VXDiffHunksBuilder struct {
	tag *h.HTMLTagBuilder
}

func VXDiffHunks() *VXDiffHunksBuilder {
	return &VXDiffHunksBuilder{tag: h.Tag("vx-diff-hunks")}
}

func (b *VXDiffHunksBuilder) Instruction(v string) *VXDiffHunksBuilder {
	b.tag.Attr("instruction", v)
	return b
}

func (b *VXDiffHunksBuilder) LeftLabel(v string) *VXDiffHunksBuilder {
	b.tag.Attr("left-label", v)
	return b
}

func (b *VXDiffHunksBuilder) RightLabel(v string) *VXDiffHunksBuilder {
	b.tag.Attr("right-label", v)
	return b
}

func (b *VXDiffHunksBuilder) LeftHTML(v string) *VXDiffHunksBuilder {
	b.tag.Attr("left-html", v)
	return b
}

func (b *VXDiffHunksBuilder) RightHTML(v string) *VXDiffHunksBuilder {
	b.tag.Attr("right-html", v)
	return b
}

func (b *VXDiffHunksBuilder) Attr(vs ...interface{}) *VXDiffHunksBuilder {
	b.tag.Attr(vs...)
	return b
}

func (b *VXDiffHunksBuilder) Write(ctx *h.Context) error {
	return b.tag.Write(ctx)
}
