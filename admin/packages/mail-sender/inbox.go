package mail_sender

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/mail"
	"strconv"
	"strings"
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

// SMTPReuse are the bits of what of the SMTP sender's account (MailSender.SMTP)
// an account takes for its own — the team mailbox that sends is the one that
// receives —: its user, its password, its host. A bit set, that value is the
// SMTP's, and the account's own is not used.
type SMTPReuse uint8

// The bits of SMTPReuse.
const (
	// SMTPReuseUser takes the SMTP's user
	SMTPReuseUser SMTPReuse = 1 << iota
	// SMTPReusePassword takes the SMTP's password
	SMTPReusePassword
	// SMTPReuseHost takes the SMTP's host (Server)
	SMTPReuseHost
)

// SMTPReuseBits are the bits of SMTPReuse, in order.
var SMTPReuseBits = []SMTPReuse{SMTPReuseUser, SMTPReusePassword, SMTPReuseHost}

// Has says every bit of b is set.
func (r SMTPReuse) Has(b SMTPReuse) bool { return r&b == b }

// InboundSettings are where the answers are read from.
type InboundSettings struct {
	// Mode is InboundManual (or "") or InboundIMAP.
	Mode string
	// The IMAP server and account: Server and Port ("imap.gmail.com", 993),
	// TLS — implicit (993) —, or STARTTLS, or none (a local server); the
	// Mailbox read ("" INBOX). Reuse takes values of the SMTP sender's
	// account (SMTPReuse: user, password, host) instead of the ones here.
	Reuse    SMTPReuse `gorm:"not null;default:0"`
	Server   string
	Port     uint16
	TLS      bool
	StartTLS bool
	User     string
	Password string
	Mailbox  string
	// After is what is done with a message taken (Inbox.Done): marked read
	// (InboundAfterRead, or ""), deleted from the server (InboundAfterDelete),
	// moved to the folder MoveTo (InboundAfterMove; empty, the domain of the
	// sending's address: "pac.com" — made when it is not there).
	After  string
	MoveTo string
}

// What is done with a message taken (InboundSettings.After). Deleting and
// moving touch only it: by UID EXPUNGE (UIDPLUS) or MOVE; a server with
// neither only marks it read — a plain EXPUNGE would delete what anything
// else marked deleted in the same mailbox.
const (
	InboundAfterRead   = "read"
	InboundAfterDelete = "delete"
	InboundAfterMove   = "move"
)

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
		s := c.Inbound.WithSMTP(&c.SMTP)
		if s.After == InboundAfterMove && strings.TrimSpace(s.MoveTo) == "" {
			s.MoveTo = c.defaultMoveTo()
		}
		return DialIMAP(s)
	}
	return nil, fmt.Errorf("mail sender: inbound mode %q", c.Inbound.Mode)
}

// defaultMoveTo is the folder the messages taken go to when the settings name
// none: the domain of the sending's address ("pac.com"), else "Imported".
func (c *MailSender) defaultMoveTo() string {
	if _, d, ok := strings.Cut(c.FromAddress(), "@"); ok && strings.TrimSpace(d) != "" {
		return strings.TrimSpace(d)
	}
	return "Imported"
}

// ManualInbox is the answers left to people, in the team's mailbox: nothing
// for the application to read.
type ManualInbox struct{}

func (ManualInbox) Unread(context.Context) ([]*InboundMessage, error) { return nil, nil }
func (ManualInbox) Done(context.Context, *InboundMessage) error       { return nil }
func (ManualInbox) Close() error                                      { return nil }

// IMAPInbox reads an IMAP account's mailbox: the messages not \Seen are the
// unread ones; Done marks one \Seen, deletes it or moves it (after).
type IMAPInbox struct {
	c      *imapclient.Client
	after  string
	moveTo string
	// made: the folder moveTo is there
	made bool
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
	return &IMAPInbox{c: c, after: s.After, moveTo: strings.TrimSpace(s.MoveTo)}, nil
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

// Done is m handled, as the settings say (InboundSettings.After): deleted —
// with UIDPLUS: \Deleted and UID EXPUNGE of it alone —, moved to its folder —
// with MOVE (or UIDPLUS) —, else marked \Seen; a server that cannot delete
// or move one message alone gets it marked \Seen, and the log says why.
func (in *IMAPInbox) Done(_ context.Context, m *InboundMessage) error {
	uid := imap.UIDSetNum(imap.UID(m.UID))
	caps := in.c.Caps()
	switch in.after {
	case InboundAfterDelete:
		if caps.Has(imap.CapUIDPlus) {
			if err := in.c.Store(uid, &imap.StoreFlags{Op: imap.StoreFlagsAdd, Silent: true, Flags: []imap.Flag{imap.FlagDeleted}}, nil).Close(); err != nil {
				return err
			}
			return in.c.UIDExpunge(uid).Close()
		}
		log.Printf("mail sender: the IMAP server has no UIDPLUS: %s marked read, not deleted", m.MessageID)
	case InboundAfterMove:
		if in.moveTo != "" && (caps.Has(imap.CapMove) || caps.Has(imap.CapUIDPlus)) {
			if err := in.ensureMailbox(in.moveTo); err != nil {
				return err
			}
			_, err := in.c.Move(uid, in.moveTo).Wait()
			return err
		}
		log.Printf("mail sender: the IMAP server has neither MOVE nor UIDPLUS: %s marked read, not moved", m.MessageID)
	}
	return in.c.Store(uid, &imap.StoreFlags{
		Op: imap.StoreFlagsAdd, Silent: true, Flags: []imap.Flag{imap.FlagSeen},
	}, nil).Close()
}

// ensureMailbox makes the folder name, once, when it is not there.
func (in *IMAPInbox) ensureMailbox(name string) error {
	if in.made {
		return nil
	}
	list, err := in.c.List("", name, nil).Collect()
	if err != nil {
		return fmt.Errorf("imap: list %s: %w", name, err)
	}
	if len(list) == 0 {
		if err := in.c.Create(name, nil).Wait(); err != nil {
			return fmt.Errorf("imap: create %s: %w", name, err)
		}
	}
	in.made = true
	return nil
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

// WithSMTP is s with the values of smtp it reuses (Reuse): its user, its
// password, its host.
func (s InboundSettings) WithSMTP(smtp *SmtpSender) InboundSettings {
	if s.Reuse.Has(SMTPReuseUser) {
		s.User = smtp.User
	}
	if s.Reuse.Has(SMTPReusePassword) {
		s.Password = smtp.Password
	}
	if s.Reuse.Has(SMTPReuseHost) {
		s.Server = smtp.Server
	}
	return s
}
