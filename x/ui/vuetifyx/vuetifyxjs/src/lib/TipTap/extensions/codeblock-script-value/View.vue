<script setup lang="ts">
import { nodeViewProps, NodeViewWrapper } from '@tiptap/vue-3'
import { computed } from 'vue'
import VXGadLangCodeEditor from '@/lib/GadLangCodeEditor'
import { Message, default as Messages } from '@/lib/Messages'

const props = defineProps({
  ...nodeViewProps,

  selected: {
    type: Boolean,
    required: true
  }
})

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

</script>

<template>
  <NodeViewWrapper as="div" style="display: inline-block">
    <div data-type="codeBlockScriptValue">
      <VXGadLangCodeEditor
        v-model="computedValue"
        single-line
        no-placeholder
      />
      <Messages class="border-t bg-grey-lighten-5" :items="messages" />
    </div>
  </NodeViewWrapper>
</template>

<style scoped>
:deep(.cm-gutters) {
  display: none !important;
}
:deep(.cm-editor) {
  display: inline-block;
}
</style>
