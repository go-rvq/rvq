package schemaform

import (
	"fmt"
	"regexp"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/microcosm-cc/bluemonday"
)

// ListMaxItems is how many items of a list the compact view (a listing cell)
// shows before it says how many are left.
const ListMaxItems = 5

// Display registers (or replaces) the component that SHOWS a value of this type
// — the detail page and the listing cell, where nothing is edited. It is the
// read side of Type: the same Context, with the value itself in Context.Data
// instead of an expression to bind, and Context.Compact set in a listing cell.
//
// FormType is always known: it is the builder itself, one level down.
func (b *Builder) Display(name string, f ComponentFunc) *Builder {
	b.displays[name] = f
	return b
}

// DisplayFunc is the component registered to show a type, and whether there is
// one.
func (b *Builder) DisplayFunc(name string) (ComponentFunc, bool) {
	if name == FormType {
		return b.formDisplayFunc, true
	}
	f, ok := b.displays[name]
	return f, ok
}

// DetailComponent shows value — decoded, as DecodeValue reads it — the way the
// schema describes it: a record as its labelled fields, a list one item under
// the other, a table layout as a table. Nothing in it is editable.
func (b *Builder) DetailComponent(ctx *web.EventContext, schema *Schema, value any) h.HTMLComponent {
	return b.showRoot(ctx, schema, value, false)
}

// ListComponent is the same value kept to a LISTING CELL: a record on one line,
// a list as the items a reader tells apart — the title of each record, the
// values of a list of plain values — up to ListMaxItems, and how many more.
func (b *Builder) ListComponent(ctx *web.EventContext, schema *Schema, value any) h.HTMLComponent {
	return b.showRoot(ctx, schema, value, true)
}

func (b *Builder) showRoot(ctx *web.EventContext, schema *Schema, value any, compact bool) h.HTMLComponent {
	if schema == nil {
		return errorComponent("schemaform: o field não tem schema")
	}
	return b.formDisplayFunc(&Context{
		Field:   &Field{Type: FormType, Schema: schema},
		Data:    value,
		Event:   ctx,
		Builder: b,
		Compact: compact,
	})
}

// displayFieldFunc is the component that shows this field — the same choice
// fieldFunc makes for editing: a display registered for the type wins, a closed
// list of values shows the label of the one held, and whatever is left shows
// as text. A read view never refuses a value for its type.
func (b *Builder) displayFieldFunc(c *Context) ComponentFunc {
	draw, registered := b.DisplayFunc(c.Field.Type)
	if registered && c.Field.Type != DefaultType {
		return draw
	}
	if items, err := c.EnumItems(); err == nil && len(items) > 0 {
		return EnumDisplayFunc
	}
	if registered {
		return draw
	}
	return TextDisplayFunc
}

// formDisplayFunc shows FormType: a record, a list of plain values, a list of
// records (as forms or as a table).
func (b *Builder) formDisplayFunc(c *Context) h.HTMLComponent {
	schema := c.Field.Schema
	if schema == nil {
		return errorComponent(fmt.Sprintf("schemaform: o field %q é um form sem schema", c.Field.Name))
	}
	if schema.Choice {
		return b.showChoice(schema, c)
	}
	if !schema.Slice {
		rec := b.showRecord(schema, c, c.Data)
		// a record inside the record: set in, behind a line, as in the form (its
		// label is above it already); the root, and a cell, as they are
		if c.Field.Name != "" && !c.Compact {
			return h.Div(rec).Class(nestedRecordClass)
		}
		return rec
	}

	items, _ := c.Data.([]any)
	if len(items) == 0 {
		return emptyDisplay(c)
	}
	switch {
	case schema.Item != nil:
		return b.showValues(schema, c, items)
	case c.Compact:
		return b.showTitles(schema, c, items)
	case schema.Layout == LayoutTable:
		return b.showTable(schema, c, items)
	case schema.Layout == LayoutGrid:
		return b.showGrid(schema, c, items)
	default:
		return b.showRecords(schema, c, items)
	}
}

