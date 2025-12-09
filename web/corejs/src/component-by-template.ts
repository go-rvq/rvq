import type { DefineComponent, Ref } from 'vue'
import * as Vue from 'vue'
import { defineComponent, inject, reactive, ref } from 'vue'

declare let window: any

interface ComponentStepsAPI {
  Config: string
  PreSetup: string
  PostSetup: string
  Data: string
  Methods: string
}

const ComponentSteps: ComponentStepsAPI = {
  Config: 'config',
  PreSetup: 'pre-setup',
  PostSetup: 'post-setup',
  Data: 'data',
  Methods: 'methods'
}

interface Context {
  form: any
  locals: any
  plaid: any
  vars: any
  closer: any
  computed: any
  updateRootTemplate: any
}

interface Config {
  name: string
  form?: any
  locals?: any
  portal: Ref
  methods?: (ctx: Context) => object
  data?: (ctx: Context) => object
}

export interface RenderSelf {
  element: HTMLElement,
  comp: any,
  window: any,
  Vue:any
}

export interface RenderOptions {
  mounted?: ((self: RenderSelf) => void)[]
  unmounted?: ((self: RenderSelf) => void)[]
  resetScroll?:boolean
}

const callCallbacks = (funcs: any[], self: RenderSelf) => {
  funcs?.forEach((fn: any) => {
    if (!fn) return
    if (typeof fn === 'string') {
      fn = new Function('self', fn)
    }
    fn(self)
  })
}

export function componentByTemplate(
  template: string,
  scope: any = {},
  portal: Ref = ref(),
  options: RenderOptions = {}
): DefineComponent {
  options = {...(options || {})}
  if (options.resetScroll) {
    options.mounted?.push(self => {
      self.element.scrollIntoView({behavior: 'smooth', block: 'start'})
    })
  }

  return defineComponent({
    setup() {
      const plaid = inject('plaid'),
        vars = inject('vars'),
        closer = inject('closer'),
        fullscreen = inject('fullscreen'),
        isFetching = inject('isFetching'),
        updateRootTemplate = inject('updateRootTemplate')

      const localScope = {
        plaid,
        vars,
        closer,
        fullscreen,
        isFetching,
        updateRootTemplate,
        options: options,
        ...scope
      }

      localScope.Vue = Vue
      localScope.window = window
      localScope.SCOPE = localScope

      return reactive(localScope)
    },
    mounted() {
      this.$nextTick(() => /**/ {
        if (this.$el) {
          if (this.$el.style && this.$el.style.height) {
            portal.value.style.height = this.$el.style.height
          }
        }

        callCallbacks(options.mounted as any[], {
          element: portal.value,
          comp: this,
          window,
          Vue
        })
      })
    },
    unmounted() {
      callCallbacks(options.unmounted as any[], {
        element: portal.value,
        comp: this,
        window,
        Vue
      })
    },
    template
  })
}
