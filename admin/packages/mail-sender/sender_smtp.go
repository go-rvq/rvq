package mail_sender

import (
	"fmt"
	"net/url"

	"github.com/wneessen/go-mail"
)

type SmtpSender struct {
	TLS      bool
	Server   string
	Port     uint16
	FromMail string
	User     string
	Password string
}

func (e *SmtpSender) IsValid() bool {
	return e.Server != "" && e.Port != 0 && e.Password != "" && e.User != ""
}

func (e *SmtpSender) HostPort() string {
	return fmt.Sprintf("%s:%d", e.Server, e.Port)
}

func (e *SmtpSender) String() string {
	u := url.URL{
		Scheme: "smtp",
		Host:   fmt.Sprintf("%s:%d", e.Server, e.Port),
		RawQuery: url.Values{
			"from": {e.User},
		}.Encode(),
	}
	return u.String()
}

func (e *SmtpSender) Send(b *MessageBuilder) (err error) {
	from := e.FromMail
	if from == "" {
		from = e.User
	}

	var m *mail.Msg
	if m, err = b.Build(); err != nil {
		return
	}

	if err = m.From(from); err != nil {
		return
	}

	var c *mail.Client
	c, err = mail.NewClient(e.Server,
		mail.WithPort(int(e.Port)), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(e.User), mail.WithPassword(e.Password))

	if err != nil {
		err = fmt.Errorf("create mail client: %s", err)
		return
	}

	if !e.TLS {
		c.SetTLSPolicy(mail.NoTLS)
	}

	c.SetSMTPAuth(mail.SMTPAuthPlainNoEnc)

	// Finally let's send out the mail
	if err = c.DialAndSend(m); err != nil {
		err = fmt.Errorf("client dial and sent: %s", err)
	}

	return
}
