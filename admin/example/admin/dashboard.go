package admin

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/x/ui/vuetify"
)

func Dashboard() h.HTMLComponent {
	return vuetify.VContainer(
		h.H1("Welcome to the RVQ demo site").Class("mt-8"),

		h.A().Text("RVQ Website").Href("https://rvq.com").Target("_blank"),
		h.A().Text("RVQ Documentation").Href("https://docs.rvq.com").Target("_blank").Class("ml-4"),
		h.A().Text("Source Code").Href("https://github.com/qor5/admin/tree/main/example").Target("_blank").Class("ml-4"),
	)
}
