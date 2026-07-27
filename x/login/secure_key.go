package login

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	rcron "github.com/robfig/cron/v3"
	"github.com/theplant/osenv"
)

// The secure key: a way in for whoever operates the server, past the form
// protection (reCAPTCHA or the built-in challenge, see challenge.go).
//
// It is a FILE on the server — `.secure_key` next to the program. Whoever can
// read it can send its content in the `X-Secure-Key` header and log in without
// the form protection standing in the way: a health check, a script, an
// operator working from a terminal. Somebody who cannot read the file has
// nothing to send.
//
// The key is rewritten on a schedule (RVQ_SECURE_KEY_RENEW, "@daily"), so a
// leaked value stops working on its own.
//
// It skips the FORM PROTECTION only. The account and the password are still
// required, and everything else — the retry count, the permissions — applies as
// usual.

const (
	// DefaultSecureKeyFile is where the key lives when nothing else is set.
	DefaultSecureKeyFile = ".secure_key"

	// SecureKeyHeader carries the key on a request.
	SecureKeyHeader = "X-Secure-Key"

	// SecureKeyFileMode is the only mode accepted: a key anybody on the machine
	// can read is not a key, so a wider mode turns the bypass OFF.
	SecureKeyFileMode fs.FileMode = 0o600
)

// SecureKeyRenew is the cron rule that rewrites the key file. See
// Builder.StartSecureKeyRenewal.
var SecureKeyRenew = osenv.Get("RVQ_SECURE_KEY_RENEW",
	"Cron rule that rewrites the .secure_key file with a fresh random value",
	"@daily")

// SecureKeyFile sets where the key file is (default DefaultSecureKeyFile).
func (b *Builder) SecureKeyFile(path string) (r *Builder) {
	b.secureKeyFile = path
	return b
}

// GetSecureKeyFile is the path of the key file.
func (b *Builder) GetSecureKeyFile() string {
	if b.secureKeyFile == "" {
		return DefaultSecureKeyFile
	}
	return b.secureKeyFile
}

// SecureKeyMatches reports whether the request carries the key the file holds.
//
// False whenever anything is off: no file, an empty one, a mode wider than
// SecureKeyFileMode, no header, or a header that does not match. The comparison
// is constant time — the answer must not depend on how much of the key is right.
func (b *Builder) SecureKeyMatches(r *http.Request) bool {
	sent := strings.TrimSpace(r.Header.Get(SecureKeyHeader))
	if sent == "" {
		return false
	}

	key, err := b.readSecureKey()
	if err != nil || key == "" {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(sent), []byte(key)) == 1
}

func (b *Builder) readSecureKey() (string, error) {
	path := b.GetSecureKeyFile()

	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if mode := info.Mode().Perm(); mode != SecureKeyFileMode {
		b.warnSecureKeyMode.Do(func() {
			log.Printf("WARNING: %s has mode %04o and is ignored — a key other "+
				"accounts on this machine can read is not a key. chmod %04o it.",
				path, mode, SecureKeyFileMode)
		})
		return "", nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// RenewSecureKey writes a fresh random key, with SecureKeyFileMode, and returns
// it. Creates the file when it is not there yet.
func (b *Builder) RenewSecureKey() (key string, err error) {
	raw := make([]byte, 32)
	if _, err = rand.Read(raw); err != nil {
		return
	}
	key = base64.RawURLEncoding.EncodeToString(raw)

	path := b.GetSecureKeyFile()
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err = os.MkdirAll(dir, 0o700); err != nil {
			return "", err
		}
	}

	// written to a temporary file and moved into place, so a reader never sees
	// half a key — and created 0600 from the start, never wider for an instant
	tmp, err := os.CreateTemp(filepath.Dir(path), ".secure_key-*")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp.Name())

	if err = tmp.Chmod(SecureKeyFileMode); err != nil {
		tmp.Close()
		return "", err
	}
	if _, err = tmp.WriteString(key + "\n"); err != nil {
		tmp.Close()
		return "", err
	}
	if err = tmp.Close(); err != nil {
		return "", err
	}
	if err = os.Rename(tmp.Name(), path); err != nil {
		return "", err
	}

	b.warnSecureKeyMode = sync.Once{} // the mode is right again
	return key, nil
}

// StartSecureKeyRenewal writes a key now (if there is none) and rewrites it on
// the schedule of SecureKeyRenew — RVQ_SECURE_KEY_RENEW, "@daily" by default,
// or any rule robfig/cron understands ("0 */6 * * *", "@every 12h"). An empty
// rule turns the renewal off, leaving whatever key is on disk.
//
// Pass a cron to schedule it on the application's own; with none, the builder
// keeps one of its own, started here and stopped by StopSecureKeyRenewal.
func (b *Builder) StartSecureKeyRenewal(c ...*rcron.Cron) (err error) {
	rule := strings.TrimSpace(SecureKeyRenew)
	if rule == "" {
		return nil
	}

	if _, err = os.Stat(b.GetSecureKeyFile()); os.IsNotExist(err) {
		if _, err = b.RenewSecureKey(); err != nil {
			return err
		}
	}

	renew := func() {
		if _, err := b.RenewSecureKey(); err != nil {
			log.Printf("could not renew %s: %v", b.GetSecureKeyFile(), err)
		}
	}

	if len(c) > 0 && c[0] != nil {
		_, err = c[0].AddFunc(rule, renew)
		return err
	}

	b.secureKeyCron = rcron.New()
	if _, err = b.secureKeyCron.AddFunc(rule, renew); err != nil {
		b.secureKeyCron = nil
		return err
	}
	b.secureKeyCron.Start()
	return nil
}

// StopSecureKeyRenewal stops the cron StartSecureKeyRenewal started. It does
// nothing when the schedule went to the application's own cron.
func (b *Builder) StopSecureKeyRenewal() {
	if b.secureKeyCron != nil {
		b.secureKeyCron.Stop()
		b.secureKeyCron = nil
	}
}
