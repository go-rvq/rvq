package i18n_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

const runtimeKey i18n.ModuleKey = "runtime"

func requestMessages(t *testing.T, b *i18n.Builder, lang string) *Messages {
	t.Helper()
	var got *Messages
	h := b.EnsureLanguage(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, _ = i18n.MustGetModuleMessages(r.Context(), runtimeKey, nil).(*Messages)
	}))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/?lang="+lang, nil))
	return got
}

func TestSetLanguage(t *testing.T) {
	b := i18n.New().
		SupportLanguages(language.English, language.SimplifiedChinese).
		RegisterForModule(language.English, runtimeKey, Messages_en_US).
		RegisterForModule(language.SimplifiedChinese, runtimeKey, Messages_zh_CN)

	db := &Messages{Update: "Refresh"}
	b.SetLanguage(language.English, map[i18n.ModuleKey]i18n.Messages{runtimeKey: db})

	if got := requestMessages(t, b, "en"); got != db {
		t.Errorf("en: %+v, want the database's", got)
	}
	if got := requestMessages(t, b, "zh"); got != Messages_zh_CN {
		t.Errorf("zh: %+v, want the code's", got)
	}
	if got := b.GetCodeModuleMessages(language.English, runtimeKey); got != Messages_en_US {
		t.Errorf("code: %+v", got)
	}
	if got := b.CodeModules()[language.English][runtimeKey]; got != Messages_en_US {
		t.Errorf("CodeModules: %+v", got)
	}
	// The default's messages (DefaultKey) are kept under the override.
	if b.GetModuleMessages(language.English, i18n.DefaultKey) == nil {
		t.Error("the other modules of the language were lost")
	}

	b.SetLanguage(language.English, nil)
	if got := requestMessages(t, b, "en"); got != Messages_en_US {
		t.Errorf("reset: %+v, want the code's", got)
	}
}

func TestSetSupportLanguagesAtRunTime(t *testing.T) {
	b := i18n.New().
		SupportLanguages(language.English).
		RegisterForModule(language.English, runtimeKey, Messages_en_US)
	zh := &Messages{Update: "更新!"}
	b.SetLanguage(language.SimplifiedChinese, map[i18n.ModuleKey]i18n.Messages{runtimeKey: zh})
	if got := requestMessages(t, b, "zh"); got != Messages_en_US {
		t.Errorf("zh not supported yet: %+v", got)
	}
	b.SetSupportLanguages(language.SimplifiedChinese, language.English)
	if got := requestMessages(t, b, "zh"); got != zh {
		t.Errorf("zh: %+v", got)
	}
	if got := b.GetSupportLanguages(); got[0] != language.SimplifiedChinese {
		t.Errorf("default %v", got)
	}
}

// Updates at run time while requests are served: run with -race.
func TestRuntimeUpdatesRace(t *testing.T) {
	b := i18n.New().
		SupportLanguages(language.English, language.SimplifiedChinese).
		RegisterForModule(language.English, runtimeKey, Messages_en_US)
	h := b.EnsureLanguage(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = i18n.MustGetModuleMessages(r.Context(), runtimeKey, Messages_en_US).(*Messages).Update
	}))
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/?lang=zh", nil))
				_ = b.GetModuleMessages(language.English, runtimeKey)
			}
		}()
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				b.SetLanguage(language.English, map[i18n.ModuleKey]i18n.Messages{runtimeKey: &Messages{Update: "x"}})
				if j%2 == i%2 {
					b.SetSupportLanguages(language.SimplifiedChinese, language.English)
				} else {
					b.SetSupportLanguages(language.English, language.SimplifiedChinese)
				}
			}
		}(i)
	}
	wg.Wait()
}

// A context of a language (LanguageContext): its messages — the nearest
// supported —, and the language asked.
func TestLanguageContext(t *testing.T) {
	b := i18n.New().SupportLanguages(language.English, language.BrazilianPortuguese).
		RegisterForModule(language.English, i18n.DefaultKey, i18n.Default_en).
		RegisterForModule(language.BrazilianPortuguese, i18n.DefaultKey, i18n.Default_pt_BR)
	ctx := b.LanguageContext(context.Background(), "pt-BR")
	if i18n.GetMessages(ctx).Yes != "Sim" || i18n.DynaFromContext(ctx).GetLanguage() != "pt-BR" {
		t.Errorf("pt-BR: %q %q", i18n.GetMessages(ctx).Yes, i18n.DynaFromContext(ctx).GetLanguage())
	}
	if i18n.GetMessages(b.LanguageContext(context.Background(), "en-US")).Yes != "Yes" {
		t.Error("en-US")
	}
	if i18n.GetMessages(b.LanguageContext(context.Background(), "xx")).Yes != "Yes" {
		t.Error("an unknown language: the first supported")
	}
}
