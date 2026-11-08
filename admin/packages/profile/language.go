package profile

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/utils/bcp47"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

// The default language of the user (user.LanguageOf): the language of the
// admin when they open it — before the browser's, after the one they
// switch to in the menu — and of the mails sent to them. None: the
// browser's.

// languageItems are the languages of the admin to choose, none first.
func languageItems(p *presets.Builder, m *Messages) []map[string]string {
	items := []map[string]string{{"title": m.LanguageAuto, "value": ""}}
	for _, tag := range p.I18n().GetSupportLanguages() {
		items = append(items, map[string]string{"title": bcp47.FlagLabel(tag.String()), "value": tag.String()})
	}
	return items
}

// languageLabel is the language lang as shown.
func languageLabel(lang string, m *Messages) string {
	if lang == "" {
		return m.LanguageAuto
	}
	return bcp47.FlagLabel(lang)
}

// configureLanguageField: the default language in the profile, chosen among
// the admin's.
func configureLanguageField(p *presets.Builder, e *presets.EditingBuilder, d *presets.DetailingBuilder) {
	e.Field("Language").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		m := GetMessages(ctx.Context())
		return VSelect().Label(field.Label).Variant(FieldVariantUnderlined).Class("mb-4").
			Hint(m.ProfileLanguageHint).PersistentHint(true).
			Items(languageItems(p, m)).ItemTitle("title").ItemValue("value").
			Attr(web.VField(field.FormKey, field.Value())...).ErrorMessages(field.Errors...).
			Attr("data-profile-language", true)
	})
	d.Field("Language").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		m := GetMessages(ctx.Context())
		lang, _ := field.Value().(string)
		return h.Div(vx.VXReadonlyField().Label(field.InputLabel()).Value(languageLabel(lang, m))).
			Attr("data-profile-language", lang)
	})
}
