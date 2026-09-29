package presets

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql/driver"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/osenv"
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
// with a model field; the name says what the value IS — a signed form value —
// and not which field it came from.
const RecordStampFormKey = "__formSign"

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

	// FormSecret is the signing key, from the environment. Falls back to
	// DefaultFormSecret, which is public — see the warning in NewHMACFormSigner.
	FormSecret = osenv.Get("RVQ_FORM_SECRET",
		"Key that signs the form values the user must not change (at least 32 bytes). Falls back to a PUBLIC default, which no deployment should keep",
		DefaultFormSecret)
)

const (
	// FormSecretMinLen is the shortest key accepted: HMAC takes any length, but
	// a secret worth the name does not fit in less.
	FormSecretMinLen = 32

	// DefaultFormSecret keeps a development server working out of the box —
	// forms survive a restart, and instances accept each other's. It is written
	// right here, so ANYBODY can forge a signed form value with it: using it
	// logs a warning, and a deployment must set RVQ_FORM_SECRET.
	DefaultFormSecret = "rvq-insecure-default-form-signing-key"
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

// NewHMACFormSigner takes the signing key, falling back to FormSecret (which
// itself falls back to the public DefaultFormSecret, with a warning).
//
// Panics on a secret shorter than FormSecretMinLen: signing with a weak key
// looks like it works, and the failure only shows up the day somebody forges a
// stamp.
func NewHMACFormSigner(key []byte) *HMACFormSigner {
	if len(key) == 0 {
		if l := len(FormSecret); l < FormSecretMinLen {
			panic(fmt.Sprintf("RVQ_FORM_SECRET must have at least %d bytes, got %d",
				FormSecretMinLen, l))
		}

		if FormSecret == DefaultFormSecret {
			warnDefaultFormSecret.Do(func() {
				log.Printf("WARNING: RVQ_FORM_SECRET is not set — forms are signed with the "+
					"DEFAULT key, which is public in the source: anybody can forge a signed "+
					"form value. Set it (at least %d bytes) before this reaches anyone else.",
					FormSecretMinLen)
			})
		}

		key = []byte(FormSecret)
	}
	return &HMACFormSigner{key: key}
}

// warned once per process: the signer is built for every presets.Builder
var warnDefaultFormSecret sync.Once

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

// RecordStateHash is the stamp of a record that has no UpdatedAt to be stamped
// by: a hash of the record as THIS FORM sees it, so a change made after the
// form was rendered is still noticed. It is what gives the optimistic lock
// something to compare on a model with no timestamps at all.
//
// What goes in are the fields the form edits — the editing FieldsBuilder, in
// order — and nothing else. A field the form does not carry cannot be part of
// what the form is locking:
//
//   - a column (a scalar, a []byte, a time.Time, or anything implementing
//     driver.Valuer, such as a JSON setting) enters as its stored value;
//   - a NESTED field enters through its own fields, element by element for a
//     list — the same walk, one level down;
//   - a RELATED record enters by its ID alone. The form carries the choice, not
//     the other record's contents, and whether it was loaded at all depends on
//     the fetcher.
func (b *EditingBuilder) RecordStateHash(obj any) string {
	h := sha256.New()
	b.WalkRecordState(obj, func(key, value string) {
		io.WriteString(h, "\x00"+key+"="+value)
	})
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

// WalkRecordState calls visit for every form key the stamp covers, with the
// value it contributes. RecordStateHash is this walk, hashed.
//
// It is what an application checks its own components against: every key the
// rendered form binds with `form["…"]` — bar the control fields — should come
// out of here, or the guard is comparing something the form does not carry.
func (b *EditingBuilder) WalkRecordState(obj any, visit func(key, value string)) {
	walkFieldsState(visit, &b.FieldsBuilder, reflect.Indirect(reflect.ValueOf(obj)), "", 0)
}

// RecordStampKeys are the form keys the stamp covers, in order.
func (b *EditingBuilder) RecordStampKeys(obj any) (keys []string) {
	b.WalkRecordState(obj, func(key, _ string) {
		keys = append(keys, key)
	})
	return
}

// maxRecordStateDepth stops the walk from following nested fields forever.
const maxRecordStateDepth = 8

func walkFieldsState(visit func(key, value string), fb *FieldsBuilder, v reflect.Value, prefix string, depth int) {
	if !v.IsValid() || v.Kind() != reflect.Struct || depth > maxRecordStateDepth {
		return
	}

	for _, f := range fb.fields {
		// the value comes through the field's own context, the same object the
		// components read — and through RawValue, which is what the form POSTS
		// BACK. The display value is not it: FieldContext.ValueOverride masks a
		// password as "***" and prints a month as its label, and hashing that
		// would hide a change instead of catching one.
		key := prefix + f.name
		fv, ok := fieldStampValue(v, f, key)
		if !ok {
			continue
		}

		if f.nested != nil {
			nested := f.nested.FieldsBuilder()
			if nested == nil {
				continue
			}
			nv := reflect.Indirect(fv)
			switch nv.Kind() {
			case reflect.Slice, reflect.Array:
				visit(key+"#", strconv.Itoa(nv.Len()))
				for i := 0; i < nv.Len(); i++ {
					walkFieldsState(visit, nested, reflect.Indirect(nv.Index(i)),
						key+"["+strconv.Itoa(i)+"].", depth+1)
				}
			default:
				walkFieldsState(visit, nested, nv, key+".", depth+1)
			}
			continue
		}

		if s, ok := columnValue(fv); ok {
			visit(key, s)
			continue
		}

		// what is left is another record: the form carries the choice, so the
		// identity is the whole of it
		if s, ok := relatedID(fv); ok {
			visit(key+"#", s)
		}
	}
}

// fieldStampValue is what this field contributes to the stamp: its RawValue,
// read through the field's own FieldContext — the same object every component
// reads the record through, so the hash and the form cannot drift apart.
//
// A name the record does not answer to is a field the form builds itself; it
// carries no record state and stays out.
func fieldStampValue(v reflect.Value, f *FieldBuilder, formKey string) (fv reflect.Value, ok bool) {
	if !v.IsValid() {
		return reflect.Value{}, false
	}
	if !v.CanAddr() {
		// RawValue needs something it can address
		p := reflect.New(v.Type())
		p.Elem().Set(v)
		v = p.Elem()
	}

	defer func() {
		if recover() != nil {
			fv, ok = reflect.Value{}, false
		}
	}()

	fc := &FieldContext{
		ToComponentOptions: &ToComponentOptions{},
		Obj:                v.Addr().Interface(),
		Field:              f,
		Name:               f.name,
		FormKey:            formKey,
		Path:               FieldPath{f.name},
	}

	value := fc.RawValue()
	if value == nil {
		return reflect.Value{}, false
	}
	return reflect.ValueOf(value), true
}

// relatedID is the identity of a related record (or of each one, in order),
// which is all of it the form carries.
func relatedID(v reflect.Value) (string, bool) {
	v = reflect.Indirect(v)
	if !v.IsValid() {
		return "nil", true
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		var ids string
		for i := 0; i < v.Len(); i++ {
			s, _ := relatedID(v.Index(i))
			ids += s + ","
		}
		return ids, true
	case reflect.Struct:
		id := v.FieldByName("ID")
		if !id.IsValid() {
			return "", false
		}
		s, ok := columnValue(id)
		return s, ok
	}

	return "", false
}

// columnValue is the field as the database stores it, and whether it is stored
// at all.
func columnValue(v reflect.Value) (string, bool) {
	if !v.IsValid() {
		return "", false
	}

	// a Valuer says for itself what it is worth, on the value or on a pointer
	// to it — which is how gorm reads it back out
	if s, ok := valuerString(v); ok {
		return s, true
	}
	if v.Kind() != reflect.Ptr && v.CanAddr() {
		if s, ok := valuerString(v.Addr()); ok {
			return s, true
		}
	}

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "nil", true
		}
		return columnValue(v.Elem())
	}

	if t, ok := v.Interface().(time.Time); ok {
		return RecordStampValue(t), true
	}

	switch v.Kind() {
	case reflect.Bool, reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return fmt.Sprint(v.Interface()), true
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			if v.IsNil() {
				return "nil", true
			}
			return string(v.Bytes()), true
		}
	}

	return "", false
}

