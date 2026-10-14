package schemaform

import (
	"errors"
	"regexp"
	"strings"

	"github.com/gad-lang/gad"
	"github.com/go-rvq/rvq/web/str_utils"
)

// TypedGadValue is v, a record of fields, as gad holds it (GadValue), each
// value of the type its field declares: a decimal — kept a float64 by the
// Decoder — a gad.Decimal, which a class's `decimal` field takes (a
// gad.Float it refuses); a record of a class, or a list of them, its fields'
// the same way.
func TypedGadValue(fields []*Field, v any) gad.Object {
	if len(fields) == 0 {
		return GadValue(v)
	}
	var m map[string]any
	switch x := v.(type) {
	case Record:
		m = x.Map()
	case map[string]any:
		m = x
	default:
		return GadValue(v)
	}
	dict := make(gad.Dict, len(m))
	for k, e := range m {
		dict[k] = typedFieldValue(fieldNamed(fields, k), e)
	}
	return dict
}

func typedFieldValue(f *Field, v any) gad.Object {
	if f == nil || v == nil {
		return GadValue(v)
	}
	if f.Type == "decimal" {
		switch n := v.(type) {
		case float64:
			return gad.DecimalFromFloat(gad.Float(n))
		case int64:
			return gad.DecimalFromFloat(gad.Float(n))
		}
	}
	if s := f.Schema; s != nil {
		if s.Slice {
			if items, ok := v.([]any); ok {
				arr := make(gad.Array, len(items))
				for i, e := range items {
					if s.Item != nil {
						arr[i] = typedFieldValue(s.Item, e)
					} else {
						arr[i] = TypedGadValue(s.Fields, e)
					}
				}
				return arr
			}
		}
		return TypedGadValue(s.Fields, v)
	}
	return GadValue(v)
}

func fieldNamed(fields []*Field, name string) *Field {
	for _, f := range fields {
		if f.Name == name {
			return f
		}
	}
	return nil
}

// fieldLabel is a field's `[label=…]`, else its name humanized.
func fieldLabel(f *Field) string {
	if l := f.Meta.String("label"); l != "" {
		return l
	}
	return str_utils.HumanizeString(f.Name)
}

// labelPath is the way to the field named name among fields — of a class's,
// the class's field before it —, by their labels; nil, none of that name.
func labelPath(fields []*Field, name string) []string {
	for _, f := range fields {
		if f.Name == name {
			return []string{fieldLabel(f)}
		}
	}
	for _, f := range fields {
		if f.Schema == nil {
			continue
		}
		sub := f.Schema.Fields
		if f.Schema.Item != nil && f.Schema.Item.Schema != nil {
			sub = f.Schema.Item.Schema.Fields
		}
		if p := labelPath(sub, name); p != nil {
			return append([]string{fieldLabel(f)}, p...)
		}
	}
	return nil
}

var fieldErrRe = regexp.MustCompile(`field "([^"]+)" (.*)$`)

// LabeledError is an error of a value of a field said by the labels of the
// way to it — the step's first —: a page shows them apart (Issue.Labels),
// in bold.
type LabeledError struct {
	Labels  []string
	Message string
}

func (e *LabeledError) Error() string {
	return strings.Join(e.Labels, " › ") + ": " + e.Message
}

// issueOf is err as the form's whole Issue: of a LabeledError, its labels
// and message apart.
func issueOf(err error) Issue {
	var le *LabeledError
	if errors.As(err, &le) {
		return Issue{Code: CodeInvalid, Message: le.Message, Labels: le.Labels}
	}
	return Issue{Code: CodeInvalid, Message: err.Error()}
}

// stepError is err, of building the record of the step as its class's
// instance, said by the labels (LabeledError): `Property Details › Asking
// price: expects decimal, got float` — not the names of the class.
func stepError(step *Step, err error) error {
	msg := err.Error()
	if m := fieldErrRe.FindStringSubmatch(msg); m != nil {
		var fields []*Field
		for _, f := range step.Fields {
			if f.Schema != nil && !step.Group {
				fields = append(fields, f.Schema.Fields...)
			} else {
				fields = append(fields, f)
			}
		}
		if p := labelPath(fields, m[1]); p != nil {
			return &LabeledError{Labels: append([]string{step.Label}, p...), Message: m[2]}
		}
	}
	return &LabeledError{Labels: []string{step.Label}, Message: msg}
}
