<script setup lang="ts">
// The dialog of the link button of vx-tiptap-editor (vuetify-pro-tiptap's
// Link: its dialogComponent, the button its activator): rvq mdkit's link
// form — an address typed, or a page of the site or a record of the admin
// found as one types (a record: its address in the admin); a new tab or
// not; the link taken out.
import type {Editor} from '@tiptap/vue-3'
import {ref} from 'vue'
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

function apply(link: { href: string, target: string | null }) {
  props.editor.chain().focus().extendMarkRange('link').setLink({href: link.href, target: link.target}).run()
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
        <LinkForm v-if="dialog" :href="value || ''" :target="value ? (target || '') : '_blank'" :linked="!!value" :words="words"
                  @apply="apply" @remove="remove" @cancel="close"/>
      </VCardText>
    </VCard>
  </VDialog>
</template>
