//go:build !formscore

package schemaform

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"unicode/utf8"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

// ItemVar is what the array sorter binds each item of a list to, so a field of
// an item is `item.<name>`.
const ItemVar = "item"

// ComponentFunc draws ONE field of the form. It is what a type is: the schema
// names a type, and the type is the component that edits a value of it.
type ComponentFunc func(c *Context) h.HTMLComponent

// FieldInfo is what the schema does not say: the words around a field.
type FieldInfo struct {
	// Label replaces the humanized field name.
	Label string
	// Hint is the line under the input.
	Hint string
	// Help is the long explanation, opened from a `?` beside the field.
	Help h.HTMLComponent
	// Placeholder is the example inside the input while it is empty.
	Placeholder string
}

// FieldInfoFunc answers for the field at PATH — the names from the root down:
// `label`, `sub.title`, `links.href`. A list adds no name, so a field of the
// records in `links` is simply `links.href`.
//
// It is optional, and so is each answer: an empty Label leaves the humanized
// name, an empty Hint leaves no hint, a nil Help leaves no `?`.
type FieldInfoFunc func(ctx *web.EventContext, path string) FieldInfo

// EnumInfoFunc answers for the enum of the field at PATH — the same path
// FieldInfoFunc answers for — with ONE ENTRY PER ITEM, in the order they should
// be offered.
//
// The path, and not the enum's name, because an enum need not have one: gad
// takes it inline, as the type of a field, as well as declared beside the
// interface. The name is there in Context.Field.Enum for whoever wants it.
//
// It is optional, and so is each entry: a nil answer leaves the values the enum
// declared, and an entry with an empty Label shows its own value.
type EnumInfoFunc func(ctx *web.EventContext, path string) []EnumItem

// EnumItemsFunc answers which values the field at PATH may hold when the schema
// does not say — a list that is not fixed, so it cannot be an `enum` in the
// schema. The field itself comes along, for whoever dispatches on its type.
//
// A nil answer means the field is not a closed list, and is drawn by its type
// as usual. An ERROR is the answer when the list should have been there and
// could not be fetched: the field is then drawn as the failure it is, and a
// save that would have been checked against the list is refused — never
// accepted unchecked.
type EnumItemsFunc func(ctx *web.EventContext, path string, field *Field) ([]EnumItem, error)

// Context is what a type's component is given.
type Context struct {
	// Field is the schema's field: its name, its type, and — when the type is
	// FormType — the schema of the form inside it.
	Field *Field
	// Value is the JS expression holding the value this input edits:
	// `form["Value"].label` in a form, `item.label` inside a list. It is what
	// the component binds with v-model.
	Value string
	// Form is the presets field the whole schema is being drawn for, and Event
	// the request being answered.
	Form  *presets.FieldContext
	Event *web.EventContext
	// Path is where this field is in the schema, from the root down:
	// `links.href` for a field of the records in the `links` list. A list adds
	// no name of its own. It is what FieldInfoFunc answers for.
	Path string
	// Builder is the one drawing, so a type that contains other fields — a form
	// inside the form — draws them the same way.
	Builder *Builder

	// Compact is a field drawn in a CELL of a table: the column header carries
	// its label and hint, so the input shows neither, and keeps to one line.
	// A component that draws something taller than a line may read it too. In a
	// read view (Display) it is also the listing cell.
	Compact bool

	// Data is the value itself, in a read view (Display): decoded, as
	// DecodeValue reads it. An editing component binds Value instead.
	Data any
	// record is the record the field is of, in a read view: its Data's
	// parent (this of its `[roundDisplay=…]`).
	record any

	// Key is the JS expression of the key the value is posted under —
	// `"Value.links[" + itemIndex + "].href"` —: what the server says is
	// wrong with it is under it (ErrorAttrs). errVar is the variable of
	// the form's scope (vars) the errors are in, by key.
	Key    string
	errVar string

	// Portal is the name of a portal embedded at the root of the component —
	// the form, the detail —, shared by all its fields: where one opens what it
	// opens (a dialog, say) without depending on the page around it. "" in a
	// listing cell (Compact), which embeds none.
	Portal string

	// noLabel marks the item of a list of PLAIN VALUES: the list already
	// carries the words, so a label on every row would only repeat them.
	noLabel bool

	info     *FieldInfo
	infoRead bool

	items     []EnumItem
	itemsErr  error
	itemsRead bool
}

// PortalName is the portal embedded at the root of the component this field
// is drawn in (Context.Portal): where a component — the builder's or one of
// the application's, registered with Type or Display — opens what it opens,
// a dialog, say. "" in a listing cell, which embeds none.
func (c *Context) PortalName() string {
	return c.Portal
}

// Info is what the builder's FieldInfoFunc says about this field, asked once;
// what it leaves unsaid, the field's own metadata says — `[label="…",
// hint="…", help="…"]`, as a class field is written.
func (c *Context) Info() FieldInfo {
	if !c.infoRead {
		c.infoRead = true
		var i FieldInfo
		if c.Builder != nil && c.Builder.info != nil {
			i = c.Builder.info(c.Event, c.Path)
		}
		if c.Field != nil {
			if l, ok := c.Field.Meta[MetaLabel].(string); ok && i.Label == "" {
				i.Label = l
			}
			if hint, ok := c.Field.Meta[MetaHint].(string); ok && i.Hint == "" {
				i.Hint = hint
			}
			if p, ok := c.Field.Meta[MetaPlaceholder].(string); ok && i.Placeholder == "" {
				i.Placeholder = p
			}
			if help, ok := c.Field.Meta[MetaHelp].(string); ok && i.Help == nil && help != "" {
				i.Help = h.Div(h.Text(help)).Style("white-space:pre-wrap")
			}
		}
		c.info = &i
	}
	return *c.info
}

// Label is what to put on the input: what the FieldInfoFunc says, else the
// label of its metadata, else the field's name humanized — for the option of a
// choice, the class's name.
func (c *Context) Label() string {
	if c.noLabel || c.Compact {
		return ""
	}
	if l := c.Info().Label; l != "" {
		return l
	}
	return presets.HumanizeString(c.Field.Name)
}

// Hint is the line under the input, or "" — also in a table cell, where the
// column header carries it.
func (c *Context) Hint() string {
	if c.Compact {
		return ""
	}
	return c.Info().Hint
}

// Placeholder is the example inside the input while it is empty, or "".
func (c *Context) Placeholder() string { return c.Info().Placeholder }

// Prefix is what goes before the value (`[prefix="R$ "]`): inside its input,
// and before it shown.
func (c *Context) Prefix() string { return c.metaString(MetaPrefix) }

// Suffix is what goes after the value (`[suffix=" kg"]`), as Prefix.
func (c *Context) Suffix() string { return c.metaString(MetaSuffix) }

func (c *Context) metaString(name string) string {
	if c.Field == nil {
		return ""
	}
	s, _ := c.Field.Meta[name].(string)
	return s
}

// affixed sets the prefix and the suffix of the field on its input.
func affixed(c *Context, f *v.VTextFieldBuilder) *v.VTextFieldBuilder {
	if p := c.Prefix(); p != "" {
		f.Prefix(p)
	}
	if s := c.Suffix(); s != "" {
		f.Suffix(s)
	}
	return f
}

