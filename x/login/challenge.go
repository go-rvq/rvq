package login

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	h "github.com/go-rvq/htmlgo"
)

// Built-in protection of the login and forget-password forms, used when no
// reCAPTCHA is configured — so an installation without a Google account still
// does not answer to a script hammering the form.
//
// Everything is decided by the SERVER. On purpose: the browser side of a proof
// of work or of a modern captcha needs `crypto.subtle`, `navigator.credentials`
// and friends, which only exist in a secure context — an admin served over plain
// HTTP inside a network would be left with no protection at all. What is used
// here works the same over HTTP and HTTPS, with no third party and no JavaScript:
//
//  1. a SIGNED, single-use challenge the form must send back — a POST that did
//     not come from a form this server rendered has none, and a captured one
//     works exactly once;
//  2. the TIME the form took: a submit in under ChallengeMinAge was not typed by
//     a person, and one after ChallengeMaxAge belongs to a page left open;
//  3. a HONEYPOT field, hidden from the reader: an automated filler completes
//     every field it finds.

const (
	// ChallengeFormKey carries the signed challenge.
	ChallengeFormKey = "__lc"
	// HoneypotFormKey names the field nobody must fill. It looks like something
	// worth filling on purpose.
	HoneypotFormKey = "email_confirmation"
)

const (
	// DefaultChallengeMinAge is how long a human takes, at the very least, to
	// read a form and type into it.
	DefaultChallengeMinAge = 900 * time.Millisecond
	// DefaultChallengeMaxAge is how long a rendered form stays acceptable.
	DefaultChallengeMaxAge = 30 * time.Minute
)

// ChallengeDisabled turns the built-in protection off. It is already off when
// reCAPTCHA is configured — the two do the same job.
func (b *Builder) ChallengeDisabled(v bool) (r *Builder) {
	b.challengeDisabled = v
	return b
}

// ChallengeAges sets how soon a form may be submitted and for how long it stays
// valid. Zero keeps the default.
func (b *Builder) ChallengeAges(min, max time.Duration) (r *Builder) {
	if min > 0 {
		b.challengeMinAge = min
	}
	if max > 0 {
		b.challengeMaxAge = max
	}
	return b
}

// ChallengeEnabled reports whether a form should carry the built-in protection:
// only when it is not disabled and reCAPTCHA is not doing the job already.
func (b *Builder) ChallengeEnabled() bool {
	return !b.challengeDisabled && !b.recaptchaEnabled
}

func (b *Builder) challengeAges() (min, max time.Duration) {
	min, max = b.challengeMinAge, b.challengeMaxAge
	if min <= 0 {
		min = DefaultChallengeMinAge
	}
	if max <= 0 {
		max = DefaultChallengeMaxAge
	}
	return
}

// NewChallenge issues the value a form carries: a random nonce, the instant it
// was issued, and the signature of both — with the same secret that signs the
// session, so there is nothing new to configure.
func (b *Builder) NewChallenge() string {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		panic(err)
	}

	payload := base64.RawURLEncoding.EncodeToString(nonce) + "." +
		strconv.FormatInt(time.Now().UnixNano(), 10)

	return payload + "." + b.signChallenge(payload)
}

func (b *Builder) signChallenge(payload string) string {
	m := hmac.New(sha256.New, []byte(b.secret))
	m.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(m.Sum(nil))
}

// VerifyFormProtection decides whether a login or forget-password POST may go
// on, and is the only place that knows which protection is in charge:
//
//	the secure key   whoever can read the file skips the rest (secure_key.go)
//	reCAPTCHA        when it is configured
//	the challenge    otherwise
//
// It returns 0 when the request is acceptable, and the fail code to show
// otherwise.
func (b *Builder) VerifyFormProtection(r *http.Request) FailCode {
	if b.SecureKeyMatches(r) {
		return 0
	}

	if b.recaptchaEnabled {
		if !recaptchaTokenCheck(b, r.FormValue("token")) {
			return FailCodeIncorrectRecaptchaToken
		}
		return 0
	}

	return b.VerifyChallenge(r)
}

