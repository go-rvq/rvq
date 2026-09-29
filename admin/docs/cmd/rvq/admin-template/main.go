package main

import (
	"fmt"
	"net/http"

	"github.com/go-rvq/rvq/admin/docs/cmd/rvq/admin-template/admin"
	"github.com/go-rvq/rvq/x/osenv"
)

func main() {
	// Setup project
	mux := admin.Initialize()

	port := osenv.Get("RVQ_PORT", "The port to serve the admin on", "9000")

	fmt.Println("Served at http://localhost:" + port + "/admin")

	http.Handle("/", mux)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
