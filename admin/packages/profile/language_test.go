package profile

import (
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

// The default language is chosen among the admin's, none (the browser's)
// first; another is not kept.
func TestLanguageChoices(t *testing.T) {
	p := presets.New(i18n.New().SupportLanguages(language.English, language.BrazilianPortuguese))
	items := languageItems(p, Messages_pt_BR)
	if len(items) != 3 || items[0]["value"] != "" || items[0]["title"] != Messages_pt_BR.LanguageAuto ||
		items[2]["value"] != "pt-BR" {
		t.Fatalf("the items: %v", items)
	}
	if !supported(p, "pt-BR") || supported(p, "fr") {
		t.Error("supported")
	}
	if languageLabel("", Messages_en_US) != Messages_en_US.LanguageAuto {
		t.Error("the label of none")
	}
}
