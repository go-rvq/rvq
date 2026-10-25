package presets

import (
	"slices"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

// The main menu has two looks: whole (the tree, CreateMenus) and its icons
// alone (the rail: vars.navRail, kept in the browser, NavRailStorageKey). In
// the rail each item of the root is its icon: a group opens, under the mouse,
// a floating menu of its items (its sub-groups within it, titled); an item,
// its page, its name the tooltip. Both come of the same nodes: the same items,
// the same permissions, the same one active.

// NavRailStorageKey is where the browser keeps the look of the main menu:
// "1", its icons alone.
const NavRailStorageKey = "rvq.presets.navRail"

// navRailWidth is the width of the main menu as its icons alone.
const navRailWidth = 64

// railMenu is the main menu as its icons alone.
func railMenu(nodes []*MenuNode, activated []string) h.HTMLComponent {
	var items h.HTMLComponents
	for _, n := range nodes {
		active := len(activated) > 0 && n.has(activated[0])
		if len(n.Children) == 0 {
			items = append(items, railLeaf(n, active))
			continue
		}
		items = append(items, railGroup(n, activated, active))
	}
	return h.Div(items...).Class("rvq-nav-rail d-flex flex-column align-center ga-1 py-2")
}

// railIcon is the icon of a node (its picture, when it has one).
func railIcon(n *MenuNode) h.HTMLComponent {
	if img, _ := n.Props["prependImage"].(string); img != "" {
		return h.Img(img).Attr("alt", "").Style("width:22px;height:22px;object-fit:contain")
	}
	icon, _ := n.Props["prependIcon"].(string)
	if icon == "" {
		icon = "mdi-circle-small"
	}
	return VIcon(icon)
}

func railButton(n *MenuNode, active bool) *VBtnBuilder {
	btn := VBtn("").Icon(true).Variant(VariantText).
		Attr("aria-label", n.Title).
		Class("rvq-nav-rail__btn").
		Children(railIcon(n))
	if active {
		btn.Color(ColorPrimary).Variant(VariantTonal)
	}
	return btn
}

// railLeaf: an item with a page: its icon, its name the tooltip.
func railLeaf(n *MenuNode, active bool) h.HTMLComponent {
	href, _ := n.Props["href"].(string)
	return VTooltip(
		web.Slot(
			railButton(n, active).Href(href).Attr("v-bind", "props"),
		).Name("activator").Scope("{ props }"),
	).Text(n.Title).Location(LocationEnd)
}

// railGroup: a group: its icon opens, under the mouse, its items.
func railGroup(n *MenuNode, activated []string, active bool) h.HTMLComponent {
	return VMenu(
		web.Slot(
			railButton(n, active).Attr("v-bind", "props"),
		).Name("activator").Scope("{ props }"),
		VCard(
			VCardTitle(h.Text(n.Title)).Class("text-subtitle-2 font-weight-bold py-2"),
			VDivider(),
			VList(railItems(n.Children, activated, 0)...).Density(DensityCompact),
		).Class("rvq-nav-rail__flyout").MinWidth(220).MaxWidth(320),
	).
		OpenOnHover(true).
		Location(LocationEnd).
		Offset(8).
		Attr(":close-delay", "150").
		Attr(":open-delay", "60")
}

// railItems are the items of a group in its floating menu: a page, its
// link; a sub-group, its title and its items under it.
func railItems(nodes []*MenuNode, activated []string, depth int) (items h.HTMLComponents) {
	for _, n := range nodes {
		href, _ := n.Props["href"].(string)
		if len(n.Children) > 0 && href == "" {
			items = append(items, VListSubheader(h.Text(n.Title)).Class("rvq-nav-rail__sub"))
			items = append(items, railItems(n.Children, activated, depth+1)...)
			continue
		}
		item := VListItem(
			web.Slot(railIcon(n)).Name(VSlotPrepend),
			VListItemTitle(h.Text(n.Title)),
		).Href(href).Active(slices.Contains(activated, n.Value))
		if depth > 0 {
			item.Class("rvq-nav-rail__nested")
		}
		items = append(items, item)
		// a page with its own pages (its children): them under it
		if len(n.Children) > 0 {
			items = append(items, railItems(n.Children, activated, depth+1)...)
		}
	}
	return
}

// navRailToggle is the button at the foot of the main menu: its icons alone,
// or whole again.
func navRailToggle(m *Messages) h.HTMLComponent {
	save := `try { window.localStorage.setItem(` + h.JSONString(NavRailStorageKey) + `, vars.navRail ? "1" : "") } catch (e) {}`
	return h.Div(
		VBtn(m.CollapseMenu).
			PrependIcon("mdi-chevron-double-left").
			Variant(VariantText).
			Block(true).
			Class("text-none justify-start").
			Attr("v-if", "!vars.navRail").
			Attr("@click", "vars.navRail = true; "+save),
		VTooltip(
			web.Slot(
				VBtn("").Icon("mdi-chevron-double-right").Variant(VariantText).
					Attr("aria-label", m.ExpandMenu).
					Attr("v-bind", "props").
					Attr("@click", "vars.navRail = false; "+save),
			).Name("activator").Scope("{ props }"),
		).Text(m.ExpandMenu).Location(LocationEnd).Attr("v-else", true),
	).Class("rvq-nav-rail__toggle d-flex justify-center pa-2 border-t").Attr("v-if", "!vars.navMobile")
}
