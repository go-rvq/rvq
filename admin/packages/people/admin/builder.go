package admin

import (
	"strings"

	orgsadmin "github.com/go-rvq/rvq/admin/packages/orgs/admin"
	"github.com/go-rvq/rvq/admin/packages/people/messages"
	"github.com/go-rvq/rvq/admin/packages/people/models"
	validators "github.com/go-rvq/rvq/admin/packages/validators/models"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"gorm.io/gorm"
)

// Builder mounts the people package (the Person registry) and holds the registry
// of document validators, one per document type. It implements presets.Plugin,
// so it is applied with pb.Use(builder).
//
// The CPF and CNPJ validators are registered by default; an application can
// register validators for other document types with RegisterDocumentValidator
// (for example a passport or a national id), routing that document type through
// its own GAD validator.
type Builder struct {
	db            *gorm.DB
	docValidators map[models.DocumentType]*validators.Validator

	// OrganizationBuilder is the orgs plugin the Person resource is nested under.
	// When nil a default one is used. It is ensured (installed) on Install.
	OrganizationBuilder *orgsadmin.Builder
}

var _ presets.Plugin = (*Builder)(nil)

// NewBuilder returns a people Builder with the default CPF/CNPJ document
// validators registered.
func NewBuilder(db *gorm.DB) *Builder {
	b := &Builder{db: db, docValidators: map[models.DocumentType]*validators.Validator{}}
	for _, v := range models.DocumentValidators() {
		switch v.Name {
		case models.CPFValidatorName:
			b.docValidators[models.DocumentCPF] = v
		case models.CNPJValidatorName:
			b.docValidators[models.DocumentCNPJ] = v
		}
	}
	return b
}

// SetOrganizationBuilder sets the orgs plugin the Person resource is nested
// under. Returns the builder for chaining.
func (b *Builder) SetOrganizationBuilder(ob *orgsadmin.Builder) *Builder {
	b.OrganizationBuilder = ob
	return b
}

// RegisterDocumentValidator registers (or overrides) the validator used for a
// document type. It is seeded on Install and applied to the Person Document
// field. Returns the builder for chaining.
func (b *Builder) RegisterDocumentValidator(docType models.DocumentType, v *validators.Validator) *Builder {
	b.docValidators[docType] = v
	return b
}

// Validators returns the registered document validators.
func (b *Builder) Validators() []*validators.Validator {
	out := make([]*validators.Validator, 0, len(b.docValidators))
	for _, v := range b.docValidators {
		out = append(out, v)
	}
	return out
}

// Install implements presets.Plugin: it migrates the people table, seeds the
// document validators, ensures the orgs package is mounted, registers the i18n
// messages and enum selects, and registers the (org-nested) Person resource.
func (b *Builder) Install(pb *presets.Builder) error {
	if err := models.AutoMigrate(b.db); err != nil {
		return err
	}
	if err := validators.AutoMigrate(b.db); err != nil {
		return err
	}
	for _, v := range b.docValidators {
		if err := validators.Seed(b.db, v); err != nil {
			return err
		}
	}
	// ensure the organization is mounted; Use is idempotent and also binds the
	// org builder so MountUnderOrg works.
	if b.OrganizationBuilder == nil {
		b.OrganizationBuilder = orgsadmin.New(b.db)
	}
	pb.Use(b.OrganizationBuilder)
	messages.Register(pb.I18n())
	registerEnumSelects(pb)
	b.configurePerson(pb)
	return nil
}

// documentValidator validates the Document field against the validator
// registered for the record's document type. Empty documents and unregistered
// types are accepted. The engine produces a language-translated error.
func (b *Builder) documentValidator() presets.FieldValidatorFunc {
	return func(field *presets.FieldContext) (verr web.ValidationErrors) {
		doc, _ := field.Value().(string)
		if strings.TrimSpace(doc) == "" {
			return
		}
		p, ok := field.Obj.(*models.Person)
		if !ok {
			return
		}
		v := b.docValidators[p.DocumentType]
		if v == nil {
			return
		}
		if err := validators.Check(b.db, v.Name, doc, requestLang(field.EventContext)); err != nil {
			verr.FieldError(field.Name, err.Error())
		}
		return
	}
}
