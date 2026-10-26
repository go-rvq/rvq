package mail_sender_test

import (
	"context"
	"errors"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	mail_sender "github.com/go-rvq/rvq/admin/packages/mail-sender"
	"github.com/go-rvq/rvq/admin/packages/mail-sender/mailtest"
)

func startMailtest(t *testing.T) *mailtest.Server {
	t.Helper()
	s, err := mailtest.Start(mailtest.Options{})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

// With the override (RVQ_MAIL_SMTP), a message goes to that server whatever
// the settings say — none here —, with the headers of its conversation.
func TestSendOverride(t *testing.T) {
	srv := startMailtest(t)
	t.Setenv(mail_sender.SMTPOverrideEnv, srv.SMTPAddr)

	var sender mail_sender.MailSender
	err := sender.Send(mail_sender.Message().
		To("Ana <ana@example.com>").
		Subject("Re: [#01A1] Recebemos sua mensagem").
		Body("<p>Olá, Ana</p>").
		MessageID("hermon.contact.1.2@example.com").
		InReplyTo("hermon.contact.1@example.com").
		References("hermon.contact.1@example.com"))
	if err != nil {
		t.Fatal(err)
	}
	if !srv.Wait(1, 5*time.Second) {
		t.Fatal("nothing arrived")
	}
	got := srv.To("ana@example.com")
	if len(got) != 1 {
		t.Fatalf("to ana: %d messages", len(got))
	}
	m := got[0]
	if m.Subject != "Re: [#01A1] Recebemos sua mensagem" || m.MessageID != "<hermon.contact.1.2@example.com>" ||
		m.Header.Get("In-Reply-To") != "<hermon.contact.1@example.com>" ||
		m.Header.Get("References") != "<hermon.contact.1@example.com>" {
		t.Errorf("message: subject %q, id %q, in-reply-to %q, references %q",
			m.Subject, m.MessageID, m.Header.Get("In-Reply-To"), m.Header.Get("References"))
	}
}

// A message to somebody outside (External) counts toward the limits of the
// recipient: over them, it is refused; one to the team does not count.
func TestSendRecipientLimits(t *testing.T) {
	srv := startMailtest(t)
	t.Setenv(mail_sender.SMTPOverrideEnv, srv.SMTPAddr)
	sender := mail_sender.MailSender{MaxPerRecipientHour: 2, MaxPerRecipientDay: 5}
	send := func(to string, external bool) error {
		return sender.Send(mail_sender.Message().To(to).Subject("x").Body("x").External(external))
	}
	for i := 0; i < 2; i++ {
		if err := send("limited@example.com", true); err != nil {
			t.Fatalf("message %d: %v", i+1, err)
		}
	}
	if err := send("limited@example.com", true); !errors.Is(err, mail_sender.ErrRateLimited) {
		t.Errorf("third in the hour: %v, want ErrRateLimited", err)
	}
	if err := send("other@example.com", true); err != nil {
		t.Errorf("another address: %v", err)
	}
	for i := 0; i < 5; i++ {
		if err := send("team@example.com", false); err != nil {
			t.Fatalf("the team, message %d: %v", i+1, err)
		}
	}
}

// The answers are read by IMAP: an answer delivered to the team's mailbox is
// unread, with the conversation it answers; once Done, it is not read
// again. Manual: nothing to read.
func TestInbox(t *testing.T) {
	srv := startMailtest(t)
	raw := strings.ReplaceAll(`From: Ana <ana@example.com>
To: team@example.com
Subject: Re: [#01A1] Recebemos sua mensagem
Message-ID: <answer-1@example.com>
In-Reply-To: <hermon.contact.1@example.com>
References: <hermon.contact.1@example.com>
Date: Mon, 26 Oct 2026 15:00:00 -0300
Content-Type: text/plain; charset=utf-8

Obrigada, pode ligar amanhã.
`, "\n", "\r\n")
	if err := srv.Deliver([]byte(raw), "team@example.com"); err != nil {
		t.Fatal(err)
	}

	host, port, _ := net.SplitHostPort(srv.IMAPAddr)
	p, _ := strconv.Atoi(port)
	settings := mail_sender.MailSender{Inbound: mail_sender.InboundSettings{
		Mode: mail_sender.InboundIMAP, Server: host, Port: uint16(p),
		User: "team@example.com", Password: mailtest.Password,
	}}
	in, err := settings.Inbox()
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()
	ctx := context.Background()
	got, err := in.Unread(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("unread: %d", len(got))
	}
	m := got[0]
	if m.From == nil || m.From.Address != "ana@example.com" || m.MessageID != "answer-1@example.com" ||
		len(m.InReplyTo) != 1 || m.InReplyTo[0] != "hermon.contact.1@example.com" ||
		!strings.Contains(m.Text, "pode ligar amanhã") {
		t.Errorf("answer: %+v", m)
	}
	if err := in.Done(ctx, m); err != nil {
		t.Fatal(err)
	}
	if got, err := in.Unread(ctx); err != nil || len(got) != 0 {
		t.Errorf("after Done: %d unread, %v", len(got), err)
	}

	manual, _ := (&mail_sender.MailSender{}).Inbox()
	if got, err := manual.Unread(ctx); err != nil || got != nil {
		t.Errorf("manual: %v %v", got, err)
	}
}
