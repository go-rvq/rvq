package mail_sender

import (
	"context"

	"github.com/go-rvq/rvq/x/lifecycle"
	"github.com/go-rvq/rvq/x/migrate"
	"gorm.io/gorm"
)

// Key is the node of the package in the graphs of the commands
// (x/lifecycle).
const Key lifecycle.Key = "rvq/mail-sender"

// its migrations, in the graph of every package's (x/migrate): `app migrate`
// and the serve run them; Install still runs them too, for the applications
// that do not run the graph yet
func init() {
	migrate.Graph.Add(lifecycle.Node{
		Key:   Key,
		Start: migrate.Func(func(_ context.Context, db *gorm.DB) error { return Migrate(db) }),
	})
}

// Migrate makes the tables of the settings of the sender and of the mails
// sent; the old switch of the SMTP's login (UseSMTPLogin) made its user and
// password (Reuse).
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&MailSender{}, &MailLog{}); err != nil {
		return err
	}
	if m := db.Migrator(); m.HasColumn(&MailSender{}, "inbound__use_smtp_login") {
		if err := db.Exec(`UPDATE mail_sender_settings SET inbound__reuse = ? WHERE inbound__use_smtp_login`,
			SMTPReuseUser|SMTPReusePassword).Error; err != nil {
			return err
		}
		if err := m.DropColumn(&MailSender{}, "inbound__use_smtp_login"); err != nil {
			return err
		}
	}
	return nil
}