// LimitAttrs are the attributes of HTML5 that limit the value — min, max,
// step, minlength, maxlength, pattern, placeholder —, those the field has.
func (c *Context) LimitAttrs() []any {
	f := c.Field
	var attrs []any
	add := func(k, v string) {
		if v != "" {
			attrs = append(attrs, k, v)
		}
	}
	if p := fixedPlaces(f); p > 0 {
		// an int of places, its input shows it ÷ 10^places
		add("min", scaleDown(f.Min, p))
		add("max", scaleDown(f.Max, p))
		add("step", strconv.FormatFloat(math.Pow10(-p), 'f', p, 64))
	} else {
		add("min", f.Min)
		add("max", f.Max)
		add("step", f.Step)
	}
	if f.MinLength > 0 {
		add("minlength", strconv.Itoa(f.MinLength))
	}
	if f.MaxLength > 0 {
		add("maxlength", strconv.Itoa(f.MaxLength))
	}
	add("pattern", f.Pattern)
	placeholder := c.Placeholder()
	if f.Mask != "" {
		// the mask: as long as it, its example when there is none
		if f.MaxLength == 0 {
			add("maxlength", strconv.Itoa(utf8.RuneCountInString(f.Mask)))
		}
		if placeholder == "" {
			placeholder = MaskPlaceholder(f.Mask)
		}
	}
	add("placeholder", placeholder)
	return attrs
}

// CompactAttrs are the attributes that keep an input to one line in a table
// cell, and nothing outside one: `hide-details` (no hint, no message line) and a
// compact density. A component spreads them on its input — `.Attr(c.CompactAttrs()...)`.
func (c *Context) CompactAttrs() []any {
	if !c.Compact {
		return nil
	}
	// `:density`, bound, is the key the builders' own Density() writes — the
	// same key replaces it instead of adding a second density beside it.
	return []any{"hide-details", true, ":density", `"` + v.DensityCompact + `"`}
}

// CheckAttrs are the attributes that check the input's value here and say
// what is wrong under it: the rules of what the schema says of it (Vuetify's
// `rules`: required, its length, its pattern and mask, its min and max —
// checked as it is left, and all at once by a VForm's validate, a step's
// next), its mask written as it is typed (v-mask), and what the server found
// wrong on save (DecodeForm, by its key: `error-messages`, gone once it
// changes). A component spreads them on its input —
// `.Attr(c.CheckAttrs()...)`.
func (c *Context) CheckAttrs() []any {
	attrs := c.RuleAttrs()
	if c.Field.Mask != "" {
		attrs = append(attrs, "v-mask", strconv.Quote(c.Field.Mask))
	}
	if c.errVar == "" || c.Key == "" {
		return attrs
	}
	errs := "vars." + c.errVar
	return append(attrs,
		":error-messages", fmt.Sprintf("(%s || {})[%s]", errs, c.Key),
		"@update:model-value", fmt.Sprintf("%s && delete %s[%s]", errs, errs, c.Key),
	)
}

// RuleAttrs are the rules of the input (Vuetify's `:rules`, checked as it is
// left — validate-on blur —), from what the schema says of the field; the
// server checks the same on save (DecodeForm).
func (c *Context) RuleAttrs() []any {
	f := c.Field
	if f.ReadOnly || c.Compact {
		return nil
	}
	m := c.Messages()
	js := func(s string) string { b, _ := json.Marshal(s); return string(b) }
	const empty = "(v === null || v === undefined || String(v).trim() === '')"
	var rules []string
	add := func(cond, msg string) {
		rules = append(rules, fmt.Sprintf("v => %s || (%s) || %s", empty, cond, js(msg)))
	}
	if !f.Nullable && f.Type != "bool" && f.Type != RangeType && f.Schema == nil {
		rules = append(rules, fmt.Sprintf("v => !%s || %s", empty, js(m.Required)))
	}
	if f.MinLength > 0 {
		add(fmt.Sprintf("String(v).length >= %d", f.MinLength), fmt.Sprintf(m.TooShort, f.MinLength))
	}
	if f.MaxLength > 0 {
		add(fmt.Sprintf("String(v).length <= %d", f.MaxLength), fmt.Sprintf(m.TooLong, f.MaxLength))
	}
	if f.Pattern != "" {
		add(fmt.Sprintf("new RegExp(%s).test(String(v))", js("^(?:"+f.Pattern+")$")), m.Pattern)
	}
	if f.Mask != "" {
		// Go's \pL is JS's \p{L} (with the u flag)
		re := strings.ReplaceAll(MaskRegexp(f.Mask).String(), `\pL`, `\p{L}`)
		add(fmt.Sprintf("new RegExp(%s, 'u').test(String(v))", js(re)), fmt.Sprintf(m.Mask, f.Mask))
	}
	number := NumberTypes[f.Type]
	// an int of places (fixed point): the input shows it ÷ 10^places; its
	// limits are of the integer kept
	fixed := fixedPlaces(f)
	shown := "Number(v)"
	if fixed > 0 {
		shown = fmt.Sprintf("Math.round(Number(v) * %s)", pow10(fixed))
	}
	bound := func(limit, op, msg string) {
		if limit == "" {
			return
		}
		if number {
			said := limit
			if fixed > 0 {
				// as the input shows it
				said = scaleDown(limit, fixed)
			}
			add(fmt.Sprintf("%s %s %s", shown, op, limit), fmt.Sprintf(msg, said))
		} else if DateTypes[f.Type] {
			add(fmt.Sprintf("String(v) %s %s", op, js(limit)), fmt.Sprintf(msg, limit))
		}
	}
	bound(f.Min, ">=", m.Min)
	bound(f.Max, "<=", m.Max)
	switch {
	case fixed > 0:
		// up to its places
		add(fmt.Sprintf("/^-?\\d+([.,]\\d{0,%d})?$/.test(String(v).trim())", fixed), m.Invalid)
	case number && (f.Type == "int" || f.Type == "uint"):
		add("/^-?\\d+$/.test(String(v).trim())", m.Invalid)
	}
	if len(rules) == 0 {
		return nil
	}
	return []any{":rules", "[" + strings.Join(rules, ", ") + "]", "validate-on", "blur"}
}

// sub is the context of the field f of the record c binds (value its
// expression: c.Value, or ItemVar of an item of a list — its key indexed).
func (c *Context) sub(f *Field, value, path string) *Context {
	key := c.Key
	if key != "" && value == ItemVar {
		key += ` + "[" + itemIndex + "]"`
	}
	if key != "" {
		key += fmt.Sprintf(" + %q", "."+f.Name)
	}
	var data any
	if value != ItemVar {
		// a read-only form's value, its field's (Data: the value itself)
		data = recordValue(c.Data, f.Name)
	}
	var record any
	if value != ItemVar {
		record = c.Data
	}
	return &Context{
		Field:  f,
		Value:  value + "." + f.Name,
		Path:   join(path, f.Name),
		Data:   data,
		record: record,
		Form:   c.Form,
		Event:  c.Event, Portal: c.Portal,
		Builder: c.Builder,
		Key:     key,
		errVar:  c.errVar,
	}
}

// Help is the long explanation, or nil.
func (c *Context) Help() h.HTMLComponent { return c.Info().Help }

// EnumItems are the values this field offers: what the builder's EnumInfoFunc
// says, or the ones the enum declared, each showing its own value.
func (c *Context) EnumItems() ([]EnumItem, error) {
	if c.itemsRead {
		return c.items, c.itemsErr
	}
	c.itemsRead = true
	if c.Builder != nil {
		c.items, c.itemsErr = c.Builder.itemsOf(c.Event, c.Path, c.Field)
	}
	return c.items, c.itemsErr
}

