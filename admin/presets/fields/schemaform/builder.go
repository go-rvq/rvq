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
	// Builder is the one drawing, so a type that contains other fields — a form
	// inside the form — draws them the same way.
	Builder *Builder
}

// Label is what to put on the input: the field's name, humanized.
func (c *Context) Label() string {
	return presets.HumanizeString(c.Field.Name)
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
	types map[string]ComponentFunc
}

func New() *Builder {
	return &Builder{
		types: map[string]ComponentFunc{
			DefaultType: TextComponentFunc,
			"int":       IntComponentFunc,
			"uint":      UintComponentFunc,
			"bool":      BoolComponentFunc,
			"color":     ColorComponentFunc,
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
		return b.record(schema, c, c.Value)
	}

	// the list: the sorter iterates, and each item is that record again — the
	// slot binds it to `item`, so a field of it is `item.<name>`
	return vx.VXArraySorter(
		web.Slot(b.record(schema, c, ItemVar)).
			Name("item").
			Scope("{ item, itemIndex }"),
	).
		Label(c.Label()).
		Density(v.DensityCompact).
		Attr("v-model", c.Value).
		Readonly(c.Form != nil && (c.Form.ReadOnly || !c.Form.Mode.IsWrite()))
}

// record renders the fields of one record, each bound under value.
func (b *Builder) record(schema *Schema, c *Context, value string) h.HTMLComponent {
	var comps []h.HTMLComponent

	for _, f := range schema.Fields {
		draw, ok := b.TypeFunc(f.Type)
		if !ok {
			comps = append(comps, errorComponent(fmt.Sprintf(
				"schemaform: o field %q pede o type %q, que não tem componente registrado (há: %v)",
				f.Name, f.Type, b.Types())))
			continue
		}

		comps = append(comps, draw(&Context{
			Field:   f,
			Value:   value + "." + f.Name,
			Form:    c.Form,
			Event:   c.Event,
			Builder: b,
		}))
	}

	return h.Div(comps...)
}

// TextComponentFunc is DefaultType: a text field, which is what an untyped
// field of the schema is.
func TextComponentFunc(c *Context) h.HTMLComponent {
	return v.VTextField().
		Label(c.Label()).
		Variant(v.FieldVariantUnderlined).
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
