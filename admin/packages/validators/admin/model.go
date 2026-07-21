package admin

import (
	"github.com/go-rvq/rvq/admin/packages/validators/messages"
	"github.com/go-rvq/rvq/admin/presets"
)

// MenuGroupID is the identifier of the validators menu group.
const MenuGroupID = "validators"

// ValidatorModelID is the presets model id of the Validator resource.
const ValidatorModelID = "validators"

// DefaultModelOptions appends the options shared by every validators model.
func DefaultModelOptions(opts ...presets.ModelBuilderOption) []presets.ModelBuilderOption {
	return append(opts, presets.ModelBuilderOptionFunc(func(mb *presets.ModelBuilder) {
		mb.SetModuleKey(messages.Key)
	}))
}

// Model registers v as a validators resource (i18n module key applied).
func Model(p *presets.Builder, v any, opts ...presets.ModelBuilderOption) *presets.ModelBuilder {
	return p.Model(v, DefaultModelOptions(opts...)...)
}
