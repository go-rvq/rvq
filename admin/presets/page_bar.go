package presets

import (
	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

// The buttons of the bar of the page (the title, its actions — above the
// breadcrumbs): one look for all of them, whatever made them — an action of
// the listing, a bulk action, an action of the detail, the layout's own —:
// text, with their icon before it, the bar's height (fix-layout.css,
// .rvq-page-header); an icon alone, as text too; the one in evidence (a new
// record, the edit), primary, with its title (PageBarPrimaryButton).

// PageBarButton is an action of the bar of the page: its title, its icon
// before it.
func PageBarButton(title, icon string) *VBtnBuilder {
	b := VBtn(title).Variant(VariantText).Class("rvq-page-bar__btn")
	if icon != "" {
		b.PrependIcon(icon)
	}
	return b
}

// PageBarIconButton is an action of the bar of the page shown by its icon
// alone (title: its tooltip, its name read aloud); primary: the one in
// evidence (a new record).
func PageBarIconButton(icon, title string, primary bool) *VBtnBuilder {
	b := VBtn("").Icon(icon).Class("rvq-page-bar__btn")
	if title != "" {
		b.Attr("title", title).Attr("aria-label", title)
	}
	if primary {
		b.Color(ColorPrimary).Variant(VariantFlat)
	} else {
		b.Variant(VariantText)
	}
	return b
}

// PageBarPrimaryButton is the action in evidence of the bar of the page (a
// new record, the edit): its title, its icon before it, primary — the
// height and the icon of the others, not a disc above them.
func PageBarPrimaryButton(title, icon string) *VBtnBuilder {
	return PageBarButton(title, icon).Color(ColorPrimary).Variant(VariantFlat)
}

// pageBarColor is the color of an action given one (not the old default,
// black: the text's own).
func pageBarColor(b *VBtnBuilder, color string) *VBtnBuilder {
	if color != "" && color != "black" {
		b.Color(color)
	}
	return b
}