// EnumInfo sets the function that says which values an enum field offers, and
// what each is called, BY THE FIELD'S PATH. It is optional; without it a field
// offers the values its enum declared, each showing its own name.
func (b *Builder) EnumInfo(f EnumInfoFunc) *Builder {
	b.enumInfo = f
	return b
}

// EnumItems sets the function that gives a field its values when the SCHEMA
// does not — a list that is not fixed, and so cannot be written as an `enum`:
// the locales in the database, the categories of this account, whatever the
// request has at hand.
//
// A field whose type names neither a registered component nor a declared enum
// is a select as soon as this function answers for it, the answer is what the
// select offers, and it is what a posted value is checked against
// (Builder.DecodeForm). It is optional.
func (b *Builder) EnumItems(f EnumItemsFunc) *Builder {
	b.enumItems = f
	return b
}

// itemsOf are the values the field at path may hold, or nil when it is not a
// closed list: the enum the schema declared (which EnumInfoFunc may relabel and
// reorder), or — when there is none — what EnumItemsFunc answers.
//
// It is the one place the two are resolved, so the select offers exactly what
// a posted value is checked against.
func (b *Builder) itemsOf(ctx *web.EventContext, path string, f *Field) ([]EnumItem, error) {
	if f == nil {
		return nil, nil
	}

	if f.Enum == nil {
		if b.enumItems == nil {
			return nil, nil
		}
		items, err := b.enumItems(ctx, path, f)
		if err != nil {
			return nil, err
		}
		return withLabels(items), nil
	}

	if b.enumInfo != nil {
		if items := b.enumInfo(ctx, path); items != nil {
			return withLabels(items), nil
		}
	}

	if f.Enum.Items != nil {
		return withLabels(append([]EnumItem(nil), f.Enum.Items...)), nil
	}
	items := make([]EnumItem, len(f.Enum.Names))
	for i, name := range f.Enum.Names {
		items[i] = EnumItem{Name: name, Label: name}
	}
	return items, nil
}

// FieldInfo sets the function that says the words around a field — its label,
// its hint and its help — by PATH. It is optional; without it a field is
// labelled by its name, humanized.
func (b *Builder) FieldInfo(f FieldInfoFunc) *Builder {
	b.info = f
	return b
}

func New() *Builder {
	return &Builder{
		types: map[string]ComponentFunc{
			DefaultType: TextComponentFunc,
			"int":       IntComponentFunc,
			"uint":      UintComponentFunc,
			"bool":      BoolComponentFunc,
			"color":     ColorComponentFunc,
			"float":     DecimalComponentFunc,
			"decimal":   DecimalComponentFunc,
			"text":      LongTextComponentFunc,
			"html":      HTMLComponentFunc,
			"time":      TimeComponentFunc,
			"date":      DateComponentFunc,
			// the time namespace's calendar types, with no zone: a date; a
			// date and a time of day; a time of day
			"cdate":     DateComponentFunc,
			"cdtime":    TimeComponentFunc,
			"ctime":     TimeComponentFunc,
			RangeType:   RangeComponentFunc,
			"duration":  DurationComponentFunc,
			FileType:    FileComponentFunc,
			ImageType:   FileComponentFunc,
			AddressType: AddressComponentFunc,
		},
		displays: defaultDisplays(),
	}
}

// TypeFunc is the component registered for a type, and whether there is one.
// FormType is always known: it is the builder itself, one level down.
func (b *Builder) TypeFunc(name string) (ComponentFunc, bool) {
	if name == FormType {
		return b.formComponentFunc, true
	}
	f, ok := b.types[name]
	return f, ok
}

// fieldFunc is the component that draws this field.
//
// A component REGISTERED for the field's type wins — that is what registering a
// type is for, an enum's own name included. What is left over goes to the
// select when the field holds one of a closed list of values: an enum the
// schema declared, or the list EnumItemsFunc answers with. A field the schema
// left untyped is the one case where the list wins over the registration, since
// `str` is what a field has when the schema said nothing about it.
func (b *Builder) fieldFunc(c *Context) (ComponentFunc, bool) {
	if c.Field.ReadOnly {
		return ReadOnlyComponentFunc, true
	}
	if c.Field.Enum != nil && c.Field.Enum.Options {
		// the field's own `[options=…]`: a select, whatever its type
		return EnumComponentFunc, true
	}
	draw, registered := b.TypeFunc(c.Field.Type)
	if registered && c.Field.Type != DefaultType {
		return draw, true
	}

	items, err := c.EnumItems()
	if err != nil {
		// The values this field may hold should have been there and were not:
		// say so where the field would be, instead of drawing a field that
		// cannot offer them.
		return func(c *Context) h.HTMLComponent {
			return errorComponent(fmt.Sprintf(c.Messages().ItemsUnavailable, c.Field.Name, err))
		}, true
	}
	if len(items) > 0 {
		return EnumComponentFunc, true
	}
	if registered {
		return draw, true
	}
	if c.Field.Enum != nil {
		return EnumComponentFunc, true // an enum that declared no value
	}
	return nil, false
}

// ComponentFunc draws the whole schema as the component of a presets field.
//
// The value the form edits hangs from the field's own form key: a record binds
// `form["<FormKey>"].<field>`, and a list is drawn by the array sorter, each
// item binding `<item>.<field>`.
func (b *Builder) ComponentFunc(schema *Schema) presets.FieldComponentFunc {
	return func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		if schema == nil {
			return errorComponent(messagesOf(ctx).NoSchema)
		}
		portal := NewPortalName()
		// what the server found wrong on save, by the keys under the field
		// (presets: FieldContext.SubErrors), in the scope: each input says its
		// own (Context.ErrorAttrs)
		errVar := errVarOf(field.FormKey)
		errs, _ := json.Marshal(field.SubErrors)
		if field.SubErrors == nil {
			errs = []byte("{}")
		}
		// a form not to be written shows its values (readOnlyField): the
		// field's, read here
		var data any
		if showsValues(field) {
			data = b.readOnlyData(field)
		}
		return h.Components(
			web.Scope(b.draw(schema, &Context{
				Field:   &Field{Name: field.Name, Type: FormType, Schema: schema},
				Value:   fmt.Sprintf("form[%q]", field.FormKey),
				Path:    "",
				Data:    data,
				Form:    field,
				Event:   ctx,
				Builder: b,
				Portal:  portal,
				Key:     strconv.Quote(field.FormKey),
				errVar:  errVar,
			})).Setup(fmt.Sprintf("vars.%s = %s", errVar, errs)),
			web.Portal().Name(portal),
		)
	}
}

// showsValues says the form of the field shows its values, not to be
// written: the field read-only, or drawn in a mode that does not write (a
// detail) — none said (a preview) still draws the inputs.
func showsValues(field *presets.FieldContext) bool {
	return field != nil && (field.ReadOnly || len(field.Mode) > 0 && !field.Mode.IsWrite())
}

