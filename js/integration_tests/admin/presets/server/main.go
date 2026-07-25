// Command server boots a preset test app (in-memory SQLite, seeded) on an HTTP
// port, for the bun/vitest integration tests to drive. The APP env var selects
// the app: "" / "listeditor" → the one-level list editor; "nested" → the
// four-level nested list editor; "nestedmodels" → two levels of nested MODELS
// (Survey → Places → Products), each with its own listing/detailing.
//
// It prints "LISTENING http://<addr>" once ready (via integration.ServeEnvMain);
// PORT=0 → ephemeral port.
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-rvq/rvq/admin/presets/integration"
)

func main() {
	var (
		handler http.Handler
		err     error
	)
	switch os.Getenv("APP") {
	case "nested":
		handler, err = integration.NewNestedSeededHandler()
	case "composite":
		handler, err = integration.NewCompositeSeededHandler()
	case "nestedmodels":
		handler, err = integration.NewNestedModelsSeededHandler()
	default:
		handler, err = integration.NewSeededHandler()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "build app:", err)
		os.Exit(1)
	}
	integration.ServeEnvMain(handler)
}
