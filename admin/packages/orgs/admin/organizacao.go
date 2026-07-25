package admin

import (
	"github.com/go-rvq/rvq/admin/packages/orgs/models"
	"github.com/go-rvq/rvq/admin/packages/shared"
	"github.com/go-rvq/rvq/admin/packages/user"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrganizacaoModelID is the presets model id of the Organizacao resource.
const OrganizacaoModelID = "organizacoes"

// OrganizacaoURIName is the URI segment of the organization in nested paths, so
// domain modules mounted under an organization live at /orgs/{id}/<module>/...
const OrganizacaoURIName = "orgs"

// configureOrganizacao registers the Organizacao CRUD. Records are owner-scoped:
// the current user only sees the organizations they own. New organizations are
// stamped with the current user as owner. Broader access (shared organizations)
// is added by the sharing manager.
func configureOrganizacao(b *presets.Builder, db *gorm.DB) {
	mb := Model(b, &models.Organizacao{}, presets.ModelWithID(OrganizacaoModelID)).
		MenuIcon("mdi-domain")
	mb.SetUriName(OrganizacaoURIName)

	mb.Listing("Nome", "Descricao").
		SearchColumns("nome", "descricao").
		OrderBy("nome")

	ed := mb.Editing("Nome", "Descricao")
	ed.Field("Nome").Required(true)
	ed.SetterFunc(func(obj interface{}, ctx *web.EventContext) {
		o := obj.(*models.Organizacao)
		if o.OwnerID == uuid.Nil {
			if u := user.GetCurrentUser(ctx.R); u != nil {
				o.OwnerID = u.GetID()
			}
		}
	})

	mb.Detailing("Nome", "Descricao")

	// the owner can share the organization with other users (SPEC: sharing is
	// mounted on Organization and Project). Verbs come from the sharing
	// permission templates (WithTemplates).
	shared.Install(mb, db, shared.WithTemplates(TemplateScope))

	// owner-scope every read (listing and detail) to the current user's orgs.
	mb.UpdateDataOperator(func(do presets.DataOperator) presets.DataOperator {
		return do.(*gorm2op.DataOperatorBuilder).
			WithReadCallbacks(func(cb *gorm2op.Callbacks[*gorm2op.DataOperatorBuilder]) {
				cb.Pre(func(state *gorm2op.CallbackState) (err error) {
					if u := user.GetCurrentUser(state.Ctx.R); u != nil {
						state.DB = state.DB.Where("owner_id = ?", u.GetID())
					}
					return nil
				})
			})
	})
}
