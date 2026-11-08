package i18n

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/text/language"
)

// A request that asks no language (no query, no cookie): the default
// (the user's), before the browser's; one that asks, its own.
func TestDefaultLanguage(t *testing.T) {
	type msgs struct{ Hi string }
	const key ModuleKey = "test"
	b := New().SupportLanguages(language.English, language.BrazilianPortuguese).
		RegisterForModule(language.English, key, &msgs{"hi"}).
		RegisterForModule(language.BrazilianPortuguese, key, &msgs{"oi"})
	b.DefaultLanguage(func(*http.Request) string { return "pt-BR" })
	var got string
	h := b.EnsureLanguage(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		got = DynaFromContext(r.Context()).GetLanguage()
	}))
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Language", "en")
	h.ServeHTTP(httptest.NewRecorder(), r)
	if got != "pt-BR" {
		t.Errorf("the default: %q", got)
	}
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/?"+b.GetQueryName()+"=en", nil))
	if got != "en" {
		t.Errorf("asked: %q", got)
	}
}
