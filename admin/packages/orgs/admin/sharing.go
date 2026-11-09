package admin

import (
	"context"
	"net/http"
	"slices"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/packages/orgs/messages"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"

	"github.com/go-rvq/rvq/admin/packages/orgs/models"
	"github.com/go-rvq/rvq/admin/packages/shared"
	"github.com/go-rvq/rvq/admin/packages/user"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ConfirmPath is the page of the admin where the owner confirms a request of
// sharing (shared.Service.InstallConfirm).
const ConfirmPath = "/shared/confirm"

// sharing is the sharing of the organizations (Builder.Sharing): nil, an
// organization is its owner's alone.
var sharing *shared.Service

// Sharing shares the organizations by s: they — and what is under them —
// seen by whom they are shared with (their roles), not by their owners
// alone; the actions of sharing in their detail, the page of the
// confirmations, the invites of the user. Before Install.
func (b *Builder) Sharing(s *shared.Service) *Builder {
	b.sharing = s
	sharing = s
	return b
}

// SharingService is the service the organizations are shared by (nil: an
// organization is its owner's alone): the resources under them may be
// shared by it too, as their children (shared.Resource.Parent).
func (b *Builder) SharingService() *shared.Service { return b.sharing }

// OrgRoles are the roles of the sharing of an organization.
var OrgRoles = []shared.RoleSeed{
	{Name: "Viewer", Actions: shared.VerbView},
	{Name: "Editor", Actions: append(append([]string{}, shared.VerbView...), shared.VerbEdit...)},
	{Name: "Manager", Actions: append(append(append([]string{}, shared.VerbView...), shared.VerbEdit...), shared.VerbDelete...)},
	{Name: "Administrator", Actions: append(append(append(append([]string{}, shared.VerbView...), shared.VerbEdit...),
		shared.VerbDelete...), shared.PermShare)},
}

// OrgResource is the organization as a resource of the sharing.
func OrgResource() shared.Resource {
	return shared.Resource{Name: OrganizacaoModelID, Table: "org_organizacoes", OwnerColumn: "owner_id",
		Label: func(ctx context.Context, db *gorm.DB, id string) string {
			var o models.Organizacao
			db.Limit(1).Find(&o, "id = ?", id)
			return o.Nome
		}}
}

// installSharing shares the organizations by s.
func installSharing(pb *presets.Builder, s *shared.Service) error {
	mb := OrganizacaoModel(pb)
	if err := s.Install(mb, shared.ModelOptions{Resource: OrgResource(), Roles: OrgRoles, ReadScope: true}); err != nil {
		return err
	}
	s.InstallConfirm(pb, ConfirmPath)
	s.InstallInvites(pb)
	s.ListenGroups()
	configureOwnerShown(mb, s)
	// read as the parent of a record shared under it (a project): its name,
	// description, owner
	s.Support(mb, shared.SupportOptions{Fields: []string{"Nome", "Descricao", "Owner", "MyDefault"}})
	return nil
}

// configureOwnerShown: an organization of another, shared with the user, says
// so — "shared" and its owner — beside its name in the listing; the detail
// says its owner.
func configureOwnerShown(mb *presets.ModelBuilder, s *shared.Service) {
	mb.Listing().Field("Nome").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		o := field.Obj.(*models.Organizacao)
		td := h.Td(h.Text(o.Nome))
		if uid := currentUserID(ctx); uid != uuid.Nil && o.OwnerID.UUID() != uid {
			m := messages.Get(ctx.Context())
			td.AppendChildren(VChip(h.Text(m.SharedBy+" "+s.UserName(ctx.R.Context(), o.OwnerID.String()))).
				Size(SizeXSmall).Variant(VariantTonal).Color(ColorPrimary).Class("ml-2").
				Attr("data-org-shared", true))
		}
		return td
	})
	dt := mb.Detailing()
	dt.Field("Owner").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		o := field.Obj.(*models.Organizacao)
		field.ValueOverride = s.UserName(ctx.R.Context(), o.OwnerID.String())
		return presets.ReadonlyComponentFunc(field, ctx)
	})
	dt.Only("Nome", "Descricao", "Owner")
}

// accessibleOrgs is db limited to the organizations of the user of r: the
// ones they see (shared), else the ones they own.
func accessibleOrgs(r *http.Request, db *gorm.DB) *gorm.DB {
	u := user.GetCurrentUser(r)
	if u == nil || u.Anonymous() {
		return db.Where("1 = 0")
	}
	if sharing == nil {
		return db.Where("owner_id = ?", u.GetID())
	}
	q, err := sharing.Scope(r.Context(), db, u.GetID().String(), OrganizacaoModelID)
	if err != nil {
		return db.Where("1 = 0")
	}
	return q
}

// orgVisible says whether the user of r sees the organization id.
func orgVisible(r *http.Request, id uuid.UUID) bool {
	if sharing == nil {
		return true
	}
	u := user.GetCurrentUser(r)
	if u == nil || u.Anonymous() {
		return false
	}
	e, err := sharing.RoleOf(r.Context(), u.GetID().String(), OrganizacaoModelID, id.String())
	if err == nil && e.Visible() {
		return true
	}
	// read as support: the parent of a record under it shared with the user
	// (what of it they read, the permissions and the supporting models say)
	ids, err := sharing.AccessibleIDs(r.Context(), u.GetID().String(), OrganizacaoModelID)
	return err == nil && slices.Contains(ids, id.String())
}
