package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-rvq/rvq/admin/example/admin"
	"github.com/go-rvq/rvq/x/osenv"
)

func main() {
	h := admin.Router(admin.ConnectDB())

	port := osenv.Get("RVQ_PORT", "The port to serve the admin on", "9000")

	fmt.Println("Served at http://localhost:" + port)

	mux := http.NewServeMux()
	mux.Handle("/",
		middleware.RequestID(
			middleware.Logger(
				middleware.Recoverer(h),
			),
		),
	)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		panic(err)
	}
}
