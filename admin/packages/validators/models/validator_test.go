package models

import (
	"testing"

	"github.com/go-rvq/rvq/thirdpart/gorm/datatypes"
)

// cpfScript validates a Brazilian CPF; it throws the message key "cpf" on any
// failure so the engine can translate it.
const cpfScript = `
param value

// keep only the digits
digits := []
for i := 0; i < len(value); i++ {
    c := value[i]
    if c >= '0' && c <= '9' {
        digits += int(c) - int('0')
    }
}
if len(digits) != 11 {
    throw "cpf"
}

// reject sequences with all equal digits (e.g. 111.111.111-11)
allEqual := true
for i := 1; i < 11; i++ {
    if digits[i] != digits[0] {
        allEqual = false
        break
    }
}
if allEqual {
    throw "cpf"
}

// two check digits: weighted sum, dv = (sum*10 % 11) % 10
for _, n in [9, 10] {
    sum := 0
    for i := 0; i < n; i++ {
        sum += digits[i] * (n + 1 - i)
    }
    dv := (sum * 10) % 11 % 10
    if dv != digits[n] {
        throw "cpf"
    }
}
return true
`

// cnpjScript validates a Brazilian CNPJ; it throws the message key "cnpj".
const cnpjScript = `
param value

digits := []
for i := 0; i < len(value); i++ {
    c := value[i]
    if c >= '0' && c <= '9' {
        digits += int(c) - int('0')
    }
}
if len(digits) != 14 {
    throw "cnpj"
}

allEqual := true
for i := 1; i < 14; i++ {
    if digits[i] != digits[0] {
        allEqual = false
        break
    }
}
if allEqual {
    throw "cnpj"
}

weights := [6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2]
for _, n in [12, 13] {
    sum := 0
    for i := 0; i < n; i++ {
        sum += digits[i] * weights[len(weights) - n + i]
    }
    dv := sum % 11
    if dv < 2 {
        dv = 0
    } else {
        dv = 11 - dv
    }
    if dv != digits[n] {
        throw "cnpj"
    }
}
return true
`

func TestEngineReturnAndThrow(t *testing.T) {
	// plain boolean return
	if ok, key, err := (&Validator{InitialValue: "param value\nreturn len(value) == 3"}).Eval("abc"); err != nil || !ok || key != "" {
		t.Fatalf("bool return: ok=%v key=%q err=%v", ok, key, err)
	}
	// throw a message key
	ok, key, err := (&Validator{InitialValue: `param value; if len(value) != 3 { throw "too_short" }; return true`}).Eval("ab")
	if err != nil || ok || key != "too_short" {
		t.Fatalf("throw: ok=%v key=%q err=%v", ok, key, err)
	}
	// compile error is a real error
	if _, _, err := (&Validator{InitialValue: "param value\nreturn ("}).Eval("x"); err == nil {
		t.Fatalf("expected compile error")
	}
}

func TestCPFValidator(t *testing.T) {
	v := &Validator{InitialValue: cpfScript}
	valid := []string{"529.982.247-25", "52998224725"}
	for _, d := range valid {
		if ok, _, err := v.Eval(d); err != nil || !ok {
			t.Errorf("CPF %q: ok=%v err=%v, want valid", d, ok, err)
		}
	}
	invalid := []string{"111.111.111-11", "529.982.247-24", "123", ""}
	for _, d := range invalid {
		ok, key, err := v.Eval(d)
		if err != nil {
			t.Errorf("CPF %q: err=%v", d, err)
		}
		if ok || key != "cpf" {
			t.Errorf("CPF %q: ok=%v key=%q, want invalid with key cpf", d, ok, key)
		}
	}
}

func TestCNPJValidator(t *testing.T) {
	v := &Validator{InitialValue: cnpjScript}
	if ok, _, err := v.Eval("11.222.333/0001-81"); err != nil || !ok {
		t.Errorf("valid CNPJ: ok=%v err=%v", ok, err)
	}
	for _, d := range []string{"11.222.333/0001-80", "00000000000000", "abc"} {
		if ok, key, _ := v.Eval(d); ok || key != "cnpj" {
			t.Errorf("CNPJ %q: ok=%v key=%q, want invalid cnpj", d, ok, key)
		}
	}
}

func TestCheckTranslatesThrownKey(t *testing.T) {
	v := &Validator{
		DefaultLang:  "pt-BR",
		InitialValue: cpfScript,
		InitialMessages: datatypes.NullJSONMap{
			"pt-BR": map[string]interface{}{"cpf": "CPF inválido."},
			"en-US": map[string]interface{}{"cpf": "Invalid CPF."},
		},
	}
	if err := v.Check("52998224725", "pt-BR"); err != nil {
		t.Errorf("valid CPF should pass: %v", err)
	}
	if err := v.Check("123", "en-US"); err == nil || err.Error() != "Invalid CPF." {
		t.Errorf("en-US message = %v, want \"Invalid CPF.\"", err)
	}
	// missing language falls back to DefaultLang
	if err := v.Check("123", "fr-FR"); err == nil || err.Error() != "CPF inválido." {
		t.Errorf("fallback message = %v, want \"CPF inválido.\"", err)
	}
}

func TestEffectiveFallbacks(t *testing.T) {
	v := &Validator{InitialValue: "iv", InitialDoc: "id", InitialMessages: datatypes.NullJSONMap{"pt-BR": map[string]interface{}{"k": "i"}}}
	if v.EffectiveValue() != "iv" || v.EffectiveDoc() != "id" || v.Message("pt-BR", "k") != "i" {
		t.Errorf("empty overrides should use Initial*: %q %q %q", v.EffectiveValue(), v.EffectiveDoc(), v.Message("pt-BR", "k"))
	}
	v.Value, v.Doc = "ov", "od"
	v.Messages = datatypes.NullJSONMap{"pt-BR": map[string]interface{}{"k": "o"}}
	if v.EffectiveValue() != "ov" || v.EffectiveDoc() != "od" || v.Message("pt-BR", "k") != "o" {
		t.Errorf("overrides should win")
	}
}

func TestValidateMessagesFormat(t *testing.T) {
	good := datatypes.NullJSONMap{"pt-BR": map[string]interface{}{"cpf": "x"}}
	if err := ValidateMessagesFormat(good); err != nil {
		t.Errorf("valid shape rejected: %v", err)
	}
	if err := ValidateMessagesFormat(nil); err != nil {
		t.Errorf("empty should be valid: %v", err)
	}
	badTop := datatypes.NullJSONMap{"pt-BR": "not-a-map"}
	if err := ValidateMessagesFormat(badTop); err == nil {
		t.Errorf("non-object language value should be rejected")
	}
	badVal := datatypes.NullJSONMap{"pt-BR": map[string]interface{}{"cpf": 123}}
	if err := ValidateMessagesFormat(badVal); err == nil {
		t.Errorf("non-string message value should be rejected")
	}
}

func TestBytecodeCacheRun(t *testing.T) {
	// pre-encode the bytecode, then eval uses the cache (no source recompile).
	bc, err := EncodeScript(cpfScript)
	if err != nil {
		t.Fatal(err)
	}
	v := &Validator{Bytecode: bc} // note: no InitialValue/Value
	if ok, _, err := v.Eval("52998224725"); err != nil || !ok {
		t.Errorf("cached bytecode valid CPF: ok=%v err=%v", ok, err)
	}
	if ok, key, _ := v.Eval("123"); ok || key != "cpf" {
		t.Errorf("cached bytecode invalid CPF: ok=%v key=%q", ok, key)
	}
}
