//go:build !formscore

package schemaform

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/go-rvq/rvq/web"
)

// Decode reads back the value the form posted for a schema drawn under key.
//
// The browser flattens what the form holds by the same two rules the builder
// binds by: an array indexes (`key[0]`) and a record names (`key.label`), so a
// list of records arrives as `key[0].label`. The schema says which shape those
// names have and what each value IS — an int is an int, a switch is a bool —
// and the result is the value itself: a Record, a list of them, or a list of
// plain values, ready to be written out.
//
// It reads; it does not check: Builder.DecodeForm does, the same reader
// (Decoder) saying what is wrong.
func (s *Schema) Decode(values url.Values, key string) any {
	v, _ := (&Decoder{Schema: s, NoValidations: true}).Value(Nest(values, key))
	return v
}

// DecodeForm reads the posted form back AND checks it, by the one logic of a
// form (Decoder) — the site's, the browser's —: a field that holds one of a
// closed list of values — an enum the schema declared, or a list
// EnumItemsFunc answered with — may only hold a value that list offered.
//
// The value comes back whole either way, so a caller may show the errors beside
// the form the user is still editing.
//
// What is wrong field by field — left empty when it is required, off its
// limits ([min, max, step, minlength, maxlength, pattern, mask]), a range
// out of order, refused by its `[validation=…]` — comes back as FieldErrors,
// by the keys the form posts (`key.a`, `key.list[0].b`), said in the
// request's language: presets has each said under its field.
func (b *Builder) DecodeForm(ctx *web.EventContext, schema *Schema, values url.Values, key string) (any, error) {
	d := b.Decoder(schema)
	d.Items = func(path string, f *Field) ([]EnumItem, error) {
		return b.itemsOf(ctx, path, f)
	}
	value, issues := d.Value(Nest(values, key))
	if len(issues) == 0 {
		return value, nil
	}
	msgs := messagesOf(ctx)
	fields := map[string][]string{}
	var errs []error
	for _, e := range issues {
		k := key
		switch {
		case e.Path == "":
		case strings.HasPrefix(e.Path, "["):
			k += e.Path
		case k == "":
			k = e.Path
		default:
			k += "." + e.Path
		}
		switch e.Code {
		case CodeItems:
			// a value never checked is not one that passed: the save fails
			err := fmt.Errorf("%s: %s: %w", where(e.Field, e.Path), msgs.ItemsUnavailableOnSave, e.Err)
			errs = append(errs, err)
			fields[k] = append(fields[k], err.Error())
			continue
		case CodeRequired:
			// of a list of values: one must be chosen
			if items, _ := b.itemsOf(ctx, SchemaPath(e.Path), e.Field); len(items) > 0 {
				fields[k] = append(fields[k], msgs.ChooseValue)
				continue
			}
		case CodeChoice:
			items, _ := b.itemsOf(ctx, SchemaPath(e.Path), e.Field)
			names := make([]string, len(items))
			for i, it := range items {
				names[i] = it.Name
			}
			err := fmt.Errorf("%s: "+msgs.NotAmongItems, where(e.Field, e.Path), e.Value, strings.Join(names, ", "))
			errs = append(errs, err)
			fields[k] = append(fields[k], err.Error())
			continue
		}
		fields[k] = append(fields[k], issueMessage(msgs, e))
	}
	return value, &FieldErrors{Fields: fields, Others: errors.Join(errs...)}
}

