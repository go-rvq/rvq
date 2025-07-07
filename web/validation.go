package web

type (
	Validator interface {
		Validate(obj interface{}, ctx *EventContext) (err ValidationErrors)
	}

	ValidatorFunc func(obj interface{}, ctx *EventContext) (err ValidationErrors)
)

func (f ValidatorFunc) Validate(obj interface{}, ctx *EventContext) (err ValidationErrors) {
	return f(obj, ctx)
}

type Validators []Validator

func (vh Validators) Validate(obj interface{}, ctx *EventContext) (err ValidationErrors) {
	for _, f := range vh {
		if err = f.Validate(obj, ctx); err.HaveErrors() {
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
