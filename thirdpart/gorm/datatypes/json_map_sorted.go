package datatypes

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type OrderedValue[T any] struct {
	Position int
	Value    T
}

type RawNullJSONOrderedMap[T any] struct {
	Values []*OrderedValue[T]
	Keys   []string
}

// NullJSONOrderedMap defined JSON data type, need to implements driver.Valuer, sql.Scanner interface
type NullJSONOrderedMap[T any] struct {
	M    map[string]*OrderedValue[T]
	keys []string
}

func (m *NullJSONOrderedMap[T]) Set(key string, value T) {
	m.M[key].Value = value
}

func (m *NullJSONOrderedMap[T]) SetMap(mv map[string]*OrderedValue[T]) {
	m.M = mv
	m.keys = make([]string, len(mv))

	var i int
	for key := range mv {
		m.keys = append(m.keys, key)
		i++
	}

	m.Sort()
}

func (m *NullJSONOrderedMap[T]) SetValues(keys []string, values []T) {
	for i, value := range values {
		m.M[keys[i]].Value = value
	}
}

func (m *NullJSONOrderedMap[T]) Clear() {
	m.M = nil
	m.keys = nil
}

func (m *NullJSONOrderedMap[T]) IsZero() bool {
	return len(m.M) == 0
}

// Value return json value, implement driver.Valuer interface
func (m NullJSONOrderedMap[T]) Value() (driver.Value, error) {
	if m.IsZero() {
		return nil, nil
	}

	ba, err := m.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *NullJSONOrderedMap[T]) Scan(val interface{}) error {
	var b []byte

	if val != nil {
		switch v := val.(type) {
		case []byte:
			b = v
		case string:
			b = []byte(v)
		default:
			return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", val))
		}
	}

	return m.UnmarshalJSON(b)
}

func (m *NullJSONOrderedMap[T]) Sort() {
	sort.Slice(m.keys, func(i, j int) bool {
		return m.M[m.keys[i]].Position < m.M[m.keys[j]].Position
	})
}

func (m *NullJSONOrderedMap[T]) Raw() (r *RawNullJSONOrderedMap[T]) {
	r = &RawNullJSONOrderedMap[T]{}
	var i int
	r.Keys = m.keys
	r.Values = make([]*OrderedValue[T], len(m.M))

	for _, key := range r.Keys {
		r.Values[i] = m.M[key]
		i++
	}
	return
}

// MarshalJSON to output non base64 encoded []byte
func (m *NullJSONOrderedMap[T]) MarshalJSON() ([]byte, error) {
	if m.IsZero() {
		return []byte("null"), nil
	}

	m.Sort()
	return json.Marshal(m.Raw())
}

// UnmarshalJSON to deserialize []byte
func (m *NullJSONOrderedMap[T]) UnmarshalJSON(b []byte) (err error) {
	if len(b) == 0 || bytes.Equal(b, []byte("null")) {
		m.Clear()
		return
	}

	var (
		r   RawNullJSONOrderedMap[T]
		dec = json.NewDecoder(bytes.NewReader(b))
	)

	dec.UseNumber()
	if err = json.Unmarshal(b, &r); err != nil {
		return
	}

	m.keys = r.Keys
	m.M = make(map[string]*OrderedValue[T], len(m.M))

	for i, key := range m.keys {
		m.M[key] = r.Values[i]
	}

	return err
}

// GormDataType gorm common data type
func (m NullJSONOrderedMap[T]) GormDataType() string {
	return "null_json_sorted_map"
}

// GormDBDataType gorm db data type
func (NullJSONOrderedMap[T]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	case "sqlserver":
		return "NVARCHAR(MAX)"
	}
	return ""
}

func (m NullJSONOrderedMap[T]) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if m.IsZero() {
		return gorm.Expr("NULL")
	}
	data, _ := m.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}
