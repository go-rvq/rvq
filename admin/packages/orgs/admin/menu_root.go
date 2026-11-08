package admin

import (
	"github.com/go-rvq/rvq/admin/packages/orgs/models"
	"github.com/go-rvq/rvq/admin/packages/user"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"gorm.io/gorm"
)

// The organization in the side menu: its record chosen in the menu, as any
// model's (presets MenuRecordSelector) — its item the group of the selector
// (the chip of the organization chosen, its clear; its resources under it,
// in their menu groups) and of the listing of the organizations —, at the
// top of the menu, for a user with organizations. The one chosen is the one
// of the page, else the one kept in the browser, else the only one of the
// user, else the default (Organizacao.Padrao).

// orgMenuRoot is the item of the organization at the top of the menu
// (presets.MenuRootFunc).
func orgMenuRoot(b *presets.Builder, db *gorm.DB) presets.MenuRootFunc {
	return func(ctx *web.EventContext) (nodes []*presets.MenuNode, opened []string, active string) {
		mb := OrganizacaoModel(b)
		u := user.GetCurrentUser(ctx.R)
		if mb == nil || u == nil || u.Anonymous() {
			return
		}
		var n int64
		if accessibleOrgs(ctx.R, db.Model(&models.Organizacao{})).Count(&n); n == 0 {
			return
		}
		node, active, opened := mb.MenuNodeOf(ctx, nil)
		return []*presets.MenuNode{node}, opened, active
	}
}

// configureMenuRoot makes the organization a record chosen in the menu, its
// item at the top of the menu (not where the model's would be).
func configureMenuRoot(b *presets.Builder, db *gorm.DB) {
	if mb := OrganizacaoModel(b); mb != nil {
		mb.MenuRecordSelectorConfig().DefaultConditions = presets.SQLConditions{{Query: "org_organizacoes.padrao"}}
		mb.MenuRecordSelector(true).InMenu(false)
	}
	b.AddMenuRootFunc(orgMenuRoot(b, db))
}
