package web

import "fmt"

type ValidationErrors struct {
	globalErrors []string
	fieldErrors  map[string][]string
}

func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{}
}

func (b *ValidationErrors) FieldError(fieldName string, message string) *ValidationErrors {
	if b.fieldErrors == nil {
		b.fieldErrors = make(map[string][]string)
	}
	b.fieldErrors[fieldName] = append(b.fieldErrors[fieldName], message)
	return b
}

func (b *ValidationErrors) FieldErrorf(fieldName string, format string, args ...any) *ValidationErrors {
	return b.FieldError(fieldName, fmt.Sprintf(format, args...))
}

func (b *ValidationErrors) GlobalError(message string) *ValidationErrors {
	b.globalErrors = append(b.globalErrors, message)
	return b
}

func (b *ValidationErrors) GetFieldErrors(fieldName string) (r []string) {
	if b.fieldErrors == nil {
		return
	}

	r = b.fieldErrors[fieldName]
	return
}

func (b *ValidationErrors) GetRemoveFieldErrors(fieldName string) (r []string) {
	if b.fieldErrors == nil {
		return
	}

	r = b.fieldErrors[fieldName]
	delete(b.fieldErrors, fieldName)
	return
}

// FieldErrorsHolder is an error that says what is wrong field by field — by
// the keys of the form, as it posts them —: a field's setter answering one
// (presets) has each said under its own field.
type FieldErrorsHolder interface {
	error
	FieldErrors() map[string][]string
}

// GetFieldErrorsUnder are the errors of the fields under key — the keys
// that go on from it: `key.a`, `key[0].b`, `key__x` —: what a field that
// draws fields of its own (a form in a field) says under each. They stay:
// a field of its own under it (a nested editor's) takes its own.
func (b *ValidationErrors) GetFieldErrorsUnder(key string) (r map[string][]string) {
	if b == nil || b.fieldErrors == nil || key == "" {
		return
	}
	for k, msgs := range b.fieldErrors {
		if len(k) <= len(key) || k[:len(key)] != key {
			continue
		}
		if c := k[len(key)]; c != '.' && c != '[' && c != '_' {
			continue
		}
		if r == nil {
			r = map[string][]string{}
		}
		r[k] = msgs
	}
	return
}

func (b *ValidationErrors) GetGlobalError() (r string) {
	if len(b.globalErrors) == 0 {
		return
	}
	return b.globalErrors[0]
}

func (b *ValidationErrors) GetGlobalErrors() (r []string) {
	return b.globalErrors
}

func (b *ValidationErrors) HaveErrors() bool {
	if len(b.globalErrors) > 0 {
		return true
	}
	if len(b.fieldErrors) > 0 {
		return true
	}
	return false
}

func (b *ValidationErrors) HaveGlobalErrors() bool {
	return len(b.globalErrors) > 0
}

func (b *ValidationErrors) HaveFieldErrors() bool {
	return len(b.fieldErrors) > 0
}

func (b *ValidationErrors) Merge(other ValidationErrors) {
	if len(other.fieldErrors) > 0 && b.fieldErrors == nil {
		b.fieldErrors = make(map[string][]string)
	}
	for name, errors := range other.fieldErrors {
		if _, ok := b.fieldErrors[name]; ok {
			b.fieldErrors[name] = append(b.fieldErrors[name], errors...)
		} else {
			b.fieldErrors[name] = errors
		}
	}

	b.globalErrors = append(b.globalErrors, other.globalErrors...)
}

func (b ValidationErrors) Error() string {
	if (b.globalErrors == nil && b.fieldErrors == nil) || (len(b.globalErrors)+len(b.fieldErrors) == 0) {
		return ""
	}
	return fmt.Sprintf("validation error global: %+v, fields: %+v", b.globalErrors, b.fieldErrors)
}
