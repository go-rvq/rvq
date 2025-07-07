package examples_vuetifyx

import (
	"net/http"

	"github.com/go-rvq/rvq/admin/docs/docsrc/assets"
	"github.com/go-rvq/rvq/admin/docs/docsrc/examples"
	"github.com/go-rvq/rvq/admin/docs/docsrc/examples/examples_vuetify"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/go-rvq/rvq/x/ui/vuetifyx"
)

func Mux(mux *http.ServeMux, prefix string) http.Handler {
	mux.Handle("/assets/main.js",
		web.PacksHandler("text/javascript",
			JSComponentsPack(),
			vuetifyx.JSComponentsPack(),
			Vuetify(),
			web.JSComponentsPack(),
		),
	)

	mux.Handle("/assets/vue.js",
		web.PacksHandler("text/javascript",
			web.JSVueComponentsPack(),
		),
	)

	HandleMaterialDesignIcons(prefix, mux)

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Write(assets.Favicon)
		return
	})

	return mux
}

func SamplesHandler(mux examples.Muxer, prefix string) {
	mux.Handle(
		VuetifyComponentsLinkageSelectPath,
		VuetifyComponentsLinkageSelectPB.Wrap(examples_vuetify.DemoVuetifyLayout),
	)
	mux.Handle(
		ExpansionPanelDemoPath,
		ExpansionPanelDemoPB.Wrap(examples_vuetify.DemoVuetifyLayout),
	)
	mux.Handle(
		KeyInfoDemoPath,
		KeyInfoDemoPB.Wrap(examples_vuetify.DemoVuetifyLayout),
	)
	mux.Handle(
		FilterDemoPath,
		FilterDemoPB.Wrap(examples_vuetify.DemoVuetifyLayout),
	)
	mux.Handle(
		DatePickersPath,
		DatePickersPB.Wrap(examples_vuetify.DemoVuetifyLayout),
	)
	return
}