// readOnlyField is a field of a read-only form: its label, its value by the
// display of its type (displayFieldFunc, as a detail shows it), its
// warnings (fieldWarnings), its hint.
func (b *Builder) readOnlyField(c *Context, warnings []string) h.HTMLComponent {
	var hint h.HTMLComponent
	if t := c.Hint(); t != "" {
		hint = h.Div(h.Text(t)).Class("text-caption text-medium-emphasis mt-1")
	}
	return h.Div(
		h.Div(h.Text(displayLabel(c))).Class("v-label theme--light text-caption"),
		h.Div(b.showField(c)).Class("pt-1"),
		warningsComponent(warnings),
		hint,
	).Class("mb-3").Attr("data-schemaform-readonly", c.Path)
}

// readOnlyData is the value a read-only form shows (Context.Data): the
// field's — its override, else its record's —, decoded when it is text.
func (b *Builder) readOnlyData(field *presets.FieldContext) (v any) {
	switch {
	case field.ValueOverride != nil:
		v = field.ValueOverride
	case field.Obj != nil:
		func() {
			defer func() { _ = recover() }()
			v = field.RawValue()
		}()
	}
	if s, ok := v.(string); ok {
		v, _ = b.DecodeValue(s)
	}
	return
}

// errVarOf is the variable of the errors of the field formKey in the scope.
func errVarOf(formKey string) string {
	return "sfErr_" + strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' {
			return r
		}
		return '_'
	}, formKey)
}

var portalSeq atomic.Uint64

// NewPortalName is a new name for the portal a component's root embeds
// (Context.Portal): one of its own, so two schema forms on a page never share
// it. A root drawn out of the builder's own — a component of the application
// that draws fields itself — makes one, embeds `web.Portal().Name(name)`, and
// gives it to the contexts it draws with.
func NewPortalName() string {
	return fmt.Sprintf("schemaform_portal_%d", portalSeq.Add(1))
}

// formComponentFunc is FormType: a field that is itself a form, or a list of
// them. It is registered for every Builder, and it is what makes the schema
// recursive.
func (b *Builder) formComponentFunc(c *Context) h.HTMLComponent {
	if c.Field.Schema == nil {
		return errorComponent(fmt.Sprintf(c.Messages().FormWithoutSchema, c.Field.Name))
	}
	body := b.draw(c.Field.Schema, c)
	// A record inside the record is a group: its label over its fields, set in
	// — a list carries its label in the sorter already, a choice in its select.
	if label := c.Label(); !c.Field.Schema.Slice && !c.Field.Schema.Choice && label != "" {
		return h.Div(
			h.Div(h.Text(label)).Class("text-subtitle-2 mb-2"),
			h.Div(body).Class(nestedRecordClass),
		).Class("mb-3")
	}
	return body
}

// draw renders a schema over the value c.Value holds: the fields of a record,
// or the array sorter of a list whose item is that same record.
func (b *Builder) draw(schema *Schema, c *Context) h.HTMLComponent {
	if schema.Choice {
		return b.choice(schema, c)
	}
	if !schema.Slice {
		return b.record(schema, c, c.Value, c.Path)
	}
	// a list of the values of an enum (`[]Kind`): one autocomplete of them
	if schema.Item != nil {
		if ic := b.itemContext(schema.Item, c); b.isEnum(ic) {
			return withHelp(c, EnumListComponentFunc(c, ic))
		}
	}
	switch schema.Layout {
	case LayoutTable:
		return b.table(schema, c)
	case LayoutGrid:
		return b.grid(schema, c)
	}

	// The list: each item is that record again, bound to `item`, so a field of
	// it is `item.<name>`. The PATH does not grow: a list adds no name, so a
	// field of its item is `links.href`, not `links.*.href` — the same thing,
	// with less to write.
	item := b.record(schema, c, ItemVar, c.Path)
	if schema.Item != nil {
		item = b.value(schema.Item, c)
	}

	// The DEFAULT slot is the list as it is edited — the sorter draws its own
	// rows (a drag handle and a title) only while sorting, through the `item`
	// slot, which is why the iteration is ours: `item` and `itemIndex` are the
	// names the fields above bind through.
	//
	// The sorter sorts; adding and removing an item are ours too.
	readOnly := c.Form != nil && (c.Form.ReadOnly || !c.Form.Mode.IsWrite())
	var row h.HTMLComponent
	if schema.Item == nil {
		// A record is a card: its form, and a button to remove it.
		row = listCard(item, c, readOnly)
	} else {
		row = h.Div(item).Class("flex-grow-1 px-2")
		if !readOnly {
			row = h.Div(
				h.Div(item).Class("flex-grow-1 px-2"),
				v.VBtn("").Icon("mdi-delete-outline").
					Variant(v.VariantText).Size(v.SizeSmall).Color("error").
					Class("mt-2 me-2").
					Attr("@click", fmt.Sprintf("%s.splice(itemIndex, 1)", c.Value)),
			).Class("d-flex align-start ga-2")
		}
	}

	// Between one value and the next, a line; a record is a card, with edges
	// of its own, and a space. Before each item but the first: the separator
	// is BETWEEN them.
	sep := h.HTMLComponent(v.VDivider().Attr("v-if", "itemIndex > 0").Class("my-2"))
	if schema.Item == nil {
		sep = h.Div().Attr("v-if", "itemIndex > 0").Class("mt-3")
	}

	sorter := vx.VXArraySorter(
		web.Slot(
			h.Div(sep, row).Attr("v-for", fmt.Sprintf("(%s, itemIndex) in %s", ItemVar, c.Value)),
			addItemButton(schema, c, readOnly),
		).Name("default"),
	).
		Label(c.Label()).
		Density(v.DensityCompact).
		Attr("v-model", c.Value).
		Readonly(readOnly)

	// While sorting, the row is the sorter's own: a handle, the position and a
	// title. A record has no `title`, so the title is the first field of the
	// schema that holds plain text — enough to tell one row from another.
	if title := firstTextField(schema); title != "" {
		// A static prop: the name of the field, not an expression to evaluate
		// in the form's scope.
		sorter.Attr("item-title", title)
	}
	return sorter
}

// value renders ONE item of a list of plain values (`[]str`): the item IS the
// value, so it binds through the list's own expression at the item's index —
// `form["x"][itemIndex]` — which a primitive needs to be written back. The path
// does not grow either: a list adds no name of its own.
func (b *Builder) value(item *Field, c *Context) h.HTMLComponent {
	ic := b.itemContext(item, c)
	draw, ok := b.fieldFunc(ic)
	if !ok {
		return errorComponent(fmt.Sprintf(c.Messages().ListTypeUnknown, item.Type, b.Types()))
	}
	return withHelp(ic, draw(ic))
}

// itemContext is the context of the item of a list of plain values: bound
// through the list's expression at the item's index, at the list's path.
func (b *Builder) itemContext(item *Field, c *Context) *Context {
	key := c.Key
	if key != "" {
		key += ` + "[" + itemIndex + "]"`
	}
	return &Context{
		Field: item,
		Value: fmt.Sprintf("%s[itemIndex]", c.Value),
		Path:  c.Path,
		Form:  c.Form,
		Event: c.Event, Portal: c.Portal,
		Builder: b,
		noLabel: true,
		Key:     key,
		errVar:  c.errVar,
	}
}

// isEnum reports whether the field of c is drawn as one of a closed list of
// values (EnumComponentFunc), as fieldFunc decides it.
func (b *Builder) isEnum(c *Context) bool {
	if c.Field.ReadOnly {
		return false
	}
	if c.Field.Enum != nil && c.Field.Enum.Options {
		return true
	}
	_, registered := b.TypeFunc(c.Field.Type)
	if registered && c.Field.Type != DefaultType {
		return false
	}
	if items, err := c.EnumItems(); err == nil && len(items) > 0 {
		return true
	}
	return !registered && c.Field.Enum != nil
}

