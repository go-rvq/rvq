import {vi} from 'vitest'
import {mount, VueWrapper} from '@vue/test-utils'
import {plaidPlugin, Root} from '../app'
import {type Ref} from 'vue'
import {createVuetify} from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

const vuetify = createVuetify({
  components,
  directives
})

global.ResizeObserver = require('resize-observer-polyfill')

// jsdom has no visualViewport (the overlays of Vuetify read it): the window's
if (!(globalThis as any).visualViewport) {
  ;(globalThis as any).visualViewport = {
    width: window.innerWidth,
    height: window.innerHeight,
    offsetLeft: 0,
    offsetTop: 0,
    pageLeft: 0,
    pageTop: 0,
    scale: 1,
    addEventListener: () => {},
    removeEventListener: () => {}
  }
}

export function mockFetchWithReturnTemplate(requestedForm: Ref<FormData>, responseFuncOrJSON: any) {
  global.fetch = vi.fn().mockImplementation((url, opts) => {
    requestedForm.value = opts.body
    let r = responseFuncOrJSON
    if (typeof responseFuncOrJSON === 'function') {
      r = responseFuncOrJSON(url, opts)
    }
    return Promise.resolve(new Response(JSON.stringify(r)))
  })
}

export function mountTemplate(template: string, global = {}): VueWrapper {
  return mount(Root, {
    props: {
      initialTemplate: template
    },
    global: {
      ...global,
      plugins: [plaidPlugin, vuetify]
    }
  })
}

export function waitUntil(
  conditionFunction: () => boolean,
  { timeout = 3000, interval = 50 } = {}
) {
  const startTime = Date.now()

  return new Promise<void>((resolve, reject) => {
    const checkCondition = async () => {
      if (conditionFunction()) {
        resolve()
      } else if (Date.now() - startTime > timeout) {
        reject(new Error('Condition not met within timeout'))
      } else {
        setTimeout(checkCondition, interval)
      }
    }

    checkCondition()
  })
}
