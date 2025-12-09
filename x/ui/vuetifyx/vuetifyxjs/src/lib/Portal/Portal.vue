<template>
  <slot name="loading" v-if="loadingModel" />
  <template v-else-if="noEmpty">
    <component v-if="props.tag" :is="props.tag" v-bind="rootProps" ref="containerEl">
      <component :is="current" v-if="noEmpty">
        <slot v-bind="internalData"></slot>
      </component>
    </component>
    <template v-else>
      <component :is="current" v-if="noEmpty">
        <slot v-bind="internalData"></slot>
      </component>
    </template>
  </template>
</template>

<script lang="ts">
export { type RenderOptions } from '@/lib/component-by-template'
import {
  buildComponentData,
  Callbacks,
  componentByTemplate,
  type RenderOptions
} from '@/lib/component-by-template'
import { type ContentType, UpdateEvent, makeProps, Portal, PortalContent, usePortals } from './types'
</script>

<script setup lang="ts">
import {
  type DefineComponent,
  onMounted,
  ref,
  shallowRef,
  useSlots,
  watch, onUnmounted, useTemplateRef
} from 'vue'
import { ComponentDefinition } from '@/lib/component-by-template'
import { definePortal } from '@/lib/Portal/types'

const containerElement = useTemplateRef('containerEl')
const props = defineProps(makeProps())

const emit = defineEmits<{
  mounted: [e: UpdateEvent],
  unmounted: [e: UpdateEvent],
  setup: [e: UpdateEvent],
  preContentUpdate: [e: UpdateEvent],
  postContentUpdate: [e: UpdateEvent],
  change: [e: UpdateEvent],
}>()

const rootProps = {
  id: props.id,
  class: props.class,
  style: props.style,
  'vx-portal': props.name || '',
  ...(props.attrs || {})
}

const model = defineModel()
const loadingModel = defineModel('loading')

const current = shallowRef<DefineComponent | null>(null)
const slots = useSlots()
const noEmpty = ref(true)
const internalData = buildComponentData({ options: { data: props.data || {} } } as ComponentDefinition)
const slotTemplate = '<slot v-bind="vx$scope"></slot>'

const mergeRenderOptions = (options: RenderOptions = {}) => {
  options.mounted ||= []
  options.unmounted ||= []
  options.setup ||= []
  options.data ??= {}

  const od = options.data as Record<string, any>

  if (props.data) {
    const pd = props.data as Record<string, any>
    Object.keys(props.data).forEach(name => {
      od[name] = Reflect.get(pd, name)
    })
  }

  return options
}

const initialValue = (() => {
  if (model.value) {
    return model.value as ContentType
  }

  if (slots.default) {
    return slotTemplate as ContentType
  }

  return null as ContentType
})()

const portals = usePortals()
const dotPortal = definePortal({
  name: props.name || '',
  content: model,
  containerElement: null,
  elements: []
} as Portal)

let internalChanging = false

const internalSetContent = (content: ContentType) => {
  const event = {
    portal: dotPortal,
    contentArg: content
  } as UpdateEvent

  emit('preContentUpdate', event)

  let c = { template: '', options: {} } as PortalContent

  if (content) {
    if (typeof content === 'string') {
      c.template = content as string
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
    } else {
      const initialValueType = typeof initialValue
      if (initialValueType === 'string') {
        c.template = initialValue as string
      } else {
        // is a portal content
        const cc = content as PortalContent
        // set template value or must clean
        c.template = cc.template || ''
        c.options = cc.options ?? {}
      }
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
    emit('postContentUpdate', event)
    emit('change', event)
    return
  }

  if (model.value !== c.template) {
    internalChanging = true
    model.value = c.template
  }

  // unreactives options
  const data: Record<string, any> = { },
    co = c.options as RenderOptions

  if (co.data) {
    const cod = co.data as Record<string, any>
    Object.keys(cod).forEach(name => {
      data[name] = Reflect.get(cod, name)
    })
  }

  const options = {
    ...c.options,
    data
  }

  noEmpty.value = true
  mergeRenderOptions(options as RenderOptions)

  const mountedHandlers: Callbacks = []

  mountedHandlers.push(({locals}) => {
    if(containerElement.value) {
      dotPortal.containerElement = (containerElement.value as any) as HTMLElement
    }
    dotPortal.elements = locals.nodes
  })

  if (c.template === slotTemplate) {
    mountedHandlers.push(({ locals }) => {
      if (dotPortal.elements?.length == 1 && (dotPortal?.elements[0] as HTMLElement).tagName === 'TEMPLATE') {
        let fragmentHTML = ''
        dotPortal.elements[0].childNodes.forEach(node => {
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

  options.setup ||= []
  options.setup.unshift((arg) => {
    arg.data.vx$portal = dotPortal
    emit('setup', { ...event, ...arg } as UpdateEvent)
  })

  options.mounted ||= []
  options.mounted.unshift(...mountedHandlers)
  options.mounted.push((arg) => emit('mounted', { ...event, ...arg } as UpdateEvent))

  options.unmounted ||= []
  options.unmounted.push((arg) => emit('unmounted', { ...event, ...arg } as UpdateEvent))

  options.mounted.push((arg) => {
    new Promise(resolve => {
      resolve(true)
      event.content = c
      emit('postContentUpdate', { ...event, ...arg } as UpdateEvent)
      emit('change', { ...event, ...arg } as UpdateEvent)
    })
  })

  current.value = componentByTemplate({
    template: c.template,
    options: options
  })
}

const setContent = (content: ContentType) => {
  loadingModel.value = false
  internalSetContent(content)
}

defineExpose(dotPortal)

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
</script>
