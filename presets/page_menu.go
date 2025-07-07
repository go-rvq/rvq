package presets

import (
	"fmt"
	"strings"

	"github.com/qor5/web/v3"
	. "github.com/qor5/x/v3/ui/vuetify"
	h "github.com/theplant/htmlgo"
)

func (m *PageBuilder) menuItem(ctx *web.EventContext, isSub bool) (r h.HTMLComponent) {
	menuIcon := m.menuIcon
	href := m.pageHandler.path
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

	item.Href(href)
	if strings.HasPrefix(href, "/") {
		funcStr := fmt.Sprintf(`(e) => {
	if (e.metaKey || e.ctrlKey) { return; }
	e.stopPropagation();
	e.preventDefault();
	%s;
}
`, web.Plaid().PushStateURL(href).Go())
		item.Attr("@click", funcStr)
	}
	// if b.isMenuItemActive(ctx, m) {
	//	item = item.Class("v-list-item--active text-primary")
	// }
	return item
}

func (m *PageBuilder) isMenuItemActive(ctx *web.EventContext) bool {
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
