# Developing a new validator

A validator is a **registration** (seeded by the application) of a GAD algorithm
plus its documentation and translated messages. This guide walks through adding
one — using a Brazilian *PIS/PASEP* validator as the example.

## 1. Write and test the GAD algorithm

The script receives the value as `param value`, `return true` when valid and
`throw "KEY"` to fail with a translatable message key. Keep it self-contained
(only the default builtins). See the [GAD reference](gad-reference.md).

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
    throw "pis"
}

weights := [3, 2, 9, 8, 7, 6, 5, 4, 3, 2]
sum := 0
for i := 0; i < 10; i++ {
    sum += digits[i] * weights[i]
}
dv := 11 - (sum % 11)
if dv >= 10 {
    dv = 0
}
if dv != digits[10] {
    throw "pis"
}
return true
```

Test the algorithm directly against the engine (no DB needed):

```go
func TestPIS(t *testing.T) {
    v := &validators.Validator{InitialValue: pisScript}
    if ok, _, err := v.Eval("120.6198.767-0"); err != nil || !ok {
        t.Fatalf("valid PIS: ok=%v err=%v", ok, err)
    }
    if ok, key, _ := v.Eval("123"); ok || key != "pis" {
        t.Fatalf("invalid PIS should throw \"pis\", got ok=%v key=%q", ok, key)
    }
}
```

## 2. Register (seed) it

Build the `Validator` with its `Initial*` fields — the algorithm, the Markdown
documentation and the per-language messages — and seed it (upsert by name,
preserving any user overrides):

```go
pis := &validators.Validator{
    Name:         "pis",
    Description:  "Brazilian PIS/PASEP",
    DefaultLang:  "pt-BR",
    InitialValue: pisScript,
    InitialDoc:   "# PIS/PASEP\n\nValidates … throws the key `pis` on failure.",
    InitialMessages: datatypes.NullJSONMap{
        "pt-BR": map[string]interface{}{"pis": "PIS/PASEP inválido."},
        "en-US": map[string]interface{}{"pis": "Invalid PIS/PASEP."},
    },
}
```

Register it through a plugin builder so it is also editable in the admin:

```go
pb.Use(validatorsadmin.New(db).Register(pis))
```

Message keys are just the strings the algorithm throws; add one entry per key,
per language. Missing languages fall back to `DefaultLang`, then to the key.

## 3. Use it

From Go — get a translated error when invalid:

```go
if err := validators.Check(db, "pis", value, lang); err != nil { … }
```

As an admin field validator (`FieldBuilder.Validators`), so the message is
translated and the algorithm stays editable:

```go
ed.Field("PIS").Validator(presets.FieldValidatorFunc(
    func(field *presets.FieldContext) (verr web.ValidationErrors) {
        s, _ := field.Value().(string)
        if s == "" { return }
        if err := validators.Check(db, "pis", s, lang); err != nil {
            verr.FieldError(field.Name, err.Error())
        }
        return
    }))
```

## 4. Per document type (people)

To route a document field by its type, register the validator on the people
`Builder` — it is seeded and applied to the Person `Document` field
automatically:

```go
peopleadmin.Configure(b, db, func(pb *peopleadmin.Builder) {
    pb.RegisterDocumentValidator(models.DocumentType("pis"), pis)
})
```

## Notes

- **Reset:** an empty `Value`/`Doc`/`Messages` override falls back to `Initial*`;
  clearing an override in the admin and saving resets it to the registered
  default.
- **Bytecode:** on save the effective algorithm is compiled and cached in the
  `Bytecode` column (not shown in the UI), so validation does not recompile.
- **Messages format:** `Messages` must be `{LANG: {KEY: VALUE}}` with string
  values (enforced by a field validator and `BeforeSave`).
