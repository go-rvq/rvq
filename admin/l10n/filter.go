package l10n

import (
	"sort"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
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
