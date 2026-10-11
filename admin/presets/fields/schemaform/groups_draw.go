//go:build !formscore

package schemaform

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// drawRows draws fields by their rows: a group's a row of columns, each at
// its share (Field.Percent); on a narrow screen they stack.
// tabAt, when given, is the JS of the tab a group of tabs or steps opens
// at: the first whose fields have what the server found wrong.
func drawRows(fields []*Field, draw func(f *Field) h.HTMLComponent, tabAt ...func(row []*Field) string) (comps []h.HTMLComponent) {
	for _, row := range Rows(fields) {
		if row[0].Row == 0 {
			comps = append(comps, draw(row[0]))
			continue
		}
		if g := row[0].Group; g != nil && (g.Kind == GroupTabs || g.Kind == GroupSteps) {
			start := "0"
			if len(tabAt) > 0 && tabAt[0] != nil {
				start = tabAt[0](row)
			}
			comps = append(comps, drawTabs(row, g, draw, start))
			continue
		}
		cols := make([]h.HTMLComponent, len(row))
		for i, f := range row {
			cols[i] = v.VCol(draw(f)).Attr("style", columnStyle(f.Percent)).
				Class("py-0 sf-col").Attr("data-schemaform-column", f.Name)
		}
		comps = append(comps, h.Style(columnCSS), v.VRow(cols...).Dense(true).Class("mb-1").
			Attr("data-schemaform-row", true))
	}
	return
}

// fieldTitle is a field as a tab names it: its label, else its name.
func fieldTitle(f *Field) string {
	if l := f.Meta.String("label"); l != "" {
		return l
	}
	return presets.HumanizeString(f.Name)
}

// drawTabs draws the fields of a group of tabs — or of steps, numbered, with
// back and next — one panel each, the one chosen shown. Next checks the
// step's fields first (a VForm of the step: their rules), and goes on only
// when nothing is wrong; a tab is chosen freely.
func drawTabs(row []*Field, g *Group, draw func(f *Field) h.HTMLComponent, start string) h.HTMLComponent {
	steps := g.Kind == GroupSteps
	id := stepsSeq.Add(1)
	tabs := make([]h.HTMLComponent, len(row))
	panels := make([]h.HTMLComponent, len(row))
	for i, f := range row {
		title := fieldTitle(f)
		if steps {
			title = strconv.Itoa(i+1) + ". " + title
		}
		tab := v.VTab(h.Text(title)).Value(i).Attr("data-schemaform-tab", f.Name)
		if icon := f.Meta.String("icon"); strings.HasPrefix(icon, "mdi-") {
			tab.PrependIcon(icon)
		}
		tabs[i] = tab
		var nav h.HTMLComponent
		if steps {
			var back, next h.HTMLComponent
			if i > 0 {
				back = v.VBtn("").PrependIcon("mdi-chevron-left").Text("‹").Variant(v.VariantText).
					Attr("@click", "locals.tab--")
			}
			if i < len(row)-1 {
				// the step's fields checked first (its VForm: their rules);
				// what is wrong said under each, the step kept
				next = v.VBtn("").AppendIcon("mdi-chevron-right").Text("›").Variant(v.VariantTonal).
					Attr("@click", fmt.Sprintf(
						// (no Promise: the template sees few globals)
						"$refs[%q] ? $refs[%q].validate().then(r => { if (r.valid) locals.tab++ }) : locals.tab++", stepRef(id, i), stepRef(id, i))).
					Attr("data-schemaform-next", f.Name)
			}
			nav = h.Div(back, v.VSpacer(), next).Class("d-flex mt-2")
		}
		var body h.HTMLComponent = draw(f)
		if steps {
			// a step's fields in a form of their own: what next checks
			// a named ref (not the scope's: a component kept in it would draw
			// it again, and again)
			body = v.VForm(body).Attr("ref", stepRef(id, i)).
				Attr("validate-on", "blur").Attr("@submit.prevent", "")
		}
		panels[i] = h.Div(body, nav).Attr("v-show", "locals.tab === "+strconv.Itoa(i)).
			Attr("data-schemaform-panel", f.Name)
	}
	return web.Scope(
		h.Div(
			v.VTabs(tabs...).Attr("v-model", "locals.tab").Density(v.DensityCompact).Class("mb-3"),
			h.Div(panels...),
		).Class("mb-3").Attr("data-schemaform-"+g.Kind, g.Name),
	).LocalsInit("{ tab: " + start + " }")
}

// columnCSS lays out the columns of a group: the whole row on a small
// screen, their share of it (--sf-col) from a medium one up.
const columnCSS = `.sf-col{flex:0 0 100%;max-width:100%}` +
	`@media (min-width:960px){.sf-col{flex:0 0 var(--sf-col);max-width:var(--sf-col)}}`

// columnStyle is the width of a column at percent of its row, from a medium
// screen up (a small one stacks it: `cols=12`).
func columnStyle(percent float64) string {
	return fmt.Sprintf("--sf-col: %s%%", strconv.FormatFloat(percent, 'f', -1, 64))
}

var stepsSeq atomic.Uint64

// stepRef is the ref of the form of the step i of the group of steps id.
func stepRef(id uint64, i int) string {
	return fmt.Sprintf("sfStep_%d_%d", id, i)
}
