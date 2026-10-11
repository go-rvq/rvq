package integration

import (
	"fmt"
	"net/http"

	"github.com/gad-lang/gad"
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/fields/schemaform"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/vue"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// A fixture for the checks of a schema form in the admin: a singleton whose
// Options are a form of steps — a contact (a phone of a mask), a project
// (rooms of limits, a word a `[validation=…]` function refuses) — kept as
// YAML.

type SFQuote struct {
	ID      uint
	Options string
}

// SchemaformQuoteURI is the singleton's.
const SchemaformQuoteURI = "quote"

// SchemaformQuoteCode is the class of its Options.
const SchemaformQuoteCode = `
[ordered]
class Contact {
	name str
	[mask="(###) ###-####"] phone str
}
[ordered]
class Project {
	[min=1, max=10] rooms int
	[validation=func(v) => v == "bad" ? "not that word" : nil] word? str
}
[ordered]
class Form {
	[steps]
	{
		[label="Contact"]
		contact Contact
		[label="Project"]
		project Project
	}
}
return Form`

// SchemaformQuoteSchema is the schema of SchemaformQuoteCode, read with its VM
// (its validations called on save).
func SchemaformQuoteSchema(b *schemaform.Builder) (*schemaform.Schema, error) {
	bi := b.Builtins()
	res, err := gad.Compile(gad.NewSymbolTable(bi.NameSet), []byte(SchemaformQuoteCode), gad.CompileOptions{})
	if err != nil {
		return nil, err
	}
	vm := gad.NewVM(bi.Build(), res.Bytecode)
	ret, err := vm.Run()
	if err != nil {
		return nil, err
	}
	cls, ok := ret.(*gad.Class)
	if !ok {
		return nil, fmt.Errorf("not a class: %T", ret)
	}
	return b.Class(vm, cls)
}

// NewSchemaformApp registers the singleton, its Options a schema form.
func NewSchemaformApp(db *gorm.DB) (*presets.Builder, error) {
	p := presets.New(i18n.New()).URIPrefix("/admin")
	p.Permission(perm.New().AllowAll())
	p.DataOperator(gorm2op.DataOperator(db))

	b := schemaform.New()
	schema, err := SchemaformQuoteSchema(b)
	if err != nil {
		return nil, err
	}
	m := p.Model(&SFQuote{}).URIName(SchemaformQuoteURI).Singleton(true)
	ed := m.Editing("Options")
	ed.Field("Options").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		value, _ := b.DecodeValue(field.Obj.(*SFQuote).Options)
		return vue.UserComponent(b.ComponentFunc(schema)(field, ctx)).Assign("form", field.FormKey, schema.Filled(value))
	}).SetterFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) error {
		value, err := b.DecodeForm(ctx, schema, ctx.R.Form, field.FormKey)
		if err != nil {
			return err
		}
		obj.(*SFQuote).Options, err = b.EncodeValue(value)
		return err
	})
	return p, nil
}

// NewSchemaformSeededHandler builds NewSchemaformApp with its one record.
func NewSchemaformSeededHandler() (http.Handler, error) {
	db, err := gorm.Open(sqlite.Open(nextMemName()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&SFQuote{}); err != nil {
		return nil, err
	}
	if err = db.Create(&SFQuote{ID: 1}).Error; err != nil {
		return nil, err
	}
	return NewSchemaformApp(db)
}
