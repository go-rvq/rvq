<script setup lang="ts">
// The button of the plugin tables that places the content of the cells under
// the cursor (or chosen): across — left, center, right — and up and down —
// top, middle, bottom —; the one chosen again, the default. The cell's
// attributes align and valign (the plugin's), the same the table builder
// writes.
import type {Editor} from '@tiptap/vue-3'
import {mdiTableCog} from '@mdi/js'
import {ActionButton} from 'vuetify-pro-tiptap'
import {computed} from 'vue'
import {useI18n} from '@/lib/locale'
import {mdkitWords} from '@/lib/locale/mdkit'

const props = defineProps<{ editor: Editor, disabled?: boolean }>()
const {t} = useI18n()
const words = mdkitWords(t).table

// inTable says the cursor is in a table
const inTable = computed(() => props.editor.isActive('tableCell') || props.editor.isActive('tableHeader'))

// attr is the value of name of the cell under the cursor
function attr(name: string) {
  const a = props.editor.isActive('tableHeader') ? props.editor.getAttributes('tableHeader') : props.editor.getAttributes('tableCell')
  return a[name] || null
}

function set(name: string, value: string | null) {
  props.editor.chain().focus().setCellAttribute(name, value || null).run()
}

const across = [
  {value: 'left', icon: 'mdi-format-align-left', title: () => words.alignLeft},
  {value: 'center', icon: 'mdi-format-align-center', title: () => words.alignCenter},
  {value: 'right', icon: 'mdi-format-align-right', title: () => words.alignRight}
]
const upDown = [
  {value: 'top', icon: 'mdi-format-vertical-align-top', title: () => words.valignTop},
  {value: 'middle', icon: 'mdi-format-vertical-align-center', title: () => words.valignMiddle},
  {value: 'bottom', icon: 'mdi-format-vertical-align-bottom', title: () => words.valignBottom}
]
</script>

<template>
  <ActionButton :rawIcon="`svg:${mdiTableCog}`" :editor="editor" :tooltip="words.cellPlace"
                :disabled="disabled || !inTable" data-tiptap-plugin="cell-place">
    <VMenu activator="parent" :close-on-content-click="false" location="bottom">
      <VCard class="pa-1" data-tiptap-cell-place>
        <div class="d-flex" data-cell-align>
          <VBtn v-for="o in across" :key="o.value" :icon="o.icon" size="small" variant="text" :active="attr('align') === o.value" :data-value="o.value" @click="set('align', attr('align') === o.value ? null : o.value)"><VIcon :icon="o.icon"></VIcon><VTooltip activator="parent" location="top" :text="o.title()"/></VBtn>
        </div>
        <div class="d-flex" data-cell-valign>
          <VBtn v-for="o in upDown" :key="o.value" :icon="o.icon" size="small" variant="text" :active="attr('valign') === o.value" :data-value="o.value" @click="set('valign', attr('valign') === o.value ? null : o.value)"><VIcon :icon="o.icon"></VIcon><VTooltip activator="parent" location="top" :text="o.title()"/></VBtn>
        </div>
      </VCard>
    </VMenu>
  </ActionButton>
</template>
