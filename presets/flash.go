package presets

import (
	"fmt"

	"github.com/qor5/web/v3"
	. "github.com/qor5/x/v3/ui/vuetify"
	h "github.com/theplant/htmlgo"
)

func RenderFlash(v any, color string) h.HTMLComponent {
	var text string
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
	default:
		text = fmt.Sprintf("%v", t)
	}

	return web.Scope(
		VSnackbar(
			h.Text(text),
			web.Slot(
				VBtn("").Variant("text").
					Attr("@click", "locals.show = false").
					Children(VIcon("mdi-close")),
			).Name("actions"),
		).Location("top").
			Timeout(-1).
			Color(color).
			Attr("v-model", "locals.show"),
	).Slot("{ locals }").LocalsInit(`{ show: true }`)
}
