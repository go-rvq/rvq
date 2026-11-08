package mailtest

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"

	mail_sender "github.com/go-rvq/rvq/admin/packages/mail-sender"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// forTests is the server of the tests of the process, started once.
var forTests struct {
	once sync.Once
	s    *Server
	err  error
}

// ForTests is the server of the tests of the process — started once — and
// every message of the application sent to it (mail_sender.SMTPOverrideEnv),
// whatever its settings say: nothing leaves the machine; a test reads what
// was sent (To, Wait, Body), or answers it (Deliver, read by IMAP).
func ForTests(t testing.TB) *Server {
	t.Helper()
	forTests.once.Do(func() {
		forTests.s, forTests.err = Start(Options{})
	})
	if forTests.err != nil {
		t.Fatalf("the mail server of the tests: %v", forTests.err)
	}
	if err := os.Setenv(mail_sender.SMTPOverrideEnv, forTests.s.SMTPAddr); err != nil {
		t.Fatal(err)
	}
	return forTests.s
}

// EnsureSettings makes the settings of the mail sender (mail_sender.MailSender)
// when db has none: what an application reads before it sends.
func EnsureSettings(db *gorm.DB) error {
	var n int64
	if err := db.Model(&mail_sender.MailSender{}).Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	return db.Create(&mail_sender.MailSender{ID: uuid.New(), Sender: "SMTP"}).Error
}

// Sender is a sender straight to the server (no settings, no central of the
// mails sent): what a test of a package that sends uses.
func (s *Server) Sender() mail_sender.Sender {
	host, port, _ := net.SplitHostPort(s.SMTPAddr)
	p, _ := strconv.ParseUint(port, 10, 16)
	return &mail_sender.SmtpSender{Server: host, Port: uint16(p), User: "test", Password: Password,
		FromMail: "noreply@mailtest.localhost"}
}

// Last is the last message received; nil when none.
func (s *Server) Last() *Message {
	msgs := s.Messages()
	if len(msgs) == 0 {
		return nil
	}
	return msgs[len(msgs)-1]
}

// Body is the text of the message — its HTML or text parts, joined —,
// decoded (quoted-printable, base64).
func (m *Message) Body() string {
	msg, err := mail.ReadMessage(bytes.NewReader(m.Raw))
	if err != nil {
		return ""
	}
	var b strings.Builder
	part(&b, msg.Header.Get("Content-Type"), msg.Header.Get("Content-Transfer-Encoding"), msg.Body)
	return b.String()
}

// part writes the text of a part (its parts', of a multipart).
func part(b *strings.Builder, ctype, enc string, r io.Reader) {
	mt, params, err := mime.ParseMediaType(ctype)
	if err != nil {
		mt = "text/plain"
	}
	if strings.HasPrefix(mt, "multipart/") {
		mr := multipart.NewReader(r, params["boundary"])
		for {
			p, err := mr.NextRawPart()
			if err != nil {
				return
			}
			part(b, p.Header.Get("Content-Type"), p.Header.Get("Content-Transfer-Encoding"), p)
		}
	}
	if !strings.HasPrefix(mt, "text/") {
		return
	}
	switch strings.ToLower(enc) {
	case "quoted-printable":
		r = quotedprintable.NewReader(r)
	case "base64":
		r = base64.NewDecoder(base64.StdEncoding, r)
	}
	data, _ := io.ReadAll(r)
	b.Write(data)
	b.WriteString("\n")
}
