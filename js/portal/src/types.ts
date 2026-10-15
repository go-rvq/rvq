import {getCurrentInstance, inject, ModelRef, PropType, provide} from 'vue'

import {Data, Locals, type RenderOptions} from './component-by-template'

const portalKey = Symbol("vx$portal")

// When value is:
// - null | { template: null } : reset to initial state.
// - '' | { template: '' } : must empty
export type ContentType = string | PortalContent | undefined | null

export interface PortalContent {
  // if is undefined, reset to initial value
  // if is blank string, empty value
  template?: string
  options?: RenderOptions
}

export interface Portal {
  name: string
  content: ModelRef<ContentType>
  elements: Node[]
  containerElement:HTMLElement|null
}

export interface UpdateEvent {
  portal: Portal
  contentArg: ContentType
  content?: PortalContent
  data?: Data
  locals?: Locals
  // the event elements: the root nodes come in and gone out
  added?: Node[]
  removed?: Node[]
}

export function usePortals(): Record<string, any> | null {
  const app = getCurrentInstance()?.appContext.app
  if (!app) return null
  const props = app.config.globalProperties
  return (props.vx$portals ??= {})
}

export function definePortal(portal:Portal) {
  provide(portalKey, portal)
  return portal
}

export function usePortal() {
  return inject<Portal>(portalKey)
}

// The props of vx-portal (README.md): the content (v-model), loading
// (v-model:loading: the slot "loading" while true), data (the variables of
// the content's templates and the props of the default slot), tag (the
// element the content is put in), name (vx$portals[name]).
export function makeProps() {
  return {
    modelValue: {
      type: Object as PropType<ContentType>
    },
    loading: {
      type: Object as PropType<Boolean | undefined>,
      default: false
    },
    data: Object,
    tag: String,
    attrs: {
      type: Object as PropType<Record<string, any>>
    },
    id: String,
    name: String,
    // the class and the style of the element of tag
    class: [String, Array, Object] as PropType<any>,
    style: [String, Array, Object] as PropType<any>
  }
}
