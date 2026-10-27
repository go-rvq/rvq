package presets

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/vue"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

// BitOption is a bit of a field of bits — an unsigned integer whose bits are
// options (`type Reuse uint8`, with its constants `1 << iota`) —: its value
// and, in the language of the request, its label and what it does, both of the
// messages of the field's module.
type BitOption struct {
	Bit         uint64
	Label       func(ctx context.Context) string
	Description func(ctx context.Context) string
}

// BitsComponentFunc draws a field of bits as chips, one an option, marked
// when its bit is set; its description opens from a `?` icon at the right of
// the chip, on hover and on click (which does not toggle the chip). To edit (a new
// record, an edit), the chips toggle and the field's value is the sum of the
// bits marked; in the detail, they only show. The field's label above, its
// hint below.
func BitsComponentFunc(options ...BitOption) FieldComponentFunc {
	return func(field *FieldContext, ctx *web.EventContext) h.HTMLComponent {
		value := bitsOf(field.Value())
		edit := !field.ReadOnly && field.Mode.Dot().HasAny(EDIT, NEW)
		c := ctx.Context()

		var rows h.HTMLComponents
		for _, o := range options {
			label, desc := "", ""
			if o.Label != nil {
				label = o.Label(c)
			}
			if o.Description != nil {
				desc = o.Description(c)
			}
			set := value&o.Bit == o.Bit
			chip := VChip(h.Text(label)).Attr("data-bit", o.Bit).Class("me-2 mb-1")
			if desc != "" {
				// the icon and its menu in a span, the menu's activator: a
				// v-icon's content is its icon's name
				chip.GetHTMLTagBuilder().AppendChildren(web.Slot(
					h.Span("").Children(
						VIcon("mdi-help-circle-outline").Size(SizeSmall),
						VMenu(
							VCard(VCardText(h.Text(desc)).Class("text-body-2")).Attr("data-bit-description", o.Bit),
						).Activator("parent").OpenOnHover(true).OpenOnClick(true).Location("top").MaxWidth(320),
					).Class("ms-1 d-inline-flex").Attr("@click.stop", "() => {}").Attr("data-bit-help", o.Bit),
				).Name("append"))
			}
			if edit {
				chip.Value(o.Bit).Filter(true).Variant(VariantOutlined).Color("primary")
			} else if set {
				chip.Color("primary").Variant(VariantTonal).PrependIcon("mdi-check").Attr("data-bit-set", true)
			} else {
				chip.Variant(VariantOutlined).Class("text-medium-emphasis")
			}
			rows = append(rows, chip)
		}

		if !edit {
			return vx.VXReadonlyField(h.Div(rows...).Class("d-flex flex-wrap")).Label(field.Label)
		}

		bits := make([]string, len(options))
		for i, o := range options {
			bits[i] = fmt.Sprint(o.Bit)
		}
		group := VChipGroup(rows...).Multiple(true).Column(true).
			Attr(":model-value", fmt.Sprintf("[%s].filter((b) => ((Number(fieldValue.value) || 0) & b) === b)", strings.Join(bits, ", "))).
			Attr("@update:model-value", "(sel) => fieldValue.value = (sel || []).reduce((a, b) => a | b, 0)")
		field.LoadHint()
		return h.Div(
			h.Div(h.Text(field.Label)).Class("text-caption text-medium-emphasis"),
			vue.FormField(group).Value(field.FormKey, value),
			h.If(len(field.Errors) > 0, h.Div(h.Text(strings.Join(field.Errors, "; "))).Class("text-caption text-error")),
			h.If(field.Hint != "", h.Div(h.Text(field.Hint)).Class("text-caption text-medium-emphasis")),
		).Class("mb-4").Attr("data-bits-field", field.FormKey)
	}
}

// bitsOf is v, an integer of bits, as an uint64; 0 for what is not one.
func bitsOf(v any) uint64 {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rv.Uint()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if n := rv.Int(); n > 0 {
			return uint64(n)
		}
	}
	return 0
}
