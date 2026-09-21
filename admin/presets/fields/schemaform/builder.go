package schemaform

import (
	"fmt"

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
}

// FieldInfoFunc answers for the field at PATH — the names from the root down:
// `label`, `sub.title`, `links.href`. A list adds no name, so a field of the
// records in `links` is simply `links.href`.
//
// It is optional, and so is each answer: an empty Label leaves the humanized
// name, an empty Hint leaves no hint, a nil Help leaves no `?`.
type FieldInfoFunc func(ctx *web.EventContext, path string) FieldInfo

// EnumItem is one value of an enum as the form offers it: the value that goes
// into the record, and what the reader sees.
type EnumItem struct {
	// Name is what goes INTO the record — the enum member's name, never the
	// number behind it.
	Name string
	// Label is what the reader sees; empty shows the Name.
	Label string
	// Hint is the line under the item in the list, when it needs one.
	Hint string
}

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
// as usual.
type EnumItemsFunc func(ctx *web.EventContext, path string, field *Field) []EnumItem

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

	// noLabel marks the item of a list of PLAIN VALUES: the list already
	// carries the words, so a label on every row would only repeat them.
	noLabel bool

	info     *FieldInfo
	infoRead bool

	items     []EnumItem
	itemsRead bool
}

// Info is what the builder's FieldInfoFunc says about this field, asked once.
func (c *Context) Info() FieldInfo {
	if !c.infoRead {
		c.infoRead = true
		if c.Builder != nil && c.Builder.info != nil {
			i := c.Builder.info(c.Event, c.Path)
			c.info = &i
		}
	}
	if c.info == nil {
		return FieldInfo{}
	}
	return *c.info
}

// Label is what to put on the input: what the FieldInfoFunc says, or the
// field's name humanized.
func (c *Context) Label() string {
	if c.noLabel {
		return ""
	}
	if l := c.Info().Label; l != "" {
		return l
	}
	return presets.HumanizeString(c.Field.Name)
}

// Hint is the line under the input, or "".
func (c *Context) Hint() string { return c.Info().Hint }

// Help is the long explanation, or nil.
func (c *Context) Help() h.HTMLComponent { return c.Info().Help }

// EnumItems are the values this field offers: what the builder's EnumInfoFunc
// says, or the ones the enum declared, each showing its own value.
func (c *Context) EnumItems() []EnumItem {
	if c.itemsRead {
		return c.items
	}
	c.itemsRead = true
	if c.Builder != nil {
		c.items = c.Builder.itemsOf(c.Event, c.Path, c.Field)
	}
	return c.items
}

// Builder maps a TYPE to the component that edits it, and draws a schema with
// them.
//
// A new Builder knows one type, DefaultType ("str"): a text field, which is
// also what an untyped field is. An application adds its own — and may replace
// that one:
//
//	schemaform.New().
//	    Type("int", myNumberField).
//	    Type("color", myColorPicker)
type Builder struct {
	types     map[string]ComponentFunc
	info      FieldInfoFunc
	enumInfo  EnumInfoFunc
	enumItems EnumItemsFunc
	encode    EncodeFunc
	decode    DecodeFunc
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
func (b *Builder) itemsOf(ctx *web.EventContext, path string, f *Field) []EnumItem {
	if f == nil {
		return nil
	}

	if f.Enum == nil {
		if b.enumItems == nil {
			return nil
		}
		items := b.enumItems(ctx, path, f)
		return withLabels(items)
	}

	if b.enumInfo != nil {
		if items := b.enumInfo(ctx, path); items != nil {
			return withLabels(items)
		}
	}

	items := make([]EnumItem, len(f.Enum.Names))
	for i, name := range f.Enum.Names {
		items[i] = EnumItem{Name: name, Label: name}
	}
	return items
}

// withLabels shows an item's own name when nothing else was given for it.
func withLabels(items []EnumItem) []EnumItem {
	for i := range items {
		if items[i].Label == "" {
			items[i].Label = items[i].Name
		}
	}
	return items
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
			"duration":  DurationComponentFunc,
		},
	}
}

