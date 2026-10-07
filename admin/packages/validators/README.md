# validators

A reusable rvq package for **named validation rules whose algorithm is a script**
written in [GAD](https://github.com/gad-lang/gad), with **per-language error
messages**. Records are seeded by the application; end users may only tune the
algorithm, its documentation and the messages — and the detail view renders the
documentation (Markdown) that explains the algorithm.

```
github.com/go-rvq/rvq/admin/packages/validators
├── models/   the Validator model, the GAD engine and the store helpers
├── admin/    the presets CRUD (Markdown detail, Messages editor, field validator)
└── messages/ the admin i18n (pt-BR default, en-US)
```

## The idea

A `Validator` couples:

- a **GAD algorithm** (`Value`, defaulting to `InitialValue`) that analyzes a
  string value and returns whether it is valid;
- **Markdown documentation** (`Doc` / `InitialDoc`) explaining the algorithm,
  rendered in the detail view;
- **translated error messages** (`Messages` / `InitialMessages`), a
  `{LANG: {KEY: VALUE}}` map, with a `DefaultLang` fallback.

The application seeds validators (their `Initial*`); users override `Value`,
`Doc` and `Messages`. **Clearing an override and saving resets it to the
registered default** (an empty override falls back to `Initial*`).

## The GAD contract

The script receives the value to analyze as its `value` parameter and:

- returns a truthy value (typically `return true`) when the value is **valid**;
- **throws a message key** — `throw "MESSAGE_KEY"` — to signal a failure with a
  translatable reason. The engine translates the key via the validator's
  `Messages` for the user's language (falling back to `DefaultLang`, then to the
  key itself).

A plain falsy return (no throw) is a generic failure (message key `invalid`).
Any other compile/runtime error is a real script error.

```gad
param value

// keep only the digits
digits := []

for i := 0; i < len(value); i++ {
    c := value[i]

    if c >= '0' && c <= '9' {
        digits = append(digits, int(c) - int('0'))
    }
}

if len(digits) != 11 {
    throw "cpf" // translated via Messages["<lang>"]["cpf"]
}

// … check digits …
return true
```

## Using validators from Go

```go
import validators "github.com/go-rvq/rvq/admin/packages/validators/models"

// look up + run, getting a translated error when invalid:
if err := validators.Check(db, "cpf", document, "pt-BR"); err != nil {
    // err.Error() == "CPF inválido."
}

// or lower level:
v, _ := validators.Get(db, "cpf")
ok, msgKey, err := v.Eval(document) // ok=false, msgKey="cpf" when invalid
```

## As an admin field validator (`FieldBuilder.Validators`)

Validate a form field through a registered validator, so the error is translated
and the algorithm is editable by users:

```go
ed.Field("Document").Validator(presets.FieldValidatorFunc(
    func(field *presets.FieldContext) (verr web.ValidationErrors) {
        doc, _ := field.Value().(string)
        if doc == "" { return }
        if err := validators.Check(db, "cpf", doc, lang); err != nil {
            verr.FieldError(field.Name, err.Error())
        }
        return
    }))
```

The **people** package does exactly this for CPF/CNPJ (see
[people](../people)); the **finance** module reuses it for holder/payer
documents.

## Mounting (the plugin builder)

`admin.Builder` implements `presets.Plugin`; apply it with `pb.Use(...)`. Use it
to seed validators and expose the CRUD:

```go
import validatorsadmin "github.com/go-rvq/rvq/admin/packages/validators/admin"

pb.Use(validatorsadmin.New(db).
    Register(myValidator1, myValidator2)) // seeds them (upsert), exposes the CRUD

// or the convenience:
validatorsadmin.Configure(pb, db)
```

`ExposeCRUD(false)` seeds without registering the admin resource.

The **people** package's `admin.Builder` is likewise a `presets.Plugin` that
registers its CPF/CNPJ document validators and lets you add more per document
type — see [people](../people).

## Admin

The `admin.Builder` registers the `validators` resource:

- **create/delete are disabled** — records are seeded by the application;
- only **Value**, **Doc** and **Messages** are editable;
- **Messages** is edited as a `{LANG: {KEY: VALUE}}` map, with a field validator
  enforcing the shape;
- the **detail view** renders the effective algorithm and the documentation
  (Markdown).

## Performance: the Bytecode cache

Each `Validator` has a `Bytecode` column (not shown in the UI). On save
(`BeforeSave`) the effective algorithm is compiled and its bytecode cached, so
validation runs the cached bytecode instead of recompiling every call. A stale
cache falls back to compiling the source.

## Developing a new validator

See [docs/developing.md](docs/developing.md) for the full guide: writing and
testing the GAD algorithm, choosing message keys, seeding the registration and
wiring it as a field validator.
