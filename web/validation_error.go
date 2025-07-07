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

func (b *ValidationErrors) Error() string {
	return fmt.Sprintf("validation error global: %+v, fields: %+v", b.globalErrors, b.fieldErrors)
}
