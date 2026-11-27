package js

import "embed"

// ECharts is the Apache ECharts of the admin's charts (js/echarts:
// dist/echarts.min.js, window.RvqECharts), with only what vx-chart draws:
// served by admin/charts, loaded by a page the first time it shows a chart.
//
//go:embed echarts/dist
var ECharts embed.FS
