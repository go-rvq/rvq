package l10n

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// newBuilder registers the locales in order, so the first registered is not the
// default — which is the case that matters here.
func newBuilder(defaultCode string, codes ...string) *Builder {
	b := New(nil)
	for _, c := range codes {
		b.RegisterLocale(c, c, c)
	}
	return b.DefaultLocaleCode(defaultCode)
}

func request(t *testing.T, query, cookie string) *http.Request {
	t.Helper()
	target := "/"
	if query != "" {
		target += "?locale=" + query
	}
	r := httptest.NewRequest(http.MethodGet, target, nil)
	if cookie != "" {
		r.AddCookie(&http.Cookie{Name: "locale", Value: cookie})
	}
	return r
}

// A request naming no locale is answered by the default, not by whichever
// locale was registered first. Registering "pt-BR" before "en-US" and making
// "en-US" the default is exactly the setup that put pt-BR on new records.
func TestGetCorrectLocaleCodeFallsBackToDefault(t *testing.T) {
	b := newBuilder("en-US", "pt-BR", "en-US")

	if got := b.GetCorrectLocaleCode(request(t, "", "")); got != "en-US" {
		t.Errorf("got %q, want %q", got, "en-US")
	}
}

// A locale that is asked for but not supported — a disabled one, for instance —
// also falls back to the default.
func TestGetCorrectLocaleCodeUnsupportedFallsBackToDefault(t *testing.T) {
	b := newBuilder("en-US", "pt-BR", "en-US")

	if got := b.GetCorrectLocaleCode(request(t, "de-DE", "")); got != "en-US" {
		t.Errorf("query: got %q, want %q", got, "en-US")
	}
	if got := b.GetCorrectLocaleCode(request(t, "", "de-DE")); got != "en-US" {
		t.Errorf("cookie: got %q, want %q", got, "en-US")
	}
}

// What the request does ask for still wins over the default.
func TestGetCorrectLocaleCodeHonoursRequest(t *testing.T) {
	b := newBuilder("en-US", "pt-BR", "en-US")

	if got := b.GetCorrectLocaleCode(request(t, "pt-BR", "")); got != "pt-BR" {
		t.Errorf("query: got %q, want %q", got, "pt-BR")
	}
	if got := b.GetCorrectLocaleCode(request(t, "", "pt-BR")); got != "pt-BR" {
		t.Errorf("cookie: got %q, want %q", got, "pt-BR")
	}
}

// A default that is not among the supported locales — it was disabled, say —
// leaves the first supported one answering, rather than nothing.
func TestGetCorrectLocaleCodeDefaultNotSupported(t *testing.T) {
	b := newBuilder("de-DE", "pt-BR", "en-US")

	if got := b.GetCorrectLocaleCode(request(t, "", "")); got != "pt-BR" {
		t.Errorf("got %q, want %q", got, "pt-BR")
	}
}

// No locales at all must not panic on an empty slice.
func TestGetCorrectLocaleCodeNoLocales(t *testing.T) {
	b := New(nil).DefaultLocaleCode("en-US")

	if got := b.GetCorrectLocaleCode(request(t, "", "")); got != "" {
		t.Errorf("got %q, want an empty string", got)
	}
}
