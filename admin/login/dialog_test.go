package login_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	plogin "github.com/go-rvq/rvq/admin/login"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/login"
	"golang.org/x/text/language"
)

// Uma sessão que morre embaixo de uma página aberta não pode levar a página
// embora: o request da própria página (plaid) volta com 401 e o endereço do
// diálogo de login, e não com um redirecionamento.
func newDialogApp(t *testing.T) http.Handler {
	t.Helper()

	i18nB := i18n.New().SupportLanguages(language.AmericanEnglish, language.BrazilianPortuguese)
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

// plaidRequest é o que plaid() manda: o cabeçalho que diz "quem está pedindo é
// a página, não o navegador".
func plaidRequest(method, url string) *http.Request {
	r := httptest.NewRequest(method, url, nil)
	r.Header.Set(web.PlaidRequestHeader, "1")
	return r
}

func TestSessionLostOnAPlaidRequestAnswers401WithLoginURI(t *testing.T) {
	app := newDialogApp(t)

	w := httptest.NewRecorder()
	app.ServeHTTP(w, plaidRequest("POST", "/admin/products?"+web.EventFuncIDName+"=presets_Update"))

	// 401 é o que o log do servidor precisa registrar
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401 (um redirecionamento levaria a página embora)", w.Code)
	}

	var res struct {
		LoginURI      string `json:"loginURI"`
		UpdatePortals []any  `json:"updatePortals"`
		RunScript     string `json:"runScript"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("resposta não é JSON: %v\n%s", err, w.Body.String())
	}

	if want := "/admin" + presets.LoginDialogURI; res.LoginURI != want {
		t.Errorf("loginURI = %q, want %q", res.LoginURI, want)
	}
	// nada mais: quem monta o diálogo é o request seguinte
	if len(res.UpdatePortals) > 0 || res.RunScript != "" {
		t.Errorf("a resposta do 401 não deve trazer nada além do loginURI:\n%s", w.Body.String())
	}
}

// Sem o cabeçalho não há página aberta para preservar — vale o redirecionamento
// de sempre, mesmo em um request de evento.
func TestEventRequestWithoutThePlaidHeaderStillRedirects(t *testing.T) {
	app := newDialogApp(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/admin/products?"+web.EventFuncIDName+"=presets_Update", nil)
	app.ServeHTTP(w, r)

	if w.Code != http.StatusFound {
		t.Errorf("status = %d, want 302", w.Code)
	}
}

// O endereço que o 401 devolveu monta o diálogo — e responde sem sessão, que é
// exatamente quando ele é pedido.
func TestLoginDialogURIRendersTheDialog(t *testing.T) {
	app := newDialogApp(t)

	w := httptest.NewRecorder()
	app.ServeHTTP(w, plaidRequest("POST", "/admin"+presets.LoginDialogURI+"?"+web.EventFuncIDName+"=__reload__"))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
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
	// dá para ampliar: a tela de login de uma aplicação pode ser alta
	if !strings.Contains(body, "expandable") {
		t.Errorf("o diálogo não é expansível:\n%s", firstLine(body))
	}

	// os diálogos que já estavam abertos saem de vista — e voltam depois
	for _, want := range []string{
		plogin.HiddenByLoginClass,
		"visibility: hidden",
		".v-overlay--active",
		"classList.remove",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("a resposta não esconde/restaura os outros diálogos (%q):\n%s", want, firstLine(body))
		}
	}

	// e, ao final, o request interrompido segue: onLoginSuccess vem no escopo de
	// quem pediu o diálogo (js/corejs/src/builder.ts, login())
	if !strings.Contains(body, "onLoginSuccess()") {
		t.Errorf("a resposta não retoma o request interrompido:\n%s", firstLine(body))
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

// O diálogo diz por que apareceu — e no idioma de quem está usando.
func TestLoginDialogTitleIsTranslated(t *testing.T) {
	app := newDialogApp(t)

	for lang, want := range map[string]string{
		"en":    login.Messages_en_US.LoginAgainTitle,
		"pt-BR": login.Messages_pt_BR.LoginAgainTitle,
	} {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, plaidRequest("POST", "/admin"+presets.LoginDialogURI+"?lang="+lang))

		body := w.Body.String()
		if !strings.Contains(body, want) {
			t.Errorf("lang=%s: o diálogo não traz o título %q:\n%s", lang, want, firstLine(body))
		}
		// e só ele: o idioma pedido é o que sai
		for other, s := range map[string]string{
			"en":    login.Messages_en_US.LoginAgainTitle,
			"pt-BR": login.Messages_pt_BR.LoginAgainTitle,
		} {
			if other != lang && strings.Contains(body, s) {
				t.Errorf("lang=%s: o diálogo traz o título de %s (%q)", lang, other, s)
			}
		}
	}

	if login.Messages_pt_BR.LoginAgainTitle != "Faça login novamente" {
		t.Errorf("título em pt-BR = %q", login.Messages_pt_BR.LoginAgainTitle)
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
