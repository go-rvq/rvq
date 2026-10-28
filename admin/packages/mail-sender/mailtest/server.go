// Package mailtest is a mail server for the tests and the local development:
// an SMTP server that takes every message — any sender, any login — and
// delivers it to the INBOX of each recipient on an IMAP server, both in
// memory. Nothing leaves the machine: an application sending to it (the
// mail sender's override, mail_sender.SMTPOverride) reaches nobody real, and
// what it sent is read back (Messages, the IMAP INBOX of each address) — or
// answered: Deliver puts a message in an address's INBOX as if it came from
// outside (a customer answering a mail).
//
// Every address is an IMAP user, made when it is first delivered to or
// asked for (User), its password Password.
package mailtest

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net"
	"net/mail"
	"strings"
	"sync"
	"time"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapserver"
	"github.com/emersion/go-imap/v2/imapserver/imapmemserver"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

// Password is the IMAP password of every address of the server.
const Password = "test"

// Message is a message the server received by SMTP.
type Message struct {
	From string
	To   []string
	// Raw is the message as it came: headers and body.
	Raw []byte
	// Header is its headers, parsed (nil when they do not parse).
	Header    mail.Header
	Received  time.Time
	Subject   string
	MessageID string
}

// Server is the SMTP server and the IMAP server, listening.
type Server struct {
	// SMTPAddr and IMAPAddr are where they listen ("127.0.0.1:2525").
	SMTPAddr, IMAPAddr string

	smtp  *smtp.Server
	imap  *imapserver.Server
	mem   *imapmemserver.Server
	mu    sync.Mutex
	users map[string]*imapmemserver.User
	msgs  []*Message
	// notify is told of each message received (Wait)
	notify chan struct{}
}

// Options say where the servers listen; empty, a free port of 127.0.0.1.
type Options struct {
	SMTPAddr, IMAPAddr string
	// Logger logs what the servers say; nil, nothing.
	Logger *log.Logger
	// Caps are capabilities the IMAP server advertises besides IMAP4rev1 —
	// imap.CapUIDPlus, imap.CapMove: a server that moves and deletes one
	// message at a time —; none, a plain one.
	Caps []imap.Cap
}

// Start starts the servers.
func Start(opts Options) (*Server, error) {
	s := &Server{
		mem:    imapmemserver.New(),
		users:  map[string]*imapmemserver.User{},
		notify: make(chan struct{}, 1),
	}
	logger := opts.Logger
	if logger == nil {
		logger = log.New(io.Discard, "", 0)
	}

	sl, err := listen(opts.SMTPAddr)
	if err != nil {
		return nil, fmt.Errorf("mailtest: smtp: %w", err)
	}
	il, err := listen(opts.IMAPAddr)
	if err != nil {
		sl.Close()
		return nil, fmt.Errorf("mailtest: imap: %w", err)
	}
	s.SMTPAddr, s.IMAPAddr = sl.Addr().String(), il.Addr().String()

	s.smtp = smtp.NewServer(smtp.BackendFunc(func(*smtp.Conn) (smtp.Session, error) {
		return &smtpSession{s: s}, nil
	}))
	s.smtp.Domain = "mailtest.localhost"
	s.smtp.AllowInsecureAuth = true
	s.smtp.MaxMessageBytes = 32 << 20
	s.smtp.ReadTimeout = time.Minute
	s.smtp.WriteTimeout = time.Minute
	s.smtp.ErrorLog = logger

	s.imap = imapserver.New(&imapserver.Options{
		NewSession: func(*imapserver.Conn) (imapserver.Session, *imapserver.GreetingData, error) {
			return &imapSession{Session: s.mem.NewSession(), s: s}, nil, nil
		},
		// IMAP4rev1: search, fetch and store by UID, what an inbox needs;
		// rev2 would need the session's own methods, which imapSession
		// (it makes the user at login) does not pass on
		Caps:         caps(opts.Caps),
		InsecureAuth: true,
		Logger:       logger,
	})

	go func() { _ = s.smtp.Serve(sl) }()
	go func() { _ = s.imap.Serve(il) }()
	return s, nil
}

func listen(addr string) (net.Listener, error) {
	if addr == "" {
		addr = "127.0.0.1:0"
	}
	return net.Listen("tcp", addr)
}

// Close stops the servers.
func (s *Server) Close() error {
	return errors.Join(s.smtp.Close(), s.imap.Close())
}

// User is the IMAP user of the address addr — its INBOX made —, made when
// it has none yet.
func (s *Server) User(addr string) *imapmemserver.User {
	addr = normalize(addr)
	s.mu.Lock()
	defer s.mu.Unlock()
	if u := s.users[addr]; u != nil {
		return u
	}
	u := imapmemserver.NewUser(addr, Password)
	_ = u.Create("INBOX", nil)
	s.mem.AddUser(u)
	s.users[addr] = u
	return u
}

