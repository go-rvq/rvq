package examples_vuetify

// @snippet_begin(VuetifyMenuSample)

import (
	. "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/docs/docsrc/examples"
	"github.com/go-rvq/rvq/admin/docs/docsrc/utils"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

type formData struct {
	EnableMessages bool
	EnableHints    bool
}

var globalFavored bool

const favoredIconPortalName = "favoredIcon"

func HelloVuetifyMenu(ctx *web.EventContext) (pr web.PageResponse, err error) {
	var fv formData
	err = ctx.UnmarshalForm(&fv)
	if err != nil {
		return
	}

	pr.Body = VContainer(
		utils.PrettyFormAsJSON(ctx),
		web.Scope(
			VMenu(
				web.Slot(
					VBtn("Menu as Popover").
						Attr("v-bind", "props").
						Color("indigo"),
				).Name("activator").Scope("{ props }"),

				VCard(
					VList(
						VListItem(
							Template(
								web.Portal(
									favoredIcon(),
								).Name(favoredIconPortalName),
							).Attr("v-slot:append", true),
						).Attr("prepend-avatar", "https://cdn.vuetifyjs.com/images/john.jpg").
							Attr("subtitle", "Founder of Vuetify").
							Title("John Leider"),
					),
					VDivider(),
					VList(
						VListItem(
							VSwitch().Attr("v-model", "form.EnableMessages").
								Attr("color", "purple").
								Attr("label", "Enable messages").
								Attr("hide-details", true),
						),
						VListItem(
							VSwitch().Attr("v-model", "form.EnableHints").
								Attr("color", "purple").
								Attr("label", "Enable hints").
								Attr("hide-details", true),
						),
					),

					VCardActions(
						VSpacer(),
						VBtn("Cancel").Variant("text").
							On("click", "locals.myMenuShow = false"),
						VBtn("Save").Color("primary").
							Variant("text").OnClick("submit"),
					),
				).MinWidth(300),
			).CloseOnContentClick(false).
				Location("end").
				Attr("v-model", "locals.myMenuShow"),
		).Locals().Form().LocalsInit("{ myMenuShow: false }").FormInit(JSONString(fv)),
	)

	return
}

func favoredIcon() HTMLComponent {
	color := ""
	if globalFavored {
		color = "text-red"
	}

	return VBtn("").Variant("text").Icon("mdi-heart").Class(color).OnClick("toggleFavored")
}

func toggleFavored(ctx *web.EventContext) (er web.EventResponse, err error) {
	globalFavored = !globalFavored
	er.updatePortals = append(er.updatePortals, &web.PortalUpdate{
		Name: favoredIconPortalName,
		Body: favoredIcon(),
	})
	return
}

func submit(ctx *web.EventContext) (er web.EventResponse, err error) {
	er.Reload = true
	er.RunScript = "locals.myMenuShow = false"
	return
}

var HelloVuetifyMenuPB = web.Page(HelloVuetifyMenu).
	EventFunc("submit", submit).
	EventFunc("toggleFavored", toggleFavored)

var HelloVuetifyMenuPath = examples.URLPathByFunc(HelloVuetifyMenu)

// @snippet_end
