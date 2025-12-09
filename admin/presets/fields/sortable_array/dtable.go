package sortable_array

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/vue"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

type Builder struct {
	mb         *presets.ModelBuilder
	wrapComp   func(ctx *presets.FieldContext, comp h.HTMLComponent) h.HTMLComponent
	wrapSorter func(ctx *presets.FieldContext, sorter *vx.VXArraySorterBuilder) h.HTMLComponent
	comp       func(itemsVar string, ctx *presets.FieldContext) h.HTMLComponent
	fieldName  string
	items      func(ctx *presets.FieldContext) (items any, err error)
}

func New(comp func(itemsVar string, ctx *presets.FieldContext) h.HTMLComponent) *Builder {
	return &Builder{
		comp: comp,
	}
}

func (b *Builder) Comp(f func(itemsVar string, ctx *presets.FieldContext) h.HTMLComponent) *Builder {
	b.comp = f
	return b
}

func (b *Builder) Model(mb *presets.ModelBuilder) *Builder {
	b.mb = mb
	return b
}

func (b *Builder) Field(fieldName string) *Builder {
	b.fieldName = fieldName
	return b
}

func (b *Builder) WrapComp(f func(ctx *presets.FieldContext, comp h.HTMLComponent) h.HTMLComponent) *Builder {
	b.wrapComp = f
	return b
}

func (b *Builder) WrapSorter(f func(ctx *presets.FieldContext, sorter *vx.VXArraySorterBuilder) h.HTMLComponent) *Builder {
	b.wrapSorter = f
	return b
}

func (b *Builder) Items(items func(ctx *presets.FieldContext) (items any, err error)) *Builder {
	b.items = items
	return b
}

func (b *Builder) TableComponentFunc(*presets.FieldContext, *web.EventContext) (r h.HTMLComponent) {
	return h.Td(v.VIcon("mdi-table"))
}

func (b *Builder) DetailComponentFunc(field *presets.FieldContext, ctx *web.EventContext) (r h.HTMLComponent) {
	var (
		comp = v.VDataTable()
	)

	if b.wrapComp != nil {
		defer func() {
			r = b.wrapComp(field, r)
		}()
	}

	comp = comp.Attr("model-value", field.StringValue())

	if b.wrapComp != nil {
		r = b.wrapComp(field, r)
	}

	return vue.FormField(comp).
		Value(field.FormKey, field.Value()).
		Bind()
}

func (b *Builder) AutoComponentFunc(field *presets.FieldContext, ctx *web.EventContext) (r h.HTMLComponent) {
	var (
		ro = field.ReadOnly || !field.Mode.IsWrite()

		items   any
		formKey string
		err     error
	)

	if items, err = b.items(field); err != nil {
		panic(err)
	}

	if b.wrapComp != nil {
		defer func() {
			r = b.wrapComp(field, r)
		}()
	}

	if ro {
		formKey = "items"
	} else {
		formKey = fmt.Sprintf("form[%q]", field.FormKey)
	}

	sorter := vx.VXArraySorter(b.comp(formKey, field)).
		Label(field.Label).
		Density(v.DensityCompact).
		Attr("v-model", formKey).
		Readonly(ro)

	if ro {
		return vue.UserComponent(sorter).Scope("items", items)
	}

	r = sorter

	if b.wrapSorter != nil {
		r = b.wrapSorter(field, sorter)
	}

	return vue.UserComponent(r).Assign("form", field.FormKey, items)
}

func (b *Builder) ComponentFunc(mode presets.FieldMode) presets.FieldComponentFunc {
	if mode.HasAny(presets.EDIT, presets.NEW) {
		return b.AutoComponentFunc
	}

	if mode.Has(presets.DETAIL) && b.mb.HasDetailing() {
		return b.DetailComponentFunc
	}

	if mode.Has(presets.LIST) {
		return b.TableComponentFunc
	}

	return nil
}

func (b *Builder) SetFieldsComponentFunc(fb *presets.FieldsBuilder, mode presets.FieldMode) {
	fb.Field(b.fieldName).ComponentFunc(b.ComponentFunc(mode))
}

func (b *Builder) Build() *Builder {
	return b.BuildMode(presets.ALL)
}

func (b *Builder) BuildMode(mode presets.FieldMode) *Builder {
	ed := b.mb.Editing()

	if mode.Has(presets.EDIT) {
		ed.WrapPostSetterFunc(func(in presets.SetterFunc) presets.SetterFunc {
			return func(obj interface{}, ctx *web.EventContext) {
				in(obj, ctx)
			}
		})
		b.SetFieldsComponentFunc(&ed.FieldsBuilder, mode)
	}

	if mode.Has(presets.NEW) && ed.HasCreatingBuilder() {
		b.SetFieldsComponentFunc(&ed.Creating().FieldsBuilder, mode)
	}

	if mode.Has(presets.DETAIL) && b.mb.HasDetailing() {
		b.SetFieldsComponentFunc(&b.mb.Detailing().FieldsBuilder, mode)
	}

	if mode.Has(presets.LIST) {
		b.SetFieldsComponentFunc(&b.mb.Listing().FieldsBuilder, mode)
	}

	return b
}
