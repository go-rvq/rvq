// The option of an ECharts chart of vx-chart: what the server says (its
// kind, its labels, its series) drawn in the theme's colors. Pure: no
// ECharts here (its gradients as plain objects), so it is tested alone.

export type ChartKind = 'line' | 'bars' | 'stacked' | 'hbars' | 'donut'

export interface ChartSeries {
  name: string
  // a color of the theme ("primary", "success"…) or a CSS color
  color?: string
  values: number[]
}

export interface ChartProps {
  kind: ChartKind
  // the labels of the axis (some empty: not written) or of the parts
  labels: string[]
  series: ChartSeries[]
  // the full label of each point (the tooltip's title); none: its label
  hints?: string[]
  // the donut's parts' colors
  colors?: string[]
  // says a value: 1234.5 → "1.234,5" (default: the locale's)
  format?: (v: number) => string
}

// Theme is how the chart finds its colors: a color of the theme (or CSS)
// as CSS, and the colors of its texts and lines.
export interface Theme {
  color: (name: string) => string
  text: string
  muted: string
  grid: string
  surface: string
}

// palette is the order of the series' colors when they have none.
export const palette = ['primary', 'info', 'warning', 'success', 'secondary', 'error']

const esc = (s: string) =>
  String(s).replace(/[&<>"']/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' })[c]!)

export const defaultFormat = (v: number) =>
  new Intl.NumberFormat(undefined, { maximumFractionDigits: 2 }).format(v)

// compact is an axis' value, short: 1.2k, 3M
export const compact = (v: number) =>
  new Intl.NumberFormat(undefined, { notation: 'compact', maximumFractionDigits: 1 }).format(v)

// seriesColor is the color of the series i
export const seriesColor = (theme: Theme, s: ChartSeries, i: number) =>
  theme.color(s.color || palette[i % palette.length])

// title is the full label of the point i
export const pointTitle = (p: ChartProps, i: number) => (p.hints && p.hints[i]) || p.labels[i] || ''

// tooltipHTML is the tooltip of the point i: its full label, and every
// series with its color and its value — the zeros too.
export function tooltipHTML(p: ChartProps, theme: Theme, i: number): string {
  const f = p.format || defaultFormat
  const rows = p.series.map((s, k) => {
    const v = s.values[i] ?? 0
    return (
      `<div style="display:flex;align-items:center;gap:8px" data-chart-tip-row>` +
      `<span style="display:inline-block;width:10px;height:10px;border-radius:2px;background:${seriesColor(theme, s, k)}"></span>` +
      `<span style="flex:1">${esc(s.name)}</span>` +
      `<b style="margin-left:12px">${esc(f(v))}</b></div>`
    )
  })
  return `<div style="font-weight:600;margin-bottom:4px" data-chart-tip-title>${esc(pointTitle(p, i))}</div>` + rows.join('')
}

const tooltipBase = (theme: Theme) => ({
  backgroundColor: theme.surface,
  borderColor: theme.grid,
  textStyle: { color: theme.text, fontSize: 12 },
  extraCssText: 'box-shadow:0 2px 8px rgba(0,0,0,.15);border-radius:6px'
})

// zoomable: the kinds by step, zoomed by a range dragged
export const zoomable = (kind: ChartKind) => kind === 'line' || kind === 'bars' || kind === 'stacked'

// chartOption is the ECharts option of p in theme.
export function chartOption(p: ChartProps, theme: Theme): Record<string, any> {
  const f = p.format || defaultFormat
  const legend =
    p.series.length > 1
      ? { show: true, type: 'scroll', top: 0, left: 0, icon: 'roundRect', itemWidth: 10, itemHeight: 10, textStyle: { color: theme.muted, fontSize: 12 } }
      : { show: false }

  if (p.kind === 'donut') {
    const s = p.series[0] || { name: '', values: [] }
    const total = s.values.reduce((a, b) => a + b, 0)
    return {
      animation: true,
      tooltip: {
        ...tooltipBase(theme),
        trigger: 'item',
        formatter: (x: any) =>
          `<div style="font-weight:600">${esc(pointTitle(p, x.dataIndex))}</div>` +
          `${esc(f(x.value))} (${total ? Math.round((x.value / total) * 100) : 0}%)`
      },
      legend: { show: true, orient: 'vertical', right: 0, top: 'middle', icon: 'roundRect', itemWidth: 10, itemHeight: 10, textStyle: { color: theme.muted, fontSize: 12 } },
      series: [
        {
          type: 'pie',
          name: s.name,
          radius: ['55%', '80%'],
          center: ['32%', '50%'],
          avoidLabelOverlap: true,
          label: { show: false },
          itemStyle: { borderColor: theme.surface, borderWidth: 2 },
          data: s.values.map((v, i) => ({
            name: p.labels[i] || '',
            value: v,
            itemStyle: { color: theme.color((p.colors && p.colors[i]) || palette[i % palette.length]) }
          }))
        }
      ]
    }
  }

  // counts (every value whole): no fractions on the axis
  const whole = p.series.every((s) => s.values.every((v) => Number.isInteger(v)))
  const valueAxis = {
    type: 'value',
    splitNumber: 3,
    minInterval: whole ? 1 : undefined,
    axisLabel: { color: theme.muted, fontSize: 11, formatter: (v: number) => compact(v) },
    splitLine: { lineStyle: { color: theme.grid } }
  }

  if (p.kind === 'hbars') {
    const s = p.series[0] || { name: '', values: [] }
    return {
      animation: true,
      grid: { left: 8, right: 48, top: 4, bottom: 4, containLabel: true },
      tooltip: { ...tooltipBase(theme), trigger: 'axis', axisPointer: { type: 'shadow' }, formatter: (xs: any[]) => tooltipHTML(p, theme, xs[0]?.dataIndex ?? 0) },
      xAxis: { ...valueAxis, show: false },
      yAxis: {
        type: 'category',
        inverse: true,
        data: p.labels,
        axisTick: { show: false },
        axisLine: { show: false },
        axisLabel: { color: theme.text, fontSize: 12, width: 160, overflow: 'truncate' }
      },
      series: [
        {
          type: 'bar',
          name: s.name,
          data: s.values,
          barMaxWidth: 16,
          itemStyle: { color: seriesColor(theme, s, 0), borderRadius: 3 },
          label: { show: true, position: 'right', color: theme.muted, fontSize: 11, formatter: (x: any) => f(x.value) }
        }
      ]
    }
  }

  // by step: line, bars, stacked
  const many = (p.labels.length || 0) > 40
  return {
    animation: true,
    legend,
    // right: room for the last label of the axis, centered at the edge
    grid: { left: 8, right: 24, top: p.series.length > 1 ? 30 : 10, bottom: 4, containLabel: true },
    tooltip: {
      ...tooltipBase(theme),
      trigger: 'axis',
      axisPointer: { type: p.kind === 'line' ? 'line' : 'shadow', lineStyle: { color: theme.muted } },
      formatter: (xs: any[]) => tooltipHTML(p, theme, xs[0]?.dataIndex ?? 0)
    },
    xAxis: {
      type: 'category',
      data: p.labels,
      boundaryGap: p.kind !== 'line',
      axisTick: { show: false },
      axisLine: { lineStyle: { color: theme.grid } },
      axisLabel: { color: theme.muted, fontSize: 11, interval: 0, hideOverlap: true }
    },
    yAxis: valueAxis,
    // a range dragged zooms in it: the toolbox's selection — shown (hidden,
    // ECharts makes no selection), its icons empty, no titles
    toolbox: {
      show: true,
      showTitle: false,
      itemSize: 1,
      right: 0,
      top: 0,
      feature: { dataZoom: { yAxisIndex: 'none', icon: { zoom: 'path://M0,0', back: 'path://M0,0' } } }
    },
    dataZoom: [{ type: 'inside', zoomOnMouseWheel: false, moveOnMouseMove: false, moveOnMouseWheel: false }],
    series: p.series.map((s, k) => {
      const color = seriesColor(theme, s, k)
      if (p.kind === 'line') {
        return {
          type: 'line',
          name: s.name,
          data: s.values,
          showSymbol: !many,
          symbolSize: 4,
          lineStyle: { width: k === 0 ? 2 : 1.25, color },
          itemStyle: { color },
          // the first series filled, its color fading
          areaStyle:
            k === 0
              ? {
                  color: {
                    type: 'linear',
                    x: 0,
                    y: 0,
                    x2: 0,
                    y2: 1,
                    colorStops: [
                      { offset: 0, color: withAlpha(color, 0.25) },
                      { offset: 1, color: withAlpha(color, 0) }
                    ]
                  }
                }
              : undefined
        }
      }
      return {
        type: 'bar',
        name: s.name,
        data: s.values,
        stack: p.kind === 'stacked' ? 'total' : undefined,
        barMaxWidth: 18,
        itemStyle: { color, borderRadius: p.kind === 'stacked' ? 0 : [3, 3, 0, 0] }
      }
    })
  }
}

// withAlpha is color ("rgb(1, 2, 3)", "#abc", a name) with alpha a: an
// rgb() made rgba(), else the color as it is
export function withAlpha(color: string, a: number): string {
  const m = /^rgb\(\s*([^)]+)\)$/.exec(color.trim())
  if (m) return `rgba(${m[1]}, ${a})`
  return color
}

// themeOf reads the theme's colors from the element el is in (Vuetify's
// --v-theme-<name>: "25,118,210").
export function themeOf(el: Element): Theme {
  const cs = getComputedStyle(el)
  const v = (name: string) => cs.getPropertyValue('--v-theme-' + name).trim()
  const rgb = (name: string, a = 1) => {
    const x = v(name)
    if (!x) return ''
    return a === 1 ? `rgb(${x})` : `rgba(${x}, ${a})`
  }
  return {
    color: (name: string) => rgb(name) || name,
    text: rgb('on-surface', 0.87) || '#333',
    muted: rgb('on-surface', 0.6) || '#777',
    grid: rgb('on-surface', 0.1) || '#eee',
    surface: rgb('surface') || '#fff'
  }
}
