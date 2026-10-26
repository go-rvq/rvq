package mail_sender

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// SMTPOverrideEnv names the SMTP server ("host:port") every message goes to
// instead of the settings' sender — the local development's and the tests'
// (package mailtest) —: an installation that has it mails nobody real,
// whatever its database (a copy of production's) says.
const SMTPOverrideEnv = "RVQ_MAIL_SMTP"

// SMTPOverride is the server of SMTPOverrideEnv; "" when there is none.
func SMTPOverride() string { return strings.TrimSpace(os.Getenv(SMTPOverrideEnv)) }

// overrideSender is the sender to the server addr ("host:port"), from the
// address the settings send from.
func (c *MailSender) overrideSender(addr string) (Sender, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, fmt.Errorf("%s=%q: %w", SMTPOverrideEnv, addr, err)
	}
	p, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("%s=%q: the port: %w", SMTPOverrideEnv, addr, err)
	}
	return &SmtpSender{Server: host, Port: uint16(p), User: "override", Password: "override", FromMail: c.FromAddress()}, nil
}

// FromAddress is the address the settings send from: the SMTP's, else the
// Gmail account's; "noreply@localhost" when none is set.
func (c *MailSender) FromAddress() string {
	switch {
	case c.SMTP.FromMail != "":
		return c.SMTP.FromMail
	case c.SMTP.User != "" && strings.Contains(c.SMTP.User, "@"):
		return c.SMTP.User
	case c.Gmail.User.Data != nil && c.Gmail.User.Data.Email != "":
		return c.Gmail.User.Data.Email
	}
	return "noreply@localhost"
}

// ErrRateLimited is a message refused because a recipient of it, outside the
// installation (MessageBuilder.External), got the most it may in the window.
var ErrRateLimited = errors.New("mail sender: too many messages to this address, try later")

// recipientLimiter counts the messages each address outside got, the last
// hour and the last day (MailSender.MaxPerRecipientHour, …PerDay), in this
// process.
type recipientLimiter struct {
	mu   sync.Mutex
	sent map[string][]time.Time
	now  func() time.Time
}

var limiter = &recipientLimiter{sent: map[string][]time.Time{}, now: time.Now}

// take reports whether a message may go to each address of to — perHour and
// perDay (0: no limit) —, and counts it when it may.
func (l *recipientLimiter) take(to []string, perHour, perDay uint) error {
	if perHour == 0 && perDay == 0 {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	keys := make([]string, len(to))
	for i, addr := range to {
		key := strings.ToLower(strings.TrimSpace(addr))
		keys[i] = key
		var kept []time.Time
		var hour uint
		for _, t := range l.sent[key] {
			if now.Sub(t) < 24*time.Hour {
				kept = append(kept, t)
				if now.Sub(t) < time.Hour {
					hour++
				}
			}
		}
		l.sent[key] = kept
		if (perHour > 0 && hour >= perHour) || (perDay > 0 && uint(len(kept)) >= perDay) {
			return fmt.Errorf("%w: %s", ErrRateLimited, addr)
		}
	}
	for _, key := range keys {
		l.sent[key] = append(l.sent[key], now)
	}
	return nil
}
