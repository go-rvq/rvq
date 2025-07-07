package presets

import "github.com/qor5/web/v3"

type (
	FieldValidatorFunc func(field *FieldContext) (err web.ValidationErrors)
	FieldValidator     interface {
		Validate(field *FieldContext) (err web.ValidationErrors)
	}
	FieldValidators []FieldValidator
	FieldSetuper    interface {
		InitField(f *FieldBuilder)
		ConfigureField(f *FieldBuilder)
	}

	FieldSetupers []FieldSetuper

	FieldStringer interface {
		FieldString(field *FieldContext) string
	}

	FieldsBuilderInterface interface {
		Field(name string) *FieldBuilder
	}
)

func (s FieldSetupers) InitField(f *FieldBuilder) {
	for _, setuper := range s {
		setuper.InitField(f)
	}
}

func (s FieldSetupers) ConfigureField(f *FieldBuilder) {
	for _, setuper := range s {
		setuper.ConfigureField(f)
	}
}

func (fv FieldValidatorFunc) Validate(field *FieldContext) (err web.ValidationErrors) {
	return fv(field)
}

func (fv FieldValidators) Validate(field *FieldContext) (err web.ValidationErrors) {
	for _, v := range fv {
		err.Merge(v.Validate(field))
	}
	return
}

func (fv *FieldValidators) Append(v ...FieldValidator) {
	*fv = append(*fv, v...)
}

func (fv *FieldValidators) AppendFunc(v ...FieldValidatorFunc) {
	for _, vf := range v {
		fv.Append(vf)
	}
}

type (
	FieldValueFormatter interface {
		FormatValue(field *FieldContext) (err error)
	}
	FieldValueFormatterFunc func(field *FieldContext) (err error)
	FieldValueFormatters    []FieldValueFormatter
)

func (f FieldValueFormatterFunc) FormatValue(field *FieldContext) (err error) {
	return f(field)
}

func (f FieldValueFormatters) FormatValue(field *FieldContext) (err error) {
	for _, f := range f {
		if err = f.FormatValue(field); err != nil {
			return
		}
	}
	return
}

func (f *FieldValueFormatters) Append(formatters ...FieldValueFormatter) {
	*f = append(*f, formatters...)
}

func (f *FieldValueFormatters) AppendFunc(formatters ...FieldValueFormatterFunc) {
	for _, formatter := range formatters {
		f.Append(formatter)
	}
}

func (b *FieldBuilder) Validator(v ...FieldValidator) *FieldBuilder {
	b.Validators.Append(v...)
	return b
}
