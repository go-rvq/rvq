<script setup lang="ts">
import { nodeViewProps, NodeViewWrapper } from '@tiptap/vue-3'
import { computed, ref } from 'vue'
import VXGadLangCodeEditor from '@/lib/GadLangCodeEditor'
import { VXBtn } from '@/lib/VXBtn'
import { useI18n } from '@/lib/locale'
import { default as Messages, Message } from '@/lib/Messages'

import insertBeforeIcon from './insert-before.svg'
import insertAfterIcon from './insert-after.svg'


const props = defineProps({
  ...nodeViewProps,

  selected: {
    type: Boolean,
    required: true
  }
})

const getPos = () => props.getPos() as number

const expanded = ref(false)
const computedValue = computed<string>({
  get() {
    return props.node.attrs.value ?? ''
  },
  set(value: string | undefined) {
    props.updateAttributes({
      value
    })
  }
})

const messages = (props.node.attrs.messages ?? []) as Message[]

const { t } = useI18n()

const insertParagraphBefore = () => {
  const { editor, node } = props
  editor.chain()
    .insertContentAt(getPos(), { type: 'paragraph' })
    .run()
}

const insertParagraphAfter = () => {
  const { editor } = props
  editor.chain()
    .insertContentAt(getPos() + 1, { type: 'paragraph' })
    .run()
}

</script>

<template>
  <NodeViewWrapper as="div" data-type="scriptBlock">
    <VXBtn
      size="small"
      :icon="expanded ? 'mdi-chevron-up' : 'mdi-chevron-down'"
      :title="expanded ? 'Ocultar': 'Expandir'" variant="text" @click="expanded = !expanded"
      density="compact"
    />
    {{ t('script') }}
    <VXBtn
      class="mx-3"
      size="small"
      variant="text"
      :icon="insertBeforeIcon"
      density="compact"
      :title="t('insertTextBlockBefore')"
      @click="insertParagraphBefore"
    />
    <VXBtn
      size="small"
      variant="text"
      :icon="insertAfterIcon"
      density="compact"
      :title="t('insertTextBlockAfter')"
      @click="insertParagraphAfter"
    />
    <VXGadLangCodeEditor
      v-if="expanded"
      v-model="computedValue"
    />
    <Messages :items="messages" class="border-t bg-grey-lighten-5" />
  </NodeViewWrapper>
</template>

<style scoped>
</style>
