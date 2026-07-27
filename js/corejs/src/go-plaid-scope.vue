<template>
  <slot
    :locals="locals"
    :form="form"
    :plaid="plaid"
    :vars="vars"
    :closer="closer"
    :setup="setup"
    :fullscreen="fullscreen"
  ></slot>
</template>

<script setup lang="ts">
import {inject, isProxy, onMounted, provide, reactive, watch} from 'vue'
import {createCloser} from '@/closer'
import debounce from 'lodash/debounce'

const props = defineProps<{
  scopeName?: string
  locals?: undefined | object | any[]
  form?: undefined | object | any[]
  closer?: undefined | object | any[]
  fullscreen?: undefined | object | any[]
  setup?: undefined | object | any[]
  useDebounce?: number
  observers?: { name: string; script: string }[]
}>()

const emit = defineEmits<{
  (e: 'change-debounced', obj: object): void
}>()

let locals = inject<object>('locals', {})

if (props.locals !== undefined) {
  let dot: object = { $parent: locals }
  if (isProxy(props.locals)) {
    dot = props.locals
  } else if (Array.isArray(props.locals)) {
    dot = reactive(Object.assign(dot, ...props.locals))
  } else {
    dot = reactive({ ...dot, ...props.locals })
  }
  locals = dot
  provide('locals', locals)
}

let form = inject<object>('form' || {})

if (props.form !== undefined) {
  let dot: object = { $parent: form }
  if (isProxy(props.form)) {
    dot = props.form
  } else if (Array.isArray(props.form)) {
    dot = reactive(Object.assign(dot, ...props.form))
  } else {
    dot = reactive({ ...dot, ...props.form })
  }
  form = dot
  provide('form', form)
}

let closer = inject<object>('closer', {})

if (props.closer !== undefined) {
  let dot: object = { $parent: closer, show: false }
  if (isProxy(props.closer)) {
    // an existing closer (a form host owns it, see ParamCloserProvided): adopt
    // it as is — re-wrapping would give this scope a second, silent copy of the
    // open/close callbacks.
    dot = props.closer
  } else if (Array.isArray(props.closer)) {
    // createCloser LAST, over the merged initializer: `show` in it is the
    // initial state and must not fire anything.
    dot = createCloser(Object.assign(dot, ...props.closer))
  } else {
    dot = createCloser({ ...dot, ...props.closer })
  }
  closer = dot
  provide('closer', closer)
}

let fullscreen = inject<object>('fullscreen', {})

if (props.fullscreen !== undefined) {
  let dot: object = { $parent: fullscreen, active: false }
  if (isProxy(props.fullscreen)) {
    dot = props.fullscreen
  } else if (Array.isArray(props.fullscreen)) {
    dot = reactive(Object.assign(dot, ...props.fullscreen))
  } else {
    dot = reactive({ ...dot, ...props.fullscreen })
  }
  fullscreen = dot
  provide('fullscreen', fullscreen)
}

const vars = inject<{ __notification?: { id: string; name: string; payload: any } }>('vars')
const plaid = inject('plaid')

interface Closer {
  show: boolean
}

if (Array.isArray(props.setup)) {
  const setupFn = new Function(
    'vars',
    'locals',
    'form',
    'plaid',
    'closer',
    'fullscreen',
    props.setup[0]
  )
  setupFn(vars, locals, form, plaid, closer, fullscreen)
}

function addObservers() {
  if (!props.observers || props.observers.length == 0) {
    return
  }
  watch(
    () => vars?.__notification,
    (newNotification) => {
      if (!newNotification) {
        return
      }
      props.observers?.forEach((observer) => {
        if (newNotification?.name === observer.name) {
          let payload
          try {
            payload =
              typeof newNotification.payload === 'string'
                ? JSON.parse(newNotification.payload)
                : newNotification.payload
          } catch (e) {
            payload = newNotification.payload
          }
          try {
            const scriptFunc = new Function(
              'name',
              'payload',
              'vars',
              'locals',
              'form',
              'plaid',
              'closer',
              'fullscreen',
              observer.script
            )
            scriptFunc(observer.name, payload, vars, locals, form, plaid, closer, fullscreen)
          } catch (error) {
            console.error('Error executing observer script:', error)
          }
        }
      })
    }
  )
}

onMounted(() => {
  setTimeout(() => {
    if (props.useDebounce) {
      const debounceWait = props.useDebounce
      const _watch = debounce((obj: any) => {
        emit('change-debounced', obj)
      }, debounceWait)
      watch(locals, (value, oldValue) => {
        _watch({
          locals: value,
          form: form,
          oldLocals: oldValue,
          oldForm: form
        })
      })
      /*watch(form, (value, oldValue) => {
        _watch({
          locals: locals,
          form: value,
          oldLocals: locals,
          oldForm: oldValue
        })
      })*/
    }
  }, 0)

  addObservers()
})
</script>
