package gorm2op

import (
	"strings"

	"github.com/qor5/admin/v3/model"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Field struct {
	s *Schema
	f *schema.Field
}

func (f *Field) Name() string {
	return f.f.Name
}

func (f *Field) DBName() string {
	return f.f.DBName
}

func (f *Field) QuotedDBName() string {
	return f.s.quote(f.f.DBName)
}

func (f *Field) String() string {
	return f.f.Name
}

func (f *Field) GormField() *schema.Field {
	return f.f
}

func (f *Field) FullDBName() string {
	return f.s.Table() + "." + f.f.DBName
}

func (f *Field) QuotedFullDBName() string {
	return f.s.QuotedTable() + "." + f.QuotedDBName()
}

type Schema struct {
	quote func(s string) string
	s     *schema.Schema
	model any
}

func (s *Schema) Fields() (fields model.Fields) {
	fields = make([]model.Field, len(s.s.Fields))
	for i, field := range s.s.Fields {
		fields[i] = &Field{s, field}
	}
	return
}

func (s *Schema) PrimaryFields() (fields model.Fields) {
	fields = make([]model.Field, len(s.s.PrimaryFields))
	for i, field := range s.s.PrimaryFields {
		fields[i] = &Field{s, field}
	}
	return
}

func (s *Schema) FieldsByName(name ...string) (fields model.Fields) {
	fields = make([]model.Field, len(name))
	for i, name := range name {
		f := s.s.FieldsByName[name]
		if f == nil {
			return nil
		}
		fields[i] = &Field{s, f}
	}
	return
}

func (s *Schema) FieldByName(name string) model.Field {
	f := s.s.FieldsByName[name]
	if f == nil {
		return nil
	}
	return &Field{s, f}
}

func (s *Schema) Table() string {
	return s.s.Table
}

func (s *Schema) String() string {
	return s.s.String()
}

func (s *Schema) QuotedTable() string {
	return s.quote(s.s.Table)
}

func (s *Schema) Quote(v string) string {
	return s.quote(v)
}

func (s *Schema) Model() any {
	return s.model
}

func (s *DataOperatorBuilder) Schema(model any) (ms model.Schema, _ error) {
	return NewSchema(s.db, model), nil
}
func NewSchema(db *gorm.DB, model any) *Schema {
	stmt := &gorm.Statement{DB: db.Model(model)}
	stmt.Parse(model)
	return &Schema{
		func(v string) string {
			var w strings.Builder
			db.Dialector.QuoteTo(&w, v)
			return w.String()
		},
		stmt.Schema,
		model,
	}
}
