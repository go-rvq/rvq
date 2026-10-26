package mail_sender

import (
	"errors"
	"io"
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