// showChoice shows a choice (Schema.Choice): the class chosen, by its label,
// and what it holds — set in, as a record inside the record is.
func (b *Builder) showChoice(schema *Schema, c *Context) h.HTMLComponent {
	m, _ := c.Data.(map[string]any)
	for _, f := range schema.Fields {
		data, ok := m[f.Name]
		if !ok {
			continue
		}
		fc := b.fieldContext(c, f, data)
		if c.Compact {
			return h.Tag("span").Children(h.Text(displayLabel(fc)))
		}
		return h.Div(
			h.Div(h.Text(displayLabel(fc))).Class("text-body-2 mb-2"),
			h.Div(b.showRecord(f.Schema, fc, data)).Class(nestedRecordClass),
		)
	}
	return emptyDisplay(c)
}

// fieldContext is the context a field of a record is shown with.
func (b *Builder) fieldContext(c *Context, f *Field, data any) *Context {
	return &Context{
		Field:   f,
		Data:    data,
		Path:    join(c.Path, f.Name),
		Event:   c.Event,
		Builder: b,
		Compact: c.Compact,
	}
}

func (b *Builder) showField(fc *Context) h.HTMLComponent {
	return b.displayFieldFunc(fc)(fc)
}

// showRecord shows a record's fields with their labels: one under the other in
// a detail, one after the other on a single line in a cell — where an empty
// field is left out, since there is no room to say it is empty.
func (b *Builder) showRecord(schema *Schema, c *Context, data any) h.HTMLComponent {
	var comps []h.HTMLComponent
	for _, f := range schema.Fields {
		fc := b.fieldContext(c, f, recordValue(data, f.Name))
		if c.Compact {
			if isEmpty(fc.Data) {
				continue
			}
			if len(comps) > 0 {
				comps = append(comps, h.Span(" · ").Class("text-medium-emphasis"))
			}
			comps = append(comps,
				h.Span(displayLabel(fc)+": ").Class("text-medium-emphasis"),
				b.showField(fc))
			continue
		}
		comps = append(comps, h.Div(
			h.Div(h.Text(displayLabel(fc))).Class("text-caption text-medium-emphasis"),
			h.Div(b.showField(fc)).Class("pt-1"),
		).Class("mb-3"))
	}
	if c.Compact {
		return h.Tag("span").Children(comps...)
	}
	return h.Div(comps...)
}

// showValues shows a list of plain values: a bulleted list in a detail, the
// values separated by commas in a cell.
func (b *Builder) showValues(schema *Schema, c *Context, items []any) h.HTMLComponent {
	show := func(data any) h.HTMLComponent {
		return b.showField(&Context{
			Field: schema.Item, Data: data, Path: c.Path,
			Event: c.Event, Builder: b, Compact: c.Compact, noLabel: true,
		})
	}
	if c.Compact {
		var comps []h.HTMLComponent
		for i, it := range limit(items) {
			if i > 0 {
				comps = append(comps, h.Text(", "))
			}
			comps = append(comps, show(it))
		}
		if more := moreItems(items); more != nil {
			comps = append(comps, more)
		}
		return h.Tag("span").Children(comps...)
	}
	lis := make([]h.HTMLComponent, len(items))
	for i, it := range items {
		lis[i] = h.Li(show(it))
	}
	return h.Ul(lis...).Class("ps-4 mb-0")
}

// showRecords shows a list of records laid out as the default list: each
// record a card of its labelled fields, one under the other — as the editor
// draws them.
func (b *Builder) showRecords(schema *Schema, c *Context, items []any) h.HTMLComponent {
	comps := make([]h.HTMLComponent, len(items))
	for i, it := range items {
		card := v.VCard(v.VCardText(b.showRecord(schema, c, it))).Variant(v.VariantOutlined)
		if i > 0 {
			card.Class("mt-3")
		}
		comps[i] = card
	}
	return h.Div(comps...)
}

