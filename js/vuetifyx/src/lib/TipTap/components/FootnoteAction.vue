<script setup lang="ts">
// The button of the plugin footnotes on the bar of vx-tiptap-editor: rvq
// mdkit's footnotes builder on the text (tiptapFootnoteAdapter), its words
// the namespace vx's (Mdkit…). The clicks on a reference or a note open it
// too (editor.storage.mdkitFootnotes.open).
import type {Editor} from '@tiptap/vue-3'
import {mdiFormatSuperscript} from '@mdi/js'
import {ActionButton} from 'vuetify-pro-tiptap'
import {markRaw, onBeforeUnmount, onMounted, ref} from 'vue'
import {FootnoteBuilder, RichEditor, configure, textStyle, tiptapFootnoteAdapter} from '@rvq/mdkit'
import {useI18n} from '@/lib/locale'
import {mdkitWords} from '@/lib/locale/mdkit'

const props = defineProps<{ editor: Editor, disabled?: boolean }>()
const {t} = useI18n()
const words = mdkitWords(t)
// mdkit's components inside (the editor of a note) speak the same words
configure({words})
textStyle()
const builder = ref<any>(null)
// the editor of the text of a note: mdkit's rich text editor
const noteEditor = markRaw(RichEditor)

function open(label?: string, pos: number | null = null) {
  builder.value?.open(tiptapFootnoteAdapter(props.editor as any, pos), {label, ref: pos != null})
}

onMounted(() => {
  (props.editor.storage as any).mdkitFootnotes.open = open
})
onBeforeUnmount(() => {
  const s = (props.editor.storage as any).mdkitFootnotes
  if (s && s.open === open) s.open = null
})
</script>

<template>
  <ActionButton :rawIcon="`svg:${mdiFormatSuperscript}`" :editor="editor" :tooltip="words.rich.footnote"
                :disabled="disabled" :action="() => open()" data-tiptap-plugin="footnotes"/>
  <FootnoteBuilder ref="builder" :words="words.footnote" :bar-words="words.rich" :note-editor="noteEditor"/>
</template>
