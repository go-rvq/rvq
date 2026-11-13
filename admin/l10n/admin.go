package l10n

import (
	"net/http"
	"strings"

	"github.com/go-rvq/rvq/utils/bcp47"
	"github.com/go-rvq/rvq/x/i18n"
	"gorm.io/gorm"
)

// AdminQueryName is the query of the locale of a model localized by the
// languages of the admin (NewAdmin): "?admin_locale=en".
const AdminQueryName = "admin_locale"

// NewAdmin is a builder of the models localized by the languages of the
// admin (i18nB's) — not the site's —, beside the admin's own (Install): a
// ModelPlugin (ModelBuilder.Use), as the site's is of pages and posts. Its
// locales are the languages of i18nB, each time a request asks for them (a
// language added to the admin is one of them at once); the locale of a
// request is its query AdminQueryName, else the language of the admin, else
// the default — the first language.
func NewAdmin(db *gorm.DB, i18nB *i18n.Builder) *Builder {
	b := New(db).ModelOnly(AdminQueryName)
	sync := func() []string {
		var codes []string
		for _, tag := range i18nB.GetSupportLanguages() {
			code := tag.String()
			codes = append(codes, code)
			if b.GetLocale(code) == nil {
				b.RegisterLocale(code, strings.ToLower(code), bcp47.Name(code))
			}
		}
		return codes
	}
	if codes := sync(); len(codes) > 0 {
		b.DefaultLocaleCode(codes[0])
	}
	b.SupportLocalesFunc(func(*http.Request) []string { return sync() })
	b.FallbackLocaleFunc(func(r *http.Request) string {
		if lang := r.FormValue(i18nB.GetQueryName()); lang != "" {
			return lang
		}
		return i18nB.GetCurrentLangFromCookie(r)
	})
	return b
}