// record renders the fields of one record, each bound under value and each
// reachable at its own path.
func (b *Builder) record(schema *Schema, c *Context, value, path string) h.HTMLComponent {
	readOnly := showsValues(c.Form)
	field := func(f *Field) h.HTMLComponent {
		fc := c.sub(f, value, path)

		var comp h.HTMLComponent
		if readOnly && (f.Schema == nil || f.Schema.Slice) {
			// not to be written: its value by the display of its type, as a
			// detail shows it, in its label and its hint; a record inside it
			// drawn as the form draws it (its cards, its steps)
			comp = withHelp(fc, b.readOnlyField(fc, fieldWarnings(c.Data, f.Name)))
		} else {
			draw, ok := b.fieldFunc(fc)
			if !ok {
				return errorComponent(fmt.Sprintf(c.Messages().FieldTypeUnknown, f.Name, f.Type, b.Types()))
			}
			comp = withHelp(fc, draw(fc))
		}
		if cond := whenCondition(f, value); cond != "" {
			comp = h.Div(comp).Attr("v-if", cond)
		}
		return comp
	}
	// a group of tabs or steps opens at the first whose fields the server
	// found something wrong with (a save refused: as back, at the step)
	tabAt := func(row []*Field) string {
		if c.errVar == "" {
			return "0"
		}
		keys := make([]string, len(row))
		for i, f := range row {
			keys[i] = c.sub(f, value, path).Key
		}
		return fmt.Sprintf("Math.max(0, [%s].findIndex(k => Object.keys(vars.%s || {}).some(e => e === k || e.startsWith(k + '.') || e.startsWith(k + '['))))",
			strings.Join(keys, ", "), c.errVar)
	}
	// its warnings (WarningsField) are not a field to fill: what was kept
	// unchecked, told apart by who shows the record
	return h.Div(drawRows(shownFields(schema.Fields), field, tabAt)...)
}

// choice draws a choice of one class (Schema.Choice): a select of the classes,
// each by its label with its hint under it, and under the select the form of
// the one chosen — its hint the chosen class's, else the field's. The value is
// an object of one key, the class: choosing another starts it empty — as its
// form is —, choosing the one there keeps what it holds, and clearing it (an
// optional field) leaves nothing.
func (b *Builder) choice(schema *Schema, c *Context) h.HTMLComponent {
	chosen := fmt.Sprintf("Object.keys(%s || {})[0]", c.Value)

	hints := map[string]string{}
	options := make([]map[string]string, len(schema.Fields))
	empties := make([]string, len(schema.Fields))
	forms := make([]h.HTMLComponent, 0, len(schema.Fields))
	for i, f := range schema.Fields {
		fc := &Context{
			Field: f,
			Value: fmt.Sprintf("%s[%q]", c.Value, f.Name),
			Path:  join(c.Path, f.Name),
			Form:  c.Form,
			Event: c.Event, Portal: c.Portal,
			Builder: b,
			errVar:  c.errVar,
		}
		if c.Key != "" {
			fc.Key = c.Key + fmt.Sprintf(" + %q", "."+f.Name)
		}
		options[i] = map[string]string{"value": f.Name, "title": fc.Label()}
		if hint := fc.Hint(); hint != "" {
			options[i]["subtitle"] = hint
			hints[f.Name] = hint
		}
		empties[i] = fmt.Sprintf("%q: %s", f.Name, emptyItemJS(f.Schema))
		forms = append(forms, h.Div(b.record(f.Schema, fc, fc.Value, fc.Path)).
			Class(nestedRecordClass).
			Attr("v-if", fmt.Sprintf("%s === %q", chosen, f.Name)))
	}

	sel := v.VSelect().
		Label(c.Label()).
		Variant(v.FieldVariantUnderlined).
		Items(options).
		ItemTitle("title").
		ItemValue("value").
		// title, value and subtitle go to each item of the list
		Attr(":item-props", "true").
		Attr(":hint", fmt.Sprintf("(%s)[%s] || %s", h.JSONString(hints), chosen, h.JSONString(c.Hint()))).
		PersistentHint(true).
		Clearable(!c.Field.Required()).
		Attr(":model-value", chosen).
		Attr("@update:model-value", fmt.Sprintf(
			"(v) => { %[1]s = v ? {[v]: (%[1]s || {})[v] || ({%[2]s})[v]} : null }",
			c.Value, strings.Join(empties, ", ")))

	return h.Div(append([]h.HTMLComponent{sel}, forms...)...).Class("mb-3")
}

// nestedRecordClass sets a record inside the record in, behind a line — the
// same in the form and in the detail.
const nestedRecordClass = "ps-5 border-s"

// MetaWhen is the field metadata that draws a field only while the record it
// is in holds some values — `[when={name: "grid"}] grid? GridConfig`: the
// field is there only while the record's name is "grid".
const MetaWhen = "when"

// whenCondition is the JS condition MetaWhen says, over the record at value;
// "" when the field is always drawn.
func whenCondition(f *Field, value string) string {
	when, ok := f.Meta[MetaWhen].(Meta)
	if !ok || len(when) == 0 {
		return ""
	}
	keys := keysOf(when)
	sort.Strings(keys)
	conds := make([]string, len(keys))
	for i, k := range keys {
		conds[i] = fmt.Sprintf("(%s || {})[%s] === %s", value, h.JSONString(k), h.JSONString(when[k]))
	}
	return strings.Join(conds, " && ")
}

// table draws a list of records as a TABLE — `[layout="table"]`: one row per
// record, one column per field Columns names (every field, when it names none),
// in that order. The header carries each column's label, with its hint behind
// it; a cell carries only the input (Context.Compact). A field left out of the
// columns keeps its value — it is in the record, only not drawn.
//
// It sits in the sorter like any list, so it sorts, adds and removes the same
// way: the sorter draws rows of its own only while sorting.
func (b *Builder) table(schema *Schema, c *Context) h.HTMLComponent {
	readOnly := c.Form != nil && (c.Form.ReadOnly || !c.Form.Mode.IsWrite())
	cols := schema.TableColumns()

	head := make([]h.HTMLComponent, 0, len(cols)+1)
	cells := make([]h.HTMLComponent, 0, len(cols)+1)

	for _, f := range cols {
		path := join(c.Path, f.Name)

		// The header is what a label and a hint are in a form.
		hc := &Context{Field: f, Path: path, Form: c.Form, Event: c.Event, Portal: c.Portal, Builder: b}
		th := h.Th(hc.Label()).Class("text-left text-no-wrap")
		if hint := hc.Hint(); hint != "" {
			th.Attr("title", hint)
		}
		head = append(head, th)

		fc := &Context{
			Field: f,
			Value: ItemVar + "." + f.Name,
			Path:  path,
			Form:  c.Form,
			Event: c.Event, Portal: c.Portal,
			Builder: b,
			errVar:  c.errVar,
			Key:     cellKey(c, f),
			Compact: true,
		}
		var cell h.HTMLComponent
		if draw, ok := b.fieldFunc(fc); ok {
			cell = draw(fc)
		} else {
			cell = errorComponent(fmt.Sprintf(c.Messages().FieldTypeUnknown, f.Name, f.Type, b.Types()))
		}
		cells = append(cells, h.Td(cell).Class("py-1"))
	}

	if !readOnly {
		head = append(head, h.Th("").Style("width: 1%"))
		cells = append(cells, h.Td(
			v.VBtn("").Icon("mdi-delete-outline").
				Variant(v.VariantText).Size(v.SizeSmall).Color("error").
				Attr("@click", fmt.Sprintf("%s.splice(itemIndex, 1)", c.Value)),
		))
	}

	table := v.VTable(
		h.Thead(h.Tr(head...)),
		h.Tbody(
			h.Tr(cells...).Attr("v-for", fmt.Sprintf("(%s, itemIndex) in %s", ItemVar, c.Value)),
		),
	).Density(v.DensityCompact)

	sorter := vx.VXArraySorter(
		web.Slot(table, addItemButton(schema, c, readOnly)).Name("default"),
	).
		Label(c.Label()).
		Density(v.DensityCompact).
		Attr("v-model", c.Value).
		Readonly(readOnly)
	if title := firstTextField(schema); title != "" {
		sorter.Attr("item-title", title)
	}
	return sorter
}

