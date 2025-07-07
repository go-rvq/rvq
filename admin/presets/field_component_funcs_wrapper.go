package presets

import (
	"fmt"
	"strconv"

	"github.com/qor5/web/v3"
	"github.com/qor5/web/v3/vue"
	. "github.com/qor5/x/v3/ui/vuetify"
	h "github.com/theplant/htmlgo"
)

func ReadOnlyFieldComponentFuncWrapper(f FieldComponentFunc) FieldComponentFunc {
	l, d := ListingFieldComponentFuncWrapper(f), FieldComponentWrapper(f)

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

func FieldComponentContainer(label any, comp h.HTMLComponent) h.HTMLComponent {
	if comp == nil {
		return nil
	}

	var labelComp *h.HTMLTagBuilder
	switch t := label.(type) {
	case string:
		labelComp = h.Label(t).Class("v-label text-caption")
	case *h.HTMLTagBuilder:
		labelComp = t
	case h.HTMLComponent:
		labelComp = h.Div(t).Class("v-label text-caption")
	}
	if labelComp == nil {
		return comp
	}

	return h.Components(
		labelComp,
		h.Div(comp).Class("pt-1 mb-4"),
	)
}

func FieldComponentWrapper(f FieldComponentFunc) FieldComponentFunc {
	return FormFieldComponentFuncWrapper(func(field *FieldContext, ctx *web.EventContext) (comp h.HTMLComponent) {
		comp = f(field, ctx)
		return FieldComponentContainer(field.Label, comp)
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

func SliceFieldComponentFuncWraper(draggable bool, itemComponentFunc func(field *FieldContext, itemsModel, deleteFunc string) h.HTMLComponent, wrap ...func(fctx *FieldContext, comp h.HTMLComponent) h.HTMLComponent) func(field *FieldContext, _ *web.EventContext) h.HTMLComponent {
	const (
		itemsModel = "fieldValue"
		deleteFunc = "removeItem"
	)

	return FieldComponentWrapper(
		func(field *FieldContext, _ *web.EventContext) h.HTMLComponent {
			v := field.RawValue()
			if v == nil {
				v = []string{}
			}

			var errComp h.HTMLComponent

			if len(field.Errors) > 0 {
				errComp = VAlert(h.Text(field.Errors[0])).Density(DensityCompact).Class("mb-2").Type("error")
			}

			formKey := strconv.Quote(field.FormKey)

			comp := itemComponentFunc(field, itemsModel, deleteFunc+"(%s)")
			if draggable {
				comp = h.Tag("vx-draggable",
					web.Slot(comp).
						Name("item").
						Scope("{element, index}"),
				).
					Attr("handle", ".handle", "animation", "300", "item-key", "$id").
					Attr("@end", "moveItem($event.oldIndex, $event.newIndex)").
					Attr("v-model", itemsModel)
			}

			comp = vue.UserComponent(
				errComp,
				h.Div(
					comp,
				),
				h.Div(
					VSpacer(),
					VBtn("").
						Icon("mdi-plus").
						Density(DensityComfortable).
						Attr("@click", "addItem"),
					VSpacer(),
				).Class("d-flex"),
			).
				Assign("form", field.FormKey, v).
				Scope(deleteFunc).
				Scope("addItem").
				Scope("moveItem").
				Scope(itemsModel).
				Setup(`({ scope }) => {
form[` + formKey + `] = form[` + formKey + `] || [];
scope.` + itemsModel + ` = form[` + formKey + `];

if (scope.fieldValue.length > 0) {
	if (typeof (scope.fieldValue[0]) != "object") {
		for (let i = 0; i < scope.fieldValue.length; i++) {
			scope.fieldValue[i] = {__value: scope.fieldValue[i]}
		}
	}
}

scope.fieldValue.forEach((item, index) => item.$id = index+1)
let lastId = scope.fieldValue.length

scope.addItem = () => {
    scope.` + itemsModel + `.push({ __value: "", $id: ++lastId })
}
scope.` + deleteFunc + ` = (index) => {
	scope.` + itemsModel + `.splice(index, 1)
}
scope.moveItem = (from, to) => {
	const item = scope.` + itemsModel + `.splice(from, 1)[0]
    scope.` + itemsModel + `.splice(to, 0, item)
}
}`)

			for _, f := range wrap {
				comp = f(field, comp)
			}
			return comp
		},
	)
}

func WriteSortableSliceFieldComponentFuncWraper(fieldComp func(field *FieldContext, itemsModel, deleteFunc string) h.HTMLComponent, wrap ...func(fctx *FieldContext, comp h.HTMLComponent) h.HTMLComponent) func(field *FieldContext, _ *web.EventContext) h.HTMLComponent {
	return SliceFieldComponentFuncWraper(true, func(field *FieldContext, itemsModel, deleteFunc string) h.HTMLComponent {
		return h.Div(
			VIcon("mdi-drag").Class("handle me-2 cursor-grab"),
			fieldComp(field, itemsModel, deleteFunc),
		).Class("d-flex")
	}, wrap...)
}

func WriteSortableTextSliceFieldComponentFuncWraper(wrap ...func(fctx *FieldContext, comp h.HTMLComponent) h.HTMLComponent) func(field *FieldContext, _ *web.EventContext) h.HTMLComponent {
	return WriteSortableSliceFieldComponentFuncWraper(func(field *FieldContext, itemsModel, deleteFunc string) h.HTMLComponent {
		return VTextField(
			web.Slot(
				VIcon("mdi-delete").
					Color("error").
					Attr("@click", fmt.Sprintf(deleteFunc, "index")),
			).Name("append-inner"),
		).
			Density(DensityCompact).
			HideDetails(true).
			Class("mb-2").
			Attr("v-model", "element.__value")
	}, wrap...)
}
