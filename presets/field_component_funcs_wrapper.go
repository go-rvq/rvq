package presets

import (
	"github.com/qor5/web/v3"
	h "github.com/theplant/htmlgo"
)

func ReadOnlyFieldComponentFuncWrapper(f FieldComponentFunc) FieldComponentFunc {
	l, d := ListingFieldComponentFuncWrapper(f), DetailingFieldComponentFuncWrapper(f)

	return func(field *FieldContext, ctx *web.EventContext) h.HTMLComponent {
		if field.Mode.Dot().Is(LIST) {
			return l(field, ctx)
		}
		return d(field, ctx)
	}
}

func ListingFieldComponentFuncWrapper(f FieldComponentFunc) FieldComponentFunc {
	return func(field *FieldContext, ctx *web.EventContext) (comp h.HTMLComponent) {
		if comp = f(field, ctx); comp == nil {
			return
		}
		return h.Td(comp)
	}
}

func DetailingFieldComponentFuncWrapper(f FieldComponentFunc) FieldComponentFunc {
	return FormFieldComponentFuncWrapper(func(field *FieldContext, ctx *web.EventContext) (comp h.HTMLComponent) {
		if comp = f(field, ctx); comp == nil {
			return
		}
		return h.Components(
			h.Label(field.Label).Class("v-label text-caption"),
			h.Div(comp).Class("pt-1"),
		)
	})
}

func FormFieldComponentFuncWrapper(f FieldComponentFunc) FieldComponentFunc {
	return func(field *FieldContext, ctx *web.EventContext) (comp h.HTMLComponent) {
		return FormFieldComponentWrapper(f(field, ctx))
	}
}

func FormFieldComponentWrapper(comp h.HTMLComponent) h.HTMLComponent {
	if comp == nil {
		return nil
	}
	var div *h.HTMLTagBuilder
	switch t := comp.(type) {
	case h.HTMLComponents:
		div = h.Div(t...)
	default:
		div = h.Div(comp)
	}
	return div.Class("mb-4")
}
