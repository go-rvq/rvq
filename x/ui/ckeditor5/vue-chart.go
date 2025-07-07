package ckeditor5

import (
	_ "embed"

	"github.com/qor5/web/v3"
)

//go:embed ckeditor5-js/dist/vue-chart.umd.cjs
var js []byte

func JSComponentsPack() web.ComponentsPack {
	return web.ComponentsPack(js)
}
