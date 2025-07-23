package presets

import (
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

func (m *HttpPageBuilder) menuItem(ctx *web.EventContext, isSub bool) (r h.HTMLComponent) {
	menuIcon := m.menuIcon
	label := m.TTitle(ctx.Context())
	item := VListItem(
		// VRow(
		// 	VCol(h.If(menuIcon != "", VIcon(menuIcon))).Cols(2),
		// 	VCol(h.Text(i18n.T(ctx.R, ModelsI18nModuleKey, m.label))).Attr("style", fmt.Sprintf("white-space: normal; font-weight: %s;font-size: 16px;", fontWeight))),

		h.If(menuIcon != "", web.Slot(VIcon(menuIcon)).Name("prepend")),
		VListItemTitle(
			h.Text(label),
		),
	).Class("rounded-lg").
		Value(label)
	// .ActiveClass("bg-red")
	// Attr("color", "primary")

	item.Href(m.fullPath)
	// if b.isMenuItemActive(ctx, m) {
	//	item = item.Class("v-list-item--active text-primary")
	// }
	return item
}

func (m *HttpPageBuilder) isMenuItemActive(ctx *web.EventContext) bool {
	href := m.pageHandler.path

	path := strings.TrimSuffix(ctx.R.URL.Path, "/")

	if path == "" && href == "/" {
		return true
	}

	if path == href {
		return true
	}

	if href != "/" && strings.HasPrefix(path, href) {
		return true
	}

	return false
}
