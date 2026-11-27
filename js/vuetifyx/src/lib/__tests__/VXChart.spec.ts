import { beforeEach, describe, expect, it } from 'vitest'
import { flushPromises } from '@vue/test-utils'
import { nextTick } from 'vue'
import VXChart from '../VXChart'
import { chartOption, tooltipHTML, withAlpha, type ChartProps, type Theme } from '../chart-option'
import { mountTemplate } from '@/lib/__tests__/testutils'

const theme: Theme = {
  color: (n) => ({ primary: 'rgb(1, 2, 3)', info: 'rgb(4, 5, 6)', success: 'rgb(7, 8, 9)' })[n] || n,
  text: '#111',
  muted: '#777',
  grid: '#eee',
  surface: '#fff'
}

const two: ChartProps = {
  kind: 'line',
  labels: ['01/11', '', '03/11'],
  hints: ['01/11/2026', '02/11/2026', '03/11/2026'],
  series: [
    { name: 'Visits', values: [3, 0, 5] },
    { name: 'Errors', color: 'success', values: [1, 0, 0] }
  ]
}

describe('chartOption', () => {
  it('the tooltip: the full label, every series with its color and value, the zeros too', () => {
    const tip = tooltipHTML({ ...two, format: String }, theme, 1)
    expect(tip).toContain('02/11/2026')
    expect(tip).toContain('Visits')
    expect(tip).toContain('Errors')
    expect(tip.match(/data-chart-tip-row/g)).toHaveLength(2)
    expect(tip).toContain('<b style="margin-left:12px">0</b>')
    expect(tip).toContain('rgb(1, 2, 3)')
    expect(tip).toContain('rgb(7, 8, 9)')
    // a name is text, never HTML
    expect(tooltipHTML({ ...two, series: [{ name: '<i>x</i>', values: [1] }] }, theme, 0)).toContain('&lt;i&gt;x&lt;/i&gt;')
  })

  it('the line: the first series filled, a legend for two, a range zoomed', () => {
    const o = chartOption(two, theme)
    expect(o.series).toHaveLength(2)
    expect(o.series[0].type).toBe('line')
    expect(o.series[0].areaStyle.color.colorStops[0].color).toBe('rgba(1, 2, 3, 0.25)')
    expect(o.series[1].areaStyle).toBeUndefined()
    expect(o.series[0].lineStyle.width).toBeGreaterThan(o.series[1].lineStyle.width)
    expect(o.legend.show).toBe(true)
    expect(o.xAxis.data).toEqual(two.labels)
    expect(o.dataZoom[0].type).toBe('inside')
    // counts: no fractions; a value with them, the axis' own
    expect(o.yAxis.minInterval).toBe(1)
    expect(chartOption({ ...two, series: [{ name: 'x', values: [0.5, 1] }] }, theme).yAxis.minInterval).toBeUndefined()
    // shown: hidden, ECharts makes no selection by a drag
    expect(o.toolbox.show).toBe(true)
    expect(o.toolbox.feature.dataZoom.yAxisIndex).toBe('none')
    expect(o.tooltip.formatter([{ dataIndex: 2 }])).toContain('03/11/2026')
  })

  it('the bars, the stacked, the horizontal bars, the ring', () => {
    expect(chartOption({ ...two, kind: 'bars' }, theme).series[0]).toMatchObject({ type: 'bar', stack: undefined })
    expect(chartOption({ ...two, kind: 'stacked' }, theme).series[1]).toMatchObject({ type: 'bar', stack: 'total' })
    const hb = chartOption({ kind: 'hbars', labels: ['a', 'b'], series: [{ name: 'n', values: [2, 1] }] }, theme)
    expect(hb.yAxis).toMatchObject({ type: 'category', inverse: true, data: ['a', 'b'] })
    expect(hb.series[0].type).toBe('bar')
    const d = chartOption({ kind: 'donut', labels: ['a', 'b'], colors: ['success'], series: [{ name: 'n', values: [3, 1] }] }, theme)
    expect(d.series[0].type).toBe('pie')
    expect(d.series[0].data[0]).toMatchObject({ name: 'a', value: 3, itemStyle: { color: 'rgb(7, 8, 9)' } })
    expect(d.series[0].data[1].itemStyle.color).toBe('rgb(4, 5, 6)')
    expect(d.tooltip.formatter({ dataIndex: 0, value: 3 })).toContain('75%')
  })

  it('withAlpha', () => {
    expect(withAlpha('rgb(1, 2, 3)', 0.5)).toBe('rgba(1, 2, 3, 0.5)')
    expect(withAlpha('#abc', 0.5)).toBe('#abc')
  })
})

