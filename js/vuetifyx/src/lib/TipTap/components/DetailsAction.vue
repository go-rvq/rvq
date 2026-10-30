<script setup lang="ts">
// The button of the plugin details on the bar of vx-tiptap-editor: the
// blocks chosen made a section that folds, or the section under the cursor
// its blocks again.
import type {Editor} from '@tiptap/vue-3'
import {mdiArrowCollapseVertical} from '@mdi/js'
import {ActionButton} from 'vuetify-pro-tiptap'
import {useI18n} from '@/lib/locale'
import {mdkitWords} from '@/lib/locale/mdkit'

const props = defineProps<{ editor: Editor, disabled?: boolean }>()
const {t} = useI18n()
const words = mdkitWords(t).blocks

function toggle() {
  const c: any = props.editor.chain().focus()
  ;(props.editor.isActive('details') ? c.unsetDetails() : c.setDetails()).run()
}
</script>

<template>
  <ActionButton :rawIcon="`svg:${mdiArrowCollapseVertical}`" :editor="editor" :tooltip="words.details" :disabled="disabled"
                :action="toggle" :is-active="() => editor.isActive('details')" data-tiptap-plugin="details"/>
</template>
