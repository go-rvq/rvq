package presets

import (
	"github.com/go-rvq/rvq/web"
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
	return g.fg.Field(name).ContextLabel(g.mi, g.ctx.Context())
}
