package mail_sender

import (
	"errors"

	"github.com/google/uuid"
)

type MailSender struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" admin:"-"`
	// SubjectPrefix is the site's name in the subjects ("PAC"; empty, the
	// application's: SiteName), SubjectRules which parts of the subject go,
	// by family of messages and to whom (subject.go).
	SubjectPrefix string
	SubjectRules  SubjectRules `gorm:"type:jsonb"`
	// IDFooter puts at the end of the body of a message of a record its
	// ids — its token, its Message-ID, its kind (MessageBuilder.Identify) —:
	// what finds it in a mailbox many things write to.
	IDFooter bool `gorm:"not null;default:true"`
	Sender   string
	Gmail    GmailSender `gorm:"embedded;embeddedPrefix:gmail__"`
	SMTP     SmtpSender  `gorm:"embedded;embeddedPrefix:smtp__"`
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
	// MailLogRetentionDays is how long the central of the mails sent keeps
	// a message (MailLog); 0, for ever.
	MailLogRetentionDays uint `gorm:"not null;default:365"`
}

// IsValid reports whether the settings send: a sender set, or the override
// (SMTPOverrideEnv: every message goes there).
func (c *MailSender) IsValid() bool {
	return SMTPOverride() != "" || c.Gmail.IsValid() || c.SMTP.IsValid()
}

// SendByMethod sends b by the sender senderName ("GMAIL", "SMTP"); every
// message, sent or not, is recorded in the central of the mails sent
// (MailLog), its sensitive content hidden.
func (c *MailSender) SendByMethod(senderName string, b *MessageBuilder) (err error) {
	method := senderName
	// the local development's and the tests' server, whatever the settings
	// say (SMTPOverrideEnv): its method even for a message the limits stop
	addr := SMTPOverride()
	if addr != "" {
		method = MethodOverride
	}
	defer func() { record(b, method, err) }()

	// the subject as it would go, in the central too when it does not
	// its parts as the settings say (subject.go)
	c.composeWith(b)
	if b.external {
		if err = limiter.take(b.to, c.MaxPerRecipientHour, c.MaxPerRecipientDay); err != nil {
			return
		}
	}

	var sender Sender

	if addr != "" {
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
