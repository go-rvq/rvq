<template>
  <div class="vx-file-upload" data-file-upload>
    <div v-if="label" class="text-caption text-medium-emphasis mb-1">{{ label }}</div>
    <div
      class="vx-file-upload-drop d-flex flex-column align-center justify-center pa-4 text-center"
      :class="{over}"
      role="button"
      tabindex="0"
      data-file-upload-drop
      @click="choose"
      @keydown.enter.prevent="choose"
      @dragover.prevent="over = true"
      @dragleave="over = false"
      @drop.prevent="dropped"
    >
      <v-icon icon="mdi-tray-arrow-up" size="large" class="mb-1" />
      <span class="text-body-2">{{ placeholder || words.uploadDrop }}</span>
      <input ref="input" type="file" class="d-none" :multiple="multiple" :accept="accept" @change="picked" />
    </div>
    <TransitionGroup name="vx-file-upload-item" tag="div">
    <div v-for="it in items" :key="it.key" class="vx-file-upload-item d-flex align-center ga-2 mt-2" :data-file-upload-item="it.state">
      <v-icon :icon="iconOf(it)" :color="colorOf(it)" size="small" />
      <div class="flex-grow-1" style="min-width: 0">
        <div class="text-body-2 text-truncate">{{ it.file.name }} <span class="text-caption text-medium-emphasis">· {{ sizeOf(it.file.size) }}</span></div>
        <div v-if="it.state === 'exists'" class="text-caption" :class="blockExisting ? 'text-error' : 'text-medium-emphasis'" data-file-upload-exists>
          {{ blockExisting ? words.uploadExists : words.uploadReused }}
          <a :href="it.existing?.href" @click.prevent="emit('open', {...it.existing, portal: portalOf(it)})" data-file-upload-existing>{{ it.existing?.title }}</a>
          <template v-if="it.existing?.createdAt"> ({{ it.existing.createdAt }})</template>
          <!-- where what the click opens goes (the detail of the file, in a dialog) -->
          <go-plaid-portal :visible="true" :portal-name="portalOf(it)" />
        </div>
        <div v-else-if="it.state === 'error'" class="text-caption text-error">{{ it.error }}</div>
        <v-progress-linear v-else-if="it.state !== 'done'" :model-value="it.progress * 100" height="3" color="primary" class="mt-1" />
      </div>
      <span class="text-caption text-medium-emphasis">{{ stateText(it) }}</span>
      <v-btn icon="mdi-close" size="x-small" variant="text" :title="words.clear" @click="remove(it)" />
    </div>
    </TransitionGroup>
    <div v-if="errorMessages && errorMessages.length" class="text-caption text-error mt-1">{{ errorMessages.join(' ') }}</div>
  </div>
</template>

<script lang="ts">
// the files of the components kept by a key (persist), and the numbers of
// their items: of the module, outlive a component
import {ref as keptRef, type Ref as KeptRef} from 'vue'
const store = new Map<string, KeptRef<any[]>>()
export function kept(key: string): KeptRef<any[]> {
  let r = store.get(key)
  if (!r) store.set(key, (r = keptRef<any[]>([])))
  return r
}
let seq = 0
// how long the line of a file sent stays
export const DONE_MS = 3000
</script>

<script lang="ts" setup>
// vx-file-upload: files dropped (or chosen) sent in parts (admin/uploadparts),
// their sha256 first, in the browser: a content the server has is not sent —
// with blockExisting (the model of the files itself), it is refused, saying
// which file it is (its detail, by its link: a click opens it in a dialog,
// the "open" event; the middle button, in a new tab); else only its hash
// goes, the content the server has used. Each file done, "uploaded" (its token,
// or {hash, name}); all done, "done". One file (not multiple): its value
// (v-model) is the token as JSON, "" when there is none. The message of a file
// that exists has its portal: "open" says it (portal), the target of what the
// click opens.
import {ref, type Ref} from 'vue'
import {useVxMessages} from '@/lib/locale'
import {call, exists, hashFile, send, type Token} from './parts'

const props = defineProps<{
  url: string
  multiple?: boolean
  blockExisting?: boolean
  label?: string
  placeholder?: string
  accept?: string
  errorMessages?: string[]
  // createEvent: the event of the model (at url) that makes a record of a
  // file sent ({token}), called by the component itself — not by a listener
  // of "uploaded", gone when the page around it is drawn again
  createEvent?: string
  // persist: the key the files and their messages are kept by — a component
  // drawn again (the listing reloaded) shows them again
  persist?: string
}>()
const model = defineModel<string>({default: ''})
const emit = defineEmits<{
  uploaded: [value: Token | {hash: string; name: string; mtime?: number}]
  done: []
  open: [existing: any]
}>()

