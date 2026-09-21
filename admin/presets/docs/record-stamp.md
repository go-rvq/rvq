# Record stamp — the edit form refuses to overwrite someone else's work

Two people open the same record. The first saves. The second saves ten minutes
later and silently throws the first one's work away: each update is valid on its
own, so nothing anywhere reports a problem.

The guard is the record's own `UpdatedAt`. The edit form carries the value the
record had **when it was rendered**, and the update refuses to run when the
stored record has moved since.

```
render   /articles/1/edit    →  form carries UpdatedAt = 10:00:00
someone else saves           →  stored UpdatedAt = 10:04:12
POST     presets_Update      →  10:00:00 ≠ 10:04:12 → refused, with a message
```

## What is guarded

Every model, by one of two stamps:

| | the stamp is | noticed change |
| --- | --- | --- |
| model with an `UpdatedAt` | the instant the record held when the form was rendered | any save over the record, whatever it touched |
| model without one | a hash of the fields THIS FORM edits | a change to something the form carries |

`UpdatedAt` is a `time.Time` or `*time.Time`, the model's own or promoted from
an embedded struct (`gorm.Model`). A model that has one is guarded by it, and
the hash never comes into play.

The check runs on the update, right after `FetchAndUnmarshal` and **before any
validation** — what the form says is only worth validating if it was written
over the record as it stands now.

## The state hash, for a model with no UpdatedAt

`EditingBuilder.RecordStateHash` hashes the record as **this form** sees it: the
editing `FieldsBuilder`, in order, and nothing else. A field the form does not
carry cannot be part of what the form is locking.

| in the form | in the hash |
| --- | --- |
| a column — scalar, `[]byte`, `time.Time`, or a `driver.Valuer` (a JSON setting, `gorm.DeletedAt`) | its stored value |
| a NESTED field | its own fields, element by element for a list, plus the list's length |
| a RELATED record | its ID alone — the form carries the choice, not the other record's contents |
| a column the form does not edit | nothing |

