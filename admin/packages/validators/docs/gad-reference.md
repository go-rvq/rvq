# GAD reference for validators

Validation scripts are written in [GAD](https://github.com/gad-lang/gad). This is
the small subset used by validators; see the GAD project for the full language.

## The contract

- The value to analyze arrives as the `value` parameter: `param value`.
- `return true` (any truthy value) → **valid**.
- `throw "MESSAGE_KEY"` → **invalid**, with a translatable reason. The engine
  translates the key via the validator's `Messages` for the user's language.
- A plain falsy return (no throw) → generic failure (key `invalid`).
- Scripts are self-contained: only the default builtins, no `import`.

## Values and characters

`value` is a string. Index it to get a character (an integer code); compare and
convert with the numeric char literals:

```gad
c := value[i] // a character

if c >= '0' && c <= '9' {
    // is it a digit?
    d := int(c) - int('0') // 0..9
}

len(value) // length
```

## Control flow

C-style `for`, plus `for _, x in array`; `if/else`; `break`/`continue`:

```gad
for i := 0; i < len(value); i++ {
    // …
}

// iterate values
for _, n in [9, 10] {
    // …
}

allEqual := true

for i := 1; i < 11; i++ {
    if digits[i] != digits[0] {
        allEqual = false

        break
    }
}
```

## Arrays

```gad
digits := [] // empty array

digits = append(digits, x) // append

weights := [6, 5, 4, 3, 2] // literal

weights[i] // index
len(weights)
```

## Operators

Arithmetic `+ - * /`, modulo `%`, comparison `== != < <= > >=`, boolean
`&& || !`. Integer math is exact — the check-digit algorithms rely on `%` and
integer division.

## Throwing a message key

```gad
if len(digits) != 11 {
    throw "cpf" // Messages["<lang>"]["cpf"] → user message
}
```

A thrown **string** becomes the message key. Register the matching entries in the
validator's `Messages` (`{LANG: {KEY: VALUE}}`); unknown languages fall back to
`DefaultLang`, then to the key itself.

## Full example (CPF)

See [`people/models/validators.go`](../../people/models/validators.go) for the
commented CPF and CNPJ algorithms used in production, and
[developing.md](developing.md) for adding your own.
