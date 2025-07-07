package presets

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
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
			val = n.Model().NewModel()
		}
		modifiedIndexes := ContextModifiedIndexesBuilder(ctx)
		fieldInfo := n.mb.Info().ChildOf(field.ModelInfo, field.Obj)
		body := n.fb.toComponentWithFormValueKey(field.ToComponentOptions, fieldInfo, val, field.Mode, field, modifiedIndexes, ctx)
		if body == nil {
			return nil
		}

		switch t := body.(type) {
		case h.HTMLComponents:
			if len(t) == 0 {
				return nil
			}
		}
		return h.Div(
			h.Label(field.Label).Class("v-label theme--light text-caption"),
			v.VCard(body).Variant("outlined").Class("mx-0 mt-1 mb-4 px-4 pb-0 pt-4"),
		)
	})
}

func (n *NestedStructBuilder) Walk(fctx *FieldContext, opts *FieldWalkHandleOptions) (s FieldWalkState) {
	fieldInfo := n.mb.Info().ChildOf(fctx.ModelInfo, fctx.Obj)
	obj := fctx.Value()
	if obj == nil {
		if opts.SkipNestedNil {
			return
		}
		obj = n.Model().NewModel()
	}
	return n.fb.walk(fieldInfo, obj, fctx.Mode, fctx.Path, fctx.FormKey, fctx.EventContext, opts)
}
