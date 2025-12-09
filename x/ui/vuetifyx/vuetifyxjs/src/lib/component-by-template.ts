import * as Vue from 'vue'
import { ComponentInternalInstance, DefineComponent, getCurrentInstance } from 'vue'
const { defineComponent } = Vue

declare let window: any

export interface ComponentDefinition {
  template: string,
  options?: RenderOptions
}

export interface Data {
  window: any
  Vue: typeof Vue

  [key: string]: any
}

export interface Locals {
  // the component definition
  componentDefinition: ComponentDefinition
  // the component element
  element: Node | null
  // the component instance
  componentInstance:ComponentInternalInstance | null
  // list of component nodes
  nodes: Node[]
  // skip calls unmounted handlers
  skipUnmountHandlers: boolean
  // other data set by any callbacks
  [key: string]: any
}

export interface CallbackArgs {
  data: Data
  locals: Locals
}

export type Callbacks = ((args:CallbackArgs) => boolean | void)[]

export interface RenderOptions {
  data?: Record<string, any>
  setup?: Callbacks
  mounted?: Callbacks
  unmounted?: Callbacks
  resetScroll?: boolean
}

const callCallbacks = (callbacks: Callbacks, args:CallbackArgs) => {
  if (!callbacks) return

  for (let i = 0; i < callbacks.length; i++) {
    let fn = callbacks[i]
    if (!fn) continue
    if (typeof fn === 'string') {
      fn = new Function('args', fn) as ((args:CallbackArgs) => void | boolean)
    }
    if (fn(args) === false) {
      break
    }
  }
}

export function buildComponentData(cd: ComponentDefinition) {
  const options = cd.options ?? {}
  const data: Data = {
    Vue: Vue,
    window
  }

  if (options.data) {
    const od = options.data as Record<string, any>
    Object.keys(options.data).forEach(name => {
      data[name] = Reflect.get(od, name)
    })
  }
  data.vx$scope = data
  return data
}

export function componentByTemplate(componentDefinition: ComponentDefinition): DefineComponent {
  const cd: ComponentDefinition = { ...componentDefinition }
  cd.options ??= {}

  const options = cd.options

  let { mounted, unmounted, setup } = options
  mounted ||= []
  unmounted ||= []
  setup ||= []

  const data = buildComponentData(cd)

  const resetScroll = cd.options?.resetScroll || false

  if (resetScroll) {
    mounted.push(({ locals }) => {
      (locals.element as HTMLElement)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
    })
  }

  const locals = {
    componentDefinition: cd,
    element: null,
    skipUnmountHandlers: false,
    componentInstance: null,
    nodes: []
  } as Locals

  return defineComponent({
    setup() {
      locals.componentInstance = getCurrentInstance()
      callCallbacks(setup, {data, locals})
      return data
    },
    mounted() {
      this.$nextTick(() => /**/ {
        const el = this.$el as Node
        if (el) {
          locals.element = el

          const nodes:Node[] = []
          let curEl:Node | null = el

          if (curEl.nodeType == Node.TEXT_NODE) {
            let nt = curEl as Text

            while (nt && nt.data == '') {
              curEl = nt.nextSibling
              if (curEl?.nodeType == Node.TEXT_NODE) {
                nt = curEl as Text
              } else {
                break
              }
            }

            while (curEl) {
              if (curEl.nodeType == Node.TEXT_NODE) {
                nt = curEl as Text
                if (nt.data == '') {
                  break
                }
              }
              nodes.push(curEl)
              curEl = curEl.nextSibling as HTMLElement
            }
            locals.nodes = nodes
          } else {
            const hel = el as HTMLElement
            if (hel.style && hel.style.height) {
              hel.style.height = this.$el.style.height
            }
            locals.nodes = [el]
          }
        }
        callCallbacks(mounted, { data, locals })
      })
    },
    unmounted() {
      if (!locals.skipUnmountHandlers) {
        callCallbacks(unmounted, { data, locals })
      }
    },
    template: cd.template
  })
}
