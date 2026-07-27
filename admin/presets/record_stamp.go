package presets

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
)

// Optimistic locking for edit forms.
//
// Two people open the same record; the first saves; the second saves minutes
// later and silently throws the first one's work away — with no error anywhere,
// because each update is valid on its own. The guard against that is the
// record's own `UpdatedAt`: the edit form carries the value the record had when
// it was RENDERED, and the update refuses to run if the record has changed since.
//
// The value travels in a hidden field, SIGNED (see FormSigner): it decides
// whether a save is allowed, so a hand-edited POST must not be able to set it.
//
// Models without an `UpdatedAt` field are not guarded — there is nothing to
// compare — and the form carries no stamp for them.

// RecordStampFormKey is the form key of the hidden stamp. The `__` prefix marks
// a control field (like the list editor's `__Deleted.…`), so it never collides
// with a model field.
const RecordStampFormKey = "__UpdatedAt"

var (
	// ErrInvalidFormSignature is returned when a signed form value does not
	// match its signature — the value was changed after the server wrote it.
	ErrInvalidFormSignature = errors.New("invalid form signature")

	// ErrRecordStampMissing is returned when a guarded model is updated without
	// the stamp. It is required: accepting the update without it would let any
	// caller skip the check by dropping the field.
	ErrRecordStampMissing = errors.New("record stamp is missing from the form")

	// ErrRecordChanged is returned when the stored record changed after the form
	// was rendered. The error's own message is the localized one shown to the
	// user (and names the author when it can); this sentinel is what
	// errors.Is matches.
	ErrRecordChanged = errors.New("record changed after the form was rendered")
)

// FormSigner signs the values a form carries that the user must NOT be able to
// change — they are read back as server decisions, not as user input.
type FormSigner interface {
	// Sign returns value plus its signature.
	Sign(value string) (signed string)
	// Unsign returns the value carried by signed, or ErrInvalidFormSignature.
	Unsign(signed string) (value string, err error)
}

// HMACFormSigner signs with HMAC-SHA256. The signature travels next to the
// value (`<value>.<signature>`), so nothing has to be kept server side.
type HMACFormSigner struct {
	key []byte
}

// NewHMACFormSigner takes the signing key. A random key is generated when key
// is empty — good enough for a single process, but forms rendered before a
// restart stop being accepted, and several instances will not accept each
// other's forms, so a real deployment should pass its own key.
func NewHMACFormSigner(key []byte) *HMACFormSigner {
	if len(key) == 0 {
		key = make([]byte, 32)
		if _, err := rand.Read(key); err != nil {
			panic(err)
		}
	}
	return &HMACFormSigner{key: key}
}

func (s *HMACFormSigner) sum(value string) string {
	m := hmac.New(sha256.New, s.key)
	m.Write([]byte(value))
	return base64.RawURLEncoding.EncodeToString(m.Sum(nil))
}

func (s *HMACFormSigner) Sign(value string) string {
	return value + "." + s.sum(value)
}

func (s *HMACFormSigner) Unsign(signed string) (value string, err error) {
	i := strings.LastIndexByte(signed, '.')
	if i < 0 {
		return "", ErrInvalidFormSignature
	}

	value, sig := signed[:i], signed[i+1:]

	// constant time: a comparison that stops at the first difference tells an
	// attacker how much of a guessed signature is right
	if !hmac.Equal([]byte(sig), []byte(s.sum(value))) {
		return "", ErrInvalidFormSignature
	}
	return value, nil
}

// RecordUser is what the stale-record message needs from whoever changed the
// record last. The application's user model usually satisfies it already.
type RecordUser interface {
	GetName() string
	GetEmail() string
}

// RecordUserFinder loads the user behind an `UpdatedByID`. Register it with
// Builder.SetRecordUserFinder — presets cannot know the application's user
// model. Without one the stale-record message simply does not name anybody.
type RecordUserFinder func(id any, ctx *web.EventContext) (RecordUser, error)

// RecordUpdatedAt reads the model's `UpdatedAt` and reports whether it has one.
// A nil `*time.Time` counts as present and zero: the field exists, the record
// was never updated.
func RecordUpdatedAt(obj any) (t time.Time, ok bool) {
	// the TYPE is what answers "does this model have an UpdatedAt": a nil
	// *time.Time has no value to read, and the field is there all the same.
	// FieldByName also finds it when it comes from an embedded struct
	// (gorm.Model and friends).
	fv := recordField(obj, "UpdatedAt")
	if !fv.IsValid() {
		return
	}

	switch value := fv.Interface().(type) {
	case time.Time:
		return value, true
	case *time.Time:
		if value == nil {
			return time.Time{}, true
		}
		return *value, true
	}
	return
}

// recordField is the model's field by name, invalid when there is none.
func recordField(obj any, name string) reflect.Value {
	rv := reflect.Indirect(reflect.ValueOf(obj))
	if !rv.IsValid() || rv.Kind() != reflect.Struct {
		return reflect.Value{}
	}

	fv := rv.FieldByName(name)
	if !fv.IsValid() || !fv.CanInterface() {
		return reflect.Value{}
	}
	return fv
}

