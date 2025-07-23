import {
  type App,
  createApp,
  type DefineComponent,
  defineComponent,
  onMounted,
  provide,
  reactive,
  ref,
  shallowRef
} from 'vue'
import { GlobalEvents } from 'vue-global-events'
import GoPlaidScope from '@/go-plaid-scope.vue'
import GoPlaidPortal from '@/go-plaid-portal.vue'
import GoPlaidRunScript from '@/go-plaid-run-script.vue'
import UserComponent from '@/user-component.vue'
import { componentByTemplate } from '@/component-by-template'
import { Builder, plaid } from '@/builder'
import { keepScroll } from '@/keepScroll'
import { assignOnMounted } from '@/assign'

export const Root = defineComponent({
  props: {
    initialTemplate: {
      type: String,
      required: true
    }
  },

  setup(props, { emit, attrs, expose }) {
    const current = shallowRef<DefineComponent | null>(null)
    const form = reactive({})
    provide('form', form)

    const locals = reactive({})
    provide('locals', locals)

    const closer = reactive({})
    provide('closer', closer)

    const fullscreen = reactive({})
    provide('fullscreen', locals)

    const vars = reactive({
      __notification: {}
    })
    provide('vars', vars)

    const _plaid = (): Builder => {
      return plaid().updateRootTemplate(updateRootTemplate).vars(vars)
    }
    provide('plaid', _plaid)

    const isFetching = ref(false)
    provide('isFetching', isFetching)

    const updateRootTemplate = (template: string) => {
      current.value = componentByTemplate(template, {
        form,
        locals
      })
    }
    provide('updateRootTemplate', updateRootTemplate)

    onMounted(() => {
      updateRootTemplate(props.initialTemplate)

      window.addEventListener('fetchStart', (e: Event) => {
        isFetching.value = true
      })
      window.addEventListener('fetchEnd', (e: Event) => {
        isFetching.value = false
      })
      window.addEventListener('popstate', (evt) => {
        if (evt && evt.state != null) {
          _plaid().onpopstate(evt)
        }
      })
    })

    return {
      current
    }
  },
  template: `
      <div id="app" v-cloak>
        <component :is="current"></component>
      </div>
    `
})

export const plaidPlugin = {
  install(app: App) {
    app.component('GoPlaidScope', GoPlaidScope)
    app.component('GoPlaidPortal', GoPlaidPortal)
    app.component('GoPlaidRunScript', GoPlaidRunScript)
    app.component('UserComponent', UserComponent)
    app.directive('keep-scroll', keepScroll)
    app.directive('assign', assignOnMounted)
    app.component('GlobalEvents', GlobalEvents)
  }
}

export interface FlashMessages {
  parent: FlashMessages
  items: any[]
}

export function createWebApp(template: string): App<Element> {
  const app = createApp(Root, {
    initialTemplate: template
  })

  const copiedToClipboard = ref(false)
  const copyToClipboard = (text: string) => {
    window.navigator.clipboard.writeText(text)
    copiedToClipboard.value = true
    window.setTimeout(() => {
      copiedToClipboard.value = false
    }, 1000)
  }

  app.config.globalProperties.copyToClipboard = copyToClipboard
  app.config.globalProperties.copiedToClipboard = copiedToClipboard
  app.config.globalProperties.presetsListing = null
  app.config.globalProperties.presetsDetailing = null
  app.config.globalProperties.presetsCreating = null
  app.config.globalProperties.presetsEditing = null
  app.use(plaidPlugin)
  return app
}