Each value is read through the field's own `FieldContext`, by `RawValue()` —
the same read the components do, and the same value the form posts back. That
is what keeps the hash and the form from drifting apart, and it is why a
component binds its own key and renders its own value (see
[fields](fields.md#the-convention-a-component-follows)).

The fallback is ON. A model turns it off — and stops carrying, and requiring,
the stamp — with:

```go
mb.SetRecordStateStamp(false)
```

A model WITH an `UpdatedAt` is guarded either way: the switch does not reach it.

## The field, and why it is signed

The form carries a hidden field, `__formSign` (`presets.RecordStampFormKey`).
The name says what the value IS — a signed form value — and not which model
field it came from; the `__` prefix marks a control field, like the list
editor's `__Deleted.…`, so it never collides with a field of the model:

```html
<input type='hidden' v-model='form["__formSign"]'
       v-assign='[form, {"__formSign": "1785159832202194000.bpblc9o8MiMy…"}]'>
```

The value decides whether a save is allowed, so it must not be user input: it is
`<nanoseconds>.<HMAC-SHA256>`, and the update rejects anything whose signature
does not match. Nothing is kept server side.

**It is required.** An update of a guarded model with no stamp is refused —
otherwise skipping the check would be as easy as dropping the field from the
POST.

```go
// presets.FormSigner: signs the form values the user must not be able to change
type FormSigner interface {
	Sign(value string) (signed string)
	Unsign(signed string) (value string, err error)
}
```

The key comes from **`RVQ_FORM_SECRET`** (at least 32 bytes). Without it forms
are signed with `presets.DefaultFormSecret`, so a development server works out of
the box — forms survive a restart and instances accept each other's — and every
process using it **logs a warning**:

```
WARNING: RVQ_FORM_SECRET is not set — forms are signed with the DEFAULT key,
which is public in the source: anybody can forge a signed form value. Set it
(at least 32 bytes) before this reaches anyone else.
```

A secret shorter than the minimum is a startup panic: signing with a weak key
looks like it works, and the failure only shows up the day somebody forges a
stamp.

The key can also be given in code, and the whole signer replaced:

```go
b.SetFormSigner(presets.NewHMACFormSigner(myKey))  // your own key
b.SetFormSigner(mySigner)                          // your own scheme
```

| | |
| --- | --- |
| `RVQ_FORM_SECRET` | the signing key; at least `presets.FormSecretMinLen` (32) bytes |
| `presets.FormSecret` | the same value, readable in Go |
| `presets.DefaultFormSecret` | what it falls back to — public, warned about, never for a deployment |
| `presets.NewHMACFormSigner(key)` | key → `FormSecret`, in that order |

Both "no stamp" and "bad signature" answer with the same message — to a
good-faith user (a restarted server, another instance) the form IS out of date;
to whoever edited the value by hand it says nothing useful.

## When, and by whom

The message always says WHEN the record was changed — the stored `UpdatedAt`,
the very value that did not match the stamp — written with the language's
`TimeFormats.DateTime` (`Messages.FormatDateTime`, which falls back to a
readable layout when a language defines none).

When the model also has an `UpdatedByID`, it also says WHO saved over the
form. presets cannot know the application's user model, so the application
supplies the lookup:

```go
b.SetRecordUserFinder(func(id any, ctx *web.EventContext) (presets.RecordUser, error) {
	var u models.User
	if err := db.First(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &u, nil // RecordUser: GetName() string, GetEmail() string
})
```

Without a finder (or when the lookup fails, or `UpdatedByID` is zero) the
message simply does not name anybody.

| Situation | Message (pt-BR) |
| --- | --- |
| changed, author unknown | *Este registro foi alterado por outra pessoa **em 27/07/2026 10:45:12 -03:00**, depois que você abriu este formulário. Recarregue-o e refaça suas alterações, para que nada do que foi salvo nesse meio tempo se perca.* |
| changed, author known | *Este registro foi alterado por **Ana Souza (ana@example.com)** em 27/07/2026 10:45:12 -03:00, depois que você abriu este formulário. …* |
| changed, model with no `UpdatedAt` | *Este registro foi alterado por outra pessoa depois que você abriu este formulário. Recarregue-o e refaça suas alterações, para que nada do que foi salvo nesse meio tempo se perca.* |
| stamp missing or forged | *Este formulário está desatualizado e não pode ser salvo. Recarregue-o e refaça suas alterações.* |

They are `Messages.ErrRecordChanged` (one `%s`: when), `ErrRecordChangedBy`
(three: name, e-mail, when), `ErrRecordChangedUnknownWhen` — the state hash says
the record moved, and nothing says when or by whom — and
`ErrRecordStampMissing`, English in `Messages_en_US`, so an application can
replace them.

## After a refusal

The form is re-rendered with the message. It then shows the STORED values and
carries a fresh stamp, so saving again is possible — deliberately, by someone who
has now seen what changed.

## Details worth knowing

- The stamp is **nanoseconds since the epoch**, or `"0"` for the zero time: an
  exact integer comparison, independent of location. A record that was never
  updated (zero) must still be zero when the update runs.
- The check reads the stored record **on its own**, it does not trust the object
  `FetchAndUnmarshal` returned: `UpdatedAt` and `UpdatedByID` are model fields
  like any other, so an application may put them in the form — and then that
  object says what the FORM sent. Deciding "is this model guarded" happens on the
  TYPE, so an unguarded model costs no extra read.
- A creation carries no stamp (there is no record yet).
- The stamp of a model with no `UpdatedAt` costs one read of the stored record
  on save — the same read the time stamp already did. A model that turned the
  fallback off costs none: that is decided on the TYPE, before any fetch.
- Errors wrap `ErrRecordChanged`, `ErrRecordStampMissing` and
  `ErrInvalidFormSignature` — `errors.Is` works, while the message the user reads
  is the localized one.

## Signing a request built by hand

A browser posts the stamp because the form it came from carried it. A test — or
in-process code submitting a form nobody rendered — has to put it there itself,
and `Builder.SignForm` is what does it:

```go
req := multipartestutils.NewMultipartBuilder().
	PageURL("/posts?__execute_event__=presets_Update&id=1").
	AddField("Title", "Updated Title").
	BuildEventFuncRequest()

var stored Post
db.First(&stored, 1) // the record AS STORED, read the way the FORM read it
req, err := pb.SignForm(req, &stored)
```

It returns a COPY of the request whose form carries `__formSign`, signed by the
builder's own signer — the one the update verifies against. The request you
passed is left alone and stays readable.

| | |
| --- | --- |
| the record | pass it as STORED, read right before signing: signing the object you are about to send says "nothing changed", which is the very thing the guard decides |
| a record with nested fields | read it through the model's own fetcher — the hash covers them, and a record read without them hashes differently |
| a model with the fallback off | nothing to sign; the request comes back unchanged |
| the model | found by the object's type, so it must be registered on the builder |
| the body | `multipart/form-data` and `application/x-www-form-urlencoded`; anything else is an error, since there is no form to add a field to |
| files | survive, with their filename and content type |

Trusted in-process code with no record to sign can skip the check altogether
with `web.WithSkipFormSign(ctx)`. It lives in the request CONTEXT, so a network
request can never ask for it. The `http_api` command exposes it as
`--skip-form-sign`, and per request in a spec file:

```json
{"uri": "/admin/things", "method": "POST", "skipFormSign": true}
```

It defaults to off — the stamp is required.

## Tests

- [`record_stamp_test.go`](../record_stamp_test.go) — the signer (round trip,
  hand-edited value, another key), the secret (two signers from the same secret
  agree, two random ones do not, a short secret panics, an explicit key wins),
  field detection (`time.Time`, `*time.Time` nil, embedded), the stamp value and
  the message (when, by whom, and the fallback date layout).
- [`tests/listeditor/conformance_test.go`](../tests/listeditor/conformance_test.go)
  — the convention, checked against a real render: it takes the edit form of a
  record with a nested list, reads every literal `form["…"]` key out of the
  rendered HTML, and requires each one (bar the control fields) to be a key the
  stamp covers. A component that binds something else fails it. `WalkRecordState`
  and `RecordStampKeys` are exported for an application to do the same over its
  own forms.
- [`record_state_hash_test.go`](../record_state_hash_test.go) — the state hash:
  the same record hashes the same; a column, a valuer column, a time, the
  related CHOICE and a nested row added, removed, edited or reordered all move
  it; a column outside the form and the related record's own fields do not; the
  guard refuses a save over a record that moved and accepts one over a record
  whose untouched columns changed; and a model with the fallback off carries no
  stamp and requires none.
- [`record_stamp_sign_test.go`](../record_stamp_sign_test.go) — `SignForm`: the
  form it was given arrives whole (fields, a repeated one, the file), the stamp
  it adds unsigns to the record's instant, the request it was given is not
  touched, an unguarded model comes back unchanged, a body that is not a form is
  refused — and, the point of it all, `VerifyRecordStamp` accepts what it signs,
  reports `ErrRecordStampMissing` without it and `ErrRecordChanged` for a stale
  record.
- [`integration/record_stamp_test.go`](../integration/record_stamp_test.go) —
  end to end over a real database: the rendered form carries the record's stamp,
  saving it back works, saving it again after somebody else saved is refused, an
  update with no stamp is refused, a hand-edited or foreign-key stamp is refused,
  the message names the author AND when the change happened, and a model with no
  `UpdatedAt` saves normally.

> A test that posts `presets_Update` by hand against a guarded model has to send
> the stamp, the way a browser does — `pb.SignForm(req, &stored)`, above.
