import * as Vue from 'vue'
import {ComponentInternalInstance, DefineComponent, getCurrentInstance} from 'vue'

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
  // updated: the root nodes before (locals.nodes, the ones now)
  previous?: Node[]
}

export type Callbacks = ((args:CallbackArgs) => boolean | void)[]

export interface RenderOptions {
  data?: Record<string, any>
  setup?: Callbacks
  mounted?: Callbacks
  // updated: its root nodes changed (an element in or out: a v-if, a
  // v-for of the root) — locals.nodes the new ones, previous the old
  updated?: Callbacks
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

// isContentNode says a root node is of the content: an element, or a text
// with something written — not a comment (the place of a v-if) nor a
// fragment's anchor (an empty text).
const isContentNode = (n: Node) =>
  n.nodeType === Node.ELEMENT_NODE || (n.nodeType === Node.TEXT_NODE && (n as Text).data.trim() !== '')

// rootNodes are the root nodes of the content the component instance
// rendered: its element; of a template of several roots (a fragment), the
// nodes between the fragment's anchors — nested fragments (v-for) too —; of
// a component at the root, its own.
export function rootNodes(inst: ComponentInternalInstance | null): Node[] {
  const tree = inst?.subTree
  if (!tree || !tree.el) {
    return []
  }
  if (tree.component) {
    return rootNodes(tree.component)
  }
  const nodes: Node[] = []
  if (tree.type === Vue.Fragment) {
    for (let n = (tree.el as Node).nextSibling; n && n !== tree.anchor; n = n.nextSibling) {
      nodes.push(n)
    }
  } else {
    nodes.push(tree.el as Node)
  }
  return nodes.filter(isContentNode)
}

const sameNodes = (a: Node[], b: Node[]) => a.length === b.length && a.every((n, i) => n === b[i])

export function componentByTemplate(componentDefinition: ComponentDefinition): DefineComponent {
  const cd: ComponentDefinition = { ...componentDefinition }
  cd.options ??= {}

  const options = cd.options

  let { mounted, updated, unmounted, setup } = options
  mounted ||= []
  updated ||= []
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
          locals.nodes = rootNodes(locals.componentInstance)
        }
        callCallbacks(mounted, { data, locals })
      })
    },
    // its root nodes after a render of its own (a v-if, a v-for at the root):
    // the updated callbacks when they changed
    updated() {
      const previous = locals.nodes
      const nodes = rootNodes(locals.componentInstance)
      if (sameNodes(previous, nodes)) {
        return
      }
      locals.nodes = nodes
      callCallbacks(updated, { data, locals, previous })
    },
    unmounted() {
      if (!locals.skipUnmountHandlers) {
        callCallbacks(unmounted, { data, locals })
      }
    },
    template: cd.template
  })
}
