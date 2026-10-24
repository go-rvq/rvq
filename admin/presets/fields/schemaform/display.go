//go:build !formscore

package schemaform

import (
	"context"
	"fmt"
	"github.com/gad-lang/gad"
	"log"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	xnumber "golang.org/x/text/number"
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

// FieldDisplay shows the fields of value one at a time, by their paths —
// their names joined by "." (a field of a record inside the record:
// "address.city") —, each as the detail shows it: by the display of its
// type (Display, read-only). show is nil of a path of no field; portal is
// where what they open goes, to be drawn once beside them.
func (b *Builder) FieldDisplay(ctx *web.EventContext, schema *Schema, value any) (show func(path string) h.HTMLComponent, portal h.HTMLComponent) {
	if schema == nil {
		return func(string) h.HTMLComponent { return nil }, nil
	}
	root := &Context{
		Field:   &Field{Type: FormType, Schema: schema},
		Data:    value,
		Event:   ctx,
		Builder: b,
		Portal:  NewPortalName(),
	}
	show = func(path string) h.HTMLComponent {
		c, s, data := root, schema, value
		names := strings.Split(path, ".")
		for i, name := range names {
			var f *Field
			for _, x := range s.Fields {
				if x.Name == name {
					f = x
					break
				}
			}
			if f == nil {
				return nil
			}
			record := data
			data = recordValue(data, name)
			c = b.fieldContext(c, f, data)
			if i == len(names)-1 {
				// its warnings under it (the record's, of it)
				if ws := warningsComponent(fieldWarnings(record, name)); ws != nil {
					return h.Div(b.showField(c), ws)
				}
				return b.showField(c)
			}
			if f.Schema == nil || f.Schema.Slice {
				return nil
			}
			s = f.Schema
		}
		return nil
	}
	return show, web.Portal().Name(root.Portal)
}

// ListComponent is the same value kept to a LISTING CELL: a record on one line,
// a list as the items a reader tells apart — the title of each record, the
// values of a list of plain values — up to ListMaxItems, and how many more.
func (b *Builder) ListComponent(ctx *web.EventContext, schema *Schema, value any) h.HTMLComponent {
	return b.showRoot(ctx, schema, value, true)
}

func (b *Builder) showRoot(ctx *web.EventContext, schema *Schema, value any, compact bool) h.HTMLComponent {
	if schema == nil {
		return errorComponent(messagesOf(ctx).NoSchema)
	}
	c := &Context{
		Field:   &Field{Type: FormType, Schema: schema},
		Data:    value,
		Event:   ctx,
		Builder: b,
		Compact: compact,
	}
	if compact {
		// a listing cell embeds no portal: one per cell would weigh on the
		// listing, and a field opens nothing there (Context.Portal is "")
		return compactCell(b.formDisplayFunc(c), c.Messages().ShowWhole)
	}
	c.Portal = NewPortalName()
	return h.Components(b.formDisplayFunc(c), web.Portal().Name(c.Portal))
}

// compactCell keeps a listing cell to one line, cut with "…"; a click on it
// shows it whole, and another cuts it again — without opening the record the
// row is. An empty cell stays empty.
func compactCell(comp h.HTMLComponent, title string) h.HTMLComponent {
	return h.ComponentFunc(func(ctx *h.Context) error {
		body, err := h.Marshal(comp, ctx.Context)
		if err != nil || strings.TrimSpace(string(body)) == "" {
			return err
		}
		return h.Tag("span").Children(h.RawHTML(body)).
			Class("schemaform-cell text-truncate d-inline-block").
			Style("max-width: 24em; vertical-align: bottom; cursor: pointer").
			Attr("title", title).
			Attr("@click.stop", `$event.currentTarget.classList.toggle("text-truncate")`).
			Write(ctx)
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
		return errorComponent(fmt.Sprintf(c.Messages().FormWithoutSchema, c.Field.Name))
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
		Field:  f,
		Data:   data,
		record: c.Data,
		Path:   join(c.Path, f.Name),
		Event:  c.Event, Portal: c.Portal,
		Builder: b,
		Compact: c.Compact,
	}
}

func (b *Builder) showField(fc *Context) h.HTMLComponent {
	return b.displayFieldFunc(fc)(fc)
}

// showRecord shows a record's fields with their labels, one under the other —
// in a cell, on a single line (showRecordCompact).
func (b *Builder) showRecord(schema *Schema, c *Context, data any) h.HTMLComponent {
	if c.Compact {
		return b.showRecordCompact(schema, c, data)
	}
	return h.Div(drawRows(shownFields(schema.Fields), func(f *Field) h.HTMLComponent {
		fc := b.fieldContext(c, f, recordValue(data, f.Name))
		return h.Div(
			// a label as the detail's (presets): .v-label
			h.Div(h.Text(displayLabel(fc))).Class("v-label theme--light text-caption"),
			h.Div(b.showField(fc)).Class("pt-1"),
			warningsComponent(fieldWarnings(data, f.Name)),
		).Class("mb-3")
	})...)
}

// fieldWarnings are the warnings of the field name of the record data: its
// WarningsField's — what was kept although it could not be checked —, the
// messages of the ones of that field.
func fieldWarnings(data any, name string) (out []string) {
	m, _ := data.(map[string]any)
	ws, _ := m[WarningsField].([]any)
	for _, w := range ws {
		wm, _ := w.(map[string]any)
		if n, _ := wm["fieldName"].(string); n != name {
			continue
		}
		if msg, _ := wm["message"].(string); msg != "" {
			out = append(out, msg)
		}
	}
	return
}

// WarningClass marks a warning of a field shown read-only (fieldWarnings):
// under an element of HideWarningsClass, hidden (the admin's style).
const (
	WarningClass      = "schemaform-warning"
	HideWarningsClass = "schemaform-hide-warnings"
)

// WarningsComponent are the warnings msgs of a field, as a read-only view
// shows them under its value (WarningClass); nil for none.
func WarningsComponent(msgs []string) h.HTMLComponent { return warningsComponent(msgs) }

// warningsComponent are the warnings of a field, under its value; nil for
// none.
func warningsComponent(msgs []string) h.HTMLComponent {
	if len(msgs) == 0 {
		return nil
	}
	var out []h.HTMLComponent
	for _, m := range msgs {
		out = append(out, h.Div(
			v.VIcon("mdi-alert-outline").Size(v.SizeSmall).Color("warning"),
			h.Span(m),
		).Class(WarningClass+" d-flex align-center ga-1 text-caption text-warning mt-1").Attr("data-schemaform-warning", true))
	}
	return h.Components(out...)
}

// shownFields are the fields a record shows: all but its warnings
// (WarningsField) — what was kept unchecked, told apart by who shows it.
func shownFields(fields []*Field) []*Field {
	out := make([]*Field, 0, len(fields))
	for _, f := range fields {
		if f.Name != WarningsField {
			out = append(out, f)
		}
	}
	return out
}

// showRecordCompact is a record in a listing cell: "label: value" of each
// field it has, by the labels, apart by ", " — "Layout das postagens:
// Serviços". An empty field is left out: there is no room to say it is empty.
func (b *Builder) showRecordCompact(schema *Schema, c *Context, data any) h.HTMLComponent {
	type pair struct {
		label string
		fc    *Context
	}
	var pairs []pair
	for _, f := range shownFields(schema.Fields) {
		fc := b.fieldContext(c, f, recordValue(data, f.Name))
		if isEmpty(fc.Data) {
			continue
		}
		pairs = append(pairs, pair{displayLabel(fc), fc})
	}
	sort.SliceStable(pairs, func(i, j int) bool { return pairs[i].label < pairs[j].label })
	var comps []h.HTMLComponent
	for i, p := range pairs {
		if i > 0 {
			comps = append(comps, h.Span(", ").Class("text-medium-emphasis"))
		}
		comps = append(comps,
			h.Span(p.label+": ").Class("text-medium-emphasis"),
			b.showField(p.fc))
	}
	return h.Tag("span").Children(comps...)
}

// showValues shows a list of plain values: a bulleted list in a detail, the
// values separated by commas in a cell.
func (b *Builder) showValues(schema *Schema, c *Context, items []any) h.HTMLComponent {
	show := func(data any) h.HTMLComponent {
		return b.showField(&Context{
			Field: schema.Item, Data: data, Path: c.Path,
			Event: c.Event, Portal: c.Portal, Builder: b, Compact: c.Compact, noLabel: true,
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
// RangeDisplayFunc is a Range[T] read: "from – to".
func RangeDisplayFunc(c *Context) h.HTMLComponent {
	return h.Span("").Attr("v-text", fmt.Sprintf("(%s ? [%s.from, %s.to].filter(x => x != null && x !== '').join(' – ') : '')", c.Value, c.Value, c.Value))
}

func TextDisplayFunc(c *Context) h.HTMLComponent {
	if isEmpty(c.Data) {
		return emptyDisplay(c)
	}
	return h.Text(c.Prefix() + displayText(c.Data) + c.Suffix())
}

// NumberDisplayFunc shows a number — "int", "uint", "float", "decimal" —
// as the language of the request writes it (pt-BR: 1.234,50), with its
// decimal places (`[places=N]`, else its step's, else its own) and its
// prefix and suffix: `[prefix="R$ ", places=2]` 10 is "R$ 10,00".
func NumberDisplayFunc(c *Context) h.HTMLComponent {
	if isEmpty(c.Data) {
		return emptyDisplay(c)
	}
	return h.Text(displayTextOf(eventContext(c.Event), c.Field, c.Data, c.record))
}

// eventContext is the context of the request of ev; Background without one.
func eventContext(ev *web.EventContext) context.Context {
	if ev != nil && ev.R != nil {
		return ev.R.Context()
	}
	return context.Background()
}

// DisplayText is the value v of the field f as the read views show it, as
// text: a number (NumberDisplayFunc), a date, a time (DateTimeDisplayFunc)
// as the language of ctx writes them, and anything else as it is — each with
// its prefix and suffix; "" for none.
//
// A float or a decimal of `[roundDisplay=…]` is rounded by it first, this an
// empty record (displayTextOf: the record when it is known).
func DisplayText(ctx context.Context, f *Field, v any) string {
	return displayTextOf(ctx, f, v, nil)
}

// displayTextOf is DisplayText, the field of the record record (this of its
// `[roundDisplay=…]`).
func displayTextOf(ctx context.Context, f *Field, v any, record any) string {
	if isEmpty(v) {
		return ""
	}
	if f != nil && f.RoundDisplayFunc != nil {
		v = roundShown(ctx, f, v, record)
	}
	var prefix, suffix, typ string
	if f != nil {
		prefix, _ = f.Meta[MetaPrefix].(string)
		suffix, _ = f.Meta[MetaSuffix].(string)
		typ = f.Type
	}
	text := displayText(v)
	switch typ {
	case "int", "uint", "float", "decimal":
		if p := fixedPlaces(f); p > 0 {
			// fixed point: the integer kept ÷ 10^places
			if n, err := strconv.ParseFloat(strings.TrimSpace(displayText(v)), 64); err == nil {
				v = n / math.Pow10(p)
			}
		}
		text = FormatNumber(contextLanguage(ctx), v, fieldDecimals(f))
	case "date", "time", "cdate", "cdtime", "ctime":
		if t, hasDate, hasTime, ok := parseDateTime(v); ok {
			text = FormatDateTime(ctx, t, hasDate, hasTime)
		}
	}
	return prefix + text + suffix
}

// roundShown is v rounded by the `[roundDisplay=…]` of f (MetaRoundDisplay),
// called with record as this (its values, as given), no parents, v and
// places — its `[places=N]`, 0 without —; v as it is when it cannot be (no
// VM, not a number, the function failed: logged).
func roundShown(ctx context.Context, f *Field, v any, record any) any {
	if f.roundVM == nil {
		return v
	}
	n, err := strconv.ParseFloat(strings.TrimSpace(displayText(v)), 64)
	if err != nil {
		return v
	}
	this := gad.Object(gad.Dict{})
	if m, ok := record.(map[string]any); ok {
		if o, err := gad.ToObject(m); err == nil {
			this = o
		}
	}
	places := fieldPlaces(f)
	if places < 0 {
		places = 0
	}
	named := gad.NewNamedArgs(gad.KeyValueArray{&gad.KeyValue{K: gad.Str("places"), V: gad.Int(places)}})
	ret, err := gad.NewInvoker(f.roundVM, f.RoundDisplayFunc).WithContext(ctx).
		Invoke(gad.Args{gad.Array{this, gad.Array{}, gad.Float(n)}}, named)
	if err != nil {
		log.Printf("ERROR: form: [%s] of %s: %v", MetaRoundDisplay, f.Name, err)
		return v
	}
	if ret == nil || ret == gad.Nil {
		return v
	}
	return ret.ToString()
}

// FormatDateTime writes t as the language of ctx does (the short forms of
// i18n's formatters): its date, its time of day, or both.
func FormatDateTime(ctx context.Context, t time.Time, hasDate, hasTime bool) string {
	m := i18n.GetMessages(ctx)
	f := m.DateTimeFormatter()
	switch {
	case hasDate && !hasTime:
		f = m.DateFormatter()
	case hasTime && !hasDate:
		f = m.TimeFormatter()
	}
	return f.Short()(ctx)(t)
}

// DateTimeDisplayFunc shows a date, a time or both — "date", "time",
// "cdate", "cdtime", "ctime" — as the language of the request writes it
// (the short forms of i18n's DateFormatter, TimeFormatter, DateTimeFormatter:
// pt-BR 10/03/2026, 08:05, 10/03/2026 08:05), with its prefix and suffix;
// what is not one of them, as text.
func DateTimeDisplayFunc(c *Context) h.HTMLComponent {
	if isEmpty(c.Data) {
		return emptyDisplay(c)
	}
	return h.Text(DisplayText(eventContext(c.Event), c.Field, c.Data))
}

// dateTimeLayouts are the forms a date, a time, or both are kept in, and
// what each has.
var dateTimeLayouts = []struct {
	layout      string
	date, clock bool
}{
	{time.RFC3339Nano, true, true},
	{"2006-01-02T15:04:05", true, true},
	{"2006-01-02 15:04:05", true, true},
	{"2006-01-02T15:04", true, true},
	{"2006-01-02 15:04", true, true},
	{"2006-01-02", true, false},
	{"15:04:05", false, true},
	{"15:04", false, true},
}

// parseDateTime reads v — a time.Time, or its text in one of
// dateTimeLayouts —: the time, and whether it has a date and a time of day.
func parseDateTime(v any) (t time.Time, hasDate, hasTime, ok bool) {
	if x, is := v.(time.Time); is {
		return x, true, true, true
	}
	s := strings.TrimSpace(fmt.Sprint(v))
	for _, l := range dateTimeLayouts {
		if t, err := time.Parse(l.layout, s); err == nil {
			return t, l.date, l.clock, true
		}
	}
	return
}

// fieldDecimals are the decimal places a number of f is shown with: an int's
// none — of `[places=N]`, fixed point, N —; a float's, a decimal's
// `[places=N]`, else its step's (0.01: 2); -1 for none said (the value's
// own).
func fieldDecimals(f *Field) int {
	if f == nil {
		return -1
	}
	if f.Type == "int" || f.Type == "uint" {
		// fixed point (`[places=N]`): its places; else none
		return fixedPlaces(f)
	}
	if n := fieldPlaces(f); n >= 0 {
		return n
	}
	if _, frac, ok := strings.Cut(strings.TrimSpace(f.Step), "."); ok {
		return len(strings.TrimRight(frac, "0"))
	}
	return -1
}

// contextLanguage is the language of ctx (i18n), English when it has none.
func contextLanguage(ctx context.Context) language.Tag {
	if d := i18n.DynaFromContext(ctx); d != nil {
		if t, err := language.Parse(d.GetLanguage()); err == nil {
			return t
		}
	}
	return language.English
}

// FormatNumber writes the number v — a number, or its text ("10.5") — as the
// language lang writes it, with decimals places (-1: as many as it has);
// what is not a number, as it is.
func FormatNumber(lang language.Tag, v any, decimals int) string {
	s := strings.TrimSpace(fmt.Sprint(v))
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return s
	}
	if decimals < 0 {
		decimals = 0
		if _, frac, ok := strings.Cut(strconv.FormatFloat(f, 'f', -1, 64), "."); ok {
			decimals = len(frac)
		}
	}
	// half away from zero (x/text rounds half to even: 0.5 would be 0)
	shift := math.Pow10(decimals)
	f = math.Round(f*shift) / shift
	return message.NewPrinter(lang).Sprint(xnumber.Decimal(f, xnumber.Scale(decimals)))
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
	if c.Data == nil && c.Field != nil && c.Field.Nullable {
		// optional, not said
		return emptyDisplay(c)
	}
	on, _ := c.Data.(bool)
	if on {
		return v.VIcon("mdi-check").Color("primary").Size(v.SizeSmall)
	}
	return v.VIcon("mdi-minus").Color("grey").Size(v.SizeSmall)
}

// FileDisplayFunc is "file" and "image": the name of the file sent — a
// record {name, size, type, file} —, and its size.
func FileDisplayFunc(c *Context) h.HTMLComponent {
	if isEmpty(c.Data) {
		return emptyDisplay(c)
	}
	m, _ := c.Data.(map[string]any)
	if m == nil {
		return h.Text(displayText(c.Data))
	}
	name := displayText(m["name"])
	if size, ok := m["size"].(float64); ok && size > 0 {
		name += " (" + FormatSize(int64(size)) + ")"
	}
	return h.Span(name)
}

// FormatSize is a size in bytes as a person reads it: "820 B", "1.4 KB",
// "3.2 MB".
func FormatSize(n int64) string {
	switch {
	case n >= 1<<30:
		return strconv.FormatFloat(float64(n)/(1<<30), 'f', 1, 64) + " GB"
	case n >= 1<<20:
		return strconv.FormatFloat(float64(n)/(1<<20), 'f', 1, 64) + " MB"
	case n >= 1<<10:
		return strconv.FormatFloat(float64(n)/(1<<10), 'f', 1, 64) + " KB"
	}
	return strconv.FormatInt(n, 10) + " B"
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
		"int":       NumberDisplayFunc,
		"uint":      NumberDisplayFunc,
		"float":     NumberDisplayFunc,
		"decimal":   NumberDisplayFunc,
		"duration":  TextDisplayFunc,
		"time":      DateTimeDisplayFunc,
		"date":      DateTimeDisplayFunc,
		"cdate":     DateTimeDisplayFunc,
		"cdtime":    DateTimeDisplayFunc,
		"ctime":     DateTimeDisplayFunc,
		RangeType:   RangeDisplayFunc,
		"text":      LongTextDisplayFunc,
		"html":      HTMLDisplayFunc,
		"bool":      BoolDisplayFunc,
		"color":     ColorDisplayFunc,
		FileType:    FileDisplayFunc,
		ImageType:   FileDisplayFunc,
		AddressType: AddressDisplayFunc,
	}
}
