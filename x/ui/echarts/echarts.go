// Package echarts is the pack of rvq's ECharts (js/echarts) for a page:
// Apache ECharts with what vx-chart (vuetifyx) draws, window.RvqECharts.
// An admin serves it as a lazy extra asset — no page loads it in its head,
// vx-chart does, the first time it shows a chart:
//
//	p.LazyExtraAsset(echarts.Path, "text/javascript", echarts.JSComponentsPack())
//	vx.VXChart().Src(p.ExtraAssetRef(echarts.Path))
package echarts

import (
	"io/fs"

	rvqjs "github.com/go-rvq/rvq/js"
	"github.com/go-rvq/rvq/web"
)

// Path is its path among an admin's extra assets:
// /!assets/extra/echarts.js.
const Path = "/echarts.js"

// JSComponentsPack is its script (echarts.min.js).
func JSComponentsPack() web.ComponentsPack {
	FS, _ := fs.Sub(rvqjs.ECharts, "echarts/dist")
	return web.ComponentsPackFromFile(FS, "echarts.min.js")
}
