package login

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-rvq/rvq/x/osenv"
	rcron "github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

// The secure key: a way in for whoever operates the server, past the form
// protection (reCAPTCHA or the built-in challenge, see challenge.go).
//
// It is a FILE on the server — `.secure_key` next to the program, or wherever
// RVQ_SECURE_KEY_FILE points. Whoever can
// read it can send its content in the `X-Secure-Key` header and log in without
// the form protection standing in the way: a health check, a script, an
// operator working from a terminal. Somebody who cannot read the file has
// nothing to send.
//
// The key is rewritten on a schedule (RVQ_SECURE_KEY_RENEW, "@daily"), so a
// leaked value stops working on its own.
//
// On a login it skips the FORM PROTECTION only: the account and the password
// are still required, and everything else — the retry count, the permissions —
// applies as usual.
//
// Where a handler takes SecureKeyBasicAuthMiddleware (the git of the admin)
// and the application turned it on (SecureKeyBasicAuth), it is a login of its
// own: the key as the password of HTTP Basic is the administrator — the
// initial account (InitialUserAccount) —, no password, no user name needed.
// Whoever reads the file operates the server already. Only over HTTPS (or
// from the machine itself): the key travels in each request.

const (
	// DefaultSecureKeyFile is where the key lives when nothing else is set.
	DefaultSecureKeyFile = ".secure_key"

	// SecureKeyHeader carries the key on a request.
	SecureKeyHeader = "X-Secure-Key"

	// SecureKeyFileMode is the only mode accepted: a key anybody on the machine
	// can read is not a key, so a wider mode turns the bypass OFF.
	SecureKeyFileMode fs.FileMode = 0o600
)

var (
	// SecureKeyPath is where the key file is, from the environment. A path given
	// in code (Builder.SecureKeyFile) wins over it.
	SecureKeyPath = osenv.Get("RVQ_SECURE_KEY_FILE",
		"Path of the file holding the login secure key, which must have mode 0600",
		DefaultSecureKeyFile)

	// TrustedProxies are the reverse proxies whose X-Forwarded-Proto is
	// believed (secureTransport), besides the loopback.
	TrustedProxies = osenv.Get("RVQ_TRUSTED_PROXIES",
		"Addresses or networks (CIDR), comma-separated, of the reverse proxies whose X-Forwarded-Proto is believed — the loopback always is; in a container, the network of its proxy (e.g. 172.16.0.0/12)",
		"")

	// SecureKeyRenew is the cron rule that rewrites the key file. See
	// Builder.StartSecureKeyRenewal.
	SecureKeyRenew = osenv.Get("RVQ_SECURE_KEY_RENEW",
		"Cron rule that rewrites the secure key file with a fresh random value",
		"@daily")
)

// SecureKeyFile sets where the key file is, over whatever the environment says.
func (b *Builder) SecureKeyFile(path string) (r *Builder) {
	b.secureKeyFile = path
	return b
}

// GetSecureKeyFile is the path of the key file: what was set in code, else
// RVQ_SECURE_KEY_FILE, else DefaultSecureKeyFile.
func (b *Builder) GetSecureKeyFile() string {
	if b.secureKeyFile != "" {
		return b.secureKeyFile
	}
	if SecureKeyPath != "" {
		return SecureKeyPath
	}
	return DefaultSecureKeyFile
}

// SecureKeyMatches reports whether the request carries the key the file holds.
//
// False whenever anything is off: no file, an empty one, a mode wider than
// SecureKeyFileMode, no header, or a header that does not match. The comparison
// is constant time — the answer must not depend on how much of the key is right.
func (b *Builder) SecureKeyMatches(r *http.Request) bool {
	return b.isSecureKey(r.Header.Get(SecureKeyHeader))
}

// isSecureKey reports whether sent is the key the file holds (as
// SecureKeyMatches: false whenever anything is off; constant time).
func (b *Builder) isSecureKey(sent string) bool {
	sent = strings.TrimSpace(sent)
	if sent == "" {
		return false
	}

	key, err := b.readSecureKey()
	if err != nil || key == "" {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(sent), []byte(key)) == 1
}

