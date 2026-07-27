package presets

import (
	"bytes"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestHMACFormSigner(t *testing.T) {
	s := NewHMACFormSigner([]byte("a key"))

	signed := s.Sign("1234567890")
	if signed == "1234567890" {
		t.Fatal("the signature is missing from the signed value")
	}

	got, err := s.Unsign(signed)
	if err != nil || got != "1234567890" {
		t.Fatalf("round trip = %q, %v; want %q, nil", got, err, "1234567890")
	}

	// what the guard exists for: a value edited by hand in the POST
	tampered := "9999999999" + signed[strings.IndexByte(signed, '.'):]
	if _, err = s.Unsign(tampered); !errors.Is(err, ErrInvalidFormSignature) {
		t.Errorf("tampered value: err = %v, want ErrInvalidFormSignature", err)
	}

	if _, err = s.Unsign("1234567890"); !errors.Is(err, ErrInvalidFormSignature) {
		t.Errorf("value with no signature: err = %v, want ErrInvalidFormSignature", err)
	}

	// another instance, another key
	if _, err = NewHMACFormSigner([]byte("another key")).Unsign(signed); !errors.Is(err, ErrInvalidFormSignature) {
		t.Errorf("another key: err = %v, want ErrInvalidFormSignature", err)
	}

	// no key given: a random one, still self-consistent
	rnd := NewHMACFormSigner(nil)
	if got, err = rnd.Unsign(rnd.Sign("42")); err != nil || got != "42" {
		t.Errorf("random key round trip = %q, %v", got, err)
	}
}

func TestFormSecret(t *testing.T) {
	original := FormSecret
	defer func() { FormSecret = original }()

	// a secret from the environment is the key: two signers accept each other's
	// values, which is what makes forms survive a restart
	FormSecret = strings.Repeat("s", FormSecretMinLen)
	signed := NewHMACFormSigner(nil).Sign("42")
	if got, err := NewHMACFormSigner(nil).Unsign(signed); err != nil || got != "42" {
		t.Errorf("two signers built from the same secret disagree: %q, %v", got, err)
	}

	// the default keeps a development server working out of the box, and says
	// so out loud
	FormSecret = DefaultFormSecret
	warnDefaultFormSecret = sync.Once{}

	var warning bytes.Buffer
	log.SetOutput(&warning)
	defer log.SetOutput(os.Stderr)

	signed = NewHMACFormSigner(nil).Sign("42")
	if got, err := NewHMACFormSigner(nil).Unsign(signed); err != nil || got != "42" {
		t.Errorf("the default secret does not sign: %q, %v", got, err)
	}
	if !strings.Contains(warning.String(), "RVQ_FORM_SECRET") {
		t.Errorf("using the default key must warn, got %q", warning.String())
	}

	// too short (an empty secret included): a startup panic, not a weak
	// signature nobody notices
	FormSecret = strings.Repeat("s", FormSecretMinLen-1)
	func() {
		defer func() {
			if recover() == nil {
				t.Error("a secret shorter than the minimum was accepted")
			}
		}()
		NewHMACFormSigner(nil)
	}()

	// an explicit key wins over the environment, whatever its length
	FormSecret = "short"
	if NewHMACFormSigner([]byte("my own key")) == nil {
		t.Error("an explicit key must be taken as it is")
	}
}

type stampedModel struct {
	ID          uint
	UpdatedAt   time.Time
	UpdatedByID uint
}

type stampedPtrModel struct {
	ID        uint
	UpdatedAt *time.Time
}

type unstampedModel struct {
	ID   uint
	Name string
}

type embeddedTimestamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type stampedEmbeddedModel struct {
	embeddedTimestamps
	ID uint
}

func TestRecordUpdatedAt(t *testing.T) {
	now := time.Date(2026, 7, 27, 10, 30, 0, 0, time.UTC)

	cases := []struct {
		name   string
		obj    any
		wantOk bool
		want   time.Time
	}{
		{name: "no such field", obj: &unstampedModel{}},
		{name: "zero", obj: &stampedModel{}, wantOk: true},
		{name: "set", obj: &stampedModel{UpdatedAt: now}, wantOk: true, want: now},
		{name: "nil pointer counts as present and zero", obj: &stampedPtrModel{}, wantOk: true},
		{name: "pointer", obj: &stampedPtrModel{UpdatedAt: &now}, wantOk: true, want: now},
		{name: "embedded (gorm.Model style)", obj: &stampedEmbeddedModel{embeddedTimestamps: embeddedTimestamps{UpdatedAt: now}}, wantOk: true, want: now},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, ok := RecordUpdatedAt(c.obj)
			if ok != c.wantOk {
				t.Fatalf("ok = %v, want %v", ok, c.wantOk)
			}
			if ok && !got.Equal(c.want) {
				t.Errorf("t = %v, want %v", got, c.want)
			}
		})
	}
}

func TestRecordUpdatedByID(t *testing.T) {
	if _, ok := RecordUpdatedByID(&unstampedModel{}); ok {
		t.Error("a model without UpdatedByID must report none")
	}
	if _, ok := RecordUpdatedByID(&stampedModel{}); ok {
		t.Error("a zero UpdatedByID records nobody")
	}
	id, ok := RecordUpdatedByID(&stampedModel{UpdatedByID: 7})
	if !ok || id != uint(7) {
		t.Errorf("id = %v, %v; want 7, true", id, ok)
	}
}

func TestRecordChangedMessage(t *testing.T) {
	at := time.Date(2026, 7, 27, 10, 45, 12, 0, time.UTC)
	when := Messages_en_US.FormatDateTime(at)

	anonymous := Messages_en_US.RecordChangedMessage("", "", at)
	if !strings.Contains(anonymous, when) {
		t.Errorf("the anonymous message does not say when: %q", anonymous)
	}

	named := Messages_en_US.RecordChangedMessage("Ana Souza", "ana@example.com", at)
	for _, want := range []string{"Ana Souza", "ana@example.com", when} {
		if !strings.Contains(named, want) {
			t.Errorf("the message does not carry %q: %q", want, named)
		}
	}

	// messages that define no layout must not silently print an empty instant
	// (time.Format("") returns "")
	empty := &Messages{}
	if got := empty.FormatDateTime(at); got == "" {
		t.Error("with no layout configured the instant came out empty")
	}
}

func TestRecordStampValue(t *testing.T) {
	if got := RecordStampValue(time.Time{}); got != "0" {
		t.Errorf("zero time = %q, want %q", got, "0")
	}

	now := time.Date(2026, 7, 27, 10, 30, 0, 123456789, time.UTC)
	if got, want := RecordStampValue(now), strconv.FormatInt(now.UnixNano(), 10); got != want {
		t.Errorf("stamp = %q, want %q", got, want)
	}

	// the same instant in another zone is the same record
	other := now.In(time.FixedZone("BRT", -3*3600))
	if RecordStampValue(other) != RecordStampValue(now) {
		t.Error("the stamp must not depend on the location")
	}
}