// grid draws a list of records as cards, schema.GridColumns() side by side (one
// under the other on a narrow screen): each card is the record's form, with a
// button to remove it; the list is sorted by the same sorter as the others.
func (b *Builder) grid(schema *Schema, c *Context) h.HTMLComponent {
	readOnly := c.Form != nil && (c.Form.ReadOnly || !c.Form.Mode.IsWrite())

	grid := v.VRow(
		gridCol(schema,
			listCard(b.record(schema, c, ItemVar, c.Path), c, readOnly).Class("h-100"),
		).Attr("v-for", fmt.Sprintf("(%s, itemIndex) in %s", ItemVar, c.Value)),
	).Dense(true)

	sorter := vx.VXArraySorter(
		web.Slot(grid, addItemButton(schema, c, readOnly)).Name("default"),
	).
		Label(c.Label()).
		Density(v.DensityCompact).
		Attr("v-model", c.Value).
		Readonly(readOnly)
	if title := firstTextField(schema); title != "" {
		sorter.Attr("item-title", title)
	}
	return sorter
}

// listCard is one record of a list as a card: its form and, unless the list is
// read only, a button to remove it — the default list's and the grid's.
func listCard(record h.HTMLComponent, c *Context, readOnly bool) *v.VCardBuilder {
	body := []h.HTMLComponent{v.VCardText(record)}
	if !readOnly {
		body = append(body, v.VCardActions(
			v.VSpacer(),
			v.VBtn("").Icon("mdi-delete-outline").
				Variant(v.VariantText).Size(v.SizeSmall).Color("error").
				Attr("@click", fmt.Sprintf("%s.splice(itemIndex, 1)", c.Value)),
		))
	}
	return v.VCard(body...).Variant(v.VariantOutlined)
}

// gridCol is one cell of a grid: the whole width on a phone, two per row on a
// small screen, and schema.GridColumns() per row from a medium one on.
func gridCol(schema *Schema, comp h.HTMLComponent) *v.VColBuilder {
	n := schema.GridColumns()
	sm := 6
	if n == 1 {
		sm = 12
	}
	return v.VCol(comp).Cols(12).Sm(sm).Md(MaxGridColumns / n)
}

// addItemButton appends one empty item of the schema's own shape to the list.
func addItemButton(schema *Schema, c *Context, readOnly bool) h.HTMLComponent {
	if readOnly {
		return nil
	}
	return v.VBtn(c.Messages().Add).
		PrependIcon("mdi-plus").
		Variant(v.VariantTonal).
		Size(v.SizeSmall).
		Class("ma-2").
		Attr("@click", fmt.Sprintf("%s.push(%s)", c.Value, emptyItemJS(schema)))
}

// cellKey is the key of the field f of an item of the list c draws as a
// table.
func cellKey(c *Context, f *Field) string {
	if c.Key == "" {
		return ""
	}
	return c.Key + ` + "[" + itemIndex + "]"` + fmt.Sprintf(" + %q", "."+f.Name)
}

