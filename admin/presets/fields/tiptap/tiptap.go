package tiptap

import (
	"context"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/vue"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

type Storer interface {
	Store(ctx context.Context, format string, data []byte, alt string) (url string, err error)
}
type StoreFunc func(ctx context.Context, format string, data []byte, alt string) (url string, err error)

func (f StoreFunc) Store(ctx context.Context, format string, data []byte, alt string) (string, error) {
	return f(ctx, format, data, alt)
}

type Builder struct {
	store             Storer
	mb                *presets.ModelBuilder
	updator           func(t *Builder, do presets.DataOperator) presets.DataOperator
	fieldNames        []string
	wrapComp          func(ctx *presets.FieldContext, comp h.HTMLComponent) h.HTMLComponent
	maxViewHeight     int
	maxEditHeight     int
	wrapReadonlyValue func(ctx *presets.FieldContext, value string) (string, error)
	wrapEditor        func(ctx *presets.FieldContext, comp *vx.VXTipTapEditorBuilder)
	validateScript    presets.FieldValidatorFunc
}

func New() *Builder {
	return &Builder{}
}

func (b *Builder) Store(saver Storer) *Builder {
	b.store = saver
	return b
}

func (b *Builder) WrapSaver(wraper func(saver Storer) Storer) *Builder {
	b.store = wraper(b.store)
	return b
}

func (b *Builder) Model(mb *presets.ModelBuilder) *Builder {
	b.mb = mb
	return b
}

func (b *Builder) Fields(fieldNames ...string) *Builder {
	b.fieldNames = fieldNames
	return b
}

func (b *Builder) Updator(f func(t *Builder, do presets.DataOperator) presets.DataOperator) *Builder {
	b.updator = f
	return b
}

func (b *Builder) WrapComp(f func(ctx *presets.FieldContext, comp h.HTMLComponent) h.HTMLComponent) *Builder {
	b.wrapComp = f
	return b
}

func (b *Builder) WrapReadonly(f func(ctx *presets.FieldContext, value string) (string, error)) *Builder {
	b.wrapReadonlyValue = f
	return b
}

func (b *Builder) WrapEditor(f func(ctx *presets.FieldContext, comp *vx.VXTipTapEditorBuilder)) *Builder {
	if b.wrapEditor == nil {
		b.wrapEditor = f
	} else {
		old := b.wrapEditor
		b.wrapEditor = func(ctx *presets.FieldContext, comp *vx.VXTipTapEditorBuilder) {
			old(ctx, comp)
			f(ctx, comp)
		}
	}
	return b
}

func (b *Builder) WithValidateScript(f presets.FieldValidatorFunc) *Builder {
	b.validateScript = f
	return b
}

func (b *Builder) MaxViewHeight(height int) *Builder {
	b.maxViewHeight = height
	return b
}

func (b *Builder) MaxEditHeight(height int) *Builder {
	b.maxEditHeight = height
	return b
}

func (b *Builder) TableComponentFunc(*presets.FieldContext, *web.EventContext) (r h.HTMLComponent) {
	return h.Td(h.Text("❮...❯"))
}

func (b *Builder) DetailComponentFunc(field *presets.FieldContext, ctx *web.EventContext) (r h.HTMLComponent) {
	value := field.StringValue()
	if len(value) == 0 {
		return nil
	}

	var (
		comp = vx.VXTipTapEditor().
			Label(field.Label).
			Lang("pt").
			Output("html").
			Readonly(true).
			ReadonlyClass(web.StdContentClass)
	)

	if b.maxViewHeight > 0 {
		comp.MaxHeight(b.maxViewHeight)
	}

	if b.wrapComp != nil {
		defer func() {
			r = b.wrapComp(field, r)
		}()
	}

	if b.wrapReadonlyValue != nil {
		var err error
		if value, err = b.wrapReadonlyValue(field, value); err != nil {
			panic(err)
		}
	}

	comp = comp.Attr("model-value", value)

	if b.wrapEditor != nil {
		b.wrapEditor(field, comp)
	}

	r = comp

	if b.wrapComp != nil {
		r = b.wrapComp(field, r)
	}

	return comp
}

func (b *Builder) AutoComponentFunc(field *presets.FieldContext, ctx *web.EventContext) (r h.HTMLComponent) {
	var (
		ro   = field.ReadOnly || !field.Mode.IsWrite()
		comp = vx.VXTipTapEditor().
			Label(field.Label).
			Lang("pt").
			Output("html").
			Readonly(ro).
			ErrorMessages(field.Errors).
			Class("mb-3")
	)

	if ro {
		if b.maxViewHeight > 0 {
			comp.MaxHeight(b.maxViewHeight)
		}
	} else if b.maxEditHeight > 0 {
		comp.MaxHeight(b.maxEditHeight)
	}

	if b.wrapEditor != nil {
		b.wrapEditor(field, comp)
	}

	if b.wrapComp != nil {
		defer func() {
			r = b.wrapComp(field, r)
		}()
	}

	if ro {
		return comp.Attr("model-value", field.StringValue())
	}

	return vue.FormField(comp).
		Value(field.FormKey, field.Value()).
		Bind()
}

func (b *Builder) ComponentFunc(mode presets.FieldMode) presets.FieldComponentFunc {
	if mode.Is(presets.EDIT, presets.NEW) {
		return b.AutoComponentFunc
	}

	if mode.Is(presets.DETAIL) && b.mb.HasDetailing() {
		return b.DetailComponentFunc
	}

	if mode.Is(presets.LIST) {
		return b.TableComponentFunc
	}

	return nil
}

func (b *Builder) SetFieldsComponentFunc(fb *presets.FieldsBuilder, mode presets.FieldMode) {
	for _, name := range b.fieldNames {
		fb.Field(name).ComponentFunc(b.ComponentFunc(mode))
	}
}

func (b *Builder) _validateScript(ed *presets.EditingBuilder) {
	for _, name := range b.fieldNames {
		ed.Field(name).Validator(b.validateScript)
	}
}

func (b *Builder) Build(mode presets.FieldMode) *Builder {
	ed := b.mb.Editing()

	if b.updator != nil && mode.IsWrite() {
		b.mb.UpdateDataOperator(func(do presets.DataOperator) presets.DataOperator {
			return b.updator(b, do)
		})
	}

	if mode.Has(presets.EDIT) {
		b.SetFieldsComponentFunc(&ed.FieldsBuilder, presets.EDIT)
		b._validateScript(ed)
	}

	if mode.Has(presets.NEW) && ed.HasCreatingBuilder() {
		nb := ed.Creating()
		b.SetFieldsComponentFunc(&nb.FieldsBuilder, presets.NEW)
		b._validateScript(nb)
	}

	if mode.Has(presets.DETAIL) && b.mb.HasDetailing() {
		b.SetFieldsComponentFunc(&b.mb.Detailing().FieldsBuilder, presets.DETAIL)
	}

	if mode.Has(presets.LIST) {
		b.SetFieldsComponentFunc(&b.mb.Listing().FieldsBuilder, presets.LIST)
	}

	return b
}
