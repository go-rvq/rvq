package mail_sender

import (
	"context"
	"html"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var rgxCode = regexp.MustCompile(`<code>([^<]*)</code>`)

func TestGmailSetupComponent(t *testing.T) {
	const callback = "http://localhost:9000/admin/admin/mail-sender/" + oauthCallbackUri

	for name, msgr := range map[string]*Messages{
		"en_US": Messages_en_US,
		"pt_BR": Messages_pt_BR,
	} {
		t.Run(name, func(t *testing.T) {
			if len(msgr.GmailSenderSetupSteps) == 0 {
				t.Fatal("no setup steps translated")
			}

			// Every GmailSenderSetup* message must be translated: an empty one
			// renders as a blank element, and a Contains assertion on it would
			// pass vacuously.
			v := reflect.ValueOf(*msgr)
			for i := 0; i < v.NumField(); i++ {
				name := v.Type().Field(i).Name
				if !strings.HasPrefix(name, "GmailSenderSetup") {
					continue
				}
				switch f := v.Field(i); f.Kind() {
				case reflect.String:
					if f.String() == "" {
						t.Errorf("%s is empty", name)
					}
				case reflect.Slice:
					if f.Len() == 0 {
						t.Errorf("%s is empty", name)
					}
				}
			}

			out, err := h.MarshallString(gmailSetupComponent(msgr, callback), context.Background())
			if err != nil {
				t.Fatal(err)
			}

			want := []string{
				`<vx-dialog`,
				`v-model='locals.setupOpen'`,
				`closable width='680'`,
				`title='` + msgr.GmailSenderSetup + `'`,
				`@click='locals.setupOpen = true'`,
				`mdi-help-circle-outline`,
				callback,
				gmailConsoleCredentialsURL,
				html.EscapeString(msgr.GmailSenderSetupBlocked),
			}
			for _, w := range want {
				if !strings.Contains(out, w) {
					t.Errorf("missing %q in:\n%s", w, out)
				}
			}
			for _, scope := range GmailScopes {
				if !strings.Contains(out, "<code>"+scope+"</code>") {
					t.Errorf("missing scope %q", scope)
				}
			}
			// Only the <code> list is an offer to declare; the prose names the
			// restricted scopes precisely to tell the reader not to add them.
			for _, code := range rgxCode.FindAllStringSubmatch(out, -1) {
				for _, restricted := range []string{"mail.google.com", "gmail.modify", "gmail.compose"} {
					if strings.Contains(code[1], restricted) {
						t.Errorf("dialog lists restricted scope %q, which the sender does not use", code[1])
					}
				}
			}
			for _, step := range msgr.GmailSenderSetupScopesSteps {
				if !strings.Contains(out, html.EscapeString(step)) {
					t.Errorf("missing scope step %q", step)
				}
			}
			for _, step := range msgr.GmailSenderSetupSteps {
				if !strings.Contains(out, html.EscapeString(step)) {
					t.Errorf("missing step %q", step)
				}
			}
		})
	}
}

func TestGmailCallbackURIPrecedence(t *testing.T) {
	const computed = "http://localhost:9000/admin/admin/mail-sender/" + oauthCallbackUri

	for _, c := range []struct {
		name, pinned, want string
	}{
		{"empty falls back", "", computed},
		{"pinned wins", "https://site.tld/x/callback", "https://site.tld/x/callback"},
	} {
		t.Run(c.name, func(t *testing.T) {
			s := &GmailSender{CallbackURI: c.pinned}
			if got := s.CallbackURIOr(computed); got != c.want {
				t.Errorf("got %q, want %q", got, c.want)
			}
		})
	}
}

func TestGmailTokenHasScope(t *testing.T) {
	const send = "https://www.googleapis.com/auth/gmail.send"

	for _, c := range []struct {
		name, scope string
		want        bool
	}{
		{"granted", "openid email " + send, true},
		{"granular consent unticked", "openid email profile", false},
		{"prefix is not a match", "openid " + send + ".readonly", false},
		{"legacy token without scope recorded", "", true},
	} {
		t.Run(c.name, func(t *testing.T) {
			if got := (&GmailToken{Scope: c.scope}).HasScope(send); got != c.want {
				t.Errorf("HasScope(%q) = %v, want %v", c.scope, got, c.want)
			}
		})
	}
}

func TestGmailScopesAreNotRestricted(t *testing.T) {
	for _, scope := range GmailScopes {
		for _, restricted := range []string{"mail.google.com", "gmail.modify", "gmail.compose"} {
			if strings.Contains(scope, restricted) {
				t.Errorf("GmailScopes asks for the restricted scope %q", scope)
			}
		}
	}
}

func TestGmailAuthCodeURL(t *testing.T) {
	u, err := url.Parse(gmailAuthCodeURL(&oauth2.Config{
		ClientID:    "cid",
		RedirectURL: "http://localhost:9000/cb",
		Scopes:      GmailScopes,
		Endpoint:    google.Endpoint,
	}))
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()

	// Without prompt=consent Google may skip the screen and reissue from the
	// remembered grant, dropping a newly added scope and the refresh token.
	if got := q.Get("prompt"); got != "consent" {
		t.Errorf("prompt = %q, want %q", got, "consent")
	}
	if got := q.Get("access_type"); got != "offline" {
		t.Errorf("access_type = %q, want %q", got, "offline")
	}
	if got := q.Get("scope"); got != strings.Join(GmailScopes, " ") {
		t.Errorf("scope = %q, want %q", got, strings.Join(GmailScopes, " "))
	}
}