// valuerString asks a driver.Valuer for its stored value. A Valuer that errors
// is left out: it has no value to compare, and the update is not the place to
// report it.
func valuerString(v reflect.Value) (string, bool) {
	if !v.CanInterface() {
		return "", false
	}
	valuer, ok := v.Interface().(driver.Valuer)
	if !ok {
		return "", false
	}
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return "nil", true
	}

	dv, err := valuer.Value()
	if err != nil {
		return "", false
	}
	if dv == nil {
		return "nil", true
	}
	if b, ok := dv.([]byte); ok {
		return string(b), true
	}
	if t, ok := dv.(time.Time); ok {
		return RecordStampValue(t), true
	}
	return fmt.Sprint(dv), true
}

// recordStamp is the value the form carries for obj, and whether the record is
// guarded at all: the instant its UpdatedAt holds, or — for a model that has
// none — the hash of the fields this form edits, unless the model turned that
// fallback off (see ModelBuilder.SetRecordStateStamp).
func (b *EditingBuilder) recordStamp(obj any) (string, bool) {
	if t, ok := RecordUpdatedAt(obj); ok {
		return RecordStampValue(t), true
	}
	if b.mb.noRecordStateStamp {
		return "", false
	}
	return b.RecordStateHash(obj), true
}

// recordStampField is the hidden field the edit form carries: the record as it
// was when the form was rendered, signed.
func (b *EditingBuilder) recordStampField(obj any) h.HTMLComponent {
	value, ok := b.recordStamp(obj)
	if !ok {
		return nil
	}
	signed := b.mb.p.FormSigner().Sign(value)
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
	// Trusted in-process code may opt out of the optimistic-lock stamp (see
	// web.WithSkipFormSign) when it submits a form without the rendered signed
	// stamp. The opt-in lives only in the request context, so a network request
	// can never set it.
	if web.SkipFormSign(ctx.R) {
		return nil
	}

	stored := b.mb.NewModel()
	// the TYPE answers whether the model is guarded, and by what — an
	// unguarded model costs no read at all
	if _, timed := RecordUpdatedAt(stored); !timed && b.mb.noRecordStateStamp {
		return nil
	}

	if err := b.Fetcher(stored, mid, ctx); err != nil {
		return err
	}

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

	if current, _ := b.recordStamp(stored); sent == current {
		return nil
	}

	at, timed := RecordUpdatedAt(stored)
	if !timed {
		// the hash says the record moved, and nothing says when or by whom
		return &recordStampError{cause: ErrRecordChanged, msg: string(msgr.ErrRecordChangedUnknownWhen)}
	}
	return &recordStampError{cause: ErrRecordChanged, msg: b.recordChangedMessage(stored, at, ctx)}
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

// SignForm returns r carrying the record stamp of obj, signed exactly as the
// rendered edit form carries it (under RecordStampFormKey).
//
// It is for a caller that builds an update request BY HAND — a unit test, or
// in-process code submitting a form nobody rendered. A guarded model REQUIRES
// the stamp and the stamp is signed, so such a request is refused without it
// (ErrRecordStampMissing); this puts in the value the form would have had.
//
// obj is the record as STORED — the one the update is about — because its
// UpdatedAt is the stamp. Read it back from the database right before signing:
// signing the object the test is about to send says the record has not changed,
// which is what the guard is there to decide. A model with no UpdatedAt is not
// guarded, and r comes back unchanged.
//
// r is left alone; the copy carries the new body. multipart/form-data and
// application/x-www-form-urlencoded bodies are supported — anything else is an
// error, since there is no form to add a field to.
//
//	var stored Post
//	db.First(&stored, id)
//	req, err := pb.SignForm(req, &stored)
func (b *Builder) SignForm(r *http.Request, obj any) (*http.Request, error) {
	mb := b.GetModel(obj)
	if mb == nil {
		return nil, fmt.Errorf("presets: no model registered for %T", obj)
	}

	value, ok := mb.Editing().recordStamp(obj)
	if !ok {
		// not guarded: there is no stamp to carry
		return r, nil
	}
	return addFormField(r, RecordStampFormKey, b.FormSigner().Sign(value))
}

// addFormField returns a copy of r whose form carries one more field, re-encoding
// the body it came with.
func addFormField(r *http.Request, key, value string) (*http.Request, error) {
	if r.Body == nil {
		return nil, fmt.Errorf("presets: cannot add %s: the request has no body", key)
	}

	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("presets: reading the request body: %w", err)
	}
	// the caller's request stays readable
	r.Body = io.NopCloser(bytes.NewReader(body))

	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("presets: parsing the request content type: %w", err)
	}

	var (
		out         bytes.Buffer
		contentType string
	)

	switch {
	case strings.HasPrefix(mediaType, "multipart/"):
		w := multipart.NewWriter(&out)
		mr := multipart.NewReader(bytes.NewReader(body), params["boundary"])
		for {
			// the RAW part: re-encoding must not decode what it copies.
			p, err := mr.NextRawPart()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("presets: reading a form part: %w", err)
			}
			// CreatePart keeps the part's own headers, so a file part keeps its
			// filename and its content type.
			pw, err := w.CreatePart(p.Header)
			if err != nil {
				return nil, err
			}
			if _, err = io.Copy(pw, p); err != nil {
				return nil, fmt.Errorf("presets: copying a form part: %w", err)
			}
		}
		if err = w.WriteField(key, value); err != nil {
			return nil, err
		}
		if err = w.Close(); err != nil {
			return nil, err
		}
		contentType = w.FormDataContentType()

	case mediaType == "application/x-www-form-urlencoded":
		values, err := url.ParseQuery(string(body))
		if err != nil {
			return nil, fmt.Errorf("presets: parsing the form body: %w", err)
		}
		values.Set(key, value)
		out.WriteString(values.Encode())
		contentType = mediaType

	default:
		return nil, fmt.Errorf("presets: cannot add %s to a %s body", key, mediaType)
	}

	signed := r.Clone(r.Context())
	newBody := out.Bytes()
	signed.Body = io.NopCloser(bytes.NewReader(newBody))
	signed.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(newBody)), nil
	}
	signed.ContentLength = int64(len(newBody))
	signed.Header.Set("Content-Type", contentType)
	// a form parsed from the old body must not travel to the copy
	signed.Form, signed.PostForm, signed.MultipartForm = nil, nil, nil
	return signed, nil
}
