package presets

import (
	"github.com/qor5/web/v3"
)

type FieldsContextLabelGetter struct {
	mi *ModelInfo
	fg interface {
		Field(name string) *FieldBuilder
	}
	ctx web.ContextValuer
}

func NewFieldsContextLabelGetter(mi *ModelInfo, fg interface {
	Field(name string) *FieldBuilder
}, ctx web.ContextValuer) *FieldsContextLabelGetter {
	return &FieldsContextLabelGetter{mi: mi, fg: fg, ctx: ctx}
}

func (g *FieldsContextLabelGetter) Get(name string) string {
	return g.fg.Field(name).ContextLabel(g.mi, g.ctx)
}
