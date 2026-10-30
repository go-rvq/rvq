<script setup lang="ts">
// The button of the plugin tabs on the bar of vx-tiptap-editor: rvq mdkit's
// tabs dialog — each tab's content written in a vx-tiptap-editor of the
// same plugins (but tabs) —, put at the cursor, or the tabs clicked changed
// (or taken out); its words the namespace vx's (Mdkit…).
import type {Editor} from '@tiptap/vue-3'
import {mdiTab} from '@mdi/js'
import {ActionButton} from 'vuetify-pro-tiptap'
import {computed, onBeforeUnmount, onMounted, ref, resolveComponent} from 'vue'
import {TabsDialog} from '@rvq/mdkit'
import {useI18n} from '@/lib/locale'
import {mdkitWords} from '@/lib/locale/mdkit'

const props = defineProps<{ editor: Editor, disabled?: boolean, plugins?: string[] }>()
const {t} = useI18n()
const words = mdkitWords(t)
const dialog = ref<any>(null)
// the editor of a tab's content: this one (registered as vx-tiptap-editor)
const editorComp = resolveComponent('vx-tiptap-editor')
const editorProps = computed(() => ({plugins: props.plugins || [], output: 'html'}))

function open(value: any = null, pos: number | null = null) {
  dialog.value?.open(value, (v: any) => {
    const c: any = props.editor.chain().focus()
    if (pos !== null) c.updateMdkitTabs(pos, v).run()
    else if (v) c.insertMdkitTabs(v).run()
  })
}

onMounted(() => {
  const s = (props.editor.storage as any).mdkitTabs
  if (s) s.open = open
})
onBeforeUnmount(() => {
  const s = (props.editor.storage as any).mdkitTabs
  if (s && s.open === open) s.open = null
})
</script>

<template>
  <ActionButton :rawIcon="`svg:${mdiTab}`" :editor="editor" :tooltip="words.widgets.tabs"
                :disabled="disabled" :action="() => open()" data-tiptap-plugin="tabs"/>
  <TabsDialog ref="dialog" :words="words.widgets" :editor="editorComp" :editor-props="editorProps"/>
</template>
