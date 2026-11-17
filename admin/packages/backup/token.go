package backup

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// sessionToken is what a client's token says: its session, the state it
// downloads, until when.
type sessionToken struct {
	Session uuid.UUID
	State   string
	Expires time.Time
}

// ErrToken: a token not signed by the server, of another state, or expired.
var ErrToken = errors.New("backup: invalid token")

// signToken is t signed (HMAC-SHA256 of secret): "<payload>.<mac>", both
// base64url.
func signToken(secret []byte, t sessionToken) string {
	payload := t.Session.String() + "|" + t.State + "|" + strconv.FormatInt(t.Expires.Unix(), 10)
	enc := base64.RawURLEncoding
	return enc.EncodeToString([]byte(payload)) + "." + enc.EncodeToString(mac(secret, payload))
}

func mac(secret []byte, payload string) []byte {
	m := hmac.New(sha256.New, secret)
	m.Write([]byte(payload))
	return m.Sum(nil)
}

// parseToken is the token s, signed by secret and not expired at now.
func parseToken(secret []byte, s string, now time.Time) (t sessionToken, err error) {
	enc := base64.RawURLEncoding
	p, m, ok := strings.Cut(s, ".")
	if !ok {
		return t, ErrToken
	}
	payload, err1 := enc.DecodeString(p)
	sig, err2 := enc.DecodeString(m)
	if err1 != nil || err2 != nil || !hmac.Equal(sig, mac(secret, string(payload))) {
		return t, ErrToken
	}
	f := strings.Split(string(payload), "|")
	if len(f) != 3 {
		return t, ErrToken
	}
	if t.Session, err = uuid.Parse(f[0]); err != nil {
		return t, ErrToken
	}
	exp, err := strconv.ParseInt(f[2], 10, 64)
	if err != nil {
		return t, ErrToken
	}
	t.State, t.Expires = f[1], time.Unix(exp, 0)
	if !now.Before(t.Expires) {
		return t, ErrToken
	}
	return t, nil
}
