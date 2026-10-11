package schemaform

import (
	"bytes"
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// Record is a decoded record: its fields IN THE ORDER THE SCHEMA DECLARES them.
// A map would sort them alphabetically and rewrite the whole value on every
// save, so the order the schema gives is kept all the way to the file.
type Record []RecordField

// RecordField is one field of a decoded record.
type RecordField struct {
	Name  string
	Value any
}

// MarshalJSON writes the record as an object, in its own order — encoding/json
// would sort a map's keys and rewrite the whole value on every save.
func (r Record) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, f := range r {
		if i > 0 {
			buf.WriteByte(',')
		}
		name, err := json.Marshal(f.Name)
		if err != nil {
			return nil, err
		}
		value, err := json.Marshal(f.Value)
		if err != nil {
			return nil, err
		}
		buf.Write(name)
		buf.WriteByte(':')
		buf.Write(value)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// MarshalYAML writes the record as a mapping, in its own order.
func (r Record) MarshalYAML() (any, error) {
	n := &yaml.Node{Kind: yaml.MappingNode}
	for _, f := range r {
		var key, value yaml.Node
		if err := key.Encode(f.Name); err != nil {
			return nil, err
		}
		if err := value.Encode(f.Value); err != nil {
			return nil, err
		}
		n.Content = append(n.Content, &key, &value)
	}
	return n, nil
}
