package admin

import (
	"strings"

	"github.com/go-rvq/rvq/admin/packages/people/models"
	validators "github.com/go-rvq/rvq/admin/packages/validators/models"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"gorm.io/gorm"
)

// configurePerson registers the people registry resource: a Person (identity)
// owned by the current organization. Its financial receiving data lives in the
// finance package and is mounted as a nested sub-resource by finance. The
// document is validated through the builder's document validators
// (FieldBuilder.Validators), routing by document type → GAD validator.
func (b *Builder) configurePerson(pb *presets.Builder) *presets.ModelBuilder {
	mb := NewModel(pb, &models.Person{}, presets.ModelWithID("persons")).
		MenuIcon("mdi-account-group")
	b.OrganizationBuilder.MountUnderOrg(mb, MenuGroupID)

	mb.Listing("Name", "DocumentType", "Document", "Active").
		SearchColumns("name", "document").
		OrderBy("name")

	ed := mb.Editing("Name", "DocumentType", "Document", "Address", "Notes", "Active")
	ed.Field("Name").Required(true)
	ed.Field("Document").Validator(b.documentValidator())

	mb.Detailing("Name", "DocumentType", "Document", "Address", "Notes", "Active")

	configureTrash(mb, b.db, models.Person{})
	return mb
}

// CPFOrCNPJValidator returns a field validator for a free-form document field
// that must be a valid CPF or CNPJ (no separate document-type field). It passes
// when the value is empty or either the CPF or the CNPJ validator accepts it;
// otherwise it reports the CPF validator's (translated) message. Reusable by
// packages with holder/payer document fields (e.g. finance).
func CPFOrCNPJValidator(db *gorm.DB) presets.FieldValidatorFunc {
	return func(field *presets.FieldContext) (verr web.ValidationErrors) {
		doc, _ := field.Value().(string)
		if strings.TrimSpace(doc) == "" {
			return
		}
		lang := requestLang(field.EventContext)
		cpfErr := validators.Check(db, models.CPFValidatorName, doc, lang)
		if cpfErr == nil {
			return
		}
		if validators.Check(db, models.CNPJValidatorName, doc, lang) == nil {
			return
		}
		verr.FieldError(field.Name, cpfErr.Error())
		return
	}
}

// requestLang returns a best-effort BCP-47 language tag for the request (the
// first Accept-Language token), used to translate validator messages. The
// validator falls back to its DefaultLang when the tag is unknown.
func requestLang(ctx *web.EventContext) string {
	if ctx == nil || ctx.R == nil {
		return ""
	}
	al := ctx.R.Header.Get("Accept-Language")
	if al == "" {
		return ""
	}
	if i := strings.IndexAny(al, ",;"); i >= 0 {
		al = al[:i]
	}
	return strings.TrimSpace(al)
}
