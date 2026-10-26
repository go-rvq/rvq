package mail_sender

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/mail"
	"strconv"
	"time"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	gomessage "github.com/emersion/go-message"
	gomail "github.com/emersion/go-message/mail"
)

// The answers to the messages sent — a customer answering the mail of a
// comment — are read back by an Inbox: by people, in the team's mailbox
// (InboundManual: nothing for the application to read), or by IMAP
// (InboundIMAP: the application reads them, and what it takes — the
// answer to a conversation it knows, by In-Reply-To — becomes, say, a comment
// of the record). The mode is the settings' (MailSender.Inbound).

// The modes of InboundSettings.
const (
	// InboundManual: the answers are the team's, in its mailbox.
	InboundManual = "manual"
	// InboundIMAP: the application reads them by IMAP.
	InboundIMAP = "imap"
)

// InboundSettings are where the answers are read from.
type InboundSettings struct {
	// Mode is InboundManual (or "") or InboundIMAP.
	Mode string
	// The IMAP server and account: Server and Port ("imap.gmail.com", 993),
	// TLS — implicit (993) —, or STARTTLS, or none (a local server); the
	// Mailbox read ("" INBOX).
	Server   string
	Port     uint16
	TLS      bool
	StartTLS bool
	User     string
	Password string
	Mailbox  string
}

// InboundMessage is a message an Inbox read.
type InboundMessage struct {
	// UID is its key in the mailbox (Done marks it).
	UID uint32
	// From, To, Subject, Date: its headers, decoded.
	From    *mail.Address
	To      []*mail.Address
	Subject string
	Date    time.Time
	// MessageID, InReplyTo and References tie it to the conversation it
	// answers: the Message-ID of a message sent (MessageBuilder.MessageID).
	MessageID  string
	InReplyTo  []string
	References []string
	// Text and HTML are its body, the plain and the HTML parts.
	Text, HTML string
	// Raw is the message as it came.
	Raw []byte
}

// Inbox reads the answers that arrived.
type Inbox interface {
	// Unread are the messages that arrived and were not handled (Done),
	// the oldest first.
	Unread(ctx context.Context) ([]*InboundMessage, error)
	// Done marks m handled: not read again.
	Done(ctx context.Context, m *InboundMessage) error
	// Close ends the inbox's connection.
	Close() error
}

// Inbox is the inbox of the settings' mode: none to read (InboundManual,
// ManualInbox), or the IMAP account.
func (c *MailSender) Inbox() (Inbox, error) {
	switch c.Inbound.Mode {
	case "", InboundManual:
		return ManualInbox{}, nil
	case InboundIMAP:
		return DialIMAP(c.Inbound)
	}
	return nil, fmt.Errorf("mail sender: inbound mode %q", c.Inbound.Mode)
}

// ManualInbox is the answers left to people, in the team's mailbox: nothing
// for the application to read.
type ManualInbox struct{}

func (ManualInbox) Unread(context.Context) ([]*InboundMessage, error) { return nil, nil }
func (ManualInbox) Done(context.Context, *InboundMessage) error       { return nil }
func (ManualInbox) Close() error                                      { return nil }

// IMAPInbox reads an IMAP account's mailbox: the messages not \Seen are the
// unread ones; Done marks one \Seen.
type IMAPInbox struct {
	c *imapclient.Client
}

