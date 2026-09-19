package admin

import (
	"fmt"
	"sort"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// advancePos moves a (1-based line, 0-based rune col) cursor over text, counting
// newlines. Used for the Equal chunks, which occupy both sides.
func advancePos(line, col int, text string) (int, int) {
	for _, r := range text {
		if r == '\n' {
			line++
			col = 0
		} else {
			col++
		}
	}
	return line, col
}

// addSideRanges records the change region text as per-line [start,end,hunk] rune
// ranges on one side's map, advancing and returning the cursor.
func addSideRanges(m map[int][][3]int, line, col int, text string, hunk int) (int, int) {
	runeStart := col
	for _, r := range text {
		if r == '\n' {
			if col > runeStart {
				m[line] = append(m[line], [3]int{runeStart, col, hunk})
			}
			line++
			col = 0
			runeStart = 0
		} else {
			col++
		}
	}
	if col > runeStart {
		m[line] = append(m[line], [3]int{runeStart, col, hunk})
	}
	return line, col
}

// hunkLineRanges walks the current→target diff (the same diffOps applyHunks uses,
// so the hunk indexes match) and returns, per side, the changed character ranges
// keyed by 1-based line, each tagged with its hunk index, plus the hunk count.
// cur holds the deletions (current side), tgt the insertions (revision side).
func hunkLineRanges(current, target string) (cur, tgt map[int][][3]int, n int) {
	cur = map[int][][3]int{}
	tgt = map[int][][3]int{}
	curLine, curCol := 1, 0
	tgtLine, tgtCol := 1, 0
	idx := -1
	inHunk := false
	for _, df := range diffOps(current, target, false) {
		switch df.Type {
		case diffmatchpatch.DiffEqual:
			curLine, curCol = advancePos(curLine, curCol, df.Text)
			tgtLine, tgtCol = advancePos(tgtLine, tgtCol, df.Text)
			inHunk = false
		case diffmatchpatch.DiffDelete:
			if !inHunk {
				idx++
				inHunk = true
			}
			curLine, curCol = addSideRanges(cur, curLine, curCol, df.Text, idx)
		case diffmatchpatch.DiffInsert:
			if !inHunk {
				idx++
				inHunk = true
			}
			tgtLine, tgtCol = addSideRanges(tgt, tgtLine, tgtCol, df.Text, idx)
		}
	}
	return cur, tgt, idx + 1
}

// rangeLines are the 1-based line numbers that have change ranges (sorted), for
// the whole-line tint (MarkLines) alongside the per-span highlight.
func rangeLines(m map[int][][3]int) []int {
	lines := make([]int, 0, len(m))
	for ln := range m {
		lines = append(lines, ln)
	}
	sort.Ints(lines)
	return lines
}

// hunkToggleSlot is the `changed` slot content: a clickable span (the changed
// text) that toggles its hunk in the selection (locals.hunks), highlighted
// stronger when selected. Bound to the slot scope's `hunk`/`text`/`added`.
func hunkToggleSlot() h.HTMLComponent {
	// A clickable changed span (like Post.Body's vx-diff-hunks): a green/red tint
	// per side (vx-hunk-add/del); when its hunk is selected, a super-highlight
	// (vx-hunk-sel) and a leading ✓ check. Styling is by class (defined globally
	// in CodeView.vue). Clicking toggles the hunk in the selection.
	return h.Tag("span").Children(
		v.VIcon("mdi-check").Size("x-small").Class("me-1 vx-hunk-check").
			Attr("v-if", "locals.hunks.includes(hunk)"),
		h.Span("").Attr("v-text", "text"),
	).
		Attr("role", "button").
		Class("vx-hunk mx-1").
		Attr(":class", "(added ? 'vx-hunk-add' : 'vx-hunk-del') + (locals.hunks.includes(hunk) ? ' vx-hunk-sel' : '')").
		Attr("@click", "locals.hunks = locals.hunks.includes(hunk) ? locals.hunks.filter(function(x){return x!==hunk}) : locals.hunks.concat([hunk])")
}

// prismCodeHunkPanel is the partial-revert panel rendered with vx-code (Prism
// syntax + intra-line change highlighting): the current and the revision value
// side by side, each changed span clickable to select its hunk; a top bar with
// the whole-field revert (leading), select/clear-all and a counter; and an apply
// button that reverts exactly the selected hunks (git checkout -p). It reuses the
// same hunk indexing and revert action as the HTML/text vx-diff-hunks panel.
func (mh *ModelHistory) prismCodeHunkPanel(recordKey, field, language, current, target, targetHash string, leading h.HTMLComponent, ctx *web.EventContext) h.HTMLComponent {
	msgr := getMessages(ctx.Context())
	cur, tgt, n := hunkLineRanges(current, target)
	if n == 0 {
		return v.VAlert(h.Text(msgr.NoHunks)).Type("info").Variant(v.VariantTonal).Density(v.DensityCompact)
	}

	toggleAll := v.VBtn("").Icon("mdi-checkbox-multiple-marked-outline").
		Variant(v.VariantText).Size(v.SizeSmall).
		Attr("title", msgr.ToggleAll).
		Attr("@click", fmt.Sprintf(
			"locals.hunks = locals.hunks.length === %d ? [] : Array.from({length: %d}, function(_,i){return i})", n, n))

	counter := h.Span("").Class("text-caption text-medium-emphasis ms-1").
		Attr("v-text", fmt.Sprintf("locals.hunks.length + '/%d'", n))

	topBar := h.Div(leading, toggleAll, counter).Class("d-flex align-center ga-1 mb-2")

	side := func(label string, code h.HTMLComponent, cls string) h.HTMLComponent {
		return h.Div(
			h.Div(h.Strong(label)).Class("text-caption text-medium-emphasis mb-1"),
			code,
		).Class("flex-1-1-0 " + cls)
	}

	curCode := vx.VXCode(current).Language(language).MarkKind("del").
		MarkLines(rangeLines(cur)).ChangeRanges(cur).ChangedSlot(hunkToggleSlot())
	tgtCode := vx.VXCode(target).Language(language).MarkKind("add").
		MarkLines(rangeLines(tgt)).ChangeRanges(tgt).ChangedSlot(hunkToggleSlot())

	apply := v.VBtn(msgr.RevertSelected).
		Color("warning").Variant(v.VariantTonal).PrependIcon("mdi-history").
		Attr("title", msgr.RevertSelected).
		Attr(":disabled", "!locals.hunks.length").
		Attr("@click", mh.revertButtonClick(recordKey, targetHash, field, "", "locals.hunks"))

	return web.Scope(
		h.Div(
			topBar,
			h.Div(h.Text(msgr.PartialRevertHint)).Class("text-caption text-medium-emphasis mb-2"),
			h.Div(
				side(msgr.Current, curCode, "pe-2").(*h.HTMLTagBuilder).Style("border-right:1px solid rgba(0,0,0,.12)"),
				side(msgr.Revision, tgtCode, "ps-2"),
			).Class("d-flex align-start"),
			h.Div(apply).Class("d-flex justify-end mt-3"),
		).Class("mt-4"),
	).LocalsInit("{ hunks: [] }")
}
