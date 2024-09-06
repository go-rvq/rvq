package presets

import (
	"reflect"

	"github.com/qor5/web/v3"
	v "github.com/qor5/x/v3/ui/vuetify"
	"github.com/sunfmin/reflectutils"
	h "github.com/theplant/htmlgo"
)

type NestedStructBuilder struct {
	mb *ModelBuilder
	fb *FieldsBuilder
}

func NestedStruct(mb *ModelBuilder, fb *FieldsBuilder) *NestedStructBuilder {
	return &NestedStructBuilder{mb: mb, fb: fb}
}

func (n *NestedStructBuilder) Model() *ModelBuilder {
	return n.mb
}

func (n *NestedStructBuilder) FieldsBuilder() *FieldsBuilder {
	return n.fb
}

func (n *NestedStructBuilder) Build(b *FieldBuilder) {
	b.ComponentFunc(func(field *FieldContext, ctx *web.EventContext) h.HTMLComponent {
		val := field.Value()
		if val == nil {
			t := reflectutils.GetType(field.Obj, field.Name).Elem()
			val = reflect.New(t).Interface()
		}
		modifiedIndexes := ContextModifiedIndexesBuilder(ctx)
		fieldInfo := n.mb.Info().ChildOf(field.ModelInfo, field.Obj)
		body := n.fb.toComponentWithFormValueKey(fieldInfo, val, field.Mode, field.FormKey, modifiedIndexes, ctx)
		return h.Div(
			h.Label(field.Label).Class("v-label theme--light text-caption"),
			v.VCard(body).Variant("outlined").Class("mx-0 mt-1 mb-4 px-4 pb-0 pt-4"),
		)
	})
}

func (n *NestedStructBuilder) Walk(fctx *FieldContext, handle FieldWalkHandle) (s FieldWalkState) {
	return n.fb.walk(fctx.ModelInfo, fctx.Obj, fctx.Mode, fctx.FormKey, fctx.EventContext, handle)
}
