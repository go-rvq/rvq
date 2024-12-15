package model

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

type Schema interface {
	Model() any
	Table() string
	QuotedTable() string
	Fields() Fields
	PrimaryFields() Fields
	FieldsByName(f ...string) Fields
	FieldByName(name string) Field
}

func HasPrimaryFields(s Schema) bool {
	return len(s.PrimaryFields()) > 0
}