// RecordUpdatedByID reads the model's `UpdatedByID`, if it has one and it is
// set (a zero id means nobody is recorded).
func RecordUpdatedByID(obj any) (id any, ok bool) {
	fv := recordField(obj, "UpdatedByID")
	if !fv.IsValid() {
		return
	}

	if fv.Kind() == reflect.Ptr {
		if fv.IsNil() {
			return
		}
		fv = fv.Elem()
	}
	if fv.IsZero() {
		return
	}
	return fv.Interface(), true
}

// RecordStampValue is the stamp of an instant: nanoseconds since the epoch, or
// "0" for the zero time. An integer — the comparison has to be exact, and the
// value is signed as text.
func RecordStampValue(t time.Time) string {
	if t.IsZero() {
		return "0"
	}
	return strconv.FormatInt(t.UTC().UnixNano(), 10)
}

// recordStampField is the hidden field the edit form carries, or nil when the
// model has no `UpdatedAt` to guard.
func (b *EditingBuilder) recordStampField(obj any) h.HTMLComponent {
	t, ok := RecordUpdatedAt(obj)
	if !ok {
		return nil
	}

	signed := b.mb.p.FormSigner().Sign(RecordStampValue(t))
	return h.Input("").Type("hidden").Attr(web.VField(RecordStampFormKey, signed)...)
}

// VerifyRecordStamp compares the stamp the form came back with against the
// record as the database holds it NOW.
//
// It reads the stored record ITSELF instead of taking the one FetchAndUnmarshal
// returned: `UpdatedAt` and `UpdatedByID` are model fields like any other, so an
// application may well put them in the form — and then the unmarshalled object
// says what the FORM sent, not what is stored. Models with no `UpdatedAt` cost
// nothing: that is decided on the type, before any fetch.
//
// Returns nil when the model is not guarded and when the record is untouched
// (both zero, or both the same instant); otherwise an error whose message is
// meant for the user, wrapping:
//
//	ErrRecordStampMissing    the guarded update came without the stamp
//	ErrInvalidFormSignature  the stamp was changed after the server wrote it
//	ErrRecordChanged         the record moved under the form
func (b *EditingBuilder) VerifyRecordStamp(mid ID, ctx *web.EventContext) error {
	stored := b.mb.NewModel()
	if _, ok := RecordUpdatedAt(stored); !ok {
		return nil
	}

	if err := b.Fetcher(stored, mid, ctx); err != nil {
		return err
	}

	current, _ := RecordUpdatedAt(stored)
	msgr := MustGetMessages(ctx.Context())

	signed := ctx.R.FormValue(RecordStampFormKey)
	if signed == "" {
		return &recordStampError{cause: ErrRecordStampMissing, msg: string(msgr.ErrRecordStampMissing)}
	}

	sent, err := b.mb.p.FormSigner().Unsign(signed)
	if err != nil {
		// Same message as a missing stamp: to a good-faith user (whose form was
		// rendered by an instance with another key, or before a restart) this IS
		// an out-of-date form. To someone who edited the value by hand, it says
		// nothing useful — which is the point.
		return &recordStampError{cause: err, msg: string(msgr.ErrRecordStampMissing)}
	}

	if sent == RecordStampValue(current) {
		return nil
	}

	return &recordStampError{cause: ErrRecordChanged, msg: b.recordChangedMessage(stored, current, ctx)}
}

// recordStampError shows the reader a localized message while keeping the
// reason machine-readable (errors.Is with the sentinels above).
type recordStampError struct {
	cause error
	msg   string
}

func (e *recordStampError) Error() string { return e.msg }
func (e *recordStampError) Unwrap() error { return e.cause }

// recordChangedMessage explains that the record moved under the form: WHEN it
// was changed (the stored UpdatedAt — the very value that did not match) and,
// when the model records it (`UpdatedByID`) and the application knows how to
// load a user (Builder.SetRecordUserFinder), by WHOM.
func (b *EditingBuilder) recordChangedMessage(stored any, at time.Time, ctx *web.EventContext) string {
	msgr := MustGetMessages(ctx.Context())

	user := b.recordUser(stored, ctx)
	if user == nil {
		return msgr.RecordChangedMessage("", "", at)
	}
	return msgr.RecordChangedMessage(user.GetName(), user.GetEmail(), at)
}

// recordUser is whoever the record's `UpdatedByID` points at, or nil when the
// model does not record it, nobody is recorded, the application registered no
// finder, or the lookup failed — the message then simply names no one.
func (b *EditingBuilder) recordUser(stored any, ctx *web.EventContext) RecordUser {
	id, ok := RecordUpdatedByID(stored)
	if !ok {
		return nil
	}

	finder := b.mb.p.RecordUserFinder()
	if finder == nil {
		return nil
	}

	user, err := finder(id, ctx)
	if err != nil {
		return nil
	}
	return user
}
