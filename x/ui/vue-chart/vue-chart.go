package vue_chart

import (
	rvqjs "github.com/go-rvq/rvq/js"
	"github.com/go-rvq/rvq/web"
)

func JSComponentsPack() web.ComponentsPack {
	v, err := rvqjs.VueChart.ReadFile("vue-chart/dist/vue-chart.umd.cjs")
	if err != nil {
		panic(err)
	}

	return web.ComponentsPack(v)
}
