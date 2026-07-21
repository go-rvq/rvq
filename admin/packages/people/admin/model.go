package admin

import (
	orgsadmin "github.com/go-rvq/rvq/admin/packages/orgs/admin"
	"github.com/go-rvq/rvq/admin/packages/people/messages"
	"github.com/go-rvq/rvq/admin/presets"
)

// MenuGroupID is the identifier of the people menu group.
const MenuGroupID = "people"

// DefaultModelOptions appends the options shared by every people model.
func DefaultModelOptions(opts ...presets.ModelBuilderOption) []presets.ModelBuilderOption {
	return append(opts, presets.ModelBuilderOptionFunc(func(mb *presets.ModelBuilder) {
		mb.SetModuleKey(messages.Key)
	}))
}

// Model registers v as a people resource nested under the current organization:
// its routes live at /orgs/{orgID}/people/<resource> and its records are scoped
// to the organization from the path (see orgsadmin.MountUnderOrg). The orgs
// package must already be configured on p (Configure ensures it).
func Model(p *presets.Builder, v any, opts ...presets.ModelBuilderOption) *presets.ModelBuilder {
	mb := presets.NewModelBuilder(p, v, DefaultModelOptions(opts...)...)
	orgsadmin.MountUnderOrg(orgsadmin.OrganizacaoModel(p), mb, MenuGroupID)
	return mb
}

// NewModel builds an unregistered people model builder, for nested/child models.
func NewModel(p *presets.Builder, v any, opts ...presets.ModelBuilderOption) *presets.ModelBuilder {
	return presets.NewModelBuilder(p, v, DefaultModelOptions(opts...)...)
}
