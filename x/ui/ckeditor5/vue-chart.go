package ckeditor5

import (
	_ "embed"

	"github.com/go-rvq/rvq/web"
)

//go:embed ckeditor5-js/dist/vue-chart.umd.cjs
var js []byte

func JSComponentsPack() web.ComponentsPack {
	return web.ComponentsPack(js)
}
