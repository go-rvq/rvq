package js

import (
	"bytes"
	"encoding/json"
	"sort"
	"strconv"
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
			b, _ := json.Marshal(k)
			out.Write(b)
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
				if b, err = json.Marshal(v); err == nil {
					out.Write(b)
				}
			}
		}
	)

	for k := range o {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	out.WriteString("{")

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

	out.WriteString("}")

	return out.Bytes(), nil
}

func (o Object) String() string {
	if len(o) == 0 {
		return "{}"
	}

	var (
		out   bytes.Buffer
		keys  []string
		write = func(k string) {
			out.WriteString(strconv.Quote(k))
			out.WriteString(": ")
			v := o[k]
			switch t := v.(type) {
			case Raw:
				out.Write([]byte(t))
			case []byte:
				out.Write(t)
			case RawSlice:
				out.WriteString(t.String())
			default:
				if vb, err := json.Marshal(v); err == nil {
					out.Write(vb)
				} else {
					panic(err)
				}
			}
		}
	)

	for k := range o {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	out.WriteString("{")

	for _, k := range keys[:len(keys)-1] {
		write(k)
		out.WriteByte(',')
	}

	write(keys[len(keys)-1])

	out.WriteString("}")

	return out.String()
}

type RawSlice []string

func (v RawSlice) String() (s string) {
	if len(v) == 0 {
		return "[]"
	}
	return "[(" + strings.Join(v, "), (") + ")]"
}
