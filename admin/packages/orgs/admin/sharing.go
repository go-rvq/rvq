package admin

import (
	"context"
	"net/http"

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
	return nil
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
	return err == nil && e.Visible()
}
