package presets

import (
	"fmt"

	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/vue"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
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
	FullscreenDisabled bool
	CloseDisabled      bool
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

func (b *ContentComponentBuilder) AddMenu(com ...h.HTMLComponent) *ContentComponentBuilder {
	b.Menu = append(b.Menu, com...)
	return b
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
		headerLeft  = b.TopLeftActions
		headerRight = b.TopRightActions
		dialog      bool
		tag         = vx.VXAdvancedExpandCloseCard("")
		width       = "500"
	)

	if b.Overlay != nil {
		dialog = b.Overlay.Mode.IsDialog()
		if b.Overlay.Width != "" {
			width = b.Overlay.Width
		}
	}

	tag.SetDensity("compact")
	tag.Width(width)
	tag.DialogOrDrawer(dialog)

	if !dialog {
		tag.Attr("location", "right")
		tag.Attr("temporary", true)
	}

	if b.Title != "" {
		tag.SetTitle(b.Title)
	}

	if b.PrimaryAction != nil {
		headerRight = append(headerRight, b.PrimaryAction)
	}

	if len(headerLeft) > 0 {
		tag.SetSlotPrependToolbar(headerLeft...)
	}

	if len(headerRight) > 0 {
		tag.SetSlotAppendToolbar(headerRight...)
	}

	body := b.JoinedBody()

	if body != nil && !b.Overlay.FullscreenDisabled {
		tag.Expandable(true)
	}

	if !b.Overlay.CloseDisabled {
		tag.Closable(true)
	}

	if b.TopBar != nil {
		tag.SetSlotTop(b.TopBar)
	}

	var bottom h.HTMLComponents

	if b.BottomBar != nil {
		bottom = append(bottom, b.BottomBar)
	}

	if len(b.BottomActions) > 0 {
		bottom = append(bottom, b.BottomActions...)
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

	if len(b.Menu) > 0 {
		tag.SlotMainMenu(b.Menu)
	}

	if body != nil {
		tag.SlotBody(body)
	}

	if len(b.MainPortals) > 0 {
		tag.SlotPortals(b.MainPortals...)
	}

	if len(bottom) > 0 {
		tag.SlotBottom(bottom)
	}

	return tag
}

func (b *ContentComponentBuilder) BuildPage() (comp h.HTMLComponent) {
	topActions := append(b.TopLeftActions, b.TopRightActions...)

	if b.PrimaryAction != nil {
		topActions = append(topActions, b.PrimaryAction)
	}

	if len(b.Menu) > 0 {
		WithMenuComponent(b.Context, b.Menu)
	}

	if len(topActions) > 0 {
		WithActionsComponent(b.Context, topActions...)
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
		b.Notices = append(b.Notices, VAlert(h.Text(text)).Color(color))
	}
	return b
}
