package models

import (
	"errors"

	"github.com/go-rvq/rvq/thirdpart/gorm/datatypes"
	"gorm.io/gorm"
)

// Validator is a named validation rule whose algorithm is a GAD script (see the
// engine): the script receives the value to analyze as its `value` parameter and
// returns True/False. Records are inserted by the application; end users may only
// override Value, Doc and Messages. When an override is empty the matching
// Initial* value is used (so clearing a field and saving resets it to the
// registered default).
//
// Messages holds the per-language error messages the algorithm may reference,
// in the shape {LANG: {MSG_KEY: MSG_VALUE}}. When a language is missing, the
// DefaultLang entry is used.
type Validator struct {
	Base

	// Name is the stable key the application looks the validator up by (e.g.
	// "cpf", "cnpj"). Unique.
	Name string `gorm:"size:100;not null;uniqueIndex:uidx_validators_name"`

	// Description is a short human description of what the validator checks.
	Description string `gorm:"size:500"`

	// DefaultLang is the fallback language for Messages (e.g. "pt-BR").
	DefaultLang string `gorm:"size:20;not null;default:pt-BR"`

	// Initial* are the registered defaults, set by the application. They are not
	// meant to be edited by end users.
	InitialValue    string                `gorm:"type:text"`
	InitialDoc      string                `gorm:"type:text"`
	InitialMessages datatypes.NullJSONMap `gorm:"type:jsonb"`

	// Value/Doc/Messages are the user overrides; empty means "use the Initial*".
	Value    string                `gorm:"type:text"`
	Doc      string                `gorm:"type:text"`
	Messages datatypes.NullJSONMap `gorm:"type:jsonb"`

	// Bytecode caches the compiled effective algorithm. It is not shown in the UI
	// and is (re)computed on save from the effective value (see BeforeSave), so
	// validation runs the cached bytecode instead of recompiling every call.
	Bytecode []byte `gorm:"type:bytea"`
}

func (Validator) TableName() string { return "validators" }

// BeforeSave validates the Messages shape and (re)compiles the effective
// algorithm into Bytecode. The bytecode only depends on InitialValue/Value, both
// known at save time, so it is refreshed on every save. A syntactically invalid
// algorithm makes the save fail.
func (v *Validator) BeforeSave(*gorm.DB) error {
	if err := ValidateMessagesFormat(v.Messages); err != nil {
		return err
	}
	if err := ValidateMessagesFormat(v.InitialMessages); err != nil {
		return err
	}
	bc, err := EncodeScript(v.EffectiveValue())
	if err != nil {
		return err
	}
	v.Bytecode = bc
	return nil
}

func (v *Validator) String() string { return v.Name }

// EffectiveValue returns the GAD algorithm in use: Value, or InitialValue when
// Value is empty.
func (v *Validator) EffectiveValue() string {
	if v.Value != "" {
		return v.Value
	}
	return v.InitialValue
}

// EffectiveDoc returns the documentation in use (markdown): Doc, or InitialDoc
// when Doc is empty.
func (v *Validator) EffectiveDoc() string {
	if v.Doc != "" {
		return v.Doc
	}
	return v.InitialDoc
}

// EffectiveMessages returns the message map in use: Messages, or InitialMessages
// when Messages is empty.
func (v *Validator) EffectiveMessages() datatypes.NullJSONMap {
	if len(v.Messages) > 0 {
		return v.Messages
	}
	return v.InitialMessages
}

// Message returns the translated message for key in lang, falling back to
// DefaultLang and then to the key itself. It reads the {LANG: {KEY: VALUE}}
// message map.
func (v *Validator) Message(lang, key string) string {
	msgs := v.EffectiveMessages()
	if s, ok := lookupMessage(msgs, lang, key); ok {
		return s
	}
	if lang != v.DefaultLang {
		if s, ok := lookupMessage(msgs, v.DefaultLang, key); ok {
			return s
		}
	}
	return key
}

func lookupMessage(m datatypes.NullJSONMap, lang, key string) (string, bool) {
	langAny, ok := m[lang]
	if !ok {
		return "", false
	}
	langMap, ok := langAny.(map[string]interface{})
	if !ok {
		return "", false
	}
	s, ok := langMap[key].(string)
	return s, ok
}

// ErrInvalidMessagesFormat is returned when Messages is not in the
// {LANG: {MSG_KEY: MSG_VALUE}} shape.
var ErrInvalidMessagesFormat = errors.New("messages must be {language: {key: text}} with string values")

// ValidateMessagesFormat checks that m matches the {LANG: {KEY: VALUE}} shape:
// each top-level value is an object whose values are strings. An empty/nil map
// is valid (it means "use the initial messages"). Used both by the admin field
// validator and the model BeforeSave.
func ValidateMessagesFormat(m datatypes.NullJSONMap) error {
	for _, langAny := range m {
		langMap, ok := langAny.(map[string]interface{})
		if !ok {
			return ErrInvalidMessagesFormat
		}
		for _, msgAny := range langMap {
			if _, ok := msgAny.(string); !ok {
				return ErrInvalidMessagesFormat
			}
		}
	}
	return nil
}
