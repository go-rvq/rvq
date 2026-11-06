<template>
  <v-dialog :model-value="true" fullscreen :scrim="false" :z-index="3200" transition="fade-transition" data-lightbox-dialog @update:model-value="close">
    <div
      class="vx-lightbox"
      tabindex="0"
      ref="root"
      @keydown="onKey"
      @click.self="close"
      @touchstart.passive="touchStart"
      @touchend.passive="touchEnd"
    >
      <div class="vx-lightbox-bar">
        <span v-if="images.length > 1" class="vx-lightbox-count" data-lightbox-count>{{ i + 1 }} / {{ images.length }}</span>
        <v-spacer />
        <v-btn icon="mdi-open-in-new" variant="text" color="white" :href="current.src" target="_blank" :title="words.lightboxOpen" />
        <v-btn icon="mdi-close" variant="text" color="white" :title="words.lightboxClose" data-lightbox-close @click="close" />
      </div>
      <div class="vx-lightbox-stage" @click.self="close">
        <div v-if="images.length > 1" class="vx-lightbox-nav prev">
          <v-btn icon="mdi-chevron-left" variant="tonal" color="white" size="large" :title="words.lightboxPrev" data-lightbox-prev @click="go(-1)" />
        </div>
        <img :key="current.src" :src="current.src" :alt="current.caption || ''" class="vx-lightbox-img" data-lightbox-img />
        <div v-if="images.length > 1" class="vx-lightbox-nav next">
          <v-btn icon="mdi-chevron-right" variant="tonal" color="white" size="large" :title="words.lightboxNext" data-lightbox-next @click="go(1)" />
        </div>
      </div>
      <div v-if="current.caption" class="vx-lightbox-caption" data-lightbox-caption>{{ current.caption }}</div>
      <div v-if="images.length > 1" class="vx-lightbox-thumbs">
        <img
          v-for="(im, k) in images"
          :key="im.src + k"
          :src="im.src"
          :class="{active: k === i}"
          data-lightbox-thumb
          @click="i = k"
        />
      </div>
    </div>
  </v-dialog>
</template>

<script lang="ts" setup>
// The images of a group, one at a time, over the page: the arrows (and the
// keys ← →, a swipe) go through them, Esc closes.
import {computed, nextTick, onMounted, ref} from 'vue'
import {useVxMessages} from '@/lib/locale'
import type {LightboxImage} from './lightbox'

const props = defineProps<{images: LightboxImage[]; index?: number; onClose?: () => void}>()

const words = useVxMessages()
const i = ref(Math.min(Math.max(0, props.index ?? 0), props.images.length - 1))
const current = computed(() => props.images[i.value])
const root = ref<HTMLElement | null>(null)

function go(d: number) {
  const n = props.images.length
  i.value = (i.value + d + n) % n
}

function close() {
  props.onClose?.()
}

function onKey(ev: KeyboardEvent) {
  if (ev.key === 'ArrowLeft') go(-1)
  else if (ev.key === 'ArrowRight') go(1)
  else if (ev.key === 'Escape') close()
  else return
  ev.preventDefault()
  ev.stopPropagation()
}

let x0 = 0
function touchStart(ev: TouchEvent) {
  x0 = ev.touches[0]?.clientX ?? 0
}
function touchEnd(ev: TouchEvent) {
  const dx = (ev.changedTouches[0]?.clientX ?? 0) - x0
  if (Math.abs(dx) > 50 && props.images.length > 1) go(dx < 0 ? 1 : -1)
}

onMounted(() => nextTick(() => root.value?.focus()))

defineExpose({go, i})
</script>

<style>
/* over every dialog: the stack of Vuetify puts a dialog mounted apart just
   over the topmost one, whatever its z-index */
.v-overlay:has(> .v-overlay__content > .vx-lightbox),
.v-overlay:has(.vx-lightbox) {
  z-index: 3200 !important;
}
.vx-lightbox {
  position: fixed;
  inset: 0;
  display: flex;
  flex-direction: column;
  background: rgba(0, 0, 0, 0.92);
  color: #fff;
  outline: none;
}
.vx-lightbox-bar {
  display: flex;
  align-items: center;
  padding: 4px 8px;
}
.vx-lightbox-count {
  padding-left: 12px;
  font-size: 0.9rem;
  opacity: 0.85;
}
.vx-lightbox-stage {
  position: relative;
  flex: 1;
  min-height: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 64px;
}
.vx-lightbox-img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.6);
  user-select: none;
}
.vx-lightbox-nav {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
}
.vx-lightbox-nav.prev {
  left: 12px;
}
.vx-lightbox-nav.next {
  right: 12px;
}
.vx-lightbox-caption {
  text-align: center;
  padding: 8px 16px 0;
  font-size: 0.9rem;
  opacity: 0.9;
}
.vx-lightbox-thumbs {
  display: flex;
  justify-content: center;
  gap: 6px;
  padding: 10px;
  overflow-x: auto;
}
.vx-lightbox-thumbs img {
  height: 56px;
  width: 56px;
  object-fit: cover;
  border-radius: 4px;
  opacity: 0.5;
  cursor: pointer;
  border: 2px solid transparent;
}
.vx-lightbox-thumbs img.active {
  opacity: 1;
  border-color: rgb(var(--v-theme-primary));
}
</style>
