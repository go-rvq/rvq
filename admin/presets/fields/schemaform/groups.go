package schemaform

import (
	"fmt"
	"strconv"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// WidthUnits are the units a row of a group is divided in: a field's
// `[width=N]` is N of them, 1 to WidthUnits — 1 is 20% —; the fields that
// say none share what is left, equally.
const WidthUnits = 5

// The kinds of a group, by its metadata: `[tabs] { … }` each field a tab,
// `[steps] { … }` each field a step (with next and back; `[steps,
// vertical]` down the page), none a row of columns.
const (
	GroupRow   = "row"
	GroupTabs  = "tabs"
	GroupSteps = "steps"
)

// Group is a group of a record: `{ … }` in its body, the field Name (`$N`).
type Group struct {
	Name string
	Kind string
	// Vertical lays steps down the page (`[steps, vertical]`); across, else.
	Vertical bool
	// Meta is the group's own metadata (written before its block).
	Meta Meta
}

// groupOf is the group f — the field `$N` — is, by its metadata.
func groupOf(f *Field) *Group {
	g := &Group{Name: f.Name, Kind: GroupRow, Meta: f.Meta}
	switch {
	case flag(f.Meta, GroupTabs):
		g.Kind = GroupTabs
	case flag(f.Meta, GroupSteps):
		g.Kind = GroupSteps
		g.Vertical = flag(f.Meta, "vertical")
	}
	return g
}

// flag says the metadata m says key, and not false: `[tabs]`, `[tabs=true]`.
func flag(m Meta, key string) bool {
	v, ok := m[key]
	if !ok {
		return false
	}
	if b, isBool := v.(bool); isBool {
		return b
	}
	return v != nil
}

// appendField adds f to s. A group — `{ … }` in the body, the field `$N` of
// a record — is not a field of the form: its fields are, in its place, each
// a column of one row (Field.Row) at its share of it (Field.Percent); what
// they hold is the record's own (the group adds no name to the value). A
// group inside a group is refused: a row has columns, not rows.
func (r *reader) appendField(s *Schema, f *Field) error {
	if !isGroupName(f.Name) || f.Schema == nil || f.Schema.Slice || f.Schema.Choice {
		s.Fields = append(s.Fields, f)
		return nil
	}
	for _, gf := range f.Schema.Fields {
		if gf.Row != 0 {
			return fmt.Errorf("schemaform: the group %s has a group inside: a group is a row of columns, one level", f.Name)
		}
	}
	g := groupOf(f)
	// the widths share a row; a tab, a step, has the whole of it
	if g.Kind == GroupRow {
		if err := groupWidths(f.Name, f.Schema.Fields); err != nil {
			return err
		}
	}
	row := maxRow(s.Fields) + 1
	for _, gf := range f.Schema.Fields {
		gf.Row, gf.Group = row, g
		s.Fields = append(s.Fields, gf)
	}
	return nil
}

// isGroupName says name is a group's: `$N` (gad's `{ … }` in a body).
func isGroupName(name string) bool {
	if len(name) < 2 || name[0] != '$' {
		return false
	}
	for _, r := range name[1:] {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// Rows are fields as they are drawn: a field alone a row of its own, the
// fields of a group (Field.Row) one row of them.
func Rows(fields []*Field) (rows [][]*Field) {
	for _, f := range fields {
		if n := len(rows); f.Row != 0 && n > 0 && rows[n-1][0].Row == f.Row {
			rows[n-1] = append(rows[n-1], f)
			continue
		}
		rows = append(rows, []*Field{f})
	}
	return
}

// drawRows draws fields by their rows: a group's a row of columns, each at
// its share (Field.Percent); on a narrow screen they stack.
func drawRows(fields []*Field, draw func(f *Field) h.HTMLComponent) (comps []h.HTMLComponent) {
	for _, row := range Rows(fields) {
		if row[0].Row == 0 {
			comps = append(comps, draw(row[0]))
			continue
		}
		if g := row[0].Group; g != nil && (g.Kind == GroupTabs || g.Kind == GroupSteps) {
			comps = append(comps, drawTabs(row, g, draw))
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
// back and next — one panel each, the one chosen shown: in the admin every
// field is edited, so a step does not wait for the one before.
func drawTabs(row []*Field, g *Group, draw func(f *Field) h.HTMLComponent) h.HTMLComponent {
	steps := g.Kind == GroupSteps
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
				next = v.VBtn("").AppendIcon("mdi-chevron-right").Text("›").Variant(v.VariantTonal).
					Attr("@click", "locals.tab++")
			}
			nav = h.Div(back, v.VSpacer(), next).Class("d-flex mt-2")
		}
		panels[i] = h.Div(draw(f), nav).Attr("v-show", "locals.tab === "+strconv.Itoa(i)).
			Attr("data-schemaform-panel", f.Name)
	}
	return web.Scope(
		h.Div(
			v.VTabs(tabs...).Attr("v-model", "locals.tab").Density(v.DensityCompact).Class("mb-3"),
			h.Div(panels...),
		).Class("mb-3").Attr("data-schemaform-"+g.Kind, g.Name),
	).LocalsInit("{ tab: 0 }")
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

// appendRows adds fields — a parent's, read on their own — to s, their rows
// numbered after the ones s has.
func appendRows(s *Schema, fields []*Field) {
	offset := maxRow(s.Fields)
	for _, f := range fields {
		if f.Row != 0 {
			f.Row += offset
		}
		s.Fields = append(s.Fields, f)
	}
}

func maxRow(fields []*Field) (n int) {
	for _, f := range fields {
		if f.Row > n {
			n = f.Row
		}
	}
	return
}

// groupWidths gives each field of the group its share of the row: its
// `[width=N]`, N of WidthUnits; the ones that say none, what is left,
// equally. The widths given are whole units, 1 to WidthUnits, and the row is
// whole: they add up to WidthUnits, every field with a share.
func groupWidths(group string, fields []*Field) error {
	given, rest := 0, 0
	for _, f := range fields {
		if _, ok := f.Meta["width"]; !ok {
			rest++
			continue
		}
		n, ok, err := f.Meta.Int("width")
		if err != nil || !ok || n < 1 || n > WidthUnits {
			return fmt.Errorf("schemaform: group %s, field %s: width is a whole number from 1 to %d", group, f.Name, WidthUnits)
		}
		given += n
	}
	names := func() string {
		out := make([]string, len(fields))
		for i, f := range fields {
			out[i] = f.Name
		}
		return strings.Join(out, ", ")
	}
	switch {
	case given > WidthUnits:
		return fmt.Errorf("schemaform: group %s (%s): the widths add up to %d, more than %d", group, names(), given, WidthUnits)
	case rest == 0 && given != WidthUnits:
		return fmt.Errorf("schemaform: group %s (%s): the widths add up to %d, not %d", group, names(), given, WidthUnits)
	case rest > 0 && given == WidthUnits:
		return fmt.Errorf("schemaform: group %s (%s): the widths given take the whole row, and %d field(s) have none left", group, names(), rest)
	}
	share := float64(WidthUnits-given) / float64(max(rest, 1))
	for _, f := range fields {
		units := share
		if n, ok, _ := f.Meta.Int("width"); ok {
			units = float64(n)
		}
		f.Percent = units * 100 / WidthUnits
	}
	return nil
}
