<script setup lang="ts">
// The button of the plugin tables on the bar of vx-tiptap-editor: rvq
// mdkit's table builder — a sheet: rows, columns, merged cells, the header,
// alignment, widths, the text of the cells formatted — on the table under the
// cursor (its HTML read: fromHTML), or a new one; done, it put in its place
// as the editor's table (toRichHTML, the cells' Markdown made HTML).
import type {Editor} from '@tiptap/vue-3'
import {mdiTableEdit} from '@mdi/js'
import {ActionButton} from 'vuetify-pro-tiptap'
import {ref} from 'vue'
import {getHTMLFromFragment} from '@tiptap/core'
import {Fragment} from '@tiptap/pm/model'
import {TableBuilder, configure, inlineHTML, tables, textStyle} from '@rvq/mdkit'
import {useI18n} from '@/lib/locale'
import {mdkitWords} from '@/lib/locale/mdkit'

const props = defineProps<{ editor: Editor, disabled?: boolean, base?: number }>()
const {t} = useI18n()
const words = mdkitWords(t)
configure({words})
textStyle()
const builder = ref<any>(null)

// tableAt is the table the cursor is in: its node and place; null, none
function tableAt() {
  const $from = props.editor.state.selection.$from
  for (let d = $from.depth; d > 0; d--) {
    const node = $from.node(d)
    if (node.type.name === 'table') return {node, pos: $from.before(d)}
  }
  return null
}

function open() {
  const found = tableAt()
  const table = found ? tables.fromHTML(getHTMLFromFragment(Fragment.from(found.node), props.editor.schema)) : null
  // done: the table put in its place (the one found, still there), or at
  // the cursor
  builder.value?.open(table, () => {
    const html = tables.toRichHTML(builder.value.t, inlineHTML, props.base || 900)
    const now = found && tableAt()
    const chain = props.editor.chain().focus()
    if (now) chain.insertContentAt({from: now.pos, to: now.pos + now.node.nodeSize}, html).run()
    else chain.insertContent(html).run()
  })
}
</script>

<template>
  <ActionButton :rawIcon="`svg:${mdiTableEdit}`" :editor="editor" :tooltip="words.markdown.table" :disabled="disabled"
                :action="open" data-tiptap-plugin="tables"/>
  <TableBuilder ref="builder" :words="words.table" :bar-words="words.rich"/>
</template>
