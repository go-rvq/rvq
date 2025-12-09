<script setup lang="ts">
import { nodeViewProps, NodeViewWrapper } from '@tiptap/vue-3'
import { default as Messages, Message } from '@/lib/Messages'
import VXGadLangCodeEditor from '@/lib/GadLangCodeEditor'
import { computed } from 'vue'

const props = defineProps(nodeViewProps)
const messages = (props.node.attrs.messages ?? []) as Message[]

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
</script>

<template>
  <NodeViewWrapper as="div" data-type="scriptValue">
    <VXGadLangCodeEditor
      v-model="computedValue"
      single-line
      no-placeholder
    />
    <Messages class="border-t bg-grey-lighten-5" :items="messages" />
  </NodeViewWrapper>
</template>

<style scoped>

:deep(.cm-gutters) {
  display: none !important;
}

:deep(.cm-editor) {
  display: inline-block;
}

:deep([data-node-view-content]) {
  min-width: 30px;
  height: 25px;
}
</style>
