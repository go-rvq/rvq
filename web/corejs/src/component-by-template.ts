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

export function componentByTemplate(
  template: string,
  scope: any = {},
  portal: Ref = ref()
): DefineComponent {
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
        ...scope
      }

      localScope.Vue = Vue
      localScope.window = window
      localScope.SCOPE = localScope

      return reactive(localScope)
    },
    mounted() {
      this.$nextTick(() => /**/ {
        if (this.$el && this.$el.style && this.$el.style.height) {
          portal.value.style.height = this.$el.style.height
        }
      })
    },
    template
  })
}
