<script setup lang="ts">
// The image of the dialog of an image of the rich text: its preview, a place
// to drop, paste or choose one (sent by upload, when given; else kept as a
// data URL), and its editor — the media library's (window.rvqImageEditor),
// when the page has it.
import {computed, ref} from 'vue'
import {useVxMessages} from '@/lib/locale'

const props = defineProps<{
  modelValue?: string
  // the image sent: its URL (none: a data URL)
  upload?: ((file: File) => Promise<string>) | null
}>()
const emit = defineEmits<{'update:modelValue': [string]}>()

const vx = useVxMessages()
const over = ref(false)
const busy = ref(false)
const input = ref<HTMLInputElement | null>(null)
const editor = computed(() => (typeof window !== 'undefined' ? (window as any).rvqImageEditor : null))

const toDataURL = (f: Blob) =>
  new Promise<string>((resolve, reject) => {
    const r = new FileReader()
    r.onload = () => resolve(String(r.result))
    r.onerror = reject
    r.readAsDataURL(f)
  })

async function take(f: File | null | undefined) {
  if (!f || !/^image\//.test(f.type)) return
  busy.value = true
  try {
    emit('update:modelValue', props.upload ? await props.upload(f) : await toDataURL(f))
  } finally {
    busy.value = false
  }
}

function onDrop(ev: DragEvent) {
  over.value = false
  take(ev.dataTransfer?.files?.[0])
}

function onPaste(ev: ClipboardEvent) {
  const item = [...(ev.clipboardData?.items || [])].find((i) => i.kind === 'file' && i.type.startsWith('image/'))
  if (item) {
    ev.preventDefault()
    take(item.getAsFile())
  }
}

async function edit() {
  if (!props.modelValue || !editor.value) return
  const r = await editor.value.open({src: props.modelValue, name: 'image.png'})
  const save = r?.action === 'save' ? r.save : null
  if (save?.file) await take(save.file)
}

defineExpose({take})
</script>

<template>
  <div class="vx-tiptap-image-source" data-tiptap-image-source>
    <div
      class="vx-tiptap-image-drop"
      :class="{over}"
      tabindex="0"
      data-tiptap-image-drop
      @click="input?.click()"
      @dragover.prevent="over = true"
      @dragleave.prevent="over = false"
      @drop.prevent="onDrop"
      @paste="onPaste"
    >
      <img v-if="modelValue" :src="modelValue" alt="" />
      <div v-if="!modelValue || over" class="vx-tiptap-image-placeholder">
        <v-icon icon="mdi-image-plus-outline" class="mb-1" />
        <div class="text-body-2">{{ vx.imageDrop }}</div>
      </div>
      <v-progress-linear v-if="busy" indeterminate color="primary" absolute location="bottom" />
      <input ref="input" type="file" accept="image/*" style="display: none" @click.stop @change="take(($event.target as HTMLInputElement).files?.[0]); ($event.target as HTMLInputElement).value = ''" />
    </div>
    <div v-if="modelValue && editor" class="d-flex justify-center mt-1">
      <v-btn variant="text" size="small" prepend-icon="mdi-image-edit-outline" data-tiptap-image-edit @click="edit">{{ vx.imageEdit }}</v-btn>
    </div>
  </div>
</template>

<style scoped>
.vx-tiptap-image-source {
  width: 310px;
  margin: 0 auto 12px;
}
.vx-tiptap-image-drop {
  position: relative;
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px dashed rgba(var(--v-border-color), var(--v-border-opacity));
  border-radius: 6px;
  cursor: pointer;
  overflow: hidden;
  outline: none;
}
.vx-tiptap-image-drop.over,
.vx-tiptap-image-drop:focus {
  border-color: rgb(var(--v-theme-primary));
}
.vx-tiptap-image-drop img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}
.vx-tiptap-image-placeholder {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 8px;
  color: rgba(var(--v-theme-on-surface), 0.6);
  background: rgba(var(--v-theme-surface), 0.85);
}
</style>