// Type registers (or replaces) the component that draws a field of this type.
func (b *Builder) Type(name string, f ComponentFunc) *Builder {
	b.types[name] = f
	return b
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
	draw, registered := b.TypeFunc(c.Field.Type)
	if registered && c.Field.Type != DefaultType {
		return draw, true
	}
	if len(c.EnumItems()) > 0 {
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

// Types are the registered type names, FormType included.
func (b *Builder) Types() (names []string) {
	names = append(names, FormType)
	for name := range b.types {
		names = append(names, name)
	}
	return
}

// ComponentFunc draws the whole schema as the component of a presets field.
//
// The value the form edits hangs from the field's own form key: a record binds
// `form["<FormKey>"].<field>`, and a list is drawn by the array sorter, each
// item binding `<item>.<field>`.
func (b *Builder) ComponentFunc(schema *Schema) presets.FieldComponentFunc {
	return func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		if schema == nil {
			return errorComponent("schemaform: o field não tem schema")
		}
		return b.draw(schema, &Context{
			Field:   &Field{Name: field.Name, Type: FormType, Schema: schema},
			Value:   fmt.Sprintf("form[%q]", field.FormKey),
			Path:    "",
			Form:    field,
			Event:   ctx,
			Builder: b,
		})
	}
}

// formComponentFunc is FormType: a field that is itself a form, or a list of
// them. It is registered for every Builder, and it is what makes the schema
// recursive.
func (b *Builder) formComponentFunc(c *Context) h.HTMLComponent {
	if c.Field.Schema == nil {
		return errorComponent(fmt.Sprintf("schemaform: o field %q é um form sem schema", c.Field.Name))
	}
	return b.draw(c.Field.Schema, c)
}

// draw renders a schema over the value c.Value holds: the fields of a record,
// or the array sorter of a list whose item is that same record.
func (b *Builder) draw(schema *Schema, c *Context) h.HTMLComponent {
	if !schema.Slice {
		return b.record(schema, c, c.Value, c.Path)
	}

	// the list: the sorter iterates, and each item is that record again — the
	// slot binds it to `item`, so a field of it is `item.<name>`. The PATH does
	// not grow: a list adds no name, so a field of its item is `links.href`,
	// not `links.*.href` — the same thing, with less to write.
	item := b.record(schema, c, ItemVar, c.Path)
	if schema.Item != nil {
		item = b.value(schema.Item, c)
	}

	return vx.VXArraySorter(
		web.Slot(item).
			Name("item").
			Scope("{ item, itemIndex }"),
	).
		Label(c.Label()).
		Density(v.DensityCompact).
		Attr("v-model", c.Value).
		Readonly(c.Form != nil && (c.Form.ReadOnly || !c.Form.Mode.IsWrite()))
}

// value renders ONE item of a list of plain values (`[]str`): the item IS the
// value, so it binds through the list's own expression at the item's index —
// `form["x"][itemIndex]` — which a primitive needs to be written back. The path
// does not grow either: a list adds no name of its own.
func (b *Builder) value(item *Field, c *Context) h.HTMLComponent {
	ic := &Context{
		Field:   item,
		Value:   fmt.Sprintf("%s[itemIndex]", c.Value),
		Path:    c.Path,
		Form:    c.Form,
		Event:   c.Event,
		Builder: b,
		noLabel: true,
	}

	draw, ok := b.fieldFunc(ic)
	if !ok {
		return errorComponent(fmt.Sprintf(
			"schemaform: a lista pede o type %q, que não tem componente registrado nem é enum (há: %v)",
			item.Type, b.Types()))
	}
	return withHelp(ic, draw(ic))
}

// record renders the fields of one record, each bound under value and each
// reachable at its own path.
func (b *Builder) record(schema *Schema, c *Context, value, path string) h.HTMLComponent {
	var comps []h.HTMLComponent

	for _, f := range schema.Fields {
		fc := &Context{
			Field:   f,
			Value:   value + "." + f.Name,
			Path:    join(path, f.Name),
			Form:    c.Form,
			Event:   c.Event,
			Builder: b,
		}

		draw, ok := b.fieldFunc(fc)
		if !ok {
			comps = append(comps, errorComponent(fmt.Sprintf(
				"schemaform: o field %q pede o type %q, que não tem componente registrado nem é enum (há: %v)",
				f.Name, f.Type, b.Types())))
			continue
		}
		comps = append(comps, withHelp(fc, draw(fc)))
	}

	return h.Div(comps...)
}

// join is the path of a field under prefix.
func join(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "." + name
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
	items := c.EnumItems()
	if len(items) == 0 {
		return errorComponent(fmt.Sprintf(
			"schemaform: o field %q não tem valores para escolher", c.Field.Name))
	}
	options := make([]map[string]string, len(items))
	var anyHint bool
	for i, it := range items {
		options[i] = map[string]string{"value": it.Name, "title": it.Label}
		if it.Hint != "" {
			options[i]["subtitle"] = it.Hint
			anyHint = true
		}
	}

	return v.VSelect().
		Label(c.Label()).
		Variant(v.FieldVariantUnderlined).
		Items(options).
		ItemTitle("title").
		ItemValue("value").
		// with a hint to show, each item carries its own props (the subtitle)
		Attr(":item-props", anyHint).
		Hint(c.Hint()).
		PersistentHint(c.Hint() != "").
		Clearable(!c.Field.Required()).
		Attr("required", c.Field.Required()).
		Attr("v-model", c.Value)
}

// TextComponentFunc is DefaultType: a text field, which is what an untyped
// field of the schema is.
func TextComponentFunc(c *Context) h.HTMLComponent {
	return v.VTextField().
		Label(c.Label()).
		Variant(v.FieldVariantUnderlined).
		Hint(c.Hint()).
		PersistentHint(c.Hint() != "").
		Attr("required", c.Field.Required()).
		Attr("v-model", c.Value)
}

// IntComponentFunc is "int": a number field. It binds with `v-model.number`, so
// what goes back into the value is a NUMBER — the value is written out again
// (as YAML, for a LocaleMessage), and a quoted number would come back a string.
func IntComponentFunc(c *Context) h.HTMLComponent {
	return numberField(c)
}

// UintComponentFunc is "uint": the same, floored at zero.
func UintComponentFunc(c *Context) h.HTMLComponent {
	return numberField(c).Attr("min", "0")
}

func numberField(c *Context) *v.VTextFieldBuilder {
	return v.VTextField().
		Type("number").
		Label(c.Label()).
		Variant(v.FieldVariantUnderlined).
		Attr("v-model.number", c.Value)
}

// DecimalComponentFunc is "float" and "decimal", which are the same thing here:
// a gad.Decimal. It is a number field bound with a plain `v-model`, NOT
// `v-model.number` — a decimal that goes through a JS number comes back a
// float, and the precision it exists for is gone.
func DecimalComponentFunc(c *Context) h.HTMLComponent {
	return v.VTextField().
		Type("number").
		Attr("step", "any").
		Label(c.Label()).
		Variant(v.FieldVariantUnderlined).
		Attr("v-model", c.Value)
}

// LongTextComponentFunc is "text": a textarea, for prose that is not markup.
func LongTextComponentFunc(c *Context) h.HTMLComponent {
	return v.VTextarea().
		Label(c.Label()).
		Variant(v.FieldVariantUnderlined).
		AutoGrow(true).
		Rows(3).
		Attr("v-model", c.Value)
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
		Attr("v-model", c.Value)
}

// DateComponentFunc is "date": an instant with no time of day.
func DateComponentFunc(c *Context) h.HTMLComponent {
	return vx.VXDatePicker().
		Label(c.Label()).
		Attr("v-model", c.Value)
}

// DurationComponentFunc is "duration", the gad time namespace's span. It is
// written the way gad reads it back — "1h30m", "250ms" — so it is a text field:
// a picker would have to choose a unit, and the value says its own.
func DurationComponentFunc(c *Context) h.HTMLComponent {
	return v.VTextField().
		Label(c.Label()).
		Variant(v.FieldVariantUnderlined).
		Placeholder("1h30m").
		Attr("v-model", c.Value)
}

// BoolComponentFunc is "bool": a switch.
func BoolComponentFunc(c *Context) h.HTMLComponent {
	return v.VSwitch().
		Label(c.Label()).
		Density(v.DensityCompact).
		Attr("v-model", c.Value)
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
