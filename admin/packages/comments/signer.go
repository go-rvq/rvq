package comments

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Signer signs the addresses of the attachments: a file reached by who was
// shown it — the comments already checked when drawn —, its link valid for
// that user, for a time. Its key is the process's (NewSigner): a restart
// makes the links of before stale (the page drawn again makes new ones).
type Signer struct{ key []byte }

// NewSigner is a Signer of a random key.
func NewSigner() *Signer {
	k := make([]byte, 32)
	if _, err := rand.Read(k); err != nil {
		panic(err)
	}
	return &Signer{key: k}
}

func (s *Signer) mac(file, user uuid.UUID, exp int64) string {
	m := hmac.New(sha256.New, s.key)
	m.Write([]byte(file.String() + "|" + user.String() + "|" + strconv.FormatInt(exp, 10)))
	return base64.RawURLEncoding.EncodeToString(m.Sum(nil))
}

// Sign is the token of the file for the user, valid for ttl: "<exp>.<mac>".
func (s *Signer) Sign(file, user uuid.UUID, ttl time.Duration) string {
	exp := time.Now().Add(ttl).Unix()
	return strconv.FormatInt(exp, 10) + "." + s.mac(file, user, exp)
}

// Verify says whether token is of the file for the user, not expired.
func (s *Signer) Verify(file, user uuid.UUID, token string) bool {
	e, sig, ok := strings.Cut(token, ".")
	if !ok {
		return false
	}
	exp, err := strconv.ParseInt(e, 10, 64)
	if err != nil || time.Now().Unix() > exp {
		return false
	}
	return hmac.Equal([]byte(sig), []byte(s.mac(file, user, exp)))
}
