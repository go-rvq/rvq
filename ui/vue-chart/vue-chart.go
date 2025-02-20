package vue_chart

import (
	_ "embed"

	"github.com/qor5/web/v3"
)

//go:embed vue-chart-js/dist/vue-chart.umd.cjs
var js []byte

func JSComponentsPack() web.ComponentsPack {
	return web.ComponentsPack(js)
}