// VerifyChallenge decides whether a POST may go on. It returns 0 when the
// request is acceptable (including when the protection is off), and the fail
// code to show otherwise.
func (b *Builder) VerifyChallenge(r *http.Request) FailCode {
	if !b.ChallengeEnabled() {
		return 0
	}

	// A hidden field nobody sees: whatever filled it was not a reader.
	if strings.TrimSpace(r.FormValue(HoneypotFormKey)) != "" {
		return FailCodeIncorrectChallenge
	}

	value := r.FormValue(ChallengeFormKey)
	if value == "" {
		return FailCodeIncorrectChallenge
	}

	i := strings.LastIndexByte(value, '.')
	if i < 0 {
		return FailCodeIncorrectChallenge
	}

	payload, sig := value[:i], value[i+1:]
	if !hmac.Equal([]byte(sig), []byte(b.signChallenge(payload))) {
		return FailCodeIncorrectChallenge
	}

	j := strings.LastIndexByte(payload, '.')
	if j < 0 {
		return FailCodeIncorrectChallenge
	}

	issued, err := strconv.ParseInt(payload[j+1:], 10, 64)
	if err != nil {
		return FailCodeIncorrectChallenge
	}

	min, max := b.challengeAges()
	age := time.Since(time.Unix(0, issued))

	switch {
	case age < 0 || age > max:
		// a form left open, or an instant from the future
		return FailCodeChallengeExpired
	case age < min:
		// nobody reads and types this fast
		return FailCodeIncorrectChallenge
	}

	// Single use: a challenge captured from a rendered page cannot be replayed.
	if !b.challenges().use(payload, time.Unix(0, issued).Add(max)) {
		return FailCodeIncorrectChallenge
	}
	return 0
}

func (b *Builder) challenges() *challengeStore {
	b.challengeStoreOnce.Do(func() {
		b.challengeStore = &challengeStore{used: map[string]time.Time{}}
	})
	return b.challengeStore
}

// challengeStore remembers the challenges already spent, until they would have
// expired anyway. It is in memory: with several instances behind a balancer a
// challenge could be replayed once per instance, which still leaves every other
// check standing.
type challengeStore struct {
	mu   sync.Mutex
	used map[string]time.Time
}

// use marks a challenge as spent and reports whether it was still unused.
func (s *challengeStore) use(payload string, expiresAt time.Time) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	if _, spent := s.used[payload]; spent {
		return false
	}

	// the map only grows with what is still valid
	for k, exp := range s.used {
		if now.After(exp) {
			delete(s.used, k)
		}
	}

	s.used[payload] = expiresAt
	return true
}

// ChallengeFields are the hidden fields a form must carry for the built-in
// protection, or nil when it is off. Both are plain HTML: no script, nothing
// that needs a secure context.
func (vh *ViewHelper) ChallengeFields() h.HTMLComponent {
	if !vh.b.ChallengeEnabled() {
		return nil
	}

	return h.Components(
		h.Input("").Type("hidden").Attr("name", ChallengeFormKey).
			Attr("value", vh.b.NewChallenge()),

		// Out of the reader's way — and out of the way of anybody using a screen
		// reader, which is what aria-hidden and tabindex are for. A field that is
		// merely `type=hidden` is too obvious to work as a honeypot.
		h.Div(
			h.Label(HoneypotFormKey).Attr("for", HoneypotFormKey),
			h.Input("").Type("text").Attr("name", HoneypotFormKey).
				Attr("id", HoneypotFormKey).
				Attr("tabindex", "-1").
				Attr("autocomplete", "off"),
		).Attr("aria-hidden", "true").
			Attr("style", "position:absolute;left:-9999px;top:-9999px;height:0;width:0;overflow:hidden"),
	)
}
