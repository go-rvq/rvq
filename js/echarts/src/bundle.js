// Apache ECharts with what vx-chart (vuetifyx) draws, and no more: the
// line, the bars, the ring (pie); the tooltip, the legend, the zoom by a
// range (dataZoom, the toolbox's selection), the grid; the SVG renderer
// (crisp texts at any size). The page's charts find it at
// window.RvqECharts once it is loaded.
import * as echarts from 'echarts/core'
import { BarChart, LineChart, PieChart } from 'echarts/charts'
import {
  DataZoomComponent,
  GridComponent,
  LegendComponent,
  ToolboxComponent,
  TooltipComponent
} from 'echarts/components'
import { SVGRenderer } from 'echarts/renderers'

echarts.use([
  LineChart,
  BarChart,
  PieChart,
  GridComponent,
  TooltipComponent,
  LegendComponent,
  DataZoomComponent,
  ToolboxComponent,
  SVGRenderer
])

window.RvqECharts = echarts
