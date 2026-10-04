// Package accesskeys is the access keys of the users: codes an automation —
// git, WebDAV, a script — authenticates with in a user's name, each request
// alone (no session). A key's permissions only restrict: a request by it is
// allowed what the user is allowed AND the key allows (perm.WithRestriction).
// Each request by a key is in its history (AccessKeyUse), and in the user's
// access log.
package accesskeys

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base32"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AccessKey is a user's access key. Its code is shown once, when made: only
// its hash is kept, and its Prefix — the public part, that finds it.
type AccessKey struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// UserID is whose key it is.
	UserID uuid.UUID `gorm:"type:uuid;index;not null"`
	// Name says what it is for ("deploy from the notebook").
	Name        string `gorm:"not null"`
	Description string
	// Prefix is the public part of the code: what finds the key.
	Prefix string `gorm:"uniqueIndex;not null"`
	// CodeHash is the hash of the code (CodeHash): the code is never kept.
	CodeHash string `gorm:"not null" json:"-"`
	// ExpiresAt is when it stops working.
	ExpiresAt time.Time `gorm:"not null"`
	// Enabled: no default in the table — a default would turn a key made
	// disabled into an enabled one; a new key's form starts enabled.
	Enabled bool `gorm:"not null"`
	// LastUsedAt and LastUsedIP are of its last request.
	LastUsedAt *time.Time
	LastUsedIP string
	// CreatedByID is who made it: the user, or an administrator.
	CreatedByID *uuid.UUID `gorm:"type:uuid"`

	// Permissions restrict it: the policies whose ReferID is its id, its
	// subject Subject(); allows only.
	Permissions []*perm.DefaultDBPolicy `gorm:"-:migration;foreignKey:ReferID"`

	// Code is the code of a key just made — shown once —, never stored.
	Code string `gorm:"-"`
}

func (k *AccessKey) String() string { return k.Name }

// BeforeCreate gives it an id.
func (k *AccessKey) BeforeCreate(*gorm.DB) error {
	if k.ID == uuid.Nil {
		k.ID = uuid.New()
	}
	return nil
}

// Subject is the subject of its policies: "key:<id>".
func (k *AccessKey) Subject() string { return SubjectOf(k.ID) }

// SubjectOf is the subject of the policies of the key id.
func SubjectOf(id uuid.UUID) string { return "key:" + id.String() }

// Usable says whether the key works now: enabled, not expired, not deleted.
func (k *AccessKey) Usable(now time.Time) bool {
	return k.Enabled && !k.DeletedAt.Valid && now.Before(k.ExpiresAt)
}

// AccessKeyUse is a request by a key: its history.
type AccessKeyUse struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `gorm:"index"`

	KeyID     uuid.UUID `gorm:"type:uuid;index;not null"`
	UserID    uuid.UUID `gorm:"type:uuid;index;not null"`
	IP        string
	UserAgent string
	// Place is where IP is, when it is known.
	Place string
	// Kind is what it reached: "git", "webdav", "api".
	Kind   string `gorm:"index"`
	Method string
	Path   string
	// Status is how it was answered (200, 403…).
	Status int
}

// BeforeCreate gives it an id.
func (u *AccessKeyUse) BeforeCreate(*gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// CodePrefix is how the codes of keys begin: "hck_<prefix>_<secret>".
const CodePrefix = "hck_"

var codeEncoding = base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567").WithPadding(base32.NoPadding)

func randomText(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return codeEncoding.EncodeToString(b)
}

// NewCode is a new code — "hck_<prefix>_<secret>" — and its prefix.
func NewCode() (code, prefix string) {
	prefix = randomText(6)[:10]
	return CodePrefix + prefix + "_" + randomText(20), prefix
}

// ErrBadCode is a code not of the format of the keys.
var ErrBadCode = errors.New("not an access key's code")

// ParseCode is the prefix of code, when it is a key's.
func ParseCode(code string) (prefix string, err error) {
	rest, ok := strings.CutPrefix(code, CodePrefix)
	if !ok {
		return "", ErrBadCode
	}
	prefix, secret, ok := strings.Cut(rest, "_")
	if !ok || len(prefix) != 10 || len(secret) < 20 {
		return "", ErrBadCode
	}
	return prefix, nil
}

// CodeHash is the hash a code is kept as.
func CodeHash(code string) string {
	sum := sha256.Sum256([]byte(code))
	return hex.EncodeToString(sum[:])
}

// Matches says whether code is the key's (in constant time).
func (k *AccessKey) Matches(code string) bool {
	return subtle.ConstantTimeCompare([]byte(CodeHash(code)), []byte(k.CodeHash)) == 1
}
