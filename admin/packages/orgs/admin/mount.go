package admin

import (
	"net/http"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/google/uuid"
	"github.com/sunfmin/reflectutils"
	"gorm.io/gorm"
)

// OrganizacaoModel returns the registered Organizacao model builder, or nil when
// the orgs module has not been configured on b.
func OrganizacaoModel(b *presets.Builder) *presets.ModelBuilder {
	return b.GetModelByID(OrganizacaoModelID)
}

// TemplateScope resolves the sharing permission-template scope from the request:
// the organization in the path (for resources nested under an org). Returns nil
// for the organization itself (top-level) or when no org is in the path, so only
// global templates apply. Use it with shared.WithTemplates.
func TemplateScope(ctx *web.EventContext) *uuid.UUID {
	if id, ok := OrgIDFromRequest(ctx.R); ok {
		return &id
	}
	return nil
}

// MountUnderOrg registers child as a resource nested under the organization, so
// its routes live at /orgs/{orgID}/<group>/<child> and its records are scoped to
// the organization taken from the request path (see ScopeToOrg). When group is
// non-empty it becomes the path segment between the org id and the resource
// (e.g. "finance" -> /orgs/{id}/finance/<child>). The child model must have an
// OrganizacaoID uuid field.
func MountUnderOrg(orgMb, child *presets.ModelBuilder, group string) {
	orgMb.AddChild(child)
	if group != "" {
		child.SetMenuGroupName(group)
	}
	ScopeToOrg(child)
}

// ScopeToOrg stamps a model's OrganizacaoID from the organization id in the
// request path on save, and scopes every read to that organization. The model
// must expose an OrganizacaoID uuid field. Use it on resources mounted under the
// organization (see MountUnderOrg).
func ScopeToOrg(mb *presets.ModelBuilder) {
	mb.Editing().PostSetterFunc(func(obj interface{}, ctx *web.EventContext) {
		if id, ok := OrgIDFromRequest(ctx.R); ok {
			reflectutils.Set(obj, "OrganizacaoID", id)
		}
	})

	mb.UpdateDataOperator(func(do presets.DataOperator) presets.DataOperator {
		return do.(*gorm2op.DataOperatorBuilder).
			WithReadCallbacks(func(cb *gorm2op.Callbacks[*gorm2op.DataOperatorBuilder]) {
				cb.Pre(func(state *gorm2op.CallbackState) (err error) {
					if id, ok := OrgIDFromRequest(state.Ctx.R); ok {
						// qualify with the model's table so the filter is
						// unambiguous when the query joins other org-scoped
						// tables (which also have an organizacao_id column).
						col := "organizacao_id"
						if tbl := statementTable(state.DB); tbl != "" {
							col = tbl + `.organizacao_id`
						}
						state.DB = state.DB.Where(col+" = ?", id)
					}
					return nil
				})
			})
	})
}

// statementTable returns the main table name of the query's model, parsing the
// statement when needed (so the org filter can be table-qualified).
func statementTable(db *gorm.DB) string {
	if db.Statement.Table != "" {
		return db.Statement.Table
	}
	if db.Statement.Model != nil {
		if err := db.Statement.Parse(db.Statement.Model); err == nil {
			return db.Statement.Table
		}
	}
	return ""
}

// OrgIDFromRequest extracts the organization UUID from the first parent id in the request path (/orgs/{orgID}/...). Exported so modules mounted under the org (e.g. finance) can resolve the current org, e.g. for the sharing template scope.
// request path (/orgs/{orgID}/...).
func OrgIDFromRequest(r *http.Request) (uuid.UUID, bool) {
	ids := presets.ParentsModelID(r)
	if len(ids) == 0 {
		return uuid.Nil, false
	}
	switch v := ids.First().Value().(type) {
	case uuid.UUID:
		return v, true
	case string:
		if u, err := uuid.Parse(v); err == nil {
			return u, true
		}
	}
	return uuid.Nil, false
}
