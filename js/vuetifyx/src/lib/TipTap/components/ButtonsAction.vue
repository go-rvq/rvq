<script setup lang="ts">
// The button of the plugin buttons on the bar of vx-tiptap-editor: rvq
// mdkit's buttons dialog, put at the cursor, or the buttons clicked changed
// (or taken out); its words the namespace vx's (Mdkit…).
import type {Editor} from '@tiptap/vue-3'
import {mdiGestureTapButton} from '@mdi/js'
import {ActionButton} from 'vuetify-pro-tiptap'
import {onBeforeUnmount, onMounted, ref} from 'vue'
import {ButtonsDialog} from '@rvq/mdkit'
import {useI18n} from '@/lib/locale'
import {mdkitWords} from '@/lib/locale/mdkit'

const props = defineProps<{ editor: Editor, disabled?: boolean }>()
const {t} = useI18n()
const words = mdkitWords(t)
const dialog = ref<any>(null)

function open(value: any = null, pos: number | null = null) {
  dialog.value?.open(value, (v: any) => {
    const c: any = props.editor.chain().focus()
    if (pos !== null) c.updateMdkitButtons(pos, v).run()
    else if (v) c.insertMdkitButtons(v).run()
  })
}

onMounted(() => {
  const s = (props.editor.storage as any).mdkitButtons
  if (s) s.open = open
})
onBeforeUnmount(() => {
  const s = (props.editor.storage as any).mdkitButtons
  if (s && s.open === open) s.open = null
})
</script>

<template>
  <ActionButton :rawIcon="`svg:${mdiGestureTapButton}`" :editor="editor" :tooltip="words.widgets.buttons"
                :disabled="disabled" :action="() => open()" data-tiptap-plugin="buttons"/>
  <ButtonsDialog ref="dialog" :words="words.widgets" :link-words="words.link"/>
</template>
