<script setup lang="ts">
// The button of the plugin diagrams on the bar of vx-tiptap-editor: the
// dialog of a diagram — its Mermaid source, its preview as it is written —,
// put at the cursor, or the one clicked changed (or taken out).
import type {Editor} from '@tiptap/vue-3'
import {mdiSitemapOutline} from '@mdi/js'
import {ActionButton} from 'vuetify-pro-tiptap'
import {onMounted, ref, watch} from 'vue'
import {drawDiagramInto} from '@rvq/mdkit'
import {useI18n} from '@/lib/locale'
import {mdkitWords} from '@/lib/locale/mdkit'

const props = defineProps<{ editor: Editor, disabled?: boolean }>()
const {t} = useI18n()
const words = mdkitWords(t).blocks
const show = ref(false)
const source = ref('')
const pos = ref<number | null>(null)
const preview = ref<HTMLElement | null>(null)
const invalid = ref(false)
let timer: any = null

const sample = 'graph TD\n  A[Start] --> B{Ok?}\n  B -->|Yes| C[Done]\n  B -->|No| A'

function open(s = '', p: number | null = null) {
  source.value = s || sample
  pos.value = p
  show.value = true
}

onMounted(() => {
  const storage = (props.editor.storage as any).mdkitDiagram
  if (storage) storage.open = open
})

// the preview, a moment after the last key
watch([source, show], () => {
  clearTimeout(timer)
  if (!show.value) return
  timer = setTimeout(async () => {
    if (preview.value) invalid.value = !(await drawDiagramInto(preview.value, source.value))
  }, 300)
})

function apply() {
  const s = source.value.trim()
  if (!s) return
  const c: any = props.editor.chain().focus()
  ;(pos.value !== null ? c.updateDiagram(pos.value, s) : c.insertDiagram(s)).run()
  show.value = false
}

function remove() {
  if (pos.value === null) return
  const node = props.editor.state.doc.nodeAt(pos.value)
  if (node) props.editor.chain().focus().deleteRange({from: pos.value, to: pos.value + node.nodeSize}).run()
  show.value = false
}
</script>

<template>
  <ActionButton :rawIcon="`svg:${mdiSitemapOutline}`" :editor="editor" :tooltip="words.diagram" :disabled="disabled"
                :action="() => open()" data-tiptap-plugin="diagrams"/>
  <VDialog v-model="show" max-width="900" data-tiptap-diagram>
    <VCard>
      <VCardTitle class="d-flex align-center ga-2">
        <VIcon icon="mdi-sitemap-outline"/>
        {{ words.diagramTitle }}
        <VSpacer/>
        <VBtn size="small" variant="text" href="https://mermaid.js.org/intro/syntax-reference.html" target="_blank"
              rel="noopener" prepend-icon="mdi-book-open-variant">{{ words.guide }}</VBtn>
      </VCardTitle>
      <VCardText>
        <div class="d-flex ga-4 flex-wrap">
          <VTextarea v-model="source" :label="words.source" rows="10" auto-grow autofocus variant="outlined" hide-details
                     style="font-family: monospace; flex: 1 1 18rem" data-diagram-source/>
          <div style="flex: 1 1 18rem; min-width: 0">
            <div class="text-caption text-medium-emphasis">{{ words.preview }}</div>
            <div ref="preview" class="pa-2 border rounded text-center" style="min-height: 10rem; overflow: auto"
                 data-diagram-preview/>
            <div v-if="invalid" class="text-caption text-error mt-1" data-diagram-invalid>{{ words.invalid }}</div>
          </div>
        </div>
      </VCardText>
      <VCardActions>
        <VBtn v-if="pos !== null" color="error" variant="text" @click="remove" data-diagram-remove>{{ words.remove }}</VBtn>
        <VSpacer/>
        <VBtn variant="text" @click="show = false">{{ words.cancel }}</VBtn>
        <VBtn color="primary" variant="flat" :disabled="!source.trim()" @click="apply" data-diagram-apply>
          {{ pos !== null ? words.update : words.insert }}
        </VBtn>
      </VCardActions>
    </VCard>
  </VDialog>
</template>