// issueMessage is what is wrong (e), said by msgs.
func issueMessage(msgs *Messages, e Issue) string {
	f := e.Field
	if f == nil {
		f = &Field{}
	}
	switch e.Code {
	case CodeRequired:
		return msgs.Required
	case CodeTooShort:
		return fmt.Sprintf(msgs.TooShort, f.MinLength)
	case CodeTooLong:
		n := f.MaxLength
		if n == 0 {
			n = MaxLine
			if f.Type == "text" || f.Type == "html" {
				n = MaxText
			}
		}
		return fmt.Sprintf(msgs.TooLong, n)
	case CodePattern:
		return msgs.Pattern
	case CodeMask:
		return fmt.Sprintf(msgs.Mask, f.Mask)
	case CodeMin:
		return fmt.Sprintf(msgs.Min, f.Min)
	case CodeMax:
		return fmt.Sprintf(msgs.Max, f.Max)
	case CodeStep:
		return fmt.Sprintf(msgs.Step, f.Step)
	case CodeRange:
		return msgs.RangeOrder
	case CodeItemsMin:
		return fmt.Sprintf(msgs.ItemsMin, f.MinItems)
	case CodeItemsMax:
		return fmt.Sprintf(msgs.ItemsMax, f.MaxItems)
	case CodeValidation:
		if e.Message != "" {
			return e.Message
		}
	}
	return msgs.Invalid
}

// FieldErrors is what is wrong with a form posted, by the keys of its
// fields, as the form posts them (web.FieldErrorsHolder); Others, what is
// of no field.
type FieldErrors struct {
	Fields map[string][]string
	Others error
}

func (e *FieldErrors) Error() string {
	keys := make([]string, 0, len(e.Fields))
	for k := range e.Fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var lines []string
	for _, k := range keys {
		for _, m := range e.Fields[k] {
			lines = append(lines, k+": "+m)
		}
	}
	if len(lines) == 0 && e.Others != nil {
		// (the errors of a field are its lines already)
		lines = append(lines, e.Others.Error())
	}
	return strings.Join(lines, "\n")
}

// FieldErrors are the errors by key (web.FieldErrorsHolder).
func (e *FieldErrors) FieldErrors() map[string][]string { return e.Fields }

func (e *FieldErrors) Unwrap() error { return e.Others }

// KeepReadOnly is value — what a form posted, decoded — with the fields
// declared `get name Type` (Field.ReadOnly) set back to what stored — the
// value as it was — holds: a post does not change what is only shown. The
// records of a list are matched by their index. stored may be what the
// decoder gave, or what YAML or JSON read (maps and slices).
func (s *Schema) KeepReadOnly(value, stored any) any {
	if s == nil || stored == nil {
		return value
	}
	if s.Slice {
		items, ok := value.([]any)
		old, ok2 := asSlice(stored)
		if !ok || !ok2 {
			return value
		}
		for i := range items {
			if i >= len(old) {
				break
			}
			if s.Item != nil {
				if s.Item.Schema != nil {
					items[i] = s.Item.Schema.KeepReadOnly(items[i], old[i])
				}
				continue
			}
			items[i] = s.keepRecord(items[i], old[i])
		}
		return items
	}
	return s.keepRecord(value, stored)
}

func (s *Schema) keepRecord(value, stored any) any {
	rec, ok := value.(Record)
	if !ok {
		return value
	}
	for i, rf := range rec {
		f := s.field(rf.Name)
		if f == nil {
			continue
		}
		old, has := fieldOf(stored, rf.Name)
		switch {
		case f.ReadOnly:
			if has {
				rec[i].Value = old
			}
		case f.Schema != nil && has:
			rec[i].Value = f.Schema.KeepReadOnly(rf.Value, old)
		}
	}
	return rec
}

// fieldOf is the field name of a record as stored: a Record, or a map YAML or
// JSON read.
func fieldOf(stored any, name string) (any, bool) {
	switch r := stored.(type) {
	case Record:
		for _, f := range r {
			if f.Name == name {
				return f.Value, true
			}
		}
	case map[string]any:
		v, ok := r[name]
		return v, ok
	case map[any]any:
		v, ok := r[name]
		return v, ok
	}
	return nil, false
}

func asSlice(v any) ([]any, bool) {
	switch s := v.(type) {
	case []any:
		return s, true
	case []map[string]any:
		out := make([]any, len(s))
		for i := range s {
			out[i] = s[i]
		}
		return out, true
	}
	return nil, false
}

// where names a field in an error: its path, and — at the root, where there is
// none — its own name.
func where(f *Field, path string) string {
	if path != "" || f == nil {
		return path
	}
	return f.Name
}
