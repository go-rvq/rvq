package model

import (
	"reflect"
	"testing"
)

func TestIDStringEscapes(t *testing.T) {
	for _, tc := range []struct {
		values []any
		want   string
	}{
		{[]any{"pt-BR", "admin/packages/user"}, "pt-BR_admin%2Fpackages%2Fuser"},
		{[]any{"a b"}, "a%20b"},
		{[]any{"9f1c2b3e-0000-4000-8000-000000000001"}, "9f1c2b3e-0000-4000-8000-000000000001"},
		{[]any{2, "en-US"}, "2_en-US"},
		{[]any{"i18n:default"}, "i18n:default"},
	} {
		if got := (ID{Values: tc.values}).String(); got != tc.want {
			t.Errorf("%v: %q, want %q", tc.values, got, tc.want)
		}
	}
}

type slugModel struct {
	Lang, Module string
}

func (m *slugModel) PrimarySlug() string { return m.Lang + "_" + m.Module }

type fakeField struct{ name string }

func (f fakeField) Name() string             { return f.name }
func (f fakeField) Type() reflect.Type       { return reflect.TypeOf("") }
func (f fakeField) DBName() string           { return f.name }
func (f fakeField) QuotedDBName() string     { return f.name }
func (f fakeField) FullDBName() string       { return f.name }
func (f fakeField) QuotedFullDBName() string { return f.name }

type fakeSchema struct{ fields Fields }

func (s fakeSchema) Model() any                      { return &slugModel{} }
func (s fakeSchema) Table() string                   { return "t" }
func (s fakeSchema) QuotedTable() string             { return "t" }
func (s fakeSchema) Fields() Fields                  { return s.fields }
func (s fakeSchema) PrimaryFields() Fields           { return s.fields }
func (s fakeSchema) FieldsByName(f ...string) Fields { return s.fields }
func (s fakeSchema) FieldByName(name string) Field   { return nil }

// A model with a slug: the ID is its slug, escaped after the encoding.
func TestIDStringSlugEncoder(t *testing.T) {
	fields := Fields{fakeField{"Lang"}, fakeField{"Module"}}
	id := ID{Values: []any{"pt-BR", "admin/packages/user"}, Fields: fields, Schema: fakeSchema{fields}}
	if got := id.Slug(); got != "pt-BR_admin/packages/user" {
		t.Errorf("Slug %q", got)
	}
	if got := id.String(); got != "pt-BR_admin%2Fpackages%2Fuser" {
		t.Errorf("String %q", got)
	}
}
