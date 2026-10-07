<script setup lang="ts">
import { ref } from 'vue'
import { useTheme } from 'vuetify'
import DiffBrowser from '@/lib/DiffBrowser.vue'
import type { Content, DiffFile } from '@/lib/diffBrowserContext'
// the diffs as the server computes them (vuetifyx.NewDiffFile, Go), saved
import samples from './diffBrowserSamples.json'

// <vx-diff-browser> with no server: the summary by v-model, each file's diff
// "loaded" as the server would answer — after a while, content.value set; for
// templates/broken.gadx, content.error.

const all = samples as DiffFile[]
const summaryOf = (f: DiffFile): DiffFile => ({ path: f.path, status: f.status, from: f.from, language: f.language })

// the summary: what the tree shows (v-model)
const summary = ref<DiffFile[]>(all.map(summaryOf))

// load answers a tab as the server would (its event's RunScript)
function load(content: Content) {
  const delay = 500 + Math.random() * 1000
  setTimeout(() => {
    if (content.file.path === 'templates/broken.gadx') {
      content.error = 'The server could not compare this file.'
      return
    }
    content.value = all.find((f) => f.path === content.file.path)
  }, delay)
}

// the discard of the actions slot: the file leaves the summary (its tab closes)
const discard = (file: DiffFile) => (summary.value = summary.value.filter((f) => f.path !== file.path))
const restore = () => (summary.value = all.map(summaryOf))

const theme = useTheme()
const toggleTheme = () => (theme.global.name.value = theme.global.current.value.dark ? 'light' : 'dark')
</script>

<template>
  <div>
    <div class="d-flex align-center ga-2 mb-3">
      <h3 class="me-auto">vx-diff-browser</h3>
      <v-btn size="small" variant="tonal" prepend-icon="mdi-restore" @click="restore">Restore the files</v-btn>
      <v-btn size="small" variant="tonal" prepend-icon="mdi-theme-light-dark" @click="toggleTheme">Theme</v-btn>
    </div>
    <p class="text-body-2 mb-3 text-medium-emphasis">
      {{ summary.length }} files in the summary (v-model). A file's diff is loaded when its tab opens, and forgotten
      when it closes; <code>templates/broken.gadx</code> fails to load.
    </p>
    <DiffBrowser v-model="summary" :load="load" height="560px">
      <template #actions="{ file }">
        <v-btn size="x-small" variant="text" icon="mdi-undo" title="Discard" @click="discard(file)" />
      </template>
    </DiffBrowser>
  </div>
</template>
