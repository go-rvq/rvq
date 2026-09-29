package presets

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

type NestedStructBuilder struct {
	mb    *ModelBuilder
	fb    *FieldsBuilder
	frame NestedStructFrame
}

// NestedStructFrame draws the fields of a nested struct — body, with the
// field's hint (nil when it has none) — as the field: by default its label over
// an outlined card holding them (DefaultNestedStructFrame).
type NestedStructFrame func(field *FieldContext, hint, body h.HTMLComponent) h.HTMLComponent

// DefaultNestedStructFrame is the field's label over an outlined card holding
// the hint and the fields.
func DefaultNestedStructFrame(field *FieldContext, hint, body h.HTMLComponent) h.HTMLComponent {
	return h.Div(
		h.Label(field.Label).Class("v-label theme--light text-caption"),
		v.VCard(hint, body).Variant("outlined").Class("mx-0 mt-1 mb-4 px-4 pb-0 pt-4"),
	)
}

func NestedStruct(mb *ModelBuilder, fb *FieldsBuilder) *NestedStructBuilder {
	return &NestedStructBuilder{mb: mb, fb: fb}
}

// Frame sets what the fields are drawn in (NestedStructFrame).
func (n *NestedStructBuilder) Frame(f NestedStructFrame) *NestedStructBuilder {
	n.frame = f
	return n
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
		switch t := body.(type) {
		case h.HTMLComponents:
			if len(t) == 0 {
				return nil
			}
		}

		var hintComp h.HTMLComponent
		if hint := field.GetOrLoadHint(); len(hint) > 0 {
			hintComp = h.Div(h.RawHTML(hint)).Class("input-fields-group_hint text-caption opacity-60 mb-3")
		}

		frame := n.frame
		if frame == nil {
			frame = DefaultNestedStructFrame
		}
		return frame(field, hintComp, body)
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
