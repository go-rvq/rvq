package mail_sender_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	mail_sender "github.com/go-rvq/rvq/admin/packages/mail-sender"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func mailLogDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:maillog_%d?mode=memory&cache=shared", time.Now().UnixNano())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&mail_sender.MailSender{}, &mail_sender.MailLog{}); err != nil {
		t.Fatal(err)
	}
	mail_sender.SetRecorder(db)
	t.Cleanup(func() { mail_sender.SetRecorder(nil) })
	return db
}

func lastLog(t *testing.T, db *gorm.DB) *mail_sender.MailLog {
	t.Helper()
	var l mail_sender.MailLog
	if err := db.Order("created_at DESC").First(&l).Error; err != nil {
		t.Fatal(err)
	}
	return &l
}

// Every message sent, or that failed, is a MailLog: whom to, what about, by
// which method, the conversation (its Message-ID, the one the server got),
// the record it is of, who and where from; the sensitive values and the
// secrets of its links hidden, then it cannot be sent again. A failure is
// recorded with why.
func TestMailLog(t *testing.T) {
	srv := startMailtest(t)
	t.Setenv(mail_sender.SMTPOverrideEnv, srv.SMTPAddr)
	db := mailLogDB(t)
	sender := mail_sender.MailSender{SubjectPrefix: "[Site] ", MaxPerRecipientHour: 1}

	err := sender.Send(mail_sender.Message().
		To("ana@example.com").
		ReplyTo("team@example.com").
		Subject("Reset 123-456").
		Body(`<p>Ana, 123-456: <a href="https://x.example/reset?user=1&token=s3cr3t">reset</a></p>`).
		Kind("user.password_reset").
		RecordPath("admin:users:<1>:").
		Sensitive("123-456").
		User("7").
		Origin("10.0.0.1", "Firefox").
		Locale("pt-BR").
		External(true))
	if err != nil {
		t.Fatal(err)
	}
	if !srv.Wait(1, 5*time.Second) {
		t.Fatal("nothing arrived")
	}
	got := srv.To("ana@example.com")[0]
	l := lastLog(t, db)
	if l.Status != mail_sender.MailLogSent || l.Kind != "user.password_reset" || l.RecordPath != "admin:users:<1>:" ||
		l.Recipients != "ana@example.com" || l.ReplyTo != "team@example.com" || l.Method != mail_sender.MethodOverride ||
		l.UserID != "7" || l.RemoteIP != "10.0.0.1" || l.UserAgent != "Firefox" || l.LocaleCode != "pt-BR" {
		t.Errorf("the log: %+v", l)
	}
	if l.MessageID == "" || "<"+l.MessageID+">" != got.MessageID {
		t.Errorf("the Message-ID: logged %q, sent %q", l.MessageID, got.MessageID)
	}
	if l.Subject != "[Site] Reset "+mail_sender.Hidden || strings.Contains(l.Body, "123-456") ||
		strings.Contains(l.Body, "s3cr3t") || !strings.Contains(l.Body, "token="+mail_sender.Hidden) ||
		!strings.Contains(l.Body, "user=1") || !l.Redacted {
		t.Errorf("redacted: subject %q, body %q, redacted %v", l.Subject, l.Body, l.Redacted)
	}
	if err := sender.Resend(l); !errors.Is(err, mail_sender.ErrRedacted) {
		t.Errorf("sending again a redacted one: %v", err)
	}

	// over the limit: failed, and why —
	if err := sender.Send(mail_sender.Message().To("ana@example.com").Subject("x").Body("x").External(true)); err == nil {
		t.Fatal("over the limit, sent")
	}
	// by the method it would have gone, with the subject it would have had
	if l := lastLog(t, db); l.Status != mail_sender.MailLogFailed || !strings.Contains(l.Error, "too many messages") ||
		l.Method != mail_sender.MethodOverride || l.Subject != "[Site] x" {
		t.Errorf("the failure: %q %q, method %q, subject %q", l.Status, l.Error, l.Method, l.Subject)
	}

	// the body hidden
	if err := sender.Send(mail_sender.Message().To("bob@example.com").Subject("Secret").Body("<p>all of it</p>").HideBody()); err != nil {
		t.Fatal(err)
	}
	if l := lastLog(t, db); l.Body != "" || !l.BodyHidden || !l.Redacted {
		t.Errorf("hidden: %+v", l)
	}
}

// A message with nothing hidden is sent again as it was — its subject not
// prefixed twice — in its conversation; the old ones go with the retention.
func TestMailLogResendAndPurge(t *testing.T) {
	srv := startMailtest(t)
	t.Setenv(mail_sender.SMTPOverrideEnv, srv.SMTPAddr)
	db := mailLogDB(t)
	sender := mail_sender.MailSender{SubjectPrefix: "[Site] "}

	if err := sender.Send(mail_sender.Message().To("ana@example.com").Subject("Hello").Body("<p>Hi</p>").
		Kind("contact.reply").RecordPath("admin:contacts:<1>:")); err != nil {
		t.Fatal(err)
	}
	first := lastLog(t, db)
	time.Sleep(10 * time.Millisecond)
	if err := sender.Resend(first); err != nil {
		t.Fatal(err)
	}
	if !srv.Wait(2, 5*time.Second) {
		t.Fatal("the second did not arrive")
	}
	again := lastLog(t, db)
	if again.ID == first.ID || again.Subject != "[Site] Hello" || again.Body != "<p>Hi</p>" ||
		again.InReplyTo != first.MessageID || again.Kind != "contact.reply" || again.RecordPath != first.RecordPath {
		t.Errorf("sent again: %+v", again)
	}
	if got := srv.To("ana@example.com"); got[1].Subject != "[Site] Hello" ||
		got[1].Header.Get("In-Reply-To") != "<"+first.MessageID+">" {
		t.Errorf("the message sent again: %q %q", got[1].Subject, got[1].Header.Get("In-Reply-To"))
	}

	db.Model(&mail_sender.MailLog{}).Where("id = ?", first.ID).Update("created_at", time.Now().AddDate(0, 0, -400))
	n, err := mail_sender.PurgeMailLogs(db, 365)
	if err != nil || n != 1 {
		t.Errorf("purged %d, %v", n, err)
	}
	if n, _ := mail_sender.PurgeMailLogs(db, 0); n != 0 {
		t.Error("0 days purged")
	}
	var left int64
	db.Model(&mail_sender.MailLog{}).Count(&left)
	if left != 1 {
		t.Errorf("left %d", left)
	}
}
