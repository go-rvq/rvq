<template>
  <slot v-if="props.scope" v-bind="localScope"></slot>
  <slot v-else></slot>
</template>
<script setup lang="ts">
import {
  computed,
  inject,
  isProxy,
  isReactive,
  isRef,
  onBeforeUnmount,
  onMounted,
  provide,
  reactive,
  ref,
  watch
} from 'vue'

import debounce from 'lodash/debounce'

declare let window: any

const props = defineProps({
  scope: {
    type: Object,
    default: () => {}
  },
  assign: {
    type: Array<Array<any>>,
    default: () => []
  },
  setup: {
    type: Array<Function>,
    default: []
  }
})

props.assign.forEach((v) => {
  Object.assign(v[0], v[1])
})

const emit = defineEmits(['mounted', 'beforeUnmount'])

let localScope = ref({})

if (props.scope) {
  localScope = ref(props.scope)
}

const arg = {
  watch,
  isProxy,
  isReactive,
  isRef,
  ref,
  reactive,
  debounce,
  computed,
  inject,
  provide,
  scope: props.scope,
  $scope: localScope,
  window
}

if (props.setup) {
  const func = (f: Function) => {
    f({ ...arg, func: () => func(f) })
  }

  props.setup.forEach(func)
}

onMounted(() => {
  emit('mounted', arg)
})

onBeforeUnmount(() => {
  emit('beforeUnmount', arg)
})
</script>
