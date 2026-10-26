package mail_sender

import (
	"errors"
	"io"
	"net"
	"strings"

	"github.com/wneessen/go-mail"
)

type Attachment struct {
	Name string
	Body io.ReadCloser
}

type MessageBuilder struct {
	to          []string
	external    bool
	messageID   string
	inReplyTo   string
	references  []string
	replyTo     string
	subject     string
	body        string
	labels      []string
	attachments []*Attachment

	// what the central of the mails sent records (MailLog)
	kind       string
	recordPath string
	sensitive  []string
	hideBody   bool
	verbatim   bool
	userID     string
	remoteIP   string
	userAgent  string
	localeCode string
}

func Message() *MessageBuilder {
	return &MessageBuilder{}
}

func (b *MessageBuilder) To(v ...string) *MessageBuilder {
	b.to = append(b.to, v...)
	return b
}

// External says the message goes to somebody outside the installation — a
// contact, the author of an issue —, not to its team: it counts toward the
// limits of each recipient (MailSender.MaxPerRecipientHour, …Day).
func (b *MessageBuilder) External(v bool) *MessageBuilder {
	b.external = v
	return b
}

// MessageID is the Message-ID of the message (`<id@host>`; the brackets are
// added when missing): what an answer to it names (In-Reply-To), the
// conversation it starts.
func (b *MessageBuilder) MessageID(v string) *MessageBuilder {
	b.messageID = v
	return b
}

// InReplyTo is the Message-ID the message answers, and References those of
// its conversation: what the mail clients group it by.
func (b *MessageBuilder) InReplyTo(v string) *MessageBuilder {
	b.inReplyTo = v
	return b
}

func (b *MessageBuilder) References(v ...string) *MessageBuilder {
	b.references = append(b.references, v...)
	return b
}

func (b *MessageBuilder) GetMessageID() string { return b.messageID }

// Kind is what the message is, for the central of the mails sent (MailLog):
// "contact.received", "user.password_reset"…
func (b *MessageBuilder) Kind(v string) *MessageBuilder {
	b.kind = v
	return b
}

// RecordPath is the record the message is of (presets record path:
// "admin:contact_messages:<…>:"): the link of its MailLog to it.
func (b *MessageBuilder) RecordPath(v string) *MessageBuilder {
	b.recordPath = v
	return b
}

// Sensitive are values the MailLog of the message never keeps: each one is
// "••••••" in its subject and body — a link with a token, a phone.
func (b *MessageBuilder) Sensitive(values ...string) *MessageBuilder {
	b.sensitive = append(b.sensitive, values...)
	return b
}

// HideBody keeps the body of the message out of its MailLog: only that it
// was hidden.
func (b *MessageBuilder) HideBody() *MessageBuilder {
	b.hideBody = true
	return b
}

// Verbatim sends the subject as it is: no labels, no prefix of the
// settings — a message sent again (MailLog), whose subject has them.
func (b *MessageBuilder) Verbatim() *MessageBuilder {
	b.verbatim = true
	return b
}

// User is who made the message be sent (the id of the user), for its
// MailLog.
func (b *MessageBuilder) User(id string) *MessageBuilder {
	b.userID = id
	return b
}

// Origin is where the request that made the message be sent came from: its
// IP and browser, for its MailLog.
func (b *MessageBuilder) Origin(ip, userAgent string) *MessageBuilder {
	if host, _, err := net.SplitHostPort(ip); err == nil {
		ip = host // a RemoteAddr: its port is not where from
	}
	b.remoteIP, b.userAgent = ip, userAgent
	return b
}

// Locale is the language of the message, for its MailLog.
func (b *MessageBuilder) Locale(code string) *MessageBuilder {
	b.localeCode = code
	return b
}

// ReplyTo is the address an answer to the mail goes to: who wrote what it
// carries (a form's visitor), not the sender.
func (b *MessageBuilder) ReplyTo(v string) *MessageBuilder {
	b.replyTo = v
	return b
}

func (b *MessageBuilder) GetReplyTo() string {
	return b.replyTo
}

func (b *MessageBuilder) GetTo() []string {
	return b.to
}

func (b *MessageBuilder) Subject(subject string) *MessageBuilder {
	b.subject = subject
	return b
}

func (b *MessageBuilder) PrependSubject(v string) *MessageBuilder {
	if len(v) > 0 {
		if b.subject == "" {
			b.subject = v
		} else {
			b.subject = v + b.subject
		}
	}
	return b
}

func (b *MessageBuilder) AppendSubject(v string) *MessageBuilder {
	if len(v) > 0 {
		if b.subject == "" {
			b.subject = v
		} else {
			b.subject += v
		}
	}
	return b
}

func (b *MessageBuilder) GetSubject() string {
	return b.subject
}

func (b *MessageBuilder) Body(body string) *MessageBuilder {
	b.body = body
	return b
}

func (b *MessageBuilder) GetBody() string {
	return b.body
}

func (b *MessageBuilder) Label(label ...string) *MessageBuilder {
	b.labels = append(b.labels, label...)
	return b
}

func (b *MessageBuilder) GetLabels() []string {
	return b.labels
}

func (b *MessageBuilder) Attachment(attachment ...*Attachment) *MessageBuilder {
	b.attachments = append(b.attachments, attachment...)
	return b
}

func (b *MessageBuilder) GetAttachments() []*Attachment {
	return b.attachments
}

func (b *MessageBuilder) FormattedSubject() string {
	if b.verbatim {
		return b.subject
	}
	var s []string
	if len(b.labels) > 0 {
		for _, label := range b.labels {
			s = append(s, "["+label+"]")
		}
	}
	s = append(s, b.subject)
	return strings.Join(s, " ")
}

func (b *MessageBuilder) Build() (m *mail.Msg, err error) {
	if len(b.to) == 0 {
		err = errors.New("no destinations")
		return
	}

	m = mail.NewMsg()
	if err = m.To(b.to...); err != nil {
		return
	}
	if b.replyTo != "" {
		if err = m.ReplyTo(b.replyTo); err != nil {
			return
		}
	}
	if b.messageID != "" {
		m.SetMessageIDWithValue(strings.Trim(b.messageID, "<>"))
	} else {
		// one of its own, kept: what its MailLog and an answer to it name
		m.SetMessageID()
		if ids := m.GetGenHeader(mail.HeaderMessageID); len(ids) > 0 {
			b.messageID = strings.Trim(ids[0], "<>")
		}
	}
	if b.inReplyTo != "" {
		m.SetGenHeader(mail.HeaderInReplyTo, msgID(b.inReplyTo))
	}
	if len(b.references) > 0 {
		refs := make([]string, len(b.references))
		for i, r := range b.references {
			refs[i] = msgID(r)
		}
		m.SetGenHeader(mail.HeaderReferences, strings.Join(refs, " "))
	}
	m.Subject(b.FormattedSubject())
	m.SetBodyString(mail.TypeTextHTML, b.body)
	return
}

type Sender interface {
	Send(b *MessageBuilder) (err error)
}

// msgID is v as a message id: `<…>`.
func msgID(v string) string { return "<" + strings.Trim(v, "<>") + ">" }
