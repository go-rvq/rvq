package l10n

import "github.com/go-rvq/rvq/admin/presets"

func ModelWithOptions(mb *presets.ModelBuilder, cfg *ModelLocalizeOptions) *presets.ModelBuilder {
	mb.SetData(LocalizeOptions, cfg)
	return mb
}
