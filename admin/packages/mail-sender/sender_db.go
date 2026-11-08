package mail_sender

import "gorm.io/gorm"

// SenderOf is the sender of the settings of db (MailSender), read at each
// message: what the admin configures, the central of the mails sent
// recording it; the server of SMTPOverrideEnv instead, when there is one.
func SenderOf(db *gorm.DB) Sender { return dbSender{db} }

type dbSender struct{ db *gorm.DB }

func (s dbSender) Send(b *MessageBuilder) error {
	var m MailSender
	s.db.Limit(1).Find(&m)
	return m.Send(b)
}
