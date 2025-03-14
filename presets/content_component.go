package presets

import (
	"fmt"

	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	"github.com/qor5/web/v3/vue"
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
		).
		Attr("@click.menu", `() => {
console.log($contentComponent);
$contentComponent.menu = !$contentComponent.menu
}`)
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
	Obj             any
	Context         *web.EventContext
	Title           string
	TopLeftActions  h.HTMLComponents
	TopRightActions h.HTMLComponents
	Tabs            []TabComponentFunc
	PreBody         h.HTMLComponents
	PostBody        h.HTMLComponents
	PrimaryAction   h.HTMLComponent
	BottomActions   h.HTMLComponents
	TopBar          h.HTMLComponent
	BottomBar       h.HTMLComponent
	Overlay         *ContentComponentBuilderOverlay
	Menu            h.HTMLComponents
	Body            h.HTMLComponent
	Scope           *web.ScopeBuilder
	MainPortals     h.HTMLComponents
	Notices         h.HTMLComponents
}

func (b *ContentComponentBuilder) scoped(comp h.HTMLComponent) h.HTMLComponent {
	if b.Scope == nil {
		return comp
	}
	b.Scope.Children(comp)
	return b.Scope
}

func (b *ContentComponentBuilder) build(comp h.HTMLComponent) h.HTMLComponent {
	comp = b.scoped(comp)
	if len(b.MainPortals) > 0 {
		comp = append(b.MainPortals, comp)
	}
	return vue.UserComponent(comp).Scope("$contentComponent", vue.Var(`{menu:false}`))
}

func (b *ContentComponentBuilder) JoinedBody() h.HTMLComponent {
	comps := append(b.Notices, b.PreBody...)
	if b.Body != nil {
		comps = append(comps, b.Body)
	}

	comps = append(comps, b.PostBody...)
	if len(comps) == 0 {
		return nil
	}
	return comps
}