// Deliver puts the message raw in the INBOX of each address of to, as one
// that came from outside: what an answer of a customer is in a test.
func (s *Server) Deliver(raw []byte, to ...string) error {
	for _, addr := range to {
		if _, err := s.User(addr).Append("INBOX", literal(raw), &imap.AppendOptions{Time: time.Now()}); err != nil {
			return fmt.Errorf("mailtest: deliver to %s: %w", addr, err)
		}
	}
	return nil
}

// Messages are the messages received by SMTP, the oldest first.
func (s *Server) Messages() []*Message {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]*Message(nil), s.msgs...)
}

// To are the messages received for addr.
func (s *Server) To(addr string) (out []*Message) {
	addr = normalize(addr)
	for _, m := range s.Messages() {
		for _, t := range m.To {
			if t == addr {
				out = append(out, m)
				break
			}
		}
	}
	return
}

// Wait waits until n messages were received, at most for d; false when
// they were not.
func (s *Server) Wait(n int, d time.Duration) bool {
	deadline := time.After(d)
	for {
		if len(s.Messages()) >= n {
			return true
		}
		select {
		case <-s.notify:
		case <-deadline:
			return len(s.Messages()) >= n
		}
	}
}

// Reset forgets the messages received (the INBOXes stay).
func (s *Server) Reset() {
	s.mu.Lock()
	s.msgs = nil
	s.mu.Unlock()
}

func (s *Server) received(from string, to []string, raw []byte) error {
	m := &Message{From: from, To: to, Raw: raw, Received: time.Now()}
	if parsed, err := mail.ReadMessage(bytes.NewReader(raw)); err == nil {
		m.Header = parsed.Header
		dec := new(mime.WordDecoder)
		if subj, err := dec.DecodeHeader(parsed.Header.Get("Subject")); err == nil {
			m.Subject = subj
		} else {
			m.Subject = parsed.Header.Get("Subject")
		}
		m.MessageID = parsed.Header.Get("Message-ID")
	}
	s.mu.Lock()
	s.msgs = append(s.msgs, m)
	s.mu.Unlock()
	select {
	case s.notify <- struct{}{}:
	default:
	}
	return s.Deliver(raw, to...)
}

func normalize(addr string) string {
	if a, err := mail.ParseAddress(addr); err == nil {
		addr = a.Address
	}
	return strings.ToLower(strings.TrimSpace(addr))
}

// smtpSession takes a message: any login, any sender, any recipient.
type smtpSession struct {
	s    *Server
	from string
	to   []string
}

var _ smtp.AuthSession = (*smtpSession)(nil)

func (ss *smtpSession) AuthMechanisms() []string { return []string{sasl.Plain} }

func (ss *smtpSession) Auth(string) (sasl.Server, error) {
	return sasl.NewPlainServer(func(string, string, string) error { return nil }), nil
}

func (ss *smtpSession) Mail(from string, _ *smtp.MailOptions) error {
	ss.from = from
	return nil
}

func (ss *smtpSession) Rcpt(to string, _ *smtp.RcptOptions) error {
	ss.to = append(ss.to, normalize(to))
	return nil
}

func (ss *smtpSession) Data(r io.Reader) error {
	raw, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return ss.s.received(ss.from, ss.to, raw)
}

func (ss *smtpSession) Reset() { ss.from, ss.to = "", nil }

func (ss *smtpSession) Logout() error { return nil }

// imapSession logs in any address — made a user then — with Password.
type imapSession struct {
	imapserver.Session
	s *Server
}

// caps are the capabilities of the IMAP server: IMAP4rev1 and extra.
func caps(extra []imap.Cap) imap.CapSet {
	set := imap.CapSet{imap.CapIMAP4rev1: {}}
	for _, c := range extra {
		set[c] = struct{}{}
	}
	return set
}

// Move is the memory session's (MOVE, when the server advertises it).
func (is *imapSession) Move(w *imapserver.MoveWriter, numSet imap.NumSet, dest string) error {
	return is.Session.(imapserver.SessionMove).Move(w, numSet, dest)
}

func (is *imapSession) Login(username, password string) error {
	is.s.User(username)
	return is.Session.Login(normalize(username), password)
}

// literalReader is a message as an IMAP literal (literal).
type literalReader struct {
	*bytes.Reader
	size int64
}

func (l literalReader) Size() int64 { return l.size }

func literal(raw []byte) imap.LiteralReader {
	return literalReader{Reader: bytes.NewReader(raw), size: int64(len(raw))}
}
