package profile

import (
	"github.com/go-rvq/rvq/admin/presets"
)

func DefaultModelOptions(opts ...presets.ModelBuilderOption) []presets.ModelBuilderOption {
	return append(opts, presets.ModelBuilderOptionFunc(func(mb *presets.ModelBuilder) {
		mb.SetModuleKey(MessagesKey)
	}))
}

func Model(p *presets.Builder, v any, opts ...presets.ModelBuilderOption) *presets.ModelBuilder {
	return p.Model(v, DefaultModelOptions(opts...)...)
}

func NewModel(p *presets.Builder, v any, opts ...presets.ModelBuilderOption) *presets.ModelBuilder {
	return presets.NewModelBuilder(p, v, DefaultModelOptions(opts...)...)
}