func (b *ContentComponentBuilder) BuildOverlay() h.HTMLComponent {
	var (
		header           h.HTMLComponent
		headerComponents = b.TopLeftActions
		headerRight      = b.TopRightActions
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

	body := b.JoinedBody()

	if body != nil && !b.Overlay.FullscreenDisabled {
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
		headerComponents = h.HTMLComponents{VToolbar(headerComponents...).Color("white").Elevation(0).Density(DensityCompact)}
	}

	if len(headerComponents) > 0 || b.TopBar != nil {
		header = h.HTMLComponents{
			headerComponents,
			b.TopBar,
			VDivider(),
		}
	}

	var bottom h.HTMLComponents

	if b.BottomBar != nil {
		bottom = append(bottom, b.BottomBar)
	}

	if len(b.BottomActions) > 0 {
		bottom = append(bottom, VCardActions(b.BottomActions...))
	}

	if len(bottom) > 0 {
		bottom = append(h.HTMLComponents{VDivider()}, bottom...)
	}

	if len(b.Tabs) > 0 {
		var tabs h.HTMLComponents
		var contents []h.HTMLComponent
		for _, panelFunc := range b.Tabs {
			tab, content := panelFunc(b.Obj, b.Context)
			if tab != nil {
				tabs = append(tabs, tab)
				contents = append(contents, content)
			}
		}
		t := h.HTMLComponents{
			VTabs(
				VTab(h.Text(MustGetMessages(b.Context.Context()).FormTitle)).Value("default"),
				h.Components(tabs...),
			).Class("v-tabs--fixed-tabs").Attr("v-model", "locals.tab"),

			VTabsWindow(
				VTabsWindowItem(
					body,
				).Value("default"),
				h.Components(contents...),
			).Attr("v-model", "locals.tab"),
		}

		if b.Scope != nil {
			b.Scope.AppendInit("{tab: 'default'}")
		} else {
			body = web.Scope(t).Slot("{ locals }").LocalsInit(`{tab: 'default'}`)
		}
	}

	if body != nil {
		body = VContainer(body).Fluid(true)
	}

	roots := h.HTMLComponents{header}

	if len(b.Menu) > 0 {
		layout := VLayout(
			VNavigationDrawer(b.Menu).
				// Attr("@input", "plaidForm.dirty && vars.presetsRightDrawer == false && !confirm('You have unsaved changes on this form. If you close it, you will lose all unsaved changes. Are you sure you want to close it?') ? vars.presetsRightDrawer = true: vars.presetsRightDrawer = $event"). // remove because drawer plaidForm has to be reset when UpdateOverlayContent
				Class("v-navigation-drawer--temporary").
				Attr("v-model", "$contentComponent.menu").
				Location(LocationLeft).
				Floating(true).
				Temporary(true),
			h.Div(body).Attr("style", `width:inherit;height:inherit;overflow:auto`),
		).
			Width("100%").
			Height("100%").
			Class("v-layout__content-component").
			SetAttr("style", "display: block;")
		body = layout
	}

	if body != nil {
		roots = append(roots, VCardText(body).Class("v-card-text__content-component").Style("padding: 0"))
	}

	roots = append(roots, bottom...)

	card := VCard(roots...).Height("inherit").Width("inherit").Class("v-card__content-component")

	return b.build(
		card, // .			Attr("style", "height:inherit,max-height:inherit"), // fix height on fullscreen
	)
}

func (b *ContentComponentBuilder) BuildPage() (comp h.HTMLComponent) {
	topActions := append(b.TopLeftActions, b.TopRightActions...)

	if b.PrimaryAction != nil {
		topActions = append(topActions, b.PrimaryAction)
	}

	if len(b.Menu) > 0 {
		b.Context.WithContextValue(CtxMenuComponent, b.Menu)
	}

	if len(topActions) > 0 {
		b.Context.WithContextValue(CtxActionsComponent, topActions)
	}

	body := b.JoinedBody()

	if len(b.Tabs) > 0 {
		var tabs h.HTMLComponents
		var contents []h.HTMLComponent
		for _, panelFunc := range b.Tabs {
			tab, content := panelFunc(b.Obj, b.Context)
			if tab != nil {
				tabs = append(tabs, tab)
				contents = append(contents, content)
			}
		}
		body = web.Scope(
			VTabs(
				VTab(h.Text(MustGetMessages(b.Context.Context()).FormTitle)).Value("default"),
				h.Components(tabs...),
			).Class("v-tabs--fixed-tabs").Attr("v-model", "locals.tab"),

			VTabsWindow(
				VTabsWindowItem(
					body,
				).Value("default"),
				h.Components(contents...),
			).Attr("v-model", "locals.tab"),
		).Slot("{ locals }").LocalsInit(`{tab: 'default'}`)
	}

	return VContainer(b.build(body)).Fluid(true)
}

func (b *ContentComponentBuilder) Build(p *Builder, r *web.PageResponse) (err error) {
	if b.Overlay.Mode.Overlayed() {
		er := &web.EventResponse{}
		comp := b.BuildOverlay()
		if b.Overlay.Mode.IsDialog() {
			p := p.Dialog().SetValidWidth(b.Overlay.Width)
			p.Respond(b.Context, er, comp)
		} else {
		}
	} else {
		r.PageTitle = b.Title
		r.Body = b.BuildPage()
	}
	return nil
}

func (b *ContentComponentBuilder) Notice(v any) *ContentComponentBuilder {
	var (
		color, text string
	)
	switch t := v.(type) {
	case *web.ValidationErrors:
		gErr := t.GetGlobalError()
		if len(gErr) > 0 {
			text = gErr
			color = "error"
		}
	case error:
		color = "error"
		text = t.Error()
	case string:
		text = t
	case h.HTMLComponent:
		b.Notices = append(b.Notices, t)
		return b
	default:
		text = fmt.Sprintf("%v", t)
	}

	if text != "" {
		b.Notices = append(b.Notices, RenderFlash(text, color))
	}
	return b
}
