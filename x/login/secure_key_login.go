package login

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-rvq/rvq/x/security"
)

// The login by the secure key: where the installation allows it (the
// `security: {secure_key_login: true}` of its app.yml, package security), the
// secure key — the file only who operates the server reads (secure_key.go)
// — logs the browser in as the administrator, the initial account
// (InitialUserAccount), with no login page: a session as a password login
// makes (completeLogin).
//
// GET of its address (SecureKeyLoginURL, "/auth/secure-key") is a form of the
// key; POST, the key in the form (`key`) or in the X-Secure-Key header. Only
// over HTTPS (or from the machine itself), with the key right; else, refused
// — logged, never said why to the browser beyond "refused".

// SecureKeyLoginURL is the address of the login by the secure key.
func (b *Builder) SecureKeyLoginURL() string { return b.secureKeyLoginURL }

var secureKeyLoginPage = template.Must(template.New("").Parse(`<!doctype html>
<html><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1">
<title>Secure key</title>
<style>body{font-family:system-ui,sans-serif;display:flex;min-height:100vh;align-items:center;justify-content:center;margin:0;background:#f5f5f5}
form{background:#fff;padding:2rem;border-radius:8px;box-shadow:0 1px 4px #0002;width:min(90vw,24rem)}
input{width:100%;box-sizing:border-box;padding:.6rem;margin:.5rem 0 1rem;border:1px solid #bbb;border-radius:4px}
button{padding:.6rem 1.2rem;border:0;border-radius:4px;background:#1e293b;color:#fff;cursor:pointer}
p.error{color:#b91c1c}</style></head>
<body><form method="post" action="{{.Action}}" autocomplete="off">
<label for="key">Secure key</label>
<input id="key" name="key" type="password" required autofocus>
{{if .Refused}}<p class="error">Refused.</p>{{end}}
<button type="submit">Enter</button>
</form></body></html>`))

func (b *Builder) secureKeyLogin(w http.ResponseWriter, r *http.Request) {
	if !security.SecureKeyLogin() {
		http.NotFound(w, r)
		return
	}
	page := func(status int, refused bool) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(status)
		_ = secureKeyLoginPage.Execute(w, map[string]any{"Action": b.secureKeyLoginURL, "Refused": refused})
	}
	switch r.Method {
	case http.MethodGet:
		page(http.StatusOK, false)
		return
	case http.MethodPost:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	sent := r.Header.Get(SecureKeyHeader)
	if sent == "" {
		sent = r.FormValue("key")
	}
	var err error
	user, uerr := b.secureKeyUser()
	switch {
	case !secureTransport(r):
		err = ErrSecureKeyTransport
	case !b.isSecureKey(sent):
		err = ErrSecureKeyWrong
	case uerr != nil:
		err = uerr
	}
	if err != nil {
		log.Printf("secure key login refused: from %s: %v", r.RemoteAddr, err)
		page(http.StatusUnauthorized, true)
		return
	}
	log.Printf("secure key login: from %s, as %s", r.RemoteAddr, b.initialUserAccount)
	if err = b.completeLogin(w, r, user); err != nil {
		http.Redirect(w, r, b.loginPageURL, http.StatusFound)
	}
}
