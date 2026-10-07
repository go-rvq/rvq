// <vx-gad-ide>: the gad IDE (js/gadide) in a page of the admin. Its code —
// CodeMirror, dockview, the IDE: large — is not in this bundle: src (and css)
// are loaded the first time one is shown, once for the page, and render
// window.GadIde.component with the page's Vue and Vuetify.
import { defineComponent, h, onMounted, ref, resolveComponent, shallowRef, type Component } from 'vue'

declare const window: Window & { GadIde?: { component: Component } }

const loading: Record<string, Promise<void>> = {}

// once is src (a script, or a stylesheet with css) loaded once for the page
const once = (src: string, css: boolean): Promise<void> =>
  (loading[src] ??= new Promise<void>((resolve, reject) => {
    const el = css ? document.createElement('link') : document.createElement('script')
    if (el instanceof HTMLLinkElement) {
      el.rel = 'stylesheet'
      el.href = src
    } else {
      el.src = src
    }
    el.onload = () => resolve()
    el.onerror = () => {
      delete loading[src]
      reject(new Error('the IDE could not be loaded: ' + src))
    }
    document.head.appendChild(el)
  }))

export default defineComponent({
  name: 'VxGadIde',
  props: {
    // the IDE's script (gadide.js) and stylesheet (gadide.css)
    src: { type: String, required: true },
    css: { type: String, default: '' },
    // where its API is: base + "api/ide/…"
    base: { type: String, required: true },
    // the height it takes ("calc(100vh - 64px)")
    height: { type: String, default: '100%' },
    // the texts of its Changes and Git panels (IdeMessages), translated
    messages: { type: Object, default: undefined }
  },
  setup(props) {
    const comp = shallowRef<Component>()
    const error = ref('')
    onMounted(async () => {
      try {
        await Promise.all([once(props.src, false), props.css ? once(props.css, true) : Promise.resolve()])
        if (!window.GadIde?.component) throw new Error('the IDE did not register: ' + props.src)
        comp.value = window.GadIde.component
      } catch (e) {
        error.value = String(e instanceof Error ? e.message : e)
      }
    })
    return () => {
      if (error.value) return h('pre', { class: 'text-error pa-4', 'data-gad-ide-error': '' }, error.value)
      if (!comp.value) return h(resolveComponent('v-progress-linear'), { indeterminate: true, color: 'primary' })
      return h(comp.value, { base: props.base, height: props.height, messages: props.messages })
    }
  }
})
