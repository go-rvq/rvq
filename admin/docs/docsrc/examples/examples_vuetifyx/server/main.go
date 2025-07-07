package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-rvq/rvq/admin/docs/docsrc/examples/examples_vuetify"
	"github.com/go-rvq/rvq/admin/docs/docsrc/examples/examples_vuetifyx"
	"github.com/go-rvq/rvq/web"
	"github.com/theplant/osenv"
)

var port = osenv.Get("PORT", "The port to serve on", "7800")

func main() {
	fmt.Println("Starting docs at :" + port)
	mux := http.NewServeMux()
	examples_vuetifyx.Mux(mux, "")
	im := &examples_vuetify.IndexMux{Mux: http.NewServeMux()}
	examples_vuetifyx.SamplesHandler(im, "/samples")
	mux.Handle("/samples/",
		middleware.Logger(
			middleware.RequestID(
				im.Mux,
			),
		),
	)
	mux.Handle("/", web.New().Page(im.Page))

	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		panic(err)
	}
}
