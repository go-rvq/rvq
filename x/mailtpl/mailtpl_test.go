package mailtpl

import (
	"slices"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/gad-lang/gad"
)

func reset() {
	mu.Lock()
	layers, appFS, renderer, fallback = nil, nil, standalone, "en-US"
	mu.Unlock()
}

const confirmEN = "@main\ntitle Confirm {= Who}\np\n    | Hello {= Who},\n    a[href=\"https://x.test/c\"] confirm\n    | .\n"
const confirmPT = "@main\ntitle Confirme {= Who}\np Olá {= Who}.\n"

// A template in the language, its base, the fallback; the subject its
// <title>, out of the body; the text of the HTML, its links with their
// addresses.
func TestRenderLanguages(t *testing.T) {
	reset()
	Register("shared", fstest.MapFS{
		"templates/mail/shared/confirm.en-US.gadx": {Data: []byte(confirmEN)},
		"templates/mail/shared/confirm.pt.gadx":    {Data: []byte(confirmPT)},
	})
	m, err := Render("shared/confirm", "pt-BR", gad.Dict{"Who": gad.Str("Ana")})
	if err != nil {
		t.Fatal(err)
	}
	if m.Subject != "Confirme Ana" || !strings.Contains(m.HTML, "Olá Ana") || strings.Contains(m.HTML, "<title") || m.Language != "pt" {
		t.Errorf("pt-BR → pt: %+v", m)
	}
	m, err = Render("shared/confirm", "de", gad.Dict{"Who": gad.Str("Ana")})
	if err != nil || m.Language != "en-US" || m.Text != "Hello Ana,confirm (https://x.test/c)." {
		t.Errorf("de → the fallback: %+v %v", m, err)
	}
	if _, err := Render("shared/none", "en", nil); err == nil {
		t.Error("no template")
	}
}

// The application's files over the defaults; the defaults it lacks.
func TestAppFSAndMissing(t *testing.T) {
	reset()
	Register("shared", fstest.MapFS{
		"templates/mail/shared/confirm.en-US.gadx": {Data: []byte(confirmEN)},
		"templates/mail/shared/invite.en-US.gadx":  {Data: []byte(confirmEN)},
	})
	app := fstest.MapFS{"templates/mail/shared/confirm.en-US.gadx": {Data: []byte("@main\ntitle Mine\np mine\n")}}
	SetFS(app)
	if m, _ := Render("shared/confirm", "en-US", gad.Dict{"Who": gad.Str("x")}); m == nil || m.Subject != "Mine" {
		t.Errorf("the application's own: %+v", m)
	}
	if m, _ := Render("shared/invite", "en-US", gad.Dict{"Who": gad.Str("x")}); m == nil || m.Subject != "Confirm x" {
		t.Errorf("the default for what it lacks: %+v", m)
	}
	miss, err := Missing(app)
	if err != nil || !slices.Equal(miss, []string{"templates/mail/shared/invite.en-US.gadx"}) {
		t.Errorf("missing: %v %v", miss, err)
	}
}

// A renderer of the application.
func TestRenderer(t *testing.T) {
	reset()
	Register("x", fstest.MapFS{"templates/mail/x/a.en-US.gadx": {Data: []byte("ignored")}})
	SetRenderer(RendererFunc(func(path string, _ []byte, _ gad.Dict) (string, error) {
		return "<title>From " + path + "</title><p>body</p>", nil
	}))
	if m, _ := Render("x/a", "en-US", nil); m == nil || m.Subject != "From templates/mail/x/a.en-US.gadx" || m.Text != "body" {
		t.Errorf("the application's renderer: %+v", m)
	}
}
