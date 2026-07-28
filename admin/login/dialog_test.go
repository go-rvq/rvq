package login_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	plogin "github.com/go-rvq/rvq/admin/login"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/login"
)

type dialogUser struct {
	ID    uint
	Name  string
	Email string
}

func (u *dialogUser) GetAccountName() string { return u.Email }

// Uma sessão que morre embaixo de uma página aberta não pode levar a página
// embora: o request da própria página (um evento) volta com o login em diálogo,
// e não com um redirecionamento.
func newDialogApp(t *testing.T) http.Handler {
	t.Helper()

	i18nB := i18n.New()
	pb := presets.New(i18nB).URIPrefix("/admin")

	lb := login.New(i18nB).
		Secret("a secret worth at least this many bytes").
		HomeURLFunc(func(r *http.Request, user any) string { return "/admin" })

	if err := plogin.New(lb).Install(pb); err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	lb.Mount(mux)
	mux.Handle("/", pb)
	return lb.Middleware()(mux)
}

func TestSessionLostOnAnEventOpensTheLoginDialog(t *testing.T) {
	app := newDialogApp(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/admin/products?"+web.EventFuncIDName+"=presets_Update", nil)
	app.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (um redirecionamento levaria a página embora)", w.Code)
	}
	body := w.Body.String()

	if !strings.Contains(body, presets.LoginPortalName) {
		t.Errorf("a resposta não atualiza o portal do login:\n%s", firstLine(body))
	}
	if !strings.Contains(body, "vars.presetsLoginDialog = true") {
		t.Errorf("a resposta não abre o diálogo:\n%s", firstLine(body))
	}
	if !strings.Contains(body, "iframe") {
		t.Errorf("o diálogo não carrega a página de login:\n%s", firstLine(body))
	}

	// e o login sabe para onde voltar: a página que avisa "pronto"
	var cont string
	for _, c := range w.Result().Cookies() {
		if strings.Contains(c.Name, "continue") {
			cont = c.Value
		}
	}
	if !strings.HasSuffix(cont, presets.LoginDoneURI) {
		t.Errorf("continue URL = %q, queria terminar em %q", cont, presets.LoginDoneURI)
	}
}

// Uma navegação comum continua indo para a tela de login: não há página aberta
// para preservar.
func TestPlainNavigationStillRedirects(t *testing.T) {
	app := newDialogApp(t)

	w := httptest.NewRecorder()
	app.ServeHTTP(w, httptest.NewRequest("GET", "/admin/products", nil))

	if w.Code != http.StatusFound {
		t.Errorf("status = %d, want 302", w.Code)
	}
	if loc := w.Header().Get("Location"); !strings.Contains(loc, "login") {
		t.Errorf("Location = %q, queria a tela de login", loc)
	}
}

// A página que o frame carrega ao final avisa quem está em volta. (Direto no
// handler: passando pelo middleware sem sessão o teste só veria o redirect, e
// na vida real o frame só chega aqui DEPOIS de logar.)
func TestLoginDonePageReportsBack(t *testing.T) {
	w := httptest.NewRecorder()
	plogin.LoginDonePage().ServeHTTP(w, httptest.NewRequest("GET", "/admin"+presets.LoginDoneURI, nil))

	body := w.Body.String()
	if !strings.Contains(body, "rvq:login-done") || !strings.Contains(body, "postMessage") {
		t.Errorf("a página de retorno não avisa a página em volta:\n%s", body)
	}
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i > 0 {
		return s[:i]
	}
	if len(s) > 300 {
		return s[:300]
	}
	return s
}
