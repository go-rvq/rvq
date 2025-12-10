package mail_sender

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-rvq/rvq/thirdpart/gorm/datatypes"
	"github.com/markbates/goth"
	"github.com/wneessen/go-mail"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type GmailToken struct {
	// AccessToken is the token that authorizes and authenticates
	// the requests.
	AccessToken string `json:"access_token"`

	// TokenType is the type of token.
	// The Type method returns either this or "Bearer", the default.
	TokenType string `json:"token_type,omitempty"`

	// RefreshToken is a token that's used by the application
	// (as opposed to the user) to refresh the access token
	// if it expires.
	RefreshToken string `json:"refresh_token,omitempty"`

	// Expiry is the optional expiration time of the access token.
	//
	// If zero, [TokenSource] implementations will reuse the same
	// token forever and RefreshToken or equivalent
	// mechanisms for that TokenSource will not be used.
	Expiry time.Time `json:"expiry,omitempty"`

	// ExpiresIn is the OAuth2 wire format "expires_in" field,
	// which specifies how many seconds later the token expires,
	// relative to an unknown time base approximately around "now".
	// It is the application's responsibility to populate
	// `Expiry` from `ExpiresIn` when required.
	ExpiresIn int64 `json:"expires_in,omitempty"`

	IDToken string `json:"id_token,omitempty"`
}

func (t *GmailToken) Token() (ot *oauth2.Token) {
	ot = &oauth2.Token{
		AccessToken:  t.AccessToken,
		TokenType:    t.TokenType,
		RefreshToken: t.RefreshToken,
		Expiry:       t.Expiry,
		ExpiresIn:    t.ExpiresIn,
	}
	ot.WithExtra(map[string]interface{}{
		"id_token": t.IDToken,
	})
	return
}

type GmailSenderCredentials struct {
	ClientID  string `json:"client_id,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
	Raw       []byte `json:"-"`
}

func (c GmailSenderCredentials) MarshalJSON() ([]byte, error) {
	return c.Raw, nil
}

func (c GmailSenderCredentials) IsZero() bool {
	return len(c.Raw) == 0
}

func (c *GmailSenderCredentials) UnmarshalJSON(data []byte) (err error) {
	var n GmailSenderCredentials
	if len(data) > 0 {
		var dot struct {
			Installed struct {
				ClientID  string `json:"client_id"`
				ProjectID string `json:"project_id"`
			} `json:"installed"`
		}
		if err = json.Unmarshal(data, &dot); err != nil {
			return
		}

		if dot.Installed.ClientID == "" {
			return Messages_en_US.ErrGmailSenderCredentialsInvalid
		}

		n.ClientID = dot.Installed.ClientID
		n.ProjectID = dot.Installed.ProjectID
		n.Raw = make([]byte, len(data))
		copy(n.Raw, data)
	}
	*c = n
	return
}

type GmailSender struct {
	User        datatypes.NullJSONType[*goth.User]             `admin:"-"`
	Token       datatypes.NullJSONType[*GmailToken]            `admin:"-"`
	Credentials datatypes.NullJSONType[GmailSenderCredentials] `admin:"-"`
	CallbackURI string                                         `admin:"-"`
}

func (e *GmailSender) IsValid() bool {
	return e.Token.Data != nil
}

func (e *GmailSender) String() string {
	if e.User.Data != nil {
		return e.User.Data.Email
	}
	return ""
}

func (e *GmailSender) Config() (config *oauth2.Config, err error) {
	config, err = google.ConfigFromJSON(e.Credentials.Data.Raw,
		"email", "profile",
		gmail.MailGoogleComScope,
		gmail.GmailModifyScope,
		gmail.GmailComposeScope,
		gmail.GmailSendScope,
	)

	if err != nil {
		err = fmt.Errorf("Unable to parse client credentials to config: %v", err)
	}

	if e.CallbackURI != "" {
		config.RedirectURL = e.CallbackURI
	}

	return
}

func (e *GmailSender) Send(b *MessageBuilder) (err error) {
	var config *oauth2.Config
	if config, err = e.Config(); err != nil {
		return
	}

	var (
		ctx    = context.Background()
		tok    = e.Token.Data.Token()
		client = config.Client(ctx, tok)
		svc    *gmail.Service
		m      *mail.Msg
	)

	if svc, err = gmail.NewService(ctx, option.WithHTTPClient(client)); err != nil {
		err = fmt.Errorf("Unable to retrieve Gmail client: %v", err)
		return
	}

	if m, err = b.Build(); err != nil {
		return
	}

	if err = m.From(e.User.Data.Email); err != nil {
		return
	}

	var buf bytes.Buffer
	if _, err = m.WriteTo(&buf); err != nil {
		return
	}

	gmsg := &gmail.Message{
		Raw: base64.RawURLEncoding.EncodeToString(buf.Bytes()),
	}

	if _, err = svc.Users.Messages.Send("me", gmsg).Do(); err != nil {
		err = fmt.Errorf("Unable to send Gmail message: %v", err)
		return
	}

	return
}
