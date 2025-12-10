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
	m.Subject(b.FormattedSubject())
	m.SetBodyString(mail.TypeTextHTML, b.body)
	return
}

type Sender interface {
	Send(b *MessageBuilder) (err error)
}
