<script setup lang="ts">
// vx-layout-watch: says the layout of the browser's width — a phone's
// (narrower than breakpoint) or a wide screen's — when it is not the one the
// server drew (mobile): at once, and each time the width crosses. Nothing to
// see: the list editors draw their form again in the layout asked (change).
import {onBeforeUnmount, onMounted} from 'vue'

const props = withDefaults(
  defineProps<{
    // the layout drawn: a phone's
    mobile?: boolean
    // a phone's: narrower than this (px)
    breakpoint?: number
  }>(),
  {mobile: false, breakpoint: 960},
)

const emit = defineEmits<{change: [mobile: boolean]}>()

const query = typeof window !== 'undefined' && window.matchMedia ? window.matchMedia(`(max-width: ${props.breakpoint - 0.02}px)`) : null

function check(matches: boolean) {
  if (matches !== props.mobile) emit('change', matches)
}

const onChange = (ev: MediaQueryListEvent) => check(ev.matches)

// at once — but only once the page is all mounted: each input puts its value
// in the form when mounted (v-assign), and a form asked again before the last
// of them would be posted without it (a row drawn after this watch)
let frame = 0
let gone = false
const later = (f: () => void) =>
  typeof requestAnimationFrame === 'function' ? (frame = requestAnimationFrame(() => (frame = requestAnimationFrame(f)))) : setTimeout(f, 0)

onMounted(() => {
  later(() => {
    if (!gone && query) check(query.matches)
  })
  query?.addEventListener?.('change', onChange)
})
onBeforeUnmount(() => {
  gone = true
  if (frame && typeof cancelAnimationFrame === 'function') cancelAnimationFrame(frame)
  query?.removeEventListener?.('change', onChange)
})

defineExpose({check})
</script>

<template>
  <span style="display: none" data-layout-watch></span>
</template>
