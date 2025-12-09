<script lang="ts" setup>
import { mdiContentCopy, mdiContentPaste } from '@mdi/js'
import { useI18n } from '@/lib/locale'
import { PropType } from 'vue'
import { Editor } from '@tiptap/vue-3'
import { type Notify } from '../../types'

const props = defineProps({
  t: {
    type: Function as PropType<(path: string) => string>,
    required: true
  },
  editor: {
    type: Object as PropType<Editor>,
    required: true
  },
  notify: {
    type: Function as PropType<Notify>,
    required: true
  }
})

const editor: Editor = props.editor
const notify = props.notify

const { t } = useI18n()

async function copy(e: any) {
  const text = JSON.stringify(editor.getJSON())
  await e.view.navigator.clipboard.writeText(text)
  notify(t('copiedToClipboard'))
}

async function paste(e: any) {
  e.view.navigator.clipboard.readText().then((text: string) => {
    if (!text) {
      notify(t('pastedIsEmpty'), 'error')
    } else if (text[0] == '{' && text[text.length - 1] == '}') {
      try {
        const data = JSON.parse(text)
        editor.commands.setContent(data)
        notify(t('pastedFromClipboard'))
      } catch (e) {
        notify(t('invalidContent'), 'error')
      }
    } else {
      notify(t('invalidContent'), 'error')
    }
  })
}
</script>

<template>
  <VListItem density="compact">
    <VListItemTitle>{{ t('copyOrPasteContent') }}</VListItemTitle>
    <template v-slot:append>
      <VIcon icon="mdi-menu-right" size="x-small"></VIcon>
    </template>

    <VMenu :open-on-focus="false" activator="parent" open-on-hover submenu>
      <VList density="compact">
        <VListItem
          @click="copy"
          :title="t('copy')"
          :subtitle="t('copyContent')"
          :prepend-icon="`svg:${mdiContentCopy}`"
        />
        <VListItem
          @click="paste"
          :title="t(`paste`)"
          :subtitle="t('pasteContent')"
          :prepend-icon="`svg:${mdiContentPaste}`"
        />
      </VList>
    </VMenu>
  </VListItem>
</template>
