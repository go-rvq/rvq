package presets

import (
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	h "github.com/theplant/htmlgo"
)

type commonPageConfig struct {
	// TODO it should be create in defaultToPage
	main h.HTMLComponent

	tabPanels []TabComponentFunc
	sidePanel ObjectComponentFunc
}

// TODO set common component which in editingBuilder or DetailingBuilder
// TODO defaultToPage build a common page
func defaultToPage(config commonPageConfig, obj interface{}, ctx *web.EventContext) h.HTMLComponent {
	msgr := MustGetMessages(ctx.Context())

	var asideContent = config.main

	if len(config.tabPanels) != 0 {
		var tabs []h.HTMLComponent
		var contents []h.HTMLComponent
		for _, panelFunc := range config.tabPanels {
			tab, content := panelFunc(obj, ctx)
			if tab != nil {
				tabs = append(tabs, tab)
				contents = append(contents, content)
			}
		}
		if len(tabs) > 0 {
			asideContent = web.Scope(
				VTabs(
					VTab(h.Text(msgr.FormTitle)).Value("default"),
					h.Components(tabs...),
				).Class("v-tabs--fixed-tabs").Attr("v-model", "locals.tab"),

				VTabsWindow(
					VTabsWindowItem(
						asideContent,
					).Value("default"),
					h.Components(contents...),
				).Attr("v-model", "locals.tab"),
			).Slot("{ locals }").LocalsInit(`{tab: 'default'}`)
		}
	}

	if config.sidePanel != nil {
		sidePanel := config.sidePanel(obj, ctx)
		if sidePanel != nil {
			asideContent = VContainer(
				VRow(
					VCol(asideContent).Cols(8),
					VCol(sidePanel).Cols(4),
				),
			)
		}
	}
	return asideContent
}
