<script setup lang="ts">
// The button of the plugin math on the bar of vx-tiptap-editor: the dialog
// of a formula — its LaTeX, in the line or on a line of its own, its preview
// (KaTeX) —, put at the cursor, or the one clicked changed (or taken out).
import type {Editor} from '@tiptap/vue-3'
import {mdiSigma} from '@mdi/js'
import {ActionButton} from 'vuetify-pro-tiptap'
import {computed, onMounted, ref} from 'vue'
import katex from 'katex'
import {loadVendor} from '@rvq/mdkit'
import {useI18n} from '@/lib/locale'
import {mdkitWords} from '@/lib/locale/mdkit'

const props = defineProps<{ editor: Editor, disabled?: boolean }>()
const {t} = useI18n()
const words = mdkitWords(t).blocks
const show = ref(false)
const latex = ref('')
const block = ref(false)
const pos = ref<number | null>(null)

// the preview: KaTeX's HTML (its errors in red, not thrown)
const preview = computed(() => {
  try {
    return katex.renderToString(latex.value || '', {displayMode: block.value, throwOnError: false})
  } catch (e) {
    return ''
  }
})

function open(l = '', p: number | null = null, b = false) {
  latex.value = l
  pos.value = p
  block.value = b
  show.value = true
}

onMounted(() => {
  // KaTeX's styles (its fonts): the page's copy (mdkit-vendor)
  loadVendor('katex-css', '').catch(() => {})
  const storage = (props.editor.storage as any).mdkitMath
  if (storage) storage.open = open
})

function apply() {
  const c: any = props.editor.chain().focus()
  const l = latex.value.trim()
  if (!l) return
  if (pos.value !== null) {
    // the one clicked: of its kind (another kind: out, the new one in its place)
    const node = props.editor.state.doc.nodeAt(pos.value)
    const wasBlock = node?.type.name === 'blockMath'
    if (wasBlock === block.value) {
      ;(block.value ? c.updateBlockMath({latex: l, pos: pos.value}) : c.updateInlineMath({latex: l, pos: pos.value})).run()
    } else {
      ;(wasBlock ? c.deleteBlockMath({pos: pos.value}) : c.deleteInlineMath({pos: pos.value}))
        [block.value ? 'insertBlockMath' : 'insertInlineMath']({latex: l, pos: pos.value}).run()
    }
  } else {
    ;(block.value ? c.insertBlockMath({latex: l}) : c.insertInlineMath({latex: l})).run()
  }
  show.value = false
}

function remove() {
  if (pos.value === null) return
  const node = props.editor.state.doc.nodeAt(pos.value)
  const c: any = props.editor.chain().focus()
  ;(node?.type.name === 'blockMath' ? c.deleteBlockMath({pos: pos.value}) : c.deleteInlineMath({pos: pos.value})).run()
  show.value = false
}
</script>

<template>
  <ActionButton :rawIcon="`svg:${mdiSigma}`" :editor="editor" :tooltip="words.formula" :disabled="disabled"
                :action="() => open()" data-tiptap-plugin="math"/>
  <VDialog v-model="show" max-width="640" data-tiptap-math>
    <VCard>
      <VCardTitle class="d-flex align-center ga-2">
        <VIcon icon="mdi-sigma"/>
        {{ words.formulaTitle }}
        <VSpacer/>
        <VBtn size="small" variant="text" href="https://katex.org/docs/supported.html" target="_blank" rel="noopener"
              prepend-icon="mdi-book-open-variant">{{ words.guide }}</VBtn>
      </VCardTitle>
      <VCardText>
        <VTextarea v-model="latex" :label="words.latex" rows="3" auto-grow autofocus variant="outlined" hide-details
                   style="font-family: monospace" data-math-latex/>
        <VBtnToggle v-model="block" mandatory density="compact" variant="outlined" class="mt-3" data-math-kind>
          <VBtn :value="false" size="small">{{ words.inline }}</VBtn>
          <VBtn :value="true" size="small">{{ words.block }}</VBtn>
        </VBtnToggle>
        <div class="text-caption text-medium-emphasis mt-3">{{ words.preview }}</div>
        <div class="pa-3 border rounded" :class="block ? 'text-center' : ''" style="min-height: 3rem; overflow-x: auto"
             v-html="preview" data-math-preview/>
      </VCardText>
      <VCardActions>
        <VBtn v-if="pos !== null" color="error" variant="text" @click="remove" data-math-remove>{{ words.remove }}</VBtn>
        <VSpacer/>
        <VBtn variant="text" @click="show = false">{{ words.cancel }}</VBtn>
        <VBtn color="primary" variant="flat" :disabled="!latex.trim()" @click="apply" data-math-apply>
          {{ pos !== null ? words.update : words.insert }}
        </VBtn>
      </VCardActions>
    </VCard>
  </VDialog>
</template>
