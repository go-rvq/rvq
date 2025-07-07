package l10n

import "github.com/qor5/admin/v3/presets"

func ModelWithOptions(mb *presets.ModelBuilder, cfg *ModelLocalizeOptions) *presets.ModelBuilder {
	mb.SetData(LocalizeOptions, cfg)
	return mb
}
