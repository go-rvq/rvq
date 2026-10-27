package mail_sender

import (
	"context"

	"github.com/go-rvq/rvq/admin/presets"
)

// smtpReuseFields are the fields of InboundSettings a bit of its Reuse takes
// from the SMTP's: shown only when it is not set.
var smtpReuseFields = map[string]SMTPReuse{
	"User":     SMTPReuseUser,
	"Password": SMTPReusePassword,
	"Server":   SMTPReuseHost,
}

// inboundOf is the InboundSettings of obj: itself, or the settings'.
func inboundOf(obj any) (s InboundSettings) {
	switch o := obj.(type) {
	case *InboundSettings:
		s = *o
	case *MailSender:
		s = o.Inbound
	}
	return
}

// smtpReuseOptions are the bits of SMTPReuse as options of a field of bits
// (presets.BitsComponentFunc): the label and the description of each, of the
// messages.
func smtpReuseOptions() []presets.BitOption {
	opt := func(bit SMTPReuse, label, desc func(*Messages) string) presets.BitOption {
		return presets.BitOption{Bit: uint64(bit),
			Label:       func(ctx context.Context) string { return label(GetMessages(ctx)) },
			Description: func(ctx context.Context) string { return desc(GetMessages(ctx)) }}
	}
	return []presets.BitOption{
		opt(SMTPReuseUser, func(m *Messages) string { return m.SMTPReuseUser }, func(m *Messages) string { return m.SMTPReuseUser_Desc }),
		opt(SMTPReusePassword, func(m *Messages) string { return m.SMTPReusePassword }, func(m *Messages) string { return m.SMTPReusePassword_Desc }),
		opt(SMTPReuseHost, func(m *Messages) string { return m.SMTPReuseHost }, func(m *Messages) string { return m.SMTPReuseHost_Desc }),
	}
}
