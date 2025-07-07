package web

import "net/url"

type Query url.Values

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (v Query) Get(key string) string {
	vs := v[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

// Set sets the key to value. It replaces any existing
// values.
func (v Query) Set(key, value string) Query {
	v[key] = []string{value}
	return v
}

// SetValid sets the key to value if value isn't blank. It replaces any existing
// values.
func (v Query) SetValid(key, value string) Query {
	if value == "" {
		return v
	}
	v[key] = []string{value}
	return v
}

// Add adds the value to key. It appends to any existing
// values associated with key.
func (v Query) Add(key, value string) Query {
	v[key] = append(v[key], value)
	return v
}

// AddValid adds the value to key if value isn't blank. It appends to any existing
// values associated with key.
func (v Query) AddValid(key, value string) Query {
	v[key] = append(v[key], value)
	return v
}

// Del deletes the values associated with key.
func (v Query) Del(key string) Query {
	delete(v, key)
	return v
}

// Has checks whether a given key is set.
func (v Query) Has(key string) bool {
	_, ok := v[key]
	return ok
}

func (v Query) URLValues() url.Values {
	return url.Values(v)
}
