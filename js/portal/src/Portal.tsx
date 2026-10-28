// <rvq-portal>: a place of the page whose content is a Vue template given at
// run time (README.md) — compiled and shown there, again each time it
// changes; its root elements followed while it is shown (elements).
import {
  type DefineComponent,
  defineComponent,
  nextTick,
  onMounted,
  onUnmounted,
  ref,
  shallowRef,
  useModel,
  watch
} from 'vue'

import {
  buildComponentData,
  linkFields,
  type Callbacks,
  type ComponentDefinition,
  componentByTemplate,
  type RenderOptions
} from './component-by-template'
import {
  type ContentType,
  definePortal,
  makeProps,
  type Portal as PortalObject,
  type PortalContent,
  type UpdateEvent,
  usePortals
} from './types'

export type { RenderOptions } from './component-by-template'

// the content of the default slot, when it is the initial content: the
// slot itself, with the portal's variables
const slotTemplate = '<slot v-bind="vx$scope"></slot>'

export default defineComponent({
  name: 'RvqPortal',
  props: makeProps(),
  emits: {
    'update:modelValue': (_v: ContentType) => true,
    'update:loading': (_v: boolean) => true,
    mounted: (_e: UpdateEvent) => true,
    unmounted: (_e: UpdateEvent) => true,
    setup: (_e: UpdateEvent) => true,
    preContentUpdate: (_e: UpdateEvent) => true,
    postContentUpdate: (_e: UpdateEvent) => true,
    change: (_e: UpdateEvent) => true,
    // the root elements of the content changed while it is shown (a v-if, a
    // v-for at its root): portal.elements the new, added and removed
    elements: (_e: UpdateEvent) => true
  },
  setup(props, { emit, slots, expose }) {
    const containerElement = ref<HTMLElement | null>(null)

    const rootProps = {
      id: props.id,
      class: props.class,
      style: props.style,
      'rvq-portal': props.name || '',
      ...(props.attrs || {})
    }

    // v-model and v-model:loading — local when the page binds none
    const model = useModel(props, 'modelValue')
    const loadingModel = useModel(props, 'loading')

    const current = shallowRef<DefineComponent | null>(null)
    const noEmpty = ref(true)
    const internalData = buildComponentData({ options: { data: props.data || {} } } as ComponentDefinition)

    const mergeRenderOptions = (options: RenderOptions = {}) => {
      options.mounted ||= []
      options.updated ||= []
      options.unmounted ||= []
      options.setup ||= []
      options.data ??= {}

      const od = options.data as Record<string, any>

      if (props.data) {
        linkFields(od, props.data as Record<string, any>)
      }

      return options
    }

    const initialValue = ((): ContentType => {
      if (model.value) {
        return model.value as ContentType
      }
      if (slots.default) {
        return slotTemplate
      }
      return null
    })()

    const portals = usePortals()
    const dotPortal = definePortal({
      name: props.name || '',
      content: model,
      containerElement: null,
      elements: []
    } as unknown as PortalObject)

    let internalChanging = false

    const internalSetContent = (content: ContentType) => {
      const event = {
        portal: dotPortal,
        contentArg: content
      } as UpdateEvent

      emit('preContentUpdate', event)

      const c = { template: '', options: {} } as PortalContent

      if (content) {
        if (typeof content === 'string') {
          c.template = content
        } else {
          const cc = content as PortalContent
          c.template = cc.template
          c.options = cc.options ?? {}
        }
      } else if (content === null || content === undefined) {
        c.template = undefined
      }

      if (c.template === undefined) {
        // reset to initial value
        if (!initialValue) {
          // must empty
          c.template = ''
        } else if (typeof initialValue === 'string') {
          c.template = initialValue
        } else {
          // is a portal content: its template, or empty
          const cc = initialValue as PortalContent
          c.template = cc.template || ''
          c.options = cc.options ?? {}
        }
      }

      if (!c.template) {
        c.template = ''

        if (model.value !== c.template) {
          internalChanging = true
          model.value = c.template
        }

        // must empty
        noEmpty.value = false
        current.value = null

        event.content = c
        // once the content shown is gone (its unmounted first)
        nextTick(() => {
          emit('postContentUpdate', event)
          emit('change', event)
        })
        return
      }

      if (model.value !== c.template) {
        internalChanging = true
        model.value = c.template
      }

      // the options' own data: the content's fields linked (its reactivity
      // kept), what is set here (the portal's data, vx$portal) its own
      const data: Record<string, any> = {},
        co = c.options as RenderOptions

      if (co.data) {
        linkFields(data, co.data as Record<string, any>)
      }

      const options = {
        ...c.options,
        data
      } as RenderOptions

      noEmpty.value = true
      mergeRenderOptions(options)

      const mountedHandlers: Callbacks = []

      mountedHandlers.push(({ locals }) => {
        if (containerElement.value) {
          dotPortal.containerElement = containerElement.value
        }
        dotPortal.elements = locals.nodes
      })

      if (c.template === slotTemplate) {
        // a slot of <template v-pre>: its HTML is the template of the portal
        mountedHandlers.push(({ locals }) => {
          if (dotPortal.elements?.length == 1 && (dotPortal.elements[0] as HTMLElement).tagName === 'TEMPLATE') {
            let fragmentHTML = ''
            dotPortal.elements[0].childNodes.forEach((node) => {
              if (node.nodeType === Node.ELEMENT_NODE) {
                fragmentHTML += (node as HTMLElement).outerHTML
              } else if (node.nodeType === Node.TEXT_NODE) {
                fragmentHTML += node.textContent
              }
            })
            locals.skipUnmountHandlers = true
            setTimeout(() => {
              internalSetContent(fragmentHTML)
            }, 1)
            return false
          }
        })
      }

      options.setup!.unshift((arg) => {
        arg.data.vx$portal = dotPortal
        emit('setup', { ...event, ...arg } as UpdateEvent)
      })

      options.mounted!.unshift(...mountedHandlers)
      options.mounted!.push((arg) => emit('mounted', { ...event, ...arg } as UpdateEvent))

      options.updated!.push((arg) => {
        const previous = arg.previous || []
        dotPortal.elements = arg.locals.nodes
        emit('elements', {
          ...event,
          ...arg,
          added: arg.locals.nodes.filter((n) => !previous.includes(n)),
          removed: previous.filter((n) => !arg.locals.nodes.includes(n))
        } as UpdateEvent)
      })

      options.unmounted!.push((arg) => {
        emit('unmounted', { ...event, ...arg } as UpdateEvent)
        // its elements gone — unless the next content's are already there
        if (dotPortal.elements === arg.locals.nodes) {
          dotPortal.elements = []
        }
      })

      options.mounted!.push((arg) => {
        Promise.resolve().then(() => {
          event.content = c
          emit('postContentUpdate', { ...event, ...arg } as UpdateEvent)
          emit('change', { ...event, ...arg } as UpdateEvent)
        })
      })

      current.value = componentByTemplate({
        template: c.template,
        options
      })
    }

    const setContent = (content: ContentType) => {
      loadingModel.value = false
      internalSetContent(content)
    }

    expose(dotPortal)

    watch(model, (v) => {
      if (internalChanging) {
        internalChanging = false
        return
      }
      setContent(v as ContentType)
    })

    onMounted(() => {
      if (props.name && portals) {
        portals[props.name] = dotPortal
      }
      setContent(null)
    })

    onUnmounted(() => {
      if (props.name && portals) {
        delete portals[props.name]
      }
    })

    return () => {
      if (loadingModel.value) {
        return slots.loading?.()
      }
      if (!noEmpty.value || !current.value) {
        return null
      }
      // the content, the default slot of the page given to it (<slot> in its
      // template), with the portal's variables
      const Content = current.value as any
      const content = <Content v-slots={{ default: () => slots.default?.(internalData) }} />
      if (props.tag) {
        const Tag = props.tag as any
        return (
          <Tag {...rootProps} ref={containerElement}>
            {content}
          </Tag>
        )
      }
      return content
    }
  }
})