const words = useVxMessages()
const over = ref(false)
const input = ref<HTMLInputElement | null>(null)

interface Item {
  key: number
  file: File
  state: 'hashing' | 'checking' | 'sending' | 'done' | 'exists' | 'error'
  progress: number
  existing?: any
  error?: string
}
const items: Ref<Item[]> = props.persist ? kept(props.persist) : ref<Item[]>([])
// one per component: the portals of its messages
const uid = Math.random().toString(36).slice(2, 10)
const portalOf = (it: Item) => `vx-file-upload-${uid}-${it.key}`

function choose() {
  input.value?.click()
}
function picked(e: Event) {
  const el = e.target as HTMLInputElement
  add([...(el.files || [])])
  el.value = ''
}
function dropped(e: DragEvent) {
  over.value = false
  add([...(e.dataTransfer?.files || [])])
}
function remove(it: Item) {
  items.value = items.value.filter((x) => x !== it)
  if (!props.multiple) model.value = ''
}

async function add(files: File[]) {
  if (!files.length) return
  if (!props.multiple) {
    files = files.slice(0, 1)
    items.value = []
    model.value = ''
  }
  const added = files.map((file) => ({key: ++seq, file, state: 'hashing', progress: 0}) as Item)
  items.value.push(...added)
  for (const it of added) await run(items.value.find((x) => x.key === it.key)!)
  emit('done')
}

async function run(it: Item) {
  try {
    const hash = await hashFile(it.file, (p) => (it.progress = p))
    it.state = 'checking'
    const existing = await exists(props.url, hash)
    if (existing) {
      it.existing = existing
      it.state = 'exists'
      if (props.blockExisting) return
      const ref = {hash, name: it.file.name, mtime: it.file.lastModified || undefined}
      if (props.createEvent) await call(props.url, props.createEvent, {token: JSON.stringify(ref)})
      if (!props.multiple) model.value = JSON.stringify(ref)
      emit('uploaded', ref)
      return
    }
    it.state = 'sending'
    it.progress = 0
    const token = await send(props.url, it.file, hash, (p) => (it.progress = p))
    if (props.createEvent) await call(props.url, props.createEvent, {token: JSON.stringify(token)})
    it.state = 'done'
    if (!props.multiple) model.value = JSON.stringify(token)
    emit('uploaded', token)
    // a file sent: its line goes after a moment (those refused, failed, stay)
    if (props.multiple) setTimeout(() => (items.value = items.value.filter((x) => x.key !== it.key)), DONE_MS)
  } catch (e: any) {
    it.state = 'error'
    it.error = String(e?.message || e)
  }
}

function iconOf(it: Item) {
  switch (it.state) {
    case 'done':
      return 'mdi-check-circle'
    case 'exists':
      return props.blockExisting ? 'mdi-alert-circle' : 'mdi-content-duplicate'
    case 'error':
      return 'mdi-alert-circle'
  }
  return 'mdi-file-upload-outline'
}
function colorOf(it: Item) {
  return it.state === 'done' ? 'success' : it.state === 'error' || (it.state === 'exists' && props.blockExisting) ? 'error' : undefined
}
function stateText(it: Item) {
  switch (it.state) {
    case 'hashing':
      return words.value.uploadHashing
    case 'checking':
      return words.value.uploadChecking
    case 'sending':
      return `${Math.round(it.progress * 100)}%`
  }
  return ''
}
function sizeOf(n: number) {
  if (n < 1024) return `${n} B`
  if (n < 1 << 20) return `${(n / 1024).toFixed(0)} KB`
  if (n < 1 << 30) return `${(n / (1 << 20)).toFixed(1)} MB`
  return `${(n / (1 << 30)).toFixed(2)} GB`
}

defineExpose({items, add})
</script>

<style scoped>
.vx-file-upload-drop {
  border: 2px dashed rgba(var(--v-border-color), 0.5);
  border-radius: 8px;
  min-height: 88px;
  cursor: pointer;
}
/* a line that goes: fades out, its height closed */
.vx-file-upload-item-leave-active,
.vx-file-upload-item-enter-active {
  transition: opacity 0.4s ease, max-height 0.4s ease, margin 0.4s ease;
  max-height: 4rem;
  overflow: hidden;
}
.vx-file-upload-item-leave-to,
.vx-file-upload-item-enter-from {
  opacity: 0;
  max-height: 0;
  margin-top: 0 !important;
}
.vx-file-upload-drop.over {
  border-color: rgb(var(--v-theme-primary));
  background: rgba(var(--v-theme-primary), 0.06);
}
</style>
