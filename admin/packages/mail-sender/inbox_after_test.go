package mail_sender_test

import (
	"context"
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	mail_sender "github.com/go-rvq/rvq/admin/packages/mail-sender"
	"github.com/go-rvq/rvq/admin/packages/mail-sender/mailtest"
)

// What is done with a message taken (InboundSettings.After), of two in the
// mailbox, the other untouched: marked read; deleted — UID EXPUNGE, with
// UIDPLUS —; moved — with MOVE — to its folder, by default the domain of the
// sending's address, made. A server with neither only marks it read.
func TestInboxAfter(t *testing.T) {
	for name, c := range map[string]struct {
		after                 string
		caps                  []imap.Cap
		inbox, unseen, folder uint32
	}{
		"read":                    {mail_sender.InboundAfterRead, nil, 2, 1, 0},
		"deleted, UIDPLUS":        {mail_sender.InboundAfterDelete, []imap.Cap{imap.CapUIDPlus}, 1, 1, 0},
		"deleted, a plain server": {mail_sender.InboundAfterDelete, nil, 2, 1, 0},
		"moved, MOVE":             {mail_sender.InboundAfterMove, []imap.Cap{imap.CapMove}, 1, 1, 1},
		"moved, a plain server":   {mail_sender.InboundAfterMove, nil, 2, 1, 0},
		"moved, UIDPLUS, no MOVE": {mail_sender.InboundAfterMove, []imap.Cap{imap.CapUIDPlus}, 1, 1, 1},
	} {
		t.Run(name, func(t *testing.T) {
			srv, err := mailtest.Start(mailtest.Options{Caps: c.caps})
			if err != nil {
				t.Fatal(err)
			}
			defer srv.Close()
			const team = "team@example.com"
			for _, id := range []string{"a-1", "a-2"} {
				raw := strings.ReplaceAll("From: Ana <ana@example.com>\nTo: "+team+"\nSubject: x\nMessage-ID: <"+id+"@example.com>\n"+
					"Date: Mon, 26 Oct 2026 15:00:00 -0300\nContent-Type: text/plain; charset=utf-8\n\nx\n", "\n", "\r\n")
				if err := srv.Deliver([]byte(raw), team); err != nil {
					t.Fatal(err)
				}
			}
			host, port, _ := net.SplitHostPort(srv.IMAPAddr)
			p, _ := strconv.Atoi(port)
			settings := mail_sender.MailSender{
				SMTP: mail_sender.SmtpSender{FromMail: "site@pac.com"},
				Inbound: mail_sender.InboundSettings{Mode: mail_sender.InboundIMAP, Server: host, Port: uint16(p),
					User: team, Password: mailtest.Password, After: c.after},
			}
			in, err := settings.Inbox()
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.Background()
			got, err := in.Unread(ctx)
			if err != nil || len(got) != 2 {
				t.Fatalf("unread: %d %v", len(got), err)
			}
			if err := in.Done(ctx, got[0]); err != nil {
				t.Fatalf("done: %v", err)
			}
			in.Close()

			cl, err := imapclient.DialInsecure(srv.IMAPAddr, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer cl.Close()
			if err := cl.Login(team, mailtest.Password).Wait(); err != nil {
				t.Fatal(err)
			}
			st, err := cl.Status("INBOX", &imap.StatusOptions{NumMessages: true, NumUnseen: true}).Wait()
			if err != nil {
				t.Fatal(err)
			}
			if *st.NumMessages != c.inbox || *st.NumUnseen != c.unseen {
				t.Errorf("the mailbox: %d messages, %d unseen", *st.NumMessages, *st.NumUnseen)
			}
			var moved uint32
			if f, err := cl.Status("pac.com", &imap.StatusOptions{NumMessages: true}).Wait(); err == nil {
				moved = *f.NumMessages
			}
			if moved != c.folder {
				t.Errorf("the folder: %d messages", moved)
			}
		})
	}
}
