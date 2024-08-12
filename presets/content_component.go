package presets

import (
	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	. "github.com/qor5/x/v3/ui/vuetify"
	h "github.com/theplant/htmlgo"
)

func MenuBtn() h.HTMLComponent {
	return VBtn("").
		Variant(VariantFlat).
		Icon(true).
		Density("compact").
		Children(
			VIcon("mdi-menu"),
		).Attr("@click.menu", "locals.menu = !locals.menu")
}

func FullScreenBtn() h.HTMLComponent {
	return h.HTMLComponents{
		VBtn("").Icon(true).
			Children(VIcon("mdi-fullscreen")).
			Size(SizeLarge).
			Density(DensityCompact).
			Attr(
				"v-if", "!closer.fullscreen",
				"@click", "closer.fullscreen = !closer.fullscreen",
			),
		VBtn("").Icon(true).
			Children(VIcon("mdi-fullscreen-exit")).
			Size(SizeLarge).
			Density(DensityCompact).
			Attr(
				"v-if", "closer.fullscreen",
				"@click", "closer.fullscreen = !closer.fullscreen",
			),
	}
}

func CloseBtn() h.HTMLComponent {
	return VBtn("").Icon(true).
		Children(VIcon("mdi-close")).
		Size(SizeLarge).
		Attr("@click.stop", "closer.show = false")
}

type ContentComponentBuilderOverlay struct {
	Mode               actions.OverlayMode
	Width              string
	Height             string
	MaxWidth           string
	MaxHeight          string
	MinWidth           string
	MinHeight          string
	FullscreenDisabled bool
	CloseDisabled      bool
	Scrollable         bool
}

type ContentComponentBuilder struct {
	Title           string
	TopLeftActions  h.HTMLComponents
	TopRightActions h.HTMLComponents
	Tabs            h.HTMLComponents
	PrimaryAction   h.HTMLComponent
	BottomActions   h.HTMLComponents
	Overlay         *ContentComponentBuilderOverlay
	Menu            h.HTMLComponents
	Body            h.HTMLComponent
	Scope           *web.ScopeBuilder
}

func (b *ContentComponentBuilder) scoped(comp h.HTMLComponent) h.HTMLComponent {
	if b.Scope == nil {
		return comp
	}
	b.Scope.Children(comp)
	return b.Scope
}

func (b *ContentComponentBuilder) BuildOverlay() h.HTMLComponent {
	var (
		header           h.HTMLComponent
		headerComponents = b.TopLeftActions
		headerRight      = b.TopRightActions
		bottomActions    h.HTMLComponent
	)

	if len(b.Menu) > 0 {
		headerComponents = append(h.HTMLComponents{MenuBtn()}, headerComponents...)
	}

	if b.Title != "" {
		headerComponents = append(headerComponents, VToolbarTitle("").
			Children(h.Text(b.Title)))
	}

	if b.PrimaryAction != nil {
		headerRight = append(headerRight, b.PrimaryAction)
	}

	if !b.Overlay.FullscreenDisabled {
		headerRight = append(headerRight, FullScreenBtn())
	}

	if !b.Overlay.CloseDisabled {
		headerRight = append(headerRight, CloseBtn())
	}

	if len(headerComponents) > 0 || len(headerRight) > 0 {
		headerComponents = append(headerComponents, VSpacer())
		headerComponents = append(headerComponents, headerRight...)
	}

	if len(headerComponents) > 0 {
		header = VToolbar(headerComponents...).Color("white").Elevation(0).Density(DensityCompact)
	}

	headerComponents = append(headerComponents, b.PrimaryAction)

	if len(b.BottomActions) > 0 {
		bottomActions = VCardActions(b.BottomActions...)
	}

	var body h.HTMLComponent = VMain(
		b.Tabs,
		b.Body,
	)

	if len(b.Menu) > 0 {
		body = h.HTMLComponents{
			VNavigationDrawer(b.Menu).
				// Attr("@input", "plaidForm.dirty && vars.presetsRightDrawer == false && !confirm('You have unsaved changes on this form. If you close it, you will lose all unsaved changes. Are you sure you want to close it?') ? vars.presetsRightDrawer = true: vars.presetsRightDrawer = $event"). // remove because drawer plaidForm has to be reset when UpdateOverlayContent
				Class("v-navigation-drawer--temporary").
				Attr("v-model", "locals.menu").
				Location(LocationLeft).
				Temporary(true).
				// Fixed(true).
				Attr(":height", `"100%"`),
			body,
		}
	}

	return b.scoped(
		VCard(
			header,
			VCardText(body),
			bottomActions,
		).
			Attr("style", "height:inherit"), // fix height on fullscreen
	)
}

func (b *ContentComponentBuilder) BuildPage(ctx *web.EventContext) (comp h.HTMLComponent) {
	topActions := append(b.TopLeftActions, b.TopRightActions...)

	if b.PrimaryAction != nil {
		topActions = append(topActions, b.PrimaryAction)
	}

	if len(b.Menu) > 0 {
		topActions = append(h.HTMLComponents{MenuBtn()}, topActions...)
	}
	ctx.WithContextValue(CtxActionsComponent, topActions)
	return b.scoped(
		VLayout(
			VMain(
				b.Tabs,
				b.Body,
			),
		))
}
