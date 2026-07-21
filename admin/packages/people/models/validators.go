package models

import (
	validators "github.com/go-rvq/rvq/admin/packages/validators/models"
	"github.com/go-rvq/rvq/thirdpart/gorm/datatypes"
	"gorm.io/gorm"
)

// Document-validator names, looked up in the validators registry.
const (
	CPFValidatorName  = "cpf"
	CNPJValidatorName = "cnpj"
)

// cpfAlgorithm validates a Brazilian CPF (11 digits, two check digits). It is
// written in GAD: it receives the value to analyze as `param value`, returns
// `true` when valid and throws the message key "cpf" on any failure, so the
// engine can translate the error to the user's language.
const cpfAlgorithm = `
param value

// Keep only the decimal digits (accepts "529.982.247-25" or "52998224725").
digits := []
for i := 0; i < len(value); i++ {
    c := value[i]
    if c >= '0' && c <= '9' {
        digits = append(digits, int(c) - int('0'))
    }
}

// A CPF has exactly 11 digits.
if len(digits) != 11 {
    throw "cpf"
}

// Reject sequences of a single repeated digit (e.g. 111.111.111-11), which pass
// the check-digit math but are not valid CPFs.
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

// Two check digits. For n in {9, 10}: sum digits[0..n) weighted by (n+1-i),
// the check digit is (sum*10 % 11) % 10 and must equal digits[n].
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

// cnpjAlgorithm validates a Brazilian CNPJ (14 digits, two check digits). Same
// contract as cpfAlgorithm; it throws the message key "cnpj" on failure.
const cnpjAlgorithm = `
param value

// Keep only the decimal digits.
digits := []
for i := 0; i < len(value); i++ {
    c := value[i]
    if c >= '0' && c <= '9' {
        digits = append(digits, int(c) - int('0'))
    }
}

// A CNPJ has exactly 14 digits.
if len(digits) != 14 {
    throw "cnpj"
}

// Reject a single repeated digit.
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

// Two check digits, using the standard weight window. For n in {12, 13}: sum
// digits[0..n) times the trailing n weights; dv = 11 - (sum % 11), clamped to 0
// when < 2, and must equal digits[n].
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

const cpfDoc = "# CPF\n\n" +
	"Validates a Brazilian **CPF** (Cadastro de Pessoas Físicas).\n\n" +
	"- Keeps only the digits, so `529.982.247-25` and `52998224725` are equivalent.\n" +
	"- Requires exactly **11 digits**.\n" +
	"- Rejects sequences of a single repeated digit (e.g. `111.111.111-11`).\n" +
	"- Verifies the **two check digits**: for `n` in {9,10}, `dv = (sum*10 % 11) % 10` " +
	"where `sum` weights `digits[0..n)` by `n+1-i`.\n\n" +
	"On failure the algorithm throws the message key `cpf`."

const cnpjDoc = "# CNPJ\n\n" +
	"Validates a Brazilian **CNPJ** (Cadastro Nacional da Pessoa Jurídica).\n\n" +
	"- Keeps only the digits.\n" +
	"- Requires exactly **14 digits**.\n" +
	"- Rejects a single repeated digit.\n" +
	"- Verifies the **two check digits** with the standard weight window `[6,5,4,3,2,9,8,7,6,5,4,3,2]`; " +
	"`dv = 11 - (sum % 11)`, clamped to 0 when `< 2`.\n\n" +
	"On failure the algorithm throws the message key `cnpj`."

// DocumentValidators returns the CPF and CNPJ validator registrations (their
// GAD algorithm, documentation and per-language messages). The application seeds
// them so they are editable in the validators admin and used to validate person
// documents (see SeedDocumentValidators).
func DocumentValidators() []*validators.Validator {
	return []*validators.Validator{
		{
			Name:         CPFValidatorName,
			Description:  "Brazilian CPF (individual taxpayer registry)",
			DefaultLang:  "pt-BR",
			InitialValue: cpfAlgorithm,
			InitialDoc:   cpfDoc,
			InitialMessages: datatypes.NullJSONMap{
				"pt-BR": map[string]interface{}{"cpf": "CPF inválido."},
				"en-US": map[string]interface{}{"cpf": "Invalid CPF."},
			},
		},
		{
			Name:         CNPJValidatorName,
			Description:  "Brazilian CNPJ (company taxpayer registry)",
			DefaultLang:  "pt-BR",
			InitialValue: cnpjAlgorithm,
			InitialDoc:   cnpjDoc,
			InitialMessages: datatypes.NullJSONMap{
				"pt-BR": map[string]interface{}{"cnpj": "CNPJ inválido."},
				"en-US": map[string]interface{}{"cnpj": "Invalid CNPJ."},
			},
		},
	}
}

// SeedDocumentValidators migrates the validators table and (re)registers the CPF
// and CNPJ validators, preserving any user overrides (see validators.Seed).
func SeedDocumentValidators(db *gorm.DB) error {
	if err := validators.AutoMigrate(db); err != nil {
		return err
	}
	for _, v := range DocumentValidators() {
		if err := validators.Seed(db, v); err != nil {
			return err
		}
	}
	return nil
}

// ValidatorNameForDocumentType returns the validator name for a document type,
// or "" when the type has no document validation (Other).
func ValidatorNameForDocumentType(t DocumentType) string {
	switch t {
	case DocumentCPF:
		return CPFValidatorName
	case DocumentCNPJ:
		return CNPJValidatorName
	}
	return ""
}
