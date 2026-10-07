<script setup lang="ts">
import { computed, markRaw, provide, reactive, shallowReactive, shallowRef, useSlots, watch } from 'vue'

import { useTheme } from 'vuetify'
import { DockviewVue } from 'dockview-vue'
import type { DockviewApi, DockviewReadyEvent } from 'dockview-vue'
import { themeDark, themeLight } from 'dockview-core'
import 'dockview-vue/dist/styles/dockview.css'
import DiffBrowserFiles from './DiffBrowserFiles.vue'
import DiffBrowserFile from './DiffBrowserFile.vue'
import { diffBrowserKey, renamedTo, type Content, type DiffFile, type SaveState } from './diffBrowserContext'

// <vx-diff-browser>: the changes of a set of files, browsed. On the left the
// tree of the files changed — each in its folders, its path kept —; clicking
// one opens it on the right, in a tab of its own: the old and the current side
// by side, the changed lines and, within them, the changed content marked (an
// IntelliJ compare), the two sides scrolling together. The panels are
// dockview's: resized, the tabs moved, split.
//
// The comparison is the server's (vuetifyx.NewDiffFile): this only shows it.
// The `actions` slot ({ file }) is drawn by each file of the tree (a discard).
// The summary (the files of the tree) is its v-model, or files.

const props = withDefaults(defineProps<{
  // the summary, by v-model — or by files
  modelValue?: DiffFile[]
  files?: DiffFile[]
  height?: string
  filesTitle?: string
  oldLabel?: string
  newLabel?: string
  binaryText?: string
  emptyText?: string
  renamedText?: string
  unchangedText?: string
  // load asks for the diff of a file opened, given the reactive object of its
  // tab (content: { file, value, error }): a plaid() call whose response
  // sets content.value — `(content) => plaid().scope({content})….go()`. With
  // none, the files have their contents already.
  load?: (content: Content) => unknown
  loadErrorText?: string
  loadingText?: string
  // save writes the file at path with what its current side has been edited
  // to: save(path, state) — state reactive, its value the text; the server's
  // answer sets state.saved = true, or state.error. With it, the current
  // sides are editable (a file readOnly aside).
  save?: (path: string, state: SaveState) => unknown
  undoText?: string
  redoText?: string
  saveText?: string
  savingText?: string
  savedText?: string
  unsavedText?: string
  revertText?: string
  prevText?: string
  nextText?: string
}>(), {
  height: '70vh',
  filesTitle: 'Files',
  oldLabel: 'Old',
  newLabel: 'Current',
  binaryText: 'A binary file: it is not compared.',
  emptyText: 'No changes.',
  renamedText: 'Renamed to',
  unchangedText: 'Its content did not change.',
  loadErrorText: 'The diff could not be loaded.',
  loadingText: 'Loading…',
  undoText: 'Undo',
  redoText: 'Redo',
  saveText: 'Save',
  savingText: 'Saving…',
  savedText: 'Saved',
  unsavedText: 'Not saved',
  revertText: 'Revert: the old part into the current',
  prevText: 'Previous change (Shift+F7)',
  nextText: 'Next change (F7)',
})

const slots = useSlots()
// the summary: v-model's, else files
const summary = computed<DiffFile[]>(() => props.modelValue ?? props.files ?? [])
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
    title: tabTitle(file),
    params: { path: file.path },
    position: other ? { referencePanel: other.id, direction: 'within' } : { referencePanel: 'files', direction: 'right' },
  })
  if (!other) dv.getPanel('files')?.api.setSize({ width: 280 })
}

// tabTitle is a file's tab's title: its name (a renamed one's both), ● when
// it has changes not saved
function tabTitle(file: DiffFile, dirty = false) {
  const name = file.from
    ? (file.from.split('/').pop() || file.from) + ' → ' + renamedTo(file)
    : file.path.split('/').pop() || file.path
  return (dirty ? '● ' : '') + name
}

function markDirty(path: string, dirty: boolean) {
  const f = summary.value.find((x) => x.path === path)
  if (f) api.value?.getPanel('file:' + path)?.api.setTitle(tabTitle(f, dirty))
}

// the contents of the files opened, by path: only those, and only while
// their tabs are open — the browser keeps no more than it shows
const contents = shallowReactive(new Map<string, Content>())

