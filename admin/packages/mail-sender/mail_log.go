package mail_sender

import (
	"errors"
	"log"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// The central of the mails sent: every message the MailSender sends, or
// fails to, is a MailLog — whom to, what about, by which method, how it
// went, the conversation it is of (MessageID, InReplyTo), the record it is
// of (RecordPath) and who made it be sent —, its sensitive content hidden
// (MessageBuilder.Sensitive, HideBody, and the values of tokens and keys in
// its links). It is only read: nobody edits a message sent.

// The statuses of a MailLog.
const (
	MailLogSent   = "sent"
	MailLogFailed = "failed"
)

// KindTest is the kind of the test message of the settings (the action
// TestSendMail).
const KindTest = "mail_sender.test"

// MethodOverride is the method of a message sent to the server of the
// override (SMTPOverrideEnv): the local development's, the tests'.
const MethodOverride = "OVERRIDE"

// Hidden is what stands for what a MailLog does not keep.
const Hidden = "••••••"

// MailLog is a message sent, or that failed to be.
type MailLog struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" admin:"-"`
	CreatedAt time.Time `gorm:"index"`
	// Kind is what it is ("contact.received"); RecordPath the record it is
	// of (presets record path)
	Kind       string `gorm:"size:64;index"`
	RecordPath string `gorm:"index"`
	// Recipients are its To, joined by ", "
	Recipients string
	ReplyTo    string
	Subject    string
	// MessageID is its Message-ID; InReplyTo the one it answers: its
	// conversation
	MessageID string `gorm:"index"`
	InReplyTo string `gorm:"index"`
	// Body is its HTML, its sensitive content hidden; BodyHidden says it
	// was not kept at all (HideBody); Redacted that something was hidden —
	// it cannot be sent again as it was
	Body       string
	BodyHidden bool
	Redacted   bool
	// Method is the sender ("GMAIL", "SMTP", MethodOverride)
	Method string `gorm:"size:16"`
	// Status is MailLogSent or MailLogFailed, Error why
	Status string `gorm:"size:16;index"`
	Error  string
	// UserID is who made it be sent; RemoteIP and UserAgent, where from
	UserID     string `gorm:"index"`
	RemoteIP   string
	UserAgent  string
	LocaleCode string `gorm:"size:16"`
}

func (MailLog) TableName() string {
	return "mail_sender_logs"
}

// String is its title: its subject, or its kind.
func (l *MailLog) String() string {
	if l.Subject != "" {
		return l.Subject
	}
	return l.Kind
}

// ConversationID is the Message-ID its conversation is grouped by: the one
// it answers, or its own.
func (l *MailLog) ConversationID() string {
	if l.InReplyTo != "" {
		return l.InReplyTo
	}
	return l.MessageID
}

var recorder atomic.Pointer[gorm.DB]

// SetRecorder is where the MailLogs are written (New sets it); nil, none
// are.
func SetRecorder(db *gorm.DB) {
	recorder.Store(db)
}

// secretParams are the values of the parameters of a link that are
// secrets: a token, a key, a code, a password.
var secretParams = regexp.MustCompile(`(?i)([?&](?:token|key|code|password|secret|signature|sig)=)[^&"'\s<>]+`)

// redact is s with the sensitive values of b, and the secrets of its links,
// hidden; changed says something was.
func (b *MessageBuilder) redact(s string) (r string, changed bool) {
	r = s
	for _, v := range b.sensitive {
		if v = strings.TrimSpace(v); v != "" {
			r = strings.ReplaceAll(r, v, Hidden)
		}
	}
	r = secretParams.ReplaceAllString(r, "${1}"+Hidden)
	return r, r != s
}

// Log is the MailLog of b, sent by method, err why it failed.
func (b *MessageBuilder) Log(method string, err error) *MailLog {
	l := &MailLog{
		ID:         uuid.New(),
		Kind:       b.kind,
		RecordPath: b.recordPath,
		Recipients: strings.Join(b.to, ", "),
		ReplyTo:    b.replyTo,
		MessageID:  strings.Trim(b.messageID, "<>"),
		InReplyTo:  strings.Trim(b.inReplyTo, "<>"),
		BodyHidden: b.hideBody,
		Method:     method,
		Status:     MailLogSent,
		UserID:     b.userID,
		RemoteIP:   b.remoteIP,
		UserAgent:  b.userAgent,
		LocaleCode: b.localeCode,
	}
	var subjectChanged, bodyChanged bool
	l.Subject, subjectChanged = b.redact(b.FormattedSubject())
	if !b.hideBody {
		l.Body, bodyChanged = b.redact(b.body)
	}
	l.Redacted = b.hideBody || subjectChanged || bodyChanged
	if err != nil {
		l.Status, l.Error = MailLogFailed, err.Error()
	}
	return l
}

// record writes the MailLog of b, when there is a recorder: a failure to
// write it is logged, never the sending's.
func record(b *MessageBuilder, method string, err error) {
	db := recorder.Load()
	if db == nil {
		return
	}
	if werr := db.Create(b.Log(method, err)).Error; werr != nil {
		log.Printf("mail sender: the log of the message: %v", werr)
	}
}

// PurgeMailLogs deletes the MailLogs older than days (0: none), the number
// deleted.
func PurgeMailLogs(db *gorm.DB, days uint) (int64, error) {
	if days == 0 {
		return 0, nil
	}
	r := db.Where("created_at < ?", time.Now().AddDate(0, 0, -int(days))).Delete(&MailLog{})
	return r.RowsAffected, r.Error
}

// ErrRedacted is the error of sending again a MailLog whose content was
// hidden: what it had is not there to send.
var ErrRedacted = errors.New("mail sender: the message had content hidden: it cannot be sent again")

// Resend sends the message of l again, by the settings c, as it was —
// subject and body —, in its conversation: its own MailLog.
func (c *MailSender) Resend(l *MailLog) error {
	if l.Redacted {
		return ErrRedacted
	}
	var to []string
	for _, a := range strings.Split(l.Recipients, ",") {
		if a = strings.TrimSpace(a); a != "" {
			to = append(to, a)
		}
	}
	b := Message().To(to...).ReplyTo(l.ReplyTo).Subject(l.Subject).Body(l.Body).Verbatim().
		Kind(l.Kind).RecordPath(l.RecordPath).Locale(l.LocaleCode)
	if id := l.ConversationID(); id != "" {
		b.InReplyTo(id).References(id)
	}
	return c.Send(b)
}
