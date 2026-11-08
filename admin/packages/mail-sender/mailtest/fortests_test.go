package mailtest

import (
	"os"
	"strings"
	"testing"

	mail_sender "github.com/go-rvq/rvq/admin/packages/mail-sender"
)

// The sender straight to the server; the body decoded — a long link
// whole, though quoted-printable breaks its lines; the server of the tests
// once, the application pointed at it.
func TestSenderAndBody(t *testing.T) {
	s, err := Start(Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	link := "https://app.test/confirm?t=" + strings.Repeat("Ab3_-", 20)
	err = s.Sender().Send(mail_sender.Message().To("ana@x.test").Subject("Olá, confirme").
		Body(`<p>Confirme: <a href="` + link + `">aqui</a> — ação</p>`))
	if err != nil {
		t.Fatal(err)
	}
	m := s.Last()
	if m == nil || m.To[0] != "ana@x.test" || m.Subject != "Olá, confirme" {
		t.Fatalf("received: %+v", m)
	}
	if b := m.Body(); !strings.Contains(b, link) || !strings.Contains(b, "ação") {
		t.Errorf("the body: %q", b)
	}
	a := ForTests(t)
	if b := ForTests(t); a != b || os.Getenv(mail_sender.SMTPOverrideEnv) != a.SMTPAddr {
		t.Error("the server of the tests: once, the application at it")
	}
}
