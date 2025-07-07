package examples_vuetify

// @snippet_begin(VuetifySnackBarsSample)
import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/docs/docsrc/examples"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

func VuetifySnackBars(ctx *web.EventContext) (pr web.PageResponse, err error) {
	pr.Body = VContainer(
		VBtn("Show Snack Bar").OnClick("showSnackBar"),
		web.Portal().Name("snackbar"),
		snackbar("bottom", "success"),
	)

	return
}

func showSnackBar(ctx *web.EventContext) (er web.EventResponse, err error) {
	er.updatePortals = append(er.updatePortals,
		&web.PortalUpdate{
			Name: "snackbar",
			Body: snackbar("top", "red"),
		},
	)

	return
}

func snackbar(pos string, color string) *web.ScopeBuilder {
	return web.Scope(
		VSnackbar().Location(pos).Timeout(-1).Color(color).
			Attr("v-model", "locals.show").
			Children(
				h.Text("Hello, I am a snackbar"),
				web.Slot(
					VBtn("").Variant("text").
						Attr("@click", "locals.show = false").
						Children(VIcon("mdi-close")),
				).Name("actions"),
			),
	).Slot("{ locals }").LocalsInit(`{ show: true }`)
}

var VuetifySnackBarsPB = web.Page(VuetifySnackBars).
	EventFunc("showSnackBar", showSnackBar)

var VuetifySnackBarsPath = examples.URLPathByFunc(VuetifySnackBars)

// @snippet_end
