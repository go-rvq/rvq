<script setup lang="ts">
import * as vue from 'vue'
import { onBeforeUnmount, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  script: {
    type: Function,
    required: true
  },
  setup: {
    type: Function,
  },
  beforeUnmount: {
    type: Function,
  },
  unmounted: {
    type: Function,
  },
})

declare let window: any

const scope = {...vue, view: window}

props.setup && props.setup(scope)

onMounted(() => {
  props.script(vue)
})

onBeforeUnmount(() => {
  props.beforeUnmount && props.beforeUnmount(scope)
})

onUnmounted(() => {
  props.unmounted && props.unmounted(scope)
})
</script>
