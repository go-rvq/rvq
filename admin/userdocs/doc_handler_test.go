package userdocs

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// The page of a guide made open: a page of its own, or the content of a
// rvq-portal (?portal=1), in the language asked for — the first, saying so,
// when it is not in it —, its languages and the other guides open links to
// it; for a site its code HTML's and nothing read by Vue (v-pre), for the
// admin its code the admin's code viewer; a node not open, none.
func TestDocHandler(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	if err := Migrate(db); err != nil {
		t.Fatal(err)
	}
	src := &Source{Package: "example.com/pkg", FS: fstest.MapFS{
		"en/guides/editor/README.md":    {Data: []byte("# Editor\n\nSee [the other](../other/README.md).\n\n```go\nx := {{ 1 }}\n```\n")},
		"pt-BR/guides/editor/README.md": {Data: []byte("# Editor (pt)\n\nTexto.\n")},
		"en/guides/other/README.md":     {Data: []byte("# Other\n")},
		"en/guides/closed/README.md":    {Data: []byte("# Closed\n")},
	}}
	forgetTitles()
	ib := i18n.New().SupportLanguages(language.English, language.MustParse("pt-BR"))
	ConfigureMessages(ib)
	b := &Builder{p: presets.New(ib).URIPrefix("/admin"), db: db, sources: []*Source{src},
		locale:  func(*web.EventContext) string { return "en" },
		locales: func() []string { return []string{"en", "pt-BR"} }}
	// the documentation's page: where its links go
	b.mb = b.p.Model(&UserDoc{}).URIName("docs").Singleton(true)
	b.Public("guides/editor", "guides/other")
	h := b.DocHandler()
	get := func(q string) (int, string) {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, b.DocPath()+"?"+q, nil))
		return w.Code, w.Body.String()
	}
	has := func(name, body string, want ...string) {
		t.Helper()
		for _, x := range want {
			if !strings.Contains(body, x) {
				t.Errorf("%s: no %s\n%s", name, x, body)
			}
		}
	}
	hasNot := func(name, body string, not ...string) {
		t.Helper()
		for _, x := range not {
			if strings.Contains(body, x) {
				t.Errorf("%s: %s\n%s", name, x, body)
			}
		}
	}

	// a page of its own
	code, page := get("node=guides/editor&locale=en-US")
	if code != 200 {
		t.Fatalf("page: %d", code)
	}
	has("page", page, `<!DOCTYPE html><html lang="en">`, `<title>Editor</title>`, `<div class="user-doc" lang="en" v-pre>`,
		`<pre data-fence><code class="language-go">x := {{ 1 }}`,
		// the other guide open: its page; the languages, links to this one
		`href="/_rvq/user-docs/doc?locale=en&amp;node=guides%2Fother"`,
		`href="/_rvq/user-docs/doc?locale=pt-BR&amp;node=guides%2Feditor" data-locale="pt-BR">Português</a>`,
		`data-locale="en" aria-current="true">English</a>`)
	has("page", page, "<style>:root{--v-border-color")
	hasNot("page", page, "<vx-code", "data-user-doc-not-yours", "<component")

	// the content of a portal: no page around it; its links, the portal's
	code, frag := get("node=guides/editor&locale=pt-br&portal=1")
	if code != 200 {
		t.Fatalf("portal: %d", code)
	}
	has("portal", frag, `<div class="user-docs" data-user-doc="guides/editor" data-user-doc-locale="pt-BR">`,
		`<component :is="'style'">`, "Editor (pt)", `>Idioma:</span>`)
	hasNot("portal", frag, "<html", "data-user-doc-not-yours", "portal=1")

	// a language it is not in: the first, saying so (in its words)
	_, frag = get("node=guides/editor&locale=de&portal=1")
	has("de", frag, `data-user-doc-locale="en"`, `data-user-doc-not-yours>This guide is not in your language`)

	// the admin's: its code the code viewer, read by Vue
	_, frag = get("node=guides/editor&locale=en&view=admin&portal=1")
	has("admin", frag, "<vx-code", `view=admin`)
	hasNot("admin", frag, "v-pre")

	if code, _ := get("node=guides/closed&locale=en"); code != http.StatusNotFound {
		t.Errorf("a node not open: %d", code)
	}
	if got := DocDataHref("/admin/user-docs/_doc", "a&b"); got != "/admin/user-docs/_doc?node=a%26b" {
		t.Errorf("DocDataHref %q", got)
	}
	if got := b.DocDataHref("guides/editor"); got != "/_rvq/user-docs/doc?node=guides%2Feditor&view=admin" {
		t.Errorf("Builder.DocDataHref %q", got)
	}
}
