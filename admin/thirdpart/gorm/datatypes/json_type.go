package datatypes

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/qor5/web/v3/zeroer"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// NullJSONType give a generic data type for json encoded data.
type NullJSONType[T any] struct {
	Data T
}

func NewJSONType[T any](data T) NullJSONType[T] {
	return NullJSONType[T]{
		Data: data,
	}
}

func (j NullJSONType[T]) IsZero() bool {
	return zeroer.IsNil(j.Data) || zeroer.IsZero(j.Data)
}

// Value return json value, implement driver.Valuer interface
func (j NullJSONType[T]) Value() (driver.Value, error) {
	if zeroer.IsNil(j.Data) || zeroer.IsZero(j.Data) {
		return nil, nil
	}
	return json.Marshal(j.Data)
}

// Scan scan value into NullJSONType[T], implements sql.Scanner interface
func (j *NullJSONType[T]) Scan(value interface{}) error {
	var b []byte
	if value != nil {
		switch v := value.(type) {
		case []byte:
			b = v
		case string:
			b = []byte(v)
		default:
			return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
		}
	}
	if len(b) == 0 || bytes.Equal(b, []byte("null")) {
		*j = NullJSONType[T]{}
		return nil
	}
	return json.Unmarshal(b, &j.Data)
}

// MarshalJSON to output non base64 encoded []byte
func (j NullJSONType[T]) MarshalJSON() ([]byte, error) {
	if j.IsZero() {
		return nil, nil
	}
	return json.Marshal(j.Data)
}

// UnmarshalJSON to deserialize []byte
func (j *NullJSONType[T]) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &j.Data)
}

// GormDataType gorm common data type
func (NullJSONType[T]) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (NullJSONType[T]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

func (js NullJSONType[T]) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if js.IsZero() {
		return gorm.Expr("NULL")
	}

	data, _ := js.MarshalJSON()

	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}

	return gorm.Expr("?", string(data))
}

// NullJSONSlice give a generic data type for json encoded slice data.
type NullJSONSlice[T any] []T

func NewJSONSlice[T any](s []T) NullJSONSlice[T] {
	return NullJSONSlice[T](s)
}

// Value return json value, implement driver.Valuer interface
func (j NullJSONSlice[T]) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan scan value into NullJSONType[T], implements sql.Scanner interface
func (j *NullJSONSlice[T]) Scan(value interface{}) error {
	var b []byte
	if value != nil {
		switch v := value.(type) {
		case []byte:
			b = v
		case string:
			b = []byte(v)
		default:
			return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
		}
	}
	if len(b) == 0 || bytes.Equal(b, []byte("null")) {
		*j = nil
		return nil
	}
	return json.Unmarshal(b, &j)
}

// GormDataType gorm common data type
func (NullJSONSlice[T]) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (NullJSONSlice[T]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

func (j NullJSONSlice[T]) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if len(j) == 0 {
		return gorm.Expr("NULL")
	}

	data, _ := json.Marshal(j)

	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}

	return gorm.Expr("?", string(data))
}

// NullJSONMap defined JSON data type, need to implements driver.Valuer, sql.Scanner interface
type NullJSONMap map[string]interface{}

// Value return json value, implement driver.Valuer interface
func (m NullJSONMap) Value() (driver.Value, error) {
	if len(m) == 0 {
		return nil, nil
	}
	ba, err := m.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *NullJSONMap) Scan(val interface{}) error {
	var b []byte

	if val != nil {
		switch v := val.(type) {
		case []byte:
			b = v
		case string:
			b = []byte(v)
		default:
			return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", val))
		}
	}

	if len(b) == 0 || bytes.Equal(b, []byte("null")) {
		*m = make(NullJSONMap)
		return nil
	}

	t := map[string]interface{}{}
	rd := bytes.NewReader(b)
	decoder := json.NewDecoder(rd)
	decoder.UseNumber()
	err := decoder.Decode(&t)
	*m = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m NullJSONMap) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	t := (map[string]interface{})(m)
	return json.Marshal(t)
}

// UnmarshalJSON to deserialize []byte
func (m *NullJSONMap) UnmarshalJSON(b []byte) error {
	t := map[string]interface{}{}
	err := json.Unmarshal(b, &t)
	*m = t
	return err
}

// GormDataType gorm common data type
func (m NullJSONMap) GormDataType() string {
	return "jsonmap"
}

// GormDBDataType gorm db data type
func (NullJSONMap) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (jm NullJSONMap) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if len(jm) == 0 {
		return gorm.Expr("NULL")
	}
	data, _ := jm.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}
