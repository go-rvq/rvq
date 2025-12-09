<template>
  <template v-if="raw">
    <component :is="current" v-if="current">
      <slot :form="form" :locals="locals" v-bind="internalScope"></slot>
    </component>
  </template>
  <template v-else>
    <div
      class="go-plaid-portal"
      v-bind="$attrs"
      :id="portalName && 'portal--' + portalName"
      v-if="visible"
      ref="portal"
    >
      <component :is="current" v-if="current">
        <slot :form="form" :locals="locals" v-bind="internalScope"></slot>
      </component>
    </div>
  </template>
</template>

<script setup lang="ts">
import {
  type DefineComponent,
  inject,
  isProxy,
  onBeforeUnmount,
  onMounted,
  onUpdated,
  provide,
  reactive,
  ref,
  shallowRef,
  useSlots
} from 'vue'

import { componentByTemplate, type RenderOptions, type RenderSelf } from '@/component-by-template'
import type { EventResponse } from '@/types'
import type { Builder } from '@/builder'

declare let window: any
window.__goplaid = window.__goplaid ?? {}
window.__goplaid.portals = window.__goplaid.portals ?? {}
window.__goplaid.index = window.__goplaid.index ?? 0

const portal = ref()

const props = defineProps<{
  loader: Builder | undefined
  content: string
  visible: boolean
  afterLoaded: () => null
  portalName: string
  form: object | undefined
  locals: object | undefined
  methods: object | undefined
  data: object | undefined
  autoReloadInterval: string | number
  raw: boolean
  scope: any
}>()

const emit = defineEmits<{
  mounted: [self:RenderSelf],
  unmounted: [self:RenderSelf]
}>()

const current = shallowRef<DefineComponent | null>(null)
const autoReloadIntervalID = ref<number>(0)
const slots = useSlots()

const name = props.portalName || 'anonymous$' + window.__goplaid.index++

let locals = inject<object>('locals', {})

if (props.locals !== undefined) {
  if (!isProxy(props.locals)) {
    locals = reactive({ $parent: locals, ...props.locals })
  }
}
provide('locals', locals)

let form = inject<object>('form', {})

if (props.form !== undefined) {
  if (!isProxy(props.form)) {
    form = reactive({ $parent: form, ...props.form })
  }
}
provide('form', form)

const internalScope = { form, locals, ...(props.scope ?? {}) }

const mergeRenderOptions = (options:RenderOptions={}) => {
  options.mounted ||= []
  options.unmounted ||= []

  options.mounted.push((self:RenderSelf) => emit('mounted', self))
  options.unmounted.push((self:RenderSelf) => emit('unmounted', self))

  return options
}

const updatePortalTemplate = (template: string, options:RenderOptions={}) => {
  current.value = componentByTemplate(template, internalScope, portal, mergeRenderOptions(options))
}

// other reactive properties and methods
const reload = () => {
  if (slots.default) {
    current.value = componentByTemplate('<slot v-bind="SCOPE"></slot>', internalScope, portal, mergeRenderOptions({}))
    return
  }

  const content = props.content
  if (content) {
    updatePortalTemplate(content)
    return
  }

  const ef = props.loader
  if (!ef) {
    return
  }
  ef.loadPortalBody(true)
    .go()
    .then((r: EventResponse) => {
      updatePortalTemplate(r.body)
    })
}

onMounted(() => {
  window.__goplaid.portals[name] = {
    portalName: name,
    updatePortalTemplate,
    reload
  }
  reload()
})

onUpdated(() => {
  if (props.autoReloadInterval && autoReloadIntervalID.value == 0) {
    const interval = parseInt(props.autoReloadInterval + '')
    if (interval == 0) {
      return
    }

    autoReloadIntervalID.value = setInterval(() => {
      reload()
    }, interval) as unknown as number
  }

  if (
    autoReloadIntervalID.value &&
    autoReloadIntervalID.value > 0 &&
    props.autoReloadInterval == 0
  ) {
    clearInterval(autoReloadIntervalID.value)
    autoReloadIntervalID.value = 0
  }
})

onBeforeUnmount(() => {
  delete window.__goplaid.portals[name]
  if (autoReloadIntervalID.value && autoReloadIntervalID.value > 0) {
    clearInterval(autoReloadIntervalID.value)
  }
})
</script>
