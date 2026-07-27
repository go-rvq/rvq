package login

import (
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-rvq/rvq/x/i18n"
	rcron "github.com/robfig/cron/v3"
)

func newSecureKeyBuilder(t *testing.T) *Builder {
	t.Helper()
	b := New(i18n.New())
	b.Secret("a secret worth at least this many bytes")
	b.SecureKeyFile(filepath.Join(t.TempDir(), DefaultSecureKeyFile))
	return b
}

// postWithKey is a login POST carrying nothing the form protection asks for,
// except the secure key header when one is given.
func postWithKey(b *Builder, key string) FailCode {
	r := httptest.NewRequest("POST", "/login", strings.NewReader(url.Values{}.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if key != "" {
		r.Header.Set(SecureKeyHeader, key)
	}
	return b.VerifyFormProtection(r)
}

func TestSecureKeyLetsTheHolderPastTheProtection(t *testing.T) {
	b := newSecureKeyBuilder(t)

	// with no key file, a POST carrying no challenge is refused as always
	if code := postWithKey(b, ""); code == 0 {
		t.Fatal("a bare POST was accepted with no protection in the way")
	}

	key, err := b.RenewSecureKey()
	if err != nil {
		t.Fatal(err)
	}

	if code := postWithKey(b, key); code != 0 {
		t.Errorf("the key holder was refused: fail code %d", code)
	}
	if code := postWithKey(b, key+"x"); code == 0 {
		t.Error("a wrong key was accepted")
	}
	if code := postWithKey(b, ""); code == 0 {
		t.Error("a POST with no key was accepted while a key file exists")
	}
}

func TestSecureKeyFileIsWritten0600(t *testing.T) {
	b := newSecureKeyBuilder(t)

	if _, err := b.RenewSecureKey(); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(b.GetSecureKeyFile())
	if err != nil {
		t.Fatal(err)
	}
	if mode := info.Mode().Perm(); mode != SecureKeyFileMode {
		t.Errorf("mode = %04o, want %04o", mode, SecureKeyFileMode)
	}
}

func TestSecureKeyIgnoredWhenOthersCanReadIt(t *testing.T) {
	b := newSecureKeyBuilder(t)

	key, err := b.RenewSecureKey()
	if err != nil {
		t.Fatal(err)
	}
	if err = os.Chmod(b.GetSecureKeyFile(), 0o644); err != nil {
		t.Fatal(err)
	}

	// a key the whole machine can read is not a key
	if code := postWithKey(b, key); code == 0 {
		t.Error("a key file readable by others was accepted")
	}

	// and writing it again puts the mode back
	if key, err = b.RenewSecureKey(); err != nil {
		t.Fatal(err)
	}
	if code := postWithKey(b, key); code != 0 {
		t.Errorf("after renewing, the key holder was refused: fail code %d", code)
	}
}

func TestSecureKeyEmptyOrMissingFile(t *testing.T) {
	b := newSecureKeyBuilder(t)

	// missing
	if code := postWithKey(b, "anything"); code == 0 {
		t.Error("a key was accepted with no file on disk")
	}

	// empty, with the right mode: still nothing to match
	if err := os.WriteFile(b.GetSecureKeyFile(), []byte("\n"), SecureKeyFileMode); err != nil {
		t.Fatal(err)
	}
	if code := postWithKey(b, ""); code == 0 {
		t.Error("an empty header matched an empty file")
	}
	if code := postWithKey(b, " "); code == 0 {
		t.Error("a blank header matched an empty file")
	}
}

func TestSecureKeySkipsRecaptchaToo(t *testing.T) {
	b := newSecureKeyBuilder(t)
	b.Recaptcha(true, RecaptchaConfig{SiteKey: "site", SecretKey: "secret"})

	key, err := b.RenewSecureKey()
	if err != nil {
		t.Fatal(err)
	}

	// reCAPTCHA would refuse this POST (no token, and no way to reach Google
	// from a test), the key goes past it
	if code := postWithKey(b, key); code != 0 {
		t.Errorf("the key holder was refused under reCAPTCHA: fail code %d", code)
	}
	if code := postWithKey(b, ""); code == 0 {
		t.Error("a POST with no token and no key was accepted under reCAPTCHA")
	}
}

func TestRenewSecureKeyWritesAFreshValue(t *testing.T) {
	b := newSecureKeyBuilder(t)

	first, err := b.RenewSecureKey()
	if err != nil {
		t.Fatal(err)
	}
	second, err := b.RenewSecureKey()
	if err != nil {
		t.Fatal(err)
	}

	if first == second {
		t.Error("renewing wrote the same key again")
	}
	if len(second) < 32 {
		t.Errorf("the key is %d chars, too short to be one", len(second))
	}

	// the file holds the current key, and the old one stops working
	if code := postWithKey(b, first); code == 0 {
		t.Error("the previous key still works")
	}
	if code := postWithKey(b, second); code != 0 {
		t.Errorf("the current key was refused: fail code %d", code)
	}
}

func TestStartSecureKeyRenewal(t *testing.T) {
	original := SecureKeyRenew
	defer func() { SecureKeyRenew = original }()

	t.Run("writes the file when there is none", func(t *testing.T) {
		b := newSecureKeyBuilder(t)
		SecureKeyRenew = "@daily"

		if err := b.StartSecureKeyRenewal(); err != nil {
			t.Fatal(err)
		}
		defer b.StopSecureKeyRenewal()

		if _, err := os.Stat(b.GetSecureKeyFile()); err != nil {
			t.Errorf("no key file after starting the renewal: %v", err)
		}
	})

	t.Run("schedules on the application's cron", func(t *testing.T) {
		b := newSecureKeyBuilder(t)
		SecureKeyRenew = "@every 1h"

		c := rcron.New()
		if err := b.StartSecureKeyRenewal(c); err != nil {
			t.Fatal(err)
		}
		if len(c.Entries()) != 1 {
			t.Errorf("entries on the given cron = %d, want 1", len(c.Entries()))
		}
	})

	t.Run("renews on schedule", func(t *testing.T) {
		b := newSecureKeyBuilder(t)
		SecureKeyRenew = "@every 100ms"

		if err := b.StartSecureKeyRenewal(); err != nil {
			t.Fatal(err)
		}
		defer b.StopSecureKeyRenewal()

		before, err := os.ReadFile(b.GetSecureKeyFile())
		if err != nil {
			t.Fatal(err)
		}

		deadline := time.Now().Add(3 * time.Second)
		for time.Now().Before(deadline) {
			after, err := os.ReadFile(b.GetSecureKeyFile())
			if err == nil && string(after) != string(before) {
				return // renewed
			}
			time.Sleep(50 * time.Millisecond)
		}
		t.Error("the key was not renewed on schedule")
	})

	t.Run("an empty rule leaves the key alone", func(t *testing.T) {
		b := newSecureKeyBuilder(t)
		SecureKeyRenew = ""

		if err := b.StartSecureKeyRenewal(); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(b.GetSecureKeyFile()); !os.IsNotExist(err) {
			t.Error("a disabled renewal wrote a key file")
		}
	})

	t.Run("a rule that makes no sense is an error", func(t *testing.T) {
		b := newSecureKeyBuilder(t)
		SecureKeyRenew = "not a cron rule"

		if err := b.StartSecureKeyRenewal(); err == nil {
			t.Error("accepted")
			b.StopSecureKeyRenewal()
		}
	})
}