// SecureKeyBasicAuthMiddleware is BasichAuthMiddleware that takes the secure
// key too: a request whose HTTP Basic password is the key is the
// administrator's — the initial account (InitialUserAccount), not locked —,
// whatever user it names, with no session (as an access key's,
// IsKeyAuthenticated). Any other request goes by BasichAuthMiddleware.
func (b *Builder) SecureKeyBasicAuthMiddleware(next http.Handler) http.Handler {
	basic := b.BasichAuthMiddleware(next)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, pass, ok := r.BasicAuth(); !ok || !b.isSecureKey(pass) {
			basic.ServeHTTP(w, r)
			return
		}
		// the key: never tried as an account's password (no retry counted)
		user, err := b.secureKeyUser()
		switch {
		case b.secureKeyBasicAuth == nil || !b.secureKeyBasicAuth():
			err = ErrSecureKeyOff
		case !secureTransport(r):
			err = ErrSecureKeyTransport
		}
		if err != nil {
			log.Printf("secure key refused: %s %s from %s: %v", r.Method, r.URL.Path, r.RemoteAddr, err)
			w.Header().Set("WWW-Authenticate", `Basic realm="secure key"`)
			http.Error(w, "ERROR: "+err.Error(), http.StatusUnauthorized)
			return
		}
		log.Printf("secure key: %s %s from %s, as %s", r.Method, r.URL.Path, r.RemoteAddr, b.initialUserAccount)
		ctx := context.WithValue(context.WithValue(WithAuth(r.Context(), AuthSecureKey), keyAuthenticatedKey{}, true), UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// SecureKeyBasicAuth sets whether the secure key is a login where a handler
// takes it (SecureKeyBasicAuthMiddleware), asked at each request: the
// application's switch. Unset, it is not.
func (b *Builder) SecureKeyBasicAuth(enabled func() bool) *Builder {
	b.secureKeyBasicAuth = enabled
	return b
}

// secureTransport reports whether r came over HTTPS: its own TLS, or — from
// the loopback or a trusted proxy (TrustedProxies) — the X-Forwarded-Proto
// the proxy says; a request from the machine itself with no proxy, too.
func secureTransport(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	ip := net.ParseIP(host)
	if ip == nil || !(ip.IsLoopback() || trustedProxy(ip)) {
		return false
	}
	proto := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-Proto"), ",")[0])
	if proto == "" {
		return ip.IsLoopback()
	}
	return strings.EqualFold(proto, "https")
}

// trustedProxy reports whether ip is of TrustedProxies.
func trustedProxy(ip net.IP) bool {
	for _, p := range strings.Split(TrustedProxies, ",") {
		if p = strings.TrimSpace(p); p == "" {
			continue
		}
		if _, n, err := net.ParseCIDR(p); err == nil {
			if n.Contains(ip) {
				return true
			}
		} else if pip := net.ParseIP(p); pip != nil && pip.Equal(ip) {
			return true
		}
	}
	return false
}

// secureKeyUser is the user of the secure key: the administrator, the
// initial account.
func (b *Builder) secureKeyUser() (any, error) {
	if b.userModel == nil || b.initialUserAccount == "" {
		return nil, ErrSecureKeyNoAdmin
	}
	user, err := b.userModel.(UserPasser).FindUser(b.db, b.newUserObject(), b.initialUserAccount)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	if user.(UserPasser).GetLocked() {
		return nil, ErrUserLocked
	}
	return user, nil
}

// The refusals of the secure key as a login.
var (
	// ErrSecureKeyNoAdmin is the secure key with no administrator to be: no
	// users, or no initial account.
	ErrSecureKeyNoAdmin = errors.New("the secure key has no administrator to log in as (InitialUserAccount)")
	// ErrSecureKeyOff is the secure key where the application has it off
	// (SecureKeyBasicAuth).
	ErrSecureKeyOff = errors.New("the access by the secure key is off")
	// ErrSecureKeyTransport is the secure key over plain HTTP.
	ErrSecureKeyTransport = errors.New("the secure key is taken over HTTPS only")
)

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

// StartSecureKeyRenewal rewrites the key file on the schedule of SecureKeyRenew
// — RVQ_SECURE_KEY_RENEW, "@daily" by default, or any rule robfig/cron
// understands ("0 */6 * * *", "@every 12h"). An empty rule turns the renewal
// off, leaving whatever key is on disk.
//
// It never CREATES the file: the bypass exists because somebody decided it
// should, so a missing key file only earns a line on the terminal saying how to
// make one. From the moment it is there, the renewal rewrites it.
//
// Pass a cron to schedule it on the application's own; with none, the builder
// keeps one of its own, started here and stopped by StopSecureKeyRenewal.
func (b *Builder) StartSecureKeyRenewal(c ...*rcron.Cron) (err error) {
	b.WarnSecureKeyMissing()

	rule := strings.TrimSpace(SecureKeyRenew)
	if rule == "" {
		return nil
	}

	if len(c) > 0 && c[0] != nil {
		_, err = c[0].AddFunc(rule, b.renewScheduled)
		return err
	}
	return b.RescheduleSecureKeyRenewal(rule)
}

// RescheduleSecureKeyRenewal renews the key by rule from now on, in place of
// the rule it was renewed by (StartSecureKeyRenewal, or a former call): the
// application's setting, changed. An empty rule stops the renewal.
func (b *Builder) RescheduleSecureKeyRenewal(rule string) error {
	b.secureKeyCronMu.Lock()
	defer b.secureKeyCronMu.Unlock()
	rule = strings.TrimSpace(rule)
	var next *rcron.Cron
	if rule != "" {
		next = rcron.New()
		if _, err := next.AddFunc(rule, b.renewScheduled); err != nil {
			return err
		}
	}
	if b.secureKeyCron != nil {
		b.secureKeyCron.Stop()
	}
	b.secureKeyCron = next
	if next != nil {
		next.Start()
	}
	return nil
}

// renewScheduled is the scheduled renewal: nothing to rewrite while there is
// no key — the file is the operator's to create (or the application's, when
// asked: RenewSecureKey).
func (b *Builder) renewScheduled() {
	if _, err := os.Stat(b.GetSecureKeyFile()); err != nil {
		return
	}
	if _, err := b.RenewSecureKey(); err != nil {
		log.Printf("could not renew %s: %v", b.GetSecureKeyFile(), err)
	}
}

// SecureKeyFileInfo is what is known of the key file, never its content.
type SecureKeyFileInfo struct {
	// Path is where it is, absolute.
	Path string
	// Exists says whether it is there.
	Exists bool
	// ModeOK says whether its mode is SecureKeyFileMode: a wider one turns
	// the key off.
	ModeOK bool
	// Mode is its mode.
	Mode fs.FileMode
	// ModTime is when it was last written (renewed).
	ModTime time.Time
}

// SecureKeyInfo is what is known of the key file (SecureKeyFileInfo).
func (b *Builder) SecureKeyInfo() SecureKeyFileInfo {
	path := b.GetSecureKeyFile()
	if p, err := filepath.Abs(path); err == nil {
		path = p
	}
	i := SecureKeyFileInfo{Path: path}
	if st, err := os.Stat(path); err == nil {
		i.Exists, i.Mode, i.ModTime = true, st.Mode().Perm(), st.ModTime()
		i.ModeOK = i.Mode == SecureKeyFileMode
	}
	return i
}

// WarnSecureKeyMissing says on the terminal that there is no secure key, and
// what a key file has to look like. Called by StartSecureKeyRenewal; the file
// is never created for you.
func (b *Builder) WarnSecureKeyMissing() {
	path := b.GetSecureKeyFile()
	if _, err := os.Stat(path); err == nil {
		return
	}

	abs := path
	if p, err := filepath.Abs(path); err == nil {
		abs = p
	}

	log.Printf("NOTE: no secure key at %s, so every login goes through the form "+
		"protection. To let automation past it, create the file with mode %04o "+
		"(it is ignored with any other mode) and send its content in the %s "+
		"header:\n    (umask 177 && head -c 32 /dev/urandom | base64 | tr -d '=' > %s)",
		abs, SecureKeyFileMode, SecureKeyHeader, abs)
}

// StopSecureKeyRenewal stops the cron StartSecureKeyRenewal started. It does
// nothing when the schedule went to the application's own cron.
func (b *Builder) StopSecureKeyRenewal() {
	b.secureKeyCronMu.Lock()
	defer b.secureKeyCronMu.Unlock()
	if b.secureKeyCron != nil {
		b.secureKeyCron.Stop()
		b.secureKeyCron = nil
	}
}
