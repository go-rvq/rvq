// Package admin mounts the organizations module into an rvq presets admin.
// It provides the Organizacao CRUD (owner-scoped) and is the base other domain
// modules (e.g. finance) depend on to scope their records to an organization.
package admin

import (
	"github.com/go-rvq/rvq/admin/packages/orgs/messages"
	"github.com/go-rvq/rvq/admin/presets"
	"gorm.io/gorm"
)

// Configure mounts the organizations module on b (via the orgs Builder plugin).
func Configure(b *presets.Builder, db *gorm.DB) {
	b.Use(New(db))
}

// DefaultModelOptions appends the options shared by every org model.
func DefaultModelOptions(opts ...presets.ModelBuilderOption) []presets.ModelBuilderOption {
	return append(opts, presets.ModelBuilderOptionFunc(func(mb *presets.ModelBuilder) {
		mb.SetModuleKey(messages.Key)
	}))
}

// Model registers v on p as an org model (i18n module key applied).
func Model(p *presets.Builder, v any, opts ...presets.ModelBuilderOption) *presets.ModelBuilder {
	return p.Model(v, DefaultModelOptions(opts...)...)
}
