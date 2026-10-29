<script setup lang="ts">
// The button of the plugin cards on the bar of vx-tiptap-editor: rvq mdkit's
// picker of a link card — the pages and the records found as one types (the
// page's search: its meta mdkit-link-cards), an address typed —; chosen, the
// card put at the cursor.
import type {Editor} from '@tiptap/vue-3'
import {mdiCardTextOutline} from '@mdi/js'
import {ActionButton} from 'vuetify-pro-tiptap'
import {ref} from 'vue'
import {LinkCardPicker} from '@rvq/mdkit'
import {useI18n} from '@/lib/locale'
import {mdkitWords} from '@/lib/locale/mdkit'

const props = defineProps<{ editor: Editor, disabled?: boolean }>()
const {t} = useI18n()
const words = mdkitWords(t)
const picker = ref<any>(null)

function open() {
  picker.value?.open((target: string) => (props.editor.chain().focus() as any).insertLinkCard(target).run())
}
</script>

<template>
  <ActionButton :rawIcon="`svg:${mdiCardTextOutline}`" :editor="editor" :tooltip="words.card.button" :disabled="disabled"
                :action="open" data-tiptap-plugin="cards"/>
  <LinkCardPicker ref="picker" :words="words.card"/>
</template>