// DialIMAP connects to the account of s and selects its mailbox.
func DialIMAP(s InboundSettings) (*IMAPInbox, error) {
	if s.Server == "" || s.User == "" {
		return nil, errors.New("mail sender: the IMAP server and account are not set")
	}
	port := s.Port
	if port == 0 {
		port = 143
		if s.TLS {
			port = 993
		}
	}
	addr := net.JoinHostPort(s.Server, strconv.Itoa(int(port)))
	var (
		c   *imapclient.Client
		err error
	)
	switch {
	case s.TLS:
		c, err = imapclient.DialTLS(addr, nil)
	case s.StartTLS:
		c, err = imapclient.DialStartTLS(addr, nil)
	default:
		c, err = imapclient.DialInsecure(addr, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("imap %s: %w", addr, err)
	}
	if err = c.Login(s.User, s.Password).Wait(); err != nil {
		c.Close()
		return nil, fmt.Errorf("imap %s: login: %w", addr, err)
	}
	mailbox := s.Mailbox
	if mailbox == "" {
		mailbox = "INBOX"
	}
	if _, err = c.Select(mailbox, nil).Wait(); err != nil {
		c.Close()
		return nil, fmt.Errorf("imap %s: select %s: %w", addr, mailbox, err)
	}
	return &IMAPInbox{c: c}, nil
}

func (in *IMAPInbox) Unread(ctx context.Context) ([]*InboundMessage, error) {
	data, err := in.c.UIDSearch(&imap.SearchCriteria{NotFlag: []imap.Flag{imap.FlagSeen}}, nil).Wait()
	if err != nil {
		return nil, fmt.Errorf("imap: search: %w", err)
	}
	uids := data.AllUIDs()
	if len(uids) == 0 {
		return nil, nil
	}
	section := &imap.FetchItemBodySection{Peek: true}
	bufs, err := in.c.Fetch(imap.UIDSetNum(uids...), &imap.FetchOptions{
		UID:         true,
		BodySection: []*imap.FetchItemBodySection{section},
	}).Collect()
	if err != nil {
		return nil, fmt.Errorf("imap: fetch: %w", err)
	}
	out := make([]*InboundMessage, 0, len(bufs))
	for _, b := range bufs {
		m, err := ParseInbound(b.FindBodySection(section))
		if err != nil {
			continue // not a message: left unread, for people
		}
		m.UID = uint32(b.UID)
		out = append(out, m)
	}
	return out, nil
}

func (in *IMAPInbox) Done(_ context.Context, m *InboundMessage) error {
	return in.c.Store(imap.UIDSetNum(imap.UID(m.UID)), &imap.StoreFlags{
		Op: imap.StoreFlagsAdd, Silent: true, Flags: []imap.Flag{imap.FlagSeen},
	}, nil).Close()
}

func (in *IMAPInbox) Close() error {
	_ = in.c.Logout().Wait()
	return in.c.Close()
}

// ParseInbound reads the message raw: its headers and its text and HTML
// parts.
func ParseInbound(raw []byte) (*InboundMessage, error) {
	r, err := gomail.CreateReader(bytes.NewReader(raw))
	if err != nil && !gomessage.IsUnknownCharset(err) {
		return nil, err
	}
	m := &InboundMessage{Raw: raw}
	h := r.Header
	if list, err := h.AddressList("From"); err == nil && len(list) > 0 {
		m.From = &mail.Address{Name: list[0].Name, Address: list[0].Address}
	}
	if list, err := h.AddressList("To"); err == nil {
		for _, a := range list {
			m.To = append(m.To, &mail.Address{Name: a.Name, Address: a.Address})
		}
	}
	m.Subject, _ = h.Subject()
	m.Date, _ = h.Date()
	m.MessageID, _ = h.MessageID()
	m.InReplyTo, _ = h.MsgIDList("In-Reply-To")
	m.References, _ = h.MsgIDList("References")
	for {
		p, err := r.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			if gomessage.IsUnknownCharset(err) {
				continue
			}
			return m, nil
		}
		ih, ok := p.Header.(*gomail.InlineHeader)
		if !ok {
			continue // an attachment
		}
		ct, _, _ := ih.ContentType()
		body, _ := io.ReadAll(p.Body)
		switch {
		case ct == "text/html" && m.HTML == "":
			m.HTML = string(body)
		case (ct == "text/plain" || ct == "") && m.Text == "":
			m.Text = string(body)
		}
	}
	return m, nil
}
