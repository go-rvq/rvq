package l10n

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/x/i18n"
)

type localeTestRecord struct {
	ID    uint
	Title string
	Locale
}

func (r *localeTestRecord) PrimarySlug() string {
	return fmt.Sprintf("%d_%s", r.ID, r.LocaleCode)
}

func (r *localeTestRecord) PrimaryColumnValuesBySlug(slug string) map[string]string {
	id, code, _ := strings.Cut(slug, "_")
	return map[string]string{"id": id, "locale_code": code}
}

// ModelInstall registers the action on the detailing builder, so a record can
// be moved to another locale from its detail page.
func TestModelInstallAddsChangeLocaleCodeAction(t *testing.T) {
	pb := presets.New(i18n.New())
	m := pb.Model(&localeTestRecord{})

	b := New(nil).DefaultLocaleCode("en-US")
	b.RegisterLocale("en-US", "en-us", "English")
	b.RegisterLocale("pt-BR", "pt-br", "Português do Brasil")

	if err := b.ModelInstall(pb, m); err != nil {
		t.Fatalf("ModelInstall: %v", err)
	}

	action := m.Detailing().GetAction(ActionChangeLocaleCode)
	if action == nil {
		t.Fatal("action not registered on the detailing builder")
	}
	if got := action.String(); got != ActionChangeLocaleCode {
		t.Errorf("action name = %q, want %q", got, ActionChangeLocaleCode)
	}
}

// Every language of the module carries the four new messages: a partial
// translation renders as blanks, not as English.
func TestChangeLocaleMessagesAreTranslated(t *testing.T) {
	for name, m := range map[string]*Messages{
		"en_US": Messages_en_US,
		"pt_BR": Messages_pt_BR,
		"zh_CN": Messages_zh_CN,
		"ja_JP": Messages_ja_JP,
	} {
		for field, value := range map[string]string{
			"ChangeLocale":               m.ChangeLocale,
			"SuccessfullyChangedLocale":  m.SuccessfullyChangedLocale,
			"ErrChangeLocaleEmpty":       m.ErrChangeLocaleEmpty,
			"ErrChangeLocaleUnavailable": m.ErrChangeLocaleUnavailable,
		} {
			if value == "" {
				t.Errorf("Messages_%s.%s is empty", name, field)
			}
		}
	}
}
