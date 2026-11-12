<script setup lang="ts">
// The dialog of the link button of vx-tiptap-editor (vuetify-pro-tiptap's
// Link: its dialogComponent, the button its activator): rvq mdkit's link
// form — an address typed, or a page or a post of the site found as one
// types; its text (empty: the title of the one chosen, else the address); a
// new tab or not; the link taken out.
import type {Editor} from '@tiptap/vue-3'
import {ref, watch} from 'vue'
import {LinkForm} from '@rvq/mdkit'
import {useI18n} from '@/lib/locale'
import {mdkitWords} from '@/lib/locale/mdkit'

interface Props {
  value?: string
  target?: '_self' | '_blank'
  editor: Editor
  destroy?: () => void
}

const props = withDefaults(defineProps<Props>(), {
  value: undefined,
  target: '_blank',
  destroy: undefined
})

const {t} = useI18n()
const words = mdkitWords(t).link
const dialog = ref<boolean>(false)
// selected is the text of the link there (its whole mark), else the one
// chosen: the text of the link, as the form opens
const selected = ref<string>('')

watch(dialog, (open) => {
  if (!open) return
  props.editor.chain().extendMarkRange('link').run()
  const {from, to} = props.editor.state.selection
  selected.value = props.editor.state.doc.textBetween(from, to, ' ')
})

// apply links the text: the one there, as it is (its marks kept); another
// — typed, or the title of the page or post chosen —, in its place
function apply(link: { href: string, target: string | null, text?: string }) {
  const attrs = {href: link.href, target: link.target}
  const chain = props.editor.chain().focus().extendMarkRange('link')
  if (link.text && link.text !== selected.value) {
    chain.insertContent({type: 'text', text: link.text, marks: [{type: 'link', attrs}]}).run()
  } else {
    chain.setLink(attrs).run()
  }
  close()
}

function remove() {
  props.editor.chain().focus().extendMarkRange('link').unsetLink().run()
  close()
}

function close() {
  dialog.value = false
  setTimeout(() => props.destroy?.(), 300)
}
</script>

<template>
  <VDialog v-model="dialog" max-width="560" activator="parent" scrollable @click:outside="close" data-tiptap-link-dialog>
    <VCard>
      <VCardTitle class="d-flex align-center ga-2">
        <VIcon icon="mdi-link-variant"/>
        {{ words.title }}
      </VCardTitle>
      <VCardText>
        <LinkForm v-if="dialog" :href="value || ''" :target="value ? (target || '') : '_blank'" :text="selected" :linked="!!value" :words="words"
                  @apply="apply" @remove="remove" @cancel="close"/>
      </VCardText>
    </VCard>
  </VDialog>
</template>
