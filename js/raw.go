package js

import (
	"bytes"
	"encoding/json"
	"sort"
	"strings"
)

type Raw string

func (v Raw) MarshalJSON() ([]byte, error) {
	return []byte(v), nil
}

type Object map[string]any

func (o Object) MarshalJSON() (_ []byte, err error) {
	if len(o) == 0 {
		return []byte("{}"), nil
	}

	var (
		out   bytes.Buffer
		keys  []string
		write = func(k string) {
			kb, _ := json.Marshal(k)
			out.Write(kb)
			out.WriteRune(':')
			v := o[k]
			switch t := v.(type) {
			case Raw:
				out.WriteByte('(')
				out.Write([]byte(t))
				out.WriteByte(')')
			case []byte:
				out.WriteByte('(')
				out.Write(t)
				out.WriteByte(')')
			case RawSlice:
				out.WriteString(t.String())
			default:
				if vb, err := json.Marshal(v); err == nil {
					out.Write(vb)
				}
			}
		}
	)

	for k := range o {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys[:len(keys)-1] {
		write(k)
		if err != nil {
			return
		}
		out.WriteByte(',')
	}

	write(keys[len(keys)-1])
	if err != nil {
		return
	}

	return out.Bytes(), nil
}

type RawSlice []string

func (v RawSlice) String() (s string) {
	if len(v) == 0 {
		return "[]"
	}
	return "[(" + strings.Join(v, "), (") + ")]"
}
