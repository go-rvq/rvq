<script setup lang="ts">
import type {Editor} from '@tiptap/vue-3'
import {mdiCodeTags} from '@mdi/js'
import {ActionButton} from 'vuetify-pro-tiptap'
import {ref} from 'vue'
import CodeMirror from '@/lib/CodeMirror'
import Dialog from '@/lib/Dialog.vue'

const props = defineProps<{
  editor: Editor
  disabled?: boolean
}>()

const open = ref(false)
const html = ref('')

// Open the dialog with the editor's current HTML source.
function action() {
  html.value = props.editor.getHTML()
  open.value = true
}

// Apply replaces the editor content with the edited HTML source.
function apply() {
  props.editor.commands.setContent(html.value, {emitUpdate: true})
  open.value = false
}

function cancel() {
  open.value = false
}
</script>

<template>
  <ActionButton
    :rawIcon="`svg:${mdiCodeTags}`"
    :editor="editor"
    tooltip="Código (HTML)"
    :disabled="disabled"
    :action="action"
  />

  <Dialog v-model="open" expandable closable title="Código (HTML)" width="60%">
    <template #body>
      <CodeMirror v-model="html" language="html" min-height="24rem" max-height="70vh" />
    </template>
    <template #bottom>
      <div class="d-flex justify-end ga-2 pa-3">
        <v-btn variant="text" @click="cancel">Cancelar</v-btn>
        <v-btn color="primary" variant="flat" @click="apply">Aplicar</v-btn>
      </div>
    </template>
  </Dialog>
</template>