// a fake ECharts: what vx-chart asks of it, recorded
let charts: any[] = []
function fakeECharts() {
  ;(window as any).RvqECharts = {
    init: () => {
      const handlers: Record<string, Function> = {}
      const zr: Record<string, Function> = {}
      const c: any = {
        actions: [] as any[],
        option: null as any,
        setOption: (o: any) => (c.option = o),
        getOption: () => c.option,
        dispatchAction: (a: any) => {
          c.actions.push(a)
          if (a.type === 'dataZoom') {
            c.option.dataZoom[0] = { ...c.option.dataZoom[0], start: a.start, end: a.end }
            handlers.datazoom?.()
          }
        },
        on: (ev: string, f: Function) => (handlers[ev] = f),
        getZr: () => ({ on: (ev: string, f: Function) => (zr[ev] = f) }),
        containPixel: () => true,
        convertFromPixel: (_: any, [x]: number[]) => [x / 100, 0],
        resize: () => {},
        dispose: () => {},
        handlers,
        zr
      }
      charts.push(c)
      return c
    }
  }
}

describe('vx-chart', () => {
  beforeEach(() => {
    charts = []
    fakeECharts()
    // the script: already loaded (its tag resolved at once)
    const append = document.head.appendChild.bind(document.head)
    document.head.appendChild = ((el: any) => {
      if (el.tagName === 'SCRIPT') {
        setTimeout(() => el.onload?.(), 0)
        return el
      }
      return append(el)
    }) as any
  })

  it('draws, zooms by a range, resets; a point opened emits its drill', async () => {
    const w = mountTemplate(VXChart, {
      src: '/admin/!charts/echarts.js?v=1',
      ...two,
      drill: [{ from: '2026-11-01', to: '2026-11-01' }, null, { from: '2026-11-03', to: '2026-11-03' }]
    })
    await new Promise((r) => setTimeout(r, 5))
    await flushPromises()
    await nextTick()
    expect(charts).toHaveLength(1)
    const c = charts[0]
    expect(c.option.series).toHaveLength(2)
    // the drag selects: the toolbox's zoom taken
    expect(c.actions[0]).toMatchObject({ type: 'takeGlobalCursor', key: 'dataZoomSelect', dataZoomSelectActive: true })

    // zoomed: the ↺ shown; clicked, all again
    c.option.dataZoom[0] = { type: 'inside', start: 20, end: 60 }
    c.handlers.datazoom()
    await nextTick()
    const reset = w.find('[data-chart-reset]')
    expect(reset.exists()).toBe(true)
    await reset.trigger('click')
    expect(c.actions.at(-1)).toMatchObject({ type: 'dataZoom', start: 0, end: 100 })
    await nextTick()
    expect(w.find('[data-chart-reset]').exists()).toBe(false)

    // a click at the point 2 (x 200: index 2): its drill; at 1, none
    c.zr.click({ offsetX: 200, offsetY: 10 })
    c.zr.click({ offsetX: 100, offsetY: 10 })
    expect(w.emitted('point')).toEqual([[{ index: 2, from: '2026-11-03', to: '2026-11-03' }]])
  })
})