function request(file: DiffFile) {
  if (!props.load) {
    contents.set(file.path, reactive({ file, value: file }))
    return
  }
  const had = contents.get(file.path)
  if (had && !had.error) return
  const content = reactive<Content>({ file: { ...file }, value: undefined, error: undefined })
  contents.set(file.path, content)
  try {
    // a promise it returns that fails is its error too
    const p = props.load(content) as any
    if (p && typeof p.catch === 'function') {
      p.catch((e: any) => {
        if (!content.value) content.error = String(e?.message || e || props.loadErrorText)
      })
    }
  } catch (e: any) {
    content.error = String(e?.message || e || props.loadErrorText)
  }
}

provide(diffBrowserKey, {
  files: summary,
  content: (path: string) => contents.get(path),
  request,
  active,
  open,
  labels: computed(() => ({
    old: props.oldLabel,
    new: props.newLabel,
    binary: props.binaryText,
    renamed: props.renamedText,
    unchanged: props.unchangedText,
    loading: props.loadingText,
    undo: props.undoText,
    redo: props.redoText,
    save: props.saveText,
    saving: props.savingText,
    saved: props.savedText,
    unsaved: props.unsavedText,
    revert: props.revertText,
    prev: props.prevText,
    next: props.nextText,
  })),
  dark: computed(() => theme.current.value.dark),
  editable: (file: DiffFile) => !!props.save && !file.readOnly && !file.binary,
  save: (path: string, state: SaveState) => {
    try {
      const p = props.save?.(path, state) as any
      if (p && typeof p.catch === 'function') p.catch((e: any) => (state.error = String(e?.message || e)))
    } catch (e: any) {
      state.error = String(e?.message || e)
    }
  },
  markDirty,
  actions: () => slots.actions,
})

function onReady(event: DockviewReadyEvent) {
  api.value = event.api
  event.api.addPanel({ id: 'files', component: 'files', title: props.filesTitle })
  // a tab closed forgets its file's content
  event.api.onDidRemovePanel((p) => {
    if (p.id.startsWith('file:')) contents.delete(p.id.slice(5))
  })
  event.api.onDidActivePanelChange((e) => {
    const id = e.panel?.id
    if (id && id.startsWith('file:')) active.path = id.slice(5)
  })
  if (summary.value.length) open(summary.value[0])
}

// fileKey is what a file of the summary is: its path, status and origin.
// The summary may be given anew on each render of the page (a literal of its
// template — the admin's renders on every request) and say the same: only a
// change of what it says counts. Watching the array itself, each answer of a
// load re-rendered the page, which gave a new array, which asked again: a
// loop of requests.
const fileKey = (f: DiffFile) => [f.path, f.status || '', f.from || ''].join('\u0001')
const summaryKey = computed(() => summary.value.map(fileKey).join('\u0002'))

watch(summaryKey, () => {
  const dv = api.value
  const byPath = new Map(summary.value.map((f) => [f.path, f]))
  // a file that is gone (discarded) closes its tab
  if (dv) {
    for (const p of [...dv.panels]) {
      if (p.id.startsWith('file:') && !byPath.has(p.id.slice(5))) p.api.close()
    }
  }
  // a file opened that changed (its status, its origin) is asked for again
  for (const [path, content] of [...contents]) {
    const f = byPath.get(path)
    if (f && fileKey(f) !== fileKey(content.file)) {
      contents.delete(path)
      request(f)
    }
  }
})

defineExpose({ open })
</script>

<template>
  <div class="vx-diff-browser" :style="{ height }">
    <div v-if="!summary.length" class="pa-4 text-medium-emphasis">{{ emptyText }}</div>
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
.vx-diff-toolbar {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 2px;
  padding: 4px 8px;
  border-bottom: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}
.vx-diff-sides-labels {
  display: grid;
  grid-template-columns: 1fr 1fr;
  border-bottom: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}
.vx-diff-sides-labels > .vx-diff-side-label {
  border-bottom: 0;
}
.vx-diff-sides-labels > .vx-diff-side-label + .vx-diff-side-label {
  border-left: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}
.vx-diff-merge {
  flex: 1 1 auto;
  min-height: 0;
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
