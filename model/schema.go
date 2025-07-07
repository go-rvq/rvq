package model

import "strconv"

type Schema interface {
	Model() any
	Table() string
	QuotedTable() string
	Fields() Fields
	PrimaryFields() Fields
	FieldsByName(f ...string) Fields
	FieldByName(name string) Field
}

type Field interface {
	Name() string
	DBName() string
	QuotedDBName() string
	FullDBName() string
	QuotedFullDBName() string
}

type Fields []Field

func (f Fields) Len() int {
	return len(f)
}

func (f Fields) First() Field {
	return f[0]
}

func (f Fields) DBNames() (s []string) {
	s = make([]string, len(f))
	for i, field := range f {
		s[i] = field.DBName()
	}
	return
}

func (f Fields) FullDBNames() (s []string) {
	s = make([]string, len(f))
	for i, field := range f {
		s[i] = field.DBName()
	}
	return
}

func (f Fields) QuotedDBNames() (s []string) {
	s = make([]string, len(f))
	for i, field := range f {
		s[i] = field.QuotedDBName()
	}
	return
}

func (f Fields) QuotedFullDBNames() (s []string) {
	s = make([]string, len(f))
	for i, field := range f {
		s[i] = field.QuotedFullDBName()
	}
	return
}

func (f Fields) Names() (s []string) {
	s = make([]string, len(f))
	for i, field := range f {
		s[i] = field.Name()
	}
	return
}

func HasPrimaryFields(s Schema) bool {
	return len(s.PrimaryFields()) > 0
}

type SingleField string

func (s SingleField) Name() string {
	return string(s)
}

func (s SingleField) DBName() string {
	return string(s)
}

func (s SingleField) QuotedDBName() string {
	return strconv.Quote(string(s))
}

func (s SingleField) FullDBName() string {
	return "[qor5/admin/model/schema/single_field]." + string(s)
}

func (s SingleField) QuotedFullDBName() string {
	return "[qor5/admin/model/schema/single_field]." + s.QuotedDBName()
}