// showGrid shows a list of records in the grid layout: each record a card of
// its labelled fields, schema.GridColumns() side by side.
func (b *Builder) showGrid(schema *Schema, c *Context, items []any) h.HTMLComponent {
	cols := make([]h.HTMLComponent, len(items))
	for i, it := range items {
		cols[i] = gridCol(schema,
			v.VCard(v.VCardText(b.showRecord(schema, c, it))).Variant(v.VariantOutlined).Class("h-100"),
		)
	}
	return v.VRow(cols...).Dense(true)
}

// showTable shows a list of records in the table layout: the columns the
// schema names, the labels in the header, each cell shown compact.
func (b *Builder) showTable(schema *Schema, c *Context, items []any) h.HTMLComponent {
	cols := schema.TableColumns()
	head := make([]h.HTMLComponent, len(cols))
	for i, f := range cols {
		hc := b.fieldContext(c, f, nil)
		th := h.Th(displayLabel(hc)).Class("text-left text-no-wrap")
		if hint := hc.Info().Hint; hint != "" {
			th.Attr("title", hint)
		}
		head[i] = th
	}

	rows := make([]h.HTMLComponent, len(items))
	for r, it := range items {
		cells := make([]h.HTMLComponent, len(cols))
		for i, f := range cols {
			fc := b.fieldContext(c, f, recordValue(it, f.Name))
			fc.Compact = true
			cells[i] = h.Td(b.showField(fc)).Class("py-1")
		}
		rows[r] = h.Tr(cells...)
	}

	return v.VTable(h.Thead(h.Tr(head...)), h.Tbody(rows...)).Density(v.DensityCompact)
}

// showTitles is a list of records in a cell: one line per record, showing the
// field a reader names it by (see firstTextField) — or the whole record, on one
// line, when it has none.
func (b *Builder) showTitles(schema *Schema, c *Context, items []any) h.HTMLComponent {
	title := firstTextField(schema)
	var tf *Field
	for _, f := range schema.Fields {
		if f.Name == title {
			tf = f
			break
		}
	}

	lines := make([]h.HTMLComponent, 0, len(items)+1)
	for _, it := range limit(items) {
		var line h.HTMLComponent
		if tf != nil {
			line = b.showField(b.fieldContext(c, tf, recordValue(it, tf.Name)))
		} else {
			line = b.showRecord(schema, c, it)
		}
		lines = append(lines, h.Div(line).Class("text-no-wrap"))
	}
	if more := moreItems(items); more != nil {
		lines = append(lines, h.Div(more))
	}
	return h.Div(lines...)
}

// limit is what a cell shows of a list.
func limit(items []any) []any {
	if len(items) > ListMaxItems {
		return items[:ListMaxItems]
	}
	return items
}

// moreItems says how many items a cell left out, or nil.
func moreItems(items []any) h.HTMLComponent {
	if n := len(items) - ListMaxItems; n > 0 {
		return h.Span(fmt.Sprintf(" +%d", n)).Class("text-medium-emphasis")
	}
	return nil
}

// displayLabel is a field's label in a read view — where a cell, too, names
// what it shows, since there is no column header to do it.
func displayLabel(c *Context) string {
	if l := c.Info().Label; l != "" {
		return l
	}
	return presets.HumanizeString(c.Field.Name)
}

// emptyDisplay is what an empty value shows: a dash in a detail, nothing in a
// cell.
func emptyDisplay(c *Context) h.HTMLComponent {
	if c.Compact {
		return h.Text("")
	}
	return h.Span("—").Class("text-disabled")
}

// recordValue is the value of a record's field, whatever the record was decoded
// into: a Record (from the form), or a map (from the stored YAML or JSON).
func recordValue(data any, name string) any {
	switch r := data.(type) {
	case Record:
		for _, f := range r {
			if f.Name == name {
				return f.Value
			}
		}
	case map[string]any:
		return r[name]
	case map[any]any:
		return r[name]
	}
	return nil
}

func isEmpty(data any) bool {
	switch d := data.(type) {
	case nil:
		return true
	case string:
		return strings.TrimSpace(d) == ""
	case []any:
		return len(d) == 0
	}
	return false
}