// emptyItemJS is a new item of the list, written as the JS the button pushes:
// a record with each of its fields empty, or the empty value of the type a list
// of plain values holds — so a new row opens with the shape the rest has.
func emptyItemJS(schema *Schema) string {
	if schema.Item != nil {
		return emptyValueJS(schema.Item)
	}

	parts := make([]string, 0, len(schema.Fields))
	for _, f := range schema.Fields {
		value := emptyValueJS(f)
		if f.Schema != nil {
			value = emptyItemJS(f.Schema)
			if f.Schema.Slice {
				value = "[]"
			}
		}
		parts = append(parts, fmt.Sprintf("%q: %s", f.Name, value))
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

// emptyValueJS is the empty value of a field's type.
func emptyValueJS(f *Field) string {
	switch f.Type {
	case "bool":
		if f.Nullable {
			return "null"
		}
		return "false"
	case "int", "uint", "float", "decimal":
		return "0"
	case AddressType:
		return "null"
	default:
		return `""`
	}
}

// firstTextField is the name of the field a row of the sorter shows while the
// list is being sorted: the one a reader would name it by — `label`, `title`,
// `name` — and otherwise the first that holds plain text. Empty when there is
// none (a list of plain values has no field at all).
func firstTextField(schema *Schema) string {
	for _, want := range []string{"label", "title", "name"} {
		for _, f := range schema.Fields {
			if f.Name == want && f.Type == DefaultType {
				return f.Name
			}
		}
	}
	for _, f := range schema.Fields {
		if f.Type == DefaultType {
			return f.Name
		}
	}
	return ""
}

// withHelp puts the field's help behind a `?` at its right, when there is one.
func withHelp(c *Context, comp h.HTMLComponent) h.HTMLComponent {
	help := c.Help()
	if help == nil {
		return comp
	}

	return h.Div(
		h.Div(comp).Class("flex-grow-1"),
		v.VMenu(
			web.Slot(
				v.VBtn("").
					Icon("mdi-help-circle-outline").
					Variant(v.VariantText).
					Size(v.SizeSmall).
					Attr("v-bind", "props"),
			).Name("activator").Scope("{ props }"),
			v.VCard(v.VCardText(help)).MaxWidth(420),
		).CloseOnContentClick(false),
	).Class("d-flex align-start ga-1")
}

// EnumComponentFunc draws a field that holds one of a closed list of values: a
// select over that list, and nothing else. The list is the enum the schema
// declared or the one EnumItemsFunc answered with — Context.EnumItems is the
// same answer the posted value is checked against.
func EnumComponentFunc(c *Context) h.HTMLComponent {
	readOnly := c.Form != nil && (c.Form.ReadOnly || !c.Form.Mode.IsWrite())
	items, err := c.EnumItems()
	if (err != nil || len(items) == 0) && readOnly {
		// not to be chosen — a record shown —, its items the ones another
		// field makes (they come as it is filled): the value as it is
		return v.VTextField().
			Label(c.Label()).
			Variant(v.FieldVariantUnderlined).
			Attr(c.CompactAttrs()...).
			Readonly(true).
			Hint(c.Hint()).
			PersistentHint(c.Hint() != "").
			Attr("v-model", c.Value)
	}
	if err != nil {
		return errorComponent(fmt.Sprintf(c.Messages().ItemsUnavailable, c.Field.Name, err))
	}
	if len(items) == 0 {
		return errorComponent(fmt.Sprintf(c.Messages().NoItems, c.Field.Name))
	}
	options := make([]map[string]any, len(items))
	var anyHint bool
	for i, it := range items {
		options[i] = map[string]any{"value": itemValue(c.Field, it.Name), "title": it.Label}
		if it.Hint != "" {
			options[i]["subtitle"] = it.Hint
			anyHint = true
		}
	}

	return v.VSelect().
		Label(c.Label()).
		Variant(v.FieldVariantUnderlined).
		Attr(c.CompactAttrs()...).
		Items(options).
		ItemTitle("title").
		ItemValue("value").
		// with a hint to show, each item carries its own props (the subtitle)
		Attr(":item-props", anyHint).
		Hint(c.Hint()).
		PersistentHint(c.Hint() != "").
		Clearable(!readOnly && !c.Field.Required()).
		Readonly(readOnly).
		Attr("required", c.Field.Required()).
		Attr("v-model", c.Value).
		Attr(c.CheckAttrs()...)
}

// EnumListComponentFunc draws a list of the values of an enum (`[]Kind`, the
// list's context c, its item's ic): an autocomplete of them — Vuetify's, as
// several are chosen —, the ones chosen chips with their ×, typed in it the
// items that match, each with its hint under it. No more than its
// `[max=N]`: then the ones not chosen are disabled.
func EnumListComponentFunc(c, ic *Context) h.HTMLComponent {
	items, err := ic.EnumItems()
	if err != nil {
		return errorComponent(fmt.Sprintf(c.Messages().ItemsUnavailable, c.Field.Name, err))
	}
	if len(items) == 0 {
		return errorComponent(fmt.Sprintf(c.Messages().NoItems, c.Field.Name))
	}
	options := make([]map[string]any, len(items))
	for i, it := range items {
		options[i] = map[string]any{"value": itemValue(ic.Field, it.Name), "title": it.Label}
		if it.Hint != "" {
			options[i]["subtitle"] = it.Hint
		}
	}
	readOnly := c.Form != nil && (c.Form.ReadOnly || !c.Form.Mode.IsWrite())
	// each item its own props — its subtitle —, disabled when the list is
	// full and it is not one of them
	itemProps := "true"
	if max := c.Field.MaxItems; max > 0 {
		itemProps = fmt.Sprintf("(it) => Object.assign({}, it, { disabled: (%[1]s || []).length >= %[2]d && !(%[1]s || []).includes(it.value) })", c.Value, max)
	}
	return v.VAutocomplete().
		Label(c.Label()).
		Variant(v.FieldVariantUnderlined).
		Items(options).
		ItemTitle("title").
		ItemValue("value").
		Attr(":item-props", itemProps).
		Multiple(true).
		Chips(true).
		ClosableChips(!readOnly).
		Readonly(readOnly).
		Hint(c.Hint()).
		PersistentHint(c.Hint() != "").
		Clearable(!readOnly && !c.Field.Required()).
		Attr("data-schemaform-enum-list", c.Path).
		Attr("v-model", c.Value).
		Attr(c.CheckAttrs()...)
}

// itemValue is the value an item puts in the record, typed as the field is: a
// number for a number field — what the decoder stores, and what a default or a
// saved value is —, so the select shows it chosen; the name otherwise.
func itemValue(f *Field, name string) any {
	switch f.Type {
	case "int", "uint":
		if n, err := strconv.ParseInt(name, 10, 64); err == nil {
			return n
		}
	case "float", "decimal":
		if n, err := strconv.ParseFloat(name, 64); err == nil {
			return n
		}
	}
	return name
}

// ReadOnlyComponentFunc draws a field declared `get name Type`: its value,
// read-only — whatever its type, it is shown, not edited.
func ReadOnlyComponentFunc(c *Context) h.HTMLComponent {
	return v.VTextField().
		Label(c.Label()).
		Variant(v.FieldVariantPlain).
		Attr(c.CompactAttrs()...).
		Hint(c.Hint()).
		PersistentHint(c.Hint() != "").
		Readonly(true).
		Attr("v-model", c.Value)
}

// TextComponentFunc is DefaultType: a text field, which is what an untyped
// field of the schema is.
func TextComponentFunc(c *Context) h.HTMLComponent {
	return affixed(c, v.VTextField()).
		Label(c.Label()).
		Variant(v.FieldVariantUnderlined).
		Attr(c.CompactAttrs()...).
		Hint(c.Hint()).
		PersistentHint(c.Hint() != "").
		Attr("required", c.Field.Required()).
		Attr(c.LimitAttrs()...).
		Attr("v-model", c.Value).
		Attr(c.CheckAttrs()...)
}

// IntComponentFunc is "int": a number field. It binds with `v-model.number`, so
// what goes back into the value is a NUMBER — the value is written out again
// (as YAML, for a LocaleMessage), and a quoted number would come back a string.
//
// With `[places=N]` it is fixed point: the integer is the value × 10^N — of
// money, its cents: `[prefix="R$ ", places=2] price int` keeps 1050 —; its
// input, as the form holds it, the value (10.50).
func IntComponentFunc(c *Context) h.HTMLComponent {
	if p := fixedPlaces(c.Field); p > 0 {
		return fixedPointField(c, p)
	}
	return numberField(c)
}

// UintComponentFunc is "uint": the same, floored at zero.
func UintComponentFunc(c *Context) h.HTMLComponent {
	var f *v.VTextFieldBuilder
	if p := fixedPlaces(c.Field); p > 0 {
		f = fixedPointField(c, p)
	} else {
		f = numberField(c)
	}
	if c.Field.Min == "" {
		f.Attr("min", "0")
	}
	return f
}

// fixedPointField is the input of an int of places p (fixed point): what the
// form holds of it is the value as shown — 10.50 —, its p places; the
// integer (× 10^p) is what is read (Decoder) and kept, made the shown value
// again as the form is filled (Schema.Filled).
func fixedPointField(c *Context, p int) *v.VTextFieldBuilder {
	return affixed(c, v.VTextField()).
		Type("number").
		Label(c.Label()).
		Variant(v.FieldVariantUnderlined).
		Attr(c.CompactAttrs()...).
		Hint(c.Hint()).
		PersistentHint(c.Hint() != "").
		Attr(c.LimitAttrs()...).
		Attr("v-model", c.Value).
		Attr(c.CheckAttrs()...)
}

func numberField(c *Context) *v.VTextFieldBuilder {
	return affixed(c, v.VTextField()).
		Type("number").
		Label(c.Label()).
		Variant(v.FieldVariantUnderlined).
		Attr(c.CompactAttrs()...).
		Hint(c.Hint()).
		PersistentHint(c.Hint() != "").
		Attr(c.LimitAttrs()...).
		Attr("v-model.number", c.Value).
		Attr(c.CheckAttrs()...)
}

// DecimalComponentFunc is "float" and "decimal", which are the same thing here:
// a gad.Decimal. It is a number field bound with a plain `v-model`, NOT
// `v-model.number` — a decimal that goes through a JS number comes back a
// float, and the precision it exists for is gone.
func DecimalComponentFunc(c *Context) h.HTMLComponent {
	f := affixed(c, v.VTextField()).
		Type("number").
		Label(c.Label()).
		Variant(v.FieldVariantUnderlined).
		Attr(c.CompactAttrs()...).
		Hint(c.Hint()).
		PersistentHint(c.Hint() != "").
		Attr(c.LimitAttrs()...).
		Attr("v-model", c.Value).
		Attr(c.CheckAttrs()...)
	switch n := fieldPlaces(c.Field); {
	case c.Field.Step != "":
	case n >= 0:
		// its decimal places (`[places=N]`): a step of the last of them
		f.Attr("step", strconv.FormatFloat(math.Pow10(-n), 'f', n, 64))
	default:
		f.Attr("step", "any")
	}
	return f
}

// LongTextComponentFunc is "text": a textarea, for prose that is not markup.
func LongTextComponentFunc(c *Context) h.HTMLComponent {
	return v.VTextarea().
		Label(c.Label()).
		Variant(v.FieldVariantUnderlined).
		AutoGrow(true).
		Rows(3).
		Attr(c.LimitAttrs()...).
		Attr("v-model", c.Value).
		Attr(c.CheckAttrs()...)
}

// HTMLComponentFunc is "html": the rich text editor, which is WYSIWYG and so is
// its own preview.
func HTMLComponentFunc(c *Context) h.HTMLComponent {
	return vx.VXTipTapEditor().
		Label(c.Label()).
		Output("html").
		Attr("v-model", c.Value)
}

// TimeComponentFunc is "time", the gad time namespace's instant: a date AND a
// time, so the picker is the datetime one.
func TimeComponentFunc(c *Context) h.HTMLComponent {
	return vx.VXDateTimePicker().
		Label(c.Label()).
		Attr("v-model", c.Value).
		Attr(c.CheckAttrs()...)
}

// DateComponentFunc is "date": an instant with no time of day.
func DateComponentFunc(c *Context) h.HTMLComponent {
	return vx.VXDatePicker().
		Label(c.Label()).
		Attr("v-model", c.Value).
		Attr(c.CheckAttrs()...)
}

// RangeComponentFunc is a Range[T]: its label, and two inputs of T — from
// and to — side by side, binding the value's `from` and `to` (the value made
// an object first).
func RangeComponentFunc(c *Context) h.HTMLComponent {
	bound := c.Field.Range
	if bound == nil {
		bound = &Field{Type: DefaultType}
	}
	msgs := c.Messages()
	draw := func(name, label string) h.HTMLComponent {
		f, ok := c.Builder.TypeFunc(bound.Type)
		if !ok {
			f = TextComponentFunc
		}
		sub := *c
		// a bound is optional when the range is, or it is the one left open
		// (its limits, [min, max, step], the bounds')
		b := *bound
		b.Name, b.Meta = name, Meta{"label": label}
		b.Nullable = c.Field.Nullable || c.Field.RangeOpen == name
		sub.Field = &b
		sub.Value = c.Value + "." + name
		sub.Path = c.Path + "." + name
		if c.Key != "" {
			sub.Key = c.Key + fmt.Sprintf(" + %q", "."+name)
		}
		return f(&sub)
	}
	inputs := h.Div(
		v.VRow(
			v.VCol(draw(RangeFrom, msgs.RangeFrom)).Cols(6),
			v.VCol(draw(RangeTo, msgs.RangeTo)).Cols(6),
		).Dense(true),
	).Attr("v-if", fmt.Sprintf("%s && typeof %s === 'object'", c.Value, c.Value))
	// what is wrong with the range whole (its from after its to)
	var rangeErr h.HTMLComponent
	if c.errVar != "" && c.Key != "" {
		errs := fmt.Sprintf("(vars.%s || {})[%s]", c.errVar, c.Key)
		rangeErr = h.Div().Class("text-error text-caption").Attr("v-if", errs).
			Attr("v-text", errs+".join(' ')").Attr("data-range-error", c.Path)
	}
	out := h.Div(h.Div(h.Text(c.Label())).Class("v-label theme--light text-caption mb-1"), inputs, rangeErr).
		Attr("data-range", c.Path)
	// the value an object, to bind its bounds: made one as it mounts — on the
	// object it is a key of (`<parent>.<key>`) —, the inputs drawn then
	if i := strings.LastIndex(c.Value, "."); i > 0 && !strings.ContainsAny(c.Value[i:], "[]") {
		parent, key := c.Value[:i], c.Value[i+1:]
		out.Attr("v-assign", fmt.Sprintf("[%s, {%s: (%s && typeof %s === 'object') ? %s : {from: null, to: null}}]",
			parent, key, c.Value, c.Value, c.Value))
	}
	return out
}

// FileComponentFunc is "file" and "image": a file input of what it takes
// ([accept=…]). The value of a file is the application's to keep — what is
// sent with the form —: the admin draws it to show the form (a preview).
func FileComponentFunc(c *Context) h.HTMLComponent {
	return v.VFileInput().
		Label(c.Label()).
		Variant(v.FieldVariantUnderlined).
		Attr(c.CompactAttrs()...).
		Attr("accept", c.Field.AcceptOf()).
		Hint(c.Hint()).
		PersistentHint(c.Hint() != "")
}

// DurationComponentFunc is "duration", the gad time namespace's span. It is
// written the way gad reads it back — "1h30m", "250ms" — so it is a text field:
// a picker would have to choose a unit, and the value says its own.
func DurationComponentFunc(c *Context) h.HTMLComponent {
	return v.VTextField().
		Label(c.Label()).
		Variant(v.FieldVariantUnderlined).
		Attr(c.CompactAttrs()...).
		Placeholder("1h30m").
		Attr("v-model", c.Value).
		Attr(c.CheckAttrs()...)
}

// BoolComponentFunc is "bool": a switch, coloured when on — an optional one
// (`x? bool`), a select of yes and no that may be left empty (null); with
// [required], one of them chosen. Coloured when it is on — a switch
// with no colour stays grey either way, which is exactly what a switch must not
// be: whether it is on is the whole of what it says.
func BoolComponentFunc(c *Context) h.HTMLComponent {
	if c.Field.Nullable {
		// optional: yes, no, or not said — a select, cleared to null —;
		// [required]: yes or no, chosen
		msgs := c.Messages()
		_, must := c.Field.Meta["required"]
		return v.VSelect().
			Label(c.Label()).
			Variant(v.FieldVariantUnderlined).
			Attr(c.CompactAttrs()...).
			Items([]map[string]any{{"title": msgs.Yes, "value": true}, {"title": msgs.No, "value": false}}).
			ItemTitle("title").ItemValue("value").
			Clearable(!must).
			Attr("required", must).
			Hint(c.Hint()).
			PersistentHint(c.Hint() != "").
			Attr("v-model", c.Value).
			Attr(c.CheckAttrs()...)
	}
	return v.VSwitch().
		Label(c.Label()).
		Color("primary").
		Density(v.DensityCompact).
		Attr(c.CompactAttrs()...).
		Attr("v-model", c.Value).
		Attr(c.CheckAttrs()...)
}

// ColorComponentFunc is "color": the color picker, showing the value it holds.
func ColorComponentFunc(c *Context) h.HTMLComponent {
	return h.Div(
		h.Label(c.Label()).Class("v-label text-caption"),
		v.VColorPicker().
			Attr("v-model", c.Value).
			Attr("mode", "hexa").
			Attr("hide-inputs", "false"),
	).Class("mb-4")
}

func errorComponent(msg string) h.HTMLComponent {
	return v.VAlert(h.Text(msg)).Type("error").Variant("tonal").Density(v.DensityCompact)
}
