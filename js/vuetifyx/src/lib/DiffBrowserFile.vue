<script setup lang="ts">
import { computed, inject, onMounted } from 'vue'
import CodeView from './CodeView.vue'
import { diffBrowserKey, statusColor, unchangedContent } from './diffBrowserContext'

// A tab of <vx-diff-browser>: a file, the old and the current side by side,
// the changed lines and content marked, the two sides scrolling together. A
// renamed (moved) one says so — from where to where —, its content before
// compared with its content now; with its content the same, that is all.

const props = defineProps<{ params: { params: { path: string } } }>()
const ctx = inject(diffBrowserKey)!
// the file's summary, and its content once had (ctx.request)
const entry = computed(() => ctx.files.value.find((f) => f.path === props.params.params.path))
const content = computed(() => ctx.content(props.params.params.path))
// the diff, once the server set it: over the summary
const file = computed(() => {
  const c = content.value
  return c?.value ? { ...c.file, ...c.value } : undefined
})
onMounted(() => {
  if (entry.value) ctx.request(entry.value)
})

// gaps align the two sides, as IntelliJ does: where a block of lines was
// changed, the side with fewer lines gets blank rows after its own, so what
// follows starts on the same row on both — and the two scroll as one.
const gaps = computed(() => {
  const f = file.value
  const oldGaps: Record<number, number> = {}
  const newGaps: Record<number, number> = {}
  if (!f) return { old: oldGaps, new: newGaps }
  const oldN = (f.old ?? '').split('\n').length
  const newN = (f.new ?? '').split('\n').length
  const removed = new Set(f.removed || [])
  const added = new Set(f.added || [])
  let i = 1
  let j = 1
  while (i <= oldN || j <= newN) {
    const r0 = i <= oldN && removed.has(i)
    const a0 = j <= newN && added.has(j)
    if (!r0 && !a0) {
      i++
      j++
      continue
    }
    let r = 0
    while (i + r <= oldN && removed.has(i + r)) r++
    let a = 0
    while (j + a <= newN && added.has(j + a)) a++
    if (r < a) oldGaps[i + r] = (oldGaps[i + r] || 0) + a - r
    else if (a < r) newGaps[j + a] = (newGaps[j + a] || 0) + r - a
    i += r
    j += a
  }
  return { old: oldGaps, new: newGaps }
})

// the two sides scroll together
function syncScroll(e: Event) {
  const src = e.target as HTMLElement
  src
    .closest('.vx-diff-sides')
    ?.querySelectorAll<HTMLElement>('.vx-diff-side-body')
    .forEach((el) => {
      if (el !== src && (el.scrollTop !== src.scrollTop || el.scrollLeft !== src.scrollLeft)) {
        el.scrollTop = src.scrollTop
        el.scrollLeft = src.scrollLeft
      }
    })
}
</script>

<template>
  <div v-if="entry" class="vx-diff-file" :data-diff-file="entry.path">
    <div class="vx-diff-file-head">
      <v-chip size="x-small" variant="tonal" :color="statusColor(entry.status)" class="me-2">{{
        (entry.status || '').trim()
      }}</v-chip>
      <template v-if="entry.from"
        ><code>{{ entry.from }}</code><span class="mx-2">→</span><code>{{ entry.path }}</code></template
      >
      <code v-else>{{ entry.path }}</code>
    </div>
    <div v-if="content?.error" class="pa-4 text-error" data-diff-error>
      {{ content.error }}
      <v-btn size="small" variant="text" icon="mdi-refresh" class="ms-1" @click="ctx.request(entry)" />
    </div>
    <div v-else-if="!file" class="vx-diff-loading" data-diff-loading>
      <v-progress-linear indeterminate color="primary" />
      <div class="pa-4 text-medium-emphasis d-flex align-center ga-2">
        <v-progress-circular indeterminate size="18" width="2" color="primary" />
        <span>{{ ctx.labels.value.loading }}</span>
      </div>
    </div>
    <div v-else-if="unchangedContent(file)" class="pa-4 vx-diff-renamed" data-renamed-only>
      {{ ctx.labels.value.renamed }} <code>{{ file.path }}</code>. {{ ctx.labels.value.unchanged }}
    </div>
    <div v-else-if="file.binary" class="pa-4 text-medium-emphasis">{{ ctx.labels.value.binary }}</div>
    <div v-else class="vx-diff-sides">
      <div class="vx-diff-side">
        <div class="vx-diff-side-label">{{ ctx.labels.value.old }}</div>
        <div class="vx-diff-side-body" @scroll="syncScroll">
          <CodeView
            :code="file.old || ''"
            :language="file.language"
            mark-kind="del"
            :mark-lines="file.removed || []"
            :change-ranges="file.delRanges || {}"
            :gaps="gaps.old"
          />
        </div>
      </div>
      <div class="vx-diff-side">
        <div class="vx-diff-side-label">{{ ctx.labels.value.new }}</div>
        <div class="vx-diff-side-body" @scroll="syncScroll">
          <CodeView
            :code="file.new || ''"
            :language="file.language"
            mark-kind="add"
            :mark-lines="file.added || []"
            :change-ranges="file.insRanges || {}"
            :gaps="gaps.new"
          />
        </div>
      </div>
    </div>
  </div>
</template>
