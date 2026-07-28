package admin

import (
	"net/http"

	"github.com/go-rvq/rvq/admin/presets"
)

func setupRouter(b *presets.Builder) (mux *http.ServeMux) {
	mux = http.NewServeMux()
	mux.Handle("/", b)
	return
}
