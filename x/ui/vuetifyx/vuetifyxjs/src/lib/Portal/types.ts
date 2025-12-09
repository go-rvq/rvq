import { getCurrentInstance, inject, ModelRef, PropType, provide } from 'vue'

import { type RenderOptions, Data, Locals } from '@/lib/component-by-template'
import { makeComponentProps } from 'vuetify/lib/composables/component'

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
    ...makeComponentProps()
  }
}
