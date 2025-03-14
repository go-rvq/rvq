package l10n

import (
	"sort"

	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/i18n"
	vx "github.com/qor5/x/v3/ui/vuetifyx"
	"gorm.io/gorm"
)

func (b *Builder) LocaleCodeListDataFilter(m *presets.ModelBuilder, multiple bool) func(ctx *web.EventContext) *vx.FilterItem {
	t := m.Translator()
	return func(ctx *web.EventContext) *vx.FilterItem {
		var (
			localeCodes []string
			locales     []*vx.SelectItem
		)

		b.db.Session(&gorm.Session{}).Model(m.Model()).
			Distinct("locale_code").
			Select("locale_code").
			Pluck("locale_code", &localeCodes)

		sort.Strings(localeCodes)

		for _, loc := range localeCodes {
			if l := b.GetLocale(loc); l != nil {
				locales = append(locales, &vx.SelectItem{
					Text:  l.Label(),
					Value: loc,
				})
			}
		}

		typ := vx.ItemTypeSelect
		if multiple {
			typ = vx.ItemTypeMultipleSelect
		}

		return &vx.FilterItem{
			Key:      "locale_code",
			Label:    i18n.Translate(t, ctx.Context(), "LocaleCode"),
			ItemType: typ,
			Options:  locales,
			Folded:   true,
		}
	}
}