// displayText is a plain value as text.
func displayText(data any) string {
	if data == nil {
		return ""
	}
	return fmt.Sprint(data)
}

// TextDisplayFunc shows a value as text — DefaultType, the numbers, the date
// and time types, and whatever type has no display of its own.
func TextDisplayFunc(c *Context) h.HTMLComponent {
	if isEmpty(c.Data) {
		return emptyDisplay(c)
	}
	return h.Text(displayText(c.Data))
}

// LongTextDisplayFunc is "text": prose, keeping its line breaks — on one line
// in a cell.
func LongTextDisplayFunc(c *Context) h.HTMLComponent {
	if isEmpty(c.Data) || c.Compact {
		return TextDisplayFunc(c)
	}
	return h.Div(h.Text(displayText(c.Data))).Style("white-space:pre-wrap")
}

var (
	htmlDisplayPolicy = bluemonday.UGCPolicy()
	htmlTags          = regexp.MustCompile(`<[^>]*>`)
)

// HTMLDisplayFunc is "html": the markup, sanitized — the value comes from a
// form, not from code. A cell shows its text only.
func HTMLDisplayFunc(c *Context) h.HTMLComponent {
	if isEmpty(c.Data) {
		return emptyDisplay(c)
	}
	s := displayText(c.Data)
	if c.Compact {
		return h.Text(strings.Join(strings.Fields(htmlTags.ReplaceAllString(s, " ")), " "))
	}
	return h.RawHTML(htmlDisplayPolicy.Sanitize(s))
}

// BoolDisplayFunc is "bool": a check when on, a dash when off — in the colour a
// switch has when it is on, like the editor's.
func BoolDisplayFunc(c *Context) h.HTMLComponent {
	on, _ := c.Data.(bool)
	if on {
		return v.VIcon("mdi-check").Color("primary").Size(v.SizeSmall)
	}
	return v.VIcon("mdi-minus").Color("grey").Size(v.SizeSmall)
}

// ColorDisplayFunc is "color": a swatch of it, and its value.
func ColorDisplayFunc(c *Context) h.HTMLComponent {
	if isEmpty(c.Data) {
		return emptyDisplay(c)
	}
	s := displayText(c.Data)
	return h.Tag("span").Children(
		h.Span("").Class("d-inline-block rounded me-1").
			Style(fmt.Sprintf("width:1em;height:1em;vertical-align:middle;border:1px solid rgba(0,0,0,.2);background:%s", cssColor(s))),
		h.Text(s),
	)
}

var cssColorRe = regexp.MustCompile(`^(#[0-9a-fA-F]{3,8}|[a-zA-Z]+|(rgb|hsl)a?\([0-9.,%\s]+\))$`)

// cssColor lets through a value that is a color and nothing else, since it goes
// into a style attribute.
func cssColor(s string) string {
	if cssColorRe.MatchString(s) {
		return s
	}
	return "transparent"
}

// EnumDisplayFunc shows the label of the value an enum field holds — the same
// label the select offered it by — or the value itself when it is not in the
// list.
func EnumDisplayFunc(c *Context) h.HTMLComponent {
	if isEmpty(c.Data) {
		return emptyDisplay(c)
	}
	name := displayText(c.Data)
	if items, err := c.EnumItems(); err == nil {
		for _, it := range items {
			if it.Name == name {
				return h.Text(it.Label)
			}
		}
	}
	return h.Text(name)
}

func defaultDisplays() map[string]ComponentFunc {
	return map[string]ComponentFunc{
		DefaultType: TextDisplayFunc,
		"int":       TextDisplayFunc,
		"uint":      TextDisplayFunc,
		"float":     TextDisplayFunc,
		"decimal":   TextDisplayFunc,
		"duration":  TextDisplayFunc,
		"time":      TextDisplayFunc,
		"date":      TextDisplayFunc,
		"text":      LongTextDisplayFunc,
		"html":      HTMLDisplayFunc,
		"bool":      BoolDisplayFunc,
		"color":     ColorDisplayFunc,
	}
}
