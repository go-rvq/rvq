<script setup lang="ts">
import { computed, markRaw, provide, reactive, shallowRef, useSlots, watch } from 'vue'
import { useTheme } from 'vuetify'
import { DockviewVue } from 'dockview-vue'
import type { DockviewApi, DockviewReadyEvent } from 'dockview-vue'
import { themeDark, themeLight } from 'dockview-core'
import 'dockview-vue/dist/styles/dockview.css'
import DiffBrowserFiles from './DiffBrowserFiles.vue'
import DiffBrowserFile from './DiffBrowserFile.vue'
import { diffBrowserKey, renamedTo, type DiffFile } from './diffBrowserContext'

// <vx-diff-browser>: the changes of a set of files, browsed. On the left the
// tree of the files changed — each in its folders, its path kept —; clicking
// one opens it on the right, in a tab of its own: the old and the current side
// by side, the changed lines and, within them, the changed content marked (an
// IntelliJ compare), the two sides scrolling together. The panels are
// dockview's: resized, the tabs moved, split.
//
// The comparison is the server's (vuetifyx.NewDiffFile): this only shows it.
// The `actions` slot ({ file }) is drawn by each file of the tree (a discard).

const props = withDefaults(defineProps<{
  files: DiffFile[]
  height?: string
  filesTitle?: string
  oldLabel?: string
  newLabel?: string
  binaryText?: string
  emptyText?: string
  renamedText?: string
  unchangedText?: string
}>(), {
  height: '70vh',
  filesTitle: 'Files',
  oldLabel: 'Old',
  newLabel: 'Current',
  binaryText: 'A binary file: it is not compared.',
  emptyText: 'No changes.',
  renamedText: 'Renamed to',
  unchangedText: 'Its content did not change.',
})

const slots = useSlots()
const theme = useTheme()
// dockview's theme follows Vuetify's, as the IDE's does
const dockTheme = computed(() => (theme.current.value.dark ? themeDark : themeLight))

// the panels, by the name addPanel gives (dockview-vue resolves them here)
const panels: Record<string, any> = { files: markRaw(DiffBrowserFiles), diff: markRaw(DiffBrowserFile) }

const api = shallowRef<DockviewApi>()
const active = reactive({ path: '' })

// open shows the file in its tab: the one it has, or a new one beside the
// others (the first, right of the tree).
function open(file: DiffFile) {
  const dv = api.value
  if (!dv) return
  const id = 'file:' + file.path
  const panel = dv.getPanel(id)
  if (panel) {
    panel.api.setActive()
    return
  }
  const other = dv.panels.find((p) => p.id.startsWith('file:'))
  dv.addPanel({
    id,
    component: 'diff',
    // a renamed file's tab says both names
    title: file.from
      ? (file.from.split('/').pop() || file.from) + ' → ' + renamedTo(file)
      : file.path.split('/').pop() || file.path,
    params: { path: file.path },
    position: other ? { referencePanel: other.id, direction: 'within' } : { referencePanel: 'files', direction: 'right' },
  })
  if (!other) dv.getPanel('files')?.api.setSize({ width: 280 })
}

provide(diffBrowserKey, {
  files: computed(() => props.files),
  active,
  open,
  labels: computed(() => ({
    old: props.oldLabel,
    new: props.newLabel,
    binary: props.binaryText,
    renamed: props.renamedText,
    unchanged: props.unchangedText,
  })),
  actions: () => slots.actions,
})

function onReady(event: DockviewReadyEvent) {
  api.value = event.api
  event.api.addPanel({ id: 'files', component: 'files', title: props.filesTitle })
  event.api.onDidActivePanelChange((e) => {
    const id = e.panel?.id
    if (id && id.startsWith('file:')) active.path = id.slice(5)
  })
  if (props.files.length) open(props.files[0])
}

// a file that is gone (discarded) closes its tab
watch(
  () => props.files.map((f) => f.path),
  (paths) => {
    const dv = api.value
    if (!dv) return
    for (const p of [...dv.panels]) {
      if (p.id.startsWith('file:') && !paths.includes(p.id.slice(5))) p.api.close()
    }
  },
)

defineExpose({ open })
</script>

<template>
  <div class="vx-diff-browser" :style="{ height }">
    <div v-if="!files.length" class="pa-4 text-medium-emphasis">{{ emptyText }}</div>
    <DockviewVue v-else :theme="dockTheme" style="height: 100%" :components="panels" @ready="onReady" />
  </div>
</template>

<style>
.vx-diff-browser {
  position: relative;
  border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  border-radius: 6px;
  overflow: hidden;
}
.vx-diff-tree {
  height: 100%;
  overflow: auto;
  padding: 4px 0;
  font-size: 0.8125rem;
  color: rgb(var(--v-theme-on-surface));
}
.vx-diff-tree-row {
  display: flex;
  align-items: center;
  height: 26px;
  padding-right: 6px;
  cursor: pointer;
  white-space: nowrap;
  user-select: none;
}
.vx-diff-tree-row:hover {
  background: rgba(var(--v-theme-on-surface), 0.05);
}
.vx-diff-tree-active {
  background: rgba(var(--v-theme-primary), 0.12);
}
.vx-diff-tree-name {
  overflow: hidden;
  text-overflow: ellipsis;
}
.vx-diff-tree-renamed {
  margin-left: 6px;
  color: rgb(var(--v-theme-info));
  overflow: hidden;
  text-overflow: ellipsis;
}
.vx-diff-tree-actions {
  margin-left: auto;
  display: inline-flex;
}
.vx-diff-file {
  display: flex;
  flex-direction: column;
  height: 100%;
  color: rgb(var(--v-theme-on-surface));
}
.vx-diff-file-head {
  display: flex;
  align-items: center;
  padding: 6px 10px;
  border-bottom: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  font-size: 0.8125rem;
}
.vx-diff-sides {
  display: grid;
  grid-template-columns: 1fr 1fr;
  flex: 1 1 auto;
  min-height: 0;
}
.vx-diff-side {
  display: flex;
  flex-direction: column;
  min-width: 0;
  min-height: 0;
}
.vx-diff-side + .vx-diff-side {
  border-left: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}
.vx-diff-side-label {
  padding: 4px 10px;
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: rgba(var(--v-theme-on-surface), 0.6);
  border-bottom: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}
.vx-diff-side-body {
  flex: 1 1 auto;
  overflow: auto;
  min-height: 0;
}
.vx-diff-side-body .vx-code {
  border-radius: 0;
  min-height: 100%;
  width: max-content;
  min-width: 100%;
}
</style>
