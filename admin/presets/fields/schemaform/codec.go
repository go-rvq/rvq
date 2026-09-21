package schemaform

import (
	"bytes"
	"encoding/json"
	"strings"

	"gopkg.in/yaml.v3"
)

// EncodeFunc writes the value the form edited to the text that is STORED, and
// DecodeFunc reads that text back into the value the form opens over. A form is
// the same form whether the value travels as YAML, as JSON or as anything else
// — only these two funcs know which.
type (
	EncodeFunc func(value any) (string, error)
	DecodeFunc func(stored string) (any, error)
)

// Codec sets how the value is carried outside the form. The default pair is
// YAML; JSONEncode/JSONDecode are here for a value stored as JSON, and an
// application may pass its own.
func (b *Builder) Codec(encode EncodeFunc, decode DecodeFunc) *Builder {
	b.encode, b.decode = encode, decode
	return b
}

// EncodeValue writes a decoded form value out.
func (b *Builder) EncodeValue(value any) (string, error) {
	if b.encode != nil {
		return b.encode(value)
	}
	return YAMLEncode(value)
}

// DecodeValue reads a stored value back. An empty text is no value at all —
// nil, which the caller reads as "start from the empty shape".
func (b *Builder) DecodeValue(stored string) (any, error) {
	if strings.TrimSpace(stored) == "" {
		return nil, nil
	}
	if b.decode != nil {
		return b.decode(stored)
	}
	return YAMLDecode(stored)
}

// YAMLEncode writes the value as YAML — the default, and what a Record writes
// itself as (in the order the schema declares, see Record.MarshalYAML).
func YAMLEncode(value any) (string, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(value); err != nil {
		return "", err
	}
	if err := enc.Close(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// YAMLDecode reads a YAML value.
func YAMLDecode(stored string) (value any, err error) {
	err = yaml.Unmarshal([]byte(stored), &value)
	return
}

// JSONEncode writes the value as indented JSON, records keeping the order the
// schema declares (see Record.MarshalJSON).
func JSONEncode(value any) (string, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(value); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// JSONDecode reads a JSON value.
func JSONDecode(stored string) (value any, err error) {
	err = json.Unmarshal([]byte(stored), &value)
	return
}
