// <vx-chart>: a chart (line, bars, stacked, horizontal bars, ring) drawn
// by Apache ECharts — a tooltip at each point with every series, a legend
// that turns a series off, a range dragged zooms in it (↺ back to all), a
// point clicked opened (point: its drill). ECharts (js/echarts) is not in
// this bundle: src is loaded the first time one is shown, once for the
// page (window.RvqECharts).
import { defineComponent, h, nextTick, onBeforeUnmount, onMounted, ref, resolveComponent, watch, type PropType } from 'vue'
import { useTheme } from 'vuetify'
import { chartOption, themeOf, zoomable, type ChartKind, type ChartSeries } from './chart-option'

declare const window: Window & { RvqECharts?: any }

const loading: Record<string, Promise<void>> = {}

// once is the script src loaded once for the page
const once = (src: string): Promise<void> =>
  (loading[src] ??= new Promise<void>((resolve, reject) => {
    const el = document.createElement('script')
    el.src = src
    el.onload = () => resolve()
    el.onerror = () => {
      delete loading[src]
      reject(new Error('the charts could not be loaded: ' + src))
    }
    document.head.appendChild(el)
  }))

export default defineComponent({
  name: 'VxChart',
  props: {
    // the charts' code (x/ui/echarts: …/!assets/extra/echarts.js?v=…)
    src: { type: String, required: true },
    kind: { type: String as PropType<ChartKind>, default: 'line' },
    labels: { type: Array as PropType<string[]>, default: () => [] },
    series: { type: Array as PropType<ChartSeries[]>, default: () => [] },
    // the full label of each point
    hints: { type: Array as PropType<string[]>, default: () => [] },
    // the donut's parts' colors
    colors: { type: Array as PropType<string[]>, default: () => [] },
    // what opens each point (emitted with point); null: not opened
    drill: { type: Array as PropType<any[]>, default: () => [] },
    // its height, in pixels (0: by its kind)
    height: { type: Number, default: 0 },
    // its title, for the readers of the screen
    title: { type: String, default: '' },
    // the reset's text
    resetLabel: { type: String, default: 'Reset zoom' }
  },
  emits: ['point'],
  setup(props, { emit }) {
    const box = ref<HTMLElement>()
    const error = ref('')
    const ready = ref(false)
    const zoomed = ref(false)
    const theme = useTheme()
    let chart: any = null
    let ro: ResizeObserver | null = null

    const heightOf = () => {
      if (props.height) return props.height
      if (props.kind === 'hbars') return Math.max(80, props.labels.length * 28 + 8)
      if (props.kind === 'donut') return 200
      return 220
    }

    const draw = () => {
      if (!chart || !box.value) return
      chart.setOption(
        chartOption(
          { kind: props.kind, labels: props.labels, series: props.series, hints: props.hints, colors: props.colors },
          themeOf(box.value)
        ),
        true
      )
      if (zoomable(props.kind)) {
        // a drag selects the range zoomed in
        chart.dispatchAction({ type: 'takeGlobalCursor', key: 'dataZoomSelect', dataZoomSelectActive: true })
      }
      zoomed.value = false
    }

    const reset = () => {
      chart?.dispatchAction({ type: 'dataZoom', start: 0, end: 100 })
      zoomed.value = false
    }

    // the point under (x, y) of the chart: its index (by step), or -1
    const indexAt = (x: number, y: number): number => {
      if (!chart || !chart.containPixel({ gridIndex: 0 }, [x, y])) return -1
      const v = chart.convertFromPixel({ gridIndex: 0 }, [x, y])
      const i = Array.isArray(v) ? Math.round(v[0]) : -1
      return i >= 0 && i < props.labels.length ? i : -1
    }

    onMounted(async () => {
      try {
        await once(props.src)
        if (!window.RvqECharts) throw new Error('the charts did not register: ' + props.src)
        ready.value = true
        await nextTick()
        if (!box.value) return
        chart = window.RvqECharts.init(box.value, null, { renderer: 'svg' })
        draw()
        chart.on('datazoom', () => {
          const o = chart.getOption()?.dataZoom?.[0]
          zoomed.value = !!o && (o.start > 0 || o.end < 100)
        })
        // a point clicked: opened, when it has where (drill)
        chart.getZr().on('click', (e: any) => {
          let i = -1
          if (zoomable(props.kind)) i = indexAt(e.offsetX, e.offsetY)
          else if (e.target && typeof e.target.dataIndex === 'number') i = e.target.dataIndex
          if (i >= 0 && props.drill[i]) emit('point', { index: i, ...props.drill[i] })
        })
        chart.on('click', (x: any) => {
          if (!zoomable(props.kind) && props.drill[x.dataIndex]) emit('point', { index: x.dataIndex, ...props.drill[x.dataIndex] })
        })
        ro = new ResizeObserver(() => chart?.resize())
        ro.observe(box.value)
      } catch (e) {
        error.value = String(e instanceof Error ? e.message : e)
      }
    })

    onBeforeUnmount(() => {
      ro?.disconnect()
      chart?.dispose()
      chart = null
    })

    watch(() => [props.kind, props.labels, props.series, props.hints, props.colors], draw, { deep: true })
    // the theme changed (light, dark): its colors read again
    watch(() => theme.global.name.value, () => nextTick(draw))

    return () => {
      if (error.value) return h('div', { class: 'text-error text-caption', 'data-chart-error': '' }, error.value)
      const children: any[] = []
      if (!ready.value) children.push(h(resolveComponent('v-progress-linear'), { indeterminate: true, color: 'primary', height: 2 }))
      children.push(
        h('div', {
          ref: box,
          role: 'img',
          'aria-label': props.title,
          style: { width: '100%', height: heightOf() + 'px', cursor: props.drill.some(Boolean) ? 'pointer' : undefined }
        })
      )
      if (zoomed.value) {
        children.push(
          h(resolveComponent('v-btn'), {
            icon: 'mdi-restore',
            size: 'x-small',
            variant: 'tonal',
            title: props.resetLabel,
            'aria-label': props.resetLabel,
            'data-chart-reset': '',
            style: { position: 'absolute', top: '2px', right: '2px', zIndex: 2 },
            onClick: reset
          })
        )
      }
      return h('div', { class: 'vx-chart', style: { position: 'relative' }, 'data-chart': props.kind }, children)
    }
  }
})
