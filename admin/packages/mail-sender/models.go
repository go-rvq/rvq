package mail_sender

import "errors"

type MailSender struct {
	ID            uint `gorm:"primarykey" admin:"-"`
	SubjectPrefix string
	Sender        string
	Gmail         GmailSender `gorm:"embedded;embeddedPrefix:gmail__"`
	SMTP          SmtpSender  `gorm:"embedded;embeddedPrefix:smtp__"`
}

func (c *MailSender) IsValid() bool {
	return c.Gmail.IsValid() || c.SMTP.IsValid()
}

func (c *MailSender) SendByMethod(senderName string, b *MessageBuilder) (err error) {
	b.PrependSubject(c.SubjectPrefix)

	var sender Sender

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
