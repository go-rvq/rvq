package login

import (
	"net/http/httptest"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/x/i18n"
)

func newChallengeBuilder() *Builder {
	b := New(i18n.New())
	b.Secret("a secret worth at least this many bytes")
	// the form is filled instantly in a test
	b.ChallengeAges(time.Nanosecond, time.Minute)
	return b
}

// formRequest runs the check over the request a browser would send.
func formRequest(b *Builder, values url.Values) (code FailCode) {
	r := httptest.NewRequest("POST", "/login", strings.NewReader(values.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return b.VerifyChallenge(r)
}

func TestChallengeAcceptsWhatTheFormCarries(t *testing.T) {
	b := newChallengeBuilder()

	if code := formRequest(b, url.Values{ChallengeFormKey: {b.NewChallenge()}}); code != 0 {
		t.Errorf("a form rendered by this server was refused: fail code %d", code)
	}
}

func TestChallengeRefusesAPostThatCameFromNowhere(t *testing.T) {
	b := newChallengeBuilder()

	cases := []struct {
		name   string
		values url.Values
	}{
		{name: "no challenge at all", values: url.Values{}},
		{name: "made up", values: url.Values{ChallengeFormKey: {"nonce.123.signature"}}},
		{name: "no signature", values: url.Values{ChallengeFormKey: {"nonce.123"}}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if code := formRequest(b, c.values); code == 0 {
				t.Error("accepted")
			}
		})
	}
}

func TestChallengeRefusesAValueSignedElsewhere(t *testing.T) {
	other := newChallengeBuilder()
	other.Secret("another secret, long enough to be one")

	b := newChallengeBuilder()
	if code := formRequest(b, url.Values{ChallengeFormKey: {other.NewChallenge()}}); code == 0 {
		t.Error("a challenge signed with another secret was accepted")
	}
}

func TestChallengeIsSingleUse(t *testing.T) {
	b := newChallengeBuilder()
	value := b.NewChallenge()

	if code := formRequest(b, url.Values{ChallengeFormKey: {value}}); code != 0 {
		t.Fatalf("the first submit was refused: fail code %d", code)
	}
	// what a script does with a challenge captured from a rendered page
	if code := formRequest(b, url.Values{ChallengeFormKey: {value}}); code == 0 {
		t.Error("the same challenge was accepted twice")
	}
}

func TestChallengeRefusesTooFastAndTooLate(t *testing.T) {
	b := newChallengeBuilder()
	b.ChallengeAges(time.Hour, 2*time.Hour) // nothing is fast enough

	if code := formRequest(b, url.Values{ChallengeFormKey: {b.NewChallenge()}}); code != FailCodeIncorrectChallenge {
		t.Errorf("a submit faster than a person = %d, want FailCodeIncorrectChallenge", code)
	}

	b = newChallengeBuilder()
	b.ChallengeAges(time.Nanosecond, time.Millisecond)
	value := b.NewChallenge()
	time.Sleep(2 * time.Millisecond)

	if code := formRequest(b, url.Values{ChallengeFormKey: {value}}); code != FailCodeChallengeExpired {
		t.Errorf("a form left open = %d, want FailCodeChallengeExpired", code)
	}
}

func TestChallengeHoneypot(t *testing.T) {
	b := newChallengeBuilder()

	values := url.Values{
		ChallengeFormKey: {b.NewChallenge()},
		HoneypotFormKey:  {"someone@example.com"},
	}
	if code := formRequest(b, values); code == 0 {
		t.Error("a filled honeypot was accepted")
	}
}

func TestChallengeOffWhenRecaptchaIsConfigured(t *testing.T) {
	b := newChallengeBuilder()
	b.Recaptcha(true, RecaptchaConfig{SiteKey: "site", SecretKey: "secret"})

	if b.ChallengeEnabled() {
		t.Error("the two protections must not both be on")
	}
	// and then the check does not stand in the way of the reCAPTCHA one
	if code := formRequest(b, url.Values{}); code != 0 {
		t.Errorf("with reCAPTCHA on, the built-in check answered %d", code)
	}
}

func TestChallengeDisabled(t *testing.T) {
	b := newChallengeBuilder()
	b.ChallengeDisabled(true)

	if b.ChallengeEnabled() {
		t.Error("still enabled")
	}
	if code := formRequest(b, url.Values{}); code != 0 {
		t.Errorf("a disabled check answered %d", code)
	}
	if fields := (&ViewHelper{b: b}).ChallengeFields(); fields != nil {
		t.Error("a disabled check must not render fields")
	}
}

func TestChallengeFieldsAreWhatTheServerVerifies(t *testing.T) {
	b := newChallengeBuilder()
	vh := &ViewHelper{b: b}

	rendered := h.MustString(vh.ChallengeFields(), t.Context())

	// the challenge travels as a plain hidden field — no script, nothing that
	// needs a secure context, so it works over HTTP as well as HTTPS
	if !strings.Contains(rendered, `name='`+ChallengeFormKey+`'`) {
		t.Errorf("the challenge field is missing:\n%s", rendered)
	}
	if !strings.Contains(rendered, `name='`+HoneypotFormKey+`'`) {
		t.Errorf("the honeypot field is missing:\n%s", rendered)
	}
	if strings.Contains(rendered, "<script") {
		t.Errorf("the protection must not depend on JavaScript:\n%s", rendered)
	}

	// what was rendered passes the check
	m := regexp.MustCompile(`name='` + ChallengeFormKey + `'[^>]*value='([^']+)'`).FindStringSubmatch(rendered)
	if m == nil {
		t.Fatalf("could not read the rendered challenge:\n%s", rendered)
	}
	value := m[1]
	if code := formRequest(b, url.Values{ChallengeFormKey: {value}}); code != 0 {
		t.Errorf("the rendered challenge was refused: fail code %d", code)
	}
}

func TestNewChallengeIsUniqueAndDated(t *testing.T) {
	b := newChallengeBuilder()

	first, second := b.NewChallenge(), b.NewChallenge()
	if first == second {
		t.Error("two challenges came out the same")
	}

	payload := first[:strings.LastIndexByte(first, '.')]
	issued, err := strconv.ParseInt(payload[strings.LastIndexByte(payload, '.')+1:], 10, 64)
	if err != nil {
		t.Fatalf("the issuing instant is not readable: %v", err)
	}
	if age := time.Since(time.Unix(0, issued)); age < 0 || age > time.Minute {
		t.Errorf("issued at %v, which is %v ago", time.Unix(0, issued), age)
	}
}
