<script setup lang="ts">
import { onMounted, onUpdated, nextTick, ref, watch } from 'vue'

// A clickable partial-revert diff. The two sides are HTML with change regions
// tagged <... class="hist-del|hist-ins" data-h="<hunk index>">. Clicking a region
// toggles its hunk index in the model value (super-highlight + check on all
// regions of that hunk). The selected indexes drive the revert.

const props = defineProps<{
  instruction?: string
  leftLabel?: string
  rightLabel?: string
  leftHtml?: string
  rightHtml?: string
}>()

const model = defineModel<number[]>({ default: () => [] })
const root = ref<HTMLElement | null>(null)

function sync() {
  const el = root.value
  if (!el) return
  const sel = new Set(model.value || [])
  el.querySelectorAll<HTMLElement>('[data-h]').forEach((node) => {
    const h = Number(node.dataset.h)
    node.classList.toggle('hist-sel', sel.has(h))
  })
}

function onClick(e: MouseEvent) {
  const target = (e.target as HTMLElement | null)?.closest('[data-h]') as HTMLElement | null
  if (!target) return
  const h = Number(target.dataset.h)
  const sel = new Set(model.value || [])
  if (sel.has(h)) sel.delete(h)
  else sel.add(h)
  model.value = Array.from(sel).sort((a, b) => a - b)
}

onMounted(() => nextTick(sync))
onUpdated(() => sync())
watch(model, () => nextTick(sync))
watch(() => [props.leftHtml, props.rightHtml], () => nextTick(sync))
</script>

<template>
  <div ref="root">
    <v-alert v-if="instruction" type="info" variant="tonal" density="compact" class="mb-3">
      {{ instruction }}
    </v-alert>
    <div class="d-flex align-start" @click="onClick">
      <div class="flex-1-1-0 pe-2 vx-hunk-col">
        <div class="text-caption text-medium-emphasis mb-1">{{ leftLabel }}</div>
        <div class="vx-hunk-side" v-html="leftHtml"></div>
      </div>
      <div class="flex-1-1-0 ps-2">
        <div class="text-caption text-medium-emphasis mb-1">{{ rightLabel }}</div>
        <div class="vx-hunk-side" v-html="rightHtml"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.vx-hunk-col {
  border-right: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}
.vx-hunk-side {
  white-space: pre-wrap;
  word-break: break-word;
}
/* Base change colors (theme). */
.vx-hunk-side :deep(.hist-del) {
  background: rgba(var(--v-theme-error), 0.14);
  text-decoration: line-through;
}
.vx-hunk-side :deep(.hist-ins) {
  background: rgba(var(--v-theme-success), 0.14);
}
.vx-hunk-side :deep(.hist-del),
.vx-hunk-side :deep(.hist-ins) {
  cursor: pointer;
  border-radius: 3px;
  position: relative;
  transition: box-shadow 0.1s ease;
}
/* Selected: super-highlight + a check to the left of the region. */
.vx-hunk-side :deep(.hist-sel) {
  outline: 2px solid rgb(var(--v-theme-primary));
  box-shadow: 0 0 0 3px rgba(var(--v-theme-primary), 0.25);
  padding-left: 1.3em;
}
.vx-hunk-side :deep(.hist-sel)::before {
  content: '✓';
  position: absolute;
  left: 0.15em;
  top: 0;
  color: rgb(var(--v-theme-primary));
  font-weight: 700;
}
</style>
