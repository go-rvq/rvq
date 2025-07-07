package presets

import "github.com/go-rvq/rvq/web"

type (
	Validator interface {
		Validate(obj interface{}, mode FieldModeStack, ctx *web.EventContext) (err web.ValidationErrors)
	}

	ValidatorFunc func(obj interface{}, mode FieldModeStack, ctx *web.EventContext) (err web.ValidationErrors)
)

func (f ValidatorFunc) Validate(obj interface{}, mode FieldModeStack, ctx *web.EventContext) (err web.ValidationErrors) {
	return f(obj, mode, ctx)
}

type Validators []Validator

func (vh Validators) Validate(obj interface{}, mode FieldModeStack, ctx *web.EventContext) (err web.ValidationErrors) {
	for _, f := range vh {
		if err = f.Validate(obj, mode, ctx); err.HaveErrors() {
			return
		}
	}
	return
}

func (vh *Validators) Append(v ...Validator) {
	*vh = append(*vh, v...)
}

func (vh *Validators) AppendFunc(v ...ValidatorFunc) {
	for _, f := range v {
		*vh = append(*vh, f)
	}
}

func (vh *Validators) Prepend(v ...Validator) {
	*vh = append(v, (*vh)...)
}
