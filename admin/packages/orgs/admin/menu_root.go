package admin

import (
	"github.com/go-rvq/rvq/admin/packages/orgs/models"
	"github.com/go-rvq/rvq/admin/packages/user"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// The organization in the side menu: its record chosen in the menu, as any
// model's (presets MenuRecordSelector) — its item the group of the selector
// (the chip of the organization chosen, its clear; its resources under it,
// in their menu groups) and of the listing of the organizations —, at the
// top of the menu, for a user who may list them: with none yet too (the
// chip "Select", the listing where the first is made). The one chosen is the
// one of the page, else the one kept in the browser, else the only one of
// the user, else their default (models.UserDefault).

// orgMenuRoot is the item of the organization at the top of the menu
// (presets.MenuRootFunc).
func orgMenuRoot(b *presets.Builder, db *gorm.DB) presets.MenuRootFunc {
	return func(ctx *web.EventContext) (nodes []*presets.MenuNode, opened []string, active string) {
		mb := OrganizacaoModel(b)
		u := user.GetCurrentUser(ctx.R)
		if mb == nil || u == nil || u.Anonymous() {
			return
		}
		// as the menu does with any model: the one who may not list them has
		// no item
		if mb.Permissioner().ReqLister(ctx.R).Denied() {
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
		mb.MenuRecordSelectorConfig().DefaultID = func(ctx *web.EventContext) string {
			if u := user.GetCurrentUser(ctx.R); u != nil && !u.Anonymous() {
				if id := models.DefaultOf(db.WithContext(ctx.R.Context()), u.GetID()); id != uuid.Nil {
					return id.String()
				}
			}
			return ""
		}
		mb.MenuRecordSelector(true).InMenu(false)
	}
	b.AddMenuRootFunc(orgMenuRoot(b, db))
}
