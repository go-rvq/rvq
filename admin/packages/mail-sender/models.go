package mail_sender

import (
	"errors"

	"github.com/google/uuid"
)

type MailSender struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" admin:"-"`
	SubjectPrefix string
	Sender        string
	Gmail         GmailSender `gorm:"embedded;embeddedPrefix:gmail__"`
	SMTP          SmtpSender  `gorm:"embedded;embeddedPrefix:smtp__"`
	// MaxPerRecipientHour and MaxPerRecipientDay are the most messages an
	// address outside the installation (MessageBuilder.External: a contact,
	// the author of an issue) gets in an hour and in a day; 0, no limit. A
	// message over them is refused (ErrRateLimited): a form of the site
	// cannot be used to flood somebody's mailbox.
	MaxPerRecipientHour uint `gorm:"not null;default:10"`
	MaxPerRecipientDay  uint `gorm:"not null;default:30"`
	// Inbound is where the answers to the messages sent are read from
	// (Inbox): the team's mailbox, by people (manual), or by IMAP.
	Inbound InboundSettings `gorm:"embedded;embeddedPrefix:inbound__"`
}

// IsValid reports whether the settings send: a sender set, or the override
// (SMTPOverrideEnv: every message goes there).
func (c *MailSender) IsValid() bool {
	return SMTPOverride() != "" || c.Gmail.IsValid() || c.SMTP.IsValid()
}

func (c *MailSender) SendByMethod(senderName string, b *MessageBuilder) (err error) {
	if b.external {
		if err = limiter.take(b.to, c.MaxPerRecipientHour, c.MaxPerRecipientDay); err != nil {
			return
		}
	}
	b.PrependSubject(c.SubjectPrefix)

	var sender Sender

	// the local development's and the tests' server, whatever the settings
	// say (SMTPOverrideEnv)
	if addr := SMTPOverride(); addr != "" {
		if sender, err = c.overrideSender(addr); err != nil {
			return
		}
		return sender.Send(b)
	}

	switch senderName {
	case "GMAIL":
		if c.Gmail.IsValid() {
			sender = &c.Gmail
		}
	case "SMTP":
		if c.SMTP.IsValid() {
			sender = &c.SMTP
		}
	}

	if sender == nil {
		//lint:ignore ST1005 shown to the user as a sentence
		return errors.New("Mail sender not configured")
	}
	return sender.Send(b)
}

func (c *MailSender) Send(b *MessageBuilder) (err error) {
	return c.SendByMethod(c.Sender, b)
}

func (MailSender) TableName() string {
	return "mail_sender_settings"
}
