package js

import "embed"

// ECharts is the Apache ECharts of the admin's charts (js/echarts:
// dist/echarts.min.js, window.RvqECharts), with only what vx-chart draws:
// x/ui/echarts' pack, a lazy extra asset of the admin (/!assets/extra/
// echarts.js), loaded by a page the first time it shows a chart.
//
//go:embed echarts/dist
var ECharts embed.FS
