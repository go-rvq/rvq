package vuetifyx

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/tag"
	"github.com/go-rvq/rvq/x/ui/vuetify"
)

type PickerBuilder struct {
	value     interface{}
	label     string
	fieldName string
	comp      h.HTMLComponent
}

// Picker wraps a picker component (a date or time picker) in a menu bound to a
// form field.
//
// It takes a plain component: the vuetify builders are generic and their
// SetAttr returns the builder, so they no longer satisfy MutableAttrHTMLComponent
// — the attribute goes through their underlying tag instead (see setAttr).
func Picker(c h.HTMLComponent) (r *PickerBuilder) {
	r = &PickerBuilder{comp: c}
	return
}

func (b *PickerBuilder) Value(v interface{}) (r *PickerBuilder) {
	b.value = v
	return b
}

func (b *PickerBuilder) Label(v string) (r *PickerBuilder) {
	b.label = v
	return b
}

func (b *PickerBuilder) FieldName(v string) (r *PickerBuilder) {
	b.fieldName = v
	return b
}

func (b *PickerBuilder) Write(ctx *h.Context) (err error) {
	menuLocal := fmt.Sprintf("picker_%s_menu", b.fieldName)
	valueLocal := fmt.Sprintf("picker_%s_value", b.fieldName)

	b.setAttr("@change", fmt.Sprintf(`locals.%s = false; locals.%s = $event; $plaid().form(plaidForm).fieldValue(%s, $event)`, menuLocal, valueLocal, h.JSONString(b.fieldName)))

	return web.Scope(
		vuetify.VMenu(
			web.Slot(
				vuetify.VTextField().
					Label(b.label).
					Readonly(true).
					PrependIcon("edit_calendar").
					Attr("v-model", fmt.Sprintf("locals.%s", valueLocal)).
					Attr("v-bind", "attrs").
					Attr("v-on", "on"),
			).Name("activator").Scope("{ on, attrs }"),

			b.comp,
		).Attr("v-model", fmt.Sprintf("locals.%s", menuLocal)).
			CloseOnContentClick(false).
			MaxWidth(290),
	).LocalsInit(fmt.Sprintf(`{%s: %s, %s: false}`, valueLocal, h.JSONString(b.value), menuLocal)).
		Slot("{ locals }").
		Write(ctx)
}

// setAttr sets an attribute on the wrapped component, whichever shape it has.
func (b *PickerBuilder) setAttr(k string, v any) {
	switch c := b.comp.(type) {
	case h.MutableAttrHTMLComponent:
		c.SetAttr(k, v)
	case tag.TagGetter:
		c.GetHTMLTagBuilder().SetAttr(k, v)
	}
}
